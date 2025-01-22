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
