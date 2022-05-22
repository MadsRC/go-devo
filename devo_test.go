package devo

import (
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
		if client.Config.HTTP == nil {
			t.Fail()
		}
		if client.Config.HTTP.Timeout != 10 * time.Second {
			t.Fail()
		}
	})
	t.Run("default Alerts config", func(t *testing.T) {
		client := NewClient(&Config{})
		if client == nil {
			t.Fail()
		}
		if client.Config.Alerts.Address != ALERTS_API_US_DEFAULT_ENDPOINT {
			t.Fail()
		}
	})
}