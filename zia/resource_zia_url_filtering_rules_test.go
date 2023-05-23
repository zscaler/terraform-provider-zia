package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zia/services/urlfilteringpolicies"
)

func TestAccResourceURLFilteringRulesBasic(t *testing.T) {
	var rules urlfilteringpolicies.URLFilteringRule
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.URLFilteringRules)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckURLFilteringRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckURLFilteringRulesConfigure(resourceTypeAndName, generatedName, variable.URLFilteringRuleDescription, variable.URLFilteringRuleAction, variable.URLFilteringRuleState),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLFilteringRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.URLFilteringRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.URLFilteringRuleAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.URLFilteringRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "url_categories.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "request_methods.#", "9"),
				),
			},

			// Update test
			{
				Config: testAccCheckURLFilteringRulesConfigure(resourceTypeAndName, generatedName, variable.FWRuleResourceDescription, variable.FWRuleResourceAction, variable.FWRuleResourceState),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLFilteringRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.URLFilteringRuleAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.URLFilteringRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "url_categories.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "request_methods.#", "9"),
				),
			},
		},
	})
}

func testAccCheckURLFilteringRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.URLFilteringRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.urlfilteringpolicies.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("url filtering rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckURLFilteringRulesExists(resource string, rule *urlfilteringpolicies.URLFilteringRule) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.urlfilteringpolicies.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckURLFilteringRulesConfigure(resourceTypeAndName, generatedName, description, action, state string) string {
	return fmt.Sprintf(`

data "zia_rule_labels" "global"{
	name = "GLOBAL"
}

data "zia_firewall_filtering_time_window" "work_hours" {
	name = "Work Hours"
}

data "zia_firewall_filtering_time_window" "off_hours" {
	name = "Off Hours"
}

data "zia_department_management" "engineering" {
	name = "Engineering"
}

data "zia_department_management" "marketing" {
	name = "Marketing"
}

data "zia_group_management" "engineering" {
	name = "Engineering"
}

data "zia_group_management" "marketing" {
	name = "Marketing"
}

data "zia_location_management" "au_sydney_branch01" {
	name = "AU-Sydney-Branch01"
}

data "zia_location_management" "au_sydney_branch02" {
	name = "AU-Sydney-Branch02"
}

data "zia_location_groups" "sdwan_can" {
	name = "SDWAN_CAN"
}

data "zia_location_groups" "sdwan_usa" {
	name = "SDWAN_USA"
}

resource "%s" "%s" {
    name = "%s"
    description = "%s"
	action = "%s"
    state = "%s"
    order = 1
	url_categories = ["ANY"]
    protocols = ["ANY_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
    locations {
		id = [data.zia_location_management.au_sydney_branch01.id, data.zia_location_management.au_sydney_branch02.id]
	}
	location_groups {
		id = [data.zia_location_groups.sdwan_can.id, data.zia_location_groups.sdwan_usa.id]
	}
	groups {
		id = [data.zia_group_management.engineering.id, data.zia_group_management.marketing.id]
	}
	departments {
		id = [data.zia_department_management.engineering.id, data.zia_department_management.marketing.id]
	}
	time_windows {
		id = [data.zia_firewall_filtering_time_window.off_hours.id, data.zia_firewall_filtering_time_window.work_hours.id]
	}
	labels {
		id = [data.zia_rule_labels.global.id]
	}
}

data "%s" "%s" {
	name = "${%s.name}"
  }
`,
		// resource variables
		resourcetype.URLFilteringRules,
		generatedName,
		generatedName,
		description,
		action,
		state,

		// data source variables
		resourcetype.URLFilteringRules,
		generatedName,
		resourceTypeAndName,
	)
}
