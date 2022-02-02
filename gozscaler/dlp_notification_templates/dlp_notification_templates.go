package dlp_notification_templates

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	dlpNotificationTemplatesEndpoint = "/dlpNotificationTemplates"
)

type DlpNotificationTemplates struct {
	ID               int    `json:"id"`
	Name             string `json:"name,omitempty"`
	Subject          string `json:"subject,omitempty"`
	AttachContent    bool   `json:"attachContent,omitempty"`
	PlainTextMessage string `json:"plainTextMessage,omitempty"`
	HtmlMessage      string `json:"htmlMessage,omitempty"`
}

func (service *Service) Get(dlpTemplateID int) (*DlpNotificationTemplates, error) {
	var dlpTemplates DlpNotificationTemplates
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpNotificationTemplatesEndpoint, dlpTemplateID), &dlpTemplates)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning dlp notification template from Get: %d", dlpTemplates.ID)
	return &dlpTemplates, nil
}

func (service *Service) GetByName(templateName string) (*DlpNotificationTemplates, error) {
	var dlpTemplates []DlpNotificationTemplates
	err := service.Client.Read(dlpNotificationTemplatesEndpoint, &dlpTemplates)
	if err != nil {
		return nil, err
	}
	for _, template := range dlpTemplates {
		if strings.EqualFold(template.Name, templateName) {
			return &template, nil
		}
	}
	return nil, fmt.Errorf("no dictionary found with name: %s", templateName)
}

func (service *Service) Create(dlpTemplateID *DlpNotificationTemplates) (*DlpNotificationTemplates, *http.Response, error) {
	resp, err := service.Client.Create(dlpNotificationTemplatesEndpoint, *dlpTemplateID)
	if err != nil {
		return nil, nil, err
	}

	createdDlpTemplate, ok := resp.(*DlpNotificationTemplates)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a dlp dictionary pointer")
	}

	log.Printf("returning new dlp notification template from create: %d", createdDlpTemplate.ID)
	return createdDlpTemplate, nil, nil
}

func (service *Service) Update(dlpTemplateID int, dlpTemplates *DlpNotificationTemplates) (*DlpNotificationTemplates, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", dlpNotificationTemplatesEndpoint, dlpTemplateID), *dlpTemplates)
	if err != nil {
		return nil, nil, err
	}
	updatedDlpTemplate, _ := resp.(*DlpNotificationTemplates)

	log.Printf("returning updates from dlp notification template from update: %d", updatedDlpTemplate.ID)
	return updatedDlpTemplate, nil, nil
}

func (service *Service) Delete(dlpTemplateID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", dlpNotificationTemplatesEndpoint, dlpTemplateID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
