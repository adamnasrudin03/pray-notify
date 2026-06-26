# pray-notify

Bot pengingat waktu shalat yang mengirim notifikasi ke Discord via webhook.

**Production:** [cron-job.org](https://cron-job.org) memicu GitHub Actions lewat API `repository_dispatch` — tanpa server sendiri.

## Arsitektur

```
cron-job.org (setiap menit, 03:30–19:59 WIB)
    │
    ▼ POST /repos/adamnasrudin03/pray-notify/dispatches
GitHub Actions workflow
    │
    ▼ prayer-notifier
Discord webhook
```

## Jadwal notifikasi

Per shalat fardhu: **5 pengingat** pada offset **-10, -5, 0, +5, +10 menit** dari waktu adzan.  
Dhuha: **1×** saat waktu masuk. Total ~**26 pesan/hari** per channel.

## Setup

### 1. GitHub repo

1. Set secret **`WEBHOOK_TARGETS`** di repo → Settings → Secrets and variables → Actions
2. Pastikan workflow `.github/workflows/schedule.yml` sudah di-push (listen event `pray-notify`)

### 2. GitHub Personal Access Token

Buat PAT untuk cron-job.org:

| Tipe | Permission |
|------|------------|
| Classic PAT | scope `repo` |
| Fine-grained | Contents: **Read and write** |

### 3. Test manual (opsional)

```bash
curl -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ghp_xxxxxxxx" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -H "Content-Type: application/json" \
  https://api.github.com/repos/adamnasrudin03/pray-notify/dispatches \
  -d '{"event_type":"pray-notify"}'
```

Sukses = HTTP **204**. Cek tab **Actions** di repo — workflow harus jalan.

### 4. Buat job di cron-job.org

Login ke [cron-job.org](https://console.cron-job.org) → **Create cronjob**.

Isi per job:

| Field | Nilai |
|-------|-------|
| **Title** | `pray-notify` (atau bebas) |
| **URL** | `https://api.github.com/repos/adamnasrudin03/pray-notify/dispatches` |
| **Schedule** | lihat tabel di bawah |
| **Timezone** | `Asia/Jakarta` |
| **Request method** | `POST` |
| **Request body** | `{"event_type":"pray-notify"}` |

**Headers** (tambahkan satu per satu):

| Header | Value |
|--------|-------|
| `Accept` | `application/vnd.github+json` |
| `Authorization` | `Bearer ghp_xxxxxxxx` |
| `X-GitHub-Api-Version` | `2022-11-28` |
| `Content-Type` | `application/json` |

Buat **3 cronjob** dengan jadwal berbeda (window shalat 03:30–19:59 WIB):

| Job | Cron expression | Jam (WIB) |
|-----|-----------------|-----------|
| 1 | `30-59 3 * * *` | 03:30 – 03:59 |
| 2 | `* 4-18 * * *` | 04:00 – 18:59 |
| 3 | `* 19 * * *` | 19:00 – 19:59 (cover Isya +10 menit) |

> Detail langkah ada di [`deploy/cron-job.org.md`](deploy/cron-job.org.md)

## Konfigurasi webhook (`WEBHOOK_TARGETS`)

Secret GitHub berisi JSON array:

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

## Development lokal

### Setup env (sekali)

```bash
make env-sync        # salin .env.example → .env, targets.example.json → targets.json
```

Edit secret lokal:
- `.env` — variabel environment
- `deploy/targets.json` — webhook Discord (jangan di-commit)

Tambah key baru di `.env.example`? Jalankan `make env-sync` lagi — key baru otomatis ditambahkan ke `.env` tanpa menimpa nilai lama.

### Push secret ke GitHub

```bash
make secret-push     # deploy/targets.json → GitHub secret WEBHOOK_TARGETS
```

Butuh [gh CLI](https://cli.github.com/) sudah login (`gh auth login`).

### Menjalankan

```bash
make run             # jalankan notifier
make run-test        # mode uji coba (-test)
make test            # unit test
make build           # build binary
```

Manual tanpa Makefile:

```bash
export WEBHOOK_TARGETS_FILE=./deploy/targets.json
go run main.go
```

## Kenapa cron-job.org + dispatch, bukan GH schedule?

GitHub Actions `on: schedule` sering telat karena antrian runner. `repository_dispatch` dipicu cron-job.org tepat waktu, workflow biasanya mulai dalam hitungan detik.

### Guard di code (anti-spam)

Selaras dengan jadwal cron-job.org, app punya 3 lapis proteksi:

| Guard | Perilaku |
|-------|----------|
| **Window operasi** | Skip run di luar **03:30–19:59** (timezone per target) |
| **Pending check** | Skip target jika tidak ada notifikasi dalam window aktif (jam 10:00 tidak proses Discord) |
| **Delivery window 90 detik** | Setiap notifikasi max 1x per siklus cron (bukan 5 menit seperti dulu) |

## Perilaku mention

- `@here` hanya di **content** webhook, untuk Dzuhur/Ashar/Maghrib
- Subuh, Isya, Dhuha: tanpa `@here`
