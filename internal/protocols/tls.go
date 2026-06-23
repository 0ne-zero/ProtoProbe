package protocols

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/internal/config"
)

type TLSResult struct {
	RTT time.Duration
}

// TestTLS dials addr over TCP and completes a TLS handshake, returning the
// total time from dial to handshake completion. host is used as the SNI
// ServerName unless overridden. insecureSkipVerify disables cert validation,
// which is useful when probing servers with self-signed or blocked certs.
func TestTLS(target *config.HostPort, insecureSkipVerify bool) (TLSResult, error) {
	addr := net.JoinHostPort(target.Host, fmt.Sprintf("%d", target.Port))
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	tlsConf := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify, //nolint:gosec
		ServerName:         target.Host,
	}
	start := time.Now()
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConf)
	rtt := time.Since(start)
	if err != nil {
		return TLSResult{}, err
	}
	defer conn.Close()
	return TLSResult{RTT: rtt}, nil
}
