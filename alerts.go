// Copyright Mads R. Havmand.
// All Rights Reserved

package devo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// Type ListAlertDefinitionsParameters contains the parameters that can be used
// provided to the ListAlertDefinitions method.
//
// The documentation for the upstream Devo API can be found here:
// https://docs.devo.com/confluence/ndt/latest/api-reference/alerting-api/working-with-alert-definitions#id-.Workingwithalertdefinitionsvv7.11.0-createalertCreateanewalertdefinition
//
// Use these parameters to group your list of alerts by a specific number
// (size) and get only one of the resulting groups (page). This comes in
// handy if you need to get only a specific set of alerts and have a long
// list.
//
// Note that the count of both the selected page and groups defined starts
// at 0, so for example, if you enter page=2 and size=5 and have 22 alerts
// in your list, the API will divide the list into groups of 5 alerts
// (0-4, 5-9, 10-14, 15-19, and 20-22) and will return the group of alerts
// 10-14.
//
// Some struct attributes here is intentionally strings, and not ints, to allow
// us to distinguish between empty/no-value and 0.
type ListAlertDefinitionsParameters struct {
	// Define the group to get. See parent struct documentation.
	Page string

	// Define the number of alerts to get. See parent struct documentation.
	Size string

	// Use this parameter to filter alerts by their names. You will only get
	// alerts that contain the terms specified in their names. The filter is
	// case insensitive.
	NameFilter string

	// Indicate an alert definition ID to get only that specific alert. You
	// will get the ID of an alert definition after creating a new alert
	// definition through the Alerting API. Note that this ID cannot be
	// found in the Devo application.
	IDFilter string
}

func (client *AlertsClient) ListAlertDefinitions(parameters *ListAlertDefinitionsParameters) ([]alert, error) {
	address, err := url.Parse(fmt.Sprintf("%s/v1/alertDefinitions", strings.TrimRight(client.Config.Address, "/")))
	if err != nil {
		return nil, err
	}

	if parameters.Page != "" {
		address.Query().Add("page", parameters.Page)
	}
	if parameters.Size != "" {
		address.Query().Add("size", parameters.Size)
	}
	if parameters.NameFilter != "" {
		address.Query().Add("nameFilter", parameters.NameFilter)
	}
	if parameters.IDFilter != "" {
		address.Query().Add("idFilter", parameters.IDFilter)
	}
	return listAlertDefinitions(client.Config.HTTP, address.String(), client.Config.Token)
}

func listAlertDefinitions(client *http.Client, address string, token string) ([]alert, error) {
	request, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	addAlertAuthentication(request, token)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	alert := []alert{}
	err = json.NewDecoder(response.Body).Decode(&alert)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	return alert, nil
}

func (client *AlertsClient) CreateAlertDefinition(alert *alert) error {
	address, err := url.Parse(fmt.Sprintf("%s/v1/alertDefinitions", strings.TrimRight(client.Config.Address, "/")))
	if err != nil {
		return err
	}

	if alert == nil {
		return errors.New("Provided alert definition cannot be nil")
	}
	if alert.Name == "" {
		return errors.New("Empty name attribute in alert definition")
	}
	if alert.Subcategory == "" {
		return errors.New("Empty subcategory attribute in alert definition")
	}
	if alert.AlertCorrelationContext.QuerySourceCode == "" {
		return errors.New("Empty querySourceCode attribute in alertCorrelationContext in alert definition")
	}
	if alert.AlertCorrelationContext.CorrelationTrigger.Kind == "" {
		return errors.New("Empty Kind attribute in CorrelationTrigger in alert definition")
	}

	return createAlertDefinition(client.Config.HTTP, address.String(), client.Config.Token, alert)
}

func createAlertDefinition(client *http.Client, address string, token string, definition *alert) error {
	body, err := json.Marshal(definition)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, address, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Set("content-type", "application/json")
	addAlertAuthentication(request, token)

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&definition)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (client *AlertsClient) UpdateAlertDefinition(alert *alert) error {
	address, err := url.Parse(fmt.Sprintf("%s/v1/alertDefinitions", strings.TrimRight(client.Config.Address, "/")))
	if err != nil {
		return err
	}

	if alert == nil {
		return errors.New("Provided alert definition cannot be nil")
	}
	if alert.ID == "" {
		return errors.New("Empty ID attribute in alert definition")
	}

	return updateAlertDefinition(client.Config.HTTP, address.String(), client.Config.Token, alert)
}

func updateAlertDefinition(client *http.Client, address string, token string, definition *alert) error {
	body, err := json.Marshal(definition)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, address, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Set("content-type", "application/json")
	addAlertAuthentication(request, token)

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&definition)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func updateAlertDefinitionStatusBulk(client *http.Client, address string, token string) error {
	request, err := http.NewRequest(http.MethodPut, address, nil)
	if err != nil {
		return err
	}

	addAlertAuthentication(request, token)

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func (client *AlertsClient) UpdateAlertDefinitionStatusBulk(alerts []string, enable bool) error {
	address, err := url.Parse(fmt.Sprintf("%s/v1/alertDefinitions", strings.TrimRight(client.Config.Address, "/")))
	if err != nil {
		return err
	}

	for i := range alerts {
		address.Query().Add("alertIds", alerts[i])
	}

	address.Query().Add("enable", strconv.FormatBool(enable))

	return updateAlertDefinitionStatusBulk(client.Config.HTTP, address.String(), client.Config.Token)
}

func deleteAlertDefinitionBulk(client *http.Client, address string, token string) error {
	request, err := http.NewRequest(http.MethodDelete, address, nil)
	if err != nil {
		return err
	}

	addAlertAuthentication(request, token)

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func (client *AlertsClient) DeleteAlertDefinitionBulk(alerts []string) error {
	address, err := url.Parse(fmt.Sprintf("%s/v1/alertDefinitions", strings.TrimRight(client.Config.Address, "/")))
	if err != nil {
		return err
	}

	for i := range alerts {
		address.Query().Add("alertIds", alerts[i])
	}

	return deleteAlertDefinitionBulk(client.Config.HTTP, address.String(), client.Config.Token)
}

// Add authentication information specific to the Devo Alerts API.
func addAlertAuthentication(request *http.Request, token string) {
	if request == nil {
		return
	}

	request.Header.Set("standAloneToken", token)
}
