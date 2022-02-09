package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFirewallFilteringRule_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceFirewallFilteringRuleConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_rule.zscaler_proxy", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_rule.o365", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceFirewallFilteringRuleConfig_basic = `
data "zia_firewall_filtering_rule" "zscaler_proxy" {
    name = "Zscaler Proxy Traffic"
}

data "zia_firewall_filtering_rule" "o365" {
    name = "Office 365 One Click Rule"
}
`
