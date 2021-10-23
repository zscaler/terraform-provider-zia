package networkservices

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	networkServiceGroupsEndpoint = "/networkServiceGroups"
)

type NetworkServiceGroups struct {
	ID          int        `json:"id"`
	Name        string     `json:"name,omitempty"`
	Services    []Services `json:"services,omitempty"`
	Description string     `json:"description,omitempty"`
}
type Services struct {
	ID            int            `json:"id"`
	Name          string         `json:"name,omitempty"`
	Tag           string         `json:"tag,omitempty"`
	SrcTCPPorts   []NetworkPorts `json:"srcTcpPorts,omitempty"`
	DestTCPPorts  []NetworkPorts `json:"destTcpPorts,omitempty"`
	SrcUDPPorts   []NetworkPorts `json:"srcUdpPorts,omitempty"`
	DestUDPPorts  []NetworkPorts `json:"destUdpPorts,omitempty"`
	Type          string         `json:"type,omitempty"`
	Description   string         `json:"description,omitempty"`
	IsNameL10nTag bool           `json:"isNameL10nTag,omitempty"`
}

func (service *Service) GetNetworkServiceGroups(serviceGroupID int) (*NetworkServiceGroups, error) {
	var networkServiceGroups NetworkServiceGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", networkServiceGroupsEndpoint, serviceGroupID), &networkServiceGroups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning network service groups from Get: %d", networkServiceGroups.ID)
	return &networkServiceGroups, nil
}

func (service *Service) GetNetworkServiceGroupsByName(serviceGroupsName string) (*NetworkServiceGroups, error) {
	var networkServiceGroups []NetworkServiceGroups
	err := service.Client.Read(networkServiceGroupsEndpoint, &networkServiceGroups)
	if err != nil {
		return nil, err
	}
	for _, networkServiceGroup := range networkServiceGroups {
		if strings.EqualFold(networkServiceGroup.Name, serviceGroupsName) {
			return &networkServiceGroup, nil
		}
	}
	return nil, fmt.Errorf("no network service groups found with name: %s", serviceGroupsName)
}

func (service *Service) CreateNetworkServiceGroups(networkServiceGroups *NetworkServiceGroups) (*NetworkServiceGroups, error) {
	resp, err := service.Client.Create(networkServiceGroupsEndpoint, *networkServiceGroups)
	if err != nil {
		return nil, err
	}

	createdNetworkServiceGroups, ok := resp.(*NetworkServiceGroups)
	if !ok {
		return nil, errors.New("object returned from api was not a network service groups pointer")
	}

	log.Printf("returning network service groups from create: %d", createdNetworkServiceGroups.ID)
	return createdNetworkServiceGroups, nil
}

func (service *Service) UpdateNetworkServiceGroups(serviceGroupID int, networkServiceGroups *NetworkServiceGroups) (*NetworkServiceGroups, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", networkServiceGroupsEndpoint, serviceGroupID), *networkServiceGroups)
	if err != nil {
		return nil, nil, err
	}
	updatedNetworkServiceGroups, _ := resp.(*NetworkServiceGroups)

	log.Printf("returning network service groups from Update: %d", updatedNetworkServiceGroups.ID)
	return updatedNetworkServiceGroups, nil, nil
}

func (service *Service) DeleteNetworkServiceGroups(serviceGroupID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", networkServiceGroupsEndpoint, serviceGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
