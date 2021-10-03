package ipdestinationgroups

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	ipDestinationGroupsEndpoint = "/ipDestinationGroups"
)

type IPDestinationGroups struct {
	ID           int      `json:"id"`
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,omitempty"`
	Addresses    []string `json:"addresses,omitempty"`
	Description  string   `json:"description,omitempty"`
	IPCategories []string `json:"ipCategories,omitempty"`
	Countries    []string `json:"countries,omitempty"`
}

func (service *Service) Get(ipGroupID int) (*IPDestinationGroups, error) {
	var ipDestinationGroups IPDestinationGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", ipDestinationGroupsEndpoint, ipGroupID), &ipDestinationGroups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning ip destination group from Get: %d", ipDestinationGroups.ID)
	return &ipDestinationGroups, nil
}

func (service *Service) GetByName(ipDestinationGroupsName string) (*IPDestinationGroups, error) {
	var ipDestinationGroups []IPDestinationGroups
	err := service.Client.Read(ipDestinationGroupsEndpoint, &ipDestinationGroups)
	if err != nil {
		return nil, err
	}
	for _, ipDestinationGroup := range ipDestinationGroups {
		if strings.EqualFold(ipDestinationGroup.Name, ipDestinationGroupsName) {
			return &ipDestinationGroup, nil
		}
	}
	return nil, fmt.Errorf("no ip destination group found with name: %s", ipDestinationGroupsName)
}

// Adds a GRE tunnel configuration.
func (service *Service) Create(ipGroupID *IPDestinationGroups) (*IPDestinationGroups, error) {
	resp, err := service.Client.Create(ipDestinationGroupsEndpoint, *ipGroupID)
	if err != nil {
		return nil, err
	}

	createdIPDestinationGroups, ok := resp.(*IPDestinationGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an ip destination group pointer")
	}

	log.Printf("returning ip destination group from create: %d", createdIPDestinationGroups.ID)
	return createdIPDestinationGroups, nil
}

func (service *Service) Update(ipGroupID string, ipGroup *IPDestinationGroups) (*IPDestinationGroups, error) {
	resp, err := service.Client.Update(ipDestinationGroupsEndpoint+"/"+ipGroupID, *ipGroup)
	if err != nil {
		return nil, err
	}
	updatedIPDestinationGroups, _ := resp.(*IPDestinationGroups)

	log.Printf("returning ip destination group from update: %d", updatedIPDestinationGroups.ID)
	return updatedIPDestinationGroups, nil
}

func (service *Service) Delete(ipGroupID string) error {
	err := service.Client.Delete(ipDestinationGroupsEndpoint + "/" + ipGroupID)
	if err != nil {
		return err
	}

	return nil
}
