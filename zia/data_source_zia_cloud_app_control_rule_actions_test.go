package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCloudAppControlRuleActions_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceCloudAppControlRuleActionsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zia_cloud_app_control_rule_actions.this1", "id"),
					resource.TestCheckResourceAttrSet("data.zia_cloud_app_control_rule_actions.this1", "available_actions.0"),
					resource.TestCheckResourceAttrSet("data.zia_cloud_app_control_rule_actions.this2", "id"),
					resource.TestCheckResourceAttrSet("data.zia_cloud_app_control_rule_actions.this2", "available_actions.0"),
				),
			},
		},
	})
}

var testAccCheckDataSourceCloudAppControlRuleActionsConfig_basic = `
data "zia_cloud_app_control_rule_actions" "this1" {
  type       = "ENTERPRISE_COLLABORATION"
  cloud_apps = ["SLACK"]
}

data "zia_cloud_app_control_rule_actions" "this2" {
  type       = "FILE_SHARE"
  cloud_apps = ["GDRIVE"]
}
`
