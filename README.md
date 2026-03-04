# JhopanStoreVPN

Lightweight VLESS VPN desktop client built with Go + Fyne + Xray-core.

## Features
- VLESS over WebSocket + TLS
- Clean light GUI with system tray
- Auto system proxy (Windows/Linux/macOS)
- Clipboard import (`vless://` URI)
- HTTP ping monitoring
- Auto-reconnect on disconnect

## Download

Go to [Releases](https://github.com/jhopan/JhopanStoreVPN/releases) to download pre-built binaries for:
- **Windows** (`.zip`)
- **Linux** (`.tar.gz`)
- **macOS** (`.tar.gz`)

## Usage

1. Extract the archive
2. Run `JhopanStoreVPN` (or `.exe` on Windows)
3. Enter server address and UUID, or paste a `vless://` link via clipboard button
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
go build -ldflags="-s -w" -o JhopanStoreVPN .
```

Place `xray` binary ([Xray-core releases](https://github.com/XTLS/Xray-core/releases)) next to the executable.

## License

MIT