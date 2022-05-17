---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_network_application"
description: |-
  Get information about ZIA firewall rule network application.

---

# Data Source zia_firewall_filtering_network_application

Use the **zia_firewall_filtering_network_application** data source to get information about a network application available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network application rule.

## Example Usage

```hcl
# ZIA Network Application Groups
data "zia_firewall_filtering_network_application" "apns"{
    id = "APNS"
    locale="en-US"
}
```

```hcl
data "zia_firewall_filtering_network_application" "dict"{
    id = "DICT"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The name of the ip source group to be exported.
* `locale` - (Optional)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `deprecated` - (Boolean)
* `description` - (String)
* `parent_category` - (String)
