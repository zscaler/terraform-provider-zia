---
subcategory: "Sandbox Policy & Settings"
layout: "zscaler"
page_title: "ZIA: sandbox_behavioral_analysis_v2"
description: |-
  Official documentation https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get
  API documentation https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get
  Updates the custom list of MD5 file hashes that are blocked by Sandbox.
---

# zia_sandbox_behavioral_analysis_v2 (Resource)

* [Official documentation](https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get)
* [API documentation](https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get)

The **zia_sandbox_behavioral_analysis_v2** resource updates the custom list of MD5 file hashes that are blocked by Sandbox. This overwrites a previously generated blocklist. If you need to completely erase the blocklist, submit an empty list.

**Note**: The ZIA API has introduced a new payload format, that required this new resource. The resource `zia_sandbox_behavioral_analysis` is still available in the provider; however, we encourage users to move to this new format. The resource `zia_sandbox_behavioral_analysis` will be eventually deprecated

~> NOTE: This an Early Access feature.

## Example Usage - Add MD5 Hashes to Sandbox

```hcl
# Add MD5 Hashes to Sandbox
resource "zia_sandbox_behavioral_analysis_v2" "this" {
  md5_hash_value_list {
    url         = "4EE43B71BB89CB9CBF7784495AE8D0DF"
    url_comment = "4EE43B71BB89CB9CBF7784495AE8D0DF"
    type        = "CUSTOM_FILEHASH_ALLOW"
  }

  md5_hash_value_list {
    url         = "8350dED6D39DF158E51D6CFBE36FB012"
    url_comment = "8350dED6D39DF158E51D6CFBE36FB012"
    type        = "CUSTOM_FILEHASH_DENY"
   }
}
```

## Example Usage - Remove All MD5 Hashes to Sandbox

```hcl
# Remove All MD5 Hashes to Sandbox
resource "zia_sandbox_behavioral_analysis_v2" "this" {
  md5_hash_value_list {}
}
```

## Example Usage - Remove All MD5 Hashes to Sandbox

```hcl
# Remove All MD5 Hashes to Sandbox
resource "zia_sandbox_behavioral_analysis_v2" "this" {
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `md5_hash_value_list` - (Optional) A custom list of MD5 hash values with metadata for sandbox blocking. Each block supports the following attributes:
  * `url` - (Optional) The MD5 hash identifier for the entry.
  * `url_comment` - (Optional) A comment describing the MD5 hash entry.
  * `type` - (Optional) The type of the entry. Supported values: `CUSTOM_FILEHASH_ALLOW`, `CUSTOM_FILEHASH_DENY`, `MALWARE`.

~> To remove all MD5 hashes from the blocklist, define the resource with an empty `md5_hash_value_list {}` block or remove all blocks entirely. This sends an empty list to the API, which clears the blocklist.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_sandbox_behavioral_analysis_v2** can be imported by using `sandbox_settings` as the import ID.

For example:

```shell
terraform import zia_sandbox_behavioral_analysis_v2.example sandbox_settings
```
