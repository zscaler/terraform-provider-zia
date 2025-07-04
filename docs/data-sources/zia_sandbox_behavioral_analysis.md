---
subcategory: "Sandbox Policy & Settings"
layout: "zscaler"
page_title: "ZIA: sandbox_behavioral_analysis"
description: |-
  Official documentation https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get
  API documentation https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get
  Gets the custom list of MD5 file hashes that are blocked by Sandbox.
---

# zia_sandbox_behavioral_analysis (Data Source)

* [Official documentation](https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get)
* [API documentation](https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get)

Use the **zia_sandbox_behavioral_analysis** data source to get get the custom list of MD5 file hashes that are blocked by Sandbox

## Example Usage

```hcl
# ZIA Security Policy Settings Data Source
data "zia_sandbox_behavioral_analysis" list_all {}
```

## Argument Reference

This data source can be executed without the need of additional parameters.

## Attribute Reference

N/A
