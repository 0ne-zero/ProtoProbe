package dns_test

import (
	"net"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	dnspkg "github.com/0ne-zero/ProtoProbe/internal/protocols/dns"
)

func TestDNSTCP_Success(t *testing.T) {
	host, port := startDNSServer(t, "tcp")

	req := &config.HostPortQuery{
		HostPort: config.HostPort{Host: host, Port: port},
		Query:    "example.com",
	}
	result, err := dnspkg.TestDNSTCP(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RTT <= 0 {
		t.Errorf("expected RTT > 0, got %v", result.RTT)
	}
}

func TestDNSTCP_Refused(t *testing.T) {
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
	_, err = dnspkg.TestDNSTCP(req)
	if err == nil {
		t.Fatal("expected error for refused connection, got nil")
	}
}
