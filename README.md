# JhopanStoreVPN

<p align="center">
  <img src="assets/icon256.png" alt="JhopanStoreVPN Logo" width="128">
</p>

<p align="center">
  <b>VLESS VPN Client — System-Wide Routing untuk semua aplikasi</b>
</p>

<p align="center">
  <a href="https://github.com/jhopan/JhopanStoreVPN/releases/latest">
    <img src="https://img.shields.io/github/v/release/jhopan/JhopanStoreVPN?style=flat-square&label=Versi%20Terbaru" alt="Release">
  </a>
  <a href="https://github.com/jhopan/JhopanStoreVPN/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/jhopan/JhopanStoreVPN?style=flat-square" alt="License">
  </a>
</p>

<p align="center">
  <b>Platform:</b> Windows • Linux • macOS • Android<br>
  <b>Dibuat oleh <a href="https://github.com/jhopan">JhopanStore</a></b>
</p>

---

## 📥 Download — Pilih yang sesuai perangkat kamu

> **Halaman download:** [Releases Page](https://github.com/jhopan/JhopanStoreVPN/releases/latest)

---

### 📱 Android

| File | Cocok untuk | Keterangan |
|---|---|---|
| `app-phone-arm64-v8a-release.apk` | **HP Android modern (2016+)** | ⭐ Paling direkomendasikan — paling kecil & cepat |
| `app-phone-armeabi-v7a-release.apk` | HP Android lama 32-bit | Untuk HP lama seperti Galaxy J1, S5 |
| `app-phone-universal-release.apk` | Semua HP ARM | Kalau tidak yakin HP kamu apa, download ini |
| `app-full-arm64-v8a-release.apk` | HP + emulator x86_64 | Versi lengkap arm64 |
| `app-full-universal-release.apk` | Semua ABI (HP + emulator) | Paling besar, untuk emulator Android Studio / BlueStacks |

> **Bingung pilih yang mana?**
> - HP biasa keluaran 2016 ke atas → **`app-phone-arm64-v8a-release.apk`**
> - Tidak tahu spesifikasi HP → **`app-phone-universal-release.apk`**
> - Emulator (Android Studio, BlueStacks, NoxPlayer) → **`app-full-universal-release.apk`**
> - HP 32-bit lama → **`app-phone-armeabi-v7a-release.apk`**

---

### 🖥️ Desktop

| File | Platform | Untuk siapa |
|---|---|---|
| `JhopanStoreVPN-windows.zip` | Windows 10/11 | Semua PC/laptop Windows 64-bit |
| `JhopanStoreVPN-linux.tar.gz` | Ubuntu, Debian, Fedora, dll | Semua PC/laptop Linux 64-bit |
| `JhopanStoreVPN-macos-arm64.tar.gz` | macOS Apple Silicon | Mac dengan chip **M1/M2/M3/M4** |
| `JhopanStoreVPN-macos-amd64.tar.gz` | macOS Intel | Mac dengan **Intel Core i5/i7/i9** |
| `JhopanStoreVPN-macos-universal.tar.gz` | macOS semua | Jalan di Mac manapun (ARM + Intel) |

> **Cara cek Mac kamu:**
> Apple menu (🍎) → About This Mac → lihat bagian **Chip** atau **Processor**
> - Tertulis "Apple M1/M2/M3/M4" → pilih **`macos-arm64`**
> - Tertulis "Intel Core" → pilih **`macos-amd64`**
> - Tidak yakin → pilih **`macos-universal`** (bisa dipakai di keduanya)

---

## 🚀 Cara Pakai

### Android

1. Download APK sesuai tabel di atas
2. Buka file APK → Install (`Allow from this source` jika diminta)
3. Buka app → masukkan **Address** (`example.com:443`) dan **UUID**
4. Atau tap ikon **Paste** (kanan atas) untuk import `vless://...` link
5. Tap **CONNECT** → setujui izin VPN
6. Status berubah "Connected" ✓

> **Penting:** Saat pertama buka, app akan meminta kamu menonaktifkan Battery Optimization.
> Lakukan ini agar VPN tidak terputus saat layar mati.

### Windows

1. Extract `JhopanStoreVPN-windows.zip`
2. Double-click `JhopanStoreVPN.exe`
3. Klik **Yes** pada UAC prompt (perlu untuk buat TUN interface)
4. Isi Address + UUID atau paste `vless://` link → CONNECT

### Linux

```bash
tar -xzf JhopanStoreVPN-linux.tar.gz
cd JhopanStoreVPN-linux
sudo ./JhopanStoreVPN
```

### macOS

```bash
# Pilih sesuai Mac kamu:
tar -xzf JhopanStoreVPN-macos-arm64.tar.gz     # M1/M2/M3/M4
# atau
tar -xzf JhopanStoreVPN-macos-amd64.tar.gz     # Intel
# atau
tar -xzf JhopanStoreVPN-macos-universal.tar.gz  # Tidak yakin

sudo open JhopanStoreVPN.app
```

---

## ⚙️ Pengaturan (Settings)

Semua pengaturan **disimpan otomatis** ke storage dan tidak hilang saat app ditutup atau di-update.

| Setting | Keterangan | Default |
|---|---|---|
| **Path** | WebSocket path di VLESS URI | `/vless` |
| **SNI** | Server Name Indication TLS (isi jika berbeda domain) | (dari Address) |
| **Host** | HTTP Host header WebSocket | (dari Address) |
| **DNS 1 / DNS 2** | DNS resolver saat VPN aktif | `8.8.8.8` / `8.8.4.4` |
| **Ping URL** | URL untuk tes latensi (HTTP HEAD request) | `https://dns.google` |
| **Auto Reconnect** | Otomatis sambung ulang jika koneksi putus | ON |
| **Allow Insecure TLS** | Skip verifikasi sertifikat TLS (untuk server self-signed) | ON |

> **DNS & HTTP Ping:**
> - Bisa diganti kapan saja — tidak ada yang di-overwrite update
> - Perubahan DNS berlaku di **koneksi berikutnya** (perlu disconnect → connect ulang)
> - Perubahan Ping URL berlaku di **ping berikutnya** setelah connect
> - Format DNS: alamat IPv4, contoh `1.1.1.1`, `8.8.8.8`, `208.67.222.222`

---

## 🌐 Bagikan VPN ke Device Lain (Hotspot Sharing) — Android

Share VPN ke laptop, HP lain, atau konsol via proxy — tanpa install apapun di device tujuan.

**Langkah:**
1. Aktifkan **Hotspot WiFi** di HP kamu
2. Sambungkan device lain ke hotspot tersebut
3. Di JhopanStoreVPN → tap ikon 📡 (kanan atas) → aktifkan **Proxy VPN**
4. Catat IP yang tertera (contoh: `192.168.43.1`)
5. Setting proxy di device lain sesuai tabel:

| Device | Tipe proxy | Cara setting |
|---|---|---|
| **Android** | HTTP (10809) | WiFi → Tahan jaringan → Ubah → Opsi Lanjutan → Proxy Manual → `[IP]:10809` |
| **iOS** | HTTP (10809) | Settings → Wi-Fi → (i) → Configure Proxy → Manual → `[IP]:10809` |
| **Windows** | HTTP (10809) | Settings → Network → Proxy → Manual → `[IP]:10809` |
| **Windows** | SOCKS5 (10808) | Via Proxifier atau SocksCap → Server: `[IP]:10808` |
| **Linux** | HTTP (10809) | Settings → Network → Proxy → HTTP: `[IP]:10809` |
| **Linux** | SOCKS5 (10808) | Terminal: `export all_proxy=socks5://[IP]:10808` |
| **macOS** | HTTP (10809) | System Settings → Network → [WiFi] → Proxies → Web Proxy: `[IP]:10809` |
| **macOS** | SOCKS5 (10808) | System Settings → Network → [WiFi] → Proxies → SOCKS Proxy: `[IP]:10808` |

> **Port 10809 = HTTP proxy** — dipakai Android, iOS, browser, sistem operasi umum
> **Port 10808 = SOCKS5 proxy** — dipakai tools lanjutan (Proxifier, terminal, dll), mendukung UDP

---

## 🔄 Auto-Reconnect

Ketika **Auto Reconnect** diaktifkan, VPN otomatis menyambung ulang dalam dua skenario:

1. **Xray crash** — proses Xray berhenti tiba-tiba → reconnect segera
2. **Jaringan kembali** — sinyal hilang lalu pulih → VPN reconnect otomatis (tidak perlu buka app)

Delay antar percobaan menggunakan **exponential backoff** (makin lama tiap gagal, max 60 detik):

| Percobaan | Delay |
|---|---|
| 1 | 3 detik |
| 2 | 6 detik |
| 3 | 12 detik |
| 4 | 24 detik |
| 5 | 48 detik |

Setelah 5 percobaan gagal berturut-turut, VPN berhenti. Perlu connect manual.

---

## 🏗️ Arsitektur Teknis

### Android

```
Semua aplikasi
      ↓
Android VpnService — TUN interface (kernel)
      ↓
tun2socks (native binary via JNI fork+exec, preserves TUN fd)
      ↓
Xray-core via libXray.aar (in-process — tidak ada subprocess)
   SOCKS5 inbound  127.0.0.1:10808  (dari tun2socks)
   HTTP   inbound  0.0.0.0:10809    (saat hotspot sharing aktif)
   VLESS outbound  WebSocket + TLS
      ↓
Internet → VLESS Server
```

### Desktop

```
Semua aplikasi
      ↓
TUN virtual interface (tun0)
      ↓
tun2socks binary (subprocess) → SOCKS5 127.0.0.1:10808
      ↓
Xray binary (subprocess)
   VLESS outbound  WebSocket + TLS
      ↓
Internet → VLESS Server
```

---

## 🛠️ Troubleshooting

### Android

<details>
<summary><b>VPN sering terputus saat layar mati</b></summary>

**Penyebab:** Android membatasi background service untuk hemat baterai.

**Solusi:**
1. Buka app → tap banner kuning → **Perbaiki** → izinkan "Unrestricted"
2. Atau: Settings → Battery → Battery Optimization → JhopanStoreVPN → **Don't Optimize**
3. HP Xiaomi/OPPO/Vivo: aktifkan **Autostart** di Settings → Apps → Autostart

</details>

<details>
<summary><b>"Connected" tapi internet tidak jalan</b></summary>

**Solusi:**
1. Cek address, UUID, path, SNI sudah benar (satu karakter salah = gagal)
2. Pastikan server VLESS kamu aktif
3. Disconnect → Connect ulang
4. Coba ganti DNS ke `1.1.1.1` di Settings

</details>

<details>
<summary><b>APK tidak bisa diinstall ("App not installed")</b></summary>

**Solusi:**
1. Settings → Security → Install Unknown Apps → izinkan
2. Coba install **`app-phone-universal-release.apk`** (support semua HP ARM)
3. Hapus dulu versi lama jika ada (uninstall dulu)

</details>

<details>
<summary><b>Hotspot sharing tidak bekerja di device lain</b></summary>

**Solusi:**
1. Pastikan switch **Proxy VPN** ON di halaman Hotspot
2. Gunakan IP yang tertera di app (bukan 192.168.1.1 atau gateway router)
3. Di device lain, matikan VPN lain yang mungkin aktif
4. Coba port 10809 (HTTP proxy) sebelum 10808 (SOCKS5)

</details>

### Desktop

<details>
<summary><b>"tun2socks not found" atau "Binary not found"</b></summary>

File binary sudah ada di dalam release zip/tar. Kalau hilang, download dari
[tun2socks releases v2.6.0](https://github.com/xjasonlyu/tun2socks/releases/tag/v2.6.0):
- Windows: `tun2socks-windows-amd64.zip` → ekstrak → `bin/tun2socks.exe`
- Linux: `tun2socks-linux-amd64.zip` → ekstrak → `bin/tun2socks`
- macOS: sesuai arsitektur → `bin/tun2socks`

</details>

<details>
<summary><b>Internet tidak pulih setelah disconnect</b></summary>

```bash
# Windows — reset routing
route -f
# Lalu reconnect WiFi / cabut-pasang kabel Ethernet

# Linux
sudo ip route flush cache
sudo systemctl restart NetworkManager

# macOS
sudo route flush
```

</details>

<details>
<summary><b>Windows: "Application already running"</b></summary>

Cek system tray (pojok kanan bawah taskbar). Kalau ada ikon app, klik kanan → Show.
Kalau tidak ada: Task Manager → cari `JhopanStoreVPN.exe` → End Task → buka lagi.

</details>

---

## 🔨 Build dari Source

### Android

```bash
# Requirements: Android Studio, JDK 17+, Android SDK 35

cd android

# APK untuk HP (ARM only — lebih kecil):
./gradlew assemblePhoneRelease

# APK untuk semua ABI (termasuk emulator):
./gradlew assembleFullRelease

# Output: android/app/build/outputs/apk/phone/release/
#         android/app/build/outputs/apk/full/release/
```

### Desktop

```bash
# Requirements: Go 1.21+, GCC/Clang, platform libs

# Linux (perlu: gcc, pkg-config, libgl1-mesa-dev, xorg-dev)
sudo apt install gcc pkg-config libgl1-mesa-dev xorg-dev
go build -ldflags="-s -w" -o JhopanStoreVPN .

# Windows (MinGW-w64, tanpa console window)
go build -ldflags="-s -w -H windowsgui" -o JhopanStoreVPN.exe .

# macOS (perlu Xcode Command Line Tools)
xcode-select --install
go build -ldflags="-s -w" -o JhopanStoreVPN .
```

File yang harus ada di sebelah executable:
- `xray` / `xray.exe` — dari [Xray-core Releases](https://github.com/XTLS/Xray-core/releases)
- `bin/tun2socks` / `bin/tun2socks.exe` — dari [tun2socks Releases](https://github.com/xjasonlyu/tun2socks/releases)

---

## Credits

- **[JhopanStore](https://github.com/jhopan)** — Developer & maintainer
- **[Xray-core](https://github.com/XTLS/Xray-core)** — VLESS proxy engine (MPL-2.0)
- **[libXray](https://github.com/XTLS/libXray)** — Android wrapper untuk Xray-core (MPL-2.0)
- **[tun2socks](https://github.com/xjasonlyu/tun2socks)** — TUN ke SOCKS5 bridge (GPL-3.0)
- **[Fyne](https://github.com/fyne-io/fyne)** — GUI toolkit Go (BSD 3-Clause)

## Lisensi

MIT — © 2026 JhopanStore

---

<p align="center">
  Ada pertanyaan atau bug? Buka <a href="https://github.com/jhopan/JhopanStoreVPN/issues">issue</a> di GitHub
</p>
