# JhopanStoreVPN

<p align="center">
  <img src="assets/icon256.png" alt="JhopanStoreVPN Logo" width="128">
</p>

<p align="center">
  Lightweight VLESS VPN desktop client built with Go + Fyne + Xray-core.
</p>

<p align="center">
  <b>Developed by <a href="https://github.com/jhopan">JhopanStore</a></b>
</p>

---

## Features

- VLESS over WebSocket + TLS
- Clean light GUI with system tray
- Auto system proxy (Windows / Linux / macOS)
- Clipboard import (`vless://` URI)
- HTTP ping monitoring
- Auto-reconnect on disconnect
- No terminal / console window — runs as a native GUI app on all platforms

## Download

Go to [**Releases**](https://github.com/jhopan/JhopanStoreVPN/releases) to download:

| Platform | File | How to Run |
|----------|------|------------|
| **Windows** | `JhopanStoreVPN-windows.zip` | Extract → double-click `JhopanStoreVPN.exe` |
| **Linux** | `JhopanStoreVPN-linux.tar.gz` | Extract → double-click `JhopanStoreVPN.desktop` or run `./launch.sh` |
| **macOS** | `JhopanStoreVPN-macos.tar.gz` | Extract → double-click `JhopanStoreVPN.app` |

## Usage

1. Download and extract the archive for your OS
2. Run the application (see table above)
3. Enter server address and UUID, or paste a `vless://` link via the clipboard button
4. Click **CONNECT**

## Build from Source

### Requirements
- Go 1.21+
- GCC (MinGW on Windows)
- Platform graphics libraries:
  - **Linux**: `sudo apt install gcc pkg-config libgl1-mesa-dev xorg-dev`
  - **macOS**: `xcode-select --install`
  - **Windows**: MinGW-w64

### Build

```bash
# Linux / macOS
go build -ldflags="-s -w" -o JhopanStoreVPN .

# Windows (no console window)
go build -ldflags="-s -w -H windowsgui" -o JhopanStoreVPN.exe .
```

Place the `xray` binary next to the executable. Download from [Xray-core Releases](https://github.com/XTLS/Xray-core/releases).

## Credits

- **[JhopanStore](https://github.com/jhopan)** — Developer & maintainer
- **[Xray-core](https://github.com/XTLS/Xray-core)** — VLESS/VMess/Trojan proxy engine (XTLS Project, MPL-2.0 license)
- **[Fyne](https://github.com/fyne-io/fyne)** — Cross-platform GUI toolkit for Go (BSD 3-Clause license)
- **[golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys)** — Go extended system packages (BSD 3-Clause license)

## License

MIT — Copyright (c) 2026 JhopanStore