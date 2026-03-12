# JhopanStoreVPN

<p align="center">
  <img src="assets/icon256.png" alt="JhopanStoreVPN Logo" width="128">
</p>

<p align="center">
  <b>Lightweight VLESS VPN Client with Full System-Wide Routing</b>
</p>

<p align="center">
  <a href="https://github.com/jhopan/JhopanStoreVPN/releases">
    <img src="https://img.shields.io/github/v/release/jhopan/JhopanStoreVPN?style=flat-square" alt="Release">
  </a>
  <a href="https://github.com/jhopan/JhopanStoreVPN/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/jhopan/JhopanStoreVPN?style=flat-square" alt="License">
  </a>
  <a href="https://github.com/jhopan/JhopanStoreVPN/stargazers">
    <img src="https://img.shields.io/github/stars/jhopan/JhopanStoreVPN?style=flat-square" alt="Stars">
  </a>
</p>

<p align="center">
  <b>Cross-Platform:</b> Windows • Linux • macOS • Android<br>
  <b>Developed by <a href="https://github.com/jhopan">JhopanStore</a></b>
</p>

---

## 🌟 Highlights

- 🚀 **True System-Wide VPN** — Routes ALL applications (WhatsApp, Telegram, Discord, games, browsers)
- 🔒 **VLESS Protocol** — Modern, efficient, and secure via Xray-core
- 🖥️ **TUN Mode** — Virtual network interface for transparent routing (Desktop)
- 📱 **Native Android VPN** — VpnService API without root
- 🎯 **Single Instance** — No duplicate processes, smart mutex/lock system
- ⚡ **Auto Admin Elevation** — UAC prompt on Windows, seamless sudo on Linux/macOS
- 🌐 **Hotspot Sharing** — Share VPN connection via WiFi (Android)
- 🎨 **Clean Modern UI** — Fyne v2 (Desktop), Material Design 3 (Android)
- 📊 **Real-Time Stats** — Ping monitoring, traffic statistics
- 🔄 **Auto-Reconnect** — Stable connection with crash recovery

---

## 📋 Features

### 🖥️ Desktop (Windows / Linux / macOS)

#### Core VPN Features

- ✅ **TUN Mode** - Virtual network interface for true system-wide routing
- ✅ **All Apps Work** - WhatsApp Desktop, Telegram, Discord, Steam, Epic Games, etc.
- ✅ **VLESS Protocol** - WebSocket + TLS via Xray-core
- ✅ **Auto Admin Elevation** - UAC prompt (Windows) / sudo (Linux/macOS) on startup
- ✅ **External tun2socks** - Reliable TUN device management via proven binary
- ✅ **Smart Routing Cleanup** - Internet automatically restored on disconnect

#### UI & Usability

- ✅ **Modern GUI** - Clean interface powered by Fyne v2
- ✅ **System Tray** - Minimize to tray, stays in background
- ✅ **Single Instance** - Prevents duplicate processes (mutex/lock)
- ✅ **Clipboard Import** - Paste `vless://` links directly
- ✅ **Real-Time Ping** - HTTP latency monitoring
- ✅ **Auto-Reconnect** - Automatic recovery on connection loss
- ✅ **No Console Window** - Pure GUI application

### 📱 Android (APK)

#### Core VPN Features

- ✅ **VpnService API** - True system-wide VPN (no root required)
- ✅ **VLESS via libXray** - Native performance with Xray-core AAR
- ✅ **Hotspot Sharing** - Share VPN with other devices (HTTP proxy port 10809)
- ✅ **Always-On VPN** - Foreground service prevents system kill
- ✅ **Config Persistence** - Settings automatically saved

#### UI & Stats

- ✅ **Material Design 3** - Modern Android UI with Jetpack Compose
- ✅ **Real-Time Statistics** - Upload/download speed, traffic counter
- ✅ **QR Code Import** - Scan config from QR codes
- ✅ **Battery Optimized** - Efficient background operation

---

## 📥 Download

---

## 📥 Download

### 💾 Latest Release

Go to [**Releases Page**](https://github.com/jhopan/JhopanStoreVPN/releases/latest) to download the latest version (v2.0.0).

### 🖥️ Desktop

| Platform          | File                          | Size   | Requirements             |
| ----------------- | ----------------------------- | ------ | ------------------------ |
| **Windows 10/11** | `JhopanStoreVPN-windows.zip`  | ~42 MB | Administrator privileges |
| **Linux (x64)**   | `JhopanStoreVPN-linux.tar.gz` | ~38 MB | Root/sudo access         |
| **macOS (ARM64)** | `JhopanStoreVPN-macos.tar.gz` | ~39 MB | macOS 11+                |

### 📱 Android

| File                 | Size   | Requirements          |
| -------------------- | ------ | --------------------- |
| `JhopanStoreVPN.apk` | ~30 MB | Android 7.0+ (API 24) |

> **⚠️ Important Notes:**
>
> - **Desktop**: TUN mode requires **administrator/root privileges** to create virtual network interfaces
> - **Windows**: UAC prompt will appear automatically on startup
> - **Linux/macOS**: Run with `sudo` or configure passwordless sudo for convenience
> - **Android**: Enable "Install from Unknown Sources" in settings

---

## 🚀 Quick Start

### 🖥️ Desktop Setup

#### Windows

```powershell
# 1. Extract downloaded zip
# 2. Double-click JhopanStoreVPN.exe
# 3. Click "Yes" on UAC prompt
# 4. Enter VLESS server details or paste vless:// link
# 5. Click CONNECT
```

**First Run**: Windows may prompt to install **Wintun driver** — click "Yes" to install.

#### Linux

```bash
# Extract archive
tar -xzf JhopanStoreVPN-linux.tar.gz
cd JhopanStoreVPN-linux

# Run with sudo
sudo ./JhopanStoreVPN

# Or use the launcher script
sudo ./launch.sh
```

#### macOS

```bash
# Extract archive
tar -xzf JhopanStoreVPN-macos.tar.gz

# Run the .app bundle
open JhopanStoreVPN.app

# Or run from terminal with sudo
sudo ./JhopanStoreVPN.app/Contents/MacOS/JhopanStoreVPN
```

### 📱 Android Setup

1. **Enable Unknown Sources**
   - Go to Settings → Security → Unknown Sources
   - Or Settings → Apps → Special Access → Install Unknown Apps

2. **Install APK**
   - Download `JhopanStoreVPN.apk`
   - Tap to install
   - Grant installation permission

3. **First Connection**
   - Open app
   - Enter VLESS server details OR scan QR code
   - Tap **CONNECT**
   - Approve VPN permission request
   - (Optional) Enable **Hotspot Sharing** in settings

4. **Hotspot Sharing** (Optional)
   - Enable "Hotspot Sharing" toggle
   - Turn on WiFi hotspot
   - Other devices connect to your hotspot
   - They automatically use VPN (proxy on port 10809)

---

## 🎯 Usage Tips

### Desktop

**System Tray**

- Close button (X) → Hides to tray
- Right-click tray icon → Show window or Exit
- Exit properly to ensure routing cleanup

**Config Management**

- Clipboard button: paste `vless://` links
- Manual entry: server address, port, UUID, path
- Settings are saved automatically

**Monitoring**

- Ping display shows latency in real-time
- Green = Connected, Red = Disconnected
- Auto-reconnect on connection loss

### Android

**Battery Optimization**

- Disable battery optimization for JhopanStoreVPN
- This prevents system from killing VPN service
- Settings → Battery → Battery Optimization → JhopanStoreVPN → Don't Optimize

**Autostart Permission** (Xiaomi, Huawei, OPPO, etc.)

- Some manufacturers require autostart permission
- Settings → Apps → Autostart → Enable for JhopanStoreVPN

**Split Tunneling** (Coming Soon)

- Currently all apps route through VPN
- Split tunneling feature in future release

---

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

## 🏗️ Technical Architecture

### 🖥️ Desktop TUN Mode

```
┌─────────────────────────────────────────────────────────────┐
│  User Applications                                          │
│  (WhatsApp, Telegram, Discord, Chrome, Games, etc.)        │
└────────────────────┬────────────────────────────────────────┘
                     │ All network traffic
                     ↓
┌────────────────────────────────────────────────────────────┐
│  TUN Virtual Network Interface (tun0)                      │
│  IP: 10.0.0.2  Gateway: 10.0.0.1                          │
└────────────────────┬───────────────────────────────────────┘
                     │ Layer 3 packets
                     ↓
┌────────────────────────────────────────────────────────────┐
│  tun2socks Process (External Binary)                       │
│  Converts TUN packets → SOCKS5 protocol                    │
└────────────────────┬───────────────────────────────────────┘
                     │ SOCKS5 protocol
                     ↓
┌────────────────────────────────────────────────────────────┐
│  Xray-core Process (localhost:10809)                       │
│  VLESS Protocol + WebSocket + TLS                          │
└────────────────────┬───────────────────────────────────────┘
                     │ Encrypted VLESS
                     ↓
                 Internet → VLESS Server
```

**Technical Details:**

- **TUN Device**: Layer 3 virtual network interface created by tun2socks
- **Routing Table**: All traffic (0.0.0.0/0) routed via TUN gateway (10.0.0.1)
- **Platform Commands**:
  - Windows: `netsh interface ip set address`, `route add`
  - Linux: `ip link set up`, `ip route add default`
  - macOS: `ifconfig`, `route add`
- **Process Management**:
  - Go manages tun2socks lifecycle via `exec.Command`
  - Single instance lock: Named mutex (Windows) / flock (Linux/macOS)
  - Admin elevation: Manifest `requireAdministrator` (Windows)
- **Cleanup**: Explicit route deletion on disconnect ensures internet restoration

### 📱 Android VpnService

```
┌─────────────────────────────────────────────────────────────┐
│  User Applications                                          │
│  (All Android apps)                                         │
└────────────────────┬────────────────────────────────────────┘
                     │ All network traffic
                     ↓
┌────────────────────────────────────────────────────────────┐
│  Android VpnService API                                     │
│  TUN interface (tun0) managed by Android OS                │
└────────────────────┬───────────────────────────────────────┘
                     │ Packets intercepted by VPNService
                     ↓
┌────────────────────────────────────────────────────────────┐
│  libXray.aar (Native Library)                              │
│  Xray-core compiled for Android (ARM64 + x86_64)          │
│  VLESS Protocol + WebSocket + TLS                          │
└────────────────────┬───────────────────────────────────────┘
                     │ Encrypted VLESS
                     ↓
                 Internet → VLESS Server
```

**Additional Features:**

- **Hotspot Sharing**: HTTP proxy on port 10809 for WiFi clients
- **Foreground Service**: Notification prevents system kill
- **Config Storage**: Shared Preferences with encryption
- **Traffic Stats**: Real-time upload/download counter

---

## 🛠️ Troubleshooting

### 🖥️ Desktop Issues

<details>
<summary><b>"Failed to start tun2socks" or "Permission denied"</b></summary>

**Cause**: Application not running with administrator/root privileges.

**Solution**:

- **Windows**: Right-click → "Run as Administrator"
- **Linux/macOS**: Run with `sudo ./JhopanStoreVPN`
- **Note**: UAC manifest should trigger admin prompt automatically on Windows

</details>

<details>
<summary><b>"tun2socks.exe not found" or "Binary not found"</b></summary>

**Cause**: tun2socks binary missing from `bin/` folder.

**Solution**:

1. Download from [tun2socks releases](https://github.com/xjasonlyu/tun2socks/releases/tag/v2.5.2)
2. Extract the appropriate binary:
   - Windows: `tun2socks-windows-amd64.exe` → `bin/tun2socks.exe`
   - Linux: `tun2socks-linux-amd64` → `bin/tun2socks`
   - macOS: `tun2socks-darwin-amd64` → `bin/tun2socks`
3. Ensure `bin/` folder is next to the executable

</details>

<details>
<summary><b>Windows: "Wintun driver installation" prompt</b></summary>

**Cause**: Wintun driver (required for TUN mode) not installed.

**Solution**:

- Click "Yes" or "Install" when prompted
- This is a one-time installation
- Driver is required for TUN interface creation
- From: [Wintun by WireGuard](https://www.wintun.net/)

</details>

<details>
<summary><b>"Connection successful but no internet"</b></summary>

**Possible Causes**:

1. Xray process not running
2. VLESS server unreachable
3. Routing table misconfigured

**Diagnostics**:

```bash
# Windows: Check if Xray is running
tasklist | findstr xray

# Linux/macOS: Check Xray process
ps aux | grep xray

# Windows: Check routing table
route print 0.0.0.0

# Linux: Check routing
ip route show

# macOS: Check routing
netstat -nr | grep default
```

**Solution**:

- Verify VLESS server address, port, UUID are correct
- Check firewall isn't blocking Xray (port 10809)
- Try disconnecting and reconnecting
- Check Xray logs in application directory

</details>

<details>
<summary><b>"Application already running" or won't start second time</b></summary>

**Cause**: Single instance lock prevents duplicate processes (this is intentional).

**Solution**:

- Check system tray for running instance
- Right-click tray icon → Show window
- Or properly exit via tray menu → Exit
- If stuck: Kill process manually:
  - Windows: Task Manager → End JhopanStoreVPN.exe
  - Linux/macOS: `killall JhopanStoreVPN`

</details>

<details>
<summary><b>"Internet doesn't restore after disconnect"</b></summary>

**Cause**: Routing cleanup failed or tun2socks crashed.

**Solution**:

```bash
# Windows: Reset routing table
route -f
# Then reconnect to WiFi/Ethernet

# Linux: Flush routing cache
sudo ip route flush cache
sudo systemctl restart NetworkManager

# macOS: Reset network
sudo route flush
# System Preferences → Network → Assist me → Diagnostics
```

</details>

### 📱 Android Issues

<details>
<summary><b>"VPN disconnects randomly" or "Service killed"</b></summary>

**Cause**: System battery optimization or aggressive task killer.

**Solution**:

1. **Disable Battery Optimization**:
   - Settings → Battery → Battery Optimization
   - Find JhopanStoreVPN → Don't Optimize

2. **Enable Autostart** (Xiaomi, Huawei, OPPO, etc.):
   - Settings → Apps → Autostart
   - Enable for JhopanStoreVPN

3. **Lock in Recent Apps**:
   - Recent apps → JhopanStoreVPN → Lock icon

4. **Remove from Battery Saver**:
   - Some ROMs require explicit exemption from battery saver

</details>

<details>
<summary><b>"Hotspot Sharing not working"</b></summary>

**Cause**: Port 10809 blocked or hotspot restrictions.

**Solution**:

1. Ensure "Enable Hotspot Sharing" is toggled ON in settings
2. Check firewall isn't blocking port 10809
3. Test from client device:
   ```bash
   # Set proxy on client device
   Host: <Your hotspot IP>
   Port: 10809
   ```
4. Some carriers/ROMs restrict hotspot tethering
5. Try disabling and re-enabling hotspot

</details>

<details>
<summary><b>"Unable to install APK" or "App not installed"</b></summary>

**Solution**:

1. Enable "Install from Unknown Sources":
   - Settings → Security → Unknown Sources (Android 7-9)
   - Settings → Apps → Special Access → Install Unknown Apps → Chrome/File Manager (Android 10+)

2. Check storage space (needs ~50 MB free)

3. Uninstall old version first if upgrading:

   ```bash
   adb uninstall com.jhopanstore.vpn
   ```

4. If signature mismatch: completely uninstall old version

</details>

<details>
<summary><b>"Connection failed" or "Server unreachable"</b></summary>

**Diagnostics**:

1. Verify VLESS config: server address, port, UUID, path
2. Test server reachability from browser
3. Check server logs for connection attempts
4. Try importing config via QR code to avoid typos
5. Ensure server supports WebSocket + TLS

</details>

---

## ❓ FAQ

<details>
<summary><b>Why does desktop require administrator/root privileges?</b></summary>

TUN mode creates virtual network interfaces and modifies system routing tables. These are privileged operations that require elevated permissions.

**Alternatives without admin**:

- Proxy mode (not implemented, TUN-only by design)
- Browser extensions (limited to browser only)

</details>

<details>
<summary><b>Is my traffic encrypted?</b></summary>

Yes! All traffic is encrypted using:

- **VLESS protocol** with TLS (Transport Layer Security)
- **WebSocket** over TLS for additional obfuscation
- End-to-end encryption between client and VLESS server

</details>

<details>
<summary><b>Does JhopanStoreVPN log any data?</b></summary>

No. The application does not log:

- Websites visited
- Traffic content
- DNS queries
- Connection history

Local logs (if enabled) only contain:

- Connection status
- Error messages for debugging
- Performance metrics (latency)

</details>

<details>
<summary><b>Can I use JhopanStoreVPN for gaming?</b></summary>

Yes! TUN mode supports:

- ✅ Steam, Epic Games, Battle.net
- ✅ Multiplayer games (CS:GO, Dota 2, League of Legends)
- ✅ Xbox/PlayStation network connectivity
- ⚠️ Note: VPN adds latency, may affect competitive gaming

</details>

<details>
<summary><b>What's the difference between proxy mode and TUN mode?</b></summary>

| Feature             | Proxy Mode            | TUN Mode (JhopanStoreVPN) |
| ------------------- | --------------------- | ------------------------- |
| Requires admin      | No                    | Yes                       |
| Works with all apps | No (only proxy-aware) | Yes (system-wide)         |
| Gaming support      | No                    | Yes                       |
| Windows apps        | No                    | Yes                       |
| Configuration       | Manual per-app        | Automatic                 |

</details>

<details>
<summary><b>Can I run JhopanStoreVPN and other VPNs simultaneously?</b></summary>

No. Only one VPN/TUN interface can be active at a time. Attempting to run multiple VPNs will cause routing conflicts.

</details>

<details>
<summary><b>How do I check if VPN is working?</b></summary>

**Desktop**:

```bash
# Check your public IP
curl ifconfig.me

# Should show your VPN server IP, not your real IP
```

**Android**:

- Check notification: should show "Connected" with upload/download stats
- Visit: https://ifconfig.me in browser (should show VPN server IP)

</details>

<details>
<summary><b>What happens if VPN disconnects?</b></summary>

- **Desktop**: Auto-reconnect attempts every 5 seconds
- **Android**: Auto-reconnect (if Always-On VPN enabled)
- **Routing**: Automatically cleaned up, internet restored
- **Kill Switch**: Not implemented (traffic leaks during disconnection)

</details>

---

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
