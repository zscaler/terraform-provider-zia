package usermanagement

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/willguibr/terraform-provider-zia/gozscaler/common"
)

const (
	usersEndpoint = "/users"
)

type Users struct {
	ID            int                    `json:"id"`
	Name          string                 `json:"name,omitempty"`
	Email         string                 `json:"email,omitempty"`
	Groups        []common.UserGroups    `json:"groups,omitempty"`
	Department    *common.UserDepartment `json:"department,omitempty"`
	Comments      string                 `json:"comments,omitempty"`
	TempAuthEmail string                 `json:"tempAuthEmail,omitempty"`
	Password      string                 `json:"password,omitempty"`
	AdminUser     bool                   `json:"adminUser"`
	Type          string                 `json:"type,omitempty"`
	Deleted       bool                   `json:"deleted"`
}

func (service *Service) Get(userID int) (*Users, error) {
	var user Users
	err := service.Client.Read(fmt.Sprintf("%s/%d", usersEndpoint, userID), &user)
	if err != nil {
		return nil, err
	}

	log.Printf("returning user from Get: %d", user.ID)
	return &user, nil
}

func (service *Service) GetUserByName(userName string) (*Users, error) {
	var users []Users
	err := service.Client.Read(fmt.Sprintf("%s?name=%s", usersEndpoint, url.QueryEscape(userName)), &users)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if strings.EqualFold(user.Name, userName) {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("no user found with name: %s", userName)
}

func (service *Service) Create(userID *Users) (*Users, error) {
	resp, err := service.Client.Create(usersEndpoint, *userID)
	if err != nil {
		return nil, err
	}

	createdUsers, ok := resp.(*Users)
	if !ok {
		return nil, errors.New("object returned from api was not a user pointer")
	}

	log.Printf("returning user from create: %v", createdUsers.ID)
	return createdUsers, nil
}

func (service *Service) Update(userID int, users *Users) (*Users, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", usersEndpoint, userID), *users)
	if err != nil {
		return nil, nil, err
	}
	updatedUser, _ := resp.(*Users)
	log.Printf("returning user from update: %d", updatedUser.ID)
	return updatedUser, nil, nil
}

func (service *Service) Delete(userID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", usersEndpoint, userID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
