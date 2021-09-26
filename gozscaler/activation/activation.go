package activation

import (
	"errors"
	"net/http"
)

const (
	activationStatusEndpoint = "/status"
	activationEndpoint       = "/status/activate"
)

type Activation struct {
	Status string `json:"status,omitempty"`
}

func (service *Service) GetActivationStatus() (*Activation, error) {
	var activation Activation
	err := service.Client.Read(activationStatusEndpoint, &activation)
	if err != nil {
		return nil, err
	}

	return &activation, nil
}

func (service *Service) CreateActivation(activation *Activation) (*Activation, *http.Response, error) {
	resp, err := service.Client.Create(activationEndpoint, activation)
	if err != nil {
		return nil, nil, err
	}

	createdActivation, ok := resp.(*Activation)
	if !ok {
		return nil, nil, errors.New("Object returned from API was not an activation Pointer")
	}

	return createdActivation, nil, nil
}
