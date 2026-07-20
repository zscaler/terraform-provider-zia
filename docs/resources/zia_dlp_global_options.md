---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_global_options"
description: |-
  Official documentation https://help.zscaler.com/legacy-apis/data-loss-prevention#/webDlpGlobalOptions-get
  API documentation https://help.zscaler.com/legacy-apis/data-loss-prevention#/webDlpGlobalOptions-get
  Updates the existing DLP Advanced Settings.
---

# zia_dlp_global_options (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-dlp-engines)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/dlpEngines-get)

Use the **zia_dlp_global_options** resource allows the management of ZIA DLP Advanced Settings in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
resource "zia_url_categories" "this" {
  super_category      = "USER_DEFINED"
  configured_name     = "MCAS Unsanctioned Apps"
  description         = "MCAS Unsanctioned Apps"
  keywords            = ["microsoft"]
  custom_category     = true
  type                = "URL_CATEGORY"
  urls = [
    ".coupons.com",
    ".resource.alaskaair.net",
    ".techrepublic.com",
    ".dailymotion.com",
    ".osiriscomm.com",
    ".uefa.com",
    ".Logz.io",
    ".alexa.com",
    ".baidu.com",
    ".cnn.com",
    ".level3.com",
  ]
}

resource "zia_dlp_global_options" "this" {
  applications = [ "ELEVENX_AI", "ONE_HUNDRED_NINE_AI"]
  urls = [ "google.com", "yahoo.com", "bing.com" ]
#   http_get_custom_url_categories = [ "Zscaler Cloud", "Zscaler Internet Access" ]
  exempt_url_encoded_data = true
  enable_npk_edm_templates = false
  enable_npk_edm_templates_for_org = false
  enable_inline_dlp_ocr = true
  enable_casb_ocr = true
  enable_email_dlp_ocr = false
  enable_evaluate_all_dlp_rules = true
  enable_edm_popular_format = false
  url_categories {
    id = [zia_url_categories.this.val]
  }
}

```

## Argument Reference

The following arguments are supported:

### Read-Only

* `applications` - (List) The list of cloud applications to which must be applied. For the complete list of supported file types refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post). To retrieve the list of cloud applications, use the data source: `zia_cloud_applications`

* `urls` - (List) List of URLs exempted from DLP evaluation

* `http_get_custom_url_categories` - (List) List of URL Categories to associate with Inspect HTTP GET Requests

* `exempt_url_encoded_data` - (Bool) Indicates whether or not URL encoded data from DLP evaluation is exempted

* `enable_npk_edm_templates` - (Bool) Indicates whether EDM with No Primary Keys is enabled. If enabled, you select the fields from the template that must match, fields that are optional, and how many optional fields are required to match. Before enabling EDM with No Primary Keys, all existing EDM schema (templates, dictionaries, engines, and policies) must be deleted or unassigned. New schema with no primary keys functionality can be created after the feature is enabled.

* `enable_npk_edm_templates_for_org` - (Bool) Read Only. Indicates whether EDM with No Primary Keys is enabled for the organization. To enable this field, you must first enable the `enable_npk_edm_templates` field.

* `enable_inline_dlp_ocr` - (Bool) Indicates whether optical character recognition (OCR) for Zscaler DLP engines to scan images for text content in data in transit is enabled

* `enable_casb_ocr` - (Bool) Indicates whether SaaS Security for Zscaler DLP engines to scan images for text content in data at rest is enabled

* `enable_email_dlp_ocr` - (Bool) Indicates whether Outbound Email DLP for Zscaler DLP engines to scan images for text content in outbound emails is sent to external domains

* `enable_evaluate_all_dlp_rules` - (Bool) Indicates whether DLP engines evaluate all rules or stop when a matching rule is found. When enabled, if multiple rules match, the rule engine selects the rule with the most restrictive action.

* `enable_edm_popular_format` - (Bool) Indicates whether EDM for popular formats is enabled. When enabled, the EDM scans SSN, SIN, and CCN dictionaries for a popular format, and the EDM match only succeeds if the entered format matches the popular data format type.

* `url_categories` - (Optional) The list of URL categories to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity
  ~> **NOTE** When associating a URL category, you can use the `zia_url_categories` resource or data source; however, you must export the attribute `val`. The `val` attribute is available on both the resource and data source, making it consistent for referencing URL categories in DLP web rules.


## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_dlp_global_options** can be imported by using `dlp_global_options` as the import ID.

For example:

```shell
terraform import zia_dlp_global_options.this "dlp_global_options"
```
