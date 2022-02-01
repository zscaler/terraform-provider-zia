---
subcategory: "Device"
layout: "zia"
page_title: "ZIA: devices"
description: |-
  Retrieve ZIA devices details.
  
---
# zia_devices (Data Source)

The **zia_devices** - data source provides details about a specific device in the Zscaler Internet Access cloud or via the API. This data source can then be associated with resources such as: URL Filtering Rules

## Example Usage

```hcl
# ZIA Admin User Data Source
data "zia_devices" "device"{
    name = "administrator"
}

output "zia_devices"{
    value = data.zia_devices.device
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the devices to be exported.
* `id` - (Optional) The unique identifer for the devices.

### Read-Only

* `id` - (String) The unique identifer for the device group.
* `name` - (String) The device name.
* `device_group_type` - (String) The device group type. i.e ``ZCC_OS``, ``NON_ZCC``, ``CBI``
* `device_model` - (String) The device model.
* `description` - (String) The device group's description.
* `os_type` - (String) The operating system (OS). ``ANY``, ``OTHER_OS``, ``IOS``, ``ANDROID_OS``, ``WINDOWS_OS``, ``MAC_OS``, ``LINUX``
* `os_version` - (String) The operating system version.
* `description` - (String) The device's description.
* `owner_user_id` - (int) The unique identifier of the device owner (i.e., user).
* `owner_name` - (String) The device owner's user name.
