package usermanagement

import (
	"errors"
	"log"
	"net/http"
)

const (
	departmentEndpoint = "/api/v1/departments"
	groupsEndpoint     = "/api/v1/groups"
	usersEndpoint      = "/api/v1/users"
)

type Department struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IdpID    string `json:"idpId,omitempty"`
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted"`
}

type Groups struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IdpID    string `json:"idpId,omitempty"`
	Comments string `json:"comments,omitempty"`
}

type Users struct {
	ID            string       `json:"id,omitempty"`
	Name          string       `json:"name,omitempty"`
	Email         string       `json:"email,omitempty"`
	Groups        []Groups     `json:"groups"`
	Department    []Department `json:"department"`
	Comments      string       `json:"comments,omitempty"`
	TempAuthEmail string       `json:"tempAuthEmail,omitempty"`
	Password      string       `json:"password,omitempty"`
	AdminUser     string       `json:"adminUser,omitempty"`
	Type          string       `json:"type,omitempty"`
}

func (service *Service) GetDepartment(departmentID string) (*Department, error) {
	var department Department
	err := service.Client.Read(departmentEndpoint+"/"+departmentID, &department)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Department from Get: %s", department.ID)
	return &department, nil
}

func (service *Service) GetGroups(groupID string) (*Groups, error) {
	var groups Groups
	err := service.Client.Read(groupsEndpoint+"/"+groupID, &groups)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Groups from Get: %s", groups.ID)
	return &groups, nil
}

func (service *Service) GetUsers(groupID string) (*Users, *http.Response, error) {
	var users Users
	err := service.Client.Read(usersEndpoint+"/"+groupID, &users)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Returning Groups from Get: %s", users.ID)
	return &users, nil, nil
}

func (service *Service) CreateUsers(userID *Users) (*Users, *http.Response, error) {
	resp, err := service.Client.Create(usersEndpoint, *userID)
	if err != nil {
		return nil, nil, err
	}

	createdUsers, ok := resp.(*Users)
	if !ok {
		return nil, nil, errors.New("Object returned from API was not a User Pointer")
	}

	log.Printf("Returning User from Create: %s", createdUsers.ID)
	return createdUsers, nil, nil
}

func (service *Service) UpdateUsers(userID string, users *Users) (*Users, error) {
	resp, err := service.Client.Update(usersEndpoint+"/"+userID, *users)
	if err != nil {
		return nil, err
	}
	updatedUsers, _ := resp.(*Users)

	log.Printf("Returning User from Update: %s", updatedUsers.ID)
	return updatedUsers, nil
}

func (service *Service) DeleteUsers(userID string) (*http.Response, error) {
	err := service.Client.Delete(usersEndpoint + "/" + userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
