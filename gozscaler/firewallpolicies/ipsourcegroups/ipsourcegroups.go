package ipsourcegroups

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	ipSourceGroupsEndpoint = "/ipSourceGroups"
)

type IPSourceGroups struct {
	ID          int      `json:"id"`
	Name        string   `json:"name,omitempty"`
	IPAddresses []string `json:"ipAddresses,omitempty"`
	Description string   `json:"description,omitempty"`
}

func (service *Service) Get(ipGroupID int) (*IPSourceGroups, error) {
	var ipSourceGroups IPSourceGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", ipSourceGroupsEndpoint, ipGroupID), &ipSourceGroups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning ip source groupfrom Get: %d", ipSourceGroups.ID)
	return &ipSourceGroups, nil
}

func (service *Service) GetByName(ipSourceGroupsName string) (*IPSourceGroups, error) {
	var ipSourceGroups []IPSourceGroups
	err := service.Client.Read(ipSourceGroupsEndpoint, &ipSourceGroups)
	if err != nil {
		return nil, err
	}
	for _, ipSourceGroup := range ipSourceGroups {
		if strings.EqualFold(ipSourceGroup.Name, ipSourceGroupsName) {
			return &ipSourceGroup, nil
		}
	}
	return nil, fmt.Errorf("no ip source group found with name: %s", ipSourceGroupsName)
}

func (service *Service) Create(ipGroupID *IPSourceGroups) (*IPSourceGroups, error) {
	resp, err := service.Client.Create(ipSourceGroupsEndpoint, *ipGroupID)
	if err != nil {
		return nil, err
	}

	createdIPSourceGroups, ok := resp.(*IPSourceGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an ip source group pointer")
	}

	log.Printf("returning ip source group from create: %d", createdIPSourceGroups.ID)
	return createdIPSourceGroups, nil
}

func (service *Service) Update(ipGroupID int, ipGroup *IPSourceGroups) (*IPSourceGroups, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", ipSourceGroupsEndpoint, ipGroupID), *ipGroup)
	if err != nil {
		return nil, err
	}
	updatedIPSourceGroups, _ := resp.(*IPSourceGroups)

	log.Printf("returning ip source group from update: %d", updatedIPSourceGroups.ID)
	return updatedIPSourceGroups, nil
}

func (service *Service) Delete(ipGroupID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", ipSourceGroupsEndpoint, ipGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
