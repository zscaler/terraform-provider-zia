---
subcategory: "NAT Control Policy"
layout: "zscaler"
page_title: "ZIA): nat_control_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/about-nat-control
  API documentation https://help.zscaler.com/zia/nat-control-policy#/dnatRules-get
  Retrieves a list of all configured and predefined DNAT Control policies.
---

# zia_nat_control_rules (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-nat-control)
* [API documentation](https://help.zscaler.com/zia/nat-control-policy#/dnatRules-get)

Use the **zia_nat_control_rules** data source to get information about a NAT Control rule available in the Zscaler Internet Access.

## Example Usage - By Name

```hcl

data "zia_nat_control_rules" "this" {
  name = "DNAT_01"
}
```

## Example Usage - By ID

```hcl

data "zia_nat_control_rules" "this" {
  id = 154658
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the forwarding rule.
* `id` - (Optional) A unique identifier assigned to the forwarding rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) - Additional information about the forwarding rule
* `order` - (string) - The order of execution for the forwarding rule order.
* `redirect_port` - (string) -  Port to which the traffic is redirected to when the DNAT rule is triggered. If not set, no redirection is done to the specific port.
* `redirect_ip` - (string) - IP address to which the traffic is redirected to when the DNAT rule is triggered. If not set, no redirection is done to the specific IP address.
* `redirect_fqdn` - (string) - FQDN to which the traffic is redirected to when the DNAT rule is triggered. This is mutually exclusive to redirect IP.
* `redirect_fqdn` - (string) - FQDN to which the traffic is redirected to when the DNAT rule is triggered. This is mutually exclusive to redirect IP.
* `trusted_resolver_rule` - (boolean) - Set to true in the predefined rule for Zscaler Trusted DNS Resolver

`Who, Where and When` supports the following attributes:

* `locations` - (Block List, Max: 1) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `location_groups` - (Block List, Max: 1) You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `users` - (Block List, Max: 1) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `groups` - (Block List, Max: 1) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `departments` - (Block List, Max: 1) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `time_windows` - (Block List, Max: 1) You can manually select up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

`network services` supports the following attributes:

* `nw_service_groups` - (Block List, Max: 1) Any number of predefined or custom network service groups to which the rule applies.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `nw_services`- (Block List, Max: 1) When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

`network applications` supports the following attributes:

* `nw_application_groups` - (Block List, Max: 1) Any number of application groups that you want to control with this rule. The service provides predefined applications that you can group, but not modify
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `nw_applications` - (Block List, Max: 1) When not used it applies the rule to all applications. The service provides predefined applications, which you can group, but not modify.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

`source ip addresses` supports the following attributes:

* `src_ip_groups` - (Block List, Max: 1) Any number of source IP address groups that you want to control with this rule.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `src_ipv6_groups` - (Block List, Max: 1) Any number of source IPv6 address groups that you want to control with this rule.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `src_ips` - (List of String) You can enter individual IP addresses, subnets, or address ranges.

`destinations` supports the following attributes:

* `dest_addresses`** - (List of String) -  IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges. If adding multiple items, hit Enter after each entry.
* `dest_countries`** - (List of String) Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.

* `dest_ip_categories`** - (List of String) IP address categories of destination for which the DNAT rule is applicable. If not set, the rule is not restricted to specific destination IP categories.

* `dest_ip_groups`** - (Block List, Max: 1) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `dest_ipv6_groups`** - (Block List, Max: 1) Any number of destination IPv6 address groups that you want to control with this rule.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `labels` (Block List, Max: 1) Labels that are applicable to the rule.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `Other Exported Arguments`
  * `id` - (int) The ID of this resource.
  * `last_modified_time` - (Number)
  * `access_control` - (String) Access privilege of this rule based on the admin's Role Based Authorization (RBA) state.
  * `enable_full_logging` (Boolean)
  * `default_rule` - (Boolean) If set to true, the default rule is applied
  * `predefined` - (Boolean) If set to true, a predefined rule is applied
