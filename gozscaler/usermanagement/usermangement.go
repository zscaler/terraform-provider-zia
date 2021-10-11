package usermanagement

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	departmentEndpoint = "/departments"
	groupsEndpoint     = "/groups"
	usersEndpoint      = "/users"
)

type User struct {
	ID            int          `json:"id"`
	Name          string       `json:"name,omitempty"`
	Email         string       `json:"email,omitempty"`
	Groups        []Groups     `json:"groups"`
	Departments   *Departments `json:"department"`
	Comments      string       `json:"comments,omitempty"`
	TempAuthEmail string       `json:"tempAuthEmail,omitempty"`
	Password      string       `json:"password,omitempty"`
	AdminUser     bool         `json:"adminUser"`
	Type          string       `json:"type,omitempty"`
	Deleted       bool         `json:"deleted"`
}
type Departments struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idpId"`
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted"`
}

type Groups struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idpId"`
	Comments string `json:"comments,omitempty"`
}

func (service *Service) GetDepartments(departmentID int) (*Departments, error) {
	var departments Departments
	err := service.Client.Read(fmt.Sprintf("%s/%d", departmentEndpoint, departmentID), &departments)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning departments from Get: %d", departments.ID)
	return &departments, nil
}

func (service *Service) GetDepartmentsByName(departmentName string) (*Departments, error) {
	var departments []Departments
	err := service.Client.Read(departmentEndpoint, &departments)
	if err != nil {
		return nil, err
	}
	for _, department := range departments {
		if strings.EqualFold(department.Name, departmentName) {
			return &department, nil
		}
	}
	return nil, fmt.Errorf("no department found with name: %s", departmentName)
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

func (service *Service) Get(userID int) (*User, error) {
	var user User
	err := service.Client.Read(fmt.Sprintf("%s/%d", usersEndpoint, userID), &user)
	if err != nil {
		return nil, err
	}

	log.Printf("returning user from Get: %d", user.ID)
	return &user, nil
}

func (service *Service) GetUserByName(userName string) (*User, error) {
	var users []User
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

func (service *Service) Create(userID *User) (*User, error) {
	resp, err := service.Client.Create(usersEndpoint, *userID)
	if err != nil {
		return nil, err
	}

	createdUsers, ok := resp.(*User)
	if !ok {
		return nil, errors.New("object returned from api was not a user pointer")
	}

	log.Printf("returning user from create: %v", createdUsers.ID)
	return createdUsers, nil
}

func (service *Service) Update(userID int, users *User) (*User, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", usersEndpoint, userID), *users)
	if err != nil {
		return nil, nil, err
	}
	updatedUser, _ := resp.(*User)
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
