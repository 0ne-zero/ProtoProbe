package protocols_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/protocols"
)

func TestHTTP_200OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	result, err := protocols.TestHTTP(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RTT <= 0 {
		t.Errorf("expected RTT > 0, got %v", result.RTT)
	}
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected StatusCode=200, got %d", result.StatusCode)
	}
}

func TestHTTP_404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	result, err := protocols.TestHTTP(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.StatusCode != http.StatusNotFound {
		t.Errorf("expected StatusCode=404, got %d", result.StatusCode)
	}
}

func TestHTTP_301NotFollowed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://example.com", http.StatusMovedPermanently)
	}))
	t.Cleanup(srv.Close)

	result, err := protocols.TestHTTP(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.StatusCode != http.StatusMovedPermanently {
		t.Errorf("expected StatusCode=301 (redirect not followed), got %d", result.StatusCode)
	}
}

func TestHTTP_Refused(t *testing.T) {
	port := freePort(t)
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	_, err := protocols.TestHTTP(url)
	if err == nil {
		t.Fatal("expected error for refused connection, got nil")
	}
}
