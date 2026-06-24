# ProtoProbe

ProtoProbe is a modular Go tool for testing internet protocol connectivity. Built for environments affected by network censorship, it measures which protocols are reachable and reports RTT, packet loss, and HTTP status codes.

---

## Example output

**Human-readable (default):**
```
./protoprobe -all
PROTOCOL         TARGET                                RTT      RESULT
───────────────  ────────────────────────────────────  ───────  ──────
ICMP             1.1.1.1                               34ms     ✅  loss: 0.00%
ICMP             8.8.8.8                               34ms     ✅  loss: 0.00%
ICMP             185.143.234.200                       62ms     ✅  loss: 0.00%
DNS/UDP          1.1.1.1:53                            62ms     ✅
DNS/UDP          8.8.8.8:53                            66ms     ✅
DNS/TCP          8.8.8.8:53                            111ms    ✅
DNS/TCP          1.1.1.1:53                            118ms    ✅
DNS/TLS (DoT)    1.1.1.1:853                           195ms    ✅
DNS/QUIC (DoQ)   dns.adguard.com:8853                  2226ms   ✅
DNS/HTTPS (DoH)  https://cloudflare-dns.com/dns-query  212ms    ✅
TCP              arvancloud.ir:443                     46ms     ✅
TCP              google.com:443                        49ms     ✅
HTTP             http://example.com                    167ms    ✅  status: 200
HTTP             http://httpforever.com                478ms    ✅  status: 200
HTTPS            https://www.arvancloud.ir             250ms    ✅  status: 200
HTTPS            https://www.google.com                554ms    ✅  status: 200
TLS              www.arvancloud.ir:443                 191ms    ✅
TLS              www.google.com:443                    200ms    ✅
ECH              crypto.cloudflare.com:443             185ms    ✅  ech: accepted
ECH              tls-ech.dev:443                       453ms    ✅  ech: accepted
QUIC             cloudflare-quic.com:443               154ms    ✅
QUIC             www.google.com:443                    170ms    ✅
QUIC             digikala.ir:443                       -        ❌  timeout: no recent network activity
WebSocket        wss://ws.postman-echo.com/raw         749ms    ✅
STUN             stun.cloudflare.com:3478              69ms     ✅
STUN             stun.l.google.com:19302               79ms     ✅
NTP              time.cloudflare.com:123               88ms     ✅
NTP              time.google.com:123                   111ms    ✅
```

**JSON (`-json`):**
```bash
./protoprobe -icmp -tcp -ech -json
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
    },
    {
      "protocol": "ECH",
      "target": "crypto.cloudflare.com:443",
      "success": true,
      "rtt_ms": 290,
      "ech_accepted": true
    }
  ]
}
```

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
    "icmp": [
        "8.8.8.8",
        "1.1.1.1",
        "185.143.234.200"
    ],
    "dns": [
        {
            "host": "1.1.1.1",
            "port": 53,
            "query": "yahoo.com"
        },
        {
            "host": "8.8.8.8",
            "port": 53,
            "query": "arvancloud.ir"
        }
    ],
    "tcp": [
        {
            "host": "google.com",
            "port": 443
        },
        {
            "host": "arvancloud.ir",
            "port": 443
        }
    ],
    "dot": [
        {
            "host": "1.1.1.1",
            "port": 853,
            "query": "yahoo.com"
        }
    ],
    "doq": [
        {
            "host": "dns.adguard.com",
            "port": 8853,
            "query": "yahoo.com"
        }
    ],
    "doh": [
        {
            "address": "https://cloudflare-dns.com/dns-query",
            "query": "yahoo.com"
        }
    ],
    "http": [
        "http://httpforever.com",
        "http://example.com"
    ],
    "https": [
        "https://www.google.com",
        "https://www.arvancloud.ir"
    ],
     "tls": [
        {
            "host": "www.google.com",
            "port": 443
        },
        {
            "host": "www.arvancloud.ir",
            "port": 443
        }
    ],
    "ech": [
        {
            "host": "crypto.cloudflare.com",
            "port": 443
        },
        {
            "host": "tls-ech.dev",
            "port": 443
        }
    ],
    "quic": [
        {
            "host": "cloudflare-quic.com",
            "port": 443
        },
        {
            "host": "www.google.com",
            "port": 443
        },
        {
            "host": "digikala.ir",
            "port": 443
        }
    ],
    "websocket": [
        "wss://ws.postman-echo.com/raw"
    ],
    "stun": [
        {
            "host": "stun.l.google.com",
            "port": 19302
        },
        {
            "host": "stun.cloudflare.com",
            "port": 3478
        }
    ],
    "ntp": [
        {
            "host": "time.cloudflare.com",
            "port": 123
        },
        {
            "host": "time.google.com",
            "port": 123
        }
    ]
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
  -ech
        Test TLS with Encrypted Client Hello
  -dot
        Test DNS over TLS
  -dot-insecure
        Skip TLS certificate verification for DoT
  -doq
        Test DNS over QUIC
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
  -stun
        Test STUN (NAT binding)
  -ntp
        Test NTP
  -json
        Output results as JSON
  -config string
        Path to config file (default "config.json")
```

---

## Notes

- **ICMP** requires root privileges or `CAP_NET_RAW` on Linux. The test will fail with a permission error without them.
- **`-tls-insecure`**, **`-dot-insecure`**, and **`-quic-insecure`** skip TLS certificate verification for their respective protocols. Useful when probing servers with self-signed certificates or in environments where the certificate chain is intercepted.
- **HTTP and HTTPS** do not follow redirects — a `3xx` response is reported as-is so a redirect to a block page does not mask censorship.
- All protocol groups run **concurrently**; output is printed in the order shown in the flags table above so results are always easy to compare across runs.
