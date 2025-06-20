---
subcategory: "Location Management"
layout: "zscaler"
page_title: "ZIA: location_groups"
description: |-
  Official documentation https://help.zscaler.com/zia/about-locations
  API documentation https://help.zscaler.com/zia/location-management#/locations-get
  Get information about Location Groups.

---

# zia_location_groups (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-locations)
* [API documentation](https://help.zscaler.com/zia/location-management#/locations-get)

Use the **zia_location_groups** data source to get information about a location group option available in the Zscaler Internet Access.

```hcl
# Retrieve ZIA Location Group
data "zia_location_groups" "example"{
    name = "Corporate User Traffic Group"
}
```

```hcl
# Retrieve ZIA Location Group
data "zia_location_groups" "example"{
    name = "Guest Wifi Group"
}
```

```hcl
# Retrieve ZIA Location Group
data "zia_location_groups" "example"{
    name = "IoT Traffic Group"
}
```

```hcl
# Retrieve ZIA Location Group
data "zia_location_groups" "example"{
    name = "Server Traffic Group"
}
```

```hcl
# Retrieve ZIA Location Group
data "zia_location_groups" "example"{
    name = "Server Traffic Group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Location group name
* `id` - (Optional) Unique identifier for the location group

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `comments` - (String) Additional information about the location group
* `deleted` - (Boolean) Indicates the location group was deleted
* `group_type` - (String) The location group's type (i.e., Static or Dynamic)
* `id` - (Number) The ID of this resource.
* `last_mod_time` - (Number) Automatically populated with the current time, after a successful `POST` or `PUT` request.
* `predefined` - (Boolean)

* `dynamic_location_group_criteria` - (Block Set) Dynamic location group information.
  * `name` - (Block List)
    * `match_string` - (String) String value to be matched or partially matched
    * `match_type` - (String) Operator that performs match action

  * `countries` - (List of String) One or more countries from a predefined set

  * `city` - (Block List)
    * `match_string` - (String) String value to be matched or partially matched
    * `match_type` - (String) Operator that performs match action

  * `managed_by` - (Block List)
    * `id` -(Number)
    * `name` -(String)
    * `extensions` -(Map of String)

  * `enforce_authentication` - (Boolean) Enforce Authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.
  * `enforce_aup` - (Boolean) Enable AUP. When set to true, AUP is enabled for the location.
  * `enforce_firewall_control` - (Boolean) Enable Firewall. When set to true, Firewall is enabled for the location.
  * `enable_xff_forwarding` - (Boolean) Enable `XFF` Forwarding. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
  * `enable_caution` - (Boolean) Enable Caution. When set to true, a caution notifcation is enabled for the location.
  * `enable_bandwidth_control` - (Boolean) Enable Bandwidth Control. When set to true, Bandwidth Control is enabled for the location.
  * `profiles` - (List of String) One or more location profiles from a predefined set

* `locations` - (List of Object) The Name-ID pairs of the locations that are assigned to the static location group. This is ignored if the groupType is Dynamic.
  * `id` -(Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `last_mod_user` - (List of Object)
  * `id` -(Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)
  * `comments` - (List of Object)

* `last_mod_time` - (List of Object) Automatically populated with the current time, after a successful POST or PUT request.
* `predefined` - (Boolean)
