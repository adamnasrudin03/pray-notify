package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adamnasrudin03/pray-notify/scheduler"
)

var (
	testingMode     bool
	prayerWindow    string
	serverMode      bool
	port            string
)

func main() {
	flag.BoolVar(&testingMode, "test", false, "Aktifkan mode uji coba.")
	flag.StringVar(&prayerWindow, "prayer", "", "Mode window: subuh, dhuha, dzuhur, ashar, maghrib, isya")
	flag.BoolVar(&serverMode, "server", false, "Jalankan sebagai HTTP server (untuk Render/health check)")
	flag.StringVar(&port, "port", "8080", "Port untuk HTTP server")
	flag.Parse()

	// Server mode untuk Render (health check + manual trigger)
	if serverMode {
		runServer(port)
		return
	}

	// CLI mode (existing behavior)
	runCLI()
}

func runCLI() {
	scheduler.RunTesting(testingMode && prayerWindow == "")

	if prayerWindow != "" {
		if err := scheduler.RunPrayerWindow(prayerWindow, testingMode); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		return
	}

	scheduler.Run(testingMode)
}

func runServer(port string) {
	mux := http.NewServeMux()

	// Health check endpoint (Render pings ini)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Manual trigger endpoint (dipanggil GitHub Actions cron)
	mux.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse query params untuk testing
		testMode := r.URL.Query().Get("test") == "true"
		prayer := r.URL.Query().Get("prayer")

		log.Printf("🚀 Trigger received: test=%v, prayer=%s", testMode, prayer)

		go func() {
			if prayer != "" {
				if err := scheduler.RunPrayerWindow(prayer, testMode); err != nil {
					log.Printf("❌ Error running prayer window %s: %v", prayer, err)
				}
			} else {
				scheduler.Run(testMode)
			}
			log.Println("✅ Scheduler completed")
		}()

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Scheduler triggered"))
	})

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Pray Notify API
Endpoints:
  GET  /healthz     - Health check
  POST /run         - Trigger scheduler (query: ?test=true&prayer=subuh)
`))
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🌐 Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("✅ Server exited")
}
