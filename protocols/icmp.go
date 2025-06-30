package protocols

import (
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type ICMPResult struct {
	AvgRtt     time.Duration
	PacketLoss float64
}

func TestICMP(host string) (*ICMPResult, error) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return nil, err
	}
	pinger.Count = 5
	pinger.Timeout = 5 * time.Second
	err = pinger.Run()
	if err != nil {
		return nil, err
	}
	stats := pinger.Statistics()
	return &ICMPResult{
		AvgRtt:     stats.AvgRtt,
		PacketLoss: stats.PacketLoss,
	}, nil
}
