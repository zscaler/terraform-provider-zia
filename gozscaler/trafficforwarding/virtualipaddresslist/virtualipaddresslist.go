package virtualipaddresslist

import (
	"fmt"
	"strings"
)

const (
	vipsEndpoint               = "/vips"
	vipRecommendedListEndpoint = "/vips/recommendedList"
)

type ZscalerVIPs struct {
	Cloudname     string   `json:"cloudName"`
	Region        string   `json:"region"`
	City          string   `json:"city"`
	Datacenter    string   `json:"dataCenter"`
	Location      string   `json:"location"`
	VPNIPs        []string `json:"vpnIps"`
	VPNDomainName string   `json:"vpnDomainName"`
	GREIPs        []string `json:"greIps"`
	GREDomainName string   `json:"greDomainName"`
	PACIPs        []string `json:"pacIps"`
	PACDomainName string   `json:"pacDomainName"`
}

type GREVirtualIPList struct {
	ID                 string `json:"id,omitempty"`
	VirtualIp          string `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool   `json:"privateServiceEdge,omitempty"`
	DataCenter         string `json:"dataCenter,omitempty"`
}

// Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud, including region and data center information. By default, the request gets all public VIPs in the cloud, but you can also include private or all VIPs in the request, if necessary.
func (service *Service) ZscalerVIPs(cloudName string) (*ZscalerVIPs, error) {
	var zscalerVips []ZscalerVIPs
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := service.Client.Read(fmt.Sprintf("%s?page=1&pageSize=1000", vipsEndpoint), &zscalerVips)
	if err != nil {
		return nil, err
	}
	for _, vips := range zscalerVips {
		if strings.EqualFold(vips.Cloudname, cloudName) {
			return &vips, nil
		}
	}
	return nil, fmt.Errorf("no cloud found with name: %s", cloudName)
}
