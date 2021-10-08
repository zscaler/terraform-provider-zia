package greinternalipranges

import (
	"fmt"
	"log"
)

const (
	greTunnelIPRangeEndpoint = "/greTunnels/availableInternalIpRanges"
)

type GREInternalIPRange struct {
	StartIPAddress string `json:"startIPAddress,omitempty"`
	EndIPAddress   string `json:"endIPAddress,omitempty"`
}

func (service *Service) GetGREInternalIPRange(count int) (*[]GREInternalIPRange, error) {
	var greInternalIPRanges []GREInternalIPRange
	err := service.Client.Read(fmt.Sprintf("%s?limit=%d", greTunnelIPRangeEndpoint, count), &greInternalIPRanges)
	if err != nil {
		return nil, err
	}
	if len(greInternalIPRanges) < count {
		return nil, fmt.Errorf("not enough internal IP range available, got %d internal IP range, required: %d", len(greInternalIPRanges), count)
	}
	log.Printf("Returning internal IP range: %s", greInternalIPRanges)
	return &greInternalIPRanges, nil
}
