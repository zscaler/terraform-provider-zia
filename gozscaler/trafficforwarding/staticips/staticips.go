package staticips

import (
	"errors"
	"log"
	"net/http"
)

const (
	staticIPEndpoint         = "/staticIP"
	staticIPValidateEndpoint = "/staticIP/validate"
)

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

func (service *Service) GetStaticIP(staticIpID string) (*StaticIP, error) {
	var staticIP StaticIP
	err := service.Client.Read(staticIPEndpoint+"/"+staticIpID, &staticIP)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Static IP from Get: %s", staticIP.ID)
	return &staticIP, nil
}

func (service *Service) CreateStaticIP(staticIpID *StaticIP) (*StaticIP, *http.Response, error) {
	resp, err := service.Client.Create(staticIPEndpoint, *staticIpID)
	if err != nil {
		return nil, nil, err
	}

	createdStaticIP, ok := resp.(*StaticIP)
	if !ok {
		return nil, nil, errors.New("Object returned from API was not a Static IP Pointer")
	}

	log.Printf("Returning Static IP from Create: %s", createdStaticIP.ID)
	return createdStaticIP, nil, nil
}

// Not sure if I want this in the code. All it does it return a "SUCCESS" message
func (service *Service) CreateStaticIPValidate(staticIpID *StaticIP) (*StaticIP, *http.Response, error) {
	resp, err := service.Client.Create(staticIPValidateEndpoint, *staticIpID)
	if err != nil {
		return nil, nil, err
	}

	createdStaticIPValidate, ok := resp.(*StaticIP)
	if !ok {
		return nil, nil, errors.New("Object returned from API was not a Static IP Pointer")
	}

	log.Printf("Returning Static IP validate from Create: %s", createdStaticIPValidate.ID)
	return createdStaticIPValidate, nil, nil
}

func (service *Service) UpdateStaticIP(staticIpID string, staticIP *StaticIP) (*StaticIP, *http.Response, error) {
	resp, err := service.Client.Update(staticIPEndpoint+"/"+staticIpID, *staticIP)
	if err != nil {
		return nil, nil, err
	}
	updatedStaticIP, _ := resp.(*StaticIP)

	log.Printf("Returning Static IP from Update: %s", updatedStaticIP.ID)
	return updatedStaticIP, nil, nil
}

func (service *Service) DeleteStaticIP(staticIpID string) (*http.Response, error) {
	err := service.Client.Delete(staticIPEndpoint + "/" + staticIpID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
