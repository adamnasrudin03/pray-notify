# pray-notify

Bot pengingat waktu shalat yang mengirim notifikasi ke Discord via webhook.

**Production:** [cron-job.org](https://cron-job.org) memanggil HTTP endpoint `/notify` — **bukan GitHub Actions** (hemat kuota Free).

## Arsitektur

```
cron-job.org (setiap menit, 03:30–19:59 WIB)
    │
    ▼ POST https://your-app.onrender.com/notify
cmd/serve (Go HTTP)
    │
    ▼ prayer-notifier logic
Discord webhook
```

> GitHub Actions (`schedule.yml`) hanya untuk **test manual**. Jangan pakai `repository_dispatch` dari cron — ~990 run/hari melebihi limit GitHub Free.

## Jadwal notifikasi

Per shalat fardhu: **5 pengingat** pada offset **-10, -5, 0, +5, +10 menit** dari waktu adzan.  
Dhuha: **1×** saat waktu masuk. Total ~**26 pesan/hari** per channel.

## Setup production

1. Deploy `cmd/serve` ke Render/Railway/VPS (lihat [`deploy/cron-job.org.md`](deploy/cron-job.org.md))
2. Set env `WEBHOOK_TARGETS` + `CRON_SECRET` di platform
3. Pasang cron-job.org → `POST /notify` dengan header `Authorization: Bearer <CRON_SECRET>`
4. **Nonaktifkan** cron lama yang hit `api.github.com/.../dispatches`

## Setup

### GitHub (opsional — test manual saja)

1. Secret **`WEBHOOK_TARGETS`** untuk workflow_dispatch manual
2. Actions → **Prayer Notification (Manual)** → Run workflow

## Development lokal

### Setup env (sekali)

```bash
make env-sync        # salin .env.example → .env, targets.example.json → targets.json
```

Edit secret lokal:
- `.env` — variabel environment
- `deploy/targets.json` — webhook Discord (jangan di-commit)

Tambah key baru di `.env.example`? Jalankan `make env-sync` lagi — key baru otomatis ditambahkan ke `.env` tanpa menimpa nilai lama.

### Push secret ke GitHub (manual test)

```bash
make secret-push     # deploy/targets.json → GitHub secret WEBHOOK_TARGETS
```

Butuh [gh CLI](https://cli.github.com/) sudah login (`gh auth login`).

### Menjalankan

```bash
make serve           # HTTP server untuk production/cron-job.org
make run             # one-shot (seperti cron)
make run-test        # mode uji coba (-test)
make test            # unit test
make build           # build binary CLI
```

Manual tanpa Makefile:

```bash
export WEBHOOK_TARGETS_FILE=./deploy/targets.json
go run main.go
```

## Konfigurasi webhook (`WEBHOOK_TARGETS`)

JSON array di env service atau `deploy/targets.json`:

```json
[
  {
    "Name": "Divisi IT",
    "URL": "https://discord.com/api/webhooks/...",
    "Lat": -6.1895,
    "Lon": 106.7649,
    "Tz": "Asia/Jakarta",
    "City": "Kebon Jeruk, Jakarta Barat",
    "SkipWeekends": true,
    "MentionMode": "here"
  }
]
```

### Field konfigurasi

| Field | Default | Keterangan |
|-------|---------|------------|
| `Name` | — | Nama channel (untuk log) |
| `URL` | — | Discord webhook URL |
| `Lat`, `Lon` | Jakarta default | Koordinat perhitungan jadwal |
| `Tz` | `Asia/Jakarta` | Timezone |
| `City` | Jakarta default | Ditampilkan di footer embed |
| `SkipWeekends` | `true` untuk "Divisi IT" | Lewati notifikasi Sabtu/Minggu |
| `MentionMode` | `here` | `here`, `none`, atau `custom` |
| `Mentions` | — | Digunakan jika `MentionMode` = `custom` |
| `EnableTTS` | `false` | Text-to-speech Discord |
| `Username` | `Reminder Adzan {shalat}` | Override nama bot |
| `UsernameBeforeAdzan` | personal untuk non-Divisi IT | Override nama bot sebelum adzan |
| `IsDisable` | `false` | Nonaktifkan target |

## Kenapa tidak pakai GitHub Actions untuk cron?

| | GitHub Actions + dispatch | cron-job.org + `/notify` |
|--|---------------------------|--------------------------|
| Run/hari | ~990 | ~990 (tapi di server sendiri) |
| Kuota GitHub | Habis (~2000 min/bulan) | **0 menit** |
| Latensi | Antrian runner | Langsung ke app |

### Guard di code (anti-spam)

| Guard | Perilaku |
|-------|----------|
| **Window operasi** | Skip run di luar **03:30–19:59** (timezone per target) |
| **Pending check** | Skip target jika tidak ada notifikasi dalam window aktif |
| **Delivery window 90 detik** | Setiap notifikasi max 1x per siklus cron |

## Perilaku mention

- `@here` hanya di **content** webhook, untuk Dzuhur/Ashar/Maghrib
- Subuh, Isya, Dhuha: tanpa `@here`
