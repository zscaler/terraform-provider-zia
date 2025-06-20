---
subcategory: "NAT Control Policy"
layout: "zscaler"
page_title: "ZIA: nat_control_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/about-nat-control
  API documentation https://help.zscaler.com/zia/nat-control-policy#/dnatRules-get
  Creates and manages ZIA NAT Control Rules.
---

# zia_nat_control_rules (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-nat-control)
* [API documentation](https://help.zscaler.com/zia/nat-control-policy#/dnatRules-get)

The **zia_nat_control_rules** resource allows the creation and management of NAT Control rules in the Zscaler Internet Access.

## Example Usage

```hcl
resource "zia_nat_control_rules" "this" {
    name = "DNAT_02"
    description = "DNAT_02"
    order=1
    rank=7
    state = "ENABLED"
    redirect_port="2000"
    redirect_ip="1.1.1.1"
    src_ips=["192.168.100.0/24", "192.168.200.1"]
    dest_addresses=["3.217.228.0-3.217.231.255", "3.235.112.0-3.235.119.255", "35.80.88.0-35.80.95.255", "server1.acme.com", "*.acme.com"]
    dest_countries=["BR", "CA", "GB"]
  departments {
    id = [8061246]
  }
  dest_ip_groups {
    id = [-4]
  }
  dest_ipv6_groups {
    id = [-5]
  }
  src_ip_groups {
    id = [18448894]
  }
  src_ipv6_groups {
    id = [-3]
  }
  time_windows {
    id = [485]
  }
  nw_services {
    id = [462370, 17472664]
  }
  locations {
    id = [256000852, -3]
  }
  location_groups {
    id = [8061257, 8061256]
  }
  labels {
    id = [1416803]
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

- `name` - (Required) Name of the Firewall Filtering policy rule
- `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.

### Optional

- `description` - (string) - Additional information about the forwarding rule
- `enable_full_logging` (Boolean)
- `redirect_port` - (string) -  Port to which the traffic is redirected to when the DNAT rule is triggered. If not set, no redirection is done to the specific port.
- `redirect_ip` - (string) - IP address to which the traffic is redirected to when the DNAT rule is triggered. If not set, no redirection is done to the specific IP address.
- `redirect_fqdn` - (string) - FQDN to which the traffic is redirected to when the DNAT rule is triggered. This is mutually exclusive to redirect IP.
- `redirect_fqdn` - (string) - FQDN to which the traffic is redirected to when the DNAT rule is triggered. This is mutually exclusive to redirect IP.

`Who, Where and When` supports the following attributes:

- `locations` - (Block List, Max: 1) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `location_groups` - (Block List, Max: 1) You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `users` - (Block List, Max: 1) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `groups` - (Block List, Max: 1) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `departments` - (Block List, Max: 1) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `time_windows` - (Block List, Max: 1) You can manually select up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

`network services` supports the following attributes:

- `nw_service_groups` - (Block List, Max: 1) Any number of predefined or custom network service groups to which the rule applies.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `nw_services`- (Block List, Max: 1) When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

`network applications` supports the following attributes:

- `nw_application_groups` - (Block List, Max: 1) Any number of application groups that you want to control with this rule. The service provides predefined applications that you can group, but not modify
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `nw_applications` - (Block List, Max: 1) When not used it applies the rule to all applications. The service provides predefined applications, which you can group, but not modify.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

`source ip addresses` supports the following attributes:

- `src_ip_groups` - (Block List, Max: 1) Any number of source IP address groups that you want to control with this rule.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `src_ipv6_groups` - (Block List, Max: 1) Any number of source IPv6 address groups that you want to control with this rule.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `src_ips` - (List of String) You can enter individual IP addresses, subnets, or address ranges.

`destinations` supports the following attributes:

- `dest_addresses`** - (List of String) -  IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges. If adding multiple items, hit Enter after each entry.
- `dest_countries`** - (List of String) Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries. Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

- `dest_ip_categories`** - (List of String) IP address categories of destination for which the DNAT rule is applicable. If not set, the rule is not restricted to specific destination IP categories.

- `dest_ip_groups`** - (Block List, Max: 1) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `dest_ipv6_groups`** - (Block List, Max: 1) Any number of destination IPv6 address groups that you want to control with this rule.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

- `labels` (Block List, Max: 1) Labels that are applicable to the rule.
      - `id` - (List of Integer) Identifier that uniquely identifies an entity

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_nat_control_rules** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import zia_nat_control_rules.example <rule_id>
```

or

```shell
terraform import zia_nat_control_rules.example <rule_name>
```
