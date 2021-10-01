package dlpdictionaries

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	dlpDictionariesEndpoint = "/dlpDictionaries"
)

type DlpDictionary struct {
	ID                    int        `json:"id"`
	Name                  string     `json:"name,omitempty"`
	Description           string     `json:"description,omitempty"`
	ConfidenceThreshold   string     `json:"confidenceThreshold,omitempty"`
	Phrases               []Phrases  `json:"phrases"`
	CustomPhraseMatchType string     `json:"customPhraseMatchType"`
	Patterns              []Patterns `json:"patterns"`
	NameL10nTag           bool       `json:"nameL10nTag"`
	Custom                bool       `json:"custom"`
	ThresholdType         string     `json:"thresholdType"`
	DictionaryType        string     `json:"dictionaryType"`
}

type Phrases struct {
	Action string `json:"action,omitempty"`
	Phrase string `json:"phrase,omitempty"`
}

type Patterns struct {
	Action  string `json:"action,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

func (service *Service) GetDlpDictionaries(dlpDictionariesID string) (*DlpDictionary, error) {
	var dlpDictionary DlpDictionary
	err := service.Client.Read(fmt.Sprintf("%s/%s", dlpDictionariesEndpoint, dlpDictionariesID), &dlpDictionary)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning dictionary from Get: %d", dlpDictionary.ID)
	return &dlpDictionary, nil
}

func (service *Service) GetDlpDictionaryByName(dictionaryName string) (*DlpDictionary, error) {
	var dictionaries []DlpDictionary
	err := service.Client.Read(dlpDictionariesEndpoint, &dictionaries)
	if err != nil {
		return nil, err
	}
	for _, dictionary := range dictionaries {
		if strings.EqualFold(dictionary.Name, dictionaryName) {
			return &dictionary, nil
		}
	}
	return nil, fmt.Errorf("no dictionary found with name: %s", dictionaryName)
}

func (service *Service) CreateDlpDictionary(dlpDictionary *DlpDictionary) (*DlpDictionary, error) {
	resp, err := service.Client.Create(dlpDictionariesEndpoint, *dlpDictionary)
	if err != nil {
		return nil, err
	}

	createdDlpDictionary, ok := resp.(*DlpDictionary)
	if !ok {
		return nil, errors.New("object returned from api was not a dlp dictionary pointer")
	}

	log.Printf("Returning new custom DLP dictionary that uses Patterns and Phrases from Create: %d", createdDlpDictionary.ID)
	return createdDlpDictionary, nil
}

func (service *Service) UpdateDlpDictionary(dlpDictionariesID string, dlpDictionaries *DlpDictionary) (*DlpDictionary, error) {
	resp, err := service.Client.Update(fmt.Sprintf("%s/%s", dlpDictionariesEndpoint, dlpDictionariesID), *dlpDictionaries)
	if err != nil {
		return nil, err
	}
	updatedDlpDictionary, _ := resp.(*DlpDictionary)

	log.Printf("Returning updates custom DLP dictionary that uses Patterns and Phrases from Update: %d", updatedDlpDictionary.ID)
	return updatedDlpDictionary, nil
}

func (service *Service) DeleteDlpDictionary(dlpDictionariesID string) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%s", dlpDictionariesEndpoint, dlpDictionariesID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
