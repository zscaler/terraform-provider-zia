package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
)

func TestAccDataSourceFirewallFilteringRule_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FirewallFilteringRules)

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
		CheckDestroy: testAccCheckFirewallFilteringRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallFilteringRuleConfigure(resourceTypeAndName, generatedName, generatedName, variable.FWRuleResourceDescription, variable.FWRuleResourceAction, variable.FWRuleResourceState, variable.FWRuleEnableLogging, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "action", resourceTypeAndName, "action"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "state", resourceTypeAndName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "device_trust_levels", resourceTypeAndName, "device_trust_levels"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enable_full_logging", strconv.FormatBool(variable.FWRuleEnableLogging)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "nw_services.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "departments.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "groups.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "time_windows.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "labels.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "src_ip_groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "dest_ip_groups.#", "1"),
				),
			},
		},
	})
}
