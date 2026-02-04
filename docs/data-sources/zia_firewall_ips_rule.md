---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA): firewall_ips_rule"
description: |-
  Official documentation https://help.zscaler.com/zia/ips-control-policy#/firewallIpsRules-get
  API documentation https://help.zscaler.com/zia/configuring-ips-control-policy
  Get information about firewall IPS Control policy rule.
---

# zia_firewall_ips_rule (Data Source)

* [Official documentation](https://help.zscaler.com/zia/ips-control-policy#/firewallIpsRules-get)
* [API documentation](https://help.zscaler.com/zia/configuring-ips-control-policy)

Use the **zia_firewall_ips_rule** data source to get information about a cloud firewall IPS rule available in the Zscaler Internet Access.

## Example Usage

```hcl
# ZIA Firewall IPS Rule by name
data "zia_firewall_ips_rule" "this" {
    name = "Default Cloud IPS Rule"
}
```

```hcl
# ZIA Firewall IPS Rule by ID
data "zia_firewall_ips_rule" "this" {
    id = "12365478"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Firewall Filtering policy rule
* `id` - (Optional) Unique identifier for the Firewall Filtering policy rule

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Enter additional notes or information. The description cannot exceed 10,240 characters.
* `order` - (Integer) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.
* `state` - (String) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule.
* `action` - (String) The action configured for the rule that must take place if the traffic matches the rule criteria, such as allowing or blocking the traffic or bypassing the rule. The following actions are accepted: `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BYPASS_IPS`
* `rank` - (Integer) By default, the admin ranking is disabled. To use this feature, you must enable admin rank. The default value is `7`.
* `enable_full_logging` - (Integer) A Boolean value that indicates whether full logging is enabled. A true value indicates that full logging is enabled, whereas a false value indicates that aggregate logging is enabled.
* `capture_pcap` - (Boolean) Value that indicates whether packet capture (PCAP) is enabled or not
* `predefined` - (Boolean) A Boolean field that indicates that the rule is predefined by using a true value
* `default_rule` - (Boolean) Value that indicates whether the rule is the Default Cloud IPS Rule or not
* `eun_enabled` - (Boolean) A Boolean value that indicates whether Web EUN is enabled for the rule
* `eun_template_id` - (Integer) The EUN template ID associated with the rule

`Devices`

* `devices` - (List of Objects) Devices to which the rule applies. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `device_groups` - (List of Objects) Device groups to which the rule applies. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

`Who, Where and When` supports the following attributes:

* `locations` - (List of Objects) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `location_groups` - (List of Objects)You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `users` - (List of Objects) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `groups` - (List of Objects) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `departments` - (List of Objects) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `time_windows` - (List of Objects) You can manually select up to `1` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

`network services` supports the following attributes:

* `nw_service_groups` - (List of Objects) Any number of predefined or custom network service groups to which the rule applies.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `nw_services`- (List of Objects) When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

`source ip addresses` supports the following attributes:

* `source_countries` (Set of String) The countries of origin of traffic for which the rule is applicable. If not set, the rule is not restricted to specific source countries.
    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `src_ip_groups` - (List of Objects)Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `src_ipv6_groups` - (List of Objects) Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `src_ips` - (Set of String) Source IP addresses or FQDNs to which the rule applies. If not set, the rule is not restricted to a specific source IP address. Each IP entry can be a single IP address, CIDR (e.g., 10.10.33.0/24), or an IP range (e.g., 10.10.33.1-10.10.33.10).

`destinations` supports the following attributes:

* `dest_addresses` (Set of String) Destination IP addresses or FQDNs to which the rule applies. If not set, the rule is not restricted to a specific destination IP address. Each IP entry can be a single IP address, CIDR (e.g., 10.10.33.0/24), or an IP range (e.g., 10.10.33.1-10.10.33.10).

* `dest_countries` (Set of String) Identify destinations based on the location of a server, select Any to apply the rule to all countries or select the countries to which you want to control traffic.
    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `res_categories` (Set of String) URL categories associated with resolved IP addresses to which the rule applies. If not set, the rule is not restricted to a specific URL category.

* `dest_ip_categories` (Set of String)  identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
* `dest_ip_groups`** - (List of Objects) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `threat_categories` (List of Objects) Advanced threat categories to which the rule applies
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `labels` (List of Objects) Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `zpa_app_segments` (List of Objects) The ZPA application segments to which the rule applies
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
