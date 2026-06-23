package dns_test

import (
	"net"
	"testing"

	dnspkg "github.com/0ne-zero/ProtoProbe/internal/protocols/dns"
	"github.com/0ne-zero/ProtoProbe/internal/config"
)

func TestDnsOverUDP_Success(t *testing.T) {
	host, port := startDNSServer(t, "udp")

	req := &config.HostPortQuery{
		HostPort: config.HostPort{Host: host, Port: port},
		Query:    "example.com",
	}
	result, err := dnspkg.TestDnsOverUDP(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RTT <= 0 {
		t.Errorf("expected RTT > 0, got %v", result.RTT)
	}
}

func TestDnsOverUDP_Refused(t *testing.T) {
	// Find a port with no listener.
	ln, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenPacket: %v", err)
	}
	port := ln.LocalAddr().(*net.UDPAddr).Port
	ln.Close()

	req := &config.HostPortQuery{
		HostPort: config.HostPort{Host: "127.0.0.1", Port: port},
		Query:    "example.com",
	}
	_, err = dnspkg.TestDnsOverUDP(req)
	if err == nil {
		t.Fatal("expected error for refused connection, got nil")
	}
}
