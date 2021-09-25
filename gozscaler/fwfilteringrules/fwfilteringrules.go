package fwfilteringrules

import (
	"fmt"
	"net/http"
)

const (
	firewallPoliciesEndpoint = "/firewallFilteringRules"
)

type FirewallFilteringPolicies struct {
	ID                  string                `json:"id,omitempty"`
	Name                string                `json:"name,omitempty"`
	Order               string                `json:"order,omitempty"`
	Rank                string                `json:"rank,omitempty"`
	Locations           []Locations           `json:"locations"`
	LocationsGroups     []LocationsGroups     `json:"locationGroups"`
	Departments         []Departments         `json:"departments"`
	Users               []Users               `json:"users"`
	TimeWindows         []TimeWindows         `json:"timeWindows"`
	Action              string                `json:"action,omitempty"`
	State               string                `json:"state,omitempty"`
	Description         string                `json:"description,omitempty"`
	LastModifiedTime    string                `json:"lastModifiedTime,omitempty"`
	LastModifiedBy      []LastModifiedBy      `json:"lastModifiedBy"`
	SrcIps              string                `json:"srcIps,omitempty"`
	SrcIpGroups         []SrcIpGroups         `json:"srcIpGroups,omitempty"`
	DestAddresses       string                `json:"destAddresses,omitempty"`
	DestIpCategories    string                `json:"destIpCategories,omitempty"`
	DestCountries       string                `json:"destCountries,omitempty"`
	DestIpGroups        []DestIpGroups        `json:"destIpGroups"`
	NwServices          []NwServices          `json:"nwServices"`
	NwServiceGroups     []NwServiceGroups     `json:"nwServiceGroups"`
	NwApplications      string                `json:"nwApplications,omitempty"`
	NwApplicationGroups []NwApplicationGroups `json:"nwApplicationGroups"`
	AppServices         []AppServices         `json:"appServices"`
	AppServiceGroups    []AppServiceGroups    `json:"appServiceGroups"`
	Labels              []Labels              `json:"labels"`
	DefaultRule         bool                  `json:"defaultRule,omitempty"`
	Predefined          bool                  `json:"predefined,omitempty"`
}

// The locations to which the Firewall Filtering policy rule applies
// This is an immutable reference to an entity. which mainly consists of id and name
type Locations struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// The location groups to which the Firewall Filtering policy rule applies
type LocationsGroups struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// The departments to which the Firewall Filtering policy rule applies
type Departments struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// The users to which the Firewall Filtering policy rule applies
type Users struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// The time interval in which the Firewall Filtering policy rule applies
type TimeWindows struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

type LastModifiedBy struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// User-defined source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
type SrcIpGroups struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// User-defined destination IP address groups on which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group.
type DestIpGroups struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

type NwServices struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.
type NwServiceGroups struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// User-defined network service application group on which the rule is applied. If not set, the rule is not restricted to a specific network service application group.
type NwApplicationGroups struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// Application services on which this rule is applied
type AppServices struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// Application service groups on which this rule is applied
type AppServiceGroups struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// Application service groups on which this rule is applied
type Labels struct {
	ID         string                 `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

func (service *Service) Get(ruleId string) (*FirewallFilteringPolicies, *http.Response, error) {
	v := new(FirewallFilteringPolicies)
	relativeURL := fmt.Sprintf("%s/%s", firewallPoliciesEndpoint, ruleId)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Create(rules FirewallFilteringPolicies) (*FirewallFilteringPolicies, *http.Response, error) {
	v := new(FirewallFilteringPolicies)
	resp, err := service.Client.NewRequestDo("POST", firewallPoliciesEndpoint, nil, rules, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(ruleId string, rule FirewallFilteringPolicies) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", firewallPoliciesEndpoint, ruleId)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, ruleId, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(ruleId string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", firewallPoliciesEndpoint, ruleId)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
