package dns

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	"github.com/miekg/dns"
	"github.com/quic-go/quic-go"
)

type DoQResult struct {
	RTT time.Duration
}

// TestDoQ sends a DNS query over a QUIC connection (RFC 9250).
// Messages are sent without a 2-octet length prefix; stream FIN signals end.
func TestDoQ(target *config.HostPortQuery) (*DoQResult, error) {
	addr := net.JoinHostPort(target.Host, fmt.Sprintf("%d", target.Port))

	tlsConf := &tls.Config{
		NextProtos: []string{"doq"},
		ServerName: target.Host,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	conn, err := quic.DialAddr(ctx, addr, tlsConf, nil)
	if err != nil {
		return nil, err
	}
	defer conn.CloseWithError(0, "")

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(target.Query), dns.TypeA)
	m.RecursionDesired = true
	m.Id = 0 // RFC 9250 §4.2.1: SHOULD be set to zero

	wire, err := m.Pack()
	if err != nil {
		return nil, err
	}

	if _, err := stream.Write(wire); err != nil {
		return nil, err
	}
	// RFC 9250 §4.2.1: close the write side to signal end of query (no length prefix)
	stream.Close()

	resp, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}
	rtt := time.Since(start)

	if len(resp) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	var reply dns.Msg
	if err := reply.Unpack(resp); err != nil {
		return nil, fmt.Errorf("unpack response: %w", err)
	}

	return &DoQResult{RTT: rtt}, nil
}
