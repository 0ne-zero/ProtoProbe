package flags_test

import (
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/flags"
)

func TestParseFlags_NoFlags_ReturnsErrorNoFlags(t *testing.T) {
	_, err := flags.ParseFlagsFrom([]string{})
	if err != flags.ErrorNoFlags {
		t.Fatalf("expected ErrorNoFlags, got %v", err)
	}
}

func TestParseFlags_AllAlone(t *testing.T) {
	opts, err := flags.ParseFlagsFrom([]string{"-all"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.All {
		t.Error("expected All=true")
	}
}

func TestParseFlags_AllWithProtocolFlag_ReturnsConflictError(t *testing.T) {
	protocolFlags := []string{"-icmp", "-tcp", "-tls", "-dou", "-dotcp", "-dot", "-doh", "-http", "-https", "-websocket", "-quic"}
	for _, pf := range protocolFlags {
		t.Run(pf, func(t *testing.T) {
			_, err := flags.ParseFlagsFrom([]string{"-all", pf})
			if err == nil {
				t.Fatalf("expected conflict error for -all %s, got nil", pf)
			}
			if err == flags.ErrorNoFlags {
				t.Fatalf("expected conflict error, not ErrorNoFlags")
			}
		})
	}
}

func TestParseFlags_IndividualProtocolFlags(t *testing.T) {
	tests := []struct {
		args  []string
		check func(opts flags.Options) bool
		name  string
	}{
		{[]string{"-icmp"}, func(o flags.Options) bool { return o.ICMP }, "icmp"},
		{[]string{"-tcp"}, func(o flags.Options) bool { return o.TCP }, "tcp"},
		{[]string{"-tls"}, func(o flags.Options) bool { return o.TLS }, "tls"},
		{[]string{"-dou"}, func(o flags.Options) bool { return o.DoUDP }, "dou"},
		{[]string{"-dotcp"}, func(o flags.Options) bool { return o.DoTCP }, "dotcp"},
		{[]string{"-dot"}, func(o flags.Options) bool { return o.DoT }, "dot"},
		{[]string{"-doh"}, func(o flags.Options) bool { return o.DoH }, "doh"},
		{[]string{"-http"}, func(o flags.Options) bool { return o.HTTP }, "http"},
		{[]string{"-https"}, func(o flags.Options) bool { return o.HTTPS }, "https"},
		{[]string{"-websocket"}, func(o flags.Options) bool { return o.WebSocket }, "websocket"},
		{[]string{"-quic"}, func(o flags.Options) bool { return o.QUIC }, "quic"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := flags.ParseFlagsFrom(tc.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.check(opts) {
				t.Errorf("expected flag %s to be true", tc.name)
			}
		})
	}
}

func TestParseFlags_TLSInsecure(t *testing.T) {
	opts, err := flags.ParseFlagsFrom([]string{"-tls", "-tls-insecure"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.TLS {
		t.Error("expected TLS=true")
	}
	if !opts.TLSInsecure {
		t.Error("expected TLSInsecure=true")
	}
}

func TestParseFlags_ICMPWithConfig(t *testing.T) {
	opts, err := flags.ParseFlagsFrom([]string{"-icmp", "-config", "/path/to/cfg.json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.ICMP {
		t.Error("expected ICMP=true")
	}
	if opts.ConfigFilePath != "/path/to/cfg.json" {
		t.Errorf("expected ConfigFilePath=/path/to/cfg.json, got %q", opts.ConfigFilePath)
	}
}

func TestParseFlags_ICMPWithJSON(t *testing.T) {
	opts, err := flags.ParseFlagsFrom([]string{"-icmp", "-json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.ICMP {
		t.Error("expected ICMP=true")
	}
	if !opts.JSON {
		t.Error("expected JSON=true")
	}
}

func TestParseFlags_DefaultConfigFilePath(t *testing.T) {
	opts, err := flags.ParseFlagsFrom([]string{"-icmp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.ConfigFilePath != flags.DefaultConfigFilePath {
		t.Errorf("expected default config path %q, got %q", flags.DefaultConfigFilePath, opts.ConfigFilePath)
	}
}
