# JhopanStoreVPN — Optimization Summary & Xray Configs

> Branch: `optimation` | Platform: Android + Desktop (Windows/Linux/macOS)

---

## Ringkasan Seluruh Optimasi

### Desktop (`main` branch — Go + Fyne)

#### Ronde 1 — Optimasi Awal

| #   | File                  | Perubahan                                                                                                  | Dampak                                                       |
| --- | --------------------- | ---------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------ |
| 1   | `core/ping/ping.go`   | Rewrite total: burst 3× via SOCKS5 port 10808, inline RFC1928 dialer, goroutine mati sendiri setelah burst | Tidak ada goroutine ping yang terus jalan, hemat CPU & radio |
| 2   | `core/xray/config.go` | Tambah `loglevel: none`, hapus HTTP inbound (port 10809), matikan sniffing                                 | Kurangi CPU per-packet, single SOCKS5 inbound                |
| 3   | `core/tun/tun.go`     | `-loglevel info` → `-loglevel silent`                                                                      | Kurangi disk I/O & write syscall tun2socks                   |
| 4   | `main.go`             | MTU 1500 → 1400, SocksAddr 10809 → 10808 (bug fix), port readiness check ke 10808                          | Kurangi fragmentasi, fix port mismatch                       |

#### Ronde 2 — Desktop Deep Optimization (Win/Linux/macOS)

| #   | File                   | Perubahan                                                                                          | Dampak                                                          |
| --- | ---------------------- | -------------------------------------------------------------------------------------------------- | --------------------------------------------------------------- |
| 5   | `core/xray/process.go` | Hapus pembuatan log file (loglevel none → output discarded), fix cleanup temp dir semua error path | Tidak ada file xray.log kosong tiap connect, tidak ada leak dir |
| 6   | `core/xray/process.go` | `onCrash` trigger untuk semua exit (hapus syarat `err != nil`)                                     | Deteksi crash termasuk xray exit code 0                         |
| 7   | `core/tun/tun.go`      | `time.Sleep(1s)` → `select { processDone / time.After(1s) }` (deteksi crash dini)                  | Koneksi gagal terdeteksi segera, bukan setelah 1 s penuh        |
| 8   | `core/tun/tun.go`      | Hapus `time.Sleep(500ms)` di `Close()`                                                             | Disconnect lebih cepat 500 ms                                   |
| 9   | `core/tun/tun.go`      | Tambah `-tcp-no-delay` dan `-udp-timeout 30s` ke args tun2socks (setara Android)                   | Latensi TCP lebih rendah, UDP memory free lebih cepat           |
| 10  | `core/tun/tun.go`      | Binary detection: `tun2socks.exe` hardcode → `runtime.GOOS`-aware (Win/Linux/macOS)                | Bisa build & run di Linux dan macOS tanpa modifikasi            |
| 11  | `main.go`              | Hapus auto-reconnect ticker 10 s (redundant — `onCrash` di process.go sudah handle ini)            | Hilangkan 1 goroutine background saat connected                 |
| 12  | `core/xray/config.go`  | Hapus outbound `block` (blackhole) yang tidak direferensikan di routing                            | Config lebih bersih, Xray tidak allocate handler unused         |

### Android (`optimation` branch — Kotlin + libXray in-process)

#### Ronde 1 — Battery & Performance Dasar

| #   | File                  | Perubahan                                         | Dampak                                       |
| --- | --------------------- | ------------------------------------------------- | -------------------------------------------- |
| 1   | `MainViewModel.kt`    | Ping burst 3× lalu berhenti (ganti infinite loop) | Tidak ada polling background setelah connect |
| 2   | `JhopanVpnService.kt` | MTU 1400 di 2 tempat (connect + reconnect)        | Kurangi fragmentasi IPv4                     |
| 3   | `Tun2socksManager.kt` | `-loglevel silent`                                | Kurangi I/O subprocess                       |
| 4   | `XrayManager.kt`      | Monitor interval 5 s → 10 s                       | Kurangi CPU wake cycles                      |

#### Ronde 2 — Deep Optimization

| #   | File                  | Perubahan                                                                                                                       | Dampak                                                              |
| --- | --------------------- | ------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------- |
| 5   | `XrayManager.kt`      | Matikan semua stats (policy/system), tambah `sockopt.tcpFastOpen`, freedom `UseIPv4`, hapus `Log.d` version check & config dump | Kurangi CPU per-paket, -1 RTT saat reconnect                        |
| 6   | `Tun2socksManager.kt` | Tutup pipe fd langsung setelah fork, tambah `-tcp-no-delay`, `-udp-timeout 30s`, monitor thread `isDaemon=true`                 | Kurangi fd leak, latensi lebih rendah, UDP memory bebas lebih cepat |
| 7   | `JhopanVpnService.kt` | Ganti `Thread.sleep(1000)` dengan active TCP probe port 10808 (250 ms × 40) di connect & reconnect                              | Connect lebih cepat, tidak buang waktu tunggu fixed                 |
| 8   | `build.gradle.kts`    | Hapus dependensi `datastore-preferences` yang tidak terpakai                                                                    | APK lebih kecil                                                     |

#### Ronde 3 — Connection Reuse, UI, Mux

| #   | File               | Perubahan                                                                                      | Dampak                                                |
| --- | ------------------ | ---------------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| 9   | `MainViewModel.kt` | Tambah `pingJob: Job?`, cancel di `disconnect()`                                               | Ping tidak lanjut jalan setelah VPN diputus           |
| 10  | `MainViewModel.kt` | `restartVpn()` ganti 15×`delay(1000)` polling → `StateFlow.first{}` + `withTimeoutOrNull(20s)` | Tidak ada polling UI, reaktif murni                   |
| 11  | `MainViewModel.kt` | `Proxy` object dibuat sekali di luar `repeat{}`, hapus `conn.disconnect()`                     | Reuse HTTP keep-alive socket antar burst ping         |
| 12  | `XrayManager.kt`   | Tambah `mux: {enabled: true, concurrency: 8}`                                                  | 1 TLS connection untuk banyak stream, hemat handshake |

#### Ronde 4 — Build Quality

| #   | File                  | Perubahan                                                        | Dampak                                                 |
| --- | --------------------- | ---------------------------------------------------------------- | ------------------------------------------------------ |
| 13  | `gradle.properties`   | `android.enableR8.fullMode=true`                                 | R8 eliminasi dead code lebih agresif, APK lebih kecil  |
| 14  | `JhopanVpnService.kt` | `IMPORTANCE_DEFAULT` → `IMPORTANCE_LOW` notifikasi               | Tidak ada suara/getar tiap update status VPN           |
| 15  | `MainViewModel.kt`    | Tambah `timeoutJob: Job?`, cancel di CONNECTED/FAILED/disconnect | Coroutine 30 s tidak nganggur setelah koneksi resolved |

#### Ronde 5 — Fix Bug, Keamanan, R8, ABI Splits

| #   | File                  | Perubahan                                                                         | Dampak                                                  |
| --- | --------------------- | --------------------------------------------------------------------------------- | ------------------------------------------------------- |
| 16  | `MainViewModel.kt`    | `disconnect()` cancel `timeoutJob`                                                | Bug fix: coroutine leak                                 |
| 17  | `MainViewModel.kt`    | `companion object val HOTSPOT_IP_172_REGEX`                                       | Regex tidak compile ulang tiap loop iterasi             |
| 18  | `JhopanVpnService.kt` | Xray retry rekursif → `while` loop iteratif                                       | Tidak ada stack frame menumpuk saat retry               |
| 19  | `XrayManager.kt`      | Hapus outbound `blackhole` yang tidak direferensikan di routing                   | Config lebih bersih, Xray tidak allocate handler unused |
| 20  | `XrayManager.kt`      | Hapus `Log.d` yang allocate `List<String>` tiap DNS resolve                       | Kurangi GC pressure                                     |
| 21  | `Tun2socksManager.kt` | `Thread.sleep(500)` → poll 10×50 ms                                               | tun2socks ready rata-rata ~50 ms bukan 500 ms           |
| 22  | `proguard-rules.pro`  | Ganti wildcard `-keep class com.jhopanstore.vpn.**` → keep spesifik per kebutuhan | R8 full mode kini bisa optimasi kode internal sendiri   |
| 23  | `AndroidManifest.xml` | `allowBackup="false"`                                                             | Credentials VLESS tidak ikut backup ADB/Google          |
| 24  | `build.gradle.kts`    | ABI splits: `arm64-v8a`, `armeabi-v7a`, `x86_64`, `x86` + `universal` APK         | 5 variant APK, distribusi ringan per arsitektur         |

---

## Alur Traffic

```
[ Android ]                              [ Desktop ]

Apps
 │
 ▼
TUN fd (MTU 1400)                       Fyne UI (port readiness check :10808)
 │                                           │
 ▼                                           ▼
tun2socks                               tun2socks
 -loglevel silent                        -loglevel silent
 -tcp-no-delay                           -tcp-no-delay
 -udp-timeout 30s                        -udp-timeout 30s
 │                                           │
 ▼ SOCKS5                               ▼ SOCKS5
Xray (libXray in-process)              Xray (subprocess)
 port 10808 → internet                  port 10808 → internet
```

---

## Xray Config — Desktop

File: `core/xray/config.go` → `GenerateConfig()`

```json
{
  "log": {
    "loglevel": "none"
  },
  "dns": {
    "servers": ["8.8.8.8", "8.8.4.4"]
  },
  "inbounds": [
    {
      "tag": "socks-in",
      "port": 10808,
      "listen": "127.0.0.1",
      "protocol": "socks",
      "settings": {
        "udp": true
      }
    }
  ],
  "outbounds": [
    {
      "tag": "proxy",
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "address": "<server-domain-or-ip>",
            "port": 443,
            "users": [
              {
                "id": "<uuid>",
                "encryption": "none",
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "security": "tls",
        "tlsSettings": {
          "serverName": "<sni>",
          "allowInsecure": false
        },
        "wsSettings": {
          "path": "/vless",
          "host": "<host-header>"
        },
        "sockopt": {
          "tcpFastOpen": true,
          "tcpNoDelay": true,
          "tcpKeepAliveIdle": 60,
          "tcpKeepAliveInterval": 30
        }
      }
    },
    {
      "tag": "direct",
      "protocol": "freedom",
      "settings": {}
    }
  ],
  "routing": {
    "domainStrategy": "AsIs",
    "rules": [
      {
        "type": "field",
        "ip": ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.0/8"],
        "outboundTag": "direct"
      }
    ]
  }
}
```

> **Catatan Desktop:**
>
> - Hanya 1 inbound (SOCKS5 port 10808). HTTP inbound (10809) dihapus.
> - **Sniffing dinonaktifkan** (tidak ada blok `sniffing` di inbound) — tidak berguna untuk full-tunnel VPN dan memakan CPU.
> - Local IP ranges langsung di-direct, tidak lewat proxy.
> - `tcpNoDelay: true` — matikan Nagle's algorithm di koneksi Xray→server → latency lebih rendah.
> - `tcpKeepAliveIdle: 60` + `tcpKeepAliveInterval: 30` — OS kirim probe TCP setiap 30 s setelah 60 s idle → deteksi koneksi mati maksimal 30 s.
> - tun2socks sekarang juga dipanggil dengan `-tcp-no-delay` dan `-udp-timeout 30s` (setara Android).
> - Binary `tun2socks` auto-detected per platform: `.exe` di Windows, tanpa ekstensi di Linux/macOS.

---

## Xray Config — Android

File: `XrayManager.kt` → `buildConfig()`

```json
{
  "log": {
    "loglevel": "none"
  },
  "policy": {
    "system": {
      "statsInboundUplink": false,
      "statsInboundDownlink": false,
      "statsOutboundUplink": false,
      "statsOutboundDownlink": false
    }
  },
  "inbounds": [
    {
      "tag": "socks-in",
      "port": 10808,
      "listen": "127.0.0.1",
      "protocol": "socks",
      "settings": {
        "udp": true
      }
    }
    // Jika hotspot sharing aktif, tambah:
    // {
    //   "tag": "http-in",
    //   "port": 10809,
    //   "listen": "0.0.0.0",
    //   "protocol": "http"
    // }
    // dan listen socks-in berubah ke "0.0.0.0"
  ],
  "outbounds": [
    {
      "tag": "proxy",
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "address": "<resolved-ip-or-domain>",
            "port": 443,
            "users": [
              {
                "id": "<uuid>",
                "encryption": "none"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "security": "tls",
        "wsSettings": {
          "path": "/vless",
          "headers": {
            "Host": "<host-header>"
          }
        },
        "tlsSettings": {
          "serverName": "<sni>",
          "allowInsecure": false
        },
        "sockopt": {
          "tcpFastOpen": true,
          "tcpNoDelay": true,
          "tcpKeepAliveIdle": 60,
          "tcpKeepAliveInterval": 30
        }
      },
      "mux": {
        "enabled": true,
        "concurrency": 8
      }
    },
    {
      "tag": "direct",
      "protocol": "freedom",
      "settings": {
        "domainStrategy": "UseIPv4"
      }
    }
    // Jika cloudflare workers (*.workers.dev / *.pages.dev), tambah:
    // { "tag": "dns-out", "protocol": "dns" }
  ],
  "routing": {
    "domainStrategy": "AsIs",
    "rules": [
      // Mode normal (non-cloudflare, non-hotspot): bypass local IP
      {
        "type": "field",
        "ip": ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.0/8"],
        "outboundTag": "direct"
      }
      // Mode cloudflare workers: DNS port 53 → dns-out (TCP DNS)
      // Mode hotspot sharing: tidak ada rule direct (semua lewat proxy)
    ]
  },
  "dns": {
    "servers": ["8.8.8.8", "8.8.4.4"],
    "queryStrategy": "UseIPv4"
    // Jika cloudflare workers: ["tcp://8.8.8.8", "tcp://8.8.4.4"]
  }
}
```

> **Catatan Android:**
>
> - **Sniffing dinonaktifkan** — tidak ada blok `sniffing` di inbound, memang sengaja (tidak berguna untuk full-tunnel, hanya buang CPU).
> - `tcpNoDelay: true` — matikan Nagle's algorithm di koneksi Xray→server, kurangi latency.
> - `tcpKeepAliveIdle: 60` + `tcpKeepAliveInterval: 30` — probe TCP tiap 30 s setelah 60 s idle → NAT timeout terdeteksi cepat, auto-reconnect lebih responsif.
> - `mux concurrency:8` — multipleks 8 stream TCP dalam 1 koneksi WebSocket/TLS, kurangi TLS handshake.
> - `tcpFastOpen: true` — kurangi 1 RTT saat reconnect (tersedia Android API 21+).
> - `statsInbound/Outbound: false` — Xray tidak tracking traffic per-koneksi, kurangi CPU.
> - `queryStrategy: UseIPv4` + `freedom UseIPv4` — paksa IPv4, hindari latency IPv6 fallback.
> - Server domain di-resolve sebelum TUN dibuat, hasilnya (IP) dipakai di config → tidak ada DNS inside TUN.
> - Hotspot sharing: `listen 0.0.0.0`, tambah HTTP inbound port 10809 (HP lain pakai WiFi manual proxy HTTP).
> - Cloudflare Workers mode: DNS otomatis pakai TCP bukan UDP (Workers tidak support UDP socket).

---

## Perbedaan Utama Android vs Desktop

| Aspek             | Desktop                                     | Android                                          |
| ----------------- | ------------------------------------------- | ------------------------------------------------ |
| Xray runtime      | Subprocess (`exec.Command`)                 | In-process (`libXray.aar` via gomobile)          |
| HTTP inbound      | Tidak ada (dihapus)                         | Ada saat hotspot sharing aktif (`0.0.0.0:10809`) |
| Mux               | Tidak ada                                   | `enabled:true, concurrency:8`                    |
| TCP Fast Open     | Tidak ada                                   | `sockopt.tcpFastOpen: true`                      |
| Stats tracking    | Tidak ada (policy tidak di-set)             | Explisit `false` semua                           |
| DNS strategy      | `servers: [ip1, ip2]` biasa                 | `+queryStrategy: UseIPv4`                        |
| Freedom outbound  | `domainStrategy` tidak di-set               | `domainStrategy: UseIPv4`                        |
| Socket protection | Tidak perlu (no TUN routing)                | `DialerController.protectFd()` via libXray       |
| Reconnect         | Auto-reconnect via `onProcessDied` callback | Iterative `while` loop dengan WakeLock           |

---

## Build Output Android

Jalankan `./gradlew assembleRelease` dari folder `android/`:

```
app/build/outputs/apk/release/
├── app-arm64-v8a-release.apk      ← Rekomendasi: HP 2016+ (paling kecil ~15-20 MB)
├── app-armeabi-v7a-release.apk    ← HP ARM 32-bit lama
├── app-x86_64-release.apk         ← Emulator / Chromebook x86_64
├── app-x86-release.apk            ← Emulator x86 legacy
└── app-universal-release.apk      ← Semua ABI (sideload-friendly, paling besar)
```

---

## Status Branch

| Branch       | Isi                                                                     |
| ------------ | ----------------------------------------------------------------------- |
| `main`       | Desktop optimizations selesai, pushed                                   |
| `optimation` | Semua Android optimizations (5 ronde) + Desktop Round 2 (8 fix), pushed |
