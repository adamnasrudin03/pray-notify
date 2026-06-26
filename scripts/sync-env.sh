#!/usr/bin/env bash
# Sync file example → lokal tanpa menimpa secret yang sudah ada.
# Usage: ./scripts/sync-env.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

log_create() { echo "✅ CREATE $1"; }
log_skip()   { echo "⏭️  SKIP   $1 (sudah ada)"; }
log_add()    { echo "➕ ADD    $1"; }

copy_if_missing() {
	local example="$1"
	local target="$2"

	if [[ ! -f "$example" ]]; then
		echo "❌ Example tidak ditemukan: $example" >&2
		exit 1
	fi

	if [[ -f "$target" ]]; then
		log_skip "$target"
		return 0
	fi

	mkdir -p "$(dirname "$target")"
	cp "$example" "$target"
	log_create "$target"
}

merge_env_keys() {
	local example="$1"
	local target="$2"

	[[ -f "$target" ]] || return 0

	while IFS= read -r line || [[ -n "$line" ]]; do
		[[ -z "$line" || "$line" =~ ^[[:space:]]*# ]] && continue
		[[ "$line" != *"="* ]] && continue

		local key="${line%%=*}"
		key="${key#"${key%%[![:space:]]*}"}"
		key="${key%"${key##*[![:space:]]}"}"

		if ! grep -q "^[[:space:]]*${key}=" "$target"; then
			printf '\n%s\n' "$line" >>"$target"
			log_add "${key} → ${target}"
		fi
	done <"$example"
}

echo "==> Sync environment files"
copy_if_missing "$ROOT/.env.example" "$ROOT/.env"
merge_env_keys "$ROOT/.env.example" "$ROOT/.env"

copy_if_missing "$ROOT/deploy/targets.example.json" "$ROOT/deploy/targets.json"

echo ""
echo "Selesai. Edit secret lokal:"
echo "  - $ROOT/.env"
echo "  - $ROOT/deploy/targets.json"
