package dlpdictionaries

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	dlpDictionariesEndpoint     = "/dlpDictionaries"
	dlpDictionariesLiteEndpoint = "/dlpDictionaries/lite"
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

func (service *Service) GetDlpDictionaries(dlpDictionariesID int) (*DlpDictionary, error) {
	var dlpDictionary DlpDictionary
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID), &dlpDictionary)
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

func (service *Service) GetDlpDictionaryLite() (*DlpDictionary, error) {
	var dlpDictionary DlpDictionary
	err := service.Client.Read(dlpDictionariesLiteEndpoint, &dlpDictionary)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning name and ID dictionary of all custom and predefined DLP dictionaries from Get: %s", dlpDictionary.ID)
	return &dlpDictionary, nil
}

func (service *Service) CreateDlpDictionary(dlpDictionariesID *DlpDictionary) (*DlpDictionary, *http.Response, error) {
	resp, err := service.Client.Create(dlpDictionariesEndpoint, *dlpDictionariesID)
	if err != nil {
		return nil, nil, err
	}

	createdDlpDictionary, ok := resp.(*DlpDictionary)
	if !ok {
		return nil, nil, errors.New("Object returned from API was not a Dlp Dictionary Pointer")
	}

	log.Printf("Returning new custom DLP dictionary that uses Patterns and Phrases from Create: %s", createdDlpDictionary.ID)
	return createdDlpDictionary, nil, nil
}

func (service *Service) UpdateDlpDictionary(dlpDictionariesID string, dlpDictionaries *DlpDictionary) (*DlpDictionary, *http.Response, error) {
	resp, err := service.Client.Update(dlpDictionariesEndpoint+"/"+dlpDictionariesID, *dlpDictionaries)
	if err != nil {
		return nil, nil, err
	}
	updatedDlpDictionary, _ := resp.(*DlpDictionary)

	log.Printf("Returning updates custom DLP dictionary that uses Patterns and Phrases from Update: %s", updatedDlpDictionary.ID)
	return updatedDlpDictionary, nil, nil
}

func (service *Service) DeleteDlpDictionary(dlpDictionariesID string) (*http.Response, error) {
	err := service.Client.Delete(dlpDictionariesEndpoint + "/" + dlpDictionariesID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
