package protocols_test

import (
	"context"
	"crypto/tls"
	"net"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	"github.com/0ne-zero/ProtoProbe/internal/protocols"
	"github.com/quic-go/quic-go"
)

func startQUICServer(t *testing.T) int {
	t.Helper()
	cert := selfSignedCert(t)
	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3"},
	}

	ln, err := quic.ListenAddr("127.0.0.1:0", tlsCfg, nil)
	if err != nil {
		t.Fatalf("startQUICServer: ListenAddr: %v", err)
	}
	port := ln.Addr().(*net.UDPAddr).Port
	t.Cleanup(func() { ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				return
			}
			conn.CloseWithError(0, "")
		}
	}()
	return port
}

func TestQUIC_SuccessInsecure(t *testing.T) {
	port := startQUICServer(t)

	result, err := protocols.TestQUIC(&config.HostPort{Host: "127.0.0.1", Port: port}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RTT <= 0 {
		t.Errorf("expected RTT > 0, got %v", result.RTT)
	}
}

func TestQUIC_Refused(t *testing.T) {
	port := freePort(t)
	_, err := protocols.TestQUIC(&config.HostPort{Host: "127.0.0.1", Port: port}, true)
	if err == nil {
		t.Fatal("expected error for no listener, got nil")
	}
}
