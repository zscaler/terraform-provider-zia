package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/dlp_web_rules"
)

func TestAccResourceDlpWebRules_basic(t *testing.T) {
	var rules dlp_web_rules.WebDLPRules
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_dlp_web_rules.test-dlp-rule"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDlpWebRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceDlpWebRulesBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDlpWebRulesExists("zia_dlp_web_rules.test-dlp-rule", &rules),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dlp-rule-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-dlp-rule-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "state", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "without_content_inspection", "false"),
					resource.TestCheckResourceAttr(resourceName, "match_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "ocr_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "min_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "zscaler_incident_reciever", "true"),
				),
			},
		},
	})
}

func testAccCheckResourceDlpWebRulesBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_dlp_web_rules" "test-dlp-rule" {
	name 						= "test-dlp-rule-%s"
	description 				= "test-dlp-rule-%s"
	action 						= "ALLOW"
	state 						= "ENABLED"
	order 						= 1
	rank 						= 7
	protocols 					= [ "HTTPS_RULE", "HTTP_RULE" ]
	without_content_inspection 	= false
	match_only 					= false
	ocr_enabled 				= true
	min_size 					= 0
	zscaler_incident_reciever 	= true
}
	`, rName, rDesc)
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

func testAccCheckDlpWebRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_dlp_web_rules" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.dlp_web_rules.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("dlp web rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
