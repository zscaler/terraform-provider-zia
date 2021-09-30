package urlcategories

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

const (
	urlCategoriesEndpoint = "/urlCategories"
)

type URLCategoryInformation struct {
	ID                               string           `json:"id"`
	ConfiguredName                   string           `json:"configuredName"`
	Urls                             []string         `json:"urls"`
	DBCategorizedUrls                []string         `json:"dbCategorizedUrls"`
	CustomCategory                   bool             `json:"customCategory"`
	Scopes                           []Scopes         `json:"scopes"`
	Editable                         bool             `json:"editable"`
	Description                      string           `json:"description"`
	Type                             string           `json:"type"`
	URLKeywordCounts                 URLKeywordCounts `json:"urlKeywordCounts"`
	Val                              int              `json:"val"`
	CustomUrlsCount                  int              `json:"customUrlsCount"`
	UrlsRetainingParentCategoryCount int              `json:"urlsRetainingParentCategoryCount"`
}
type Scopes struct {
	ScopeGroupMemberEntities []ScopeGroupMemberEntities `json:"scopeGroupMemberEntities"`
	Type                     string                     `json:"Type"`
	ScopeEntities            []ScopeEntities            `json:"ScopeEntities"`
}
type ScopeGroupMemberEntities struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}
type ScopeEntities struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions"`
}

type URLKeywordCounts struct {
	TotalURLCount            int `json:"totalUrlCount"`
	RetainParentURLCount     int `json:"retainParentUrlCount"`
	TotalKeywordCount        int `json:"totalKeywordCount"`
	RetainParentKeywordCount int `json:"retainParentKeywordCount"`
}

func (service *Service) GetURLCategories(urlCategoryInfoID string) (*URLCategoryInformation, error) {
	var urlCategoryInfo URLCategoryInformation
	err := service.Client.Read(fmt.Sprintf("%s/%s", urlCategoriesEndpoint, urlCategoryInfoID), &urlCategoryInfo)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning custom url category from Get: %s", urlCategoryInfo.ID)
	return &urlCategoryInfo, nil
}

func (service *Service) GetCustomURLCategories(customName string) (*URLCategoryInformation, error) {
	var urlCategories []URLCategoryInformation
	err := service.Client.Read(fmt.Sprintf("%s?customOnly=%s", urlCategoriesEndpoint, url.QueryEscape(customName)), &urlCategories)
	if err != nil {
		return nil, err
	}
	for _, custom := range urlCategories {
		if strings.EqualFold(custom.ID, customName) {
			return &custom, nil
		}
	}
	return nil, fmt.Errorf("no custom url category found with name: %s", customName)
}
