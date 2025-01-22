---
subcategory: "Advanced Threat Protection"
layout: "zscaler"
page_title: "ZIA: atp_malicious_urls"
description: |-
  Updates the malicious URLs added to the denylist in ATP policy
---

# Resource: zia_atp_malicious_urls

The **zia_atp_malicious_urls** resource alows you to Updates the malicious URLs added to the denylist in ATP policy. To learn more see [Advanced Threat Protection](https://help.zscaler.com/unified/configuring-security-exceptions-advanced-threat-protection-policy)

## Example Usage

```hcl
resource "zia_atp_malicious_urls" "this" {
    malicious_urls = [
        "test1.malicious.com",
        "test2.malicious.com",
        "test3.malicious.com",
    ]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `malicious_urls` - (Set of String) List of malicious URLs that are blocked by the ATP policy. At least 1 URL must be specified.

### Optional

There are no optional parameters supported by this resource.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_atp_malicious_urls** can be imported by using `all_urls` as the import ID.

For example:

```shell
terraform import zia_atp_malicious_urls.this all_urls
```
