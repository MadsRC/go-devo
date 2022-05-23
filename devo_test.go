// Copyright Mads R. Havmand.
// All Rights Reserved

package devo

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("nil Config", func(t *testing.T) {
		client := NewClient(nil)
		if client == nil {
			t.Fail()
		}
	})
	t.Run("default HTTP client", func(t *testing.T) {
		client := NewClient(&Config{})
		if client == nil {
			t.Fail()
		}
		if client.config.HTTP == nil {
			t.Fail()
		}
		if client.config.HTTP.Timeout != 10*time.Second {
			t.Fail()
		}
	})
	t.Run("custom HTTP client", func(t *testing.T) {
		client := NewClient(&Config{HTTP: &http.Client{Timeout: 20 * time.Second}})
		if client == nil {
			t.Fail()
		}
		if client.config.HTTP == nil {
			t.Fail()
		}
		if client.config.HTTP.Timeout != 20*time.Second {
			t.Fail()
		}
	})
}
