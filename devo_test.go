// Copyright Mads R. Havmand.
// All Rights Reserved

package devo

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestDefaultAlertsEndpoints(t *testing.T) {
	u, err := url.Parse(ALERTS_API_US_DEFAULT_ENDPOINT)
	if err != nil {
		t.Fatalf("Error while parsing ALERTS_API_US_DEFAULT_ENDPOINT: %v", err)
	}
	if u.String() != ALERTS_API_US_DEFAULT_ENDPOINT {
		t.Fatalf("Parsed ALERTS_API_US_DEFAULT_ENDPOINT does not match ALERTS_API_US_DEFAULT_ENDPOINT")
	}
	u, err = url.Parse(ALERTS_API_EU_DEFAULT_ENDPOINT)
	if err != nil {
		t.Fatalf("Error while parsing ALERTS_API_EU_DEFAULT_ENDPOINT: %v", err)
	}
	if u.String() != ALERTS_API_EU_DEFAULT_ENDPOINT {
		t.Fatalf("Parsed ALERTS_API_EU_DEFAULT_ENDPOINT does not match ALERTS_API_EU_DEFAULT_ENDPOINT")
	}
}

func testClientServices(t *testing.T, c *Client) {
	services := []string{
		"Alerts",
	}

	cp := reflect.ValueOf(c)
	cv := reflect.Indirect(cp)

	for _, s := range services {
		if cv.FieldByName(s).IsNil() {
			t.Errorf("c.%s shouldn't be nil", s)
		}
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	if c.UserAgent != defaultUserAgent {
		t.Errorf("Client UserAgent = %v, expected %v", c.UserAgent, defaultUserAgent)
	}
}

func testClientDefaultHTTPClient(t *testing.T, c *Client) {
	if c.client == nil {
		t.Errorf("Client HTTP client is nil, expected %v", http.DefaultClient)
	}
	if c.client != http.DefaultClient {
		t.Errorf("Client HTTP client = %v, expected %v", c.client, http.DefaultClient)
	}
}

func testClientDefaults(t *testing.T, c *Client) {
	testClientDefaultUserAgent(t, c)
	testClientDefaultHTTPClient(t, c)
	testClientServices(t, c)
}

func TestNewClient(t *testing.T) {
	c, err := New(nil)
	if err != nil {
		t.Fatalf("New(): %v", err)
	}
	testClientDefaults(t, c)
}

func TestCustomUserAgent(t *testing.T) {
	userAgent := "iLikeCake/42.0.0"
	c, err := New(nil, SetUserAgent(userAgent))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if c.UserAgent != userAgent {
		t.Errorf("New() UserAgent = %s; expected %s", c.UserAgent, userAgent)
	}
}
