# ProtoProbe

ProtoProbe is a modular Go tool to test and measure the connectivity and performance of key internet protocols.  
Designed for environments with network censorship, it helps check which protocols work and gather stats like RTT and packet loss.

---

## Features
- ICMP (ping)
- TCP
- UDP
- DNS over UDP (Plain DNS)
- DNS over TLS (DoT)
- DNS over HTTPS (DoH)
- WebSocket

Each test is modular and easy to extend.

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