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
	CloudName     string   `json:"cloudName"`
	Region        string   `json:"region"`
	City          string   `json:"city"`
	DataCenter    string   `json:"dataCenter"`
	Location      string   `json:"location"`
	VPNIPs        []string `json:"vpnIps"`
	VPNDomainName string   `json:"vpnDomainName"`
	GREIPs        []string `json:"greIps"`
	GREDomainName string   `json:"greDomainName"`
	PACIPs        []string `json:"pacIps"`
	PACDomainName string   `json:"pacDomainName"`
}

type GREVirtualIPList struct {
	ID                 int    `json:"id"`
	VirtualIp          string `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool   `json:"privateServiceEdge,omitempty"`
	DataCenter         string `json:"dataCenter,omitempty"`
}

// Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud, including region and data center information. By default, the request gets all public VIPs in the cloud, but you can also include private or all VIPs in the request, if necessary.
func (service *Service) GetZscalerVIPs(datacenter string) (*ZscalerVIPs, error) {
	var zscalerVips []ZscalerVIPs

	err := service.Client.Read(vipsEndpoint, &zscalerVips)
	if err != nil {
		return nil, err
	}
	for _, vips := range zscalerVips {
		if strings.EqualFold(vips.DataCenter, datacenter) {
			return &vips, nil
		}
	}
	return nil, fmt.Errorf("no datacenter found with name: %s", datacenter)
}

// Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud by sourceIP
func (service *Service) GetZSGREVirtualIPList(sourceIP string) (*[]GREVirtualIPList, error) {
	var zscalerVips []GREVirtualIPList
	err := service.Client.Read(fmt.Sprintf("%s?sourceIp=%s", vipRecommendedListEndpoint, sourceIP), &zscalerVips)
	if err != nil {
		return nil, err
	}
	return &zscalerVips, nil
}
