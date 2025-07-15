# ProtoProbe

ProtoProbe is a modular and extensible Go tool designed to test and measure the connectivity of internet protocols.
Built specifically for environments affected by network censorship, it helps identify which protocols are functional and collects statistics such as RTT and packet loss.

---

## Features
- ICMP (ping)
- TCP
- DNS over UDP (Normal DNS)
- DNS over TCP
- DNS over TLS (DoT)
- DNS over HTTPS (DoH)
- WebSocket

---

## Configuration
ProtoProbe uses a JSON configuration file named [`config.json`](https://github.com/0ne-zero/ProtoProbe/blob/main/cmd/config.json) to define which hosts, ports, and queries to test.

---
## Downlaod/Build
You can either clone and build the tool from source and view the code, or download the binary release.

- **Clone and build:**

```bash
git clone https://github.com/0ne-zero/ProtoProbe.git
cd ProtoProbe/cmd/
go build -o protoprobe .
./protoprobe
```
- **Download binary release:**
Download the appropriate binary for your OS and architecture from the [releases](https://github.com/0ne-zero/ProtoProbe/releases) page. Then, either place the config.json file next to the executable or provide its path using the -config flag.
```bash
./protoprobe -config /path/to/config.json
```

---
## Help
```
Usage of protoprobe:
  -all
        Test all protocols
  -config string
        Path to config file (default "config.json")
  -doh
        Test DNS over HTTPS
  -dot
        Test DNS over TLS
  -dotcp
        Test DNS over TCP
  -dou
        Test DNS over UDP
  -icmp
        Test ICMP
  -tcp
        Test TCP
  -websocket
        Test WebSocket
```
---
## Example output
```bash
./protoprobe -all 

ICMP:
2025/07/15 20:22:57 [ICMP] | 8.8.8.8 | avg-rtt: 75ms | packet-loss: 0.00% ✅
2025/07/15 20:22:57 [ICMP] | 1.1.1.1 | avg-rtt: 135ms | packet-loss: 0.00% ✅

TCP:
2025/07/15 20:22:58 [TCP] | google.com:443 | rtt: 273ms ✅

Dns over UDP:
2025/07/15 20:22:58 [DNS/UDP] | 8.8.8.8:53 | rtt: 68ms ✅
2025/07/15 20:22:58 [DNS/UDP] | 1.1.1.1:53 | rtt: 141ms ✅

Dns over TCP:
2025/07/15 20:23:00 [DNS/TCP] | 8.8.8.8:53 | EOF ❌
2025/07/15 20:23:06 [DNS/TCP] | 1.1.1.1:53 | read tcp 192.168.1.2:46330->1.1.1.1:53: i/o timeout ❌

DoT:
2025/07/15 20:23:14 [DNS/TLS (DoT)] | 1.1.1.1:853 | context deadline exceeded ❌

DoH:
2025/07/15 20:23:15 [DNS/HTTPS (DoH)] | https://cloudflare-dns.com/dns-query | rtt: 723ms ✅

WebSocket:
[WebSocket] | wss://echo.websocket.events | rtt: 921.0533ms ✅
```
