---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: dc_exclusions"
description: |-
    Official documentation https://help.zscaler.com/zia/excluding-data-center-based-traffic-forwarding-method
    API documentation https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/dcExclusions-get
    Retrieves the list of Zscaler data centers (DCs) that are currently excluded from service to your organization based on configured exclusions
---

# zia_dc_exclusions (Data Source)

* [Official documentation](https://help.zscaler.com/zia/excluding-data-center-based-traffic-forwarding-method)
* [API documentation](https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/dcExclusions-get)

Use the **zia_dc_exclusions** data source to retrieve the list of Zscaler data centers (DCs) that are currently excluded from service to your organization based on configured exclusions in the ZIA Admin Portal

~> NOTE: This an Early Access feature.

## Example Usage - Retrieve All DC Exclusions

```hcl
data "zia_dc_exclusions" "all" {
}
```

## Example Usage - Filter by Name

```hcl
data "zia_dc_exclusions" "this" {
    name = "ADL"
}
```

## Example Usage - Filter by Resource ID

```hcl
data "zia_dc_exclusions" "example" {
  id = zia_dc_exclusions.example.id
}
```

## Example Usage - Filter by Datacenter ID

```hcl
data "zia_dc_exclusions" "example" {
  datacenter_id = 1221
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `id` - (String) Filter exclusions by the resource id (datacenter ID as string). Use to look up a specific exclusion by its zia_dc_exclusions resource id.
* `datacenter_id` - (Integer) Filter exclusions by datacenter ID (dcid).
* `name` - (String) Filter exclusions by datacenter name (case-insensitive partial match on dcName.name).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The data source identifier. Set to the datacenter ID when exactly one exclusion is returned; otherwise `id-<datacenter_id>`, `name-<name>`, or `all`.
* `name` - (String, Optional + Computed) Filter value when provided; when filtering by `datacenter_id` and exactly one exclusion is returned, set to that exclusion's datacenter name.

### exclusions

* `exclusions` - (List) List of DC exclusion entries.
  * `id` - (String) The exclusion identifier (datacenter ID as string). Matches the zia_dc_exclusions resource id.
  * `dc_id` - (Integer) Datacenter ID (dcid) for the exclusion.
  * `expired` - (Boolean) Whether the exclusion has expired.
  * `start_time` - (Integer) Unix timestamp when the exclusion window starts.
  * `end_time` - (Integer) Unix timestamp when the exclusion window ends.
  * `description` - (String) Description of the DC exclusion.
  * `dc_name_id` - (Integer) Datacenter ID from the dcName reference.
  * `dc_name` - (String) Datacenter name from the dcName reference.
