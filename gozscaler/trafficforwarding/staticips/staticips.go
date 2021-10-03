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
	ID                   int             `json:"id,omitempty"`
	IpAddress            string          `json:"ipAddress"`
	GeoOverride          bool            `json:"geoOverride"`
	Latitude             float64         `json:"latitude,omitempty"`
	Longitude            float64         `json:"longitude,omitempty"`
	RoutableIP           bool            `json:"routableIP,omitempty"`
	LastModificationTime int             `json:"lastModificationTime"`
	Comment              string          `json:"comment,omitempty"`
	ManagedBy            *ManagedBy      `json:"managedBy,omitempty"`      // Should probably move this to a common package. Used by multiple resources
	LastModifiedBy       *LastModifiedBy `json:"lastModifiedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources
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

func (service *Service) Get(staticIpID int) (*StaticIP, error) {
	var staticIP StaticIP
	err := service.Client.Read(fmt.Sprintf("%s/%d", staticIPEndpoint, staticIpID), &staticIP)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning static ip from Get: %d", staticIP.ID)
	return &staticIP, nil
}

func (service *Service) Create(staticIpID *StaticIP) (*StaticIP, *http.Response, error) {
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

func (service *Service) Update(staticIpID int, staticIP *StaticIP) (*StaticIP, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", staticIPEndpoint, staticIpID), *staticIP)
	if err != nil {
		return nil, nil, err
	}
	updatedStaticIP, _ := resp.(*StaticIP)

	log.Printf("returning static ip from update: %d", updatedStaticIP.ID)
	return updatedStaticIP, nil, nil
}
func (service *Service) Delete(staticIpID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", staticIPEndpoint, staticIpID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
