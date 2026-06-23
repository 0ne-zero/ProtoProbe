package dns_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	dnspkg "github.com/0ne-zero/ProtoProbe/internal/protocols/dns"
)

func TestDoH_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/dns-message")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	req := &config.URLQuery{
		URL:   srv.URL,
		Query: "example.com",
	}
	result, err := dnspkg.TestDoH(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RTT <= 0 {
		t.Errorf("expected RTT > 0, got %v", result.RTT)
	}
}

func TestDoH_Refused(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	req := &config.URLQuery{
		URL:   fmt.Sprintf("http://127.0.0.1:%d", port),
		Query: "example.com",
	}
	_, err = dnspkg.TestDoH(req)
	if err == nil {
		t.Fatal("expected error for refused connection, got nil")
	}
}
