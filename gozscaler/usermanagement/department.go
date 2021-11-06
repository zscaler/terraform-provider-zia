package usermanagement

import (
	"fmt"
	"log"
	"strings"
)

const (
	departmentEndpoint = "/departments"
)

type Department struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idpId"`
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted"`
}

func (service *Service) GetDepartments(departmentID int) (*Department, error) {
	var departments Department
	err := service.Client.Read(fmt.Sprintf("%s/%d", departmentEndpoint, departmentID), &departments)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning departments from Get: %d", departments.ID)
	return &departments, nil
}

func (service *Service) GetDepartmentsByName(departmentName string) (*Department, error) {
	var departments []Department
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
