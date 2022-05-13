package rule_labels

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zscaler/terraform-provider-zia/gozscaler/common"
)

const (
	ruleLabelsEndpoint = "/ruleLabels"
)

type RuleLabels struct {
	ID                  int                      `json:"id"`
	Name                string                   `json:"name,omitempty"`
	Description         string                   `json:"description,omitempty"`
	LastModifiedTime    int                      `json:"lastModifiedTime,omitempty"`
	LastModifiedBy      *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`
	CreatedBy           *common.IDNameExtensions `json:"createdBy,omitempty"`
	ReferencedRuleCount int                      `json:"referencedRuleCount,omitempty"`
}

func (service *Service) Get(ruleLabelID int) (*RuleLabels, error) {
	var ruleLabel RuleLabels
	err := service.Client.Read(fmt.Sprintf("%s/%d", ruleLabelsEndpoint, ruleLabelID), &ruleLabel)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning rule label from Get: %d", ruleLabel.ID)
	return &ruleLabel, nil
}

func (service *Service) GetRuleLabelByName(labelName string) (*RuleLabels, error) {
	var ruleLabels []RuleLabels
	err := service.Client.Read(ruleLabelsEndpoint, &ruleLabels)
	if err != nil {
		return nil, err
	}
	for _, ruleLabel := range ruleLabels {
		if strings.EqualFold(ruleLabel.Name, labelName) {
			return &ruleLabel, nil
		}
	}
	return nil, fmt.Errorf("no rule label found with name: %s", labelName)
}

func (service *Service) Create(ruleLabelID *RuleLabels) (*RuleLabels, *http.Response, error) {
	resp, err := service.Client.Create(ruleLabelsEndpoint, *ruleLabelID)
	if err != nil {
		return nil, nil, err
	}

	createdRuleLabel, ok := resp.(*RuleLabels)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a rule label pointer")
	}

	log.Printf("returning new rule label from create: %d", createdRuleLabel.ID)
	return createdRuleLabel, nil, nil
}

func (service *Service) Update(ruleLabelID int, ruleLabels *RuleLabels) (*RuleLabels, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", ruleLabelsEndpoint, ruleLabelID), *ruleLabels)
	if err != nil {
		return nil, nil, err
	}
	updatedRuleLabel, _ := resp.(*RuleLabels)

	log.Printf("returning updates rule label from update: %d", updatedRuleLabel.ID)
	return updatedRuleLabel, nil, nil
}

func (service *Service) Delete(ruleLabelID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", ruleLabelsEndpoint, ruleLabelID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
