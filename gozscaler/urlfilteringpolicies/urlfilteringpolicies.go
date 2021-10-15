package urlfilteringpolicies

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	urlFilteringPoliciesEndpoint = "/urlFilteringRules"
)

type URLFilteringRule struct {
	ID                     int                `json:"id,omitempty"`
	Name                   string             `json:"name"`
	Order                  int                `json:"order,omitempty"`
	Protocols              []string           `json:"protocols,omitempty"`
	Locations              []IDNameExtensions `json:"locations,omitempty"`
	Groups                 []IDNameExtensions `json:"groups,omitempty"`
	Departments            []IDNameExtensions `json:"departments,omitempty"`
	Users                  []IDNameExtensions `json:"users,omitempty"`
	URLCategories          []string           `json:"urlCategories"`
	State                  string             `json:"state"`
	TimeWindows            []IDNameExtensions `json:"timeWindows"`
	Rank                   int                `json:"rank,omitempty"`
	RequestMethods         []string           `json:"requestMethods"`
	EndUserNotificationURL string             `json:"endUserNotificationUrl"`
	OverrideUsers          []IDNameExtensions `json:"overrideUsers,omitempty"`
	OverrideGroups         []IDNameExtensions `json:"overrideGroups,omitempty"`
	BlockOverride          bool               `json:"blockOverride,omitempty"`
	TimeQuota              int                `json:"timeQuota,omitempty"`
	SizeQuota              int                `json:"sizeQuota,omitempty"`
	Description            string             `json:"description"`
	LocationGroups         []IDNameExtensions `json:"locationGroups,omitempty"`
	Labels                 []IDNameExtensions `json:"labels,omitempty"`
	ValidityStartTime      int                `json:"validityStartTime"`
	ValidityEndTime        int                `json:"validityEndTime"`
	ValidityTimeZoneID     string             `json:"validityTimeZoneId"`
	LastModifiedTime       int                `json:"lastModifiedTime"`
	LastModifiedBy         *IDNameExtensions  `json:"lastModifiedBy,omitempty"`
	EnforceTimeValidity    bool               `json:"enforceTimeValidity,omitempty"`
	Action                 string             `json:"action"`
	Ciparule               bool               `json:"ciparule,omitempty"`
}

type IDNameExtensions struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func (service *Service) Get(ruleID int) (*URLFilteringRule, error) {
	var urlFilteringPolicies URLFilteringRule
	err := service.Client.Read(fmt.Sprintf("%s/%d", urlFilteringPoliciesEndpoint, ruleID), &urlFilteringPolicies)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning url filtering rules from Get: %d", urlFilteringPolicies.ID)
	return &urlFilteringPolicies, nil
}

func (service *Service) GetByName(urlFilteringPolicyName string) (*URLFilteringRule, error) {
	var urlFilteringPolicies []URLFilteringRule
	err := service.Client.Read(urlFilteringPoliciesEndpoint, &urlFilteringPolicies)
	if err != nil {
		return nil, err
	}
	for _, urlFilteringPolicy := range urlFilteringPolicies {
		if strings.EqualFold(urlFilteringPolicy.Name, urlFilteringPolicyName) {
			return &urlFilteringPolicy, nil
		}
	}
	return nil, fmt.Errorf("no url filtering rule found with name: %s", urlFilteringPolicyName)
}

func (service *Service) Create(ruleID *URLFilteringRule) (*URLFilteringRule, error) {
	resp, err := service.Client.Create(urlFilteringPoliciesEndpoint, *ruleID)
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
