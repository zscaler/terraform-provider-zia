---
subcategory: "Sandbox Policy & Settings"
layout: "zscaler"
page_title: "ZIA: sandbox_report"
description: |-
  Official documentation https://help.zscaler.com/zia/sandbox-report-use-cases
  API documentation https://help.zscaler.com/zia/sandbox-policy-settings#/behavioralAnalysisAdvancedSettings-get
   Gets a full or summary detail report for an MD5 hash of a file that was analyzed by Sandbox.
---

# zia_sandbox_report (Data Source)

* [Official documentation](https://help.zscaler.com/zia/sandbox-report-use-cases)
* [API documentation](https://help.zscaler.com/zia/sandbox-report#/sandbox/report/{md5Hash}-get)

Use the **zia_sandbox_report** data source gets a full (i.e., complete) or summary detail report for an MD5 hash of a file that was analyzed by Sandbox.

## Example Usage - Obtain Full Sandbox Report

```hcl
data "zia_sandbox_report" "this" {
  md5_hash = "F69CA01D65E6C8F9E3540029E5F6AB92"
  details = "full"
}
```

## Example Usage - Obtain Summarized Sandbox Report

```hcl
data "zia_sandbox_report" "this" {
  md5_hash = "F69CA01D65E6C8F9E3540029E5F6AB92"
  details = "summary"
}
```

## Attributes Reference

### Required

* `md5_hash` - (Required) MD5 hash of the file that was analyzed by Sandbox.

* `details` - (Required) Type of report, full or summary.
  * `full` - A Complete detail report for an MD5 hash of a file that was analyzed by Sandbox
  * `summary` - Summary detail report for an MD5 hash of a file that was analyzed by Sandbox
