package filteringrules

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zia/gozscaler/common"
)

const (
	firewallRulesEndpoint = "/firewallFilteringRules"
)

type FirewallFilteringRules struct {
	ID                  int                       `json:"id,omitempty"`
	Name                string                    `json:"name,omitempty"`
	Order               int                       `json:"order,omitempty"`
	Rank                int                       `json:"rank"`
	AccessControl       string                    `json:"accessControl,omitempty"`
	EnableFullLogging   bool                      `json:"enableFullLogging"`
	Action              string                    `json:"action,omitempty"`
	State               string                    `json:"state,omitempty"`
	Description         string                    `json:"description,omitempty"`
	LastModifiedTime    int                       `json:"lastModifiedTime,omitempty"`
	LastModifiedBy      *common.IDNameExtensions  `json:"lastModifiedBy,omitempty"`
	SrcIps              []string                  `json:"srcIps,omitempty"`
	DestAddresses       []string                  `json:"destAddresses,omitempty"`
	DestIpCategories    []string                  `json:"destIpCategories,omitempty"`
	DestCountries       []string                  `json:"destCountries,omitempty"`
	NwApplications      []string                  `json:"nwApplications,omitempty"`
	DefaultRule         bool                      `json:"defaultRule"`
	Predefined          bool                      `json:"predefined"`
	Locations           []common.IDNameExtensions `json:"locations,omitempty"`
	LocationsGroups     []common.IDNameExtensions `json:"locationGroups,omitempty"`      // The location groups to which the Firewall Filtering policy rule applies
	Departments         []common.IDNameExtensions `json:"departments,omitempty"`         // The departments to which the Firewall Filtering policy rule applies
	Groups              []common.IDNameExtensions `json:"groups,omitempty"`              // The groups to which the Firewall Filtering policy rule applies
	Users               []common.IDNameExtensions `json:"users,omitempty"`               // The users to which the Firewall Filtering policy rule applies
	TimeWindows         []common.IDNameExtensions `json:"timeWindows,omitempty"`         // The time interval in which the Firewall Filtering policy rule applies
	NwApplicationGroups []common.IDNameExtensions `json:"nwApplicationGroups,omitempty"` // User-defined network service application group on which the rule is applied. If not set, the rule is not restricted to a specific network service application group.
	AppServices         []common.IDNameExtensions `json:"appServices,omitempty"`         // Application services on which this rule is applied
	AppServiceGroups    []common.IDNameExtensions `json:"appServiceGroups,omitempty"`    // Application service groups on which this rule is applied
	Labels              []common.IDNameExtensions `json:"labels,omitempty"`              // Labels that are applicable to the rule.
	DestIpGroups        []common.IDNameExtensions `json:"destIpGroups,omitempty"`        // User-defined destination IP address groups on which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group.
	NwServices          []common.IDNameExtensions `json:"nwServices,omitempty"`
	NwServiceGroups     []common.IDNameExtensions `json:"nwServiceGroups,omitempty"` // User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.
	SrcIpGroups         []common.IDNameExtensions `json:"srcIpGroups,omitempty"`     // User-defined source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
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
