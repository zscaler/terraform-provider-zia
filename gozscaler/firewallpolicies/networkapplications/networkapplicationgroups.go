package networkapplications

import (
	"fmt"
	"log"
	"strings"
)

const (
	networkAppGroupsEndpoint = "/networkApplicationGroups"
)

type NetworkApplicationGroups struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name,omitempty"`
	NetworkApplications []string `json:"networkApplications,omitempty"`
	Description         string   `json:"description,omitempty"`
}

func (service *Service) GetNetworkApplicationGroups(groupID int) (*NetworkApplicationGroups, error) {
	var networkApplicationGroups NetworkApplicationGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", networkAppGroupsEndpoint, groupID), &networkApplicationGroups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning network application group from Get: %d", networkApplicationGroups.ID)
	return &networkApplicationGroups, nil
}

func (service *Service) GetNetworkApplicationGroupsByName(appGroupsName string) (*NetworkApplicationGroups, error) {
	var networkApplicationGroups []NetworkApplicationGroups
	err := service.Client.Read(networkAppGroupsEndpoint, &networkApplicationGroups)
	if err != nil {
		return nil, err
	}
	for _, networkAppGroup := range networkApplicationGroups {
		if strings.EqualFold(networkAppGroup.Name, appGroupsName) {
			return &networkAppGroup, nil
		}
	}
	return nil, fmt.Errorf("no network application group found with name: %s", appGroupsName)
}
