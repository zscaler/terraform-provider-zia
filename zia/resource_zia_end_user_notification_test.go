package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccResourceEndUserNotificationBasic(t *testing.T) {
	resourceName := "zia_end_user_notification.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceEndUserNotificationDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create the resource with specific values
			{
				Config: testAccResourceEndUserNotificationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "aup_custom_frequency", "0"),
					resource.TestCheckResourceAttr(resourceName, "aup_day_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "aup_frequency", "NEVER"),
					resource.TestCheckResourceAttr(resourceName, "caution_again_after", "300"),
					resource.TestCheckResourceAttr(resourceName, "caution_custom_text", "Proceeding to visit the site may violate your company policy. Press the \"Continue\" button to access the site anyway or press the \"Back\" button on your browser to go back"),
					resource.TestCheckResourceAttr(resourceName, "caution_per_domain", "true"),
					resource.TestCheckResourceAttr(resourceName, "custom_text", "Website blocked"),
					resource.TestCheckResourceAttr(resourceName, "display_company_logo", strconv.FormatBool(variable.DisplayCompLogo)),
					resource.TestCheckResourceAttr(resourceName, "display_company_name", strconv.FormatBool(variable.DisplayCompName)),
					resource.TestCheckResourceAttr(resourceName, "display_reason", strconv.FormatBool(variable.DisplayCompReason)),
					resource.TestCheckResourceAttr(resourceName, "notification_type", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "org_policy_link", "http://24326813.zscalerthree.net/policy.html"),
					resource.TestCheckResourceAttr(resourceName, "quarantine_custom_notification_text", "We are checking this file for a potential security risk. The file you attempted to download is being analyzed for your protection.\nIt is not blocked. The analysis can take up to 10 minutes, depending on the size and type of the file. If safe, your file downloads automatically.\nIf unsafe, the file will be blocked.\n"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://redirect.acme.com"),
					resource.TestCheckResourceAttr(resourceName, "security_review_enabled", strconv.FormatBool(variable.EunSecurityReviewEnabled)),
					resource.TestCheckResourceAttr(resourceName, "security_review_text", "Click to request security review."),
					resource.TestCheckResourceAttr(resourceName, "security_review_custom_location", "https://redirect.acme.com"),
					resource.TestCheckResourceAttr(resourceName, "support_email", "support@24326813.zscalerthree.net"),
					resource.TestCheckResourceAttr(resourceName, "support_phone", "+91-9000000000"),
					resource.TestCheckResourceAttr(resourceName, "url_cat_review_enabled", strconv.FormatBool(variable.UrlCatReviewEnabled)),
					resource.TestCheckResourceAttr(resourceName, "url_cat_review_submit_to_security_cloud", strconv.FormatBool(variable.EunUrlCatReviewSubmitToSecurityCloud)),
					resource.TestCheckResourceAttr(resourceName, "url_cat_review_text", "If you believe you received this message in error, please click here to request a review of this site."),
					resource.TestCheckResourceAttr(resourceName, "url_cat_review_custom_location", "https://redirect.acme.com"),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_custom_location", "https://redirect.acme.com"),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_enabled", strconv.FormatBool(variable.EunWebDlpReviewEnabled)),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_submit_to_security_cloud", strconv.FormatBool(variable.EunWebDlpReviewSubmitToSecurityCloud)),
					resource.TestCheckResourceAttr(resourceName, "web_dlp_review_text", "Click to request policy review."),
				),
			},
			// Step 2: Import the resource and verify the state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckResourceEndUserNotificationDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}

// Helper function to generate test configuration for the resource
func testAccResourceEndUserNotificationConfig() string {
	return `
resource "zia_end_user_notification" "test" {
  aup_custom_frequency                = 0
  aup_day_offset                      = 0
  aup_frequency                       = "NEVER"
  caution_again_after                 = 300
  caution_custom_text                 = "Proceeding to visit the site may violate your company policy. Press the \"Continue\" button to access the site anyway or press the \"Back\" button on your browser to go back"
  caution_per_domain                  = true
  custom_text                         = "Website blocked"
  display_company_logo                = true
  display_company_name                = true
  display_reason                      = true
  notification_type                   = "DEFAULT"
  redirect_url                        = "https://redirect.acme.com"
  org_policy_link                     = "http://24326813.zscalerthree.net/policy.html"
  quarantine_custom_notification_text = <<-EOT
We are checking this file for a potential security risk. The file you attempted to download is being analyzed for your protection.
It is not blocked. The analysis can take up to 10 minutes, depending on the size and type of the file. If safe, your file downloads automatically.
If unsafe, the file will be blocked.

EOT
  security_review_custom_location          = "https://redirect.acme.com"
  security_review_enabled                  = true
  security_review_submit_to_security_cloud = true
  security_review_text                     = "Click to request security review."
  support_email                            = "support@24326813.zscalerthree.net"
  support_phone                            = "+91-9000000000"
  url_cat_review_enabled                   = true
  url_cat_review_submit_to_security_cloud  = false
  url_cat_review_text                      = "If you believe you received this message in error, please click here to request a review of this site."
  url_cat_review_custom_location           = "https://redirect.acme.com"
  web_dlp_review_custom_location           = "https://redirect.acme.com"
  web_dlp_review_enabled                   = true
  web_dlp_review_submit_to_security_cloud  = false
  web_dlp_review_text                      = "Click to request policy review."
}
`
}
