package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEndUserNotification_Basic(t *testing.T) {
	resourceName := "data.zia_end_user_notification.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceEndUserNotificationConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "aup_frequency", "NEVER"),
					resource.TestCheckResourceAttr(resourceName, "aup_day_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "aup_message", ""),
					resource.TestCheckResourceAttr(resourceName, "notification_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "display_reason", "false"),
					resource.TestCheckResourceAttr(resourceName, "display_company_name", "false"),
					resource.TestCheckResourceAttr(resourceName, "display_company_logo", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_text", "Website blocked"),
					resource.TestCheckResourceAttr(resourceName, "url_cat_review_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "url_cat_review_submit_to_security_cloud", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_review_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_custom_location", "https://redirect.acme.com"),
				),
			},
		},
	})
}

var testAccCheckDataSourceEndUserNotificationConfig_basic = `
data "zia_end_user_notification" "this" {}
`
