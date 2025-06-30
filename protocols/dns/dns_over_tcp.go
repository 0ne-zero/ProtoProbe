package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/config"
	"github.com/miekg/dns"
)

type DNSOverTCPResult struct {
	RTT time.Duration
}

func TestDNSTCP(dnsRequest *config.DNS_Host_Port_Query) (*DNSOverTCPResult, error) {
	m := new(dns.Msg)
	m.SetQuestion(fmt.Sprintf("%s.", dnsRequest.Query), dns.TypeA)
	c := new(dns.Client)
	c.Net = "tcp"
	c.Timeout = timeout
	start := time.Now()
	serverAddr := net.JoinHostPort(dnsRequest.Host, fmt.Sprintf("%d", dnsRequest.Port))
	_, _, err := c.Exchange(m, serverAddr)
	if err != nil {
		return nil, err
	}
	rtt := time.Since(start)
	return &DNSOverTCPResult{RTT: rtt}, nil
}
