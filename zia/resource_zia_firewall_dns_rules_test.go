package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewalldnscontrolpolicies"
)

func TestAccResourceFirewallDNSRulesBasic(t *testing.T) {
	var rules firewalldnscontrolpolicies.FirewallDNSRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FirewallDNSRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, "tf-acc-test-"+ruleLabelGeneratedName, variable.RuleLabelDescription)

	// Generate Source IP Group HCL Resource
	sourceIPGroupTypeAndName, _, sourceIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringSourceGroup)
	sourceIPGroupHCL := testAccCheckFWIPSourceGroupsConfigure(sourceIPGroupTypeAndName, "tf-acc-test-"+sourceIPGroupGeneratedName, variable.FWSRCGroupDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallDNSRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallDNSRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.FWDNSRuleDescription, variable.FWDNSAction, variable.FWDNSState, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDNSRulesExists(resourceTypeAndName, &rules),
					testAccCheckFirewallDNSRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWDNSRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWDNSAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWDNSState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_countries.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "source_countries.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "departments.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.0.id.#", "2"),
				),
			},

			// Update test
			{
				Config: testAccCheckFirewallDNSRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.FWDNSRuleDescription, variable.FWDNSAction, variable.FWDNSState, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDNSRulesExists(resourceTypeAndName, &rules),
					testAccCheckFirewallDNSRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWDNSRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWDNSAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWDNSState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_countries.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "source_countries.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "departments.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.0.id.#", "2"),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckFirewallDNSRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FirewallDNSRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := firewalldnscontrolpolicies.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("firewall dns rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckFirewallDNSRulesExists(resource string, rule *firewalldnscontrolpolicies.FirewallDNSRules) resource.TestCheckFunc {
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
		service := apiClient.Service

		receivedRule, err := firewalldnscontrolpolicies.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFirewallDNSRulesConfigure(resourceTypeAndName, generatedName, name, description, action, state string, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s

// source ip group resource
%s

// firewall dns rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		sourceIPGroupHCL,
		getFirewallDNSRulesResourceHCL(generatedName, name, description, action, state, ruleLabelTypeAndName, sourceIPGroupTypeAndName),

		// data source variables
		resourcetype.FirewallDNSRules,
		generatedName,
		resourceTypeAndName,
	)
}

func getFirewallDNSRulesResourceHCL(generatedName, name, description, action, state string, ruleLabelTypeAndName, sourceIPGroupTypeAndName string) string {
	return fmt.Sprintf(`

data "zia_location_groups" "sdwan_can" {
	name = "SDWAN_CAN"
}

data "zia_location_groups" "sdwan_usa" {
	name = "SDWAN_USA"
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

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "%s"
	action = "%s"
	state = "%s"
	order = 11
	redirect_ip = "1.2.3.4"
	dest_countries = ["CA", "US"]
	source_countries = ["CA", "US"]
	protocols = ["ANY_RULE"]
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
		id = ["${%s.id}"]
	}
	src_ip_groups {
		id = ["${%s.id}"]
	}
	depends_on = [ %s, %s ]
}
		`,
		// resource variables
		resourcetype.FirewallDNSRules,
		generatedName,
		name,
		description,
		action,
		state,
		ruleLabelTypeAndName,
		sourceIPGroupTypeAndName,
		ruleLabelTypeAndName,
		sourceIPGroupTypeAndName,
	)
}
