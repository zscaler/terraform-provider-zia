---
subcategory: "Firewall Filtering Network Application"
layout: "zia"
page_title: "ZIA: firewall_filtering_network_application"
description: |-
  Retrieve ZIA firewall rule network application.
  
---

# zia_firewall_filtering_network_application (Data Source)

The **zia_firewall_filtering_network_application** data source provides details about a specific network application available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network application rule.

## Example Usage

```hcl
# ZIA Network Application Groups
data "zia_firewall_filtering_network_application" "apns"{
    id = "APNS"
    locale="en-US"
}

output "zia_firewall_filtering_network_application_apns"{
    value = data.zia_firewall_filtering_network_application.apns
}
```

```hcl
data "zia_firewall_filtering_network_application" "dict"{
    id = "DICT"
}

output "zia_firewall_filtering_network_application_dict"{
    value = data.zia_firewall_filtering_network_application.dict
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The name of the ip source group to be exported.
* `locale` - (Optional)

### Read-Only

* `deprecated` - (Boolean)
* `description` - (String)
* `parent_category` - (String)
