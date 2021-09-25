package dlpdictionaries

import (
	"fmt"
	"net/http"
)

const (
	dlpDictionariesEndpoint = "/dlpDictionaries"
)

type DlpDictionary struct {
	ID                    string     `json:"id,omitempty"`
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

func (service *Service) Get(dlpDictId string) (*DlpDictionary, *http.Response, error) {
	v := new(DlpDictionary)
	relativeURL := fmt.Sprintf("%s/%s", dlpDictionariesEndpoint, dlpDictId)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Create(dlpDict DlpDictionary) (*DlpDictionary, *http.Response, error) {
	v := new(DlpDictionary)
	resp, err := service.Client.NewRequestDo("POST", dlpDictionariesEndpoint, nil, dlpDict, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(dlpDictId string, dlp DlpDictionary) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", dlpDictionariesEndpoint, dlpDictId)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, dlpDictId, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(dlpDictId string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", dlpDictionariesEndpoint, dlpDictId)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
