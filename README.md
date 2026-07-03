# pray-notify

Bot pengingat waktu shalat yang mengirim notifikasi ke Discord via webhook.

**Production (Opsi A):** [cron-job.org](https://cron-job.org) memicu GitHub Actions **6×/hari**. Workflow **sleep** sampai tiap offset, lalu kirim — tanpa server lain.

> Repo disarankan **public** agar kuota GitHub Actions Free cukup (~120 menit/hari).

## Arsitektur

```
cron-job.org (6×/hari)
    │
    ▼ repository_dispatch  {"prayer":"ashar"}
GitHub Actions (~20 menit: sleep −10…+10 menit)
    │
    ▼ prayer-notifier -prayer ashar
Discord webhook
```

## Jadwal notifikasi

Per shalat fardhu: **5 pengingat** (−10, −5, 0, +5, +10 menit). Dhuha: **1×**.  
Total ~**26 pesan/hari** per channel, **6 workflow/hari**.

## Setup

1. `make env-sync` → edit **`deploy/targets.json`** (isi URL webhook Discord, bukan `REPLACE_ME`)
2. `make secret-push` → sync ke GitHub secret **`WEBHOOK_TARGETS`**
3. Pasang **6 cronjob** di cron-job.org → [`deploy/cron-job.org.md`](deploy/cron-job.org.md)
4. **Nonaktifkan** cron hit tiap menit / dispatch lama
5. Repo **public** (disarankan)

## Development lokal

```bash
make env-sync
make secret-push     # optional: sync targets.json → GitHub secret
make run-window-test PRAYER=ashar    # test window tanpa sleep
make run-window PRAYER=ashar         # test dengan sleep
make test
```

## Konfigurasi webhook (`WEBHOOK_TARGETS`)

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

| Field | Default | Keterangan |
|-------|---------|------------|
| `SkipWeekends` | `true` untuk "Divisi IT" | Lewati Sabtu/Minggu |
| `MentionMode` | `here` | `here`, `none`, `custom` |

## Kuota GitHub Actions

| Skema | Run/hari | Menit/hari (kira-kira) |
|-------|----------|------------------------|
| Cron tiap menit (lama) | ~990 | habis limit |
| **Opsi A (6 window)** | **6** | **~120** |

## Perilaku mention

- `@here` untuk Dzuhur, Ashar, Maghrib
- Subuh, Isya, Dhuha: tanpa `@here`
