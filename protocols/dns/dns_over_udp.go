package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/config"
	"github.com/miekg/dns"
)

type DNSResult struct {
	RTT time.Duration
}

func TestDnsOverUDP(dnsRequest *config.DNS_Host_Port_Query) (*DNSResult, error) {
	m := new(dns.Msg)
	m.SetQuestion(fmt.Sprintf("%s.", dnsRequest.Query), dns.TypeA)
	c := new(dns.Client)
	c.Timeout = timeout

	serverAddr := net.JoinHostPort(dnsRequest.Host, fmt.Sprintf("%d", dnsRequest.Port))
	start := time.Now()
	_, _, err := c.Exchange(m, serverAddr)
	if err != nil {
		return nil, err
	}
	rtt := time.Since(start)
	return &DNSResult{RTT: rtt}, nil
}
