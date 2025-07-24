package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceBandwdithControlRules_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.BandwdithControlRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBandwdithControlRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBandwdithControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.BandwdithControlRuleDescription, variable.BandwdithControlRulestate, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "state", resourceTypeAndName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "min_bandwidth", resourceTypeAndName, "min_bandwidth"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "max_bandwidth", resourceTypeAndName, "max_bandwidth"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "bandwidth_classes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "labels.#", "1"),
				),
			},
		},
	})
}
