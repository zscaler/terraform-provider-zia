---
subcategory: "Firewall Filtering - Network Services"
layout: "zia"
page_title: "ZIA: firewall_filtering_network_service"
description: |-
        Gets a list of all network services. The search parameters find matching values within the name or description attributes.
---


# zia_firewall_filtering_network_service (Resource)

The **zia_firewall_filtering_network_service** - data source retrieves information about all network services.

```hcl
data "zia_firewall_filtering_network_service" "example" {
  name = "ICMP_ANY"
}

output "zia_firewall_filtering_network_service" {
  value = data.zia_firewall_filtering_network_service.example
}
```

## Argument Reference

The following arguments are supported:

* `name` -(String)

### Read-Only

* `description` - (String)

* `is_name_l10n_tag` -(Boolean
* `tag` -(String)
* `type` -(String)

`src_tcp_ports` -(Block List)

* `start` - (Set of Number)
* `end` - (Set of Number)

`src_udp_ports` -(Block List)

* `start` - (Set of Number)
* `end` - (Set of Number)

`dest_tcp_ports` - (Block List)

* `start` - (Set of Number)
* `end` - (Set of Number)

`dest_udp_ports` -(Block List)

* `start` - (Set of Number)
* `end` - (Set of Number)
