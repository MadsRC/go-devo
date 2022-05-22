package devo

import (
	"net/http"
	"time"
)

type Config struct {
	HttpClient *http.Client
	Alerts     struct {
		Token string
	}
}

type Client struct {
	http   *http.Client
	Alerts alertClient
}

type alertClient struct {
	http  *http.Client
	token string
}

func NewClient(config *Config) *Client {
	if config.HttpClient == nil {
		config.HttpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &Client{
		Alerts: alertClient{token: config.Alerts.Token, http: config.HttpClient},
	}
}
