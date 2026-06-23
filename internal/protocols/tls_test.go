package protocols_test

import (
	"crypto/tls"
	"net"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	"github.com/0ne-zero/ProtoProbe/internal/protocols"
)

func startTLSServer(t *testing.T) int {
	t.Helper()
	cert := selfSignedCert(t)
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	ln, err := tls.Listen("tcp", "127.0.0.1:0", tlsCfg)
	if err != nil {
		t.Fatalf("startTLSServer: listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	t.Cleanup(func() { ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			// Complete the TLS handshake before closing.
			go func(c net.Conn) {
				tlsConn := c.(*tls.Conn)
				tlsConn.Handshake() //nolint:errcheck
				c.Close()
			}(conn)
		}
	}()
	return port
}

func TestTLS_SuccessInsecure(t *testing.T) {
	port := startTLSServer(t)

	result, err := protocols.TestTLS(&config.HostPort{Host: "127.0.0.1", Port: port}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RTT <= 0 {
		t.Errorf("expected RTT > 0, got %v", result.RTT)
	}
}

func TestTLS_CertVerificationFails(t *testing.T) {
	port := startTLSServer(t)

	_, err := protocols.TestTLS(&config.HostPort{Host: "127.0.0.1", Port: port}, false)
	if err == nil {
		t.Fatal("expected TLS cert verification error, got nil")
	}
}

func TestTLS_Refused(t *testing.T) {
	port := freePort(t)
	_, err := protocols.TestTLS(&config.HostPort{Host: "127.0.0.1", Port: port}, true)
	if err == nil {
		t.Fatal("expected error for refused connection, got nil")
	}
}
