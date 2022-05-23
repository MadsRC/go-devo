// Copyright Mads R. Havmand.
// All Rights Reserved

package devo

import (
	"encoding/json"
	"net/http"
	"time"
)

type AlertsConfig struct {
	HTTP    *http.Client
	Token   string
	Address string
}

type AlertsClient struct {
	Config *AlertsConfig
}

func NewAlertsClient(config *AlertsConfig) *AlertsClient {
	if config == nil {
		config = &AlertsConfig{}
	}
	if config.HTTP == nil {
		config.HTTP = &http.Client{Timeout: 10 * time.Second}
	}

	if config.Address == "" {
		config.Address = ALERTS_API_US_DEFAULT_ENDPOINT
	}

	return &AlertsClient{Config: config}
}

type alertCorrelationTrigger struct {
	Kind              string `json:"kind"`
	ExternalOffset    string `json:"externalOffset"`
	InternalPeriod    string `json:"internalPeriod"`
	InternalOffset    string `json:"internalOffset"`
	Period            string `json:"period"`
	Threshold         string `json:"threshold"`
	BackPeriod        string `json:"backPeriod"`
	Absolute          string `json:"absolute"`
	AggregationColumn string `json:"aggregationColumn"`
}

type alertCorrelationContext struct {
	QuerySourceCode    string                  `json:"querySourceCode"`
	Priority           int                     `json:"priority"`
	CorrelationTrigger alertCorrelationTrigger `json:"correlationTrigger"`
}

type alert struct {
	ID                      string `json:"id"`
	CreationDate            int    `json:"creationDate"`
	Name                    string `json:"name"`
	Message                 string `json:"message"`
	Description             string `json:"description"`
	Subcategory             string `json:"subcategory"`
	CategoryID              string `json:"categoryId"`
	SubcategoryID           string `json:"subcategoryId"`
	IsActive                bool   `json:"isActive"`
	IsAlertChain            bool   `json:"isAlertChain"`
	AlertCorrelationContext alertCorrelationContext
	ActionPolicyID          []interface{} `json:"actionPolicyId"`
}

func listAlertDefinitions(client *http.Client, address string, token string) ([]alert, error) {
	request, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("content-type", "application/json")
	request.Header.Set("standAloneToken", token)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	alert := []alert{}
	err = json.NewDecoder(response.Body).Decode(&alert)
	if err != nil {
		return nil, err
	}
	return alert, nil
}
