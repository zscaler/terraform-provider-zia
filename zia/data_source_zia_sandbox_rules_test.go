package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceSandboxRules_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.SandboxRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuleLabelsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSandboxRulesConfigure(
					resourceTypeAndName, generatedName, generatedName, variable.SandboxRuleDescription,
					variable.SandboxState, variable.SandboxAction, ruleLabelTypeAndName, ruleLabelHCL,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "state", resourceTypeAndName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ba_rule_action", resourceTypeAndName, "ba_rule_action"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "first_time_enable", resourceTypeAndName, "first_time_enable"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "first_time_operation", resourceTypeAndName, "first_time_operation"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ml_action_enabled", resourceTypeAndName, "ml_action_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "order", resourceTypeAndName, "order"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "ba_policy_categories.#", "4"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "file_types.#", "19"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "protocols.#", "4"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "labels.#", "1"),
				),
			},
		},
	})
}
