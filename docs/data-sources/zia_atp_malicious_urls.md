---
subcategory: "Advanced Threat Protection"
layout: "zscaler"
page_title: "ZIA: atp_malicious_urls"
description: |-
  Retrieves the malicious URLs added to the denylist in the Advanced Threat Protection (ATP) policy
---

# Data Source: zia_atp_malicious_urls

Use the **zia_atp_malicious_urls** data source to Retrieves the malicious URLs added to the denylist in the Advanced Threat Protection (ATP) policy. To learn more see [Advanced Threat Protection](https://help.zscaler.com/unified/configuring-security-exceptions-advanced-threat-protection-policy)

## Example Usage

```hcl
data "zia_atp_malicious_urls" "this" {}
```

## Argument Reference

This data source can be executed without the need of additional parameters.

## Attribute Reference

N/A
