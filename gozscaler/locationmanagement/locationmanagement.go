package locationmanagement

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	locationsEndpoint = "/locations"
)

// Gets locations only, not sub-locations. When a location matches the given search parameter criteria only its parent location is included in the result set, not its sub-locations.
type Locations struct {
	ID                                  int              `json:"id"`
	Name                                string           `json:"name,omitempty"`
	ParentID                            int              `json:"parentId,omitempty"`
	UpBandwidth                         int              `json:"upBandwidth,omitempty"`
	DnBandwidth                         int              `json:"dnBandwidth,omitempty"`
	Country                             string           `json:"country"`
	TZ                                  string           `json:"tz"`
	IPAddresses                         []string         `json:"ipAddresses"`
	Ports                               string           `json:"ports"`
	VPNCredentials                      []VPNCredentials `json:"vpnCredentials"`
	AuthRequired                        bool             `json:"authRequired"`
	SSLScanEnabled                      bool             `json:"sslScanEnabled"`
	ZappSSLScanEnabled                  bool             `json:"zappSSLScanEnabled"`
	XFFForwardEnabled                   bool             `json:"xffForwardEnabled"`
	SurrogateIP                         bool             `json:"surrogateIP"`
	IdleTimeInMinutes                   int              `json:"idleTimeInMinutes"`
	DisplayTimeUnit                     string           `json:"displayTimeUnit"`
	SurrogateIPEnforcedForKnownBrowsers bool             `json:"surrogateIPEnforcedForKnownBrowsers"`
	SurrogateRefreshTimeInMinutes       int              `json:"surrogateRefreshTimeInMinutes"`
	SurrogateRefreshTimeUnit            string           `json:"surrogateRefreshTimeUnit"`
	OFWEnabled                          bool             `json:"ofwEnabled"`
	IPSControl                          bool             `json:"ipsControl"`
	AUPEnabled                          bool             `json:"aupEnabled"`
	CautionEnabled                      bool             `json:"cautionEnabled"`
	AUPBlockInternetUntilAccepted       bool             `json:"aupBlockInternetUntilAccepted"`
	AUPForceSSLInspection               bool             `json:"aupForceSslInspection"`
	AUPTimeoutInDays                    int              `json:"aupTimeoutInDays"`
	Profile                             string           `json:"profile"`
	Description                         string           `json:"description"`
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

type VPNCredentials struct {
	ID           int         `json:"id"`
	Type         string      `json:"type,omitempty"`
	FQDN         string      `json:"fqdn"`
	PreSharedKey string      `json:"preSharedKey,omitempty"`
	Comments     string      `json:"comments,omitempty"`
	Location     []Location  `json:"location"`
	ManagedBy    []ManagedBy `json:"managedBy"`
}

// Gets locations only, not sub-locations. When a location matches the given search parameter criteria only its parent location is included in the result set, not its sub-locations
func (service *Service) GetLocations(locationID string) (*Locations, error) {
	var location Locations
	err := service.Client.Read(fmt.Sprintf("%s/%s", locationsEndpoint, locationID), &location)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Location from Get: %d", location.ID)
	return &location, nil
}

func (service *Service) GetLocationByName(locationName string) (*Locations, error) {
	var locations []Locations
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := service.Client.Read(fmt.Sprintf("%s?page=1&pageSize=1000", locationsEndpoint), &locations)
	if err != nil {
		return nil, err
	}
	for _, location := range locations {
		if strings.EqualFold(location.Name, locationName) {
			return &location, nil
		}
	}
	return nil, fmt.Errorf("no location found with name: %s", locationName)
}

func (service *Service) CreateLocations(locations *Locations) (*Locations, error) {
	resp, err := service.Client.Create(locationsEndpoint, *locations)
	if err != nil {
		return nil, err
	}

	createdLocations, ok := resp.(*Locations)
	if !ok {
		return nil, errors.New("object returned from api was not a location pointer")
	}

	log.Printf("returning locations from create: %d", createdLocations.ID)
	return createdLocations, nil
}

func (service *Service) UpdateLocations(locationsID string, locations *Locations) (*Locations, error) {
	resp, err := service.Client.Update(locationsEndpoint+"/"+locationsID, *locations)
	if err != nil {
		return nil, err
	}
	updatedLocations, _ := resp.(*Locations)

	log.Printf("returning locations from Update: %d", updatedLocations.ID)
	return updatedLocations, nil
}

func (service *Service) DeleteLocations(locationsID string) error {
	err := service.Client.Delete(locationsEndpoint + "/" + locationsID)
	if err != nil {
		return err
	}

	return nil
}

// Gets a name and ID dictionary of locations.
func (service *Service) GetSublocations(sublocations string) (*Locations, error) {
	var subLocations Locations
	err := service.Client.Read(locationsEndpoint, "/"+"%s"+"/sublocations")
	if err != nil {
		return nil, err
	}

	return &subLocations, nil
}
