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
	"net/url"
)

const (
	// Default User Agent for HTTP requests
	defaultUserAgent = "go-devo"
)

type Client struct {
	client    *http.Client
	UserAgent string

	Alerts         AlertsService
	AlertsEndpoint *url.URL
	AlertsToken    string
}

type ClientOpt func(*Client) error

func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{client: httpClient, UserAgent: defaultUserAgent}
	c.Alerts = &AlertsServiceOp{client: c}
	u, _ := url.Parse(ALERTS_API_US_DEFAULT_ENDPOINT)
	c.AlertsEndpoint = u

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func SetUserAgent(userAgent string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil

	}
}

func SetAlertsEndpoint(endpoint string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		c.AlertsEndpoint = u
		return nil

	}
}

func SetAlertsToken(token string) ClientOpt {
	return func(c *Client) error {
		c.AlertsToken = token
		return nil
	}
}
