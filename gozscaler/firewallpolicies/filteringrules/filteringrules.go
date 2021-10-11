package filteringrules

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	firewallRulesEndpoint = "/firewallFilteringRules"
)

type FirewallFilteringRules struct {
	ID                  int                   `json:"id,omitempty"`
	Name                string                `json:"name,omitempty"`
	Order               int                   `json:"order,omitempty"`
	Rank                int                   `json:"rank,omitempty"`
	AccessControl       string                `json:"accessControl,omitempty"`
	EnableFullLogging   bool                  `json:"enableFullLogging"`
	Locations           []Locations           `json:"locations"`
	LocationsGroups     []LocationsGroups     `json:"locationGroups"`
	Departments         []Departments         `json:"departments"`
	Groups              []Groups              `json:"groups"`
	Users               []Users               `json:"users"`
	TimeWindows         []TimeWindows         `json:"timeWindows"`
	Action              string                `json:"action,omitempty"`
	State               string                `json:"state,omitempty"`
	Description         string                `json:"description,omitempty"`
	LastModifiedTime    int                   `json:"lastModifiedTime,omitempty"`
	LastModifiedBy      []LastModifiedBy      `json:"lastModifiedBy"`
	SrcIps              []string              `json:"srcIps,omitempty"`
	SrcIpGroups         []SrcIpGroups         `json:"srcIpGroups,omitempty"`
	DestAddresses       []string              `json:"destAddresses,omitempty"`
	DestIpCategories    []string              `json:"destIpCategories,omitempty"`
	DestCountries       []string              `json:"destCountries,omitempty"`
	DestIpGroups        []DestIpGroups        `json:"destIpGroups"`
	NwServices          []NwServices          `json:"nwServices"`
	NwServiceGroups     []NwServiceGroups     `json:"nwServiceGroups"`
	NwApplications      []string              `json:"nwApplications,omitempty"`
	NwApplicationGroups []NwApplicationGroups `json:"nwApplicationGroups"`
	AppServices         []AppServices         `json:"appServices"`
	AppServiceGroups    []AppServiceGroups    `json:"appServiceGroups"`
	Labels              []Labels              `json:"labels"`
	DefaultRule         bool                  `json:"defaultRule"`
	Predefined          bool                  `json:"predefined"`
}

// The locations to which the Firewall Filtering policy rule applies
// This is an immutable reference to an entity. which mainly consists of id and name
type Locations struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// The location groups to which the Firewall Filtering policy rule applies
type LocationsGroups struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// The departments to which the Firewall Filtering policy rule applies
type Departments struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// The groups to which the Firewall Filtering policy rule applies
type Groups struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// The users to which the Firewall Filtering policy rule applies
type Users struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// The time interval in which the Firewall Filtering policy rule applies
type TimeWindows struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

type LastModifiedBy struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// User-defined source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
type SrcIpGroups struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// User-defined destination IP address groups on which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group.
type DestIpGroups struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}
type NwServices struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.
type NwServiceGroups struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// User-defined network service application group on which the rule is applied. If not set, the rule is not restricted to a specific network service application group.
type NwApplicationGroups struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// Application services on which this rule is applied
type AppServices struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// Application service groups on which this rule is applied
type AppServiceGroups struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// Labels that are applicable to the rule.
type Labels struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

func (service *Service) Get(ruleID int) (*FirewallFilteringRules, error) {
	var rule FirewallFilteringRules
	err := service.Client.Read(fmt.Sprintf("%s/%d", firewallRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning firewall rule from Get: %d", rule.ID)
	return &rule, nil
}

func (service *Service) GetByName(ruleName string) (*FirewallFilteringRules, error) {
	var rules []FirewallFilteringRules
	err := service.Client.Read(firewallRulesEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no firewall rule found with name: %s", ruleName)
}

func (service *Service) Create(rule *FirewallFilteringRules) (*FirewallFilteringRules, error) {
	resp, err := service.Client.Create(firewallRulesEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*FirewallFilteringRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	log.Printf("returning rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func (service *Service) Update(ruleID int, rules *FirewallFilteringRules) (*FirewallFilteringRules, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", firewallRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*FirewallFilteringRules)
	log.Printf("returning firewall rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func (service *Service) Delete(ruleID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", firewallRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
