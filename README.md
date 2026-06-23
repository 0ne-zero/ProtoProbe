# ProtoProbe

ProtoProbe is a modular Go tool for testing internet protocol connectivity. Built for environments affected by network censorship, it measures which protocols are reachable and reports RTT, packet loss, and HTTP status codes.

---

## Supported protocols

| Flag | Protocol | Notes |
|---|---|---|
| `-icmp` | ICMP (ping) | Requires root / `CAP_NET_RAW` on Linux |
| `-dou` | DNS over UDP | Standard DNS |
| `-dotcp` | DNS over TCP | DNS fallback transport |
| `-tcp` | TCP | Generic transport connectivity |
| `-tls` | TLS handshake | Detects SNI-based blocking |
| `-dot` | DNS over TLS (DoT) | |
| `-doh` | DNS over HTTPS (DoH) | |
| `-http` | HTTP | |
| `-https` | HTTPS | |
| `-quic` | QUIC (HTTP/3 handshake) | |
| `-websocket` | WebSocket | |

---

## Build

**From source:**
```bash
git clone https://github.com/0ne-zero/ProtoProbe.git
cd ProtoProbe
go build -o protoprobe ./cmd/protoprobe
./protoprobe -all
```

**Pre-built binaries:** Download the appropriate binary for your OS and architecture from the [Releases](https://github.com/0ne-zero/ProtoProbe/releases) page.

---

## Configuration

ProtoProbe reads targets from a JSON config file (default: `config.json` in the working directory). Pass a custom path with `-config`.

```json
{
    "icmp": ["8.8.8.8", "1.1.1.1"],
    "dns": [
        {"host": "1.1.1.1", "port": 53, "query": "yahoo.com"},
        {"host": "8.8.8.8", "port": 53, "query": "yahoo.com"}
    ],
    "tcp": [
        {"host": "google.com", "port": 443}
    ],
    "tls": [
        {"host": "www.google.com", "port": 443}
    ],
    "dot": [
        {"host": "1.1.1.1", "port": 853, "query": "yahoo.com"}
    ],
    "doh": [
        {"address": "https://cloudflare-dns.com/dns-query", "query": "yahoo.com"}
    ],
    "http":  ["http://neverssl.com", "http://example.com"],
    "https": ["https://www.google.com"],
    "quic": [
        {"host": "cloudflare-quic.com", "port": 443},
        {"host": "www.google.com",      "port": 443}
    ],
    "websocket": ["wss://ws.postman-echo.com/raw"]
}
```

---

## Flags

```
Usage of protoprobe:
  -all
        Test all protocols
  -icmp
        Test ICMP
  -dou
        Test DNS over UDP
  -dotcp
        Test DNS over TCP
  -tcp
        Test TCP
  -tls
        Test TLS handshake
  -tls-insecure
        Skip TLS certificate verification for TLS
  -dot
        Test DNS over TLS
  -dot-insecure
        Skip TLS certificate verification for DoT
  -doh
        Test DNS over HTTPS
  -http
        Test HTTP
  -https
        Test HTTPS
  -quic
        Test QUIC (HTTP/3 handshake)
  -quic-insecure
        Skip TLS certificate verification for QUIC
  -websocket
        Test WebSocket
  -json
        Output results as JSON
  -config string
        Path to config file (default "config.json")
```

---

## Example output

**Human-readable (default):**
```
./protoprobe -all

ICMP:
2025/07/15 20:22:57 [ICMP] | 8.8.8.8 | avg-rtt: 75ms | packet-loss: 0.00% ✅
2025/07/15 20:22:57 [ICMP] | 1.1.1.1 | avg-rtt: 135ms | packet-loss: 0.00% ✅

Dns over UDP:
2025/07/15 20:22:58 [DNS/UDP] | 8.8.8.8:53 | rtt: 68ms ✅
2025/07/15 20:22:58 [DNS/UDP] | 1.1.1.1:53 | rtt: 141ms ✅

Dns over TCP:
2025/07/15 20:23:00 [DNS/TCP] | 8.8.8.8:53 | EOF ❌
2025/07/15 20:23:06 [DNS/TCP] | 1.1.1.1:53 | i/o timeout ❌

TCP:
2025/07/15 20:22:58 [TCP] | google.com:443 | rtt: 273ms ✅

TLS:
2025/07/15 20:22:59 [TLS] | www.google.com:443 | rtt: 310ms ✅

DoT:
2025/07/15 20:23:14 [DNS/TLS (DoT)] | 1.1.1.1:853 | context deadline exceeded ❌

DoH:
2025/07/15 20:23:15 [DNS/HTTPS (DoH)] | https://cloudflare-dns.com/dns-query | rtt: 723ms ✅

HTTP:
2025/07/15 20:23:15 [HTTP] | http://neverssl.com | rtt: 210ms | status: 200 ✅

HTTPS:
2025/07/15 20:23:16 [HTTPS] | https://www.google.com | rtt: 289ms | status: 200 ✅

QUIC:
2025/07/15 20:23:16 [QUIC] | cloudflare-quic.com:443 | rtt: 180ms ✅

WebSocket:
2025/07/15 20:23:17 [WebSocket] | wss://ws.postman-echo.com/raw | rtt: 921ms ✅
```

**JSON (`-json`):**
```bash
./protoprobe -icmp -tcp -json
```
```json
{
  "results": [
    {
      "protocol": "ICMP",
      "target": "8.8.8.8",
      "success": true,
      "rtt_ms": 75,
      "packet_loss": 0
    },
    {
      "protocol": "TCP",
      "target": "google.com:443",
      "success": true,
      "rtt_ms": 273
    },
    {
      "protocol": "TCP",
      "target": "arvancloud.ir:443",
      "success": false,
      "error": "dial tcp: connection refused"
    }
  ]
}
```

---

## Notes

- **ICMP** requires root privileges or `CAP_NET_RAW` on Linux. The test will fail with a permission error without them.
- **`-tls-insecure`**, **`-dot-insecure`**, and **`-quic-insecure`** skip TLS certificate verification for their respective protocols. Useful when probing servers with self-signed certificates or in environments where the certificate chain is intercepted.
- **HTTP and HTTPS** do not follow redirects — a `3xx` response is reported as-is so a redirect to a block page does not mask censorship.
- Results within each protocol group run **concurrently**; protocol groups run sequentially in the order shown in the flags table above.
