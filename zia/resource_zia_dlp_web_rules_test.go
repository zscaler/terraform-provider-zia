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
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_web_rules"
)

func TestAccResourceDlpWebRulesBasic(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceTypeAndName, "match_only", strconv.FormatBool(variable.DLPMatchOnly)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ocr_enabled", strconv.FormatBool(variable.DLPOCREnabled)),
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
					resource.TestCheckResourceAttr(resourceTypeAndName, "match_only", strconv.FormatBool(variable.DLPMatchOnly)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ocr_enabled", strconv.FormatBool(variable.DLPOCREnabled)),
				),
			},
		},
	})
}

func testAccCheckDlpWebRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPWebRules {
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
		receivedRule, err := apiClient.dlp_web_rules.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckDlpWebRulesConfigure(resourceTypeAndName, generatedName, description, action, state string) string {
	return fmt.Sprintf(`

data "zia_url_categories" "corporate_marketing"{
	id = "CORPORATE_MARKETING"
}

data "zia_url_categories" "finance"{
	id = "FINANCE"
}

data "zia_dlp_engines" "pci" {
	name = "PCI"
}

data "zia_dlp_engines" "glba" {
	name = "GLBA"
}
data "zia_rule_labels" "can"{
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

data "zia_location_groups" "sdwan_can" {
	name = "SDWAN_CAN"
}

data "zia_location_groups" "sdwan_usa" {
	name = "SDWAN_USA"
}

resource "%s" "%s" {
	name 						= "tf-acc-test-%s"
	description 				= "%s"
    action 						= "%s"
    state 						= "%s"
	order 						= "%d"
	rank 						= 7
	protocols                 = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
	without_content_inspection 	= false
	match_only 					= false
	ocr_enabled 				= false
	file_types                  = []
	min_size 					= 20
	zscaler_incident_receiver 	= true
	location_groups {
		id = [data.zia_location_groups.sdwan_usa.id, data.zia_location_groups.sdwan_can.id]
	}
	groups {
		id = [data.zia_group_management.engineering.id, data.zia_group_management.marketing.id]
	}
	departments {
		id = [data.zia_department_management.engineering.id, data.zia_department_management.marketing.id]
	}
	time_windows {
		id = [data.zia_firewall_filtering_time_window.work_hours.id, data.zia_firewall_filtering_time_window.off_hours.id]
	}
	labels {
		id = [data.zia_rule_labels.can.id]
	}
	dlp_engines {
		id = [data.zia_dlp_engines.pci.id, data.zia_dlp_engines.glba.id]
	}
	url_categories {
		id = [data.zia_url_categories.corporate_marketing.val, data.zia_url_categories.finance.val]
	}
}

data "%s" "%s" {
	id = "${%s.id}"
}

`,
		// resource variables
		resourcetype.DLPWebRules,
		generatedName,
		generatedName,
		description,
		action,
		state,
		variable.DLPRuleOrder,

		// data source variables
		resourcetype.DLPWebRules,
		generatedName,
		resourceTypeAndName,
	)
}
