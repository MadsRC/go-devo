// Copyright Mads R. Havmand.
// All Rights Reserved

package devo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type AlertCorrelationTrigger struct {
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

type AlertCorrelationContext struct {
	QuerySourceCode    string                  `json:"querySourceCode"`
	Priority           int                     `json:"priority"`
	CorrelationTrigger AlertCorrelationTrigger `json:"correlationTrigger"`
}

type Alert struct {
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
	AlertCorrelationContext AlertCorrelationContext
	ActionPolicyID          []interface{} `json:"actionPolicyId"`
}

type AlertsService interface {
	List(parameters *AlertListRequest) ([]Alert, error)
	Create(createRequest *AlertCreateRequest) (*Alert, error)
	Update(updateRequest *AlertUpdateRequest) (*Alert, error)
	Delete(deleteRequest *AlertDeleteRequest) error
	Status(statusRequest *AlertStatusUpdateRequest) error
}

type AlertsServiceOp struct {
	client *Client
}

const (
	// Default endpoint for US based Devo domains.
	ALERTS_API_US_DEFAULT_ENDPOINT = "https://api-us.devo.com/alerts"

	// Default endpoint for EU based Devo domains.
	ALERTS_API_EU_DEFAULT_ENDPOINT = "https://api-eu.devo.com/alerts"

	// Default path for API Alerting
	ALERTS_API_PATH_ALERT_DEFINITIONS = "/v1/alertDefinitions"

	// Default path for API Alerting status
	ALERTS_API_PATH_ALERT_DEFINITIONS_STATUS = "/v1/alertDefinitions/status"
)

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
type AlertListRequest struct {
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

func (s *AlertsServiceOp) List(parameters *AlertListRequest) ([]Alert, error) {
	u, err := s.client.AlertsEndpoint.Parse(ALERTS_API_PATH_ALERT_DEFINITIONS)
	if err != nil {
		return nil, err
	}

	if parameters.Page != "" {
		query := u.Query()
		query.Add("page", parameters.Page)
		u.RawQuery = query.Encode()
	}
	if parameters.Size != "" {
		query := u.Query()
		query.Add("size", parameters.Size)
		u.RawQuery = query.Encode()
	}
	if parameters.NameFilter != "" {
		query := u.Query()
		query.Add("nameFilter", parameters.NameFilter)
		u.RawQuery = query.Encode()
	}
	if parameters.IDFilter != "" {
		query := u.Query()
		query.Add("idFilter", parameters.IDFilter)
		u.RawQuery = query.Encode()
	}

	request, err := alertsNewRequest(s.client, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	alerts := []Alert{}
	_, err = alertsDo(s.client, request, alerts)

	if err != nil {
		return nil, err
	}

	return alerts, nil
}

type AlertCreateRequest struct {
	Name                    string                  `json:"name"`
	Message                 string                  `json:"message,omitempty"`
	Description             string                  `json:"description,omitempty"`
	Subcategory             string                  `json:"subcategory"`
	AlertCorrelationContext AlertCorrelationContext `json:"alertCorrelationContext"`
}

func (s *AlertsServiceOp) Create(createRequest *AlertCreateRequest) (*Alert, error) {
	if createRequest == nil {
		return nil, errors.New("Create request cannot be empty")
	}

	u, err := s.client.AlertsEndpoint.Parse(ALERTS_API_PATH_ALERT_DEFINITIONS)
	if err != nil {
		return nil, err
	}

	request, err := alertsNewRequest(s.client, "POST", u.String(), createRequest)
	if err != nil {
		return nil, err
	}

	alert := Alert{}
	_, err = alertsDo(s.client, request, alert)

	if err != nil {
		return nil, err
	}

	return &alert, nil
}

type AlertUpdateRequest struct {
	Name                    string                  `json:"name"`
	Message                 string                  `json:"message,omitempty"`
	Description             string                  `json:"description,omitempty"`
	Subcategory             string                  `json:"subcategory,omitempty"`
	AlertCorrelationContext AlertCorrelationContext `json:"alertCorrelationContext,omitempty"`
}

func (s *AlertsServiceOp) Update(updateRequest *AlertUpdateRequest) (*Alert, error) {
	if updateRequest == nil {
		return nil, errors.New("Update request cannot be empty")
	}

	u, err := s.client.AlertsEndpoint.Parse(ALERTS_API_PATH_ALERT_DEFINITIONS)
	if err != nil {
		return nil, err
	}

	request, err := alertsNewRequest(s.client, "PUT", u.String(), updateRequest)
	if err != nil {
		return nil, err
	}

	alert := Alert{}
	_, err = alertsDo(s.client, request, alert)

	if err != nil {
		return nil, err
	}

	return &alert, nil
}

type AlertDeleteRequest struct {
	AlertIDs []string
}

func (s *AlertsServiceOp) Delete(deleteRequest *AlertDeleteRequest) error {
	if deleteRequest == nil {
		return errors.New("Delete request cannot be empty")
	}
	if len(deleteRequest.AlertIDs) < 1 {
		return errors.New("no alert IDs in delete request")
	}

	u, err := s.client.AlertsEndpoint.Parse(ALERTS_API_PATH_ALERT_DEFINITIONS)
	if err != nil {
		return err
	}

	for i := range deleteRequest.AlertIDs {
		query := u.Query()
		query.Add("alertIds", deleteRequest.AlertIDs[i])
		u.RawQuery = query.Encode()
	}

	_, err = alertsNewRequest(s.client, "DELETE", u.String(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

type AlertStatusUpdateRequest struct {
	AlertIDs []string
	Enable   bool
}

func (s *AlertsServiceOp) Status(statusRequest *AlertStatusUpdateRequest) error {
	if statusRequest == nil {
		return errors.New("Delete request cannot be empty")
	}
	if len(statusRequest.AlertIDs) < 1 {
		return errors.New("no alert IDs in delete request")
	}

	u, err := s.client.AlertsEndpoint.Parse(ALERTS_API_PATH_ALERT_DEFINITIONS_STATUS)
	if err != nil {
		return err
	}

	for i := range statusRequest.AlertIDs {
		query := u.Query()
		query.Add("alertIds", statusRequest.AlertIDs[i])
		u.RawQuery = query.Encode()
	}

	query := u.Query()
	query.Add("enable", strconv.FormatBool(statusRequest.Enable))
	u.RawQuery = query.Encode()

	_, err = alertsNewRequest(s.client, "PUT", u.String(), statusRequest)
	if err != nil {
		return err
	}

	return nil
}

// alertsNewRequest create a new API request for the Alerting API. A relative URL can be provided in urlStr
// which will be resolved to the AlertsEndpoint of the Client. Relative URLS should always be specified without
// a preceding slash. If specified, the value pointed to by body is JSON encoded and included in as the
// request body.
func alertsNewRequest(client *Client, method string, urlStr string, body interface{}) (*http.Request, error) {
	u, err := client.AlertsEndpoint.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var request *http.Request
	switch method {
	case http.MethodGet:
		request, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		request, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Content-Type", "application/json")
	}

	request.Header.Set("User-Agent", client.UserAgent)
	request.Header.Set("standAloneToken", client.AlertsToken)

	return request, nil
}

// alertsDo will send an API request to the Alerting API. The API response is JSON decoded and
// stored in the value pointed to by value. If value implements the io.Writer interface, the
// raw response will be written to value, without attempting to decode it.
func alertsDo(client *Client, request *http.Request, value interface{}) (*http.Response, error) {
	response, err := client.client.Do(request)
	if err != nil {
		return nil, err
	}

	if value != nil {
		w, ok := value.(io.Writer)
		if ok {
			_, err = io.Copy(w, response.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(response.Body).Decode(value)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}
