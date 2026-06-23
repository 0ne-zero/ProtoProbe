package flags

import (
	"flag"
	"fmt"
)

const DefaultConfigFilePath = "config.json"

var (
	ErrorNoFlags         = fmt.Errorf("error: there should be at least one flag")
	errorALLFlagConflict = fmt.Errorf("error: -a cannot be used together with individual protocol flags")
)

type Options struct {
	ConfigFilePath string
	All            bool
	ICMP           bool
	TCP            bool
	TLS            bool
	DoUDP          bool // DNS over UDP
	DoTCP          bool // DNS over TCP
	DoT            bool
	DoH            bool
	HTTP           bool
	HTTPS          bool
	WebSocket      bool
	QUIC           bool
	DoTInsecure    bool // skip TLS certificate verification for DoT
	TLSInsecure    bool // skip TLS certificate verification for TLS
	QUICInsecure   bool // skip TLS certificate verification for QUIC
	JSON           bool // emit results as JSON instead of human-readable text
}

func ParseFlags() (Options, error) {
	var opts Options
	flag.StringVar(&opts.ConfigFilePath, "config", DefaultConfigFilePath, "Path to config file")
	flag.BoolVar(&opts.All, "all", false, "Test all protocols")
	flag.BoolVar(&opts.ICMP, "icmp", false, "Test ICMP")
	flag.BoolVar(&opts.TCP, "tcp", false, "Test TCP")
	flag.BoolVar(&opts.TLS, "tls", false, "Test TLS handshake")
	flag.BoolVar(&opts.TLSInsecure, "tls-insecure", false, "Skip TLS certificate verification for TLS")
	flag.BoolVar(&opts.DoUDP, "dou", false, "Test DNS over UDP")
	flag.BoolVar(&opts.DoTCP, "dotcp", false, "Test DNS over TCP")
	flag.BoolVar(&opts.DoT, "dot", false, "Test DNS over TLS")
	flag.BoolVar(&opts.DoTInsecure, "dot-insecure", false, "Skip TLS certificate verification for DoT")
	flag.BoolVar(&opts.DoH, "doh", false, "Test DNS over HTTPS")
	flag.BoolVar(&opts.HTTP, "http", false, "Test HTTP")
	flag.BoolVar(&opts.HTTPS, "https", false, "Test HTTPS")
	flag.BoolVar(&opts.WebSocket, "websocket", false, "Test WebSocket")
	flag.BoolVar(&opts.QUIC, "quic", false, "Test QUIC (HTTP/3 handshake)")
	flag.BoolVar(&opts.QUICInsecure, "quic-insecure", false, "Skip TLS certificate verification for QUIC")
	flag.BoolVar(&opts.JSON, "json", false, "Output results as JSON")

	flag.Parse()

	// Validate: -all cannot be used with other protocol flags
	if opts.All && (opts.ICMP || opts.TCP || opts.TLS || opts.DoUDP || opts.DoTCP || opts.DoT || opts.DoH || opts.HTTP || opts.HTTPS || opts.WebSocket || opts.QUIC) {
		return opts, errorALLFlagConflict
	}

	// Validate: at least one protocol flag must be present
	if !opts.All && !opts.ICMP && !opts.TCP && !opts.TLS && !opts.DoUDP && !opts.DoTCP && !opts.DoT && !opts.DoH && !opts.HTTP && !opts.HTTPS && !opts.WebSocket && !opts.QUIC {
		return opts, ErrorNoFlags
	}

	return opts, nil
}
