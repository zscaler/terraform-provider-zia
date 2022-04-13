package urlfilteringpolicies

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zia/gozscaler/common"
)

const (
	urlFilteringPoliciesEndpoint = "/urlFilteringRules"
)

type URLFilteringRule struct {
	ID                     int                       `json:"id,omitempty"`
	Name                   string                    `json:"name,omitempty"`
	Order                  int                       `json:"order,omitempty"`
	Protocols              []string                  `json:"protocols,omitempty"`
	URLCategories          []string                  `json:"urlCategories"`
	State                  string                    `json:"state,omitempty"`
	UserAgentTypes         []string                  `json:"userAgentTypes,omitempty"`
	Rank                   int                       `json:"rank,omitempty"`
	RequestMethods         []string                  `json:"requestMethods,omitempty"`
	EndUserNotificationURL string                    `json:"endUserNotificationUrl,omitempty"`
	BlockOverride          bool                      `json:"blockOverride"`
	TimeQuota              int                       `json:"timeQuota,omitempty"`
	SizeQuota              int                       `json:"sizeQuota,omitempty"`
	Description            string                    `json:"description,omitempty"`
	ValidityStartTime      int                       `json:"validityStartTime,omitempty"`
	ValidityEndTime        int                       `json:"validityEndTime,omitempty"`
	ValidityTimeZoneID     string                    `json:"validityTimeZoneId,omitempty"`
	LastModifiedTime       int                       `json:"lastModifiedTime,omitempty"`
	EnforceTimeValidity    bool                      `json:"enforceTimeValidity"`
	Action                 string                    `json:"action,omitempty"`
	Ciparule               bool                      `json:"ciparule"`
	DeviceGroups           []common.IDNameExtensions `json:"deviceGroups"`
	Devices                []common.IDNameExtensions `json:"devices"`
	LastModifiedBy         *common.IDNameExtensions  `json:"lastModifiedBy,omitempty"`
	OverrideUsers          []common.IDNameExtensions `json:"overrideUsers,omitempty"`
	OverrideGroups         []common.IDNameExtensions `json:"overrideGroups,omitempty"`
	LocationGroups         []common.IDNameExtensions `json:"locationGroups,omitempty"`
	Labels                 []common.IDNameExtensions `json:"labels,omitempty"`
	Locations              []common.IDNameExtensions `json:"locations,omitempty"`
	Groups                 []common.IDNameExtensions `json:"groups,omitempty"`
	Departments            []common.IDNameExtensions `json:"departments,omitempty"`
	Users                  []common.IDNameExtensions `json:"users,omitempty"`
	TimeWindows            []common.IDNameExtensions `json:"timeWindows,omitempty"`
}

func (service *Service) Get(ruleID int) (*URLFilteringRule, error) {
	var rule URLFilteringRule
	err := service.Client.Read(fmt.Sprintf("%s/%d", urlFilteringPoliciesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning url filtering rules from Get: %d", rule.ID)
	return &rule, nil
}

func (service *Service) GetByName(ruleName string) (*URLFilteringRule, error) {
	var rules []URLFilteringRule
	err := service.Client.Read(urlFilteringPoliciesEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no url filtering rule found with name: %s", ruleName)
}

func (service *Service) Create(rule *URLFilteringRule) (*URLFilteringRule, error) {
	resp, err := service.Client.Create(urlFilteringPoliciesEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdURLFilteringRule, ok := resp.(*URLFilteringRule)
	if !ok {
		return nil, errors.New("object returned from api was not a url filtering rule pointer")
	}

	log.Printf("returning url filtering rule from create: %d", createdURLFilteringRule.ID)
	return createdURLFilteringRule, nil
}

func (service *Service) Update(ruleID int, rules *URLFilteringRule) (*URLFilteringRule, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", urlFilteringPoliciesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, nil, err
	}
	updatedURLFilteringRule, _ := resp.(*URLFilteringRule)

	log.Printf("returning url filtering rule from update: %d", updatedURLFilteringRule.ID)
	return updatedURLFilteringRule, nil, nil
}

func (service *Service) Delete(ruleID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", urlFilteringPoliciesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAll returns the all rules
func (service *Service) GetAll() ([]URLFilteringRule, error) {
	var urlFilteringPolicies []URLFilteringRule
	err := service.Client.Read(urlFilteringPoliciesEndpoint, &urlFilteringPolicies)
	if err != nil {
		return nil, err
	}
	return urlFilteringPolicies, nil
}

// RulesCount returns the number of rules
func (service *Service) RulesCount() int {
	rules, _ := service.GetAll()
	return len(rules)
}

// Reorder chanegs the order of the rule
func (service *Service) Reorder(ruleID, order int) (int, error) {
	resp, err := service.Get(ruleID)
	if err != nil {
		return 0, err
	}
	resp.Order = order
	_, _, err = service.Update(ruleID, resp)
	return order, err
}
