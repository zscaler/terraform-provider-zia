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

type Departments struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idpId"`
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted"`
}

type Group struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idpId"`
	Comments string `json:"comments,omitempty"`
}

type User struct {
	ID            int         `json:"id"`
	Name          string      `json:"name,omitempty"`
	Email         string      `json:"email,omitempty"`
	Groups        []Group     `json:"groups"`
	Departments   Departments `json:"department"`
	Comments      string      `json:"comments,omitempty"`
	TempAuthEmail string      `json:"tempAuthEmail,omitempty"`
	Password      string      `json:"password,omitempty"`
	AdminUser     bool        `json:"adminUser"`
	Type          string      `json:"type,omitempty"`
	Deleted       bool        `json:"deleted"`
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

func (service *Service) GetGroups(groupID int) (*Group, error) {
	var groups Group
	err := service.Client.Read(fmt.Sprintf("%s/%d", groupsEndpoint, groupID), &groups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Groups from Get: %d", groups.ID)
	return &groups, nil
}

func (service *Service) GetGroupByName(groupName string) (*Group, error) {
	var groups []Group
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

func (service *Service) GetUser(userID int) (*User, error) {
	var user User
	err := service.Client.Read(fmt.Sprintf("%s/%d", usersEndpoint, userID), &user)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Groups from Get: %d", user.ID)
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

func (service *Service) CreateUser(user *User) (*User, *http.Response, error) {
	resp, err := service.Client.Create(usersEndpoint, *user)
	if err != nil {
		return nil, nil, err
	}

	createdUsers, ok := resp.(*User)
	if !ok {
		return nil, nil, errors.New("object returned from API was not a User Pointer")
	}

	log.Printf("Returning User from Create: %v", createdUsers.ID)
	return createdUsers, nil, nil
}

func (service *Service) UpdateUser(userID int, user *User) (*User, error) {
	resp, err := service.Client.Update(fmt.Sprintf("%s/%d", usersEndpoint, userID), *user)
	if err != nil {
		return nil, err
	}
	updatedUser, _ := resp.(*User)
	log.Printf("Returning User from Update: %d", updatedUser.ID)
	return updatedUser, nil
}

func (service *Service) DeleteUser(userID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", usersEndpoint, userID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
