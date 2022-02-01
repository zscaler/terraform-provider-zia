---
subcategory: "Device Groups"
layout: "zia"
page_title: "ZIA: device_groups"
description: |-
  Retrieve ZIA device group details.
  
---
# zia_device_groups (Data Source)

The **zia_device_groups** - data source provides details about a specific device group in the Zscaler Internet Access cloud or via the API. This data source can then be associated with resources such as: URL Filtering Rules

## Example Usage

```hcl
# ZIA Admin User Data Source
data "zia_device_groups" "ios"{
    name = "IOS"
}

output "zia_device_groups_ios"{
    value = data.zia_device_groups.ios
}

data "zia_device_groups" "android"{
    name = "Android"
}

output "zia_device_groups_android"{
    value = data.zia_device_groups.android
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the device group to be exported.
* `id` - (Optional) The unique identifer for the device group.

### Read-Only

* `id` - (String) The unique identifer for the device group.
* `name` - (String) The device group name.
* `group_type` - (String) The device group type. i.e ``ZCC_OS``, ``NON_ZCC``, ``CBI``
* `description` - (String) The device group's description.
* `os_type` - (String) The operating system (OS).
* `predefined` - (Boolean) Indicates whether this is a predefined device group. If this value is set to true, the group is predefined.
* `device_names` - (String) The names of devices that belong to the device group. The device names are comma-separated.
* `device_count` - (int) The number of devices within the group.
