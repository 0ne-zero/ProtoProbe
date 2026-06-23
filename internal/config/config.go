package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ICMP      []string        `json:"icmp"`
	TCP       []HostPort      `json:"tcp"`
	DNS       []HostPortQuery `json:"dns"`
	DoT       []HostPortQuery `json:"dot"`
	DoH       []URLQuery      `json:"doh"`
	WebSocket []string        `json:"websocket"`
	QUIC      []HostPort      `json:"quic"`
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
