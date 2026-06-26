---
subcategory: "HTTP Header Control"
layout: "zscaler"
page_title: "ZIA: http_header_action_profile"
description: |-
  Official documentation https://help.zscaler.com/zia/about-http-header-profile
  API documentation https://help.zscaler.com/legacy-apis/http-header-control#/httpHeaderActionProfile-get
  Creates and manages ZIA HTTP header action profiles.
---

# zia_http_header_action_profile (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-http-header-profile)
* [API documentation](https://help.zscaler.com/legacy-apis/http-header-control#/httpHeaderActionProfile-get)

The **zia_http_header_action_profile** resource allows the creation and management of HTTP header action profiles in the Zscaler Internet Access cloud or via the API. The profile defines the header key/value pairs applied to matching traffic.

## Example Usage

```hcl
# ZIA HTTP Header Action Profile Resource
resource "zia_http_header_action_profile" "this" {
  name = "ActionProfile01"
  description = "Example header action profile"

  http_header_action_profile_keys {
    key = "X-Forwarded-For"
    value = "10.0.0.1"
  }

  http_header_action_profile_keys {
    key = "User-Agent"
    value = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
  }

  http_header_action_profile_keys {
    key = "Accept"
    value = "application/json"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The HTTP header action profile name.

### Optional

* `description` - (String) Additional information about the HTTP header action profile.
* `slot_id` - (Number) The slot ID assigned to the HTTP header action profile.
* `profile_ready_for_use` - (Boolean) Indicates whether the HTTP header action profile is ready for use.
* `http_header_action_profile_keys` - (Block List) The list of header key/value pairs applied by the action profile.
  * `key` - (Required, String) The header key.
  * `value` - (Required, String) The header value.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_http_header_action_profile** can be imported by using `<PROFILE_ID>` or `<PROFILE_NAME>` as the import ID.

For example:

```shell
terraform import zia_http_header_action_profile.example <profile_id>
```

or

```shell
terraform import zia_http_header_action_profile.example <profile_name>
```
