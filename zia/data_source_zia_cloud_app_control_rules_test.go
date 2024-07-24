package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/variable"
)

func TestAccDataSourceCloudAppControlRules_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.CloudAppControlRule)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudAppControlRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCloudAppControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.CloudAppControlRuleDescription, variable.CloudAppControlRuleState, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "actions", resourceTypeAndName, "actions"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "state", resourceTypeAndName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "rank", resourceTypeAndName, "rank"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "applications.#", "2"),
				),
			},
		},
	})
}
