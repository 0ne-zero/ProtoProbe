package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/0ne-zero/ProtoProbe/cmd/flags"
	"github.com/0ne-zero/ProtoProbe/config"
	"github.com/0ne-zero/ProtoProbe/protocols"
	"github.com/0ne-zero/ProtoProbe/protocols/dns"
)

func main() {
	opts, err := flags.ParseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
			icmp, err := protocols.TestICMP(h)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("[%s] avg-rtt: %v, packet-loss: %v\n", h, icmp.AvgRtt, icmp.PacketLoss)
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
			tcp, err := protocols.TestTCP(hp)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("[%s:%d] rtt: %v\n", hp.Host, hp.Port, tcp.RTT)
			}
		}(hostPort)
	}
	tcpWg.Wait()
}

func runDnsOverUDPTest(cfg *config.Config) {
	var dnsOverUDPWg sync.WaitGroup
	for _, server := range cfg.NormalDNSHostPort {
		dnsOverUDPWg.Add(1)
		go func(s config.DNS_Host_Port_Query) {
			defer dnsOverUDPWg.Done()
			dns, err := dns.TestDnsOverUDP(&s)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("[%s:%d] rtt: %v\n", s.Host, s.Port, dns.RTT)
			}
		}(server)
	}
	dnsOverUDPWg.Wait()
}

func runDnsOverTCPTest(cfg *config.Config) {
	var dnsOverTCPWg sync.WaitGroup
	for _, server := range cfg.NormalDNSHostPort {
		dnsOverTCPWg.Add(1)
		go func(s config.DNS_Host_Port_Query) {
			defer dnsOverTCPWg.Done()
			dns, err := dns.TestDNSTCP(&s)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("[%s:%d] rtt: %v\n", s.Host, s.Port, dns.RTT)
			}
		}(server)
	}
	dnsOverTCPWg.Wait()
}

func runDoTTest(cfg *config.Config) {
	var dnsOverTlsWg sync.WaitGroup
	for _, server := range cfg.DoT {
		dnsOverTlsWg.Add(1)
		go func(s config.DNS_Host_Port_Query) {
			defer dnsOverTlsWg.Done()
			dns, err := dns.TestDoT(&s)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("[%s:%d] rtt: %v\n", s.Host, s.Port, dns.RTT)
			}
		}(server)
	}
	dnsOverTlsWg.Wait()
}

func runDoHTest(cfg *config.Config) {
	var dnsOverHttpsWg sync.WaitGroup
	for _, server := range cfg.DoH {
		dnsOverHttpsWg.Add(1)
		go func(s config.DNS_URL_Query) {
			defer dnsOverHttpsWg.Done()
			dns, err := dns.TestDoH(&s)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("[%s] rtt: %v\n", s, dns.RTT)
			}
		}(server)
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
				fmt.Println(err)
			} else {
				fmt.Printf("[%s] rtt: %v\n", s, ws.RTT)
			}
		}(server)
	}
	wsWg.Wait()
}

func printHeader(proto string) {
	fmt.Println()
	fmt.Printf("%s:\n", proto)
}
