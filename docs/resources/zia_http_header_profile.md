---
subcategory: "HTTP Header Control"
layout: "zscaler"
page_title: "ZIA: http_header_profile"
description: |-
  Official documentation https://help.zscaler.com/zia/about-http-header-profile
  API documentation https://help.zscaler.com/legacy-apis/http-header-control#/httpHeaderProfile-get
  Creates and manages ZIA HTTP header insertion profiles.
---

# zia_http_header_profile (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-http-header-profile)
* [API documentation](https://help.zscaler.com/legacy-apis/http-header-control#/httpHeaderProfile-get)

The **zia_http_header_profile** resource allows the creation and management of HTTP header insertion profiles in the Zscaler Internet Access cloud or via the API. The profile defines the matching criteria evaluated against traffic.

## Example Usage

```hcl
# ZIA HTTP Header Profile Resource
resource "zia_http_header_profile" "this" {
  name        = "Profile01"
  description = "Example header profile"

  http_header_profile_criteria {
    header   = "ORIGIN"
    cloud_app_bitmap = ["CHATGPT_AI"]
    category_bitmap = ["GENERAL_AI_ML", "AI_ML_APPS"]
  }
  http_header_profile_criteria {
    header   = "REFERER"
    cloud_app_bitmap = ["CHATGPT_AI"]
    category_bitmap = ["GENERAL_AI_ML", "AI_ML_APPS"]
  }
  http_header_profile_criteria {
    header   = "USERAGENT"
    user_agent_bitmap = "FIREFOX"
    operator = "UAVERSIONEQ"
    user_agent_version = "123.0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The HTTP header profile name.

### Optional

* `description` - (String) Additional information about the HTTP header profile.
* `slot_id` - (Number) The slot ID assigned to the HTTP header profile.
* `profile_ready_for_use` - (Boolean) Indicates whether the HTTP header profile is ready for use.
* `http_header_profile_criteria` - (Block List) The list of matching criteria evaluated by the HTTP header profile.
  * `header` - (String) The header evaluated by the criteria. Supported Values: `USERAGENT`, `REFERER`, `ORIGIN`
  * `operator` - (String) The operator applied to the header criteria. Supported Values: `UAVERSIONGT`, `UAVERSIONLT`, `UAVERSIONEQ`, `UAVERSIONNEQ`, `UAVERSIONANY`
  * `user_agent` - (String) The user agent evaluated by the criteria.
  * `user_agent_bitmap` - (String) The user agent bitmap evaluated by the criteria. Supportee Values: `OPERA`, `FIREFOX`, `MSIE`, `MSEDGE`, `CHROME`, `SAFARI`, `OTHER`, `MSCHREDGE`, `BRAVE`
  * `user_agent_version` - (String) The user agent version evaluated by the criteria.
  * `id` - (Number) Identifier that uniquely identifies the criteria entry.
  * `category_bitmap` - (List of String) The URL category bitmap evaluated by the criteria.
  * `cloud_app_bitmap` - (List of String) The cloud application bitmap evaluated by the criteria.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_http_header_profile** can be imported by using `<PROFILE_ID>` or `<PROFILE_NAME>` as the import ID.

For example:

```shell
terraform import zia_http_header_profile.example <profile_id>
```

or

```shell
terraform import zia_http_header_profile.example <profile_name>
```
