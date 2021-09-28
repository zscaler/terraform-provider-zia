package greinternalipranges

import "log"

const (
	greTunnelIPRangeEndpoint = "/api/v1/greTunnels/availableInternalIpRanges"
)

type GREInternalIPRanges struct {
	StartIPAddress string `json:"startIPAddress,omitempty"`
	EndIPAddress   string `json:"endIPAddress,omitempty"`
}

func (service *Service) GetGREInternalIPRanges() (*GREInternalIPRanges, error) {
	var greInternalIPRanges GREInternalIPRanges
	err := service.Client.Read(greTunnelIPRangeEndpoint, &greInternalIPRanges)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning VPN Credentials from Get: %s", greInternalIPRanges)
	return &greInternalIPRanges, nil
}
