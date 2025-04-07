package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceFileTypeControlRules_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FileTypeControlRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuleLabelsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFileTypeControlRulesConfigure(
					resourceTypeAndName, generatedName, generatedName, variable.FileTypeControlRuleDescription,
					variable.FileTypeControlRuleAction, variable.FileTypeControlRuleState, ruleLabelTypeAndName, ruleLabelHCL,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "state", resourceTypeAndName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "order", resourceTypeAndName, "order"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "filtering_action", resourceTypeAndName, "filtering_action"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "operation", resourceTypeAndName, "operation"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "active_content", resourceTypeAndName, "active_content"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "unscannable", resourceTypeAndName, "unscannable"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "device_trust_levels.#", "4"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "file_types.#", "4"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "protocols.#", "4"),
					// resource.TestCheckResourceAttr(dataSourceTypeAndName, "departments.#", "2"),
					// resource.TestCheckResourceAttr(dataSourceTypeAndName, "groups.#", "2"),
					// resource.TestCheckResourceAttr(dataSourceTypeAndName, "time_windows.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "labels.#", "1"),
				),
			},
		},
	})
}
