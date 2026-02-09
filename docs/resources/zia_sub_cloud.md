---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: sub_cloud"
description: |-
    Official documentation https://help.zscaler.com/zia/understanding-subclouds
    API documentation https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/subclouds-get
    Update a subcloud information available in the Zscaler Internet Access cloud.
---

# zia_sub_cloud (Resource)

* [Official documentation](https://help.zscaler.com/zia/understanding-subclouds)
* [API documentation](https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/subclouds-get)

Use the **zia_sub_cloud** resource to update the subcloud and excluded data centers based on the specified ID.

~> NOTE: This an Early Access feature.

## Example Usage

```hcl
data "zia_sub_cloud" "lookup" {
    name = "BIZDevZSThree01"
}

data "zia_datacenters" "this" {
    name = "YVR1"
}

data "zia_datacenters" "this1" {
    name = "SEA1"
}

resource "zia_sub_cloud" "this" {
    cloud_id = data.zia_sub_cloud.lookup.id
    name     = "BIZDevZSThree01"

    # Using Unix timestamps
    exclusions {
        datacenter {
            id   = data.zia_datacenters.this.datacenters[0].id
            name = data.zia_datacenters.this.datacenters[0].name
        }
        country    = "CANADA"
        end_time   = 1770422399
    }

    # Using human-readable UTC date/time (same as UI "Data Center Disabled Until")
    exclusions {
        datacenter {
            id   = data.zia_datacenters.this1.datacenters[0].id
            name = data.zia_datacenters.this1.datacenters[0].name
        }
        country      = "UNITED_STATES"
        end_time_utc = "02/19/2026 11:59:00 pm"
    }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `cloud_id` - (Integer) Unique identifier for the subcloud as an integer.

### Optional

* `name` - (String) Subcloud name. This attribute is read-only and cannot be updated via the API. If provided, it will be ignored and the name from the API response will be used instead.
* `exclusions` - (List) List of data centers excluded from the subcloud.
  * `datacenter` - (List, Required, MaxItems: 1) The excluded datacenter reference.
    * `id` - (Integer, Required) Unique identifier for the datacenter.
    * `name` - (String, Optional) Datacenter name.
    * `country` - (String, Optional) Country where the datacenter is located.

  **NOTE:**  Use the datasource `zia_datacenters` to retrieve the list of available datacenter IDs

  * `country` - (String, Required) Country where the excluded data center is located.
  * `end_time` - (Integer, Optional) Exclusion end time (Unix timestamp). Either `end_time` or `end_time_utc` must be set.
  * `end_time_utc` - (String, Optional) Data center disabled until (UTC). Format: `MM/DD/YYYY HH:MM:SS am/pm`. If set, overrides `end_time`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The unique identifier for the subcloud as a string (Terraform resource ID).
* `cloud_id` - (Integer) Unique identifier for the subcloud as an integer.
* `name` - (String) Subcloud name. This attribute is read-only and cannot be updated via the API.

### Exclusions

* `exclusions` - (List) List of data centers excluded from the subcloud.
  * `datacenter` - (List) The excluded datacenter reference.
    * `id` - (Integer) Unique identifier for the datacenter.
    * `name` - (String) Datacenter name.
    * `country` - (String) Country where the datacenter is located.

  **NOTE:**  Use the datasource `zia_datacenters` to retrieve the list of available datacenter IDs

  * `country` - (String) Country where the excluded data center is located.
  * `end_time` - (Integer, Optional) Exclusion end time (Unix timestamp). Either `end_time` or `end_time_utc` must be set.
  * `end_time_utc` - (String, Optional) Data center disabled until (UTC). Format: `MM/DD/YYYY HH:MM:SS am/pm`. If set, overrides `end_time`.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_sub_cloud** can be imported by using `<SUB_CLOUD ID>` or `<SUB_CLOUD NAME>` as the import ID.

For example:

```shell
terraform import zia_sub_cloud.example <sub_cloud_id>
```

or

```shell
terraform import zia_sub_cloud.example <sub_cloud_name>
```
