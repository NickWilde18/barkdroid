# barkdroid

**Bark-compatible Android push client for self-hosters in China.**

Send push notifications to your Android phone using the Bark API format — via JPush (极光推送), no FCM required. You bring your own JPush credentials, build your own APK, and self-host the server. No public service, no shared credentials.

## Why

- **Bark** is great for iOS, but doesn't solve Android
- **FCM** is unreliable in China — phones kill Google services aggressively
- **Existing solutions** (PushDeer, PushMe, ntfy, Gotify) either rely on WebSocket background connections, depend on FCM, or don't support pluggable Chinese push providers

barkdroid combines **Bark API compatibility** + **JPush's system-level push channel** + **self-hosted server** + **BYO credentials**.

## Quick Start

```bash
# 1. Register a JPush account at https://portal.jiguang.cn
# 2. Configure your credentials
cp server/config.example.yaml server/config.yaml    # edit
cp android/config.example.properties android/config.properties  # edit

# 3. Start the server
cd server && docker compose up -d

# 4. Build and install the Android app
cd android && ./gradlew assembleRelease

# 5. Open the app to get your device key, then:
curl "http://your-server:8080/YOUR_KEY/Server Alert/CPU > 90%"
```

Full guide: [docs/quickstart.md](docs/quickstart.md)

## Architecture

```
curl GET /YOUR_KEY/Title/Body
        │              Bark-compatible API
   ┌────▼─────────────┐
   │  barkdroid server │  Go + GoFrame + SQLite
   └────────┬──────────┘
            │  JPush REST API
   ┌────────▼──────────┐
   │     极光 JPush      │  System-level push channel
   └────────┬──────────┘
            │  Proprietary protocol
   ┌────────▼──────────┐
   │   Android App      │  Kotlin + Jetpack Compose
   └───────────────────┘
```

## Features

- ✅ **Bark API compatible** — existing Bark scripts work with just a URL change
- ✅ **JPush provider** — system-level push, survives background app killing
- ✅ **Self-hosted** — Docker deployment, single ~10MB binary
- ✅ **BYO credentials** — each user owns their JPush account, package name, and signing key
- ✅ **iOS Bark forwarding** — optionally forward pushes to your iOS Bark device
- ✅ **Pluggable providers** — JPush today, Getui/TPNS architecturally ready
- ✅ **Zero FCM dependency** — no Google services required

## Tech Stack

| Component | Technology |
|-----------|------------|
| Server | Go + [GoFrame](https://github.com/gogf/gf) |
| Database | SQLite (pure Go, WAL mode) |
| Android App | Kotlin + Jetpack Compose |
| Push Provider | JPush REST API (no SDK needed server-side) |
| Deployment | Docker + docker-compose |

## FAQ

**Why not just use ntfy/Gotify/PushDeer?**

ntfy and Gotify are excellent general-purpose notification platforms, but they don't solve the core problem: reliable push delivery on Chinese Android phones where background processes are aggressively killed. PushDeer was close but is abandoned and all users shared a single MiPush credential (which eventually failed).

**Why do I need to bring my own JPush account?**

Because a shared public credential is a single point of failure. When PushDeer's MiPush authorization was revoked, all Android users lost push notifications simultaneously. With your own credentials, your push channel is yours alone.

**Does JPush cost money?**

JPush has a free tier that includes plenty of pushes for personal/self-hosted use. Check their [pricing page](https://www.jiguang.cn/push) for current details.

## License

[MIT](LICENSE)

> Built with Go + GoFrame · Android Kotlin/Compose · JPush
