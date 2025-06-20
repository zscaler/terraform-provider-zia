---
subcategory: "Advanced Threat Protection"
layout: "zscaler"
page_title: "ZIA: atp_security_exceptions"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-advanced-threat-protection-policy
  API documentation https://help.zscaler.com/zia/advanced-threat-protection-policy#/cyberThreatProtection/advancedThreatSettings-put
  Updates security exceptions for the ATP policy.
---

# zia_atp_security_exceptions (Resource)

* [Official documentation](https://help.zscaler.com/zia/configuring-advanced-threat-protection-policy)
* [API documentation](https://help.zscaler.com/zia/advanced-threat-protection-policy#/)

The **zia_atp_security_exceptions** resource alows you to updates security exceptions for the ATP policy. To learn more see [Advanced Threat Protection](https://help.zscaler.com/unified/configuring-security-exceptions-advanced-threat-protection-policy)

## Example Usage

```hcl
resource "zia_atp_security_exceptions" "this" {
    bypass_urls = [
        "site1.example.com",
        "site2.example.com",
        "site3.example.com",
    ]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `bypass_urls` - (Set of String) Allowlist URLs that are not inspected by the ATP policy.

**NOTE** To delete all current urls, submit a empty list.

### Optional

There are no optional parameters supported by this resource.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_atp_security_exceptions** can be imported by using `all_urls` as the import ID.

For example:

```shell
terraform import zia_atp_security_exceptions.this all_urls
```
