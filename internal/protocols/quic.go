package protocols

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	"github.com/quic-go/quic-go"
)

type QUICResult struct {
	RTT time.Duration
}

// TestQUIC dials a QUIC (HTTP/3) endpoint and returns the handshake RTT.
// insecureSkipVerify skips TLS certificate validation (useful for censorship probing).
func TestQUIC(target *config.HostPort, insecureSkipVerify bool) (QUICResult, error) {
	addr := net.JoinHostPort(target.Host, fmt.Sprintf("%d", target.Port))
	tlsConf := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify, //nolint:gosec
		NextProtos:         []string{"h3"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	conn, err := quic.DialAddr(ctx, addr, tlsConf, nil)
	if err != nil {
		return QUICResult{}, err
	}
	defer conn.CloseWithError(0, "")
	rtt := time.Since(start)
	return QUICResult{RTT: rtt}, nil
}
