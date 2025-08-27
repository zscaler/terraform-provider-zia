---
subcategory: "End User Notification"
layout: "zscaler"
page_title: "ZIA: end_user_notification"
description: |-
  Official documentation https://help.zscaler.com/zia/understanding-browser-based-end-user-notifications
  API documentation https://help.zscaler.com/zia/end-user-notifications#/eun-get
  Updates the advanced threat configuration settings.
---

# zia_end_user_notification (Resource)

* [Official documentation](https://help.zscaler.com/zia/understanding-browser-based-end-user-notifications)
* [API documentation](https://help.zscaler.com/zia/end-user-notifications#/eun-get)

The **zia_end_user_notification** resource allows you to update the browser-based end user notification (EUN) configuration details. To learn more see [Understanding Browser-Based End User Notifications](https://help.zscaler.com/unified/understanding-browser-based-end-user-notifications)

## Example Usage - NOTIFICATION TYPE - DEFAULT

```hcl
# END USER NOTIFICATION TYPE - DEFAULT
resource "zia_end_user_notification" "this" {
  aup_frequency     = "NEVER"
  aup_message       = "Please review and accept the terms."
  notification_type = "DEFAULT"
  custom_text       = "Website blocked"

  url_cat_review_enabled                  = true
  url_cat_review_submit_to_security_cloud = true
  url_cat_review_text                     = "Click here to submit a review request."

  security_review_enabled                  = true
  security_review_submit_to_security_cloud = true
  security_review_text                     = "Request a security review if this message appears in error."

  web_dlp_review_enabled                  = true
  web_dlp_review_submit_to_security_cloud = false
  web_dlp_review_custom_location          = "https://dlp-review-location.com"
  web_dlp_review_text                     = "This file is being reviewed for security reasons."

  redirect_url    = "https://dlp-review-location.com"
  support_email   = "support@8061240.zscalerbeta.net"
  support_phone   = "+91-9000000000"
  org_policy_link = "http://8061240.zscalerbeta.net/policy.html"

  caution_again_after = 300
  caution_per_domain  = true
  caution_custom_text = "Access to this site is restricted. Proceed with caution."

  idp_proxy_notification_text         = "Your connection is being proxied through the organization's secure access service."
  quarantine_custom_notification_text = "The file is being analyzed for potential security risks. Please wait while the process completes."
}
```

## Example Usage - NOTIFICATION TYPE - CUSTOM

```hcl
# END USER NOTIFICATION TYPE - CUSTOM
resource "zia_end_user_notification" "this" {
  aup_frequency     = "ON_WEEKDAY"
  aup_day_offset    = "1"
  aup_message       = "Please review and accept the terms."
  notification_type = "CUSTOM"
  display_reason    = true
  display_comp_name = true
  display_comp_logo = true
  custom_text       = "Website blocked"

  url_cat_review_enabled                  = true
  url_cat_review_submit_to_security_cloud = true
  url_cat_review_custom_location          = "https://custom-review-location.com"
  url_cat_review_text                     = "Click here to submit a review request."

  security_review_enabled                  = true
  security_review_submit_to_security_cloud = true
  security_review_custom_location          = "https://security-review-location.com"
  security_review_text                     = "Request a security review if this message appears in error."

  web_dlp_review_enabled                  = true
  web_dlp_review_submit_to_security_cloud = false
  web_dlp_review_custom_location          = "https://dlp-review-location.com"
  web_dlp_review_text                     = "This file is being reviewed for security reasons."

  redirect_url    = "https://dlp-review-location.com"
  support_email   = "support@8061240.zscalerbeta.net"
  support_phone   = "+91-9000000000"
  org_policy_link = "http://8061240.zscalerbeta.net/policy.html"

  caution_again_after = 300
  caution_per_domain  = true
  caution_custom_text = "Access to this site is restricted. Proceed with caution."

  idp_proxy_notification_text         = "Your connection is being proxied through the organization's secure access service."
  quarantine_custom_notification_text = "The file is being analyzed for potential security risks. Please wait while the process completes."
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `aup_frequency` (String) - The frequency at which the Acceptable Use Policy (AUP) is shown to end users. Supported values: NEVER, SESSION, ONLOGIN, CUSTOM, DAILY, WEEKLY, ON_DATE, ON_WEEKDAY.
* `aup_custom_frequency` (Integer) - The custom frequency (in days) for showing the AUP. Valid range: 1-180.
* `aup_day_offset` (Integer) - Specifies which day of the week or month the AUP is shown. Valid range: 1-31.
* `aup_message` (String) - The acceptable use statement that appears in the AUP.
* `notification_type` (String) - The type of EUN, either DEFAULT or CUSTOM.
* `display_reason` (Boolean) - Indicates whether the reason for blocking access is displayed in the EUN.
* `display_comp_name` (Boolean) - Indicates whether the organization's name is displayed in the EUN.
* `display_comp_logo` (Boolean) - Indicates whether the organization's logo is displayed in the EUN.
* `custom_text` (String) - Custom text displayed in the EUN.
* `url_cat_review_enabled` (Boolean) - Indicates whether URL Categorization notifications are enabled.
* `url_cat_review_submit_to_security_cloud` (Boolean) - Indicates whether review requests are submitted to Zscaler Security Cloud.
* `url_cat_review_custom_location` (String) - Custom URL location for review requests of blocked URLs.
* `url_cat_review_text` (String) - The message displayed in the URL Categorization notification.
* `security_review_enabled` (Boolean) - Indicates whether Security Violation notifications are enabled.
* `security_review_submit_to_security_cloud` (Boolean) - Indicates whether review requests are submitted to Zscaler Security Cloud.
* `security_review_custom_location` (String) - Custom URL for review requests of blocked URLs.
* `security_review_text` (String) - The message displayed in the Security Violation notification.
* `web_dlp_review_enabled` (Boolean) - Indicates whether Web DLP Violation notifications are enabled.
* `web_dlp_review_submit_to_security_cloud` (Boolean) - Indicates whether review requests for DLP policy violations are sent to Zscaler Security Cloud.
* `web_dlp_review_custom_location` (String) - Custom URL for review requests related to DLP violations.
* `web_dlp_review_text` (String) - The message displayed in the Web DLP Violation notification.
* `redirect_url` (String) - Redirect URL when the custom notification type is selected.
* `support_email` (String) - IT Support email address.
* `support_phone` (String) - IT Support phone number.
* `org_policy_link` (String) - URL of the organization's policy page (required for default notification type).
* `caution_again_after` (Integer) - Time interval (in minutes) for showing the caution notification when browsing restricted sites. Minimum: 5.
* `caution_per_domain` (Boolean) - Specifies whether the caution notification is displayed per domain for URLs in the Miscellaneous or Unknown category.
* `caution_custom_text` (String) - Custom message displayed in the caution notification.
* `idp_proxy_notification_text` (String) - Message displayed in the IdP Proxy notification.
* `quarantine_custom_notification_text` (String) - Message displayed in the quarantine notification.

## Important Notes

**Text Attributes and CSS Styling**: When setting attributes such as `custom_text`, `url_cat_review_text`, `security_review_text`, `web_dlp_review_text`, `caution_custom_text`, `idp_proxy_notification_text`, and `quarantine_custom_notification_text`, we recommend using heredocs (EOT) especially when including CSS stylesheets. This ensures proper formatting and readability of complex text content.

**JavaScript Limitation**: The ZIA API currently does not accept JavaScript tags in notification text attributes. Using JavaScript tags will result in an HTTP 406 Rejected error. For more information on customizing EUN CSS styles, see the [Zscaler documentation](https://help.zscaler.com/zia/customizing-euns-css-styles).

### Example with Heredoc for CSS Styling

```hcl
resource "zia_end_user_notification" "this" {
  notification_type = "CUSTOM"
  custom_text = <<EOT
    <div style="background-color: #f0f0f0; padding: 20px; border-radius: 5px;">
      <h2 style="color: #333;">Access Blocked</h2>
      <p style="color: #666;">This website has been blocked for security reasons.</p>
    </div>
  EOT
  
  url_cat_review_text = <<EOT
    <div style="text-align: center; margin: 10px 0;">
      <button style="background-color: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 3px;">
        Request Review
      </button>
    </div>
  EOT
}
```

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_end_user_notification** can be imported by using `enduser_notification` as the import ID.

For example:

```shell
terraform import zia_end_user_notification.this "enduser_notification"
```
