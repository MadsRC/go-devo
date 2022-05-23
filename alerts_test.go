// Copyright Mads R. Havmand.
// All Rights Reserved

package devo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAlerts(t *testing.T) {
	t.Run("nil Config", func(t *testing.T) {
		client := NewAlertsClient(nil)
		if client == nil {
			t.Fail()
		}
	})
	t.Run("default HTTP client", func(t *testing.T) {
		client := NewAlertsClient(&AlertsConfig{})
		if client == nil {
			t.Fail()
		}
		if client.Config.HTTP == nil {
			t.Fail()
		}
		if client.Config.HTTP.Timeout != 10*time.Second {
			t.Fail()
		}
	})
	t.Run("default configs", func(t *testing.T) {
		client := NewAlertsClient(&AlertsConfig{})
		if client == nil {
			t.Fail()
		}
		if client.Config.Address != ALERTS_API_US_DEFAULT_ENDPOINT {
			t.Fail()
		}
	})
	t.Run("custom HTTP client", func(t *testing.T) {
		client := NewAlertsClient(&AlertsConfig{HTTP: &http.Client{Timeout: 20 * time.Second}})
		if client == nil {
			t.Fail()
		}
		if client.Config.HTTP == nil {
			t.Fail()
		}
		if client.Config.HTTP.Timeout != 20*time.Second {
			t.Fail()
		}
	})
	t.Run("listAlertDefinitions", func(t *testing.T) {
		responseData := `[
			{
				"id": "70736",
				"creationDate": 1604567644173,
				"name": "Alert_API_each",
				"message": "$eventdate $username - $count - API",
				"description": "Alert created by API",
				"categoryId": "7",
				"subcategory": "lib.my.testfake.AlertAPI_v660",
				"subcategoryId": "133",
				"isActive": true,
				"isFavorite": false,
				"isAlertChain": false,
				"alertCorrelationContext": {
					"id": "622",
					"nameId": "my.alert.testfake.Alert_API_each_Staging_1604567643224",
					"ownerEmail": "john@xx.com",
					"querySourceCode": "from siem.logtrust.web.activity group every 1m by username, url every 1m select count() as count",
					"priority": 5,
					"correlationTrigger": {
						"kind": "each"
					}
				},
				"actionPolicyId": []
			}
		]`
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Standalonetoken") != "nothinginparticular" {
				t.Fail()
			}
			fmt.Fprintf(w, responseData)
		}))
		defer svr.Close()
		_, err := listAlertDefinitions(&http.Client{}, svr.URL, "nothinginparticular")
		if err != nil {
			t.Errorf("%+v\n", err)
		}
	})
	t.Run("createAlertDefinition", func(t *testing.T) {
		responseData := `{
			"id": "70736",
			"creationDate": 1604567644173,
			"name": "Alert_API_each",
			"message": "$eventdate $username - $count - API",
			"description": "Alert created by API",
			"categoryId": "7",
			"subcategory": "lib.my.testfake.AlertAPI_v660",
			"subcategoryId": "133",
			"isActive": true,
			"isFavorite": false,
			"isAlertChain": false,
			"alertCorrelationContext": {
				"id": "622",
				"nameId": "my.alert.testfake.Alert_API_each_Staging_1604567643224",
				"ownerEmail": "john@xx.com",
				"querySourceCode": "from siem.logtrust.web.activity group every 1m by username, url every 1m select count() as count",
				"priority": 5,
				"correlationTrigger": {
					"kind": "each"
				}
			},
			"actionPolicyId": []
		}`
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Type") != "application/json" {
				t.Fail()
			}
			if r.Header.Get("Standalonetoken") != "nothinginparticular" {
				t.Fail()
			}
			fmt.Fprintf(w, responseData)
		}))
		defer svr.Close()
		alert := alert{
			Name:        "Alert_API_each",
			Subcategory: "lib.my.testfake.AlertAPI_v660",
			AlertCorrelationContext: alertCorrelationContext{
				QuerySourceCode: "from siem.logtrust.web.activity group every 1m by username, url every 1m select count() as count",
				CorrelationTrigger: alertCorrelationTrigger{
					Kind: "each",
				},
			},
		}
		err := createAlertDefinition(&http.Client{}, svr.URL, "nothinginparticular", &alert)
		if err != nil {
			t.Errorf("%+v\n", err)
		}
	})
}
