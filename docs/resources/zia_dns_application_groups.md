---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: dns_application_groups"
description: |-
  Official documentation https://help.zscaler.com/zia/about-dns-application-groups
  API documentation https://help.zscaler.com/legacy-apis/firewall-policies#/dnsApplicationGroups/{groupId}-get
  Creates and manages DNS Application Groups.
---

# zia_dns_application_groups (Resource)

* [Official documentation](https://help.zscaler.com/zia/firewall-policies#/ipDestinationGroups-post)
* [API documentation](https://help.zscaler.com/zia/firewall-policies#/ipDestinationGroups-post)

The **zia_dns_application_groups** resource allows the creation and management of DNS Application Groups in the Zscaler Internet Access. This resource can then be associated with a ZIA DNS Control rule `zia_firewall_dns_rule`.

## Example Usage

```hcl
resource "zia_dns_application_groups" "this" {
	name        = "DNSAppGroup01"
	description = "DNSAppGroup01"
    dns_applications = ["RACKSPACE", "MCAFEE", "SOPHOS", "BRIGHTSPCACE", "CSWG", "INTECH", "SECURESERVER", "FRANCETELECOM"]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Destination IP group name
* `description` (Optional) Additional information about the destination IP group
* `applications` - (List) The list of cloud applications to which must be added to the group. FTo retrieve the list of cloud applications, use the data source: `zia_cloud_applications`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_dns_application_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import zia_dns_application_groups.example <group_id>
```

or

```shell
terraform import zia_dns_application_groups.example <group_name>
```
