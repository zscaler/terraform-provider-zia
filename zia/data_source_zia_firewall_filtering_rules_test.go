package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccDataSourceFirewallFilteringRule_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FirewallFilteringRules)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallFilteringRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallFilteringRuleConfigure(resourceTypeAndName, generatedName, variable.FWRuleResourceDescription, variable.FWRuleResourceAction, variable.FWRuleResourceState, variable.FWRuleEnableLogging),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "action", resourceTypeAndName, "action"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "state", resourceTypeAndName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "order", resourceTypeAndName, "order"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enable_full_logging", strconv.FormatBool(variable.FWRuleEnableLogging)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "nw_services.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "departments.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "time_windows.#", "1"),
				),
			},
		},
	})
}
