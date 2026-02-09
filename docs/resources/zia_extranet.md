---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: extranet"
description: |-
    Official documentation https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/extranet-get
    API documentation https://help.zscaler.com/zia/understanding-extranet-application-support
    Retrieves the list of extranets configured for the organization.
---

# zia_extranet (Resource)

* [Official documentation](https://help.zscaler.com/legacy-apis/traffic-forwarding-0#/extranet-get)
* [API documentation](https://help.zscaler.com/zia/understanding-extranet-application-support)

Use the **zia_extranet** resource Adds a new extranet for the organization in the Zscaler Internet Access cloud. Extranets are configured as part of Zscaler Extranet Application Support which allows an organization to connect its internal network with another organization’s network (e.g., partners, third-party vendors, etc.) that does not use the Zscaler service. Extranet Application Support enables Zscaler-managed organization users to securely access extranet resources through an IPSec VPN tunnel established between the Zscaler data center and the external organization’s data center, without requiring additional hardware or software installations.

~> NOTE: This an Early Access feature.

## Example Usage - Retrieve by Name

```hcl
resource "zia_extranet" "this" {
    name        = "Extranet01"
    description = "Extranet01"

    extranet_dns_list {
        name                 = "DNS01"
        primary_dns_server   = "8.8.8.8"
        secondary_dns_server = "4.4.4.4"
        use_as_default       = true
    }

    extranet_dns_list {
        name                 = "DNS02"
        primary_dns_server   = "192.168.1.1"
        secondary_dns_server = "192.168.1.2"
        use_as_default       = false
    }

    extranet_ip_pool_list {
        name           = "TFS01"
        ip_start       = "10.0.0.1"
        ip_end         = "10.0.0.21"
        use_as_default = true
    }

    extranet_ip_pool_list {
        name           = "TFS02"
        ip_start       = "10.0.0.22"
        ip_end         = "10.0.0.43"
        use_as_default = false
    }
}
```

## Argument Reference

The following arguments are supported:

### Required

At least one of the following must be provided:

* `id` - (Integer) Unique identifier for the extranet. Used to look up a single extranet when provided.
* `name` - (String) Extranet name. Used to search for an extranet by name when provided.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Integer) The unique identifier for the extranet.
* `name` - (String) The name of the extranet.
* `description` - (String) The description of the extranet.

### extranet_dns_list

* `extranet_dns_list` - (List) Information about the DNS servers specified for the extranet.
  * `name` - (String) The name of the DNS server.
  * `primary_dns_server` - (String) The IP address of the primary DNS server.
  * `secondary_dns_server` - (String) The IP address of the secondary DNS server.
  * `use_as_default` - (Boolean) Whether the DNS servers specified in the extranet are the designated default servers.

### extranet_ip_pool_list

* `extranet_ip_pool_list` - (List) Information about the traffic selectors (IP pools) specified for the extranet.
  * `name` - (String) The name of the IP pool.
  * `ip_start` - (String) The starting IP address of the pool.
  * `ip_end` - (String) The ending IP address of the pool.
  * `use_as_default` - (Boolean) Whether this IP pool is the designated default.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_extranet** can be imported by using `<EXTRANET ID>` or `<EXTRANET NAME>` as the import ID.

For example:

```shell
terraform import zia_extranet.example <extranet_id>
```

or

```shell
terraform import zia_extranet.example <extranet_name>
```
