package networkservicegroup

import (
	"fmt"
	"log"
	"strings"
)

const (
	networkServiceGroupsEndpoint = "/networkServiceGroups/"
)

type NetworkServiceGroups struct {
	ID          int        `json:"id"`
	Name        string     `json:"name,omitempty"`
	Services    []Services `json:"services,omitempty"`
	Description string     `json:"description,omitempty"`
}
type Services struct {
	ID            int    `json:"id"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	IsNameL10nTag bool   `json:"isNameL10nTag"`
}

func (service *Service) GetNetworkServiceGroups(NetworkServiceGroupsID int) (*NetworkServiceGroups, error) {
	var networkServiceGroups NetworkServiceGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", networkServiceGroupsEndpoint, NetworkServiceGroupsID), &networkServiceGroups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning network service groups from Get: %d", networkServiceGroups.ID)
	return &networkServiceGroups, nil
}

func (service *Service) GetNetworkServiceGroupsByName(networkserviceGroupsName string) (*NetworkServiceGroups, error) {
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
