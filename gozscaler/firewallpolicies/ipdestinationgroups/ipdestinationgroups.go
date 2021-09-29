package ipdestinationgroups

import (
	"fmt"
	"log"
	"strings"
)

const (
	ipDestinationGroupsEndpoint     = "/ipDestinationGroups"
	ipDestinationGroupsLiteEndpoint = "/ipDestinationGroups/lite"
)

type IPDestinationGroups struct {
	ID           int      `json:"id"`
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,omitempty"`
	Addresses    []string `json:"addresses,omitempty"`
	Description  string   `json:"description"`
	IPCategories []string `json:"ipCategories"`
	Countries    []string `json:"countries"`
}

type IPDestinationGroupsLite struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

func (service *Service) GetIPDestinationGroups(ipDestinationGroupID int) (*IPDestinationGroups, error) {
	var ipDestinationGroups IPDestinationGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", ipDestinationGroupsEndpoint, ipDestinationGroupID), &ipDestinationGroups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning ip destination group from Get: %d", ipDestinationGroups.ID)
	return &ipDestinationGroups, nil
}

func (service *Service) GetIPDestinationGroupsByName(ipDestinationGroupsName string) (*IPDestinationGroups, error) {
	var ipDestinationGroups []IPDestinationGroups
	err := service.Client.Read(ipDestinationGroupsEndpoint, &ipDestinationGroups)
	if err != nil {
		return nil, err
	}
	for _, ipDestinationGroup := range ipDestinationGroups {
		if strings.EqualFold(ipDestinationGroup.Name, ipDestinationGroupsName) {
			return &ipDestinationGroup, nil
		}
	}
	return nil, fmt.Errorf("no dictionary found with name: %s", ipDestinationGroupsName)
}

func (service *Service) GetIPDestinationGroupsLite(ipDestinationGroupLiteID int) (*IPDestinationGroups, error) {
	var ipDestinationGroupsLite IPDestinationGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", ipDestinationGroupsLiteEndpoint, ipDestinationGroupLiteID), &ipDestinationGroupsLite)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning ip destination group lite from Get: %d", ipDestinationGroupsLite.ID)
	return &ipDestinationGroupsLite, nil
}

func (service *Service) GetIPDestinationLiteGroupsByName(ipDestinationLiteGroupsName string) (*IPDestinationGroups, error) {
	var ipDestinationLiteGroups []IPDestinationGroups
	err := service.Client.Read(ipDestinationGroupsEndpoint, &ipDestinationLiteGroups)
	if err != nil {
		return nil, err
	}
	for _, ipDestinationLiteGroup := range ipDestinationLiteGroups {
		if strings.EqualFold(ipDestinationLiteGroup.Name, ipDestinationLiteGroupsName) {
			return &ipDestinationLiteGroup, nil
		}
	}
	return nil, fmt.Errorf("no dictionary found with name: %s", ipDestinationLiteGroupsName)
}
