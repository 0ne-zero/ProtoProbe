package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/config"
)

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeConfigFile: %v", err)
	}
	return path
}

func TestLoadConfig_ValidFull(t *testing.T) {
	path := writeConfigFile(t, `{
		"icmp": ["8.8.8.8", "1.1.1.1"],
		"tcp": [{"host": "example.com", "port": 80}],
		"tls": [{"host": "example.com", "port": 443}],
		"dns": [{"host": "8.8.8.8", "port": 53, "query": "example.com"}],
		"dot": [{"host": "1.1.1.1", "port": 853, "query": "example.com"}],
		"doh": [{"address": "https://cloudflare-dns.com/dns-query", "query": "example.com"}],
		"http": ["http://example.com"],
		"https": ["https://example.com"],
		"quic": [{"host": "quic.tech", "port": 8443}],
		"websocket": ["ws://example.com/ws"]
	}`)

	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.ICMP) != 2 || cfg.ICMP[0] != "8.8.8.8" || cfg.ICMP[1] != "1.1.1.1" {
		t.Errorf("ICMP: got %v", cfg.ICMP)
	}
	if len(cfg.TCP) != 1 || cfg.TCP[0].Host != "example.com" || cfg.TCP[0].Port != 80 {
		t.Errorf("TCP: got %v", cfg.TCP)
	}
	if len(cfg.TLS) != 1 || cfg.TLS[0].Host != "example.com" || cfg.TLS[0].Port != 443 {
		t.Errorf("TLS: got %v", cfg.TLS)
	}
	if len(cfg.DNS) != 1 || cfg.DNS[0].Host != "8.8.8.8" || cfg.DNS[0].Port != 53 || cfg.DNS[0].Query != "example.com" {
		t.Errorf("DNS: got %v", cfg.DNS)
	}
	if len(cfg.DoT) != 1 || cfg.DoT[0].Host != "1.1.1.1" || cfg.DoT[0].Port != 853 || cfg.DoT[0].Query != "example.com" {
		t.Errorf("DoT: got %v", cfg.DoT)
	}
	if len(cfg.DoH) != 1 || cfg.DoH[0].URL != "https://cloudflare-dns.com/dns-query" || cfg.DoH[0].Query != "example.com" {
		t.Errorf("DoH: got %v", cfg.DoH)
	}
	if len(cfg.HTTP) != 1 || cfg.HTTP[0] != "http://example.com" {
		t.Errorf("HTTP: got %v", cfg.HTTP)
	}
	if len(cfg.HTTPS) != 1 || cfg.HTTPS[0] != "https://example.com" {
		t.Errorf("HTTPS: got %v", cfg.HTTPS)
	}
	if len(cfg.QUIC) != 1 || cfg.QUIC[0].Host != "quic.tech" || cfg.QUIC[0].Port != 8443 {
		t.Errorf("QUIC: got %v", cfg.QUIC)
	}
	if len(cfg.WebSocket) != 1 || cfg.WebSocket[0] != "ws://example.com/ws" {
		t.Errorf("WebSocket: got %v", cfg.WebSocket)
	}
}

func TestLoadConfig_EmptyJSON(t *testing.T) {
	path := writeConfigFile(t, `{}`)

	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.ICMP) != 0 {
		t.Errorf("expected empty ICMP, got %v", cfg.ICMP)
	}
	if len(cfg.TCP) != 0 {
		t.Errorf("expected empty TCP, got %v", cfg.TCP)
	}
	if len(cfg.TLS) != 0 {
		t.Errorf("expected empty TLS, got %v", cfg.TLS)
	}
	if len(cfg.DNS) != 0 {
		t.Errorf("expected empty DNS, got %v", cfg.DNS)
	}
	if len(cfg.DoT) != 0 {
		t.Errorf("expected empty DoT, got %v", cfg.DoT)
	}
	if len(cfg.DoH) != 0 {
		t.Errorf("expected empty DoH, got %v", cfg.DoH)
	}
	if len(cfg.HTTP) != 0 {
		t.Errorf("expected empty HTTP, got %v", cfg.HTTP)
	}
	if len(cfg.HTTPS) != 0 {
		t.Errorf("expected empty HTTPS, got %v", cfg.HTTPS)
	}
	if len(cfg.QUIC) != 0 {
		t.Errorf("expected empty QUIC, got %v", cfg.QUIC)
	}
	if len(cfg.WebSocket) != 0 {
		t.Errorf("expected empty WebSocket, got %v", cfg.WebSocket)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := writeConfigFile(t, `not valid json {{{`)

	_, err := config.LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := config.LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfig_PartialFields(t *testing.T) {
	path := writeConfigFile(t, `{
		"icmp": ["192.168.1.1"]
	}`)

	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.ICMP) != 1 || cfg.ICMP[0] != "192.168.1.1" {
		t.Errorf("ICMP: got %v", cfg.ICMP)
	}
	if len(cfg.TCP) != 0 {
		t.Errorf("expected empty TCP, got %v", cfg.TCP)
	}
}
