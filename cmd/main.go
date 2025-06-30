package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/0ne-zero/ProtoProbe/config"
	"github.com/0ne-zero/ProtoProbe/protocols"
	"github.com/0ne-zero/ProtoProbe/protocols/dns"
)

const configPath = "config.json"

func main() {
	//fmt.Println("Loading config file...")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printHeader("ICMP")
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

	printHeader("TCP")
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

	printHeader("Dns over UDP")
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

	printHeader("Dns over TCP")
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

	printHeader("DoT")
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

	printHeader("DoH")
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

	printHeader("WebSocket")
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
