---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_destination_groups"
description: |-
  Get information about IP destination groups.

---

# Data Source: zia_firewall_filtering_destination_groups

Use the **zia_firewall_filtering_destination_groups** data source to get information about IP destination groups option available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering rule.

## Example Usage

```hcl
# ZIA Destination Groups
data "zia_firewall_filtering_destination_groups" "example" {
    name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the destination group to be exported.
* `id` - (Optional) The ID of the destination group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Additional information about the destination IP group
* `addresses` - (List of String) Destination IP addresses within the group
* `countries` - (List of String) Destination IP address counties. You can identify destinations based on the location of a server.
* `ip_categories` - (List of String) Destination IP address URL categories. You can identify destinations based on the URL category of the domain.
* `type` - (String) Destination IP group type (i.e., the group can contain destination IP addresses or FQDNs)
