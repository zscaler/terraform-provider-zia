package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudappcontrol"
)

func TestAccResourceCloudAppControlRulesBasic(t *testing.T) {
	var rules cloudappcontrol.WebApplicationRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.CloudAppControlRule)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, "tf-acc-test-"+ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudAppControlRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCloudAppControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.CloudAppControlRuleDescription, variable.CloudAppControlRuleState, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudAppControlRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.CloudAppControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.CloudAppControlRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.CloudAppControlRuleType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "applications.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "departments.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.0.id.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckCloudAppControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.CloudAppControlRuleDescription, variable.CloudAppControlRuleState, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudAppControlRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.CloudAppControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.CloudAppControlRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.CloudAppControlRuleType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "applications.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "departments.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.0.id.#", "1"),
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

func testAccCheckCloudAppControlRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.CloudAppControlRule {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		ruleType := rs.Primary.Attributes["type"]
		rule, err := cloudappcontrol.GetByRuleID(context.Background(), service, ruleType, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("cloud app control rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckCloudAppControlRulesExists(resource string, rule *cloudappcontrol.WebApplicationRules) resource.TestCheckFunc {
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

		ruleType := rs.Primary.Attributes["type"]

		apiClient := testAccProvider.Meta().(*Client)
		service := apiClient.Service

		receivedRule, err := cloudappcontrol.GetByRuleID(context.Background(), service, ruleType, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckCloudAppControlRulesConfigure(resourceTypeAndName, generatedName, name, description, state, ruleLabelTypeAndName, ruleLabelHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s

// cloud app control rule resource
%s

data "%s" "%s" {
	type = "STREAMING_MEDIA"
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		getCloudAppControlRuleResourceHCL(generatedName, name, description, state, ruleLabelTypeAndName),

		// data source variables
		resourcetype.CloudAppControlRule,
		generatedName,
		resourceTypeAndName,
	)
}

func getCloudAppControlRuleResourceHCL(generatedName, name, description, state, ruleLabelTypeAndName string) string {
	// Generate current time + 5 minutes for validity_start_time
	validityStartTime := time.Now().Add(5 * time.Minute).UTC().Format(time.RFC1123)

	// Generate time 365 days + 5 minutes from now for validity_end_time
	validityEndTime := time.Now().AddDate(1, 0, 0).Add(5 * time.Minute).UTC().Format(time.RFC1123)

	return fmt.Sprintf(`
data "zia_firewall_filtering_time_window" "work_hours" {
	name = "Work Hours"
}

data "zia_department_management" "engineering" {
	name = "Engineering"
}

data "zia_group_management" "engineering" {
	name = "Engineering"
}

data "zia_location_groups" "sdwan_can" {
	name = "SDWAN_CAN"
}

resource "%s" "%s" {
    name 					= "tf-acc-test-%s"
    description 			= "%s"
    state 					= "%s"
	type 					= "STREAMING_MEDIA"
	actions 				= ["ALLOW_STREAMING_VIEW_LISTEN", "ALLOW_STREAMING_UPLOAD"]
	applications 			= ["YOUTUBE", "GOOGLE_STREAMING"]
    order 					= 1
	rank					= 7
	device_trust_levels 	= [	"UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST" ]
	user_agent_types 		= [	"OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE", "OTHER" ]
	user_risk_score_levels 	= ["LOW", "MEDIUM", "HIGH", "CRITICAL"]
	location_groups {
		id = [data.zia_location_groups.sdwan_can.id]
	}
	groups {
		id = [data.zia_group_management.engineering.id]
	}
	departments {
		id = [data.zia_department_management.engineering.id]
	}
	time_windows {
		id = [data.zia_firewall_filtering_time_window.work_hours.id]
	}
	labels {
		id = ["${%s.id}"]
	}
	enforce_time_validity = true
	validity_time_zone_id = "US/Pacific"
    validity_start_time = "%s"
    validity_end_time = "%s"
}
`,
		// resource variables
		resourcetype.CloudAppControlRule,
		generatedName,
		name,
		description,
		state,
		ruleLabelTypeAndName,
		validityStartTime,
		validityEndTime,
	)
}
