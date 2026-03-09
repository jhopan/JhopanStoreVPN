# JhopanStoreVPN

<p align="center">
  <img src="assets/icon256.png" alt="JhopanStoreVPN Logo" width="128">
</p>

<p align="center">
  Lightweight VLESS VPN client with <b>TUN mode</b> for system-wide routing (Desktop) and VpnService (Android)
</p>

<p align="center">
  <b>Developed by <a href="https://github.com/jhopan">JhopanStore</a></b>
</p>

---

## Features

### Desktop (Windows / Linux / macOS)
- ✅ **TUN Mode** - System-wide VPN routing (all applications work: WhatsApp, Telegram, Discord, games, etc.)
- ✅ VLESS over WebSocket + TLS via Xray-core
- ✅ Clean GUI with system tray support
- ✅ External tun2socks integration for reliable TUN device management
- ✅ Clipboard import (`vless://` URI)
- ✅ HTTP ping monitoring with latency display
- ✅ Auto-reconnect on disconnect
- ✅ No console window — pure native GUI

### Android (APK)
- ✅ **VpnService API** - True system-wide VPN (no root required)
- ✅ VLESS protocol via libXray
- ✅ **Hotspot Sharing** - Share VPN with WiFi hotspot (HTTP proxy on port 10809)
- ✅ **Config Persistence** - Settings saved automatically to prevent data loss
- ✅ **Service Stability** - Foreground service prevents system from killing app
- ✅ Material Design 3 UI with Jetpack Compose
- ✅ Real-time statistics (upload/download speed)
- ✅ Import config from QR code or manual input

## Download

Go to [**Releases**](https://github.com/jhopan/JhopanStoreVPN/releases) to download:

### Desktop

| Platform | File | How to Run |
|----------|------|------------|
| **Windows** | `JhopanStoreVPN-windows.zip` | Extract → **Run as Administrator** → `JhopanStoreVPN.exe` |
| **Linux** | `JhopanStoreVPN-linux.tar.gz` | Extract → `sudo ./JhopanStoreVPN` or `sudo ./launch.sh` |
| **macOS** | `JhopanStoreVPN-macos.tar.gz` | Extract → Run `JhopanStoreVPN.app` |

> **Note**: Desktop TUN mode requires **administrator/root privileges** to create virtual network interfaces and configure routing tables.

### Android

| File | Minimum Version | Install |
|------|-----------------|---------|
| `JhopanStoreVPN.apk` | Android 7.0 (API 24) | Download → Enable "Unknown Sources" → Install |

## Quick Start

### Desktop

1. Download and extract the archive for your OS
2. **Run as Administrator** (Windows: right-click → "Run as Administrator" | Linux/macOS: use `sudo`)
3. Enter server address and UUID, or paste a `vless://` link via the clipboard button
4. Click **CONNECT**
5. All applications will route through VPN automatically (no manual proxy configuration needed)

### Android

1. Download and install `JhopanStoreVPN.apk`
2. Open app and enter server details or scan QR code
3. Tap **CONNECT** and approve VPN permission
4. Enable **Hotspot Sharing** if you want to share VPN with other devices via WiFi hotspot

## Build from Source

### Desktop Requirements
- Go 1.21+
- GCC (MinGW on Windows)
- Platform graphics libraries:
  - **Linux**: `sudo apt install gcc pkg-config libgl1-mesa-dev xorg-dev`
  - **macOS**: `xcode-select --install`
  - **Windows**: MinGW-w64

### Desktop Build

```bash
# 1. Clone repository
git clone https://github.com/jhopan/JhopanStoreVPN.git
cd JhopanStoreVPN

# 2. Download tun2socks binary
# Windows: Download from https://github.com/xjasonlyu/tun2socks/releases
# Place tun2socks.exe in bin/ folder

# 3. Build
# Linux / macOS
go build -ldflags="-s -w" -o JhopanStoreVPN .

# Windows (no console window)
go build -ldflags="-s -w -H windowsgui" -o JhopanStoreVPN.exe .
```

**Required files next to executable:**
- `xray` or `xray.exe` - Download from [Xray-core Releases](https://github.com/XTLS/Xray-core/releases)
- `bin/tun2socks.exe` (Windows) - Download from [tun2socks Releases](https://github.com/xjasonlyu/tun2socks/releases)

### Android Build

```bash
# 1. Prerequisites
# - Android Studio Hedgehog+ (2023.1.1+)
# - JDK 17+
# - Android SDK 34+

# 2. Build APK
cd android
./gradlew assembleRelease

# Output: android/app/build/outputs/apk/release/app-release.apk
```

## Technical Architecture

### Desktop TUN Mode
```
Applications → TUN Device (tun0)
            ↓
        tun2socks process → SOCKS5 (Xray) → VLESS Server → Internet
            ↓
    Routing Table (0.0.0.0/0 via TUN)
```

- **TUN Device**: Created by external `tun2socks` binary
- **Routing**: Configured via `netsh` (Windows), `ip` (Linux), `route` (macOS)
- **Process Management**: Go exec wrapper manages tun2socks lifecycle
- **Requirements**: Admin/root privileges for interface creation

### Android VpnService
```
Applications → Android VpnService API
            ↓
        TUN Device (tun0) → libXray → VLESS Server → Internet
            ↓
    Packet interception (no root required)
```

- **VpnService**: Android framework API handles TUN creation
- **libXray**: Native library (AAR) for high-performance VLESS protocol
- **Hotspot Sharing**: HTTP proxy on port 10809 for WiFi clients
- **Foreground Service**: Prevents system from killing VPN connection

## Troubleshooting

### Desktop

**"Failed to start tun2socks"**
- Solution: Run as Administrator/root
- Make sure `bin/tun2socks.exe` exists next to the executable

**"tun2socks.exe not found"**
- Download from [tun2socks releases](https://github.com/xjasonlyu/tun2socks/releases/tag/v2.5.2)
- Extract and place in `bin/` folder

**Windows may prompt "Wintun driver installation"**
- Click "Yes" or "Allow" to install the TUN driver
- This is required for TUN mode to work

**Connection successful but no internet**
- Check if Xray is running: look for `xray.exe` process
- Verify VLESS server is accessible
- Check routing table: `route print 0.0.0.0` should show route via TUN

### Android

**VPN disconnects randomly**
- App uses foreground service to prevent system kill
- Check battery optimization settings: disable for JhopanStoreVPN
- Some OEMs (Xiaomi, Huawei) may require additional "autostart" permission

**Hotspot not sharing VPN**
- Make sure "Enable Hotspot Sharing" is checked in settings
- Check if port 10809 is accessible from hotspot clients
- Some devices may have restrictions on sharing VPN connections

## Credits

- **[JhopanStore](https://github.com/jhopan)** — Developer & maintainer
- **[Xray-core](https://github.com/XTLS/Xray-core)** — VLESS/VMess/Trojan proxy engine (MPL-2.0 license)
- **[libXray](https://github.com/XTLS/libXray)** — Android wrapper for Xray-core (MPL-2.0 license)
- **[tun2socks](https://github.com/xjasonlyu/tun2socks)** — TUN device to SOCKS5 proxy adapter (GPL-3.0 license)
- **[Fyne](https://github.com/fyne-io/fyne)** — Cross-platform GUI toolkit for Go (BSD 3-Clause license)
- **[golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys)** — Go extended system packages (BSD 3-Clause license)

## License

MIT — Copyright (c) 2026 JhopanStore

---

<p align="center">
  <b>Questions or issues?</b> Open an <a href="https://github.com/jhopan/JhopanStoreVPN/issues">issue</a> on GitHub
</p>