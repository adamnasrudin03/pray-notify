# Setup cron-job.org untuk pray-notify

Panduan membuat cron job di [cron-job.org](https://cron-job.org) yang memicu workflow GitHub Actions via `repository_dispatch`.

## Prasyarat

- Repo: `adamnasrudin03/pray-notify`
- Secret `WEBHOOK_TARGETS` sudah diset di GitHub
- PAT GitHub dengan akses write ke repo

## Langkah

### 1. Buat cronjob

1. Buka [console.cron-job.org](https://console.cron-job.org)
2. Klik **Cronjobs** → **Create cronjob**

### 2. Tab Common

| Field | Nilai |
|-------|-------|
| Title | `pray-notify` |
| URL | `https://api.github.com/repos/adamnasrudin03/pray-notify/dispatches` |
| Enable job | ✅ |

### 3. Tab Schedule

- **Timezone:** `Asia/Jakarta`
- **Crontab expression:** pilih salah satu (buat 3 job terpisah):

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
| Request body | `{"event_type":"pray-notify"}` |

**Custom request headers:**

```
Accept: application/vnd.github+json
Authorization: Bearer ghp_xxxxxxxx
X-GitHub-Api-Version: 2022-11-28
Content-Type: application/json
```

Ganti `ghp_xxxxxxxx` dengan PAT Anda.

### 5. Simpan & test

Klik **Create** → **Run now** → cek **History** (harus HTTP 204).

Verifikasi di GitHub repo → tab **Actions** → workflow "Prayer Notification" jalan.

## Ulangi untuk 3 window waktu

Duplikat job di atas dengan cron expression berbeda (3 job total):

| Job name | Cron | Window WIB | Keterangan |
|----------|------|------------|------------|
| pray-notify-03 | `30-59 3 * * *` | 03:30–03:59 | Subuh -10 menit |
| pray-notify-04-18 | `* 4-18 * * *` | 04:00–18:59 | Subuh s.d Maghrib |
| pray-notify-19 | `* 19 * * *` | 19:00–19:59 | Isya + offset +10 menit |

Job 3 memakai **seluruh jam 19** agar pengingat **Isya +10 menit** tidak terlewat.

Di sisi code, window operasi juga dibatasi **03:30–19:59** (timezone target) + delivery window **90 detik** agar tidak spam meski cron jalan tiap menit.

## Troubleshooting

| Gejala | Solusi |
|--------|--------|
| HTTP 401 | PAT salah atau expired |
| HTTP 404 | Cek `GITHUB_REPO` / URL path |
| HTTP 422 | `event_type` harus `pray-notify` (sama dengan workflow) |
| Workflow tidak jalan | Pastikan workflow sudah di-push ke branch default |
| Notifikasi telat | Cek timezone job = `Asia/Jakarta` |
