package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/0ne-zero/ProtoProbe/internal/flags"
	"github.com/0ne-zero/ProtoProbe/internal/config"
	"github.com/0ne-zero/ProtoProbe/internal/protocols"
	"github.com/0ne-zero/ProtoProbe/internal/protocols/dns"
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

	var allResults []ProbeResult

	if opts.All || opts.ICMP {
		if !opts.JSON {
			printHeader("ICMP")
		}
		results := runICMPTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.DoUDP {
		if !opts.JSON {
			printHeader("Dns over UDP")
		}
		results := runDnsOverUDPTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.DoTCP {
		if !opts.JSON {
			printHeader("Dns over TCP")
		}
		results := runDnsOverTCPTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.TCP {
		if !opts.JSON {
			printHeader("TCP")
		}
		results := runTCPTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.TLS {
		if !opts.JSON {
			printHeader("TLS")
		}
		results := runTLSTest(cfg, opts.TLSInsecure)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.ECH {
		if !opts.JSON {
			printHeader("ECH")
		}
		results := runECHTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.DoT {
		if !opts.JSON {
			printHeader("DoT")
		}
		results := runDoTTest(cfg, opts.DoTInsecure)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.DoQ {
		if !opts.JSON {
			printHeader("DoQ")
		}
		results := runDoQTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.DoH {
		if !opts.JSON {
			printHeader("DoH")
		}
		results := runDoHTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.HTTP {
		if !opts.JSON {
			printHeader("HTTP")
		}
		results := runHTTPTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.HTTPS {
		if !opts.JSON {
			printHeader("HTTPS")
		}
		results := runHTTPSTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.QUIC {
		if !opts.JSON {
			printHeader("QUIC")
		}
		results := runQUICTest(cfg, opts.QUICInsecure)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.WebSocket {
		if !opts.JSON {
			printHeader("WebSocket")
		}
		results := runWebSocketTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.STUN {
		if !opts.JSON {
			printHeader("STUN")
		}
		results := runSTUNTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.All || opts.NTP {
		if !opts.JSON {
			printHeader("NTP")
		}
		results := runNTPTest(cfg)
		if !opts.JSON {
			printHuman(results)
		}
		allResults = append(allResults, results...)
	}

	if opts.JSON {
		printJSON(allResults)
	}
}

func runICMPTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, host := range cfg.ICMP {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			r := ProbeResult{Protocol: "ICMP", Target: h}
			res, err := protocols.TestICMP(h)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.AvgRtt)
				r.PacketLoss = ptrFloat64(res.PacketLoss)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(host)
	}
	wg.Wait()
	return results
}

func runTCPTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.TCP {
		wg.Add(1)
		go func(hp config.HostPort) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "TCP", Target: target}
			res, err := protocols.TestTCP(hp)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func runDnsOverUDPTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.DNS {
		wg.Add(1)
		go func(hp config.HostPortQuery) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "DNS/UDP", Target: target}
			res, err := dns.TestDnsOverUDP(&hp)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func runDnsOverTCPTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.DNS {
		wg.Add(1)
		go func(hp config.HostPortQuery) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "DNS/TCP", Target: target}
			res, err := dns.TestDNSTCP(&hp)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func runDoTTest(cfg *config.Config, insecureSkipVerify bool) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.DoT {
		wg.Add(1)
		go func(hp config.HostPortQuery) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "DNS/TLS (DoT)", Target: target}
			res, err := dns.TestDoT(&hp, insecureSkipVerify)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func runDoHTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, uq := range cfg.DoH {
		wg.Add(1)
		go func(uq config.URLQuery) {
			defer wg.Done()
			r := ProbeResult{Protocol: "DNS/HTTPS (DoH)", Target: uq.URL}
			res, err := dns.TestDoH(&uq)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(uq)
	}
	wg.Wait()
	return results
}

func runWebSocketTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, server := range cfg.WebSocket {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			r := ProbeResult{Protocol: "WebSocket", Target: s}
			res, err := protocols.TestWebSocket(s)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(server)
	}
	wg.Wait()
	return results
}

func runQUICTest(cfg *config.Config, insecureSkipVerify bool) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.QUIC {
		wg.Add(1)
		go func(hp config.HostPort) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "QUIC", Target: target}
			res, err := protocols.TestQUIC(&hp, insecureSkipVerify)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func runTLSTest(cfg *config.Config, insecureSkipVerify bool) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.TLS {
		wg.Add(1)
		go func(hp config.HostPort) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "TLS", Target: target}
			res, err := protocols.TestTLS(&hp, insecureSkipVerify)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func runHTTPTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, url := range cfg.HTTP {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			r := ProbeResult{Protocol: "HTTP", Target: u}
			res, err := protocols.TestHTTP(u)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
				r.StatusCode = ptrInt(res.StatusCode)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(url)
	}
	wg.Wait()
	return results
}

func runHTTPSTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, url := range cfg.HTTPS {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			r := ProbeResult{Protocol: "HTTPS", Target: u}
			res, err := protocols.TestHTTP(u)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
				r.StatusCode = ptrInt(res.StatusCode)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(url)
	}
	wg.Wait()
	return results
}

func runECHTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.ECH {
		wg.Add(1)
		go func(hp config.HostPort) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "ECH", Target: target}
			res, err := protocols.TestECH(&hp)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
				r.ECHAccepted = ptrBool(res.ECHAccepted)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func runDoQTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.DoQ {
		wg.Add(1)
		go func(hp config.HostPortQuery) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "DNS/QUIC (DoQ)", Target: target}
			res, err := dns.TestDoQ(&hp)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func runSTUNTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.STUN {
		wg.Add(1)
		go func(hp config.HostPort) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "STUN", Target: target}
			res, err := protocols.TestSTUN(&hp)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func runNTPTest(cfg *config.Config) []ProbeResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []ProbeResult
	for _, hp := range cfg.NTP {
		wg.Add(1)
		go func(hp config.HostPort) {
			defer wg.Done()
			target := net.JoinHostPort(hp.Host, fmt.Sprintf("%d", hp.Port))
			r := ProbeResult{Protocol: "NTP", Target: target}
			res, err := protocols.TestNTP(&hp)
			if err != nil {
				r.Error = err.Error()
			} else {
				r.Success = true
				r.RTTMs = rttMillis(res.RTT)
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(hp)
	}
	wg.Wait()
	return results
}

func printHeader(proto string) {
	fmt.Println()
	fmt.Printf("%s:\n", proto)
}
