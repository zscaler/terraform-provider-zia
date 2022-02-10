package dlp_engines

import (
	"fmt"
	"log"
	"strings"
)

const (
	dlpEnginesEndpoint = "/dlpEngines"
)

type DlpEngines struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name,omitempty"`
	PredefinedEngineName string `json:"predefinedEngineName,omitempty"`
	EngineExpression     string `json:"engineExpression,omitempty"`
	CustomDlpEngine      bool   `json:"customDlpEngine,omitempty"`
	Description          string `json:"description,omitempty"`
}

func (service *Service) Get(dlpEngineID int) (*DlpEngines, error) {
	var dlpEngines DlpEngines
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpEnginesEndpoint, dlpEngineID), &dlpEngines)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning dlp engine from Get: %d", dlpEngines.ID)
	return &dlpEngines, nil
}

func (service *Service) GetByName(engineName string) (*DlpEngines, error) {
	var dlpEngines []DlpEngines
	err := service.Client.Read(dlpEnginesEndpoint, &dlpEngines)
	if err != nil {
		return nil, err
	}
	for _, engine := range dlpEngines {
		if strings.EqualFold(engine.Name, engineName) {
			return &engine, nil
		}
	}
	return nil, fmt.Errorf("no dlp engine found with name: %s", engineName)
}
