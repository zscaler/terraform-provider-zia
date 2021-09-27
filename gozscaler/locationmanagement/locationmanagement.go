package locationmanagement

import (
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/vpncredentials"
)

const (
	locationsEndpoint     = "/locations"
	locationsLiteEndpoint = "/locations/lite"
)

// Gets locations only, not sub-locations. When a location matches the given search parameter criteria only its parent location is included in the result set, not its sub-locations.
type Locations struct {
	ID                                  int                              `json:"id,omitempty"`
	Name                                string                           `json:"name,omitempty"`
	ParentID                            int                              `json:"parentId,omitempty"`
	UpBandwidth                         int                              `json:"upBandwidth,omitempty"`
	DnBandwidth                         int                              `json:"dnBandwidth,omitempty"`
	Country                             map[string]interface{}           `json:"country"`
	TZ                                  map[string]interface{}           `json:"tz"`
	IPAddresses                         map[string]interface{}           `json:"ipAddresses"`
	Ports                               map[int]interface{}              `json:"ports"`
	VPNCredentials                      []*vpncredentials.VPNCredentials `json:"vpnCredentials,omitempty"`
	AuthRequired                        bool                             `json:"authRequired"`
	SSLScanEnabled                      bool                             `json:"sslScanEnabled"`
	ZappSSLScanEnabled                  bool                             `json:"zappSSLScanEnabled"`
	XFFForwardEnabled                   bool                             `json:"xffForwardEnabled"`
	SurrogateIP                         bool                             `json:"surrogateIP"`
	IdleTimeInMinutes                   int                              `json:"idleTimeInMinutes"`
	DisplayTimeUnit                     int                              `json:"displayTimeUnit"`
	SurrogateIPEnforcedForKnownBrowsers bool                             `json:"surrogateIPEnforcedForKnownBrowsers"`
	SurrogateRefreshTimeInMinutes       int                              `json:"surrogateRefreshTimeInMinutes"`
	SurrogateRefreshTimeUnit            string                           `json:"surrogateRefreshTimeUnit"`
	OFWEnabled                          bool                             `json:"ofwEnabled"`
	IPSControl                          bool                             `json:"ipsControl"`
	AUPEnabled                          bool                             `json:"aupEnabled"`
	CautionEnabled                      bool                             `json:"cautionEnabled"`
	AUPBlockInternetUntilAccepted       bool                             `json:"aupBlockInternetUntilAccepted"`
	AUPForceSSLInspection               bool                             `json:"aupForceSslInspection"`
	AUPTimeoutInDays                    string                           `json:"aupTimeoutInDays"`
	ManagedBy                           []ManagedBy                      `json:"managedBy"`
	Profile                             string                           `json:"profile"`
	Description                         string                           `json:"description"`
}

type ManagedBy struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// Gets locations only, not sub-locations. When a location matches the given search parameter criteria only its parent location is included in the result set, not its sub-locations
func (service *Service) GetLocations(locationID string) (*Locations, error) {
	var location Locations
	err := service.Client.Read(locationsEndpoint+"/"+locationID, &location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

// Gets a name and ID dictionary of locations.
func (service *Service) GetLocationLite(locationLite string) (*Locations, error) {
	var lite Locations
	err := service.Client.Read(locationsLiteEndpoint, &locationLite)
	if err != nil {
		return nil, err
	}

	return &lite, nil
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
