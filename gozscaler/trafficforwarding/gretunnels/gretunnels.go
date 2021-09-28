package gretunnels

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	greTunnelsEndpoint       = "/api/v1/greTunnels"
	ipGreTunnelInfoEndpoint  = "/api/v1/orgProvisioning/ipGreTunnelInfo"
	greTunnelIPRangeEndpoint = "/api/v1/greTunnels/availableInternalIpRanges"
)

type GreTunnels struct {
	ID                   int              `json:"id,omitempty"`
	SourceIP             string           `json:"sourceIp,omitempty"`
	InternalIpRange      string           `json:"internalIpRange,omitempty"`
	LastModificationTime int              `json:"lastModificationTime,omitempty"`
	WithinCountry        bool             `json:"withinCountry"`
	Comment              string           `json:"comment,omitempty"`
	IpUnnumbered         bool             `json:"ipUnnumbered"`
	ManagedBy            ManagedBy        `json:"managedBy,omitempty"`      // Should probably move this to a common package. Used by multiple resources
	LastModifiedBy       LastModifiedBy   `json:"lastModifiedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources
	PrimaryDestVip       PrimaryDestVip   `json:"primaryDestVip,omitempty"`
	SecondaryDestVip     SecondaryDestVip `json:"secondaryDestVip,omitempty"`
}

type PrimaryDestVip struct {
	ID                 int    `json:"id,omitempty"`
	VirtualIP          string `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool   `json:"privateServiceEdge"`
	Datacenter         string `json:"datacenter,omitempty"`
}

type SecondaryDestVip struct {
	ID                 int    `json:"id,omitempty"`
	VirtualIP          string `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool   `json:"privateServiceEdge"`
	Datacenter         string `json:"datacenter,omitempty"`
}

type ManagedBy struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

type LastModifiedBy struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// Gets all provisioned GRE tunnel information.

func (service *Service) GetGreTunnels(greTunnelID int) (*GreTunnels, error) {
	var greTunnels GreTunnels
	err := service.Client.Read(fmt.Sprintf("%s/%d", greTunnelsEndpoint, greTunnelID), &greTunnels)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning VPN Credentials from Get: %d", greTunnels.ID)
	return &greTunnels, nil
}

// Adds a GRE tunnel configuration.
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
