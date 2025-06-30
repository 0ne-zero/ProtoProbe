package dns

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/config"
	"github.com/miekg/dns"
)

type DoTResult struct {
	RTT time.Duration
}

func TestDoT(dnsRequest *config.DNS_Host_Port_Query) (DoTResult, error) {
	m := new(dns.Msg)
	m.SetQuestion(fmt.Sprintf("%s.", dnsRequest.Query), dns.TypeA)
	c := new(dns.Client)
	c.Net = "tcp-tls"
	c.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	c.Timeout = timeout
	start := time.Now()
	serverAddr := net.JoinHostPort(dnsRequest.Host, fmt.Sprintf("%d", dnsRequest.Port))
	_, _, err := c.Exchange(m, serverAddr)
	if err != nil {
		return DoTResult{}, err
	}
	rtt := time.Since(start)
	return DoTResult{RTT: rtt}, nil
}
