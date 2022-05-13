package gretunnels

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	greTunnelsEndpoint = "/greTunnels"
	// ipGreTunnelInfoEndpoint = "/orgProvisioning/ipGreTunnelInfo"
)

type GreTunnels struct {
	ID                   int               `json:"id,omitempty"`
	SourceIP             string            `json:"sourceIp,omitempty"`
	InternalIpRange      string            `json:"internalIpRange,omitempty"`
	LastModificationTime int               `json:"lastModificationTime,omitempty"`
	WithinCountry        bool              `json:"withinCountry"`
	Comment              string            `json:"comment,omitempty"`
	IPUnnumbered         bool              `json:"ipUnnumbered"`
	ManagedBy            *ManagedBy        `json:"managedBy,omitempty"`      // Should probably move this to a common package. Used by multiple resources
	LastModifiedBy       *LastModifiedBy   `json:"lastModifiedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources
	PrimaryDestVip       *PrimaryDestVip   `json:"primaryDestVip,omitempty"`
	SecondaryDestVip     *SecondaryDestVip `json:"secondaryDestVip,omitempty"`
}

type PrimaryDestVip struct {
	ID                 int     `json:"id,omitempty"`
	VirtualIP          string  `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool    `json:"privateServiceEdge"`
	Datacenter         string  `json:"datacenter,omitempty"`
	Latitude           float64 `json:"latitude,omitempty"`
	Longitude          float64 `json:"longitude,omitempty"`
	City               string  `json:"city,omitempty"`
	CountryCode        string  `json:"countryCode,omitempty"`
	Region             string  `json:"region,omitempty"`
}

type SecondaryDestVip struct {
	ID                 int     `json:"id,omitempty"`
	VirtualIP          string  `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool    `json:"privateServiceEdge"`
	Datacenter         string  `json:"datacenter,omitempty"`
	Latitude           float64 `json:"latitude,omitempty"`
	Longitude          float64 `json:"longitude,omitempty"`
	City               string  `json:"city,omitempty"`
	CountryCode        string  `json:"countryCode,omitempty"`
	Region             string  `json:"region,omitempty"`
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

	log.Printf("returning gre tunnel from get: %d", greTunnels.ID)
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
		return nil, nil, errors.New("object returned from api was not a gre tunnel pointer")
	}

	log.Printf("returning gre tunnels from create: %d", createdGreTunnels.ID)
	return createdGreTunnels, nil, nil
}

func (service *Service) UpdateGreTunnels(greTunnelID int, greTunnels *GreTunnels) (*GreTunnels, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", greTunnelsEndpoint, greTunnelID), *greTunnels)
	if err != nil {
		return nil, nil, err
	}
	updatedGreTunnels, _ := resp.(*GreTunnels)

	log.Printf("returning gre tunnels from update: %d", updatedGreTunnels.ID)
	return updatedGreTunnels, nil, nil
}

func (service *Service) DeleteGreTunnels(greTunnelID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", greTunnelsEndpoint, greTunnelID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
