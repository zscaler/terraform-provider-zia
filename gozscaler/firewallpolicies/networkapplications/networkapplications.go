package networkapplications

import (
	"fmt"
)

const (
	networkApplicationsEndpoint = "/networkApplications"
)

type NetworkApplications struct {
	ID             string `json:"id"`
	ParentCategory string `json:"parentCategory,omitempty"`
	Description    string `json:"description,omitempty"`
	Deprecated     bool   `json:"deprecated"`
}

func (service *Service) GetNetworkApplication(id, locale string) (*NetworkApplications, error) {
	var networkApplications NetworkApplications
	url := fmt.Sprintf("%s/%s", networkApplicationsEndpoint, id)
	if locale != "" {
		url = fmt.Sprintf("%s?locale=%s", url, locale)
	}
	err := service.Client.Read(url, &networkApplications)
	if err != nil {
		return nil, err
	}
	return &networkApplications, nil
}
