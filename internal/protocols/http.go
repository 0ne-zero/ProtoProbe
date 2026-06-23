package protocols

import (
	"net/http"
	"time"
)

type HTTPResult struct {
	RTT        time.Duration
	StatusCode int
}

// TestHTTP performs a GET request to url and returns the RTT and HTTP status
// code. Redirects are not followed so a 3xx response is reported as-is —
// a redirect could otherwise mask a censorship block page.
func TestHTTP(url string) (HTTPResult, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	start := time.Now()
	resp, err := client.Get(url)
	rtt := time.Since(start)
	if err != nil {
		return HTTPResult{}, err
	}
	defer resp.Body.Close()
	return HTTPResult{RTT: rtt, StatusCode: resp.StatusCode}, nil
}
