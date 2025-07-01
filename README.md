# ProtoProbe

ProtoProbe is a modular extendable Go tool to test and measure the connectivity of internet protocols.  
Designed for environments with network censorship, it helps check which protocols work and gather stats like RTT and packet loss.

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
ProtoProbe uses a JSON configuration file named `config.json` to define which hosts, ports, and queries to test.

---
## Usage

1. **Clone and build:**
```bash
git clone https://github.com/0ne-zero/ProtoProbe.git
cd ProtoProbe/cmd/
go build -o ProtoProbe .
./ProtoProbe
```
