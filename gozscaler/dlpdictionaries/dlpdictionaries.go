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
	ID                    int               `json:"id"`
	Name                  string            `json:"name,omitempty"`
	Description           string            `json:"description,omitempty"`
	ConfidenceThreshold   string            `json:"confidenceThreshold,omitempty"`
	CustomPhraseMatchType string            `json:"customPhraseMatchType,omitempty"`
	NameL10nTag           bool              `json:"nameL10nTag"`
	Custom                bool              `json:"custom"`
	ThresholdType         string            `json:"thresholdType,omitempty"`
	DictionaryType        string            `json:"dictionaryType,omitempty"`
	Proximity             int               `json:"proximity,omitempty"`
	Phrases               []Phrases         `json:"phrases"`
	Patterns              []Patterns        `json:"patterns"`
	EDMMatchDetails       []EDMMatchDetails `json:"exactDataMatchDetails"`
	// IDMProfileMatchAccuracy []IDMProfileMatchAccuracy `json:"idmProfileMatchAccuracy"`
}

type Phrases struct {
	Action string `json:"action,omitempty"`
	Phrase string `json:"phrase,omitempty"`
}

type Patterns struct {
	Action  string `json:"action,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

type EDMMatchDetails struct {
	DictionaryEdmMappingID int    `json:"dictionaryEdmMappingId,omitempty"`
	SchemaID               int    `json:"schemaId,omitempty"`
	PrimaryField           int    `json:"primaryField,omitempty"`
	SecondaryFields        []int  `json:"secondaryFields,omitempty"`
	SecondaryFieldMatchOn  string `json:"secondaryFieldMatchOn,omitempty"`
}

type IDMProfileMatchAccuracy struct {
	AdpIdmProfile string `json:"adpIdmProfile,omitempty"`
	MatchAccuracy string `json:"matchAccuracy,omitempty"`
}

func (service *Service) Get(dlpDictionariesID int) (*DlpDictionary, error) {
	var dlpDictionary DlpDictionary
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID), &dlpDictionary)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning dictionary from Get: %d", dlpDictionary.ID)
	return &dlpDictionary, nil
}

func (service *Service) GetByName(dictionaryName string) (*DlpDictionary, error) {
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

func (service *Service) Create(dlpDictionariesID *DlpDictionary) (*DlpDictionary, *http.Response, error) {
	resp, err := service.Client.Create(dlpDictionariesEndpoint, *dlpDictionariesID)
	if err != nil {
		return nil, nil, err
	}

	createdDlpDictionary, ok := resp.(*DlpDictionary)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a dlp dictionary pointer")
	}

	log.Printf("returning new custom dlp dictionary that uses patterns and phrases from create: %d", createdDlpDictionary.ID)
	return createdDlpDictionary, nil, nil
}

func (service *Service) Update(dlpDictionariesID int, dlpDictionaries *DlpDictionary) (*DlpDictionary, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID), *dlpDictionaries)
	if err != nil {
		return nil, nil, err
	}
	updatedDlpDictionary, _ := resp.(*DlpDictionary)

	log.Printf("returning updates custom dlp dictionary that uses patterns and phrases from ppdate: %d", updatedDlpDictionary.ID)
	return updatedDlpDictionary, nil, nil
}

func (service *Service) DeleteDlpDictionary(dlpDictionariesID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
