package protocols

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/internal/config"
)

type STUNResult struct {
	RTT time.Duration
}

const stunMagicCookie uint32 = 0x2112A442

// TestSTUN sends a STUN Binding Request (RFC 8489) over UDP and validates
// the Binding Success Response.
func TestSTUN(target *config.HostPort) (*STUNResult, error) {
	addr := net.JoinHostPort(target.Host, fmt.Sprintf("%d", target.Port))

	conn, err := net.DialTimeout("udp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second)) //nolint:errcheck

	// STUN message header: type(2) + length(2) + magic(4) + txID(12) = 20 bytes
	req := make([]byte, 20)
	binary.BigEndian.PutUint16(req[0:2], 0x0001) // Binding Request
	binary.BigEndian.PutUint16(req[2:4], 0)      // no attributes
	binary.BigEndian.PutUint32(req[4:8], stunMagicCookie)
	rand.Read(req[8:20]) //nolint:errcheck

	start := time.Now()
	if _, err := conn.Write(req); err != nil {
		return nil, err
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	rtt := time.Since(start)
	if err != nil {
		return nil, err
	}

	if n < 20 {
		return nil, fmt.Errorf("response too short (%d bytes)", n)
	}
	if msgType := binary.BigEndian.Uint16(buf[0:2]); msgType != 0x0101 {
		return nil, fmt.Errorf("unexpected message type 0x%04x (want 0x0101)", msgType)
	}

	return &STUNResult{RTT: rtt}, nil
}
