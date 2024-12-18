package client

import (
	"bytes"
	"fmt"
	"github.com/thoas/go-funk"
	"net/http"
	"net/url"
	"os"

	json "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

const Blanc = ""

const ContentType = "Content-Type"
const ApplicationJSON = "application/json"

const BaseURL = "https://uptime.betterstack.com"
const APIV2Group = BaseURL + "/api/v2"
const Monitors = APIV2Group + "/monitors"
const MonitorID = APIV2Group + "/monitors/%s"
const MonitorGroupID = APIV2Group + "/monitor-groups/%s"
const MonitorGroups = APIV2Group + "/monitor-groups"

type BetterstackClient struct {
	headers http.Header
}

func NewClient(apiToken string) *BetterstackClient {
	var headers = getDefaultHeaders()
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	return &BetterstackClient{
		headers: headers,
	}
}

func NewClientFromENV() *BetterstackClient {
	var token = os.Getenv("BETTERSTACK_TOKEN")
	if funk.IsEmpty(token) {
		log.Fatal("BETTERSTACK_TOKEN environment variable not set")
	}
	var headers = getDefaultHeaders()
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return &BetterstackClient{
		headers: headers,
	}
}

func (c *BetterstackClient) ListMonitors(page int, filterType, filterValue string) (MonitorsResponse, error) {
	var result MonitorsResponse

	if page < 1 {
		page = 1
	}

	params := url.Values{}
	params.Add("per_page", "250")
	params.Add("page", fmt.Sprintf("%d", page))

	if funk.NotEmpty(filterType) && funk.NotEmpty(filterValue) {
		switch filterType {
		case "url":
			params.Add("url", filterValue)
		case "pronounceable_name":
			params.Add("pronounceable_name", filterValue)
		default:
			return result, fmt.Errorf("invalid filter type: %s", filterType)
		}
	}

	var targetURL = fmt.Sprintf("%s?%s", Monitors, params.Encode())

	var monitorsRequest, monsErr = http.NewRequest(http.MethodGet, targetURL, nil)
	if monsErr != nil {
		return result, fmt.Errorf("failed to create request: %v", monsErr)
	}

	monitorsRequest.Header = c.headers

	var monitorsResponse, monsRespErr = http.DefaultClient.Do(monitorsRequest)
	if monsRespErr != nil {
		return result, fmt.Errorf("failed to execute request: %v", monsRespErr)
	}

	var unmErr = json.NewDecoder(monitorsResponse.Body).Decode(&result)
	if unmErr != nil {
		return result, fmt.Errorf("failed to unmarshal response: %v", unmErr)
	}

	if funk.NotEmpty(result.Errors) {
		return result, fmt.Errorf("failed to list monitors: %v", result.Errors)
	}

	return result, nil
}

func (c *BetterstackClient) ListAllMonitors() ([]Monitor, error) {
	var result []Monitor
	var monitorResponses MonitorsResponse
	var monsErr error
	var page = 1

	monitorResponses, monsErr = c.ListMonitors(page, Blanc, Blanc)
	if monsErr != nil {
		return result, monsErr
	}

	var lastPage, paginationErr = monitorResponses.Pagination.GetLastPage()
	if paginationErr != nil {
		return result, paginationErr
	}

	for _, mon := range monitorResponses.Data {
		mon.Attributes.ID = mon.ID
		result = append(result, mon.Attributes)
	}

	if page == lastPage {
		return result, nil
	}

	page++

	for i := page; i <= lastPage; i++ {
		tempMonitors, tempErr := c.ListMonitors(i, Blanc, Blanc)
		if tempErr != nil {
			return result, tempErr
		}
		for _, mon := range tempMonitors.Data {
			mon.Attributes.ID = mon.ID
			result = append(result, mon.Attributes)
		}
	}

	return result, nil
}

func (c *BetterstackClient) FindMonitor(kind, val string) ([]Monitor, error) {
	var result []Monitor
	var monitorResponses MonitorsResponse
	var monsErr error
	var page = 1

	monitorResponses, monsErr = c.ListMonitors(page, kind, val)
	if monsErr != nil {
		return result, monsErr
	}

	for _, mon := range monitorResponses.Data {
		mon.Attributes.ID = mon.ID
		result = append(result, mon.Attributes)
	}

	return result, nil
}

func (c *BetterstackClient) CreateMonitor(monitor Monitor) (MonitorResponse, error) {
	var result MonitorResponse
	var serializedBody, serErr = json.Marshal(monitor)
	if serErr != nil {
		return result, serErr
	}

	var postBody = bytes.NewReader(serializedBody)

	var monitorRequest, monsErr = http.NewRequest(http.MethodPost, Monitors, postBody)
	if monsErr != nil {
		return result, fmt.Errorf("failed to create request: %v", monsErr)
	}

	monitorRequest.Header = c.headers

	var monitorResponse, monsRespErr = http.DefaultClient.Do(monitorRequest)
	if monsRespErr != nil || monitorResponse.StatusCode != http.StatusCreated {
		return result, fmt.Errorf("failed to execute request: %v", monsRespErr)
	}

	var unmErr = json.NewDecoder(monitorResponse.Body).Decode(&result)
	if unmErr != nil {
		return result, fmt.Errorf("failed to unmarshal response: %v", unmErr)
	}

	if funk.NotEmpty(result.Errors) {
		return result, fmt.Errorf("failed to list monitors: %v", result.Errors)
	}

	result.Data.Attributes.ID = result.Data.ID

	return result, nil
}

func (c *BetterstackClient) GetMonitor(id string) (MonitorResponse, error) {
	var result MonitorResponse
	var targetURL = fmt.Sprintf(MonitorID, id)

	var monitorRequest, monErr = http.NewRequest(http.MethodGet, targetURL, nil)
	if monErr != nil {
		return result, fmt.Errorf("failed to create request: %v", monErr)
	}

	monitorRequest.Header = c.headers

	var monitorResponse, monRespErr = http.DefaultClient.Do(monitorRequest)
	if monRespErr != nil {
		return result, fmt.Errorf("failed to execute request: %v", monRespErr)
	}

	var unmErr = json.NewDecoder(monitorResponse.Body).Decode(&result)
	if unmErr != nil {
		return result, fmt.Errorf("failed to unmarshal response: %v", unmErr)
	}

	if funk.NotEmpty(result.Errors) {
		return result, fmt.Errorf("failed to get monitor: %v", result.Errors)
	}

	result.Data.Attributes.ID = result.Data.ID

	return result, nil
}

func (c *BetterstackClient) UpdateMonitor(id string, monitor Monitor) (MonitorResponse, error) {
	var result MonitorResponse
	var serializedBody, serErr = json.Marshal(monitor)
	if serErr != nil {
		return result, serErr
	}

	var postBody = bytes.NewReader(serializedBody)
	var targetURL = fmt.Sprintf(MonitorID, id)

	var monitorRequest, monErr = http.NewRequest(http.MethodPatch, targetURL, postBody)
	if monErr != nil {
		return result, fmt.Errorf("failed to create request: %v", monErr)
	}

	monitorRequest.Header = c.headers

	var monitorResponse, monRespErr = http.DefaultClient.Do(monitorRequest)
	if monRespErr != nil {
		return result, fmt.Errorf("failed to execute request: %v", monRespErr)
	}

	if monitorResponse.StatusCode != http.StatusOK {
		return result, fmt.Errorf("failed to execute request: %v", monitorResponse.Status)
	}

	var unmErr = json.NewDecoder(monitorResponse.Body).Decode(&result)
	if unmErr != nil {
		return result, fmt.Errorf("failed to unmarshal response: %v", unmErr)
	}

	if funk.NotEmpty(result.Errors) {
		return result, fmt.Errorf("failed to update monitor: %v", result.Errors)
	}

	result.Data.Attributes.ID = result.Data.ID
	return result, nil
}

func (c *BetterstackClient) DeleteMonitor(id string) error {

	var targetURL = fmt.Sprintf(MonitorID, id)

	var monitorRequest, monErr = http.NewRequest(http.MethodDelete, targetURL, nil)
	if monErr != nil {
		return fmt.Errorf("failed to create request: %v", monErr)
	}

	monitorRequest.Header = c.headers

	var monitorResponse, monRespErr = http.DefaultClient.Do(monitorRequest)
	if monRespErr != nil || monitorResponse.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to execute request: %v", monRespErr)
	}

	return nil
}

func getDefaultHeaders() http.Header {
	var headers = http.Header{}
	headers.Add(ContentType, ApplicationJSON)
	return headers
}
