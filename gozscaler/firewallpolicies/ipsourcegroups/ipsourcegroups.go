package ipsourcegroups

import (
	"fmt"
	"log"
	"strings"
)

const (
	ipSourceGroupsLiteEndpoint = "/ipSourceGroups/lite"
)

type IPSourceGroupsLite struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

func (service *Service) GetIPSourceGroupsLite(ipSourceGroupsLiteID int) (*IPSourceGroupsLite, error) {
	var ipSourceGroupsLite IPSourceGroupsLite
	err := service.Client.Read(fmt.Sprintf("%s/%d", ipSourceGroupsLiteEndpoint, ipSourceGroupsLiteID), &ipSourceGroupsLite)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning ip source group lite from Get: %d", ipSourceGroupsLiteID)
	return &ipSourceGroupsLite, nil
}

func (service *Service) GetIPSourceGroupsLiteByName(IPSourceGroupsLitesName string) (*IPSourceGroupsLite, error) {
	var ipSourceGroupsLite []IPSourceGroupsLite
	err := service.Client.Read(ipSourceGroupsLiteEndpoint, &ipSourceGroupsLite)
	if err != nil {
		return nil, err
	}
	for _, ipSourceGroupLite := range ipSourceGroupsLite {
		if strings.EqualFold(ipSourceGroupLite.Name, IPSourceGroupsLitesName) {
			return &ipSourceGroupLite, nil
		}
	}
	return nil, fmt.Errorf("no ip source group found with name: %s", IPSourceGroupsLitesName)
}
