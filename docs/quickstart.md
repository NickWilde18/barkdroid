# barkdroid Quick Start

Bark-compatible Android push client for self-hosters in China.

## What is this?

barkdroid lets you send push notifications to your Android phone using the Bark API format, via JPush (极光推送). You bring your own JPush credentials, build your own APK, and host your own server — no public service, no shared credentials.

## Architecture

```
curl GET /YOUR_KEY/Title/Body
        │
   ┌────▼─────────────┐
   │ barkdroid server  │  (Go, self-hosted, ~10 MB Docker image)
   │  - Bark-compat API│
   │  - SQLite device DB│
   │  - JPush provider  │
   └────────┬──────────┘
            │ JPush REST API
   ┌────────▼──────────┐
   │    极光 JPush       │
   └────────┬──────────┘
            │ system-level push
   ┌────────▼──────────┐
   │  Android App       │  (Kotlin, self-built APK)
   │  - registers device│
   │  - receives pushes │
   │  - taps to open URL│
   └───────────────────┘
```

## 7-Step Setup

### 1. Register a JPush account

Go to https://portal.jiguang.cn, create an account, and create an Android app.
You'll get:
- **AppKey** — put this in Android `config.properties`
- **Master Secret** — put this in server `config.yaml` (NEVER in APK!)

### 2. Clone & configure

```bash
git clone <this-repo> barkdroid
cd barkdroid
```

**Server config:** Copy and edit:
```bash
cp server/config.example.yaml server/config.yaml
vim server/config.yaml
```
Fill in `jpush.enabled: true`, `jpush.app_key`, `jpush.master_secret`.

**Android config:** Copy and edit:
```bash
cp android/config.example.properties android/config.properties
vim android/config.properties
```
Fill in `JPUSH_APP_KEY` and `SERVER_BASE_URL`.

### 3. Build & deploy server

```bash
# With Docker (recommended):
cd server
docker compose up -d

# Or without Docker (requires Go 1.22+):
go build -o barkdroid .
./barkdroid
```

### 4. Build Android APK

```bash
cd android

# Generate a keystore (one-time):
keytool -genkey -v -keystore barkdroid.keystore \
  -alias barkdroid -keyalg RSA -keysize 2048 -validity 36500

# Optional: configure signing (copy signing.example.gradle → signing.gradle)

# Build:
./gradlew assembleRelease

# APK at: app/build/outputs/apk/release/app-release.apk
```

### 5. Install APK on your phone

Transfer the APK to your phone and install it. Sign in to JPush console and verify your app is registered with the same package name and signature.

### 6. Get your device key

Open the barkdroid app on your phone. After a few seconds it will show:
```
Your Device Key
a1b2c3d4
```
This is your Bark-compatible key. Copy it.

### 7. Send a push!

```bash
# Bark-style GET:
curl "http://your-server:8080/a1b2c3d4/Test Title/Hello from barkdroid"

# Or POST:
curl -X POST http://your-server:8080/push \
  -H 'Content-Type: application/json' \
  -d '{"key":"a1b2c3d4","title":"Alert","body":"Something happened","url":"https://example.com"}'
```

Your Android phone should get a push notification. Tapping it opens the URL.

## Bark API Reference

| Method | Path | Description |
|--------|------|-------------|
| GET | `/:key/:body` | Push with body only (no title) |
| GET | `/:key/:title/:body` | Push with title and body |
| GET | `/:key/:title/:body?url=...` | Push with tap-to-open URL |
| POST | `/push` | JSON push (see below) |

**POST /push** request body:
```json
{
  "key": "a1b2c3d4",
  "title": "Alert",
  "body": "Server CPU > 90%",
  "url": "https://grafana.example.com"
}
```

**POST /register** request body:
```json
{
  "platform": "android",
  "push_provider": "jpush",
  "registration_id": "..."
}
```

## Security Notes

- **Master Secret stays on the server** — it's in `config.yaml`, never in the APK
- **AppKey goes in the APK** — it's only used for device registration, safe to embed
- Each user owns their own JPush account, package name, and signing key — no shared credentials
- If one user's JPush account has problems, it only affects them

## License

[MIT](LICENSE)
