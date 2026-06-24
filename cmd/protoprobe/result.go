package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type ProbeResult struct {
	Protocol    string   `json:"protocol"`
	Target      string   `json:"target"`
	Success     bool     `json:"success"`
	RTTMs       *int64   `json:"rtt_ms,omitempty"`
	PacketLoss  *float64 `json:"packet_loss,omitempty"`
	StatusCode  *int     `json:"status_code,omitempty"`
	ECHAccepted *bool    `json:"ech_accepted,omitempty"`
	Error       string   `json:"error,omitempty"`
}

func rttMillis(d time.Duration) *int64 {
	ms := d.Milliseconds()
	return &ms
}

func ptrFloat64(f float64) *float64 { return &f }
func ptrInt(i int) *int             { return &i }
func ptrBool(b bool) *bool          { return &b }

// streamTable prints the header immediately, then prints each row as it
// arrives on resultCh. wProto and wTarget are pre-computed from the config
// so the header can be aligned before any test completes.
func streamTable(wProto, wTarget int, resultCh []<-chan ProbeResult) {
	const wRTT = 7 // wide enough for "99999ms"

	line := func(proto, target, rtt, result string) {
		fmt.Printf("%-*s  %-*s  %-*s  %s\n", wProto, proto, wTarget, target, wRTT, rtt, result)
	}
	dash := func(n int) string { return strings.Repeat("─", n) }

	line("PROTOCOL", "TARGET", "RTT", "RESULT")
	line(dash(wProto), dash(wTarget), dash(wRTT), dash(6))

	for _, ch := range resultCh {
		for r := range ch {
		rttStr := "-"
		if r.RTTMs != nil {
			rttStr = fmt.Sprintf("%dms", *r.RTTMs)
		}

		var result string
		if r.Success {
			switch {
			case r.PacketLoss != nil:
				result = fmt.Sprintf("✅  loss: %.2f%%", *r.PacketLoss)
			case r.StatusCode != nil:
				result = fmt.Sprintf("✅  status: %d", *r.StatusCode)
			case r.ECHAccepted != nil:
				if *r.ECHAccepted {
					result = "✅  ech: accepted"
				} else {
					result = "✅  ech: not accepted"
				}
			default:
				result = "✅"
			}
		} else {
			result = "❌  " + r.Error
		}

		line(r.Protocol, r.Target, rttStr, result)
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
