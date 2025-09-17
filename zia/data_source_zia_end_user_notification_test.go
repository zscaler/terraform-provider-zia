package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
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
					resource.TestCheckResourceAttr(resourceName, "notification_type", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "display_reason", strconv.FormatBool(variable.DisplayCompReason)),
					resource.TestCheckResourceAttr(resourceName, "display_company_name", strconv.FormatBool(variable.DisplayCompName)),
					resource.TestCheckResourceAttr(resourceName, "display_company_logo", strconv.FormatBool(variable.DisplayCompLogo)),
					resource.TestCheckResourceAttr(resourceName, "custom_text", "Website blocked"),
					resource.TestCheckResourceAttr(resourceName, "url_cat_review_enabled", strconv.FormatBool(variable.UrlCatReviewEnabled)),
					resource.TestCheckResourceAttr(resourceName, "url_cat_review_submit_to_security_cloud", strconv.FormatBool(variable.EunUrlCatReviewSubmitToSecurityCloud)),
					resource.TestCheckResourceAttr(resourceName, "security_review_enabled", strconv.FormatBool(variable.EunSecurityReviewEnabled)),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_enabled", strconv.FormatBool(variable.EunWebDlpReviewEnabled)),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_submit_to_security_cloud", strconv.FormatBool(variable.EunWebDlpReviewSubmitToSecurityCloud)),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_custom_location", "https://redirect.acme.com"),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_text", "Click to request policy review."),
				),
			},
		},
	})
}

var testAccCheckDataSourceEndUserNotificationConfig_basic = `
data "zia_end_user_notification" "this" {}
`
