package urlfilteringpolicies

import (
	"fmt"
	"log"
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
	LastModifiedBy         LastModifiedBy   `json:"lastModifiedBy"`
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

func (service *Service) GetURLFilteringRules(urlFilteringPoliciesID int) (*URLFilteringRule, error) {
	var urlFilteringPolicies URLFilteringRule
	err := service.Client.Read(fmt.Sprintf("%s/%d", urlFilteringPoliciesEndpoint, urlFilteringPoliciesID), &urlFilteringPolicies)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning url filtering rules from Get: %d", urlFilteringPolicies.ID)
	return &urlFilteringPolicies, nil
}

func (service *Service) GetURLFilteringRulesByName(urlFilteringPolicyName string) (*URLFilteringRule, error) {
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
