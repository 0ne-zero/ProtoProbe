package protocols

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/internal/config"
	"github.com/miekg/dns"
)

type ECHResult struct {
	RTT         time.Duration
	ECHAccepted bool
}

func fetchECHConfigList(host string) ([]byte, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), dns.TypeHTTPS)
	m.RecursionDesired = true

	c := &dns.Client{Timeout: 5 * time.Second}
	r, _, err := c.Exchange(m, "1.1.1.1:53")
	if err != nil {
		return nil, fmt.Errorf("HTTPS record lookup: %w", err)
	}

	for _, ans := range r.Answer {
		if https, ok := ans.(*dns.HTTPS); ok {
			for _, param := range https.Value {
				if param.Key() == dns.SVCB_ECHCONFIG {
					if echParam, ok := param.(*dns.SVCBECHConfig); ok {
						return echParam.ECH, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("no ECH config in HTTPS record for %s", host)
}

// TestECH fetches the ECH config via DNS and performs a TLS handshake with
// Encrypted Client Hello enabled. ECHAccepted reports whether the server
// negotiated ECH (false means the connection succeeded but ECH was not used).
func TestECH(target *config.HostPort) (*ECHResult, error) {
	echConfig, err := fetchECHConfigList(target.Host)
	if err != nil {
		return nil, err
	}

	addr := net.JoinHostPort(target.Host, fmt.Sprintf("%d", target.Port))
	tlsConf := &tls.Config{
		ServerName:                     target.Host,
		EncryptedClientHelloConfigList: echConfig,
	}

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	start := time.Now()
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConf)
	rtt := time.Since(start)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return &ECHResult{
		RTT:         rtt,
		ECHAccepted: conn.ConnectionState().ECHAccepted,
	}, nil
}
