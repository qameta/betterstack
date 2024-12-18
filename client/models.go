package client

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/thoas/go-funk"
)

// Monitors

type Monitor struct {

	// ID represents the unique identifier for the Monitor, used to distinguish and reference the monitor entity.
	// Do not use on creation
	ID string `json:"id,omitempty"`

	// Required if using global API token to specify the team which should own the resource
	TeamName string `json:"team_name,omitempty"`

	// Valid values:
	// status — We will check your website for a 2XX HTTP status code.
	// expected_status_code — We will check if your website returned one of the values in expected_status_codes.
	// keyword — We will check if your website contains the required_keyword.
	// keyword_absence — We will check if your website doesn't contain the required_keyword.
	// ping — We will ping your host specified in the url parameter.
	// tcp — We will test a TCP port at your host specified in the url parameter (port is required).
	// udp — We will test a UDP port at your host specified in the url parameter (port and required_keyword are
	// required).
	// smtp — We will check for SMTP server at the host specified in the url parameter (port is required, and can be
	// one of 25, 465, 587, or a combination of those ports separated by a comma).
	// pop — We will check for a POP3 server at the host specified in the url parameter (port is required, and can be
	// 110, 995, or both).
	// imap — We will check for an IMAP server at the host specified in the url parameter (port is required, and can be
	// 143, 993, or both).
	// dns — We will check for a DNS server at the host specified in the url parameter (request_body is required, and
	// should contain the domain to query the DNS server with).
	// playwright — We will run the scenario defined by playwright_script, identified in the UI by scenario_name.
	MonitorType string `json:"monitor_type"`

	// The URL of your website or the host you want to ping. See monitor_type below.
	URL string `json:"url"`

	// The name of the monitor
	PronounceableName string `json:"pronounceable_name"`

	// Send email alerts
	Email bool `json:"email"`

	// Send SMS alerts
	SMS bool `json:"sms"`

	// Phone call alerts
	Call bool `json:"call"`

	// Should we send a push notification to the on-call person?
	Push bool `json:"push"`

	// Check frequency (in seconds)
	CheckFrequency int `json:"check_frequency,omitempty"`

	// The request headers that will be sent with the check
	RequestHeaders []RequestHeader `json:"request_headers,omitempty"`

	// An array of status codes you expect to receive from your website. These status codes are considered only if the
	// monitor_type is expected_status_code.
	ExpectedStatusCodes []int `json:"expected_status_codes,omitempty"`

	// How many days before the domain expires do you want to be alerted? Valid values are 1, 2, 3, 7, 14, 30, and 60.
	DomainExpiration int `json:"domain_expiration,omitempty"`

	// How many days before the SSL certificate expires do you want to be alerted? Valid values are 1, 2, 3, 7, 14, 30,
	// and 60.
	SSLExpiration int `json:"ssl_expiration,omitempty"`

	// Set the escalation policy for the monitor.
	PolicyID string `json:"policy_id,omitempty"`

	// Should we automatically follow redirects when sending the HTTP request?
	FollowRedirects bool `json:"follow_redirects,omitempty"`

	// Required if monitor_type is set to keyword or udp. We will create a new incident if this keyword is missing
	// on your page.
	RequiredKeyword string `json:"required_keyword,omitempty"`

	// How long to wait before escalating the incident alert to the team. Leave blank to disable escalating to the
	// entire team. In seconds.
	TeamWait int `json:"team_wait,omitempty"`

	// Set to true to pause monitoring — we won't notify you about downtime. Set to 'false' to resume monitoring.
	Paused bool `json:"paused,omitempty"`

	// Required if monitor_type is set to tcp, udp, smtp, pop, or imap. tcp and udp monitors accept any ports,
	// while smtp, pop, and imap accept only the specified ports corresponding with their
	// servers (e.g. 25,465,587 for smtp).
	Port int `json:"port,omitempty"`

	// An array of regions to set. Allowed values are ['us', 'eu', 'as', 'au'] or any subset of these regions.
	Regions []string `json:"regions"`

	// Set this attribute if you want to add this monitor to a monitor group.
	MonitorGroupID any `json:"monitor_group_id,omitempty"`

	// How long the monitor must be up to automatically mark an incident as resolved after being down. In seconds.
	RecoveryPeriod int `json:"recovery_period,omitempty"`

	// Should we verify SSL certificate validity?
	VerifySSL bool `json:"verify_ssl,omitempty"`

	// How long should we wait after observing a failure before we start a new incident? In seconds.
	ConfirmationPeriod int `json:"confirmation_period,omitempty"`

	// HTTP Method used to make a request. Valid options: GET, HEAD, POST, PUT, PATCH
	HTTPMethod string `json:"http_method,omitempty"`

	// How long to wait before timing out the request? In seconds. When monitor_type is set to playwright,
	// this determines the Playwright scenario timeout instead.
	RequestTimeout int `json:"request_timeout,omitempty"`

	// Request body for POST, PUT, PATCH requests. Required if monitor_type is set to dns
	// (domain to query the DNS server with).
	RequestMethod string `json:"request_method,omitempty"`

	// Basic HTTP authentication username to include with the request.
	AuthUsername string `json:"auth_username,omitempty"`

	// Basic HTTP authentication password to include with the request.
	AuthPassword string `json:"auth_password,omitempty"`

	// An array of maintenance days to set. If a maintenance window is overnight both affected days should be set.
	// Allowed values are ['mon', 'tue', 'wed', 'thu', 'fri', 'sat', 'sun'] or any subset of these days.
	MaintenanceDays []string `json:"maintenance_days,omitempty"`

	// Start of the maintenance window each day. We won't check your website during this window. Example: '01:00:00'
	MaintenanceFrom string `json:"maintenance_from,omitempty"`

	// End of the maintenance window each day. Example: '03:00:00'
	MaintenanceTo string `json:"maintenance_to,omitempty"`

	// The timezone to use for the maintenance window each day. Defaults to UTC. The accepted values can be found
	// in the Rails TimeZone documentation. https://api.rubyonrails.org/classes/ActiveSupport/TimeZone.html
	MaintenanceTimezone string `json:"maintenance_timezone,omitempty"`

	// Do you want to keep cookies when redirecting?
	RememberCookies bool `json:"remember_cookies,omitempty"`

	// For Playwright monitors, the JavaScript source code of the scenario.
	PlaywrightScript string `json:"playwright_script,omitempty"`

	// For Playwright monitors, the scenario name identifying the monitor in the UI.
	ScenarioName string `json:"scenario_name,omitempty"`

	// Status represents the current status of the monitor, indicating its operational state or health.
	Status string `json:"status,omitempty"`
}

// Monitor Groups

type MonitorGroup struct {
	Name      string     `json:"name"`
	TeamName  string     `json:"team_name"`
	SortIndex int        `json:"sort_index"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Paused    bool       `json:"paused"`
}

// REST Models

type MonitorResponse ResponseWrapper[Monitor]
type MonitorsResponse ListWrapper[Monitor]
type MonitorGroupResponse ResponseWrapper[MonitorGroup]
type MonitorGroupsResponse ListWrapper[MonitorGroup]

// Commons

type ResponseWrapper[T Monitor | MonitorGroup] struct {
	Data       EntityWrapper[T] `json:"data,omitempty"`
	Errors     any              `json:"errors,omitempty"`
	Pagination Pagination       `json:"pagination,omitempty"`
}

type ListWrapper[T Monitor | MonitorGroup] struct {
	Data       []EntityWrapper[T] `json:"data,omitempty"`
	Errors     any                `json:"errors,omitempty"`
	Pagination Pagination         `json:"pagination,omitempty"`
}

type EntityWrapper[T Monitor | MonitorGroup] struct {
	ID         string `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	Attributes T      `json:"attributes,omitempty"`
}

type Pagination struct {
	First    string `json:"first"`
	Last     string `json:"last"`
	Previous string `json:"prev"`
	Next     string `json:"next"`
}

func (p *Pagination) HasNext() bool {
	return funk.NotEmpty(p.Next)
}

func (p *Pagination) HasPrevious() bool {
	return funk.NotEmpty(p.Previous)
}

// Useful when you need to iterate over all pages collecting entities

func (p *Pagination) GetLastPage() (int, error) {
	parsedURL, err := url.Parse(p.Last)
	if err != nil {
		return 0, fmt.Errorf("failed to parse URL: %v", err)
	}

	queryParams := parsedURL.Query()
	pageStr := queryParams.Get("page")
	if funk.IsEmpty(pageStr) {
		return 0, fmt.Errorf("'page' parameter not found in URL")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'page' value: %v", err)
	}

	return page, nil
}

type RequestHeader struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
