package networkapplications

/*
import (
	"fmt"
	"log"
)

const (
	networkApplicationsEndpoint = "/networkApplication/"
)

type NetworkApplications struct {
	ID             string `json:"id"`
	ParentCategory string `json:"parentCategory,omitempty"`
	Description    string `json:"description,omitempty"`
	Deprecated     bool   `json:"deprecated"`
}

func (service *Service) GetNetworkApplications() (*NetworkApplications, error) {
	var networkApplication ([]NetworkApplications)
	err := service.Client.Read(fmt.Sprintf("%s?locale=%s", networkApplicationsEndpoint), &networkApplication)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning departments from Get: %s", networkApplication)
	return &networkApplication
}
*/
