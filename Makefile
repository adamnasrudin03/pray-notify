.PHONY: help env-sync env-check secret-push build build-serve run serve run-test test clean

ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SCRIPTS := $(ROOT)scripts

help:
	@echo "pray-notify — perintah tersedia:"
	@echo ""
	@echo "  make env-sync      Salin .env.example & targets.example.json jika belum ada"
	@echo "  make env-check     Cek file env lokal sudah ada"
	@echo "  make secret-push   Push deploy/targets.json → GitHub secret (manual test)"
	@echo "  make build         Build binary prayer-notifier (CLI)"
	@echo "  make build-serve   Build binary HTTP server (production)"
	@echo "  make serve         Jalankan HTTP server untuk cron-job.org"
	@echo "  make run           Jalankan one-shot (butuh .env)"
	@echo "  make run-test      Jalankan mode uji coba (-test)"
	@echo "  make test          go test ./..."
	@echo "  make clean         Hapus binary build"

env-sync:
	@chmod +x $(SCRIPTS)/sync-env.sh
	@$(SCRIPTS)/sync-env.sh

env-check:
	@test -f .env || (echo "❌ .env belum ada. Jalankan: make env-sync" && exit 1)
	@test -f deploy/targets.json || (echo "❌ deploy/targets.json belum ada. Jalankan: make env-sync" && exit 1)
	@echo "✅ env OK"

secret-push: env-check
	@chmod +x $(SCRIPTS)/push-secret.sh
	@$(SCRIPTS)/push-secret.sh

build:
	go build -o prayer-notifier main.go

build-serve:
	go build -o serve cmd/serve/main.go

serve: env-check build-serve
	@set -a && . ./.env && set +a && ./serve

run: env-check
	@set -a && . ./.env && set +a && go run main.go

run-test: env-check
	@set -a && . ./.env && set +a && go run main.go -test

test:
	go test ./...

clean:
	rm -f prayer-notifier pray-notify serve
