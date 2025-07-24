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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_control_rules"
)

func TestAccResourceBandwdithControlRules_Basic(t *testing.T) {
	var rules bandwidth_control_rules.BandwidthControlRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.BandwdithControlRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, "tf-acc-test-"+ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBandwdithControlRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBandwdithControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.BandwdithControlRuleDescription, variable.BandwdithControlRulestate, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBandwdithControlRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.BandwdithControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.BandwdithControlRulestate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "zia_bandwidth_classes.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.0.id.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckBandwdithControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.BandwdithControlRuleDescription, variable.BandwdithControlRulestate, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBandwdithControlRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.BandwdithControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.BandwdithControlRulestate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "zia_bandwidth_classes.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.0.id.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceTypeAndName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceTypeAndName)
					}
					if rs.Primary.ID == "" {
						return "", fmt.Errorf("no record ID is set")
					}

					ruleType := rs.Primary.Attributes["type"]
					return fmt.Sprintf("%s:%s", ruleType, rs.Primary.ID), nil
				},
			},
		},
	})
}

func testAccCheckBandwdithControlRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.BandwdithControlRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := bandwidth_control_rules.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("bandwidth control rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckBandwdithControlRulesExists(resource string, rule *bandwidth_control_rules.BandwidthControlRules) resource.TestCheckFunc {
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

		receivedRule, err := bandwidth_control_rules.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckBandwdithControlRulesConfigure(resourceTypeAndName, generatedName, name, description, state, ruleLabelTypeAndName, ruleLabelHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s

// bandwitdh control rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		getBandwdithControlRulesResourceHCL(generatedName, name, description, state, ruleLabelTypeAndName),

		// data source variables
		resourcetype.BandwdithControlRules,
		generatedName,
		resourceTypeAndName,
	)
}

func getBandwdithControlRulesResourceHCL(generatedName, name, description, state, ruleLabelTypeAndName string) string {
	return fmt.Sprintf(`
data "zia_firewall_filtering_time_window" "work_hours" {
	name = "Work Hours"
}

data "zia_bandwidth_classes" "this" {
    name = "WEBCONF"
}
resource "%s" "%s" {
    name 					= "tf-acc-test-%s"
    description 			= "%s"
    state 					= "%s"
    order 					= 1
	rank					= 7
    min_bandwidth = 5
    max_bandwidth = 100
	protocols = ["ANY_RULE"]
    bandwidth_classes  {
        id = [data.zia_bandwidth_classes.this.id]
    }
	time_windows {
		id = [data.zia_firewall_filtering_time_window.work_hours.id]
	}
	labels {
		id = ["${%s.id}"]
	}
}
`,
		// resource variables
		resourcetype.BandwdithControlRules,
		generatedName,
		name,
		description,
		state,
		ruleLabelTypeAndName,
	)
}
