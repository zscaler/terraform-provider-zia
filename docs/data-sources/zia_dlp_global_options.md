---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_global_options"
description: |-
  Official documentation https://help.zscaler.com/legacy-apis/data-loss-prevention#/webDlpGlobalOptions-get
  API documentation https://help.zscaler.com/legacy-apis/data-loss-prevention#/webDlpGlobalOptions-get
  Get information about DLP Advanced Settings information
---

# zia_dlp_global_options (Data Source)

* [Official documentation](https://help.zscaler.com/legacy-apis/data-loss-prevention#/webDlpGlobalOptions-get)
* [API documentation](https://help.zscaler.com/legacy-apis/data-loss-prevention#/webDlpGlobalOptions-get)

Use the **zia_dlp_global_options** data source to get information about a ZIA DLP Advanced Settings information in the Zscaler Internet Access cloud or via the API.

## Example Usage - Retrieves DLP Advanced Settings information

```hcl
#
data "zia_dlp_global_options" "this"{
}
```

## Argument Reference

The following arguments are supported:

### Read-Only

* `applications` - (List) List of cloud applications exempted from DLP evaluation
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
* `url_categories` - (List) List of custom URL categories exempted from DLP evaluation. Supports the following attributes:
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` - (String) Identifier that uniquely identifies an entity
  * `extensions` - (Map of String)
