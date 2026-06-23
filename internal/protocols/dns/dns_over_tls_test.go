package dns_test

import (
	"crypto/tls"
	"net"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	dnspkg "github.com/0ne-zero/ProtoProbe/internal/protocols/dns"
	"github.com/miekg/dns"
)

func startDoTServer(t *testing.T) (host string, port int) {
	t.Helper()
	cert := selfSignedCert(t)
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	ln, err := tls.Listen("tcp", "127.0.0.1:0", tlsCfg)
	if err != nil {
		t.Fatalf("startDoTServer: Listen: %v", err)
	}
	addr := ln.Addr().(*net.TCPAddr)
	host = addr.IP.String()
	port = addr.Port

	started := make(chan struct{})
	srv := &dns.Server{
		Listener:          ln,
		Net:               "tcp-tls",
		Handler:           dns.HandlerFunc(dnsHandler),
		NotifyStartedFunc: func() { close(started) },
	}

	go func() {
		srv.ActivateAndServe() //nolint:errcheck
	}()

	<-started
	t.Cleanup(func() { srv.Shutdown() }) //nolint:errcheck

	return host, port
}

func TestDoT_SuccessInsecure(t *testing.T) {
	host, port := startDoTServer(t)

	req := &config.HostPortQuery{
		HostPort: config.HostPort{Host: host, Port: port},
		Query:    "example.com",
	}
	result, err := dnspkg.TestDoT(req, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RTT <= 0 {
		t.Errorf("expected RTT > 0, got %v", result.RTT)
	}
}

func TestDoT_Refused(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	req := &config.HostPortQuery{
		HostPort: config.HostPort{Host: "127.0.0.1", Port: port},
		Query:    "example.com",
	}
	_, err = dnspkg.TestDoT(req, true)
	if err == nil {
		t.Fatal("expected error for refused connection, got nil")
	}
}
