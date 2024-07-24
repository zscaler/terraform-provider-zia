---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: cloud_app_control_rule"
description: |-
  Creates and manages ZIA Cloud Application Control rule.
---

# Resource: zia_cloud_app_control_rule

The **zia_cloud_app_control_rule** resource allows the creation and management of ZIA Cloud Application Control rules in the Zscaler Internet Access.

**NOTE** Resources or DataSources to retrieve Tenant Profile or Cloud Application Risk Profile ID information are not currently available.

## Example Usage - Basic Rule Configuration

```hcl
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
    device_trust_levels   = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    user_agent_types      = ["OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"]
}
```

## Example Usage - With Cloud Risk Profile Configuration

```hcl
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
```

## Example Usage - With Tenant Profile Configuration

**NOTE** Tenant profile is supported only for specific applications depending on the type

```hcl
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
```

## Example Usage - With ISOLATE ACTION

⚠️ **WARNING 1:**: Creating a Cloud Application Control Rule with the actions containing `ISOLATE_` Cloud Browser Isolation subscription is required. See the "Cloud Application Control - Rule Types vs Actions Matrix" below. To learn more, contact Zscaler Support or your local account team.

```hcl
data "zia_cloud_browser_isolation_profile" "this" {
    name = "BD_SA_Profile1_ZIA"
}

resource "zia_cloud_app_control_rule" "this" {
    name                         = "Example"
    description                  = "Example"
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
    order                 = 1
    enforce_time_validity = true
    validity_start_time   = "Mon, 17 Jun 2024 23:30:00 UTC"
    validity_end_time     = "Tue, 17 Jun 2025 23:00:00 UTC"
    validity_time_zone_id = "US/Pacific"
    time_quota            = 15
    size_quota            = 10
    device_trust_levels   = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    cbi_profile {
        id = data.zia_cloud_browser_isolation_profile.this.id
        name = data.zia_cloud_browser_isolation_profile.this.name
        url = data.zia_cloud_browser_isolation_profile.this.url
    }
    user_agent_types = [ "OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE" ]
}
```

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Policy access rule can be imported by using `<RULE_TYPE:RULE_ID>` or `<RULE_TYPE:RULE_NAME>` as the import ID.

For example:

```shell
terraform import zia_cloud_app_control_rule.this <rule_type:rule_id>
```

```shell
terraform import zia_cloud_app_control_rule.this <"rule_type:rule_name">
```

## Cloud Application Control - Rule Types vs Actions Matrix

**Note**: Refer to this matrix when configuring types vs actions for each specific rules

|            Types                     |                    Actions                      |
|:------------------------------------:|:-----------------------------------------------:|
|--------------------------------------|-------------------------------------------------|
|           `AI_ML`                    |          `ALLOW_AI_ML_WEB_USE`                  |
|           `AI_ML`                    |          `CAUTION_AI_ML_WEB_USE`                |
|           `AI_ML`                    |          `DENY_AI_ML_WEB_USE`                   |
|           `AI_ML`                    |          `ISOLATE_AI_ML_WEB_USE`                |
|--------------------------------------|-------------------------------------------------|
|     `BUSINESS_PRODUCTIVITY`          |     `ALLOW_BUSINESS_PRODUCTIVITY_APPS`          |
|     `BUSINESS_PRODUCTIVITY`          |     `BLOCK_BUSINESS_PRODUCTIVITY_APPS`          |
|     `BUSINESS_PRODUCTIVITY`          |     `CAUTION_BUSINESS_PRODUCTIVITY_APPS`        |
|     `BUSINESS_PRODUCTIVITY`          |     `ISOLATE_BUSINESS_PRODUCTIVITY_APPS`        |
|--------------------------------------|-------------------------------------------------|
|     `CONSUMER`                       |          `ALLOW_CONSUMER_APPS`                  |
|     `CONSUMER`                       |          `BLOCK_CONSUMER_APPS`                  |
|     `CONSUMER`                       |          `CAUTION_CONSUMER_APPS`                |
|     `CONSUMER`                       |          `ISOLATE_CONSUMER_APPS`                |
|--------------------------------------|-------------------------------------------------|
|     `DNS_OVER_HTTPS`                 |          `ALLOW_DNS_OVER_HTTPS_USE`             |
|     `DNS_OVER_HTTPS`                 |          `DENY_DNS_OVER_HTTPS_USE`              |
|--------------------------------------|-------------------------------------------------|
|     `ENTERPRISE_COLLABORATION`       |      `ALLOW_ENTERPRISE_COLLABORATION_APPS`      |
|     `ENTERPRISE_COLLABORATION`       |      `BLOCK_ENTERPRISE_COLLABORATION_APPS`      |
|     `ENTERPRISE_COLLABORATION`       |      `CAUTION_ENTERPRISE_COLLABORATION_APPS`    |
|     `ENTERPRISE_COLLABORATION`       |      `ISOLATE_ENTERPRISE_COLLABORATION_APPS`    |
|--------------------------------------|-------------------------------------------------|
|     `FILE_SHARE`                     |          `ALLOW_FILE_SHARE_VIEW`                |
|     `FILE_SHARE`                     |          `ALLOW_FILE_SHARE_UPLOAD`              |
|     `FILE_SHARE`                     |          `CAUTION_FILE_SHARE_VIEW`              |
|     `FILE_SHARE`                     |          `DENY_FILE_SHARE_VIEW`                 |
|     `FILE_SHARE`                     |          `DENY_FILE_SHARE_UPLOAD`               |
|     `FILE_SHARE`                     |          `ISOLATE_FILE_SHARE_VIEW`              |
|--------------------------------------|-------------------------------------------------|
|     `FINANCE`                        |          `ALLOW_FINANCE_USE`                    |
|     `FINANCE`                        |          `CAUTION_FINANCE_USE`                  |
|     `FINANCE`                        |          `DENY_FINANCE_USE`                     |
|     `FINANCE`                        |          `ISOLATE_FINANCE_USE`                  |
|--------------------------------------|-------------------------------------------------|
|     `HEALTH_CARE`                    |          `ALLOW_HEALTH_CARE_USE`                |
|     `HEALTH_CARE`                    |          `CAUTION_HEALTH_CARE_USE`              |
|     `HEALTH_CARE`                    |          `DENY_HEALTH_CARE_USE`                 |
|     `HEALTH_CARE`                    |          `ISOLATE_HEALTH_CARE_USE`              |
|--------------------------------------|-------------------------------------------------|
|     `HOSTING_PROVIDER`               |          `ALLOW_HOSTING_PROVIDER_USE`           |
|     `HOSTING_PROVIDER`               |          `CAUTION_HOSTING_PROVIDER_USE`         |
|     `HOSTING_PROVIDER`               |          `DENY_HOSTING_PROVIDER_USE`            |
|     `HOSTING_PROVIDER`               |          `ISOLATE_HOSTING_PROVIDER_USE`         |
|--------------------------------------|-------------------------------------------------|
|     `HUMAN_RESOURCES`                |          `ALLOW_HUMAN_RESOURCES_USE`            |
|     `HUMAN_RESOURCES`                |          `CAUTION_HUMAN_RESOURCES_USE`          |
|     `HUMAN_RESOURCES`                |          `DENY_HUMAN_RESOURCES_USE`             |
|     `HUMAN_RESOURCES`                |          `ISOLATE_HUMAN_RESOURCES_USE`          |
|--------------------------------------|-------------------------------------------------|
|     `INSTANT_MESSAGING`              |          `ALLOW_CHAT`                           |
|     `INSTANT_MESSAGING`              |          `ALLOW_FILE_TRANSFER_IN_CHAT`          |
|     `INSTANT_MESSAGING`              |          `BLOCK_CHAT`                           |
|     `INSTANT_MESSAGING`              |          `BLOCK_FILE_TRANSFER_IN_CHAT`          |
|     `INSTANT_MESSAGING`              |          `CAUTION_CHAT`                         |
|     `INSTANT_MESSAGING`              |          `ISOLATE_CHAT`                         |
|--------------------------------------|-------------------------------------------------|
|     `IT_SERVICES`                    |          `ALLOW_IT_SERVICES_USE`                |
|     `IT_SERVICES`                    |          `CAUTION_LEGAL_USE`                    |
|     `IT_SERVICES`                    |          `DENY_IT_SERVICES_USE`                 |
|     `IT_SERVICES`                    |          `ISOLATE_IT_SERVICES_USE`              |
|--------------------------------------|-------------------------------------------------|
|     `LEGAL`                          |          `ALLOW_LEGAL_USE`                      |
|     `LEGAL`                          |          `DENY_DNS_OVER_HTTPS_USE`              |
|     `LEGAL`                          |          `DENY_LEGAL_USE`                       |
|     `LEGAL`                          |          `ISOLATE_LEGAL_USE`                    |
|--------------------------------------|-------------------------------------------------|
|     `SALES_AND_MARKETING`            |          `ALLOW_SALES_MARKETING_APPS`           |
|     `SALES_AND_MARKETING`            |          `BLOCK_SALES_MARKETING_APPS`           |
|     `SALES_AND_MARKETING`            |          `CAUTION_SALES_MARKETING_APPS`         |
|     `SALES_AND_MARKETING`            |          `ISOLATE_SALES_MARKETING_APPS`         |
|--------------------------------------|-------------------------------------------------|
|     `STREAMING_MEDIA`                |          `ALLOW_STREAMING_VIEW_LISTEN`          |
|     `STREAMING_MEDIA`                |          `ALLOW_STREAMING_UPLOAD`               |
|     `STREAMING_MEDIA`                |          `BLOCK_STREAMING_UPLOAD`               |
|     `STREAMING_MEDIA`                |          `CAUTION_STREAMING_VIEW_LISTEN`        |
|     `STREAMING_MEDIA`                |          `ISOLATE_STREAMING_VIEW_LISTEN`        |
|--------------------------------------|-------------------------------------------------|
|     `SOCIAL_NETWORKING`              |          `ALLOW_SOCIAL_NETWORKING_VIEW`         |
|     `SOCIAL_NETWORKING`              |          `ALLOW_SOCIAL_NETWORKING_POST`         |
|     `SOCIAL_NETWORKING`              |          `BLOCK_SOCIAL_NETWORKING_VIEW`         |
|     `SOCIAL_NETWORKING`              |          `BLOCK_SOCIAL_NETWORKING_POST`         |
|     `SOCIAL_NETWORKING`              |          `CAUTION_SOCIAL_NETWORKING_VIEW`       |
|--------------------------------------|-------------------------------------------------|
|     `SYSTEM_AND_DEVELOPMENT`         |          `ALLOW_SYSTEM_DEVELOPMENT_APPS`        |
|     `SYSTEM_AND_DEVELOPMENT`         |          `ALLOW_SYSTEM_DEVELOPMENT_UPLOAD`      |
|     `SYSTEM_AND_DEVELOPMENT`         |          `BLOCK_SYSTEM_DEVELOPMENT_APPS`        |
|     `SYSTEM_AND_DEVELOPMENT`         |          `BLOCK_SYSTEM_DEVELOPMENT_UPLOAD`      |
|     `SYSTEM_AND_DEVELOPMENT`         |          `CAUTION_SYSTEM_DEVELOPMENT_APPS`      |
|     `SYSTEM_AND_DEVELOPMENT`         |          `ISOLATE_SALES_MARKETING_APPS`         |
|--------------------------------------|-------------------------------------------------|
|     `WEBMAIL`                        |          `ALLOW_WEBMAIL_VIEW`                   |
|     `WEBMAIL`                        |          `ALLOW_WEBMAIL_ATTACHMENT_SEND`        |
|     `WEBMAIL`                        |          `ALLOW_WEBMAIL_SEND`                   |
|     `WEBMAIL`                        |          `CAUTION_WEBMAIL_VIEW`                 |
|     `WEBMAIL`                        |          `BLOCK_WEBMAIL_VIEW`                   |
|     `WEBMAIL`                        |          `BLOCK_WEBMAIL_ATTACHMENT_SEND`        |
|     `WEBMAIL`                        |          `BLOCK_WEBMAIL_SEND`                   |
|     `WEBMAIL`                        |          `ISOLATE_WEBMAIL_VIEW`                   |
|--------------------------------------|-------------------------------------------------|

## Cloud Application Control - Rule Types vs Tenant Profile Support

**Note**: Refer to this matrix when configuring a Cloud App Control rule with Tenant Profile

[Reference](https://help.zscaler.com/zia/documentation-knowledgebase/policies/cloud-apps/cloud-app-control-policies)

|               Type               |         Applications          | tenancy_profile_ids |
|:--------------------------------:|:-----------------------------:|:-------------------:|
|----------------------------------|-------------------------------|---------------------|
| `BUSINESS_PRODUCTIVITY`          | `"GOOGLEANALYTICS"`           |          ✅         |
|----------------------------------|-------------------------------|---------------------|
| `ENTERPRISE_COLLABORATION`       | `"GOOGLECALENDAR"`            |          ✅         |
| `ENTERPRISE_COLLABORATION`       | `"GOOGLEKEEP"`                |          ✅         |
| `ENTERPRISE_COLLABORATION`       | `"GOOGLEMEET"`                |          ✅         |
| `ENTERPRISE_COLLABORATION`       | `"GOOGLESITES"`               |          ✅         |
| `ENTERPRISE_COLLABORATION`       | `"WEBEX"`                     |          ✅         |
| `ENTERPRISE_COLLABORATION`       | `"SLACK"`                     |          ✅         |
| `ENTERPRISE_COLLABORATION`       | `"WEBEX_TEAMS"`               |          ✅         |
| `ENTERPRISE_COLLABORATION`       | `"ZOOM"`                      |          ✅         |
|----------------------------------|-------------------------------|---------------------|
| `FILE_SHARE`                     | `"DROPBOX"`                   |          ✅         |
| `FILE_SHARE`                     | `"GDRIVE"`                    |          ✅         |
| `FILE_SHARE`                     | `"GPHOTOS"`                   |          ✅         |
|----------------------------------|-------------------------------|---------------------|
| `HOSTING_PROVIDER`               | `"GCLOUDCOMPUTE"`             |          ✅         |
| `HOSTING_PROVIDER`               | `"AWS"`                       |          ✅         |
| `HOSTING_PROVIDER`               | `"IBMSMARTCLOUD"`             |          ✅         |
| `HOSTING_PROVIDER`               | `"GAPPENGINE"`                |          ✅         |
| `HOSTING_PROVIDER`               | `"GOOGLE_CLOUD_PLATFORM"`     |          ✅         |
|----------------------------------|-------------------------------|---------------------|
| `IT_SERVICES`                    | `"MSLOGINSERVICES"`           |          ✅         |
| `IT_SERVICES`                    | `"GOOGLOGINSERVICE"`          |          ✅         |
| `IT_SERVICES`                    | `"WEBEX_LOGIN_SERVICES"`      |          ✅         |
| `IT_SERVICES`                    | `"ZOHO_LOGIN_SERVICES"`       |          ✅         |
|----------------------------------|-------------------------------|---------------------|
| `SOCIAL_NETWORKING`              | `"GOOGLE_GROUPS"`             |          ✅         |
| `SOCIAL_NETWORKING`              | `"GOOGLE_PLUS"`               |          ✅         |
|----------------------------------|-------------------------------|---------------------|
| `STREAMING_MEDIA`                | `"YOUTUBE"`                   |          ✅         |
| `STREAMING_MEDIA`                | `"GOOGLE_STREAMING"`          |          ✅         |
|----------------------------------|-------------------------------|---------------------|
| `SYSTEM_AND_DEVELOPMENT`         | `"GOOGLE_DEVELOPERS"`         |          ✅         |
| `SYSTEM_AND_DEVELOPMENT`         | `"GOOGLEAPPMAKER"`            |          ✅         |
|----------------------------------|-------------------------------|---------------------|
| `WEBMAIL`                        | `"GOOGLE_WEBMAIL"`            |          ✅         |
|----------------------------------|-------------------------------|---------------------|
