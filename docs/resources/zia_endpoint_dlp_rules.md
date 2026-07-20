---
subcategory: "Endpoint Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: endpoint_dlp_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-endpoint-dlp-policy-rules#Rules
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpRules-get
  Creates a new Endpoint DLP policy rule
---

# zia_endpoint_dlp_rules (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-endpoint-dlp-policy-rules#Rules)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpRules-get)

The **zia_dlp_web_rules** resource allows the creation and management of ZIA Endpoint DLP Rules in the Zscaler Internet Access cloud or via the API.

~> **NOTE:** Rule orders must always be contiguous (no gaps). Deleting a rule must be followed by order number re-adjustment of the remaining rules to ensure the API honours the required order.

~> **NOTE:** The `order` attribute must always be a positive whole number starting at 1. Negative numbers and zero are **not supported** and will result in an error.

## Example Usage

```hcl
resource "zia_endpoint_dlp_rules" "this" {
  name                       = "Rule_01"
  description                = "Rule_01"
  action                     = "ALLOW"
  state                      = "ENABLED"
  data_transfer_method       = "APPLICATION_FILE_ACCESS"
  order                      = 1
  rank                       = 0
  severity = "RULE_SEVERITY_HIGH"
  file_types = ["FTCATEGORY_SCZIP", "FTCATEGORY_ACCDB"]
  device_trust_levels = ["LOW_TRUST", "HIGH_TRUST"]
  user_risk_score_levels = ["HIGH", "CRITICAL"]
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

  end_point_application_groups {
    group_id = ["366"]
  }

  auditor {
    id = 43859078
  }

  notification_template {
    id = 5855
  }

  eun_enabled  = true
  eun_template_id  = 11
  receiver {
    id = 1664
  }
}
```

## Example Usage - "Creating Parent Rules and SubRules"

~> **NOTE:** Exception (sub-) rules are managed with the dedicated [`zia_endpoint_dlp_sub_rules`](https://registry.terraform.io/providers/zscaler/zia/latest/docs/resources/zia_endpoint_dlp_sub_rules) resource. Reference the parent’s `rule_id` from the sub-rule so Terraform destroys sub-rules before their parent. See that resource’s documentation for the full example.

## Example Usage - "Configuring Receiver for Endpoint DLP Policy Rule"

```hcl
resource "zia_endpoint_dlp_rules" "this" {
  name                       = "Rule_01"
  description                = "Rule_01"
  action                     = "ALLOW"
  state                      = "ENABLED"
  data_transfer_method       = "APPLICATION_FILE_ACCESS"
  order                      = 1
  rank                       = 0
  severity = "RULE_SEVERITY_HIGH"
  file_types = ["FTCATEGORY_SCZIP", "FTCATEGORY_ACCDB"]
  device_trust_levels = ["LOW_TRUST", "HIGH_TRUST"]
  user_risk_score_levels = ["HIGH", "CRITICAL"]
  without_content_inspection = false
  min_size                   = 10
  end_point_applications {
    zapp_id = ["500000085"]
  }
  receiver {
    id = 1664
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

resource "zia_endpoint_dlp_rules" "this" {
  name                       = "Rule_01"
  description                = "Rule_01"
  action                     = "ALLOW"
  state                      = "ENABLED"
  data_transfer_method       = "APPLICATION_FILE_ACCESS"
  order                      = 1
  rank                       = 0
  end_point_applications {
    zapp_id = ["500000085"]
  }
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

* `name` - (String) The DLP policy rule name.
* `order` - (Number) The rule order of execution for the DLP policy rule with respect to other rules. Must be a positive whole number starting at `1`.

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

* `file_types` - (List of String) The list of file type categories to which the Endpoint DLP policy rule must be applied. For the complete list of supported file type categories refer to the you can use the resource [zia_file_type_categories](https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_file_type_categories) or to the [API Documentation](https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zia/file-type-control-policy/file-type-category-resource-get-file-type-categories)

* `device_trust_levels` - (List of String) The list of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. Supported values: `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`.

* `user_risk_score_levels` - (List of String) The user risk score levels to which the DLP policy rule must be applied. Supported values: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.

### Read-Only

* `sub_rules` - (List of String) The IDs of the exception (sub-) rules associated with this parent rule. This attribute is **computed** and populated from the API after read. Sub-rules are created and managed with the dedicated [`zia_endpoint_dlp_sub_rules`](https://registry.terraform.io/providers/zscaler/zia/latest/docs/resources/zia_endpoint_dlp_sub_rules) resource; they are no longer authored on this resource.

### Optional

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

    ~> **KNOWN ISSUE (API):** The ZIA API currently does not persist `end_point_application_groups` when the rule is created or updated programmatically (via the API/Terraform). The value is accepted without error but is returned empty on subsequent reads, which causes Terraform to report perpetual drift. Configuring the group through the ZIA Admin Portal works as expected. Until the API is fixed, add `end_point_application_groups` to `lifecycle.ignore_changes` to suppress the drift:

    ```hcl
    resource "zia_endpoint_dlp_rules" "this" {
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
