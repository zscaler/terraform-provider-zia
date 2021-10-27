---
subcategory: "Location Management"
layout: "zia"
page_title: "ZIA: gre_internal_ip_range_list"
description: |-
  Gets the next available GRE tunnel internal IP address ranges
  
---

# zia_gre_internal_ip_range_list (Data Source)

The **zia_gre_internal_ip_range_list** - data source retrieves details about available Zscaler GRE tunnel internal ip ranges for the purposes of GRE tunnel creation in the Zscaler Internet Access when the `ip_unnumbered` parameter is set to `false`

## Example Usage

```hcl
# Retrieve GRE available Internal IP Ranges
# By default it will return the first 10 available internal ip ranges
data "zia_gre_internal_ip_range_list" "example"{
}

output "zia_gre_internal_ip_range_list_example"{
    value = data.zia_gre_internal_ip_range_list.example
}
```

## Argument Reference

The following arguments are supported:

* `required_count` - (Optional)

### Read-Only

`list`

* `end_ip_address` - (String) Starting IP address in the range
* `start_ip_address` - (String) Ending IP address in the range
