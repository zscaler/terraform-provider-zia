---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: dns_application_groups"
description: |-
  Official documentation https://help.zscaler.com/zia/about-dns-application-groups
  API documentation https://help.zscaler.com/legacy-apis/firewall-policies#/dnsApplicationGroups/{groupId}-get
  Get information about DNS Application Groups.
---


# zia_dns_application_groups (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-dns-application-groups)
* [API documentation](https://help.zscaler.com/legacy-apis/firewall-policies#/dnsApplicationGroups/{groupId}-get)

Use the **zia_dns_application_groups** data source to get information about dns application groups available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA DNS Control rule `zia_firewall_dns_rule`.

## Example Usage - DNS Application Group by Name

```hcl
data "zia_dns_application_groups" "example" {
    name = "DNSGroup01"
}
```

## Example Usage - DNS Application Group by ID

```hcl
data "zia_dns_application_groups" "example" {
    id = 12113780
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ip source group to be exported.
* `id` - (Optional) The ID of the ip source group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of this resource.
* `description` - (String)
* `ip_addresses` - (List of String)
