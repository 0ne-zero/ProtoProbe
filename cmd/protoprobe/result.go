package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type ProbeResult struct {
	Protocol   string   `json:"protocol"`
	Target     string   `json:"target"`
	Success    bool     `json:"success"`
	RTTMs      *int64   `json:"rtt_ms,omitempty"`
	PacketLoss *float64 `json:"packet_loss,omitempty"`
	Error      string   `json:"error,omitempty"`
}

func rttMillis(d time.Duration) *int64 {
	ms := d.Milliseconds()
	return &ms
}

func ptrFloat64(f float64) *float64 { return &f }

func printHuman(results []ProbeResult) {
	for _, r := range results {
		if r.Success {
			if r.PacketLoss != nil {
				log.Printf("[%s] | %s | avg-rtt: %dms | packet-loss: %.2f%% ✅\n", r.Protocol, r.Target, *r.RTTMs, *r.PacketLoss)
			} else {
				log.Printf("[%s] | %s | rtt: %dms ✅\n", r.Protocol, r.Target, *r.RTTMs)
			}
		} else {
			log.Printf("[%s] | %s | %s ❌\n", r.Protocol, r.Target, r.Error)
		}
	}
}

func printJSON(results []ProbeResult) {
	out := struct {
		Results []ProbeResult `json:"results"`
	}{Results: results}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		log.Fatal(err)
	}
}
