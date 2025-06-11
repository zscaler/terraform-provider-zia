package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceNatControlRules_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.NatControlRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, ruleLabelGeneratedName, variable.RuleLabelDescription)

	// Generate Source IP Group HCL Resource
	sourceIPGroupTypeAndName, _, sourceIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringSourceGroup)
	sourceIPGroupHCL := testAccCheckFWIPSourceGroupsConfigure(sourceIPGroupTypeAndName, sourceIPGroupGeneratedName, variable.FWSRCGroupDescription)

	// Generate Destination IP Group HCL Resource
	dstIPGroupTypeAndName, _, dstIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringDestinationGroup)
	dstIPGroupHCL := testAccCheckFWIPDestinationGroupsConfigure(dstIPGroupTypeAndName, dstIPGroupGeneratedName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatControlRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNatControlRulesConfigure(
					resourceTypeAndName, generatedName, generatedName, variable.FWIPSRuleDescription,
					variable.FWIPSState, variable.FWRuleEnableLogging, ruleLabelTypeAndName, ruleLabelHCL,
					sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "state", resourceTypeAndName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "order", resourceTypeAndName, "order"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "redirect_ip", resourceTypeAndName, "redirect_ip"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "redirect_port", resourceTypeAndName, "redirect_port"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enable_full_logging", strconv.FormatBool(variable.FWRuleEnableLogging)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "nw_services.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "labels.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "src_ip_groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "dest_ip_groups.#", "1"),
				),
			},
		},
	})
}
