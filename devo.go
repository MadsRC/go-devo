package devo

import (
	"net/http"
	"time"
)

const ALERTS_API_US_DEFAULT_ENDPOINT = "https://api-us.devo.com/alerts"
const ALERTS_API_EU_DEFAULT_ENDPOINT = "https://api-eu.devo.com/alerts"

type Config struct {
	HTTP *http.Client
	Alerts struct {
		Token string
		Address string
	}
}

type Client struct {
	Config *Config
}

func NewClient(config *Config) *Client {
	if config == nil {
		config = &Config{}
	}
	if config.HTTP == nil {
		config.HTTP = &http.Client{Timeout: 10 * time.Second}
	}

	if config.Alerts.Address == "" {
		config.Alerts.Address = ALERTS_API_US_DEFAULT_ENDPOINT
	}

	return &Client{Config: config}
}