.PHONY: help env-sync env-check secret-push build run run-window run-window-test run-test test clean

ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SCRIPTS := $(ROOT)scripts

help:
	@echo "pray-notify — perintah tersedia:"
	@echo ""
	@echo "  make env-sync              Salin .env.example & targets.example.json"
	@echo "  make env-check             Cek file env lokal"
	@echo "  make secret-push           Push targets.json → GitHub secret"
	@echo "  make build                 Build binary prayer-notifier"
	@echo "  make run-window PRAYER=ashar       Window mode (sleep + kirim)"
	@echo "  make run-window-test PRAYER=ashar  Window mode tanpa sleep"
	@echo "  make run-test              Legacy test mode (-test)"
	@echo "  make test                  go test ./..."
	@echo "  make clean                 Hapus binary"

env-sync:
	@chmod +x $(SCRIPTS)/sync-env.sh
	@$(SCRIPTS)/sync-env.sh

env-check:
	@test -f .env || (echo "❌ .env belum ada. Jalankan: make env-sync" && exit 1)
	@test -f deploy/targets.json || (echo "❌ deploy/targets.json belum ada. Jalankan: make env-sync" && exit 1)
	@! grep -q 'REPLACE_ME' deploy/targets.json || (echo "❌ deploy/targets.json masih berisi REPLACE_ME — isi webhook Discord asli" && exit 1)
	@echo "✅ env OK"

secret-push: env-check
	@chmod +x $(SCRIPTS)/push-secret.sh
	@$(SCRIPTS)/push-secret.sh

build:
	go build -o prayer-notifier main.go

run-window: env-check
	@test -n "$(PRAYER)" || (echo "❌ Set PRAYER=subuh|dhuha|dzuhur|ashar|maghrib|isya" && exit 1)
	@set -a && . ./.env && set +a && go run main.go -prayer "$(PRAYER)"

run-window-test: env-check
	@test -n "$(PRAYER)" || (echo "❌ Set PRAYER=subuh|dhuha|dzuhur|ashar|maghrib|isya" && exit 1)
	@set -a && . ./.env && set +a && go run main.go -prayer "$(PRAYER)" -test

run-test: env-check
	@set -a && . ./.env && set +a && go run main.go -test

test:
	go test ./...

clean:
	rm -f prayer-notifier pray-notify
