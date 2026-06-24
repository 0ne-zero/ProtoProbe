package protocols

import (
	"fmt"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/internal/config"
)

type NTPResult struct {
	RTT time.Duration
}

// TestNTP sends an NTPv3 client request over UDP (RFC 5905) and validates
// the server response.
func TestNTP(target *config.HostPort) (*NTPResult, error) {
	addr := net.JoinHostPort(target.Host, fmt.Sprintf("%d", target.Port))

	conn, err := net.DialTimeout("udp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second)) //nolint:errcheck

	// 48-byte NTP packet; first byte: LI=0, VN=3, Mode=3 (client)
	req := make([]byte, 48)
	req[0] = 0x1B

	start := time.Now()
	if _, err := conn.Write(req); err != nil {
		return nil, err
	}

	resp := make([]byte, 48)
	n, err := conn.Read(resp)
	rtt := time.Since(start)
	if err != nil {
		return nil, err
	}

	if n < 48 {
		return nil, fmt.Errorf("response too short (%d bytes)", n)
	}
	if mode := resp[0] & 0x07; mode != 4 {
		return nil, fmt.Errorf("unexpected NTP mode %d (want 4)", mode)
	}

	return &NTPResult{RTT: rtt}, nil
}
