package dlp_web_rules

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/willguibr/terraform-provider-zia/gozscaler/common"
)

const (
	webDlpRulesEndpoint = "/webDlpRules"
)

type WebDLPRules struct {
	ID                       int                       `json:"id"`
	Order                    int                       `json:"order,omitempty"`
	Protocols                []string                  `json:"protocols,omitempty"`
	Rank                     int                       `json:"rank,omitempty"`
	Name                     string                    `json:"name,omitempty"`
	Description              string                    `json:"description,omitempty"`
	FileTypes                []string                  `json:"fileTypes,omitempty"`
	CloudApplications        []string                  `json:"cloudApplications,omitempty"`
	MinSize                  int                       `json:"minSize,omitempty"`
	Action                   string                    `json:"action,omitempty"`
	State                    string                    `json:"state,omitempty"`
	MatchOnly                bool                      `json:"matchOnly,omitempty"`
	LastModifiedTime         int                       `json:"lastModifiedTime,omitempty"`
	WithoutContentInspection bool                      `json:"withoutContentInspection,omitempty"`
	OcrEnabled               bool                      `json:"ocrEnabled,omitempty"`
	ZscalerIncidentReciever  bool                      `json:"zscalerIncidentReciever,omitempty"`
	ExternalAuditorEmail     string                    `json:"externalAuditorEmail,omitempty"`
	Auditor                  *common.IDNameExtensions  `json:"auditor,omitempty"`
	LastModifiedBy           *common.IDNameExtensions  `json:"lastModifiedBy,omitempty"`
	NotificationTemplate     *common.IDNameExtensions  `json:"notificationTemplate,omitempty"`
	IcapServer               *common.IDNameExtensions  `json:"icapServer,omitempty"`
	Locations                []common.IDNameExtensions `json:"locations,omitempty"`
	LocationGroups           []common.IDNameExtensions `json:"locationGroups,omitempty"`
	Groups                   []common.IDNameExtensions `json:"groups,omitempty"`
	Departments              []common.IDNameExtensions `json:"departments,omitempty"`
	Users                    []common.IDNameExtensions `json:"users,omitempty"`
	URLCategories            []common.IDNameExtensions `json:"urlCategories,omitempty"`
	DLPEngines               []common.IDNameExtensions `json:"dlpEngines,omitempty"`
	TimeWindows              []common.IDNameExtensions `json:"timeWindows,omitempty"`
	Labels                   []common.IDNameExtensions `json:"labels,omitempty"`
}

func (service *Service) Get(ruleID int) (*WebDLPRules, error) {
	var webDlpRules WebDLPRules
	err := service.Client.Read(fmt.Sprintf("%s/%d", webDlpRulesEndpoint, ruleID), &webDlpRules)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning web dlp rule from Get: %d", webDlpRules.ID)
	return &webDlpRules, nil
}

func (service *Service) GetByName(ruleName string) (*WebDLPRules, error) {
	var webDlpRules []WebDLPRules
	err := service.Client.Read(webDlpRulesEndpoint, &webDlpRules)
	if err != nil {
		return nil, err
	}
	for _, rule := range webDlpRules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no web dlp rule found with name: %s", ruleName)
}

func (service *Service) Create(ruleID *WebDLPRules) (*WebDLPRules, *http.Response, error) {
	resp, err := service.Client.Create(webDlpRulesEndpoint, *ruleID)
	if err != nil {
		return nil, nil, err
	}

	createdWebDlpRules, ok := resp.(*WebDLPRules)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a web dlp rule pointer")
	}

	log.Printf("returning new web dlp rule from create: %d", createdWebDlpRules.ID)
	return createdWebDlpRules, nil, nil
}

func (service *Service) Update(ruleID int, webDlpRules *WebDLPRules) (*WebDLPRules, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", webDlpRulesEndpoint, ruleID), *webDlpRules)
	if err != nil {
		return nil, nil, err
	}
	updatedWebDlpRules, _ := resp.(*WebDLPRules)

	log.Printf("returning updates from web dlp rule from update: %d", updatedWebDlpRules.ID)
	return updatedWebDlpRules, nil, nil
}

func (service *Service) Delete(ruleID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", webDlpRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
