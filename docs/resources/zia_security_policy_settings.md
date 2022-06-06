---
subcategory: "Security Policy Settings"
layout: "zscaler"
page_title: "ZIA: security_policy_settings"
description: |-
  Add or Remove a URL to and from the allow and deny list in the advanced threat protection.
---

# Resource: zia_security_policy_settings

The **zia_security_policy_settings** resource alows you to add or remove a URL to the allow and denylist under the Advanced Threat Protection policy in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Add URLs to Whitelist
resource "zia_security_settings" "test"{
    whitelist_urls = []
}
```

```hcl
# Add URLs to Blacklist
resource "zia_security_settings" "test"{
    blacklist_urls = []
}
```

```hcl
# Add URLs to both Whitelist and Blacklist
resource "zia_security_settings" "test"{
    whitelist_urls = []
    blacklist_urls = []
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `whitelist_urls` - (Required) Allowlist URLs whose contents will not be scanned. Allows up to 255 URLs
* `blacklist_urls` - (Optional) URLs on the denylist for your organization. Allow up to 25000 URLs.
