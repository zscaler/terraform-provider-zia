package staticips

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	staticIPEndpoint = "/staticIP"
)

// Gets all provisioned static IP addresses.
type StaticIP struct {
	ID                   int            `json:"id,omitempty"`
	IpAddress            []string       `json:"ipAddress"`
	GeoOverride          bool           `json:"geoOverride"`
	Latitude             int            `json:"latitude,omitempty"`
	Longitude            int            `json:"longitude,omitempty"`
	RoutableIP           bool           `json:"routableIP,omitempty"`
	LastModificationTime int            `json:"lastModificationTime"`
	Comment              string         `json:"comment,omitempty"`
	ManagedBy            ManagedBy      `json:"managedBy,omitempty"`      // Should probably move this to a common package. Used by multiple resources
	LastModifiedBy       LastModifiedBy `json:"lastModifiedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources
}

type ManagedBy struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type LastModifiedBy struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func (service *Service) GetStaticIP(staticIpID string) (*StaticIP, error) {
	var staticIP StaticIP
	err := service.Client.Read(fmt.Sprintf("%s/%s", staticIPEndpoint, staticIpID), &staticIP)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning static ip from Get: %d", staticIP.ID)
	return &staticIP, nil
}

/*
func (service *Service) GetStaticByIP(staticIP string) (*StaticIP, error) {
	var staticips []StaticIP
	err := service.Client.Read(staticIPEndpoint, &staticips)
	if err != nil {
		return nil, err
	}
	for _, static := range staticips {
		if strings.EqualFold(static.IpAddress, staticIP) {
			return &static, nil
		}
	}
	return nil, fmt.Errorf("no static ip found with name: %s", staticIP)
}
*/

func (service *Service) CreateStaticIP(staticIpID *StaticIP) (*StaticIP, *http.Response, error) {
	resp, err := service.Client.Create(staticIPEndpoint, *staticIpID)
	if err != nil {
		return nil, nil, err
	}

	createdStaticIP, ok := resp.(*StaticIP)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a static ip pointer")
	}

	log.Printf("returning static ip from create: %d", createdStaticIP.ID)
	return createdStaticIP, nil, nil
}

func (service *Service) UpdateStaticIP(staticIpID string, staticIP *StaticIP) (*StaticIP, *http.Response, error) {
	resp, err := service.Client.Update(staticIPEndpoint+"/"+staticIpID, *staticIP)
	if err != nil {
		return nil, nil, err
	}
	updatedStaticIP, _ := resp.(*StaticIP)

	log.Printf("returning static ip from update: %d", updatedStaticIP.ID)
	return updatedStaticIP, nil, nil
}

func (service *Service) DeleteStaticIP(staticIpID string) (*http.Response, error) {
	err := service.Client.Delete(staticIPEndpoint + "/" + staticIpID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
