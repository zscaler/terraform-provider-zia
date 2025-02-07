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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sslinspection"
)

func TestAccResourceSSLInspectionRules_Basic(t *testing.T) {
	var rules sslinspection.SSLInspectionRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.SSLInspectionRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, "tf-acc-test-"+ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSSLInspectionRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSSLInspectionRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.SSLInspectionRuleDescription, variable.SSLInspectionRuleState, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSLInspectionRulesExists(resourceTypeAndName, &rules),
					testAccCheckSSLInspectionRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.SSLInspectionRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.SSLInspectionRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "road_warrior_for_kerberos", strconv.FormatBool(variable.RoadWarriorKerberos)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "departments.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.0.id.#", "2"),
				),
			},

			// Update test
			{
				Config: testAccCheckSSLInspectionRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.SSLInspectionRuleDescription, variable.SSLInspectionRuleState, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSLInspectionRulesExists(resourceTypeAndName, &rules),
					testAccCheckSSLInspectionRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.SSLInspectionRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.SSLInspectionRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "road_warrior_for_kerberos", strconv.FormatBool(variable.RoadWarriorKerberos)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
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

func testAccCheckSSLInspectionRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.SSLInspectionRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := sandbox_rules.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("ssl inspection rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckSSLInspectionRulesExists(resource string, rule *sslinspection.SSLInspectionRules) resource.TestCheckFunc {
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

		receivedRule, err := sslinspection.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckSSLInspectionRulesConfigure(resourceTypeAndName, generatedName, name, description, state string, ruleLabelTypeAndName, ruleLabelHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s

// ssl inspection rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		getSSLInspectionRulesResourceHCL(generatedName, name, description, state, ruleLabelTypeAndName),

		// data source variables
		resourcetype.SSLInspectionRules,
		generatedName,
		resourceTypeAndName,
	)
}

func getSSLInspectionRulesResourceHCL(generatedName, name, description, state string, ruleLabelTypeAndName string) string {
	return fmt.Sprintf(`


data "zia_location_groups" "sdwan_can" {
	name = "SDWAN_CAN"
}

data "zia_location_groups" "sdwan_usa" {
	name = "SDWAN_USA"
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
	state = "%s"
	order = 6
	rank  = 7
	road_warrior_for_kerberos 	 = true
	cloud_applications           = ["CHATGPT_AI", "ANDI"]
	platforms                    = ["SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX"]
	action {
		type                         = "DECRYPT"

		decrypt_sub_actions {
			server_certificates                 = "ALLOW"
			ocsp_check                          = true
			block_ssl_traffic_with_no_sni_enabled = true
			min_client_tls_version              = "CLIENT_TLS_1_0"
			min_server_tls_version              = "SERVER_TLS_1_0"
			block_undecrypt                    = true
			http2_enabled                       = false
		}
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
	labels {
		id = ["${%s.id}"]
	}
	depends_on = [ %s ]
}
		`,
		// resource variables
		resourcetype.SSLInspectionRules,
		generatedName,
		name,
		description,
		state,
		ruleLabelTypeAndName,
		ruleLabelTypeAndName,
	)
}
