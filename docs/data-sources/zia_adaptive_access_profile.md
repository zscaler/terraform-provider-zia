---
subcategory: "Adaptive Access Profiles"
layout: "zscaler"
page_title: "ZIA: adaptive_access_profile"
description: |-
  Official documentation https://help.zscaler.com/zia/about-adaptive-access
  Get information about ZIA Adaptive Access profiles.
---

# zia_adaptive_access_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-adaptive-access)

Use the **zia_adaptive_access_profile** data source to get information about Adaptive Access profiles configured in the Zscaler Internet Access cloud or via the API.

The data source supports two lookup modes:

* Provide `name` to look up a single Adaptive Access profile by name.
* Provide `iam_aap_ids` and/or `org_id` to return the matching profile rules in the `profiles` block.

## Example Usage

### Look up a single profile by name

```hcl
data "zia_adaptive_access_profile" "this" {
  name = "Profile01"
}
```

### Filter the profile rules

```hcl
data "zia_adaptive_access_profile" "by_ids" {
  iam_aap_ids = ["1234", "5678"]
}
```

```hcl
data "zia_adaptive_access_profile" "by_org" {
  org_id = 123456
}
```

## Argument Reference

The following arguments are supported. At least one of `name`, `iam_aap_ids`, or `org_id` must be provided.

* `name` - (Optional) The Adaptive Access profile name. Used to look up a single profile.
* `iam_aap_ids` - (Optional) Filters the profile rules by one or more Adaptive Access profile IDs. Setting this attribute returns the matching profile rules in the `profiles` block.
* `org_id` - (Optional) Filters the profile rules by organization ID. Setting this attribute returns the matching profile rules in the `profiles` block.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` (Number) The Adaptive Access profile ID
* `type` (String) The Adaptive Access profile type
* `aap_index` (Number) The Adaptive Access profile index
* `iam_aap_id` (String) The Adaptive Access profile ID that is used by the API for policy configuration. This field allows you to specify which Adaptive Access profiles are applied in the access policy criteria.
* `deleted` (Boolean) A Boolean value that indicates whether the Adaptive Access profile is deleted

The `profiles` block is populated when `iam_aap_ids` or `org_id` is set. Each entry exports:

* `id` (Number) The Adaptive Access profile ID
* `name` (String) The Adaptive Access profile name
* `type` (String) The Adaptive Access profile type
* `aap_index` (Number) The Adaptive Access profile index
* `iam_aap_id` (String) The Adaptive Access profile ID that is used by the API for policy configuration. This field allows you to specify which Adaptive Access profiles are applied in the access policy criteria.
* `deleted` (Boolean) A Boolean value that indicates whether the Adaptive Access profile is deleted
