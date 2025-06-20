---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: gre_internal_ip_range_list"
description: |-
  Official documentation https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels/availableInternalIpRanges-get
  API documentation https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels/availableInternalIpRanges-get
  Gets the next available GRE tunnel internal IP address ranges.
---

# zia_gre_internal_ip_range_list (Data Source)

* [Official documentation](https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels/availableInternalIpRanges-get)
* [API documentation](https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels/availableInternalIpRanges-get)

Use the **zia_gre_internal_ip_range_list** data source to get information about the next available GRE tunnel internal ip ranges for the purposes of GRE tunnel creation in the Zscaler Internet Access when the `ip_unnumbered` parameter is set to `false`

## Example Usage

```hcl
# Retrieve GRE available Internal IP Ranges
# By default it will return the first 10 available internal ip ranges
data "zia_gre_internal_ip_range_list" "example"{
}
```

```hcl
# Retrieve GRE available Internal IP Ranges
# By using the `required_count` parameter it will return the indicated number of IP ranges.
data "zia_gre_internal_ip_range_list" "example"{
  required_count = 20
}
```

## Argument Reference

The following arguments are supported:

* `required_count` - (Required)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `end_ip_address` - (String) Starting IP address in the range
* `start_ip_address` - (String) Ending IP address in the range
