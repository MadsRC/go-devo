package alerts

import (
	"net/http"
	"time"
)

type Config struct {
	HTTP    *http.Client
	Token   string
	Address string
}

type Client struct {
	Config *Config
}

func NewClient(config *Config) *Client {
	if config.HTTP == nil {
		config.HTTP = &http.Client{Timeout: 10 * time.Second}
	}

	if config.Address == "" {
		config.Address = "https://api-us.devo.com/alerts"
	}

	return &Client{Config: config}
}
