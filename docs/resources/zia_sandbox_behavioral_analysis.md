---
subcategory: "Sandbox Policy & Settings"
layout: "zscaler"
page_title: "ZIA: sandbox_behavioral_analysis"
description: |-
  Updates the custom list of MD5 file hashes that are blocked by Sandbox.
---

# Resource: zia_sandbox_behavioral_analysis

The **zia_sandbox_behavioral_analysis** resource updates the custom list of MD5 file hashes that are blocked by Sandbox. This overwrites a previously generated blocklist. If you need to completely erase the blocklist, submit an empty list.

**Note**: Only the file types that are supported by Sandbox analysis can be blocked using MD5 hashes.

## Example Usage - Add MD5 Hashes to Sandbox

```hcl
# Add MD5 Hashes to Sandbox
resource "zia_sandbox_behavioral_analysis" "this" {
  file_hashes_to_be_blocked = [
        "42914d6d213a20a2684064be5c80ffa9",
        "c0202cf6aeab8437c638533d14563d35",
  ]
}
```

## Example Usage - Remove All MD5 Hashes to Sandbox

```hcl
# Remove All MD5 Hashes to Sandbox
resource "zia_sandbox_behavioral_analysis" "this" {
  file_hashes_to_be_blocked = []
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `file_hashes_to_be_blocked` - (Required) A custom list of unique MD5 file hashes that must be blocked by Sandbox. A maximum of 10000 MD5 file hashes can be blocked.

**Note 3**: The Sandbox only supports MD5 hashes. The provider will validate the MD5 format prior to submission.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_sandbox_behavioral_analysis** can be imported by using `sandbox_settings` as the import ID.

For example:

```shell
terraform import zia_sandbox_behavioral_analysis.example sandbox_settings
```
