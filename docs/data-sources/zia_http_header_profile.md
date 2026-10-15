---
subcategory: "HTTP Header Control"
layout: "zscaler"
page_title: "ZIA: http_header_profile"
description: |-
  Official documentation https://help.zscaler.com/zia/about-http-header-profile
  API documentation https://help.zscaler.com/legacy-apis/http-header-control#/httpHeaderProfile-get
  Gets information about a ZIA HTTP header insertion profile.
---

# zia_http_header_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-http-header-profile)
* [API documentation](https://help.zscaler.com/legacy-apis/http-header-control#/httpHeaderProfile-get)

Use the **zia_http_header_profile** data source to get information about an HTTP header insertion profile in the Zscaler Internet Access cloud. The profile can then be referenced in URL filtering rules.

## Example Usage

```hcl
# Retrieve an HTTP header profile by name
data "zia_http_header_profile" "this" {
  name = "Profile01"
}
```

```hcl
# Retrieve an HTTP header profile by ID
data "zia_http_header_profile" "this" {
  id = 1234567
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, Number) The unique identifier for the HTTP header profile.
* `name` - (Optional, String) The HTTP header profile name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Additional information about the HTTP header profile.
* `slot_id` - (Number) The slot ID assigned to the HTTP header profile.
* `deleted` - (Boolean) Indicates whether the HTTP header profile is deleted.
* `profile_ready_for_use` - (Boolean) Indicates whether the HTTP header profile is ready for use.
* `http_header_profile_criteria` - (List) The list of matching criteria evaluated by the HTTP header profile.
  * `id` - (Number) Identifier that uniquely identifies the criteria entry.
  * `header` - (String) The header evaluated by the criteria. Returned values: `USERAGENT`, `REFERER`, `ORIGIN`
  * `operator` - (String) The operator applied to the header criteria. Returned values: `UAVERSIONGT`, `UAVERSIONLT`, `UAVERSIONEQ`, `UAVERSIONNEQ`, `UAVERSIONANY`
  * `user_agent_bitmap` - (String) The user agent (browser) evaluated by the criteria. Returned values: `OPERA`, `FIREFOX`, `MSIE`, `MSEDGE`, `CHROME`, `SAFARI`, `OTHER`, `MSCHREDGE`, `BRAVE`
  * `user_agent` - (String) The custom user agent string evaluated by the criteria. Set when `user_agent_bitmap` is `OTHER`.
  * `user_agent_version` - (String) The user agent version evaluated by the criteria.
  * `category_bitmap` - (List of String) The URL categories evaluated by the criteria.
  * `cloud_app_bitmap` - (List of String) The cloud applications evaluated by the criteria.
