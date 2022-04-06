package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFirewallFilteringRule_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_firewall_filtering_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFirewallFilteringRuleBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFirewallFilteringRule(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-rule-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-rule-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "state", "ENABLED"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDataSourceFirewallFilteringRuleBasic(rName, rDesc string) string {
	return fmt.Sprintf(`
data "zia_firewall_filtering_network_service" "zscaler_proxy_nw_services" {
	name = "ZSCALER_PROXY_NW_SERVICES"
}
data "zia_department_management" "engineering" {
	name = "Engineering"
}
data "zia_group_management" "normal_internet" {
	name = "Normal_Internet"
}

data "zia_firewall_filtering_time_window" "work_hours" {
	name = "Work hours"
}
resource "zia_firewall_filtering_rule" "test" {
	name = "test-fw-rule-%s"
	description = "test-fw-rule-%s"
	action = "ALLOW"
	state = "ENABLED"
	order = 1
	nw_services {
		id = [ data.zia_firewall_filtering_network_service.zscaler_proxy_nw_services.id ]
	}
	departments {
		id = [ data.zia_department_management.engineering.id ]
	}
	groups {
		id = [ data.zia_group_management.normal_internet.id ]
	}
	time_windows {
		id = [ data.zia_firewall_filtering_time_window.work_hours.id ]
	}
}
data "zia_firewall_filtering_rule" "test" {
	name = zia_firewall_filtering_rule.test.name
}
	`, rName, rDesc)
}

func testAccDataSourceFirewallFilteringRule(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
