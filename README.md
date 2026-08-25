# Clean IP Scanner

Find the fastest clean Cloudflare IPs for your network.

## Features

- **Normal Mode**: TCP ping + download/upload speed test
- **Xray Mode**: Protocol-aware testing through VLESS/VMess/Trojan/Shadowsocks tunnels
- **Checkpoint**: Resume interrupted scans
- **Progress Bar**: Real-time progress with EWMA speed calculation

## Requirements

### Normal Mode
- Go 1.21+
- Network access

### Xray Mode (Optional)
- Xray binary at `./xray/xray`
- Config at `config/xray_config.txt` (URL format) or `config/xray_config.json`

## Build

```bash
# Install Go (iSH/Alpine)
apk add go git build-base

# Clone and build
git clone https://github.com/surfrpt1/clean-ip-scanner.git
cd clean-ip-scanner
go build -o clean-ip-scanner .
```

## Usage

```bash
# Basic scan
./clean-ip-scanner -save

# Download only test
./clean-ip-scanner -mode 1 -n 10 -save

# With Xray mode
./clean-ip-scanner -xray -save

# Custom workers
./clean-ip-scanner -workers 100 -save

# Reset checkpoint and start fresh
./clean-ip-scanner -reset -save
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-mode` | Scan mode: 0=full, 1=download, 2=upload | 0 |
| `-workers` | Concurrent workers | 200 |
| `-test-count` | IPs to speed test | 5 |
| `-n` | Shorthand for -test-count | 5 |
| `-seed` | Random seed (0=current time) | 0 |
| `-save` | Save results to files | false |
| `-reset` | Reset checkpoint | false |
| `-xray` | Use Xray mode | false |

## Output

Results are saved to:
- `results.txt` - Full results with speeds
- `simple_results.txt` - Simple IP list

## Disclaimer

This software is provided for **educational and research purposes only**. Use it at your own risk. The authors are not responsible for any misuse, legal consequences, or damages arising from the use of this tool. Users are solely responsible for ensuring their use complies with all applicable laws and regulations in their jurisdiction. By using this software, you agree that you will only use it in a lawful manner and accept all responsibility for your actions.

## License

MIT
