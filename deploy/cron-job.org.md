# Setup cron-job.org — Opsi A (GitHub Actions + window loop)

**6 hit/hari** ke GitHub `repository_dispatch`. Workflow **sleep** sampai tiap offset (−10 s.d +10), lalu kirim Discord.

> Repo disarankan **public** agar kuota Actions Free cukup (~120 menit/hari, bukan ~3.600).

## Arsitektur

```
cron-job.org (6×/hari, perkiraan sebelum shalat)
    │
    ▼ POST .../dispatches  {"event_type":"pray-notify","client_payload":{"prayer":"ashar"}}
GitHub Actions (~20 menit/job: sleep + kirim 5 offset)
    │
    ▼ prayer-notifier -prayer ashar
Discord webhook
```

## Prasyarat GitHub

1. Repo **public** (disarankan)
2. Secret **`WEBHOOK_TARGETS`** — dari `deploy/targets.json` via `make secret-push`
3. PAT untuk cron-job.org (scope `repo` atau Contents write)

> `make env-check` gagal jika `deploy/targets.json` masih berisi `REPLACE_ME`.

## Nonaktifkan cron lama

Hapus/disable job yang:
- Hit **tiap menit** (`* 4-18 * * *`, dll.)

## Buat 6 cronjob di cron-job.org

Timezone: **Asia/Jakarta**  
Method: **POST**  
URL: `https://api.github.com/repos/adamnasrudin03/pray-notify/dispatches`

Headers:

```
Accept: application/vnd.github+json
Authorization: Bearer ghp_xxxxxxxx
X-GitHub-Api-Version: 2022-11-28
Content-Type: application/json
```

### Tabel cronjob

| Job | Cron (WIB) | Body `client_payload.prayer` | Keterangan |
|-----|------------|------------------------------|------------|
| subuh | `15 4 * * *` | `subuh` | ~10 menit sebelum Subuh musim panas |
| dhuha | `15 6 * * *` | `dhuha` | sebelum Dhuha |
| dzuhur | `45 11 * * *` | `dzuhur` | sebelum Dzuhur |
| ashar | `0 15 * * *` | `ashar` | sebelum Ashar |
| maghrib | `45 17 * * *` | `maghrib` | sebelum Maghrib |
| isya | `0 19 * * *` | `isya` | sebelum Isya |

Cron di atas **perkiraan** (waktu shalat bergeser). Workflow menghitung jadwal **hari ini** lalu sleep sampai offset −10, −5, 0, +5, +10.

### Body request (contoh Ashar)

```json
{
  "event_type": "pray-notify",
  "client_payload": {
    "prayer": "ashar"
  }
}
```

Ganti `"ashar"` sesuai job.

## Test manual

**GitHub UI:** Actions → Prayer Notification → Run workflow → pilih shalat.

**curl:**

```bash
curl -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ghp_xxx" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -H "Content-Type: application/json" \
  https://api.github.com/repos/adamnasrudin03/pray-notify/dispatches \
  -d '{"event_type":"pray-notify","client_payload":{"prayer":"ashar"}}'
```

**Lokal:**

```bash
make run-window-test PRAYER=ashar   # tanpa sleep
```

## Perilaku window

| Shalat | Offset | Durasi job (kira-kira) |
|--------|--------|------------------------|
| Fardhu | −10, −5, 0, +5, +10 | ~20 menit |
| Dhuha | 0 saja | ~1 menit |

**6 job × ~20 menit ≈ 120 menit/hari** (~3.600/bulan). Muat di repo **public**; tight untuk repo **private** (2.000 menit/bulan).

## Troubleshooting

| Gejala | Solusi |
|--------|--------|
| HTTP 401 | PAT salah/expired |
| HTTP 405 Discord | URL webhook salah atau masih `REPLACE_ME` — cek `deploy/targets.json` |
| HTTP 422 | `event_type` harus `pray-notify` |
| Workflow kosong prayer | Cek `client_payload.prayer` |
| Notif terlewat | Cron terlalu telat; geser cron lebih awal |
| Limit Actions | Pastikan repo public + cron lama dimatikan |
| Job > 30 menit | Workflow timeout; cek jam shalat |
