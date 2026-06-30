# Setup cron-job.org untuk pray-notify

Production **tidak lagi** memicu GitHub Actions (`repository_dispatch`) karena ~990 run/hari melebihi kuota **GitHub Actions Free** (~2.000 menit/bulan).

cron-job.org sekarang memanggil **HTTP endpoint** `/notify` dari service yang menjalankan binary (Render, Railway, VPS, dll.).

## Arsitektur

```
cron-job.org (tiap menit, 03:30–19:59 WIB)
    │
    ▼ POST https://your-app.onrender.com/notify
cmd/serve (Go HTTP)
    │
    ▼ scheduler.Run()
Discord webhook
```

GitHub Actions (`schedule.yml`) hanya untuk **test manual** via workflow_dispatch.

## Deploy service (contoh Render)

1. Buat **Web Service** dari repo ini
2. **Build command:** `go build -o serve cmd/serve/main.go`
3. **Start command:** `./serve`
4. **Environment variables:**
   - `WEBHOOK_TARGETS` — JSON array (sama seperti GitHub secret)
   - `CRON_SECRET` — token acak untuk auth cron-job.org
   - `PORT` — biasanya diset otomatis oleh platform

## Prasyarat cron-job.org

- URL service sudah live (`/health` → `ok`)
- `CRON_SECRET` sudah diset di service

## Langkah cron-job.org

### 1. Buat cronjob

1. Buka [console.cron-job.org](https://console.cron-job.org)
2. **Cronjobs** → **Create cronjob**

### 2. Tab Common

| Field | Nilai |
|-------|-------|
| Title | `pray-notify` |
| URL | `https://your-app.onrender.com/notify` |
| Enable job | ✅ |

### 3. Tab Schedule

- **Timezone:** `Asia/Jakarta`
- Buat **3 job** dengan cron:

```
30-59 3 * * *
```

```
* 4-18 * * *
```

```
* 19 * * *
```

### 4. Tab Request

| Field | Nilai |
|-------|-------|
| Request method | `POST` |

**Custom request headers:**

```
Authorization: Bearer your-cron-secret
```

### 5. Test

**Run now** → HTTP **204**. Cek Discord channel.

Test manual:

```bash
curl -X POST \
  -H "Authorization: Bearer your-cron-secret" \
  https://your-app.onrender.com/notify
```

## Nonaktifkan job lama (penting)

Hapus atau disable cron-job yang masih memanggil:

```
https://api.github.com/repos/adamnasrudin03/pray-notify/dispatches
```

Itu memicu GitHub Actions dan menghabiskan kuota.

## Window & anti-spam di code

| Guard | Perilaku |
|-------|----------|
| Window operasi | 03:30–19:59 (timezone target) |
| Pending check | Skip jika tidak ada notifikasi aktif |
| Delivery window | 90 detik per pengingat |

## Troubleshooting

| Gejala | Solusi |
|--------|--------|
| HTTP 401 | `CRON_SECRET` tidak cocok |
| HTTP 204 tapi tidak ada Discord | Cek `WEBHOOK_TARGETS` di service |
| Service sleep (Render free) | Request pertama lambat; pertimbangkan paid tier atau platform lain |
| Masih kena limit GitHub | Pastikan job dispatch ke GitHub sudah dihapus |
