package networkservices

import (
	"fmt"
	"log"
	"strings"
)

const (
	networkServicesEndpoint = "/networkServices/"
)

type NetworkServices struct {
	ID            int            `json:"id"`
	Name          string         `json:"name,omitempty"`
	Tag           string         `json:"tag"`
	SrcTCPPorts   []SrcTCPPorts  `json:"srcTcpPorts"`
	DestTCPPorts  []DestTCPPorts `json:"destTcpPorts"`
	SrcUDPPorts   []SrcUDPPorts  `json:"srcUdpPorts"`
	DestUDPPorts  []DestUDPPorts `json:"destUdpPorts"`
	Type          string         `json:"type"`
	Description   string         `json:"description"`
	IsNameL10nTag string         `json:"isNameL10nTag"`
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

func (service *Service) GetNetworkServicesLite(serviceID int) (*NetworkServices, error) {
	var networkServices NetworkServices
	err := service.Client.Read(fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID), &networkServices)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning network services from Get: %d", networkServices.ID)
	return &networkServices, nil
}

func (service *Service) GetNetworkServicesLiteByName(networkServiceName string) (*NetworkServices, error) {
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
