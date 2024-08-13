## Example Usage - Basic Rule Configuration
resource "zia_cloud_app_control_rule" "this" {
    name                         = "Example_WebMail_Rule"
    description                  = "Example_WebMail_Rule"
    order                        = 1
    rank                         = 7
    state                        = "ENABLED"
    type                         = "WEBMAIL"
    actions                      = [
        "ALLOW_WEBMAIL_VIEW",
        "ALLOW_WEBMAIL_ATTACHMENT_SEND",
        "ALLOW_WEBMAIL_SEND",
    ]
    applications          = [
        "GOOGLE_WEBMAIL",
        "YAHOO_WEBMAIL",
    ]
}

## Example Usage - With Cloud Risk Profile Configuration
resource "zia_cloud_app_control_rule" "this" {
    name                         = "Example_WebMail_Rule"
    description                  = "Example_WebMail_Rule"
    order                        = 1
    rank                         = 7
    state                        = "ENABLED"
    type                         = "WEBMAIL"
    actions                      = [
        "ALLOW_WEBMAIL_VIEW",
        "ALLOW_WEBMAIL_ATTACHMENT_SEND",
        "ALLOW_WEBMAIL_SEND",
    ]
    applications          = [
        "GOOGLE_WEBMAIL",
        "YAHOO_WEBMAIL",
    ]
    cloud_app_risk_profile {
      id = 318
    }
}


## Example Usage - With Tenant Profile Configuration
# NOTE Tenant profile is supported only for specific applications depending on the type
resource "zia_cloud_app_control_rule" "this" {
    name                         = "Example_WebMail_Rule"
    description                  = "Example_WebMail_Rule"
    order                        = 1
    rank                         = 7
    state                        = "ENABLED"
    type                         = "WEBMAIL"
    actions                      = [
        "ALLOW_WEBMAIL_VIEW",
        "ALLOW_WEBMAIL_ATTACHMENT_SEND",
        "ALLOW_WEBMAIL_SEND",
    ]
    applications          = [
        "GOOGLE_WEBMAIL",
        "YAHOO_WEBMAIL",
    ]
    tenancy_profile_ids {
        id = [ 19016237 ]
    }
}

