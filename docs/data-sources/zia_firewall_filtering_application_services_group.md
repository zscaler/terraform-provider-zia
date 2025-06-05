---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_network_application_group"
description: |-
  Get information about firewall rule network application groups.

---


# Data Source - zia_firewall_filtering_application_services_group

Use the **zia_firewall_filtering_application_services_group** data source to get information about a network application group available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network application rule.

## Example Usage

```hcl
# ZIA Network Application Groups
data "zia_firewall_filtering_application_services_group" "example" {
    name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ip source group to be exported.
* `id` - (Optional) The ID of the ip source group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String)
* `network_applications` - (List of String)
* `id` - The ID of this resource.
