package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_dlp_rules"
)

func TestAccResourceEndpointDLPSubRules_Basic(t *testing.T) {
	var rule endpoint_dlp_rules.EndpointDlpRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPEndpointSubRules)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDLPSubRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEndpointDLPSubRulesConfigure(resourceTypeAndName, generatedName, variable.DLPEndpointSubRuleDescription, variable.DLPEndpointSubRuleAction, variable.DLPEndpointSubRuleState),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPSubRulesExists(resourceTypeAndName, &rule),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointSubRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.DLPEndpointSubRuleAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.DLPEndpointSubRuleState),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "parent_rule"),
				),
			},

			// Update test
			{
				Config: testAccCheckEndpointDLPSubRulesConfigure(resourceTypeAndName, generatedName, variable.DLPEndpointSubRuleDescription, variable.DLPEndpointSubRuleAction, variable.DLPEndpointSubRuleState),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPSubRulesExists(resourceTypeAndName, &rule),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointSubRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.DLPEndpointSubRuleAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.DLPEndpointSubRuleState),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "parent_rule"),
				),
			},
		},
	})
}

func testAccCheckEndpointDLPSubRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPEndpointSubRules {
			continue
		}

		subID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		parentID, err := strconv.Atoi(rs.Primary.Attributes["parent_rule"])
		if err != nil {
			// Without a parent reference there is nothing to look up.
			continue
		}

		parent, err := endpoint_dlp_rules.Get(context.Background(), service, parentID)
		if err != nil {
			// Parent is gone, so its sub-rules are gone too.
			continue
		}

		for _, sr := range parent.SubRules {
			if sr.ID == subID {
				return fmt.Errorf("endpoint dlp sub-rule with id %d still exists under parent %d", subID, parentID)
			}
		}
	}

	return nil
}

func testAccCheckEndpointDLPSubRulesExists(resource string, rule *endpoint_dlp_rules.EndpointDlpRules) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		subID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		parentID, err := strconv.Atoi(rs.Primary.Attributes["parent_rule"])
		if err != nil {
			return fmt.Errorf("failed to read parent_rule for resource %s: %s", resource, err)
		}

		apiClient := testAccProvider.Meta().(*Client)
		service := apiClient.Service

		parent, err := endpoint_dlp_rules.Get(context.Background(), service, parentID)
		if err != nil {
			return fmt.Errorf("failed fetching parent rule %d for resource %s. Received error: %s", parentID, resource, err)
		}

		for i := range parent.SubRules {
			if parent.SubRules[i].ID == subID {
				*rule = parent.SubRules[i]
				return nil
			}
		}
		return fmt.Errorf("endpoint dlp sub-rule %d not found under parent %d", subID, parentID)
	}
}

func testAccCheckEndpointDLPSubRulesConfigure(resourceTypeAndName, name, description, action, state string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`

resource "zia_endpoint_dlp_rules" "parent" {
	name                 = "tf-acc-test-%s-parent"
	description          = "parent rule for endpoint dlp sub-rule acc test"
	action               = "%s"
	state                = "%s"
	order                = 1
	rank                 = 0
	data_transfer_method = "REMOVABLE_DRIVE_TRANSFER"
	severity             = "RULE_SEVERITY_HIGH"
}

resource "%s" "%s" {
	name                 = "tf-acc-test-%s"
	description          = "%s"
	action               = "%s"
	state                = "%s"
	order                = 1
	rank                 = 0
	parent_rule          = zia_endpoint_dlp_rules.parent.rule_id
	data_transfer_method = "REMOVABLE_DRIVE_TRANSFER"
	severity             = "RULE_SEVERITY_HIGH"
}
`,
		// Parent rule name
		name,
		action,
		state,

		// Sub-rule resource type and name
		resourcetype.DLPEndpointSubRules,
		resourceName,
		name,
		description,
		action,
		state,
	)
}
