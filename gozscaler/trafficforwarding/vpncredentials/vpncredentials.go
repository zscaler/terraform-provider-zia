package vpncredentials

import (
	"fmt"
	"log"
	"strings"
)

const (
	vpnCredentialsEndpoint = "/vpnCredentials"
)

type VPNCredentials struct {
	ID           int       `json:"id"`
	Type         string    `json:"type,omitempty"`
	FQDN         string    `json:"fqdn"`
	PreSharedKey string    `json:"preSharedKey,omitempty"`
	Comments     string    `json:"comments,omitempty"`
	Location     Location  `json:"location"`
	ManagedBy    ManagedBy `json:"managedBy"`
}
type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ManagedBy struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (service *Service) GetVPNCredentials(vpnCredentialID int) (*VPNCredentials, error) {
	var vpnCredentials VPNCredentials
	err := service.Client.Read(fmt.Sprintf("%s/%d", vpnCredentialsEndpoint, vpnCredentialID), &vpnCredentials)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning VPN Credentials from Get: %d", vpnCredentials.ID)
	return &vpnCredentials, nil
}

func (service *Service) GetVPNCredentialsByFQDN(vpnCredentialName string) (*VPNCredentials, error) {
	var vpnCredentials []VPNCredentials

	err := service.Client.Read(vpnCredentialsEndpoint, &vpnCredentials)
	if err != nil {
		return nil, err
	}
	for _, vpnCredential := range vpnCredentials {
		if strings.EqualFold(vpnCredential.FQDN, vpnCredentialName) {
			return &vpnCredential, nil
		}
	}
	return nil, fmt.Errorf("no vpn credentials found with fqdn: %s", vpnCredentialName)
}

/*
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
*/
