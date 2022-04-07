package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDlpWebRules_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_dlp_web_rules.test-dlp-rule"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDlpWebRulesBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDlpWebRules(resourceName),
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

func testAccCheckDataSourceDlpWebRulesBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_dlp_web_rules" "test-dlp-rule" {
	name = "test-dlp-rule-%s"
	description = "test-dlp-rule-%s"
	action = "ALLOW"
	state = "ENABLED"
	order = 1
	rank = 7
	protocols = [ "HTTPS_RULE", "HTTP_RULE" ]
	without_content_inspection = false
	match_only = false
	ocr_enabled = true
	min_size = 0
	zscaler_incident_reciever = true
}

data "zia_dlp_web_rules" "test-dlp-rule" {
	name = zia_dlp_web_rules.test-dlp-rule.name
}
	`, rName, rDesc)
}

func testAccDataSourceDlpWebRules(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
