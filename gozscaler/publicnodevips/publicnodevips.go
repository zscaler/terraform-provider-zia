package publicnodevips

const (
	publicVipEndpoint          = "/vip"
	vipRecommendedListEndpoint = "/vip/recommendedList"
)

// Zscaler Public node details.
type PublicNodes struct {
	CloudName     string                 `json:"cloudName,omitempty"`
	Region        string                 `json:"region,omitempty"`
	City          string                 `json:"city,omitempty"`
	DataCenter    string                 `json:"dataCenter,omitempty"`
	Location      string                 `json:"location,omitempty"`
	VpnDomainName string                 `json:"vpnDomainName,omitempty"`
	GreDomainName string                 `json:"greDomainName,omitempty"`
	VpnIps        map[string]interface{} `json:"vpnIps"`
	GreIps        map[string]interface{} `json:"greIps"`
	PacIps        map[string]interface{} `json:"pacIps"`
	PacDomainName map[string]interface{} `json:"pacDomainName"`
}

func (service *Service) GetPublicNodeVipAddresses() (*PublicNodes, error) {
	var publicNodes PublicNodes
	err := service.Client.Read(publicVipEndpoint, &publicNodes)
	if err != nil {
		return nil, err
	}

	return &publicNodes, nil
}
