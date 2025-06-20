---
subcategory: "Security Policy Settings"
layout: "zscaler"
page_title: "ZIA: zia_security_settings"
description: |-
  Official documentation https://help.zscaler.com/zia/security-policy-settings#/security-put
  API documentation https://help.zscaler.com/zia/security-policy-settings#/security-put
  Gets a list of URLs that are on the allow and denylist.
---

# zia_security_settings (Data Source)

* [Official documentation](https://help.zscaler.com/zia/security-policy-settings#/security-put)
* [API documentation](https://help.zscaler.com/zia/security-policy-settings#/security-put)

Use the **zia_security_settings** data source to get a list of URLs that were added to the allow and denylist under the Advanced Threat Protection policy in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# ZIA Security Policy Settings Data Source
data "zia_security_settings" "example"{}
```

## Argument Reference

This data source can be executed without the need of additional parameters.

## Attribute Reference

N/A
