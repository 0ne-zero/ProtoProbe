package flags

import (
	"flag"
	"fmt"
	"os"
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
	ECH            bool
	DoUDP          bool // DNS over UDP
	DoTCP          bool // DNS over TCP
	DoT            bool
	DoQ            bool
	DoH            bool
	HTTP           bool
	HTTPS          bool
	WebSocket      bool
	QUIC           bool
	STUN           bool
	NTP            bool
	DoTInsecure    bool // skip TLS certificate verification for DoT
	TLSInsecure    bool // skip TLS certificate verification for TLS
	QUICInsecure   bool // skip TLS certificate verification for QUIC
	JSON           bool // emit results as JSON instead of human-readable text
}

// ParseFlags parses command-line flags from os.Args[1:].
func ParseFlags() (Options, error) {
	return ParseFlagsFrom(os.Args[1:])
}

// ParseFlagsFrom parses flags from the given args slice using a fresh FlagSet.
// This is the testable core; ParseFlags delegates to it with os.Args[1:].
func ParseFlagsFrom(args []string) (Options, error) {
	var opts Options
	fs := flag.NewFlagSet("protoprobe", flag.ContinueOnError)

	fs.StringVar(&opts.ConfigFilePath, "config", DefaultConfigFilePath, "Path to config file")
	fs.BoolVar(&opts.All, "all", false, "Test all protocols")
	fs.BoolVar(&opts.ICMP, "icmp", false, "Test ICMP")
	fs.BoolVar(&opts.TCP, "tcp", false, "Test TCP")
	fs.BoolVar(&opts.TLS, "tls", false, "Test TLS handshake")
	fs.BoolVar(&opts.TLSInsecure, "tls-insecure", false, "Skip TLS certificate verification for TLS")
	fs.BoolVar(&opts.ECH, "ech", false, "Test TLS with Encrypted Client Hello")
	fs.BoolVar(&opts.DoUDP, "dou", false, "Test DNS over UDP")
	fs.BoolVar(&opts.DoTCP, "dotcp", false, "Test DNS over TCP")
	fs.BoolVar(&opts.DoT, "dot", false, "Test DNS over TLS")
	fs.BoolVar(&opts.DoTInsecure, "dot-insecure", false, "Skip TLS certificate verification for DoT")
	fs.BoolVar(&opts.DoQ, "doq", false, "Test DNS over QUIC")
	fs.BoolVar(&opts.DoH, "doh", false, "Test DNS over HTTPS")
	fs.BoolVar(&opts.HTTP, "http", false, "Test HTTP")
	fs.BoolVar(&opts.HTTPS, "https", false, "Test HTTPS")
	fs.BoolVar(&opts.WebSocket, "websocket", false, "Test WebSocket")
	fs.BoolVar(&opts.QUIC, "quic", false, "Test QUIC (HTTP/3 handshake)")
	fs.BoolVar(&opts.QUICInsecure, "quic-insecure", false, "Skip TLS certificate verification for QUIC")
	fs.BoolVar(&opts.STUN, "stun", false, "Test STUN (NAT binding)")
	fs.BoolVar(&opts.NTP, "ntp", false, "Test NTP")
	fs.BoolVar(&opts.JSON, "json", false, "Output results as JSON")

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	// Validate: -all cannot be used with other protocol flags
	if opts.All && (opts.ICMP || opts.TCP || opts.TLS || opts.ECH || opts.DoUDP || opts.DoTCP || opts.DoT || opts.DoQ || opts.DoH || opts.HTTP || opts.HTTPS || opts.WebSocket || opts.QUIC || opts.STUN || opts.NTP) {
		return opts, errorALLFlagConflict
	}

	// Validate: at least one protocol flag must be present
	if !opts.All && !opts.ICMP && !opts.TCP && !opts.TLS && !opts.ECH && !opts.DoUDP && !opts.DoTCP && !opts.DoT && !opts.DoQ && !opts.DoH && !opts.HTTP && !opts.HTTPS && !opts.WebSocket && !opts.QUIC && !opts.STUN && !opts.NTP {
		return opts, ErrorNoFlags
	}

	return opts, nil
}
