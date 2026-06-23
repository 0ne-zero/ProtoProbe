package protocols_test

import (
	"os"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/protocols"
)

func TestICMP_Loopback(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("requires root")
	}

	result, err := protocols.TestICMP("127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PacketLoss != 0 {
		t.Errorf("expected PacketLoss=0, got %v", result.PacketLoss)
	}
}

func TestICMP_InvalidHost(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("requires root")
	}

	_, err := protocols.TestICMP("something.invalid")
	if err == nil {
		t.Fatal("expected error for invalid host, got nil")
	}
}
