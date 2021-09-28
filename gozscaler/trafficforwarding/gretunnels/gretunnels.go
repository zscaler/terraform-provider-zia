package gretunnels

import (
	"errors"
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
	PrimaryDestVip       PrimaryDestVip   `json:"primaryDestVip,omitempty"`
	SecondaryDestVip     SecondaryDestVip `json:"secondaryDestVip,omitempty"`
	InternalIpRange      string           `json:"internalIpRange,omitempty"`
	ManagedBy            ManagedBy        `json:"managedBy,omitempty"`      // Should probably move this to a common package. Used by multiple resources
	LastModifiedBy       LastModifiedBy   `json:"lastModifiedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources
	LastModificationTime string           `json:"lastModificationTime,omitempty"`
	WithinCountry        bool             `json:"withinCountry"`
	Comment              string           `json:"comment,omitempty"`
	IpUnnumbered         bool             `json:"ipUnnumbered"`
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
	Extensions map[string]interface{} `json:"extensions"`
}

type LastModifiedBy struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// Gets a list of IP addresses with GRE tunnel details.
type IPGreTunnelInfo struct {
	IpAddress         string `json:"ipAddress,omitempty"`
	GreEnabled        bool   `json:"greEnabled,omitempty"`
	GreTunnelIP       string `json:"greTunnelIP,omitempty"`
	PrimaryGW         string `json:"primaryGW,omitempty"`
	SecondaryGW       string `json:"secondaryGW,omitempty"`
	TunID             string `json:"tunID,omitempty"`
	GreRangePrimary   string `json:"greRangePrimary,omitempty"`
	GreRangeSecondary string `json:"greRangeSecondary,omitempty"`
}

// Gets the next available GRE tunnel internal IP address ranges.
type GRETunnelIPRange struct {
	StartIPAddress string `json:"startIPAddress,omitempty"`
	EndIPAddress   bool   `json:"endIPAddress,omitempty"`
}

// Gets all provisioned static IP addresses.
type StaticIP struct {
	ID                   string           `json:"id,omitempty"`
	IpAddress            string           `json:"ipAddress,omitempty"`
	GeoOverride          bool             `json:"geoOverride"`
	Latitude             float64          `json:"latitude,omitempty"`
	Longitude            float64          `json:"longitude,omitempty"`
	RoutableIP           bool             `json:"routableIP,omitempty"`
	LastModificationTime string           `json:"lastModificationTime,omitempty"`
	ManagedBy            []ManagedBy      `json:"managedBy,omitempty"`      // Should probably move this to a common package. Used by multiple resources
	LastModifiedBy       []LastModifiedBy `json:"lastModifiedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources
	Comment              string           `json:"comment,omitempty"`
}

// Gets all provisioned GRE tunnel information.
func (service *Service) GetGreTunnels(greTunnelID string) (*GreTunnels, error) {
	var greTunnels GreTunnels
	err := service.Client.Read(greTunnelsEndpoint+"/"+greTunnelID, &greTunnels)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning GRE Tunnels from Get: %s", greTunnels.ID)
	return &greTunnels, nil
}

// Gets a list of IP addresses with GRE tunnel details.
func (service *Service) GetIPGreTunnelInfo() ([]IPGreTunnelInfo, error) {
	var ipGreTunnelInfo []IPGreTunnelInfo
	err := service.Client.Read(ipGreTunnelInfoEndpoint, &ipGreTunnelInfo)
	if err != nil {
		return nil, err
	}

	return ipGreTunnelInfo, nil
}

func (service *Service) GetGRETunnelIPRange() ([]GRETunnelIPRange, error) {
	var greTunnelIPRange []GRETunnelIPRange
	err := service.Client.Read(greTunnelIPRangeEndpoint, &greTunnelIPRange)
	if err != nil {
		return nil, err
	}

	return greTunnelIPRange, nil
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
