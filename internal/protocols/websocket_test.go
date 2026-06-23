package protocols_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0ne-zero/ProtoProbe/internal/protocols"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func TestWebSocket_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.Close()
	}))
	t.Cleanup(srv.Close)

	// Convert http:// URL to ws://
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	result, err := protocols.TestWebSocket(wsURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RTT <= 0 {
		t.Errorf("expected RTT > 0, got %v", result.RTT)
	}
}

func TestWebSocket_Refused(t *testing.T) {
	port := freePort(t)
	url := fmt.Sprintf("ws://127.0.0.1:%d", port)
	_, err := protocols.TestWebSocket(url)
	if err == nil {
		t.Fatal("expected error for refused connection, got nil")
	}
}
