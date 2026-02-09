---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: sub_cloud"
description: |-
    Official documentation https://help.zscaler.com/zia/understanding-subclouds
    API documentation https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/subclouds-get
    Gets information about a subcloud available in the Zscaler Internet Access cloud.
---

# zia_sub_cloud (Data Source)

* [Official documentation](https://help.zscaler.com/zia/understanding-subclouds)
* [API documentation](https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/subclouds-get)

Use the **zia_sub_cloud** data source to get information about a subcloud available in the Zscaler Internet Access cloud. Summary of the subcloud associated with an organization. This also represents the data centers that are excluded and associated with the subcloud.

~> NOTE: This an Early Access feature.

## Example Usage - Retrieve by Name

```hcl
data "zia_sub_cloud" "this" {
    name = "SubCloud01"
}
```

## Example Usage - Retrieve by ID

```hcl
data "zia_sub_cloud" "this" {
    id = 1254674585
}
```

## Argument Reference

The following arguments are supported:

### Required

At least one of the following must be provided:

* `id` - (Integer) Unique identifier for the subcloud. Used to look up a single subcloud when provided.
* `name` - (String) Subcloud name. Used to search for a subcloud by name when provided.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Integer) Unique identifier for the subcloud
* `name` - (String) Subcloud name

### Data Centers (dcs)

* `dcs` - (List) List of data centers associated with the subcloud
  * `id` - (Integer) Data center ID for the country
  * `name` - (String) Data center name
  * `country` - (String) Country name. Enum with 245 predefined country values.

### Exclusions

* `exclusions` - (List) List of data centers that are excluded from the subcloud. Information about data center excluded from the subcloud.
  * `datacenter` - (List) The data center associated with the subcloud. An immutable reference to an entity that mainly consists of `id` and `name`.
    * `id` - (Integer) A unique identifier for an entity
    * `name` - (String) The configured name of the entity (read-only)
    * `extensions` - (Map of String) Extension attributes
  * `country` - (String) Country where the data center is located. Enum with 245 predefined country values.
  * `expired` - (Boolean) The subcloud data center exclusion is disabled
  * `disabled_by_ops` - (Boolean) If set to true, this data center exclusion is disabled by Zscaler CloudOps
  * `create_time` - (Integer) Timestamp when the data center exclusion was created
  * `start_time` - (Integer) Timestamp when the data center exclusion was started
  * `end_time` - (Integer) Timestamp when the data center exclusion was stopped
  * `last_modified_time` - (Integer) Timestamp when the data center exclusion entry was last modified
  * `last_modified_user` - (List) Last user that modified the data center. An immutable reference to an entity that mainly consists of `id` and `name`.
    * `id` - (Integer) A unique identifier for an entity
    * `name` - (String) The configured name of the entity (read-only)
    * `extensions` - (Map of String) Extension attributes
