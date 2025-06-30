package dns

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/0ne-zero/ProtoProbe/config"
	"github.com/miekg/dns"
)

type DoHResult struct {
	RTT time.Duration
}

func TestDoH(dnsRequest *config.DNS_URL_Query) (DoHResult, error) {
	m := new(dns.Msg)
	m.SetQuestion(fmt.Sprintf("%s.", dnsRequest.Query), dns.TypeA)
	raw, _ := m.Pack()
	req, _ := http.NewRequest("POST", dnsRequest.Addr, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/dns-message")
	client := &http.Client{Timeout: timeout}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return DoHResult{}, err
	}
	resp.Body.Close()
	rtt := time.Since(start)
	return DoHResult{RTT: rtt}, nil
}
