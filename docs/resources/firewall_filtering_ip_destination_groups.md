---
subcategory: "Firewall Policies"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): firewall_filtering_destination_groups"
sidebar_current: "docs-resource-zia-firewall-filtering-destination-group"
description: |-
  Creates and manages ZIA Cloud firewall IP destination groups.
---

# Resource: zia_firewall_filtering_destination_groups

The **zia_firewall_filtering_destination_groups** resource allows the creation and management of ZIA Cloud Firewall IP destination groups in the Zscaler Internet Access. This resource can then be associated with a ZIA cloud firewall filtering rule.

```hcl
# IP Destination Group of Type DSTN_FQDN
resource "zia_firewall_filtering_destination_groups" "dstn_fqdn" {
  name        = "Example Destination FQDN"
  description = "Example Destination FQDN"
  type        = "DSTN_FQDN"
  addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
}
```

```hcl
# IP Destination Group of Type DSTN_IP
resource "zia_firewall_filtering_destination_groups" "example_dstn_ip" {
  name        = "Example Destination IP"
  description = "Example Destination IP"
  type        = "DSTN_IP"
  addresses = ["3.217.228.0-3.217.231.255",
    "3.235.112.0-3.235.119.255",
    "52.23.61.0-52.23.62.25",
    "35.80.88.0-35.80.95.255"]
}
```

```hcl
# IP Destination Group of Type DSTN_DOMAIN
resource "zia_firewall_filtering_destination_groups" "example_dstn_domain" {
  name          = "Example Destination Domain"
  description   = "Example Destination Domain"
  type          = "DSTN_DOMAIN"
  addresses     = ["acme.com", "acme1.com"]
}
```

```hcl
# IP Destination Group of Type DSTN_OTHER
resource "zia_firewall_filtering_destination_groups" "example_dstn_other" {
  name          = "Example Destination Other"
  description   = "Example Destination Other"
  type          = "DSTN_OTHER"
  countries     = ["COUNTRY_CA"]
  categories    = ["COUNTRY_CA"]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Destination IP group name
* `type` - (Required) Destination IP group type (i.e., the group can contain destination IP addresses or FQDNs). The supported values are:
  * `DSTN_IP`
  * `DSTN_FQDN`
  * `DSTN_DOMAIN`
  * `DSTN_OTHER`
* `addresses` (Required) Destination IP addresses, domains or FQDNs within the group

### Optional

* `countries` (Optional) Destination IP address counties. You can identify destinations based on the location of a server.
* `description` (Optional) Additional information about the destination IP group
* `ip_categories` (Optional) Destination IP address URL categories. You can identify destinations based on the URL category of the domain.
