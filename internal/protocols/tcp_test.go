package protocols_test

import (
	"net"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	"github.com/0ne-zero/ProtoProbe/internal/protocols"
)

func TestTCP_Success(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	t.Cleanup(func() { ln.Close() })

	// Accept connections in background so dial doesn't hang.
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	result, err := protocols.TestTCP(config.HostPort{Host: "127.0.0.1", Port: port})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RTT <= 0 {
		t.Errorf("expected RTT > 0, got %v", result.RTT)
	}
}

func TestTCP_Refused(t *testing.T) {
	port := freePort(t)
	_, err := protocols.TestTCP(config.HostPort{Host: "127.0.0.1", Port: port})
	if err == nil {
		t.Fatal("expected error for refused connection, got nil")
	}
}
