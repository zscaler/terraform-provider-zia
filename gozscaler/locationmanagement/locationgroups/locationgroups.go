package locationgroups

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	locationGroupEndpoint = "/locations/groups"
)

type LocationGroup struct {
	ID                           int                           `json:"id,omitempty"`
	Name                         string                        `json:"name,omitempty"`
	Deleted                      bool                          `json:"deleted,omitempty"`
	GroupType                    string                        `json:"groupType,omitempty"`
	DynamicLocationGroupCriteria *DynamicLocationGroupCriteria `json:"dynamicLocationGroupCriteria,omitempty"`
	Comments                     string                        `json:"comments"`
	Locations                    []Locations                   `json:"locations"`
	LastModUser                  *LastModUser                  `json:"lastModUser"`
	LastModTime                  int                           `json:"lastModTime"`
	Predefined                   bool                          `json:"predefined"`
}

type DynamicLocationGroupCriteria struct {
	Name                   *Name       `json:"name,omitempty"`
	Countries              []string    `json:"countries,omitempty"`
	City                   *City       `json:"city,omitempty"`
	ManagedBy              []ManagedBy `json:"managedBy,omitempty"`
	EnforceAuthentication  bool        `json:"enforceAuthentication"`
	EnforceAup             bool        `json:"enforceAup"`
	EnforceFirewallControl bool        `json:"enforceFirewallControl"`
	EnableXffForwarding    bool        `json:"enableXffForwarding"`
	EnableCaution          bool        `json:"enableCaution"`
	EnableBandwidthControl bool        `json:"enableBandwidthControl"`
	Profiles               []string    `json:"profiles"`
}

type Name struct {
	MatchString string `json:"matchString,omitempty"`
	MatchType   string `json:"matchType,omitempty"`
}

type City struct {
	MatchString string `json:"matchString,omitempty"`
	MatchType   string `json:"matchType,omitempty"`
}

type ManagedBy struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type Locations struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type LastModUser struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func (service *Service) GetLocationGroup(groupID int) (*LocationGroup, error) {
	var locationGroup LocationGroup
	err := service.Client.Read(fmt.Sprintf("%s/%d", locationGroupEndpoint, groupID), &locationGroup)
	if err != nil {
		return nil, err
	}

	log.Printf("returning location group from Get: %d", locationGroup.ID)
	return &locationGroup, nil
}

func (service *Service) GetLocationGroupByName(locationGroupName string) (*LocationGroup, error) {
	var locationGroups []LocationGroup
	err := service.Client.Read(fmt.Sprintf("%s?name=%s", locationGroupEndpoint, url.QueryEscape(locationGroupName)), &locationGroups)
	if err != nil {
		return nil, err
	}
	for _, locationGroup := range locationGroups {
		if strings.EqualFold(locationGroup.Name, locationGroupName) {
			return &locationGroup, nil
		}
	}
	return nil, fmt.Errorf("no location group found with name: %s", locationGroupName)
}

func (service *Service) CreateLocationGroup(locationGroups *LocationGroup) (*LocationGroup, error) {
	resp, err := service.Client.Create(locationGroupEndpoint, *locationGroups)
	if err != nil {
		return nil, err
	}

	createdLocationGroup, ok := resp.(*LocationGroup)
	if !ok {
		return nil, errors.New("object returned from api was not a location group pointer")
	}

	log.Printf("returning location group from create: %d", createdLocationGroup.ID)
	return createdLocationGroup, nil
}

func (service *Service) UpdateLocationGroup(groupID int, locationGroups *LocationGroup) (*LocationGroup, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", locationGroupEndpoint, groupID), *locationGroups)
	if err != nil {
		return nil, nil, err
	}
	updatedLocationGroup, _ := resp.(*LocationGroup)

	log.Printf("returning location group from update: %d", updatedLocationGroup.ID)
	return updatedLocationGroup, nil, nil
}

func (service *Service) DeleteLocationGroup(groupID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", locationGroupEndpoint, groupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
