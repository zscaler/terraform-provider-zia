package usermanagement

import (
	"fmt"
	"log"
	"strings"
)

const (
	groupsEndpoint = "/groups"
)

type Groups struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idpId"`
	Comments string `json:"comments,omitempty"`
}

func (service *Service) GetGroups(groupID int) (*Groups, error) {
	var groups Groups
	err := service.Client.Read(fmt.Sprintf("%s/%d", groupsEndpoint, groupID), &groups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Groups from Get: %d", groups.ID)
	return &groups, nil
}

func (service *Service) GetGroupByName(groupName string) (*Groups, error) {
	var groups []Groups
	err := service.Client.Read(groupsEndpoint, &groups)
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		if strings.EqualFold(group.Name, groupName) {
			return &group, nil
		}
	}
	return nil, fmt.Errorf("no group found with name: %s", groupName)
}
