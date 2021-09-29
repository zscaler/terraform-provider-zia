package ipdestinationgroups

import (
	"fmt"
	"log"
	"strings"
)

const (
	ipDestinationGroupsLiteEndpoint = "/ipDestinationGroups/lite"
)

type IPDestinationGroupsLite struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

func (service *Service) GetIPDestinationGroupsLite(ipDestinationGroupLiteID int) (*IPDestinationGroupsLite, error) {
	var ipDestinationGroupsLite IPDestinationGroupsLite
	err := service.Client.Read(fmt.Sprintf("%s/%d", ipDestinationGroupsLiteEndpoint, ipDestinationGroupLiteID), &ipDestinationGroupsLite)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning ip destination group lite from Get: %d", ipDestinationGroupsLite.ID)
	return &ipDestinationGroupsLite, nil
}

func (service *Service) GetIPDestinationGroupsLiteByName(ipDestinationLiteGroupsName string) (*IPDestinationGroupsLite, error) {
	var ipDestinationLiteGroups []IPDestinationGroupsLite
	err := service.Client.Read(ipDestinationGroupsLiteEndpoint, &ipDestinationLiteGroups)
	if err != nil {
		return nil, err
	}
	for _, ipDestinationLiteGroup := range ipDestinationLiteGroups {
		if strings.EqualFold(ipDestinationLiteGroup.Name, ipDestinationLiteGroupsName) {
			return &ipDestinationLiteGroup, nil
		}
	}
	return nil, fmt.Errorf("no ip destination group found with name: %s", ipDestinationLiteGroupsName)
}
