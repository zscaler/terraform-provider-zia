---
subcategory: "Advanced Threat Protection"
layout: "zscaler"
page_title: "ZIA: atp_malicious_urls"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-advanced-threat-protection-policy
  API documentation https://help.zscaler.com/zia/advanced-threat-protection-policy#/cyberThreatProtection/advancedThreatSettings-put
  Retrieves the malicious URLs added to the denylist in the Advanced Threat Protection (ATP) policy
---

# zia_atp_malicious_urls (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-advanced-threat-protection-policy)
* [API documentation](https://help.zscaler.com/zia/advanced-threat-protection-policy#/)

Use the **zia_atp_malicious_urls** data source to Retrieves the malicious URLs added to the denylist in the Advanced Threat Protection (ATP) policy. To learn more see [Advanced Threat Protection](https://help.zscaler.com/unified/configuring-security-exceptions-advanced-threat-protection-policy)

## Example Usage

```hcl
data "zia_atp_malicious_urls" "this" {}
```

## Argument Reference

This data source can be executed without the need of additional parameters.

## Attribute Reference

N/A
