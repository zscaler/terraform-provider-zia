package grevirtualiplist

const (
	vipRecommendedListEndpoint = "/vip/recommendedList"
)

type GREVirtualIPList struct {
	ID                 string `json:"id,omitempty"`
	VirtualIp          string `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool   `json:"privateServiceEdge,omitempty"`
	DataCenter         string `json:"dataCenter,omitempty"`
}

func (service *Service) GetZSGREVirtualIPList() (*GREVirtualIPList, error) {
	var greVirtualIps GREVirtualIPList
	err := service.Client.Read(vipRecommendedListEndpoint, &greVirtualIps)
	if err != nil {
		return nil, err
	}

	return &greVirtualIps, nil
}
