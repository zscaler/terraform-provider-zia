package gretunnels

import (
	"errors"
	"log"
	"net/http"
)

const (
	greTunnelsEndpoint = "/greTunnels"
)

type GreTunnels struct {
	ID                   string             `json:"id,omitempty"`
	SourceIP             string             `json:"sourceIp,omitempty"`
	PrimaryDestVip       []PrimaryDestVip   `json:"primaryDestVip,omitempty"`
	SecondaryDestVip     []SecondaryDestVip `json:"secondaryDestVip,omitempty"`
	InternalIpRange      string             `json:"internalIpRange,omitempty"`
	ManagedBy            []ManagedBy        `json:"managedBy,omitempty"`
	LastModifiedBy       []LastModifiedBy   `json:"lastModifiedBy,omitempty"`
	LastModificationTime string             `json:"lastModificationTime,omitempty"`
	WithinCountry        bool               `json:"withinCountry"`
	Comment              string             `json:"comment,omitempty"`
	IpUnnumbered         bool               `json:"ipUnnumbered"`
}

type PrimaryDestVip struct {
	ID                 string `json:"id,omitempty"`
	VirtualIP          string `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool   `json:"privateServiceEdge"`
	Datacenter         string `json:"datacenter,omitempty"`
}

type SecondaryDestVip struct {
	ID                 string `json:"id,omitempty"`
	VirtualIP          string `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool   `json:"privateServiceEdge"`
	Datacenter         string `json:"datacenter,omitempty"`
}

type ManagedBy struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type LastModifiedBy struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func (service *Service) GetGreTunnels(greTunnelID string) (*GreTunnels, error) {
	var greTunnels GreTunnels
	err := service.Client.Read(greTunnelsEndpoint+"/"+greTunnelID, &greTunnels)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning GRE Tunnels from Get: %s", greTunnels.ID)
	return &greTunnels, nil
}

func (service *Service) CreateGreTunnels(greTunnelID *GreTunnels) (*GreTunnels, *http.Response, error) {
	resp, err := service.Client.Create(greTunnelsEndpoint, *greTunnelID)
	if err != nil {
		return nil, nil, err
	}

	createdGreTunnels, ok := resp.(*GreTunnels)
	if !ok {
		return nil, nil, errors.New("Object returned from API was not a GRE Tunnel Pointer")
	}

	log.Printf("Returning GRE Tunnel from Create: %s", createdGreTunnels.ID)
	return createdGreTunnels, nil, nil
}

func (service *Service) UpdateGreTunnels(greTunnelID string, greTunnels *GreTunnels) (*GreTunnels, *http.Response, error) {
	resp, err := service.Client.Update(greTunnelsEndpoint+"/"+greTunnelID, *greTunnels)
	if err != nil {
		return nil, nil, err
	}
	updatedGreTunnels, _ := resp.(*GreTunnels)

	log.Printf("Returning GRE Tunnels from Update: %s", updatedGreTunnels.ID)
	return updatedGreTunnels, nil, nil
}

func (service *Service) DeleteGreTunnels(greTunnelID string) (*http.Response, error) {
	err := service.Client.Delete(greTunnelsEndpoint + "/" + greTunnelID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
