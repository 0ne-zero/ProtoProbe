package protocols

import (
	"fmt"
	"net"
	"time"

	"github.com/0ne-zero/ProtoProbe/config"
)

type TCPResult struct {
	RTT time.Duration
}

func TestTCP(hostPort config.DNS_Host_Port_Query) (*TCPResult, error) {
	addr := net.JoinHostPort(hostPort.Host, fmt.Sprintf("%d", hostPort.Port))

	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	elapsed := time.Since(start)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return &TCPResult{RTT: elapsed}, nil

}
