package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/filteringrules"
)

func TestAccFirewallFilteringRule_basic(t *testing.T) {
	var rules filteringrules.FirewallFilteringRules
	rName := "test-fw-rule-" + acctest.RandString(5)
	rDesc := "test-fw-rule-" + acctest.RandString(20)
	resourceName := "zia_firewall_filtering_rule.test-fw-rule"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallFilteringRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallFilteringRuleBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallFilteringRuleExists("zia_firewall_filtering_rule.test-fw-rule", &rules),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "state", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "order", "1"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFirewallFilteringRuleBasic(rName, rDesc string) string {
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
resource "zia_firewall_filtering_rule" "test-fw-rule" {
	name = "%s"
	description = "%s"
	action = "ALLOW"
	state = "ENABLED"
	order = 1
	src_ips=[]
	dest_addresses=[]
	dest_ip_categories=[]
	dest_countries=[]
	nw_applications=[]
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
	locations {
		id = []
	}
	location_groups {
		id = []
	}
	users {
		id = []
	}
	labels {
		id = []
	}
	src_ip_groups {
		id = []
	}
	dest_ip_groups {
		id = []
	}
	app_service_groups {
		id = []
	}
	app_services {
		id = []
	}
	nw_application_groups {
		id = []
	}
	nw_service_groups {
		id = []
	}

}
	`, rName, rDesc)
}

func testAccCheckFirewallFilteringRuleExists(resource string, rule *filteringrules.FirewallFilteringRules) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("firewall rule not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no firewall rule ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedRule, err := apiClient.filteringrules.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFirewallFilteringRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_firewall_filtering_rule" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		foundRule, err := apiClient.filteringrules.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if foundRule != nil {
			return fmt.Errorf("firewall filtering rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
