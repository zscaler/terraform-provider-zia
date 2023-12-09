---
subcategory: "Sandbox Settings"
layout: "zscaler"
page_title: "ZIA: sandbox_file_submission"
description: |-
  Submits raw or archive files (e.g., ZIP) to Sandbox for analysis.
---

# Resource: zia_sandbox_file_submission

The **zia_sandbox_file_submission** resource submits raw or archive files (e.g., ZIP) to Zscaler's Sandbox for analysis. You can submit up to 100 files per day and it supports all file types that are currently supported by Sandbox. The resource also allows the submissions of raw or archive files to the Zscaler service for out-of-band file inspection to generate real-time verdicts for known and unknown files. It leverages capabilities such as Malware Prevention, Advanced Threat Prevention, Sandbox cloud effect, AI/ML-driven file analysis, and integrated third-party threat intelligence feeds to inspect files and classify them as benign or malicious instantaneously.

⚠️ **WARNING 1:**: Zscaler Cloud Sandbox is a subscription service and requires additional license. To learn more, contact Zscaler Support or your local account team.

⚠️ **WARNING 2:**: The ZIA Terraform provider requires both the `ZIA_CLOUD` and `ZIA_SANDBOX_TOKEN` in order to authenticate to the Zscaler Cloud Sandbox environment. For details on how obtain the API Token visit the Zscaler help portal [About Sandbox API Token](https://help.zscaler.com/zia/about-sandbox-api-token)

**Note 1**: After files are sent for analysis, you must use GET /sandbox/report/{md5Hash} in order to retrieve the verdict. You can get the Sandbox report 10 minutes after a file is sent for analysis.

**Note 2**: All file types that are currently supported by the Malware Protection policy and Advanced Threat Protection policy are supported for inspection, and each file is limited to a size of 400 MB.

## Example Usage - Submit raw or archive files

```hcl
# Submit raw EXE file to Zscaler Sandbox
locals {
  files = toset([
    "zs-test-pe-file.exe"
  ])
}

resource "zia_sandbox_file_submission" "this" {
  for_each = local.files
  file_path     = each.key
  submission_method = "submit"
  force = true
}
```

## Example Usage - Submits raw or archive for out-of-band file inspection

```hcl
# Submit raw EXE file to Zscaler Sandbox
locals {
  files = toset([
    "zs-test-pe-file.exe"
  ])
}

resource "zia_sandbox_file_submission" "this" {
  for_each = local.files
  file_path     = each.key
  submission_method = "discan"
  force = true
}
```

## Attributes Reference

### Required

* `file_path` - (Required) The path where the raw or archive files for submission are located.

* `submission_method` - (Required) The submission method to be used. Supportedd values are: `submit` and `discan`
  * `submit` - Submits raw or archive files (e.g., ZIP) to Sandbox for analysis.
  * `discan` - Submits raw or archive files (e.g., ZIP) to the Zscaler service for out-of-band file inspection to generate real-time verdicts for known and unknown files.

* `force` - (Optional) Submit file to sandbox even if found malicious during AV scan and a verdict already exists. Supported values are `true` or `false`
  * `true` - If a verdict already exists for the file, you can use set force = `true` to make the sandbox reanalyze the file.
  * `false` - By default, files are scanned by Zscaler antivirus (AV) and submitted directly to the sandbox in order to obtain a verdict.

**Note 3**: After files are sent for analysis, you must use the following data source `zia_sandbox_report` in order to retrieve the verdict. You can get the Sandbox report ~10 minutes after a file is sent for analysis.
