---
subcategory: "End User Notification"
layout: "zscaler"
page_title: "ZIA: end_user_notification"
description: |-
  Official documentation https://help.zscaler.com/zia/understanding-browser-based-end-user-notifications
  API documentation https://help.zscaler.com/zia/end-user-notifications#/eun-get
  Retrieves browser-based end user notification (EUN) configuration details
---

# zia_end_user_notification (Data Source)

* [Official documentation](https://help.zscaler.com/zia/understanding-browser-based-end-user-notifications)
* [API documentation](https://help.zscaler.com/zia/end-user-notifications#/eun-get)

Use the **zia_end_user_notification** data source to get information about browser-based end user notification (EUN) configuration details.

## Example Usage

```hcl
data "zia_end_user_notification" "example"{}
```

## Argument Reference

### Read-Only

* `aup_frequency` (String) - The frequency at which the Acceptable Use Policy (AUP) is shown to end users. Supported values: NEVER, SESSION, ONLOGIN, CUSTOM, DAILY, WEEKLY, ON_DATE, ON_WEEKDAY.
* `aup_custom_frequency` (Integer) - The custom frequency (in days) for showing the AUP. Valid range: 1-180.
* `aup_day_offset` (Integer) - Specifies which day of the week or month the AUP is shown. Valid range: 1-31.
* `aup_message` (String) - The acceptable use statement that appears in the AUP.
* `notification_type` (String) - The type of EUN, either DEFAULT or CUSTOM.
* `display_reason` (Boolean) - Indicates whether the reason for blocking access is displayed in the EUN.
* `display_company_name` (Boolean) - Indicates whether the organization's name is displayed in the EUN.
* `display_company_logo` (Boolean) - Indicates whether the organization's logo is displayed in the EUN.
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
