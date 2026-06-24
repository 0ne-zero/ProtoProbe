package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	"github.com/0ne-zero/ProtoProbe/internal/flags"
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

	// Each protocol gets its own channel so results can be drained in
	// protocol order while all tests run concurrently.
	var ordered []<-chan ProbeResult

	enqueue := func(fn func(chan<- ProbeResult)) {
		ch := make(chan ProbeResult, 16)
		ordered = append(ordered, ch)
		go func() { fn(ch); close(ch) }()
	}

	if opts.All || opts.ICMP {
		enqueue(func(ch chan<- ProbeResult) { runICMPTest(cfg, ch) })
	}
	if opts.All || opts.DoUDP {
		enqueue(func(ch chan<- ProbeResult) { runDnsOverUDPTest(cfg, ch) })
	}
	if opts.All || opts.DoTCP {
		enqueue(func(ch chan<- ProbeResult) { runDnsOverTCPTest(cfg, ch) })
	}
	if opts.All || opts.TCP {
		enqueue(func(ch chan<- ProbeResult) { runTCPTest(cfg, ch) })
	}
	if opts.All || opts.TLS {
		enqueue(func(ch chan<- ProbeResult) { runTLSTest(cfg, opts.TLSInsecure, ch) })
	}
	if opts.All || opts.ECH {
		enqueue(func(ch chan<- ProbeResult) { runECHTest(cfg, ch) })
	}
	if opts.All || opts.DoT {
		enqueue(func(ch chan<- ProbeResult) { runDoTTest(cfg, opts.DoTInsecure, ch) })
	}
	if opts.All || opts.DoQ {
		enqueue(func(ch chan<- ProbeResult) { runDoQTest(cfg, ch) })
	}
	if opts.All || opts.DoH {
		enqueue(func(ch chan<- ProbeResult) { runDoHTest(cfg, ch) })
	}
	if opts.All || opts.HTTP {
		enqueue(func(ch chan<- ProbeResult) { runHTTPTest(cfg, ch) })
	}
	if opts.All || opts.HTTPS {
		enqueue(func(ch chan<- ProbeResult) { runHTTPSTest(cfg, ch) })
	}
	if opts.All || opts.QUIC {
		enqueue(func(ch chan<- ProbeResult) { runQUICTest(cfg, opts.QUICInsecure, ch) })
	}
	if opts.All || opts.WebSocket {
		enqueue(func(ch chan<- ProbeResult) { runWebSocketTest(cfg, ch) })
	}
	if opts.All || opts.STUN {
		enqueue(func(ch chan<- ProbeResult) { runSTUNTest(cfg, ch) })
	}
	if opts.All || opts.NTP {
		enqueue(func(ch chan<- ProbeResult) { runNTPTest(cfg, ch) })
	}

	if opts.JSON {
		var all []ProbeResult
		for _, ch := range ordered {
			for r := range ch {
				all = append(all, r)
			}
		}
		printJSON(all)
	} else {
		wProto, wTarget := configWidths(cfg, opts)
		streamTable(wProto, wTarget, ordered)
	}
}

// configWidths computes column widths from the config before any tests run,
// so the table header can be printed immediately.
func configWidths(cfg *config.Config, opts flags.Options) (wProto, wTarget int) {
	wProto = len("PROTOCOL")
	wTarget = len("TARGET")

	upd := func(proto string, targets ...string) {
		if n := len(proto); n > wProto {
			wProto = n
		}
		for _, t := range targets {
			if n := len(t); n > wTarget {
				wTarget = n
			}
		}
	}
	hp := func(host string, port int) string {
		return net.JoinHostPort(host, fmt.Sprintf("%d", port))
	}

	if opts.All || opts.ICMP {
		upd("ICMP", cfg.ICMP...)
	}
	if opts.All || opts.DoUDP {
		var ts []string
		for _, h := range cfg.DNS {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("DNS/UDP", ts...)
	}
	if opts.All || opts.DoTCP {
		var ts []string
		for _, h := range cfg.DNS {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("DNS/TCP", ts...)
	}
	if opts.All || opts.TCP {
		var ts []string
		for _, h := range cfg.TCP {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("TCP", ts...)
	}
	if opts.All || opts.TLS {
		var ts []string
		for _, h := range cfg.TLS {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("TLS", ts...)
	}
	if opts.All || opts.ECH {
		var ts []string
		for _, h := range cfg.ECH {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("ECH", ts...)
	}
	if opts.All || opts.DoT {
		var ts []string
		for _, h := range cfg.DoT {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("DNS/TLS (DoT)", ts...)
	}
	if opts.All || opts.DoQ {
		var ts []string
		for _, h := range cfg.DoQ {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("DNS/QUIC (DoQ)", ts...)
	}
	if opts.All || opts.DoH {
		var ts []string
		for _, h := range cfg.DoH {
			ts = append(ts, h.URL)
		}
		upd("DNS/HTTPS (DoH)", ts...)
	}
	if opts.All || opts.HTTP {
		upd("HTTP", cfg.HTTP...)
	}
	if opts.All || opts.HTTPS {
		upd("HTTPS", cfg.HTTPS...)
	}
	if opts.All || opts.QUIC {
		var ts []string
		for _, h := range cfg.QUIC {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("QUIC", ts...)
	}
	if opts.All || opts.WebSocket {
		upd("WebSocket", cfg.WebSocket...)
	}
	if opts.All || opts.STUN {
		var ts []string
		for _, h := range cfg.STUN {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("STUN", ts...)
	}
	if opts.All || opts.NTP {
		var ts []string
		for _, h := range cfg.NTP {
			ts = append(ts, hp(h.Host, h.Port))
		}
		upd("NTP", ts...)
	}

	return wProto, wTarget
}

func runICMPTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(host)
	}
	wg.Wait()
}

func runTCPTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}

func runDnsOverUDPTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}

func runDnsOverTCPTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}

func runTLSTest(cfg *config.Config, insecureSkipVerify bool, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}

func runECHTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}

func runDoTTest(cfg *config.Config, insecureSkipVerify bool, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}

func runDoQTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}

func runDoHTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(uq)
	}
	wg.Wait()
}

func runHTTPTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(url)
	}
	wg.Wait()
}

func runHTTPSTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(url)
	}
	wg.Wait()
}

func runQUICTest(cfg *config.Config, insecureSkipVerify bool, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}

func runWebSocketTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(server)
	}
	wg.Wait()
}

func runSTUNTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}

func runNTPTest(cfg *config.Config, out chan<- ProbeResult) {
	var wg sync.WaitGroup
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
			out <- r
		}(hp)
	}
	wg.Wait()
}
