---
subcategory: "Firewall Filtering - Network Services"
layout: "zia"
page_title: "ZIA: firewall_filtering_network_service"
description: |-
      Adds a new network service.
---

# zia_firewall_filtering_network_service (Resource)

The **zia_firewall_filtering_network_service** - data source retrieves information about all network services.

```hcl
resource "zia_firewall_filtering_network_service" "example" {
  name        = "example"
  description = "example"
  src_tcp_ports {
  }
  dest_tcp_ports {
    start = 5000
    end = 5005
  }
    dest_tcp_ports {
  }
  type = "CUSTOM"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)
* `description` - (Optional)

* `is_name_l10n_tag` - (Optional
* `tag` - (Optional)
* `type` - (Optional) Supported values: `STANDARD`, `PREDEFINED`, `CUSTOM`

`src_tcp_ports` - (Block List)

* `start` - (Optional)
* `end` - (Optional)

`src_udp_ports` -(Block List)

* `start` - (Set of Number)
* `end` - (Set of Number)

`dest_tcp_ports` - (Block List)

* `start` - (Required)
* `end` - (Required)

`dest_udp_ports` -(Block List)

* `start` - (Set of Number)
* `end` - (Set of Number)
