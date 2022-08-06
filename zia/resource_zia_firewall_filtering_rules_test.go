package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/filteringrules"
)

func TestAccResourceFirewallFilteringRuleBasic(t *testing.T) {
	var rules filteringrules.FirewallFilteringRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FirewallFilteringRules)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallFilteringRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallFilteringRuleConfigure(resourceTypeAndName, generatedName, variable.FWRuleResourceDescription, variable.FWRuleResourceAction, variable.FWRuleResourceState, variable.FWRuleEnableLogging),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallFilteringRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWRuleResourceState),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "order"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "nw_services.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "departments.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enable_full_logging", strconv.FormatBool(variable.FWRuleEnableLogging)),
				),
			},

			// Update test
			{
				Config: testAccCheckFirewallFilteringRuleConfigure(resourceTypeAndName, generatedName, variable.FWRuleResourceDescription, variable.FWRuleResourceAction, variable.FWRuleResourceState, variable.FWRuleEnableLogging),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallFilteringRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWRuleResourceState),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "order"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "nw_services.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "departments.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enable_full_logging", strconv.FormatBool(variable.FWRuleEnableLogging)),
				),
			},
		},
	})
}

func testAccCheckFirewallFilteringRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FirewallFilteringRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.filteringrules.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("firewall filtering rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckFirewallFilteringRuleExists(resource string, rule *filteringrules.FirewallFilteringRules) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
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

func testAccCheckFirewallFilteringRuleConfigure(resourceTypeAndName, generatedName, description, action, state string, enableLogging bool) string {
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
resource "%s" "%s" {
    name = "tf-acc-test-%s"
    description = "%s"
    action = "%s"
    state = "%s"
    order = 1
	enable_full_logging = "%s"
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

data "%s" "%s" {
	id = "${%s.id}"
  }

`,
		// resource variables
		resourcetype.FirewallFilteringRules,
		generatedName,
		generatedName,
		description,
		action,
		state,
		strconv.FormatBool(enableLogging),

		// data source variables
		resourcetype.FirewallFilteringRules,
		generatedName,
		resourceTypeAndName,
	)
}
