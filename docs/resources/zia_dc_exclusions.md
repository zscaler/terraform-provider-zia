---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: dc_exclusions"
description: |-
    Official documentation https://help.zscaler.com/zia/excluding-data-center-based-traffic-forwarding-method
    API documentation https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/dcExclusions-get
    Adds a data center (DC) exclusion to disable the tunnels terminating at a virtual IP address of a Zscaler DC
---

# zia_dc_exclusions (Resource)

* [Official documentation](https://help.zscaler.com/zia/excluding-data-center-based-traffic-forwarding-method)
* [API documentation](https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/dcExclusions-get)

Use the **zia_dc_exclusions** Resource to add a data center (DC) exclusion to disable the tunnels terminating at a virtual IP address of a Zscaler DC, triggering a failover from primary to secondary tunnels in the event of service disruptions, Zscaler Trust Portal incidents, disasters, etc. You can configure to exclude a specific DC based on the traffic forwarding method for a designated time period.

~> NOTE: This an Early Access feature.

## Example Usage

```hcl
data "zia_datacenters" "this" {
  name = "SJC4"
}

# Using Unix timestamps
resource "zia_dc_exclusions" "this" {
  datacenter_id = data.zia_datacenters.this.datacenter_id
  start_time   = 1770422399
  end_time     = 1770508799
  description  = "Optional description"
}

# Using human-readable UTC date/time (same as zia_sub_cloud exclusions)
resource "zia_dc_exclusions" "utc" {
  datacenter_id  = data.zia_datacenters.this.datacenter_id
  start_time_utc = "02/05/2026 12:00 am"
  end_time_utc   = "02/19/2026 11:59 pm"
  description    = "Optional description"
}
```

## Argument Reference

The following arguments are supported:

### Required (on create)

* `datacenter_id` - (Integer, Optional/Computed) Datacenter ID (dcid) to exclude. Required when creating; can be omitted when importing by numeric ID. Prefer `data.zia_datacenters.this.datacenter_id` when filtering by name with a single result. Force new if changed.
* `start_time` - (Integer, Optional/Computed) Unix timestamp when the exclusion window starts. Either `start_time` or `start_time_utc` must be set.
* `start_time_utc` - (String, Optional/Computed) Exclusion window start (UTC). Format: `MM/DD/YYYY HH:MM am/pm`. If set, overrides `start_time`.
* `end_time` - (Integer, Optional/Computed) Unix timestamp when the exclusion window ends. Either `end_time` or `end_time_utc` must be set.
* `end_time_utc` - (String, Optional/Computed) Exclusion window end (UTC). Format: `MM/DD/YYYY HH:MM am/pm`. If set, overrides `end_time`.

### Optional

* `description` - (String) Description of the DC exclusion.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) Terraform state ID; matches the API datacenter ID (numeric, as string).
* `start_time_utc` - (String) Exclusion window start in UTC, formatted as `MM/DD/YYYY HH:MM am/pm` (computed from API).
* `end_time_utc` - (String) Exclusion window end in UTC, formatted as `MM/DD/YYYY HH:MM am/pm` (computed from API).
* `expired` - (Boolean) Whether the exclusion has expired (read-only from API).
