package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ICMP      []string        `json:"icmp"`
	TCP       []HostPort      `json:"tcp"`
	TLS       []HostPort      `json:"tls"`
	ECH       []HostPort      `json:"ech"`
	DNS       []HostPortQuery `json:"dns"`
	DoT       []HostPortQuery `json:"dot"`
	DoQ       []HostPortQuery `json:"doq"`
	DoH       []URLQuery      `json:"doh"`
	HTTP      []string        `json:"http"`
	HTTPS     []string        `json:"https"`
	QUIC      []HostPort      `json:"quic"`
	WebSocket []string        `json:"websocket"`
	STUN      []HostPort      `json:"stun"`
	NTP       []HostPort      `json:"ntp"`
}

type HostPort struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type HostPortQuery struct {
	HostPort
	Query string `json:"query"`
}

type URLQuery struct {
	URL   string `json:"address"`
	Query string `json:"query"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
