package protocols

import (
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketResult struct {
	RTT time.Duration
}

func TestWebSocket(url string) (WebSocketResult, error) {
	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	start := time.Now()
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return WebSocketResult{}, err
	}
	defer conn.Close()
	rtt := time.Since(start)
	return WebSocketResult{RTT: rtt}, nil
}
