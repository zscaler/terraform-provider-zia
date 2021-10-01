package activation

import (
	"errors"
)

const (
	activationStatusEndpoint = "/status"
	activationEndpoint       = "/status/activate"
)

type Activation struct {
	Status string `json:"status"`
}

func (service *Service) GetActivationStatus() (*Activation, error) {
	var activation Activation
	err := service.Client.Read(activationStatusEndpoint, &activation)
	if err != nil {
		return nil, err
	}

	return &activation, nil
}

func (service *Service) CreateActivation(activation Activation) (*Activation, error) {
	resp, err := service.Client.Create(activationEndpoint, activation)
	if err != nil {
		return nil, err
	}

	createdActivation, ok := resp.(*Activation)
	if !ok {
		return nil, errors.New("object returned from api was not an activation pointer")
	}

	return createdActivation, nil
}
