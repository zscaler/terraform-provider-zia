---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: datacenters"
description: |-
    Official documentation https://help.zscaler.com/zia/understanding-subclouds
    API documentation https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/datacenters-get
    Retrieves the list of Zscaler data centers (DCs) that can be excluded from service to your organization.
---

# zia_datacenters (Data Source)

* [Official documentation](https://help.zscaler.com/zia/understanding-subclouds)
* [API documentation](https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/datacenters-get)

Use the **zia_datacenters** data source to retrieve the list of Zscaler data centers (DCs) that can be excluded from service to your organization.

~> NOTE: This an Early Access feature.

## Example Usage - Retrieve All Datacenters

```hcl
data "zia_datacenters" "all" {
}
```

## Example Usage - Filter by Name

```hcl
data "zia_datacenters" "filtered" {
    name = "CA Client Node DC"
}
```

## Example Usage - Filter by Multiple Criteria

```hcl
data "zia_datacenters" "filtered" {
    city            = "San Jose"
    dc_provider     = "Zscaler Internal"
    gov_only        = false
    third_party_cloud = false
    virtual         = false
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `datacenter_id` - (Integer) Filter datacenters by ID. When exactly one result is returned, this is set to that datacenter's ID.
* `name` - (String) Filter datacenters by name (case-insensitive partial match). When exactly one result is returned, this is set to that datacenter's name.
* `city` - (String) Filter datacenters by city (case-insensitive partial match). When exactly one result is returned, this is set to that datacenter's city.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String, Computed) When exactly one datacenter matches, the numeric ID as string (e.g. `"515"`). Otherwise a filter-based identifier (e.g. `"name-SJC4"`, `"id-515"`, `"all"`).
* `datacenter_id` - (Integer, Computed) When filtering by ID or when exactly one datacenter matches, the datacenter's numeric ID. Prefer this over `datacenters[0].id` when you expect a single result (e.g. `name = "SJC4"`).
* `name` - (String, Computed) Set from the filter or from the single matching datacenter's name.
* `city` - (String, Computed) Set from the filter or from the single matching datacenter's city.

### Datacenters

* `datacenters` - (List) List of datacenters matching the filter criteria.
  * `id` - (Integer) Unique identifier for the datacenter.
  * `name` - (String) Zscaler data center name.
  * `provider` - (String) Provider of the datacenter.
  * `city` - (String) City where the datacenter is located.
  * `timezone` - (String) Timezone of the datacenter.
  * `lat` - (Integer) Latitude coordinate (legacy field).
  * `longi` - (Integer) Longitude coordinate (legacy field).
  * `latitude` - (Float) Latitude coordinate.
  * `longitude` - (Float) Longitude coordinate.
  * `gov_only` - (Boolean) Whether this is a government-only datacenter.
  * `third_party_cloud` - (Boolean) Whether this is a third-party cloud datacenter.
  * `upload_bandwidth` - (Integer) Upload bandwidth in bytes per second.
  * `download_bandwidth` - (Integer) Download bandwidth in bytes per second.
  * `owned_by_customer` - (Boolean) Whether the datacenter is owned by the customer.
  * `managed_bcp` - (Boolean) Whether the datacenter is managed by BCP.
  * `dont_publish` - (Boolean) Whether the datacenter should not be published.
  * `dont_provision` - (Boolean) Whether the datacenter should not be provisioned.
  * `not_ready_for_use` - (Boolean) Whether the datacenter is not ready for use.
  * `for_future_use` - (Boolean) Whether the datacenter is reserved for future use.
  * `regional_surcharge` - (Boolean) Whether there is a regional surcharge for this datacenter.
  * `create_time` - (Integer) Timestamp when the datacenter was created.
  * `last_modified_time` - (Integer) Timestamp when the datacenter was last modified.
  * `virtual` - (Boolean) Whether this is a virtual datacenter.
