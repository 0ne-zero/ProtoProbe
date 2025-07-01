package flags

import (
	"flag"
	"fmt"
)

const DefaultConfigFilePath = "config.json"

type Options struct {
	ConfigFilePath string
	All            bool
	ICMP           bool
	TCP            bool
	DoUDP          bool // DNS over UDP
	DoTCP          bool // DNS over TCP
	DoT            bool
	DoH            bool
	WebSocket      bool
}

func ParseFlags() (*Options, error) {
	var opts Options
	flag.StringVar(&opts.ConfigFilePath, "config", DefaultConfigFilePath, "Path to config file")
	flag.BoolVar(&opts.All, "all", false, "Test all protocols")
	flag.BoolVar(&opts.ICMP, "icmp", false, "Test ICMP")
	flag.BoolVar(&opts.TCP, "tcp", false, "Test TCP")
	flag.BoolVar(&opts.DoUDP, "dou", false, "Test DNS over UDP")
	flag.BoolVar(&opts.DoTCP, "dotcp", false, "Test DNS over TCP")
	flag.BoolVar(&opts.DoT, "dot", false, "Test DNS over TLS")
	flag.BoolVar(&opts.DoH, "doh", false, "Test DNS over HTTPS")
	flag.BoolVar(&opts.WebSocket, "websocket", false, "Test WebSocket")

	flag.Parse()

	// Validate: -a cannot be used with other protocol flags
	if opts.All && (opts.ICMP || opts.TCP || opts.DoUDP || opts.DoTCP || opts.DoT || opts.DoH || opts.WebSocket) {
		return nil, fmt.Errorf("error: -a cannot be used together with individual protocol flags")
	}

	return &opts, nil
}
