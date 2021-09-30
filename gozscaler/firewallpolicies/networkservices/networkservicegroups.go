package networkservices

import (
	"fmt"
	"log"
	"strings"
)

const (
	networkServiceGroupsEndpoint = "/networkServiceGroups/lite"
)

type NetworkServiceGroups struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
	//Services    []Services `json:"services"`
	Description string `json:"description"`
}

func (service *Service) GetNetworkServiceGroupsLite(NetworkServiceGroupsID int) (*NetworkServiceGroups, error) {
	var networkServiceGroups NetworkServiceGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", networkServiceGroupsEndpoint, NetworkServiceGroupsID), &networkServiceGroups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning network application group from Get: %d", networkServiceGroups.ID)
	return &networkServiceGroups, nil
}

func (service *Service) GetNetworkServiceGroupsLiteByName(networkserviceGroupsName string) (*NetworkServiceGroups, error) {
	var networkServiceGroups []NetworkServiceGroups
	err := service.Client.Read(networkServiceGroupsEndpoint, &networkServiceGroups)
	if err != nil {
		return nil, err
	}
	for _, networkServiceGroup := range networkServiceGroups {
		if strings.EqualFold(networkServiceGroup.Name, networkserviceGroupsName) {
			return &networkServiceGroup, nil
		}
	}
	return nil, fmt.Errorf("no network service groups found with name: %s", networkserviceGroupsName)
}
