---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_web_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-dlp-policy-rules-content-inspection#Rules
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-get
  Creates and manages ZIA DLP Web Rules.
---

# zia_dlp_web_rules (Resource)

* [Official documentation](https://help.zscaler.com/zia/configuring-dlp-policy-rules-content-inspection#Rules)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-get)

The **zia_dlp_web_rules** resource allows the creation and management of ZIA DLP Web Rules in the Zscaler Internet Access cloud or via the API.

⚠️ **WARNING:** Zscaler Internet Access DLP supports a maximum of 127 Web DLP Rules to be created via API.

~> **NOTE:** Predefined rules can be managed via the Terraform provider for reordering purposes; however, `destroy` operations are not supported for predefined rules, and not all attributes available on custom rules apply to them. When deleting existing custom rules, use the Terraform `-target` flag to target the specific rule to be removed.

~> **NOTE:** Rule orders must always be contiguous (no gaps). Deleting a rule must be followed by order number re-adjustment of the remaining rules to ensure the API honours the required order.

~> **NOTE:** The `order` attribute must always be a positive whole number starting at 1. Negative numbers and zero are **not supported** and will result in an error.

## Example Usage - "FTCATEGORY_ALL_OUTBOUND" File Type"

```hcl
data "zia_dlp_engines" "this" {
  predefined_engine_name = "EXTERNAL"
}

resource "zia_dlp_web_rules" "this" {
  name                       = "Example"
  description                = "Example"
  action                     = "BLOCK"
  order                      = 1
  rank                       = 7
  state                      = "ENABLED"
  protocols                  = [ "FTP_RULE", "HTTPS_RULE", "HTTP_RULE" ]
  file_types                 = [ "FTCATEGORY_ALL_OUTBOUND" ]
  zscaler_incident_receiver  = false
  without_content_inspection = true
  user_risk_score_levels     = [ "LOW", "MEDIUM", "HIGH", "CRITICAL" ]
  severity                   = "RULE_SEVERITY_HIGH"
  dlp_engines {
    id = [ data.zia_dlp_engines.this.id ]
  }
}
```

```hcl
// Example 1: Using data source to reference existing URL category
data "zia_url_categories" "existing_category" {
    configured_name = "Example"
}

// Example 2: Creating new URL category and referencing it
resource "zia_url_categories" "new_category" {
  configured_name = "Custom_Category"
  description     = "Custom category for DLP rules"
  custom_category = true
  super_category  = "USER_DEFINED"
  type            = "URL_CATEGORY"
}

// Retrieve an ICAP Server by Name
data "zia_dlp_icap_servers" "this" {
  name = "ZS_ICAP_01"
}

resource "zia_dlp_web_rules" "this" {
  name                      = "Terraform_Test"
  description               = "Terraform_Test"
  action                    = "BLOCK"
  order                     = 1
  protocols                 = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
  rank                      = 7
  state                     = "ENABLED"
  zscaler_incident_receiver = true
  without_content_inspection = false
  url_categories {
    id = [ data.zia_url_categories.existing_category.val ]
  }
  icap_server {
    id = data.zia_dlp_icap_servers.this.id
  }
}

resource "zia_dlp_web_rules" "with_new_category" {
  name                      = "Terraform_Test_New_Category"
  description               = "Terraform_Test with new category"
  action                    = "BLOCK"
  order                     = 2
  protocols                 = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
  rank                      = 7
  state                     = "ENABLED"
  zscaler_incident_receiver = true
  without_content_inspection = false
  url_categories {
    id = [ zia_url_categories.new_category.val ]
  }
  icap_server {
    id = data.zia_dlp_icap_servers.this.id
  }
}
```

## Example Usage - "FTCATEGORY_ALL_OUTBOUND" File Type - New"

```hcl
data "zia_dlp_engines" "this" {
  predefined_engine_name = "EXTERNAL"
}

data "zia_file_type_categories" "this" {
    name = "FileType01"
}

resource "zia_dlp_web_rules" "this" {
  name                       = "Example"
  description                = "Example"
  action                     = "BLOCK"
  order                      = 1
  rank                       = 7
  state                      = "ENABLED"
  protocols                  = [ "FTP_RULE", "HTTPS_RULE", "HTTP_RULE" ]
  zscaler_incident_receiver  = false
  without_content_inspection = true
  user_risk_score_levels     = [ "LOW", "MEDIUM", "HIGH", "CRITICAL" ]
  severity                   = "RULE_SEVERITY_HIGH"
  file_type_categories {
    id = [ data.zia_file_type_categories.this.id ]
  }
  dlp_engines {
    id = [ data.zia_dlp_engines.this.id ]
  }
}
```

## Example Usage - "Specify Incident Receiver Setting"

```hcl
// Retrieve a custom URL Category by Name
data "zia_url_categories" "this"{
    configured_name = "Example"
}

// Retrieve a Incident Receiver by Name
data "zia_dlp_incident_receiver_servers" "this" {
  name = "ZS_INC_RECEIVER_01"
}

resource "zia_dlp_web_rules" "this" {
  name                      = "Terraform_Test"
  description               = "Terraform_Test"
  action                    = "BLOCK"
  order                     = 1
  protocols                 = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
  rank                      = 7
  state                     = "ENABLED"
  zscaler_incident_receiver = true
  without_content_inspection = false
  url_categories {
    id = [ data.zia_url_categories.this.val ]
  }
  icap_server {
    id = data.zia_dlp_incident_receiver_servers.this.id
  }
  notification_template {
    id = data.zia_dlp_notification_templates.this.id
  }
}
```

## Example Usage - "Creating Parent Rules and SubRules"

⚠️ **WARNING:** Destroying a parent rule will also destroy all sub-rules.

~> **NOTE:** Exception (sub-) rules are **separate** `zia_dlp_web_rules` resources. Set `parent_rule` to the parent’s numeric rule ID (`rule_id`). The parent’s `sub_rules` attribute is **computed** after apply—it lists child rule IDs returned by the API; you do not author nested rule blocks inside the parent resource.

 **NOTE** Exception rules can be configured only when the inline DLP rule evaluation type is set
 to evaluate all DLP rules in the DLP Advanced Settings.
 To learn more, see [Configuring DLP Advanced Settings](https://help.zscaler.com/%22/zia/configuring-dlp-advanced-settings/%22)

```hcl
resource "zia_dlp_web_rules" "parent_rule" {
  name                       = "ParentRule1"
  description                = "ParentRule1"
  action                     = "ALLOW"
  state                      = "ENABLED"
  order                      = 1
  rank                       = 0
  protocols                  = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
  cloud_applications         = ["GOOGLE_WEBMAIL", "WINDOWS_LIVE_HOTMAIL"]
  without_content_inspection = false
  match_only                 = false
  min_size                   = 20
  zscaler_incident_receiver  = true
}

resource "zia_dlp_web_rules" "subrule1" {
  name                       = "SubRule1"
  description                = "SubRule1"
  action                     = "ALLOW"
  state                      = "ENABLED"
  order                      = 1
  rank                       = 0
  protocols                  = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
  cloud_applications         = ["GOOGLE_WEBMAIL", "WINDOWS_LIVE_HOTMAIL"]
  without_content_inspection = false
  match_only                 = false
  parent_rule                = zia_dlp_web_rules.parent_rule.rule_id
}
```

## Example Usage - "Configuring Receiver for DLP Policy Rule"

```hcl
resource "zia_dlp_web_rules" "with_receiver" {
  name                       = "Terraform_Test_with_Receiver"
  description                = "DLP rule with receiver configuration"
  action                     = "ALLOW"
  state                      = "ENABLED"
  order                      = 1
  rank                       = 0
  protocols                  = [
    "WEBSOCKETSSL_RULE",
    "WEBSOCKET_RULE",
    "FTP_RULE",
    "HTTPS_RULE",
    "HTTP_RULE"
  ]
  severity = "RULE_SEVERITY_HIGH"

  # Basic receiver configuration with just ID
  receiver {
    id = "23136553"
  }
}
```

## Example Usage - Configure Cloud to Cloud Forwarding

```hcl
# Retrieve Cloud-to-Cloud Incident Receiver (C2CIR) information
data "zia_dlp_cloud_to_cloud_ir" "this" {
  name = "AzureTenant01"
}

# Output the retrieved C2CIR information for reference
output "zia_dlp_cloud_to_cloud_ir" {
  value = data.zia_dlp_cloud_to_cloud_ir.this
}

resource "zia_dlp_web_rules" "this" {
  name                       = "Terraform_Test_policy_prod_tf"
  description                = "Terraform_Test_policy_prod_tf"
  action                     = "ALLOW"
  state                      = "ENABLED"
  order                      = 1
  rank                       = 0
  protocols                  = [
        "WEBSOCKETSSL_RULE",
        "WEBSOCKET_RULE",
        "FTP_RULE",
        "HTTPS_RULE",
        "HTTP_RULE"
    ]
  severity = "RULE_SEVERITY_HIGH"

  # Configure receiver using values from the C2CIR data source
  receiver {
    id   = tostring(data.zia_dlp_cloud_to_cloud_ir.this.onboardable_entity[0].tenant_authorization_info[0].smir_bucket_config[0].id)
    name = data.zia_dlp_cloud_to_cloud_ir.this.onboardable_entity[0].tenant_authorization_info[0].smir_bucket_config[0].config_name
    type = data.zia_dlp_cloud_to_cloud_ir.this.onboardable_entity[0].type
    tenant {
      id   = tostring(data.zia_dlp_cloud_to_cloud_ir.this.id)
      name = data.zia_dlp_cloud_to_cloud_ir.this.name
    }
  }
}
```

**Note:** The receiver configuration uses values from the C2CIR data source:

* `id`: Uses the SMIR bucket configuration ID (converted to string)
* `name`: Uses the SMIR bucket configuration name
* `type`: Uses the onboardable entity type (e.g., "C2CIR")
* `tenant.id`: Uses the C2CIR tenant ID (converted to string)
* `tenant.name`: Uses the C2CIR tenant name

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP policy rule name.
* `order` - (Required) The rule order of execution for the DLP policy rule with respect to other rules.

### Optional

* `description` - (Optional) The description of the DLP policy rule.
* `external_auditor_email` - (Optional) The email address of an external auditor to whom DLP email notifications are sent.
* `match_only` - (Optional) The match only criteria for DLP engines.
* `without_content_inspection` - (Optional) Indicates a DLP policy rule without content inspection, when the value is set to true.
  * `without_content_inspection` must be set to false if `file_types` is not defined.

* `ocr_enabled` - (Optional) Enables or disables image file scanning. When OCR is enabled only the following ``file_types`` are supported: ``WINDOWS_META_FORMAT``, ``BITMAP``, ``JPEG``, ``PNG``, ``TIFF``
* `zscaler_incident_receiver` - (Optional) Indicates whether a Zscaler Incident Receiver is associated to the DLP policy rule.

* `action` - (Optional) The action taken when traffic matches the DLP policy rule criteria. The supported values are:
  * `ANY`
  * `NONE`
  * `BLOCK`
  * `CONFIRM`
  * `ALLOW`
  * `ICAP_RESPONSE`

* `state` - (Optional) Enables or disables the DLP policy rule.. The supported values are:
  * `DISABLED`
  * `ENABLED`

* `file_types` - (Optional) The list of file types to which the DLP policy rule must be applied. For the complete list of supported file types refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post)

  * ~> Note: `BITMAP`, `JPEG`, `PNG`, and `TIFF` file types are exclusively supported when optical character recognition `ocr_enabled` is set to `true` for DLP rules with content inspection.

  * ~> Note: `FTCATEGORY_ALL_OUTBOUND` file type is applicable only when the predefined DLP engine called `EXTERNAL` is used and when the attribute `without_content_inspection` is set to `false`.

  * ~> Note: `FTCATEGORY_ALL_OUTBOUND` file type cannot be used alongside any other file type.

* `cloud_applications` - (Optional) The list of cloud applications to which the DLP policy rule must be applied. For the complete list of supported file types refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post). To retrieve the list of cloud applications, use the data source: `zia_cloud_applications`

* `severity` - (Optional) Indicates the severity selected for the DLP rule violation: Returned values are:  `RULE_SEVERITY_HIGH`, `RULE_SEVERITY_MEDIUM`, `RULE_SEVERITY_LOW`, `RULE_SEVERITY_INFO`

* `user_risk_score_levels` (Optional) - Indicates the user risk score level selected for the DLP rule violation: Returned values are: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

* `parent_rule` (Optional) - The unique identifier of the parent rule under which an exception (sub-) rule is added. Use the parent resource’s `rule_id` (integer). The rule rank must be set to `0`.

    ~> **Note**: Exception rules can be configured only when the inline DLP rule evaluation type is set to evaluate all DLP rules in the DLP Advanced Settings. To learn more, see [Configuring DLP Advanced Settings](https://help.zscaler.com/%22/zia/configuring-dlp-advanced-settings/%22)

    ~> **Note**: It is not possible to add existing rules as sub-rules under the parent rule.

* `sub_rules` (Optional, Computed) - Set of sub-rule IDs (strings), populated from the API for a **parent** rule after read. Sub-rules are managed as their own `zia_dlp_web_rules` resources with `parent_rule` set; do not model sub-rules as nested blocks in the parent. When sending updates, the provider may pass sub-rule references as ID-only entries as required by the API. The rule rank must be set to `0` where applicable.

    ~> **Note**: All attributes within the Web DLP rule model apply to sub-rules. Exception rules can be configured only when the inline DLP rule evaluation type is set to evaluate all DLP rules in the DLP Advanced Settings. To learn more, see [Configuring DLP Advanced Settings](https://help.zscaler.com/%22/zia/configuring-dlp-advanced-settings/%22)

* `notification_template` - (Optional) The template used for DLP notification emails.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `auditor` - (Optional) The auditor to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `url_categories` - (Optional) The list of URL categories to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity
  ~> **NOTE** When associating a URL category, you can use the `zia_url_categories` resource or data source; however, you must export the attribute `val`. The `val` attribute is available on both the resource and data source, making it consistent for referencing URL categories in DLP web rules.

* `dlp_engines` - (Optional) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity. Maximum of up to `4` dlp engines. When not used it implies `Any` to apply the rule to all locations.

* `locations` - (Optional) The Name-ID pairs of locations to which the DLP policy rule must be applied. Maximum of up to `32` locations. When not used it implies `Any` to apply the rule to all locations.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `location_groups` - (Optional) The Name-ID pairs of locations groups to which the DLP policy rule must be applied. Maximum of up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `users` - (Optional) The Name-ID pairs of users to which the DLP policy rule must be applied. Maximum of up to `4` users. When not used it implies `Any` to apply the rule to all users.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `groups` - (Optional) The Name-ID pairs of groups to which the DLP policy rule must be applied. Maximum of up to `32` groups. When not used it implies `Any` to apply the rule to all groups.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `departments` - (Optional) The name-ID pairs of the departments that are excluded from the DLP policy rule.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `excluded_users` - (Optional) The name-ID pairs of the users that are excluded from the DLP policy rule. Maximum of up to `256` users.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `excluded_departments` - (Optional) The name-ID pairs of the groups that are excluded from the DLP policy rule. Maximum of up to `256` departments.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `excluded_groups` - (Optional) The name-ID pairs of the groups that are excluded from the DLP policy rule. Maximum of up to `256` groups.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `time_windows` - (Optional) The Name-ID pairs of time windows to which the DLP policy rule must be applied. Maximum of up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `labels` - (List of Object) The Name-ID pairs of rule labels associated to the DLP policy rule.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `source_ip_groups` - (List of Object) The source ip groups to which the DLP policy rule applies
  * `id` - (Optional) Source IP address groups for which the rule is applicable.

* `icap_server` The DLP server, using ICAP, to which the transaction content is forwarded.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `receiver` - (Optional) The receiver information for the DLP policy rule.
  * `id` - (Required) Unique identifier for the receiver
  * `name` - (Optional) Name of the receiver
  * `type` - (Optional) Type of the receiver
  * `tenant` - (Optional) Tenant information for the receiver
    * `id` - (Optional) Unique identifier for the tenant
    * `name` - (Optional) Name of the tenant

* `workload_groups` (Optional) The list of preconfigured workload groups to which the policy must be applied
  * `id` - (Optional) A unique identifier assigned to the workload group
  * `name` - (Optional) The name of the workload group

* `file_type_categories` - (Optional) File type categories to which the rule applies (IDs from `zia_file_type_categories`). Zscaler recommends this over legacy `file_types` where possible. The API allows either `fileTypes` or `fileTypeCategories`, but **not both** in the same request; use only one of `file_types` or `file_type_categories` in configuration.
  * `id` - (Optional) File type category ID.
    **NOTE** Use the data source `zia_file_type_categories` to retrieve file type categories.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_dlp_web_rules** can be imported using either:

* `<RULE ID>` — numeric ID of the rule (parent or exception/sub-rule), or
* `<RULE NAME>` — exact rule name; resolution includes **exception rules** nested under parents in the API (not only top-level rule names).

For example:

```shell
terraform import zia_dlp_web_rules.example <rule_id>
```

or

```shell
terraform import zia_dlp_web_rules.example <rule_name>
```

After import, run `terraform plan` and align `parent_rule` and other attributes with your intended configuration.
