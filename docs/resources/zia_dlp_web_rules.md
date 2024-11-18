---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_web_rules"
description: |-
  Creates and manages ZIA DLP Web Rules.
---

# Resource: zia_dlp_web_rules

The **zia_dlp_web_rules** resource allows the creation and management of ZIA DLP Web Rules in the Zscaler Internet Access cloud or via the API.

⚠️ **WARNING:** Zscaler Internet Access DLP supports a maximum of 127 Web DLP Rules to be created via API.

## Example Usage - "ALL_OUTBOUND" File Type"

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
  file_types                 = [ "ALL_OUTBOUND" ]
  zscaler_incident_receiver  = false
  without_content_inspection = false
  user_risk_score_levels     = [ "LOW", "MEDIUM", "HIGH", "CRITICAL" ]
  severity                   = "RULE_SEVERITY_HIGH"
  dlp_engines {
    id = [ data.zia_dlp_engines.this.id ]
  }
}
```

```hcl
// Retrieve a custom URL Category by Name
data "zia_url_categories" "this"{
    configured_name = "Example"
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
    id = [ data.zia_url_categories.this.val ]
  }
  icap_server {
    id = data.zia_dlp_icap_servers.this.id
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

⚠️ **WARNING:** Destroying a parent rule will also destroy all subrules

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
  parent_rule = zia_dlp_web_rules.parent_rule.id
}
```

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
  * `ALLOW`
  * `ICAP_RESPONSE`

* `state` - (Optional) Enables or disables the DLP policy rule.. The supported values are:
  * `DISABLED`
  * `ENABLED`

* `file_types` - (Optional) The list of file types to which the DLP policy rule must be applied. For the complete list of supported file types refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post)

  * ~> Note: `BITMAP`, `JPEG`, `PNG`, and `TIFF` file types are exclusively supported when optical character recognition `ocr_enabled` is set to `true` for DLP rules with content inspection.

  * ~> Note: `ALL_OUTBOUND` file type is applicable only when the predefined DLP engine called `EXTERNAL` is used and when the attribute `without_content_inspection` is set to `false`.

  * ~> Note: `ALL_OUTBOUND` file type cannot be used alongside any other file type.

* `cloud_applications` - (Optional) The list of cloud applications to which the DLP policy rule must be applied.. For the complete list of supported file types refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post)

* `severity` - (Optional) Indicates the severity selected for the DLP rule violation: Returned values are:  `RULE_SEVERITY_HIGH`, `RULE_SEVERITY_MEDIUM`, `RULE_SEVERITY_LOW`, `RULE_SEVERITY_INFO`

* `user_risk_score_levels` (Optional) - Indicates the user risk score level selectedd for the DLP rule violation: Returned values are: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

* `parent_rule`(Optional) - The unique identifier of the parent rule under which an exception rule is added. The rule rank must be set to `0`

    ~> **Note**: Exception rules can be configured only when the inline DLP rule evaluation type is set to evaluate all DLP rules in the DLP Advanced Settings. To learn more, see [Configuring DLP Advanced Settings](https://help.zscaler.com/%22/zia/configuring-dlp-advanced-settings/%22)

    ~> **Note**: It is not possible to add existing rules as as subrules under the parent rule.

* `sub_rules`(List) - The list of exception rules added to a parent rule. The rule rank must be set to `0`

    ~> **Note**: All attributes within the WebDlpRule model are applicable to the sub-rules. Values for each rule are specified by using the WebDlpRule object Exception rules can be configured only when the inline DLP rule evaluation type is set to evaluate all DLP rules in the DLP Advanced Settings. To learn more, see [Configuring DLP Advanced Settings](https://help.zscaler.com/%22/zia/configuring-dlp-advanced-settings/%22)

* `notification_template` - (Optional) The template used for DLP notification emails.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `auditor` - (Optional) The auditor to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `url_categories` - (Optional) The list of URL categories to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity
  ~> **NOTE** When associating a URL category, you can use the `zia_url_categories` resource or data source; however, you must export the attribute `val`

* `dlp_engines` - (Optional) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity. Maximum of up to `4` dlp engines. When not used it implies `Any` to apply the rule to all locations.

* `locations` - (Optional) The Name-ID pairs of locations to which the DLP policy rule must be applied. Maximum of up to `8` locations. When not used it implies `Any` to apply the rule to all locations.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `location_groups` - (Optional) The Name-ID pairs of locations groups to which the DLP policy rule must be applied. Maximum of up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `users` - (Optional) The Name-ID pairs of users to which the DLP policy rule must be applied. Maximum of up to `4` users. When not used it implies `Any` to apply the rule to all users.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `groups` - (Optional) The Name-ID pairs of groups to which the DLP policy rule must be applied. Maximum of up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
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

* `workload_groups` (Optional) The list of preconfigured workload groups to which the policy must be applied
  * `id` - (Optional) A unique identifier assigned to the workload group
  * `name` - (Optional) The name of the workload group

|              Inspection Type              |                   File Types                                  |
|:------------------------------------------|:--------------------------------------------------------------|
|--------------------------------|-------------------------------------------------------|
|                                           | `FTCATEGORY_ENCRYPT`, `FTCATEGORY_P7Z`,                       |
|                                           | `FTCATEGORY_BZIP2`, `FTCATEGORY_CAB`,                         |
|                                           | `FTCATEGORY_FCL`, `FTCATEGORY_GZIP`,                          |
|       `WITHOUT INSPECTION`                | `FTCATEGORY_ISO`, `FTCATEGORY_LZH`,                           |
|                                           | `FTCATEGORY_RAR`,`FTCATEGORY_STUFFIT`,                        |
|                                           | `FTCATEGORY_TAR`, `FTCATEGORY_XZ`,                            |
|                                           | `FTCATEGORY_ZIP`, `FTCATEGORY_SCZIP`,                         |
|                                           | `FTCATEGORY_ZIPX`                                             |
|--------------------------------|--------------------------------------------------------|
|--------------------------------|--------------------------------------------------------|
|                                           | `FTCATEGORY_ACCDB`, `FTCATEGORY_APPLE_DOCUMENTS`,             |
|                                           | `FTCATEGORY_ASM`, `FTCATEGORY_AU3`,                           |
|                                           | `FTCATEGORY_BASH_SCRIPTS`, `FTCATEGORY_BASIC_SOURCE_CODE`,    |
|                                           | `FTCATEGORY_BCP`,`FTCATEGORY_BITMAP`,                         |
|                                           | `FTCATEGORY_BORLAND_CPP_FILES`, `FTCATEGORY_C_FILES`,         |
|                                           | `FTCATEGORY_COBOL`,`FTCATEGORY_CSV`,                          |
|                                           | `FTCATEGORY_CSX`,`FTCATEGORY_DAT`,                            |
|                                           | `FTCATEGORY_DCM`, `FTCATEGORY_DELPHI`,                        |
|                                           | `FTCATEGORY_F_FILES`,`FTCATEGORY_FOR`,                        |
|                                           | `FTCATEGORY_FORM_DATA_POST`, `FTCATEGORY_DSP`,                |
|                                           | `FTCATEGORY_EML_FILES`, `FTCATEGORY_GO_FILES`,                |
|                                           | `FTCATEGORY_HTTP`,`FTCATEGORY_IFC`,                           |
|                                           | `FTCATEGORY_INCLUDE_FILES`,`FTCATEGORY_INF`,                  |
|                                           | `FTCATEGORY_JAVA_FILES`,`FTCATEGORY_JPEG`,                    |
|                                           | `FTCATEGORY_JSON`,`FTCATEGORY_LOG_FILES`,                     |
|                                           | `FTCATEGORY_MAKE_FILES`, `FTCATEGORY_MATLAB_FILES`,           |
|                                           | `FTCATEGORY_MS_EXCEL`,`FTCATEGORY_MS_MDB`,                    |
|                                           | `FTCATEGORY_MS_MSG`,`FTCATEGORY_MS_POWERPOINT`,               |
|          `WITH INSPECTION`                | `FTCATEGORY_MS_PUB`,`FTCATEGORY_MS_RTF`,                      |
|                                           | `FTCATEGORY_MS_WORD`,`FTCATEGORY_MSC`,                        |
|                                           | `FTCATEGORY_NATVIS`,`FTCATEGORY_OLM`,                         |
|                                           | `FTCATEGORY_OPEN_OFFICE_PRESENTATIONS`,                       |
|                                           | `FTCATEGORY_OPEN_OFFICE_SPREADSHEETS`,                        |
|                                           | `FTCATEGORY_MS_CPP_FILES`, `FTCATEGORY_PDF_DOCUMENT`,         |
|                                           | `FTCATEGORY_PERL_FILES`,`FTCATEGORY_PNG`,                     |
|                                           | `FTCATEGORY_POD`,`FTCATEGORY_POWERSHELL`,                     |
|                                           | `FTCATEGORY_PYTHON`, `FTCATEGORY_RES_FILES`,                  |
|                                           | `FTCATEGORY_RPY`,`FTCATEGORY_RSP`,                            |
|                                           | `FTCATEGORY_RUBY_FILES`,`FTCATEGORY_SAS`,                     |
|                                           | `FTCATEGORY_SC`,`FTCATEGORY_SCALA`,                           |
|                                           | `FTCATEGORY_SCT`,`FTCATEGORY_SHELL_SCRAP`,                    |
|                                           | `FTCATEGORY_SQL`,`FTCATEGORY_TABLEAU_FILES`,                  |
|                                           | `FTCATEGORY_TIFF`, `FTCATEGORY_TLH`                           |
|                                           | `FTCATEGORY_TLI`,`FTCATEGORY_TXT`,                            |
|                                           | `FTCATEGORY_UNK_TXT`,`FTCATEGORY_XAML`                        |
|                                           | `FTCATEGORY_VISUAL_BASIC_SCRIPT`,                             |
|                                           | `FTCATEGORY_VISUAL_CPP_FILES`,`FTCATEGORY_VSDX`,              |
|                                           | `FTCATEGORY_X1B`,`FTCATEGORY_VISUAL_BASIC_FILES`,             |
|                                           | `FTCATEGORY_XML`, `FTCATEGORY_YAML_FILES`,                    |
|                                           | `FTCATEGORY_JAVA_APPLET`,`FTCATEGORY_JAVASCRIPT`              |
|                                           | `FTCATEGORY_WINDOWS_SCRIPT_FILES`                             |
|------------------------------|----------------------------------------------------------|

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_dlp_web_rules** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import zia_dlp_web_rules.example <rule_id>
```

or

```shell
terraform import zia_dlp_web_rules.example <rule_name>
```
