package networkapplications

import (
	"errors"
	"fmt"
	"log"
	"net/http"
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

	log.Printf("Returning network application groups from Get: %d", networkApplicationGroups.ID)
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
	return nil, fmt.Errorf("no network application groups found with name: %s", appGroupsName)
}

func (service *Service) Create(applicationGroup *NetworkApplicationGroups) (*NetworkApplicationGroups, error) {
	resp, err := service.Client.Create(networkAppGroupsEndpoint, *applicationGroup)
	if err != nil {
		return nil, err
	}

	createdApplicationGroups, ok := resp.(*NetworkApplicationGroups)
	if !ok {
		return nil, errors.New("object returned from api was not a network application groups pointer")
	}

	log.Printf("returning network application groups from create: %d", createdApplicationGroups.ID)
	return createdApplicationGroups, nil
}

func (service *Service) Update(groupID int, applicationGroup *NetworkApplicationGroups) (*NetworkApplicationGroups, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", networkAppGroupsEndpoint, groupID), *applicationGroup)
	if err != nil {
		return nil, nil, err
	}
	updatedApplicationGroups, _ := resp.(*NetworkApplicationGroups)

	log.Printf("returning network application groups from Update: %d", updatedApplicationGroups.ID)
	return updatedApplicationGroups, nil, nil
}

func (service *Service) Delete(groupID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", networkAppGroupsEndpoint, groupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
