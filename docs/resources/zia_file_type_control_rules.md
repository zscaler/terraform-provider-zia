---
subcategory: "File Type Control Policy"
layout: "zscaler"
page_title: "ZIA: file_type_control_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/about-file-type-control
  API documentation https://help.zscaler.com/zia/file-type-control-policy#/fileTypeRules-post
  Creates and manages ZIA Cloud firewall filtering rule.
---

# zia_file_type_control_rules (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-file-type-control)
* [API documentation](https://help.zscaler.com/zia/file-type-control-policy#/fileTypeRules-post)

The **zia_file_type_control_rules** resource allows the creation and management of ZIA file type control rules in the Zscaler Internet Access.

## Example Usage

```hcl
data "zia_department_management" "engineering" {
 name = "Engineering"
}

data "zia_group_management" "normal_internet" {
    name = "Normal_Internet"
}

data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_policy"
  app_class   = ["AI_ML"]
}

resource "zia_file_type_control_rules" "this" {
    name               = "Terraform_File_Type01"
    description        = "Terraform_File_Type01"
    state              = "ENABLED"
    order              = 1
    rank               = 7
    filtering_action   = "BLOCK"
    operation          = "DOWNLOAD"
    active_content     = true
    unscannable        = false
    device_trust_levels = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    file_types         = ["FTCATEGORY_MS_WORD", "FTCATEGORY_MS_POWERPOINT", "FTCATEGORY_PDF_DOCUMENT", "FTCATEGORY_MS_EXCEL"]
    protocols          = ["FOHTTP_RULE", "FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
    cloud_applications = tolist([for app in data.zia_cloud_applications.this.applications : app["app"]])

    departments {
        id = [ data.zia_department_management.engineering.id ]
    }
    groups {
        id = [ data.zia_group_management.normal_internet.id ]
    }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the Firewall Filtering policy rule
* `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.

### Optional

* `description` - (String) Additional information about the File Type rule.
* `state` - (String) Rule State. Supported Values: `ENABLED` and `DISABLED`
* `order` - (Integer) Order of policy execution with respect to other file-type policies.
* `filtering_action` - (String) Action taken when traffic matches policy. Supported values: `BLOCK`, `CAUTION`, `ALLOW`.
* `rank` - (Integer) Admin rank of the admin who creates this rule. Supported values: Range `1` to `7`
* `capture_pcap` - (Boolean) A Boolean value that indicates whether packet capture (PCAP) is enabled.
* `operation` - (String) File operation performed. Supported Values: `UPLOAD`, `DOWNLOAD` or `UPLOAD_DOWNLOAD`
* `active_content` - (Boolean) Flag to check whether a file has active content.
    **NOTE** The attribute can only be set when the `file_types` list contain the following values: `FTCATEGORY_MS_WORD`, `FTCATEGORY_MS_POWERPOINT`, `FTCATEGORY_PDF_DOCUMENT`, `FTCATEGORY_MS_EXCEL`.
* `unscannable` - (Boolean) Flag to check whether a file is scannable.
* `file_types` - (List of Strings) File type categories for which the policy is applied. If not set, the rule is applied across all file types.
* `min_size` - (Integer) Minimum file size (in KB) used for evaluation of the rule. Values between: `0` and `409600`
* `max_size` - (Integer) Maximum file size (in KB) used for evaluation of the rule. Values between: `0` and `409600`
* `protocols` - (List of Strings) Protocol for the given rule. Supported Values are: `ANY_RULE`, `SMRULEF_CASCADING_ALLOWED`, `FOHTTP_RULE`, `FTP_RULE`, `HTTPS_RULE`, `HTTP_RULE`
* `cloud_applications` - (List of Strings) The list of cloud applications to which the File Type Control policy rule must be applied. To retrieve the list of cloud applications, use the data source: `zia_cloud_applications`
* `device_trust_levels` - (List of Strings) List of device trust levels for which the rule must be applied. While the High Trust, Medium Trust, or Low Trust evaluation is applicable only to Zscaler Client Connector traffic, Unknown evaluation applies to all traffic. Supported values: `ANY`, `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`

* `url_categories` - (Set of Strings) The list of URL categories to which the DLP policy rule must be applied.

* `locations` - (Optional) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (String) Identifier that uniquely identifies an entity

* `location_groups` - (Optional) You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (String) Identifier that uniquely identifies an entity

* `users` - (Optional) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (String) Identifier that uniquely identifies an entity

* `groups` - (Optional) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (String) Identifier that uniquely identifies an entity

* `departments` - (Optional) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (String) Identifier that uniquely identifies an entity

* `time_windows` - (Optional) You can manually select up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (String) Identifier that uniquely identifies an entity

* `labels` Labels that are applicable to the rule.
      - `id` - (String) Identifier that uniquely identifies an entity

* `zpa_app_segments` List of Source IP Anchoring-enabled ZPA Application Segments for which this rule is applicable
      - `id` - (String) Identifier that uniquely identifies an entity

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_file_type_control_rules** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import zia_file_type_control_rules.example <rule_id>
```

or

```shell
terraform import zia_file_type_control_rules.example <rule_name>
```
