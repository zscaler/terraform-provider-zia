---
subcategory: "Device Groups"
layout: "zia"
page_title: "ZIA: devices"
description: |-
  Official documentation https://help.zscaler.com/zia/device-groups#/deviceGroups-get
  API documentation https://help.zscaler.com/zia/device-groups#/deviceGroups-get
  Get information about ZIA devices.
---

# zia_devices (Data Source)

* [Official documentation](https://help.zscaler.com/zia/device-groups#/deviceGroups-get)
* [API documentation](https://help.zscaler.com/zia/device-groups#/deviceGroups-get)

Use the **zia_devices** data source to get information about a device in the Zscaler Internet Access cloud or via the API. This data source can then be associated with resources such as: URL Filtering Rules

## Example Usage

```hcl
# ZIA Admin User Data Source
data "zia_devices" "device"{
    name = "administrator"
}
```

### Example Usage - With JMESPath Search

```hcl
# Use JMESPath to pre-filter devices by OS type
data "zia_devices" "windows_device" {
    name   = "workstation-01"
    search = "[?osType == 'WINDOWS_OS']"
}
```

```hcl
# Filter devices by model
data "zia_devices" "device" {
    name   = "administrator"
    search = "[?contains(deviceModel, 'MacBook')]"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the devices to be exported.
* `id` - (Optional) The unique identifer for the devices.
* `search` - (Optional) A [JMESPath](https://jmespath.org/) expression to filter results client-side after all pages have been retrieved from the API. The expression is applied to the list of devices before name or ID matching. This is useful in large environments to narrow down the candidate set. Field names in expressions must use the API's camelCase names (e.g., `name`, `osType`, `osVersion`, `deviceModel`, `ownerName`).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

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
