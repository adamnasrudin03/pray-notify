# AGENTS.md — pray-notify

> **Panduan untuk AI Agent / Developer** yang bekerja pada project ini.
> Jangan commit/push tanpa persetujuan eksplisit dari user.

---

## 📦 Project Overview

**pray-notify** — Bot pengingat waktu shalat ke Discord via webhook.

- **Language**: Go 1.25+
- **Scheduler**: `github.com/hablullah/go-prayer` (Kemenag, Shafii)
- **Deployment**: Dual-mode (GitHub Actions **atau** Render)
- **Config**: JSON via `WEBHOOK_TARGETS` env var

---

## 🏗 Arsitektur Dual-Mode

```
┌─────────────────────────────────────────────────────────────┐
│                    cron-job.org (6×/hari)                    │
└───────────────────────────┬─────────────────────────────────┘
                             │
              ┌──────────────┴──────────────┐
              ▼                             ▼
     ┌─────────────────┐            ┌─────────────────┐
     │ schedule.yml    │            │  Render HTTP    │
     │ (GitHub Actions)│            │  Server /run    │
     │ - build lokal   │            │ - curl POST     │
     │ - run scheduler │            │ - Render exec   │
     │ - pakai GH      │            │ - pakai Render  │
     │   WEBHOOK_TARGETS            │   WEBHOOK_TARGETS
     └────────┬────────┘            └────────┬────────┘
              │                              │
              └──────────────────┬───────────┘
                                 ▼
                          Discord Webhooks
```

### Workflow Files

| File | Trigger | Runs On | Env Source |
|------|---------|---------|------------|
| `.github/workflows/schedule.yml` | `pray-notify` | GitHub Actions (ubuntu-latest) | GitHub Secret `WEBHOOK_TARGETS` |

**Render HTTP Server** dijalankan terpisah, dipanggil langsung oleh cron-job.org ke endpoint `/run`.

---

## 🚀 Quick Start (Development)

```bash
# 1. Sync env files
make env-sync

# 2. Edit deploy/targets.json (isi webhook Discord asli)
# 3. Push ke GitHub secret (optional)
make secret-push

# 4. Test lokal
make run-window-test PRAYER=ashar   # tanpa sleep
make run-window PRAYER=ashar        # dengan sleep

# 5. Test server mode (Render-like)
./pray-notify -server -port=8080
curl http://localhost:8080/healthz
curl -X POST "http://localhost:8080/run?prayer=subuh&test=true"
```

---

## 📁 Project Structure

```
pray-notify/
├── main.go                 # Entry point (CLI + HTTP server)
├── Dockerfile              # Multi-stage build (~15MB)
├── render.yaml             # Render free tier config
├── Makefile                # Build, test, deploy helpers
├── go.mod / go.sum
├── .env.example            # Template env
├── deploy/
│   ├── targets.example.json    # Template webhook config
│   └── cron-job.org.md         # Setup cron-job.org guide
├── .github/workflows/
│   └── schedule.yml            # GitHub Actions lokal (pray-notify)
├── config/                   # Config loading (WEBHOOK_TARGETS)
├── scheduler/                # Core scheduling logic
├── notification/             # Discord webhook sender
├── models/                   # Data structures
├── helpers/                  # Utilities
├── quotes/                   # Quran/hadits quotes
├── reminders/                # Offset reminders (-10..+10 min)
└── scripts/                  # sync-env.sh, push-secret.sh
```

---

## ⚙️ Configuration (WEBHOOK_TARGETS)

Format JSON array (set via env var `WEBHOOK_TARGETS`):

```json
[
  {
    "Name": "Divisi IT",
    "URL": "https://discord.com/api/webhooks/ID/TOKEN",
    "Lat": -6.1895,
    "Lon": 106.7649,
    "Tz": "Asia/Jakarta",
    "City": "Kebon Jeruk, Jakarta Barat",
    "SkipWeekends": true,
    "MentionMode": "here",
    "EnableTTS": false,
    "Username": "",
    "UsernameBeforeAdzan": "",
    "IsDisable": false
  }
]
```

| Field | Required | Default | Notes |
|-------|----------|---------|-------|
| `Name` | ✅ | — | Label untuk logging |
| `URL` | ✅ | — | Discord webhook URL |
| `Lat`/`Lon` | ✅ | Jakarta | Koordinat lokasi |
| `Tz` | ✅ | Asia/Jakarta | IANA timezone |
| `City` | ❌ | Jakarta | Nama tampilan footer |
| `SkipWeekends` | ❌ | `true` (Divisi IT) | Skip Sabtu/Minggu |
| `MentionMode` | ❌ | `here` | `here`, `none`, `custom` |
| `EnableTTS` | ❌ | `false` | Discord text-to-speech |

---

## 🔐 Secrets & Environment Variables

### GitHub Secrets (Settings → Secrets → Actions)

| Secret | Used By | Description |
|--------|---------|-------------|
| `WEBHOOK_TARGETS` | `schedule.yml` | JSON array webhook (untuk GH Actions lokal) |

### Render Dashboard (Environment Variables)

| Env Var | Description |
|---------|-------------|
| `WEBHOOK_TARGETS` | JSON array webhook (untuk Render execution) |
| `PORT` | Auto-set to `8080` |

---

## 🐳 Docker & Render Deploy

```bash
# Build image
make docker-build

# Run locally (mount targets.json)
make docker-run

# Deploy ke Render:
# 1. Push ke GitHub (auto-deploy via render.yaml)
# 2. Set WEBHOOK_TARGETS di Render Dashboard
```

**render.yaml** sudah konfigurasi:
- Free tier (750 jam/bulan)
- Region: Singapore
- Health check: `/healthz`
- Start command: `./pray-notify -server -port=8080`

---

## ⏰ cron-job.org Setup (6 jobs/hari)

Buat 6 job di cron-job.org, tiap 15 menit (03:30–19:59 WIB). Setiap waktu shalat butuh **2 job** (satu untuk GitHub Actions, satu untuk Render):

| Prayer | GitHub Actions (repository_dispatch) | Render (direct POST) |
|--------|--------------------------------------|----------------------|
| subuh  | `pray-notify` + `{"prayer":"subuh"}` | `POST /run?prayer=subuh` |
| dhuha  | `pray-notify` + `{"prayer":"dhuha"}` | `POST /run?prayer=dhuha` |
| dzuhur | `pray-notify` + `{"prayer":"dzuhur"}` | `POST /run?prayer=dzuhur` |
| ashar  | `pray-notify` + `{"prayer":"ashar"}` | `POST /run?prayer=ashar` |
| maghrib| `pray-notify` + `{"prayer":"maghrib"}` | `POST /run?prayer=maghrib` |
| isya   | `pray-notify` + `{"prayer":"isya"}` | `POST /run?prayer=isya` |

**GitHub API endpoint**: `POST https://api.github.com/repos/adamnasrudin03/pray-notify/dispatches`
**Headers**: `Authorization: Bearer <GH_TOKEN>`, `Accept: application/vnd.github.everest-preview+json`

**Render endpoint**: `POST https://YOUR-RENDER.onrender.com/run?prayer=<prayer>&test=false`

Detail: [`deploy/cron-job.org.md`](deploy/cron-job.org.md)

---

## 🧪 Testing

```bash
# Unit tests
make test

# Integration test (manual trigger)
gh workflow run "Prayer Notification" -f prayer=subuh

# Direct Render test
curl -X POST "https://YOUR-RENDER.onrender.com/run?prayer=subuh&test=true"
```

---

## 📝 Development Guidelines

### Code Style
- Go standard formatting (`gofmt`, `go vet`)
- Run `make test` before commit
- Keep `main.go` thin — logic in `scheduler/`, `notification/`

### Adding New Prayer/Reminder
1. Update `reminders/` offsets
2. Update `scheduler/scheduler.go` prayer list
3. Update `models/adzan.go` labels if needed
4. Test dengan `make run-window-test PRAYER=<new>`

### Modifying Workflows
- **schedule.yml** → only for GitHub Actions lokal execution
- Render HTTP server dipanggil langsung oleh cron-job.org (bukan via GitHub Actions)
- Jangan tambahkan workflow baru untuk trigger Render

### Commit Messages
```
feat: <fitur baru>
fix: <bug fix>
chore: <maintenance>
docs: <dokumentasi>
refactor: <refactor tanpa behavior change>
```

---

## ⚠️ Important Rules for Agents

1. **NO AUTO COMMIT/PUSH** — Always ask user before `git commit` / `git push`
2. **Read AGENTS.md first** — Before making changes
3. **Test locally** — Run `make test` and `make run-window-test` before proposing changes
4. **Dual-mode awareness** — Changes to scheduler affect both GH Actions and Render
5. **Secrets handling** — Never log/print `WEBHOOK_TARGETS` or Discord URLs
6. **Timezone** — All times in WIB (Asia/Jakarta), cron in UTC

---

## 🔗 Useful Commands

```bash
make help              # List all commands
make env-check         # Validate local env
make build             # Build binary
make clean             # Remove binaries
make render-deploy     # Show deploy steps
```

---

## 📚 References

- [go-prayer docs](https://github.com/hablullah/go-prayer)
- [Render Web Service docs](https://render.com/docs/web-services)
- [GitHub Actions workflow_dispatch](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#workflow_dispatch)
- [cron-job.org](https://cron-job.org)

---

*Generated for pray-notify project. Update saat arsitektur berubah.*
