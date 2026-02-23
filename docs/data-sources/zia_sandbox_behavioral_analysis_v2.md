---
subcategory: "Sandbox Policy & Settings"
layout: "zscaler"
page_title: "ZIA: sandbox_behavioral_analysis_v2"
description: |-
  Official documentation https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get
  API documentation https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get
  Gets the custom list of MD5 file hashes that are blocked by Sandbox.
---

# zia_sandbox_behavioral_analysis_v2 (Data Source)

* [Official documentation](https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get)
* [API documentation](https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get)

Use the **zia_sandbox_behavioral_analysis_v2** data source to get the custom list of MD5 file hashes that are blocked by Sandbox.

## Example Usage

```hcl
data "zia_sandbox_behavioral_analysis_v2" "this" {}
```

## Argument Reference

This data source does not require any arguments.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `md5_hash_value_list` - A list of MD5 hash values with metadata for sandbox blocking. Each entry contains:
  * `url` - The MD5 hash identifier for the entry.
  * `url_comment` - A comment describing the MD5 hash entry.
  * `type` - The type of the entry (e.g. `CUSTOM_FILEHASH_ALLOW`, `CUSTOM_FILEHASH_DENY`, `MALWARE`).
