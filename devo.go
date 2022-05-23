// Copyright Mads R. Havmand.
// All Rights Reserved

// Package devo provides a client library for working with the REST API
// interfaces available with https://devo.com.
//
// This package is NOT production ready as is.
//
// The official API documentation for the upstream Devo API's can be found at
// https://docs.devo.com/confluence/ndt/latest/api-reference.
package devo

import (
	"net/http"
	"time"
)

const (
	// Default endpoint for US based Devo domains.
	ALERTS_API_US_DEFAULT_ENDPOINT = "https://api-us.devo.com/alerts"

	// Default endpoint for EU based Devo domains.
	ALERTS_API_EU_DEFAULT_ENDPOINT = "https://api-eu.devo.com/alerts"
)

type Config struct {
	HTTP   *http.Client
	Alerts *AlertsConfig
}

type Client struct {
	config *Config
	Alerts *AlertsClient
}

func NewClient(config *Config) *Client {
	if config == nil {
		config = &Config{}
	}
	if config.HTTP == nil {
		config.HTTP = &http.Client{Timeout: 10 * time.Second}
	}

	alerts := NewAlertsClient(config.Alerts)

	return &Client{config: config, Alerts: alerts}
}
