package vpncredentials

import (
	"errors"
	"fmt"
	"log"
	"net/http"
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

func (service *Service) Get(vpnCredentialID int) (*VPNCredentials, error) {
	var vpnCredentials VPNCredentials
	err := service.Client.Read(fmt.Sprintf("%s/%d", vpnCredentialsEndpoint, vpnCredentialID), &vpnCredentials)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning VPN Credentials from Get: %d", vpnCredentials.ID)
	return &vpnCredentials, nil
}

func (service *Service) GetByFQDN(vpnCredentialName string) (*VPNCredentials, error) {
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

func (service *Service) Create(vpnCredentials *VPNCredentials) (*VPNCredentials, error) {
	resp, err := service.Client.Create(vpnCredentialsEndpoint, *vpnCredentials)
	if err != nil {
		return nil, err
	}

	createdVpnCredentials, ok := resp.(*VPNCredentials)
	if !ok {
		return nil, errors.New("object returned from api was not a vpn credential pointer")
	}

	log.Printf("returning vpn credential from create: %d", createdVpnCredentials.ID)
	return createdVpnCredentials, nil
}

func (service *Service) Update(vpnCredentialID string, vpnCredentials *VPNCredentials) (*VPNCredentials, *http.Response, error) {
	resp, err := service.Client.Update(vpnCredentialsEndpoint+"/"+vpnCredentialID, *vpnCredentials)
	if err != nil {
		return nil, nil, err
	}
	updatedVpnCredentials, _ := resp.(*VPNCredentials)

	log.Printf("returning vpn credential from Update: %d", updatedVpnCredentials.ID)
	return updatedVpnCredentials, nil, nil
}

func (service *Service) Delete(vpnCredentialID string) error {
	err := service.Client.Delete(vpnCredentialsEndpoint + "/" + vpnCredentialID)
	if err != nil {
		return err
	}

	return nil
}
