package usermanagement

import (
	"fmt"
	"net/http"
)

const (
	userMgmtEndpoint = "/api/v1/adminUsers"
)

type Department struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IdpID    string `json:"idpId,omitempty"`
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted"`
}

func (service *Service) Get(departmentID string) (*Department, *http.Response, error) {
	v := new(Department)
	relativeURL := fmt.Sprintf("%s/%s", userMgmtEndpoint, departmentID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
