package vpncredentials

import (
	"errors"
	"log"
	"net/http"
)

const (
	vpnCredentialsEndpoint = "/api/v1/vpnCredentials"
)

type VPNCredentials struct {
	ID           int         `json:"id"`
	Type         string      `json:"type,omitempty"`
	FQDN         []string    `json:"fqdn"`
	PreSharedKey string      `json:"preSharedKey,omitempty"`
	Comments     string      `json:"comments,omitempty"`
	Location     []Location  `json:"location"`
	ManagedBy    []ManagedBy `json:"managedBy"`
}

type Location struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

type ManagedBy struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

func (service *Service) GetVPNCredentials(vpnCredentialID string) (*VPNCredentials, error) {
	var vpnCredentials VPNCredentials
	err := service.Client.Read(vpnCredentialsEndpoint+"/"+vpnCredentialID, &vpnCredentials)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning VPN Credentials from Get: %s", vpnCredentials.ID)
	return &vpnCredentials, nil
}

func (service *Service) CreateVPNCredentials(vpnCredentialID *VPNCredentials) (*VPNCredentials, *http.Response, error) {
	resp, err := service.Client.Create(vpnCredentialsEndpoint, *vpnCredentialID)
	if err != nil {
		return nil, nil, err
	}

	createdVpnCredentials, ok := resp.(*VPNCredentials)
	if !ok {
		return nil, nil, errors.New("Object returned from API was not a VPN Credential Pointer")
	}

	log.Printf("Returning VPN Credential from Create: %s", createdVpnCredentials.ID)
	return createdVpnCredentials, nil, nil
}

func (service *Service) UpdateVPNCredentials(vpnCredentialID string, vpnCredentials *VPNCredentials) (*VPNCredentials, *http.Response, error) {
	resp, err := service.Client.Update(vpnCredentialsEndpoint+"/"+vpnCredentialID, *vpnCredentials)
	if err != nil {
		return nil, nil, err
	}
	updatedVpnCredentials, _ := resp.(*VPNCredentials)

	log.Printf("Returning VPN Credential from Update: %s", updatedVpnCredentials.ID)
	return updatedVpnCredentials, nil, nil
}

func (service *Service) DeleteVPNCredentials(vpnCredentialID string) (*http.Response, error) {
	err := service.Client.Delete(vpnCredentialsEndpoint + "/" + vpnCredentialID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
