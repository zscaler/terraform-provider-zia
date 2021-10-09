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
	ID                     int              `json:"id"`
	Name                   string           `json:"name"`
	Order                  int              `json:"order"`
	Protocols              []string         `json:"protocols"`
	Locations              []Locations      `json:"locations"`
	Groups                 []Groups         `json:"groups"`
	Departments            []Departments    `json:"departments"`
	Users                  []Users          `json:"users"`
	URLCategories          []string         `json:"urlCategories"`
	State                  string           `json:"state"`
	TimeWindows            []TimeWindows    `json:"timeWindows"`
	Rank                   int              `json:"rank"`
	RequestMethods         []string         `json:"requestMethods"`
	EndUserNotificationURL string           `json:"endUserNotificationUrl"`
	OverrideUsers          []OverrideUsers  `json:"overrideUsers"`
	OverrideGroups         []OverrideGroups `json:"overrideGroups"`
	BlockOverride          bool             `json:"blockOverride"`
	TimeQuota              int              `json:"timeQuota"`
	SizeQuota              int              `json:"sizeQuota"`
	Description            string           `json:"description"`
	LocationGroups         []LocationGroups `json:"locationGroups"`
	Labels                 []Labels         `json:"labels"`
	ValidityStartTime      int              `json:"validityStartTime"`
	ValidityEndTime        int              `json:"validityEndTime"`
	ValidityTimeZoneID     string           `json:"validityTimeZoneId"`
	LastModifiedTime       int              `json:"lastModifiedTime"`
	LastModifiedBy         *LastModifiedBy  `json:"lastModifiedBy"`
	EnforceTimeValidity    bool             `json:"enforceTimeValidity"`
	Action                 string           `json:"action"`
	Ciparule               bool             `json:"ciparule"`
}

type Locations struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type Groups struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type Departments struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type Users struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type TimeWindows struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type OverrideUsers struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type OverrideGroups struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type LocationGroups struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type Labels struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type LastModifiedBy struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
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
