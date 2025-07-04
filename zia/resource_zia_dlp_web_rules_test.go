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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
)

func TestAccResourceDlpWebRules_Basic(t *testing.T) {
	var rules dlp_web_rules.WebDLPRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPWebRules)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDlpWebRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDlpWebRulesConfigure(resourceTypeAndName, generatedName, variable.DLPWebRuleDesc, variable.DLPRuleResourceAction, variable.DLPRuleResourceState),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDlpWebRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPWebRuleDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.DLPRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.DLPRuleResourceState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "3"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "without_content_inspection", strconv.FormatBool(variable.DLPRuleContentInspection)),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "match_only", strconv.FormatBool(variable.DLPMatchOnly)),
				),
			},

			// Update test
			{
				Config: testAccCheckDlpWebRulesConfigure(resourceTypeAndName, generatedName, variable.DLPWebRuleDesc, variable.DLPRuleResourceAction, variable.DLPRuleResourceState),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDlpWebRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPWebRuleDesc),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.DLPRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.DLPRuleResourceState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "3"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "without_content_inspection", strconv.FormatBool(variable.DLPRuleContentInspection)),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "match_only", strconv.FormatBool(variable.DLPMatchOnly)),
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

func testAccCheckDlpWebRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPWebRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := dlp_web_rules.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("dlp web rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckDlpWebRulesExists(resource string, rule *dlp_web_rules.WebDLPRules) resource.TestCheckFunc {
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

		receivedRule, err := dlp_web_rules.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule
		return nil
	}
}

func testAccCheckDlpWebRulesConfigure(resourceTypeAndName, name, description, action, state string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`

	data "zia_url_categories" "corporate_marketing"{
	id = "CORPORATE_MARKETING"
}

data "zia_url_categories" "finance"{
	id = "FINANCE"
}

data "zia_firewall_filtering_time_window" "work_hours" {
	name = "Work Hours"
}

data "zia_firewall_filtering_time_window" "off_hours" {
	name = "Off Hours"
}

data "zia_location_groups" "sdwan_can" {
	name = "SDWAN_CAN"
}

data "zia_location_groups" "sdwan_usa" {
	name = "SDWAN_USA"
}

data "zia_dlp_notification_templates" "this" {
  name = "BD_DLP_Notification_Template"
}

resource "%s" "%s" {
	name 						= "tf-acc-test-%s"
	description 				= "%s"
    action 						= "%s"
    state 						= "%s"
	order 						= 1
	rank 						= 0
	protocols                 = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
	without_content_inspection 	= false
	zscaler_incident_receiver = true
	user_risk_score_levels = ["LOW", "MEDIUM", "HIGH", "CRITICAL"]
	severity = "RULE_SEVERITY_HIGH"
	external_auditor_email     = "auditor@acme.com"
	location_groups {
		id = [data.zia_location_groups.sdwan_usa.id, data.zia_location_groups.sdwan_can.id]
	}
	time_windows {
		id = [data.zia_firewall_filtering_time_window.work_hours.id, data.zia_firewall_filtering_time_window.off_hours.id]
	}
	url_categories {
		id = [data.zia_url_categories.corporate_marketing.val, data.zia_url_categories.finance.val]
	}
    notification_template {
     	id = data.zia_dlp_notification_templates.this.id
    }
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }

`,
		// Resource type and name for the dlp web rule
		resourcetype.DLPWebRules,
		resourceName,
		name,
		description,
		action,
		state,

		// Data source type and name
		resourcetype.DLPWebRules,
		resourceName,

		// Reference to the resource
		resourcetype.DLPWebRules,
		resourceName,
	)
}
