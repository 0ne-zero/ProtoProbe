package dns

import (
	"context"
	"crypto/tls"
	"encoding/binary"
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

// TestDoQ sends a DNS query over a QUIC connection.
// Although RFC 9250 specifies no framing, all deployed servers (Cloudflare,
// AdGuard, etc.) were built on earlier drafts that use the same 2-octet
// length prefix as DNS-over-TCP. We match that de-facto standard.
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
	defer stream.Close()

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(target.Query), dns.TypeA)
	m.RecursionDesired = true
	m.Id = 0 // RFC 9250 §4.2.1: SHOULD be zero

	wire, err := m.Pack()
	if err != nil {
		return nil, err
	}

	// 2-byte length prefix (de-facto standard used by deployed servers)
	buf := make([]byte, 2+len(wire))
	binary.BigEndian.PutUint16(buf, uint16(len(wire)))
	copy(buf[2:], wire)

	if _, err := stream.Write(buf); err != nil {
		return nil, err
	}

	// Read 2-byte response length, then the DNS payload
	var lenBuf [2]byte
	if _, err := io.ReadFull(stream, lenBuf[:]); err != nil {
		return nil, fmt.Errorf("read response length: %w", err)
	}
	resp := make([]byte, binary.BigEndian.Uint16(lenBuf[:]))
	if _, err := io.ReadFull(stream, resp); err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	rtt := time.Since(start)

	var reply dns.Msg
	if err := reply.Unpack(resp); err != nil {
		return nil, fmt.Errorf("unpack response: %w", err)
	}

	return &DoQResult{RTT: rtt}, nil
}
