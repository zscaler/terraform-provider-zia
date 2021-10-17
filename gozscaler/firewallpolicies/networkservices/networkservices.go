package networkservices

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	networkServicesEndpoint = "/networkServices"
)

type NetworkServices struct {
	ID            int            `json:"id"`
	Name          string         `json:"name,omitempty"`
	Tag           string         `json:"tag,omitempty"`
	SrcTCPPorts   []SrcTCPPorts  `json:"srcTcpPorts,omitempty"`
	DestTCPPorts  []DestTCPPorts `json:"destTcpPorts,omitempty"`
	SrcUDPPorts   []SrcUDPPorts  `json:"srcUdpPorts,omitempty"`
	DestUDPPorts  []DestUDPPorts `json:"destUdpPorts,omitempty"`
	Type          string         `json:"type,omitempty"`
	Description   string         `json:"description,omitempty"`
	IsNameL10nTag bool           `json:"isNameL10nTag,omitempty"`
}
type SrcTCPPorts struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type DestTCPPorts struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}
type SrcUDPPorts struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type DestUDPPorts struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

func (service *Service) Get(serviceID int) (*NetworkServices, error) {
	var networkServices NetworkServices
	err := service.Client.Read(fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID), &networkServices)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning network services from Get: %d", networkServices.ID)
	return &networkServices, nil
}

func (service *Service) GetByName(networkServiceName string) (*NetworkServices, error) {
	var networkServices []NetworkServices
	err := service.Client.Read(networkServicesEndpoint, &networkServices)
	if err != nil {
		return nil, err
	}
	for _, networkService := range networkServices {
		if strings.EqualFold(networkService.Name, networkServiceName) {
			return &networkService, nil
		}
	}
	return nil, fmt.Errorf("no network services found with name: %s", networkServiceName)
}

func (service *Service) Create(networkService *NetworkServices) (*NetworkServices, error) {
	resp, err := service.Client.Create(networkServicesEndpoint, *networkService)
	if err != nil {
		return nil, err
	}

	createdNetworkServices, ok := resp.(*NetworkServices)
	if !ok {
		return nil, errors.New("object returned from api was not a network service pointer")
	}

	log.Printf("returning network service from create: %d", createdNetworkServices.ID)
	return createdNetworkServices, nil
}

func (service *Service) Update(serviceID int, networkService *NetworkServices) (*NetworkServices, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID), *networkService)
	if err != nil {
		return nil, nil, err
	}
	updatedNetworkServices, _ := resp.(*NetworkServices)

	log.Printf("returning network service from Update: %d", updatedNetworkServices.ID)
	return updatedNetworkServices, nil, nil
}

func (service *Service) Delete(serviceID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
