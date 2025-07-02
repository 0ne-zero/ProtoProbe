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