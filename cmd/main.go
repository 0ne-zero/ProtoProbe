package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/0ne-zero/ProtoProbe/cmd/flags"
	"github.com/0ne-zero/ProtoProbe/config"
	"github.com/0ne-zero/ProtoProbe/protocols"
	"github.com/0ne-zero/ProtoProbe/protocols/dns"
)

func main() {
	opts, err := flags.ParseFlags()
	if err != nil {
		if !errors.Is(err, flags.ErrorNoFlags) {
			fmt.Println(err)
			os.Exit(1)
		}
		opts.All = true
	}

	cfg, err := config.LoadConfig(opts.ConfigFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if opts.All || opts.ICMP {
		printHeader("ICMP")
		runICMPTest(cfg)
	}

	if opts.All || opts.TCP {
		printHeader("TCP")
		runTCPTest(cfg)
	}

	if opts.All || opts.DoUDP {
		printHeader("Dns over UDP")
		runDnsOverUDPTest(cfg)
	}

	if opts.All || opts.DoTCP {
		printHeader("Dns over TCP")
		runDnsOverTCPTest(cfg)
	}

	if opts.All || opts.DoT {
		printHeader("DoT")
		runDoTTest(cfg)
	}

	if opts.All || opts.DoH {
		printHeader("DoH")
		runDoHTest(cfg)
	}

	if opts.All || opts.WebSocket {
		printHeader("WebSocket")
		runWebSocketTest(cfg)
	}
}

func runICMPTest(cfg *config.Config) {
	var icmpWg sync.WaitGroup
	for _, host := range cfg.ICMPHost {
		icmpWg.Add(1)
		go func(h string) {
			defer icmpWg.Done()
			res, err := protocols.TestICMP(h)
			if err != nil {
				log.Printf("[ICMP] | %s | %v ❌\n", h, err)
			} else {
				log.Printf("[ICMP] | %s | avg-rtt: %v | packet-loss: %.2f%% ✅\n", h, res.AvgRtt.Round(time.Millisecond), res.PacketLoss)
			}
		}(host)
	}
	icmpWg.Wait()
}

func runTCPTest(cfg *config.Config) {
	var tcpWg sync.WaitGroup
	for _, hostPort := range cfg.TCPHostPort {
		tcpWg.Add(1)
		go func(hp config.DNS_Host_Port_Query) {
			defer tcpWg.Done()
			res, err := protocols.TestTCP(hp)
			if err != nil {
				log.Printf("[TCP] | %s:%d | %v ❌\n", hp.Host, hp.Port, err)
			} else {
				log.Printf("[TCP] | %s:%d | rtt: %v ✅\n", hp.Host, hp.Port, res.RTT.Round(time.Millisecond))
			}
		}(hostPort)
	}
	tcpWg.Wait()
}

func runDnsOverUDPTest(cfg *config.Config) {
	var dnsOverUDPWg sync.WaitGroup
	for _, hostPort := range cfg.NormalDNSHostPort {
		dnsOverUDPWg.Add(1)
		go func(hp config.DNS_Host_Port_Query) {
			defer dnsOverUDPWg.Done()
			res, err := dns.TestDnsOverUDP(&hp)
			if err != nil {
				log.Printf("[DNS/UDP] | %s:%d | %v ❌\n", hp.Host, hp.Port, err)
			} else {
				log.Printf("[DNS/UDP] | %s:%d | rtt: %v ✅\n", hp.Host, hp.Port, res.RTT.Round(time.Millisecond))
			}
		}(hostPort)
	}
	dnsOverUDPWg.Wait()
}

func runDnsOverTCPTest(cfg *config.Config) {
	var dnsOverTCPWg sync.WaitGroup
	for _, hostPort := range cfg.NormalDNSHostPort {
		dnsOverTCPWg.Add(1)
		go func(hp config.DNS_Host_Port_Query) {
			defer dnsOverTCPWg.Done()
			res, err := dns.TestDNSTCP(&hp)
			if err != nil {
				log.Printf("[DNS/TCP] | %s:%d | %v ❌\n", hp.Host, hp.Port, err)
			} else {
				log.Printf("[DNS/TCP] | %s:%d | rtt: %v ✅\n", hp.Host, hp.Port, res.RTT.Round(time.Millisecond))
			}
		}(hostPort)
	}
	dnsOverTCPWg.Wait()
}

func runDoTTest(cfg *config.Config) {
	var dnsOverTlsWg sync.WaitGroup
	for _, hostPort := range cfg.DoT {
		dnsOverTlsWg.Add(1)
		go func(hp config.DNS_Host_Port_Query) {
			defer dnsOverTlsWg.Done()
			res, err := dns.TestDoT(&hp)
			if err != nil {
				log.Printf("[DNS/TLS (DoT)] | %s:%d | %v ❌\n", hp.Host, hp.Port, err)
			} else {
				log.Printf("[DNS/TLS (DoT)] | %s:%d | rtt: %v ✅\n", hp.Host, hp.Port, res.RTT.Round(time.Millisecond))
			}
		}(hostPort)
	}
	dnsOverTlsWg.Wait()
}

func runDoHTest(cfg *config.Config) {
	var dnsOverHttpsWg sync.WaitGroup
	for _, urlQuery := range cfg.DoH {
		dnsOverHttpsWg.Add(1)
		go func(urlQuery config.DNS_URL_Query) {
			defer dnsOverHttpsWg.Done()
			res, err := dns.TestDoH(&urlQuery)
			if err != nil {
				log.Printf("[DNS/HTTPS (DoH)] | %s | %v ❌\n", urlQuery.Addr, err)
			} else {
				log.Printf("[DNS/HTTPS (DoH)] | %s | rtt: %v ✅\n", urlQuery.Addr, res.RTT.Round(time.Millisecond))
			}
		}(urlQuery)
	}
	dnsOverHttpsWg.Wait()
}

func runWebSocketTest(cfg *config.Config) {
	var wsWg sync.WaitGroup
	for _, server := range cfg.WebSocket {
		wsWg.Add(1)
		go func(s string) {
			defer wsWg.Done()
			ws, err := protocols.TestWebSocket(s)
			if err != nil {
				log.Printf("[WebSocket] | %s | %v ❌\n", s, err)
			} else {
				fmt.Printf("[WebSocket] | %s | rtt: %v ✅\n", s, ws.RTT)
			}
		}(server)
	}
	wsWg.Wait()
}

func printHeader(proto string) {
	fmt.Println()
	fmt.Printf("%s:\n", proto)
}
