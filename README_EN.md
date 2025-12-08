# Latency Testing Utility

A dual-implementation HTTP latency measurement tool available in both Go and Python, designed to measure response times while emulating a Chrome browser to bypass anti-bot protection.

## Features

- 🌐 Chrome browser emulation to bypass bot detection
- 🔥 Warm-up request to exclude cold start effects
- 📊 Comprehensive statistics (average, min, max, 95th/99th percentiles)
- ⚙️ Configurable test count and delay between requests

## Quick Start

### Go Version

```bash
# Install dependencies
go mod tidy

# Build
go build -o ping_test.exe test_ping.go

# Run
./ping_test.exe -url "https://example.com" -num_tests 5 -delay 1.0
```

**Using Pre-built Linux Binary:**

A ready-to-use compiled binary for Linux (AMD64) is available: `ping_test_linux`

```bash
# 1. Make it executable
chmod +x ping_test_linux

# 2. Run
./ping_test_linux -url "https://example.com" -num_tests 5 -delay 1.0
```

> **Note:** The binary is statically compiled and doesn't require Go or any dependencies on the Linux system.

### Python Version

```bash
# Install dependencies
pip install curl_cffi

# Run
python test_ping.py "https://example.com" --num_tests 5 --delay 1.0
```

## Command Line Arguments

**Go version:**
- `-url` (required): Target URL to test
- `-num_tests` (optional, default: 5): Number of tests to run
- `-delay` (optional, default: 1.0): Delay in seconds between requests

**Python version:**
- `url` (required): Target URL to test
- `--num_tests` (optional, default: 5): Number of tests to run
- `--delay` (optional, default: 1.0): Delay in seconds between requests

## Example Output

```
Performing warm-up request...
Test 1: Latency = 245.32 ms
Test 2: Latency = 198.45 ms
Test 3: Latency = 203.12 ms
Test 4: Latency = 215.67 ms
Test 5: Latency = 189.23 ms

Average latency: 210.36 ms
Min latency: 189.23 ms
Max latency: 245.32 ms
95th percentile latency: 241.09 ms
99th percentile latency: 244.48 ms
```

## Full Documentation

For complete documentation in Russian, see [README.md](README.md)
