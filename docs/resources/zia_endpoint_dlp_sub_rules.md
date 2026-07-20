---
subcategory: "Endpoint Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: endpoint_dlp_sub_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-endpoint-dlp-policy-rules#Rules
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpRules-get
  Creates and manages an exception (sub) rule under an Endpoint DLP policy rule
---

# zia_endpoint_dlp_sub_rules (Resource)

* [Official documentation](https://help.zscaler.com/zia/configuring-endpoint-dlp-policy-rules#Rules)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpRules-get)

The **zia_endpoint_dlp_sub_rules** resource allows the creation and management of exception (sub) rules under an existing ZIA Endpoint DLP policy rule. Each sub-rule belongs to a single parent rule (`zia_endpoint_dlp_rules`) and is created, updated, and deleted through the parent rule's dedicated sub-rule endpoint.

⚠️ **WARNING:** A sub-rule is bound to its parent. Deleting the parent rule on the service side removes all of its sub-rules automatically. Always reference the parent's `rule_id` (see the example) so Terraform destroys the sub-rules **before** the parent and never leaves orphaned resources in state.

~> **NOTE:** Changing `parent_rule` forces the sub-rule to be recreated under the new parent, because the service binds each sub-rule to its parent through the request path.

~> **NOTE:** Sub-rule orders must always be contiguous (no gaps) within a parent. Deleting a sub-rule must be followed by order number re-adjustment of the remaining sub-rules of that parent to ensure the API honours the required order.

~> **NOTE:** The `order` attribute must always be a positive whole number starting at 1. Negative numbers and zero are **not supported** and will result in an error.

~> **NOTE:** Exception rules can be configured only when the inline DLP rule evaluation type is set to evaluate all DLP rules in the DLP Advanced Settings. To learn more, see [Configuring DLP Advanced Settings](https://help.zscaler.com/zia/configuring-dlp-advanced-settings) or leverage the resource [zia_dlp_global_options](https://registry.terraform.io/providers/zscaler/zia/latest/docs/resources/zia_dlp_global_options) to enable this feature.

## Example Usage - "Minimal Sub-rule"

```hcl
# Parent Endpoint DLP rule
resource "zia_endpoint_dlp_rules" "parent" {
  name                 = "Rule_01"
  description          = "Rule_01"
  action               = "ALLOW"
  state                = "ENABLED"
  data_transfer_method = "REMOVABLE_DRIVE_TRANSFER"
  order                = 1
  rank                 = 0
  severity             = "RULE_SEVERITY_HIGH"
  file_types           = ["FTCATEGORY_SCZIP", "FTCATEGORY_ACCDB"]
}

# Exception (sub) rule managed under the parent above
resource "zia_endpoint_dlp_sub_rules" "subrule1" {
  name                 = "SubRule01"
  description          = "SubRule01"
  action               = "ALLOW"
  state                = "ENABLED"
  data_transfer_method = "REMOVABLE_DRIVE_TRANSFER"
  order                = 1
  rank                 = 0
  severity             = "RULE_SEVERITY_HIGH"

  # Reference the parent's rule_id so Terraform orders create/destroy correctly.
  parent_rule = zia_endpoint_dlp_rules.parent.rule_id
}
```

## Example Usage - "Multiple Sub-rules with Endpoint Applications"

The parent rule can be resolved with the `zia_endpoint_dlp_rules` data source, and
each sub-rule declares its own contiguous `order` within that parent. This mirrors
a full production configuration with DLP engines, notification templates, an
Incident Receiver, and endpoint applications.

```hcl
# Look up the parent Endpoint DLP rule by name
data "zia_endpoint_dlp_rules" "this" {
  name = "Rule_01"
}

resource "zia_endpoint_dlp_sub_rules" "subrule1" {
  name                       = "Rule_01_Subrule1"
  description                = "Rule_01_Subrule1"
  action                     = "ALLOW"
  state                      = "ENABLED"
  data_transfer_method       = "APPLICATION_FILE_ACCESS"
  order                      = 1
  rank                       = 0
  severity                   = "RULE_SEVERITY_HIGH"
  file_types                 = ["FTCATEGORY_SCZIP", "FTCATEGORY_ACCDB"]
  device_trust_levels        = ["LOW_TRUST", "HIGH_TRUST"]
  user_risk_score_levels     = ["HIGH", "CRITICAL"]
  parent_rule                = data.zia_endpoint_dlp_rules.this.id
  without_content_inspection = false
  min_size                   = 10

  dlp_engines {
    id = [3, 4]
  }

  device_groups {
    id = [35227159]
  }

  devices {
    id = [36452290]
  }

  departments {
    id = [68759309]
  }

  groups {
    id = [165437858]
  }

  users {
    id = [165214882]
  }

  labels {
    id = [4900918]
  }

  end_point_applications {
    zapp_id = ["500000085"]
  }

  auditor {
    id = 43859078
  }

  notification_template {
    id = 5855
  }

  eun_enabled     = true
  eun_template_id = 11

  receiver {
    id = 1664
  }
}

resource "zia_endpoint_dlp_sub_rules" "subrule2" {
  name                       = "Rule_01_Subrule2"
  description                = "Rule_01_Subrule2"
  action                     = "ALLOW"
  state                      = "ENABLED"
  data_transfer_method       = "APPLICATION_FILE_ACCESS"
  order                      = 2
  rank                       = 0
  severity                   = "RULE_SEVERITY_HIGH"
  file_types                 = ["FTCATEGORY_SCZIP", "FTCATEGORY_ACCDB"]
  device_trust_levels        = ["LOW_TRUST", "HIGH_TRUST"]
  user_risk_score_levels     = ["HIGH", "CRITICAL"]
  parent_rule                = data.zia_endpoint_dlp_rules.this.id
  without_content_inspection = false
  min_size                   = 10

  dlp_engines {
    id = [3, 4]
  }

  end_point_applications {
    zapp_id = ["500000085"]
  }

  receiver {
    id = 1664
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (String) The DLP policy rule name.
* `order` - (Number) The rule order of execution for the sub-rule with respect to other sub-rules of the same parent. Must be a positive whole number starting at `1`.
* `parent_rule` - (Number) The unique identifier of the parent rule under which this exception (sub) rule is created. Use the parent resource's `rule_id` (integer). Changing this value **recreates** the sub-rule under the new parent.

### Optional

* `description` - (String) The description of the DLP policy rule.
* `rank` - (Number) Admin rank of the admin who creates this rule. Supported range: `0` to `7`.
* `state` - (String) Enables or disables the DLP policy rule. Supported values: `ENABLED`, `DISABLED`.
* `action` - (String) The action taken when traffic matches the DLP policy rule criteria. Supported values: `ANY`, `NONE`, `BLOCK`, `CONFIRM`, `ALLOW`.
* `severity` - (String) Indicates the severity selected for the DLP rule violation. Supported values: `RULE_SEVERITY_HIGH`, `RULE_SEVERITY_MEDIUM`, `RULE_SEVERITY_LOW`, `RULE_SEVERITY_INFO`.
* `data_transfer_method` - (String) The data transfer method to which the DLP policy rule must be applied (e.g. `APPLICATION_FILE_ACCESS`). This value must be set to `APPLICATION_FILE_ACCESS` when configuring `end_point_applications` or `end_point_application_groups`.
* `network_type` - (String) The network type to which the DLP policy rule must be applied.
* `min_size` - (Number) The minimum file size (in KB) used for evaluation of the DLP policy rule. Supported range: `0` to `96000`.
* `without_content_inspection` - (Boolean) Indicates a DLP policy rule without content inspection, when the value is set to `true`.
* `external_auditor_email` - (String) The email address of an external auditor to whom DLP email notifications are sent. Cannot be combined with `auditor`.
* `eun_enabled` - (Boolean) Indicates whether the End User Notification (EUN) is enabled for the DLP policy rule.
* `eun_template_id` - (Number) The unique identifier of the End User Notification (EUN) template.
* `uc_template_id` - (Number) The unique identifier of the user confirmation template.

* `file_types` - (List of String) The list of file type categories to which the sub-rule must be applied. For the complete list of supported file type categories refer to the data source [zia_file_type_categories](https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_file_type_categories).

* `device_trust_levels` - (List of String) The list of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. Supported values: `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`.

* `user_risk_score_levels` - (List of String) The user risk score levels to which the DLP policy rule must be applied. Supported values: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.

* `dlp_engines` - (Block) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity. Maximum of up to `4` DLP engines.

* `users` - (Block) The Name-ID pairs of users to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `groups` - (Block) The Name-ID pairs of groups to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `departments` - (Block) The Name-ID pairs of departments to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `devices` - (Block) The Name-ID pairs of devices to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `device_groups` - (Block) The Name-ID pairs of device groups to which the DLP policy rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `resources` - (Block) The Name-ID pairs of resources to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `resource_groups` - (Block) The Name-ID pairs of resource groups to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `labels` - (Block) The Name-ID pairs of rule labels associated to the DLP policy rule.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `auditor` - (Block) The auditor to which the DLP policy rule must be applied. Cannot be combined with `external_auditor_email`; when set, `notification_template` must also be set.
  * `id` - (Number) Identifier that uniquely identifies an entity.

* `notification_template` - (Block) The template used for DLP notification emails.
  * `id` - (Number) Identifier that uniquely identifies an entity.

* `end_point_applications` - (Block) The endpoint applications to which the DLP policy rule must be applied. Can only be set when `data_transfer_method` is `APPLICATION_FILE_ACCESS`.
  * `zapp_id` - (List of String) The unique identifiers of the endpoint applications.

* `end_point_application_groups` - (Block) The endpoint application groups to which the DLP policy rule must be applied. Can only be set when `data_transfer_method` is `APPLICATION_FILE_ACCESS`.
  * `group_id` - (List of String) The unique identifiers of the endpoint application groups.

    ~> **KNOWN ISSUE (API):** The ZIA API currently does not persist `end_point_application_groups` when the sub-rule is created or updated programmatically (via the API/Terraform). The value is accepted without error but is returned empty on subsequent reads, which causes Terraform to report perpetual drift. Configuring the group through the ZIA Admin Portal works as expected. Until the API is fixed, add `end_point_application_groups` to `lifecycle.ignore_changes` to suppress the drift:

    ```hcl
    resource "zia_endpoint_dlp_sub_rules" "this" {
      # ...
      end_point_application_groups {
        group_id = ["366"]
      }

      lifecycle {
        ignore_changes = [end_point_application_groups]
      }
    }
    ```

* `receiver` - (Block) The Zscaler Incident Receiver associated with the DLP policy rule. Only the `id` must be set; the receiver `type` is always sent as `ZIR`.
  * `id` - (Required, String) The unique identifier of the Zscaler Incident Receiver.

### Read-Only

* `id` - (String) The Terraform resource ID (the sub-rule's numeric ID as a string).
* `sub_rule_id` - (Number) The service-assigned identifier of the sub-rule.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Because a sub-rule can only be located through its parent, the import ID must carry both identifiers in the form `<parentRuleID>:<subRuleID>` or `<parentRuleID>:<subRuleName>`.

For example:

```shell
terraform import zia_endpoint_dlp_sub_rules.example 1839792:1839815
```

or

```shell
terraform import zia_endpoint_dlp_sub_rules.example 1839792:SubRule01
```

After import, run `terraform plan` and align `parent_rule` and other attributes with your intended configuration.
