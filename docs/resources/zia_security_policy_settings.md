---
subcategory: "Security Policy Settings"
layout: "zscaler"
page_title: "ZIA: zia_security_settings"
description: |-
  Add or Remove URLs for whitelisting or blacklisting in the Malware and Advanced Threat Protection.
---

# Resource: zia_security_settings

The **zia_security_settings** resource alows you to add or remove a URL to the allow and denylist under the Advanced Threat Protection policy in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Add URLs to ZIA Whitelist - Malware Protection
resource "zia_security_settings" "this" {
  whitelist_urls = [
    "resource5.acme.net",
    "resource6.acme.net",
    "resource7.acme.net",
    "resource8.acme.net",
  ]
}
```

```hcl
# Add URLs to ZIA Blacklist - Advanced Threat Protection
resource "zia_security_settings" "this" {
  blacklist_urls = [
    "resource1.acme.net",
    "resource2.acme.net",
    "resource3.acme.net",
    "resource4.acme.net",
  ]
}
```

```hcl
# Add URLs to both Whitelist and Blacklist
# Advanced Threat Protection & Malware Protection
resource "zia_security_settings" "this" {
  whitelist_urls = [
    "resource5.acme.net",
    "resource6.acme.net",
    "resource7.acme.net",
    "resource8.acme.net",
  ]
  blacklist_urls = [
    "resource1.acme.net",
    "resource2.acme.net",
    "resource3.acme.net",
    "resource4.acme.net",
  ]
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `whitelist_urls` - (Required) Allowlist URLs whose contents will not be scanned. Allows up to `255` URLs
* `blacklist_urls` - (Optional) URLs on the denylist for your organization. Allow up to `25000` URLs.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_security_settings** can be imported by using `all_urls` as the import ID.

For example:

```shell
terraform import zia_security_settings.example all_urls
```
