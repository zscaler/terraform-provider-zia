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
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "2"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "file_types.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "cloud_applications.#", "4"),
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
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "2"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "file_types.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "cloud_applications.#", "4"),
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

resource "%s" "%s" {
	name 						= "tf-acc-test-%s"
	description 				= "%s"
    action 						= "%s"
    state 						= "%s"
	order 						= 1
	rank 						= 7
	protocols 					= [ "HTTPS_RULE", "HTTP_RULE" ]
	cloud_applications 			= ["ZENDESK", "LUCKY_ORANGE", "MICROSOFT_POWERAPPS", "MICROSOFTLIVEMEETING"]
	without_content_inspection 	= false
	match_only 					= false
	ocr_enabled 				= false
	file_types                  = []
	min_size 					= 20
	zscaler_incident_reciever 	= true
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

		// data source variables
		resourcetype.DLPWebRules,
		generatedName,
		resourceTypeAndName,
	)
}
