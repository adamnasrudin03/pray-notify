#!/usr/bin/env bash
# Push deploy/targets.json ke GitHub Actions secret WEBHOOK_TARGETS.
# Usage: ./scripts/push-secret.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGETS_FILE="${WEBHOOK_TARGETS_FILE:-$ROOT/deploy/targets.json}"

if [[ -f "$ROOT/.env" ]]; then
	# shellcheck disable=SC1090
	set -a
	source "$ROOT/.env"
	set +a
fi

if [[ ! -f "$TARGETS_FILE" ]]; then
	echo "❌ File tidak ditemukan: $TARGETS_FILE" >&2
	echo "   Jalankan: make env-sync" >&2
	exit 1
fi

if ! command -v gh >/dev/null 2>&1; then
	echo "❌ gh CLI tidak terinstall. Install: https://cli.github.com/" >&2
	exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
	if ! python3 -c "import json,sys; json.load(open(sys.argv[1]))" "$TARGETS_FILE" 2>/dev/null; then
		echo "❌ JSON tidak valid: $TARGETS_FILE" >&2
		exit 1
	fi
elif ! jq empty "$TARGETS_FILE" >/dev/null 2>&1; then
	echo "❌ JSON tidak valid: $TARGETS_FILE" >&2
	exit 1
fi

REPO="${GITHUB_REPO:-}"
if [[ -z "$REPO" ]]; then
	REPO="$(gh repo view --json nameWithOwner -q .nameWithOwner 2>/dev/null || true)"
fi

if [[ -z "$REPO" ]]; then
	echo "❌ GITHUB_REPO tidak diset di .env dan tidak bisa deteksi dari git remote" >&2
	exit 1
fi

echo "==> Push WEBHOOK_TARGETS ke $REPO"
gh secret set WEBHOOK_TARGETS --repo "$REPO" <"$TARGETS_FILE"
echo "✅ Secret WEBHOOK_TARGETS berhasil diupdate"
