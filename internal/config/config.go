package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ICMPHost          []string              `json:"icmp"`
	TCPHostPort       []DNS_Host_Port_Query `json:"tcp"`
	NormalDNSHostPort []DNS_Host_Port_Query `json:"dns"`       // Both DNS over UDP and DNS over TCP server address (with query)
	DoT               []DNS_Host_Port_Query `json:"dot"`       // DNS over TLS server addres (with query)
	DoH               []DNS_URL_Query       `json:"doh"`       // DNS over HTTPS server address (with query)
	WebSocket         []string              `json:"websocket"` // Used for websocker server address
}

type DNS_Host_Port_Query struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Query string `json:"query"`
}
type DNS_URL_Query struct {
	Addr  string `json:"address"`
	Query string `json:"query"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(file, &cfg)
	return &cfg, err
}
