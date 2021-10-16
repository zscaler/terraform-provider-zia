package urlcategories

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	urlCategoriesEndpoint = "/urlCategories"
)

type URLCategory struct {
	ID                               string            `json:"id,omitempty"`
	ConfiguredName                   string            `json:"configuredName,omitempty"`
	Keywords                         []string          `json:"keywords"`
	KeywordsRetainingParentCategory  []string          `json:"keywordsRetainingParentCategory"`
	Urls                             []string          `json:"urls"`
	DBCategorizedUrls                []string          `json:"dbCategorizedUrls"`
	CustomCategory                   bool              `json:"customCategory"`
	Scopes                           []Scopes          `json:"scopes"`
	Editable                         bool              `json:"editable"`
	Description                      string            `json:"description,omitempty"`
	Type                             string            `json:"type,omitempty"`
	URLKeywordCounts                 *URLKeywordCounts `json:"urlKeywordCounts,omitempty"`
	Val                              int               `json:"val,omitempty"`
	CustomUrlsCount                  int               `json:"customUrlsCount,omitempty"`
	SuperCategory                    string            `json:"superCategory,omitempty"`
	UrlsRetainingParentCategoryCount int               `json:"urlsRetainingParentCategoryCount"`
}
type Scopes struct {
	ScopeGroupMemberEntities []ScopeGroupMemberEntities `json:"scopeGroupMemberEntities"`
	Type                     string                     `json:"Type,omitempty"`
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
	TotalURLCount            int `json:"totalUrlCount,omitempty"`
	RetainParentURLCount     int `json:"retainParentUrlCount,omitempty"`
	TotalKeywordCount        int `json:"totalKeywordCount,omitempty"`
	RetainParentKeywordCount int `json:"retainParentKeywordCount,omitempty"`
}

func (service *Service) Get(categoryID string) (*URLCategory, error) {
	var urlCategory URLCategory
	err := service.Client.Read(fmt.Sprintf("%s/%s", urlCategoriesEndpoint, categoryID), &urlCategory)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning custom url category from Get: %s", urlCategory.ID)
	return &urlCategory, nil
}

func (service *Service) GetCustomURLCategories(customName string) (*URLCategory, error) {
	var urlCategory []URLCategory
	err := service.Client.Read(fmt.Sprintf("%s?customOnly=%s", urlCategoriesEndpoint, url.QueryEscape(customName)), &urlCategory)
	if err != nil {
		return nil, err
	}
	for _, custom := range urlCategory {
		if strings.EqualFold(custom.ID, customName) {
			return &custom, nil
		}
	}
	return nil, fmt.Errorf("no custom url category found with name: %s", customName)
}

func (service *Service) CreateURLCategories(category *URLCategory) (*URLCategory, error) {
	resp, err := service.Client.Create(urlCategoriesEndpoint, *category)
	if err != nil {
		return nil, err
	}

	createdUrlCategory, ok := resp.(*URLCategory)
	if !ok {
		return nil, errors.New("object returned from API was not a url category Pointer")
	}

	log.Printf("Returning url category from Create: %v", createdUrlCategory.ID)
	return createdUrlCategory, nil
}

func (service *Service) UpdateURLCategories(categoryID string, category *URLCategory) (*URLCategory, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%s", urlCategoriesEndpoint, categoryID), *category)
	if err != nil {
		return nil, nil, err
	}
	updatedUrlCategory, _ := resp.(*URLCategory)
	log.Printf("Returning url category from Update: %s", updatedUrlCategory.ID)
	return updatedUrlCategory, nil, nil
}

func (service *Service) DeleteURLCategories(categoryID string) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%s", urlCategoriesEndpoint, categoryID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
