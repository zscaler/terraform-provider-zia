---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_ips_rule"
description: |-
  Official documentation https://help.zscaler.com/zia/ips-control-policy#/firewallIpsRules-get
  API documentation https://help.zscaler.com/zia/configuring-ips-control-policy
  Creates and manages ZIA Cloud firewall IPS rule.
---

# zia_firewall_ips_rule (Resource)

* [Official documentation](https://help.zscaler.com/zia/ips-control-policy#/firewallIpsRules-get)
* [API documentation](https://help.zscaler.com/zia/configuring-ips-control-policy)

The **zia_firewall_ips_rule** resource allows the creation and management of ZIA Cloud Firewall IPS rules in the Zscaler Internet Access.

**NOTE 1** Zscaler Cloud Firewall contain default and predefined rules which are placed in their respective orders. These rules `CANNOT` be deleted and NOT all attributes are suppported. When configuring your rules make sure that the `order` attributue value consider these pre-existing rules so that Terraform can place the new rules in the correct position, and drifts can be avoided. i.e If there are 2 pre-existing rules and intend to manage those rules via Terraform, you must first import those rules into the state and start the ordering accordingly. However, if DO NOT intend to manage predefined rules via Terraform, the provider will reorder the rules automatically while ignoring the order of pre-existing rules, as the API is responsible for moving these rules to their respective positions as API calls are made.

## Example Usage

```hcl
data "zia_firewall_filtering_network_service" "zscaler_proxy_nw_services" {
    name = "ZSCALER_PROXY_NW_SERVICES"
}

data "zia_department_management" "engineering" {
 name = "Engineering"
}

data "zia_group_management" "normal_internet" {
    name = "Normal_Internet"
}

data "zia_firewall_filtering_time_window" "work_hours" {
    name = "Work hours"
}

resource "zia_firewall_ips_rule" "example" {
    name = "Example_IPS_Rule01"
    description = "Example_IPS_Rule01"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    dest_countries = ["CA", "US"]
    source_countries = ["CA", "US"]
    threat_categories {
        id = [ 66 ]
    }
    nw_services {
        id = [ data.zia_firewall_filtering_network_service.zscaler_proxy_nw_services.id ]
    }
    departments {
        id = [ data.zia_department_management.engineering.id ]
    }
    groups {
        id = [ data.zia_group_management.normal_internet.id ]
    }
    time_windows {
        id = [ data.zia_firewall_filtering_time_window.work_hours.id ]
    }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the Firewall IPS policy rule
* `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Enter additional notes or information. The description cannot exceed 10,240 characters.
* `order` - (Integer) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.

* `state` - (Optional) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule. Supported Values: `ENABLED`, `DISABLED`

* `action` - (String) The action configured for the rule that must take place if the traffic matches the rule criteria, such as allowing or blocking the traffic or bypassing the rule. The following actions are accepted: `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BYPASS_IPS`

* `rank` - (Integer) By default, the admin ranking is disabled. To use this feature, you must enable admin rank in UI first. The default value is `7`. Visit to learn more [About Admin Rank](https://help.zscaler.com/zia/about-admin-rank)

* `enable_full_logging` - (Integer) A Boolean value that indicates whether full logging is enabled. A true value indicates that full logging is enabled, whereas a false value indicates that aggregate logging is enabled.
* `capture_pcap` - (Boolean) Value that indicates whether packet capture (PCAP) is enabled or not
* `predefined` - (Boolean) A Boolean field that indicates that the rule is predefined by using a true value
* `default_rule` - (Boolean) Value that indicates whether the rule is the Default Cloud IPS Rule or not

`Devices`

* `devices` - (List of Objects) Devices to which the rule applies. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `device_groups` - (List of Objects) Device groups to which the rule applies. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (Integer) Identifier that uniquely identifies an entity

`Who, Where and When` supports the following attributes:

* `locations` - (List of Objects) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `location_groups` - (List of Objects)You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `users` - (List of Objects) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `groups` - (List of Objects) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `departments` - (List of Objects) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `time_windows` - (List of Objects) You can manually select up to `1` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (Integer) Identifier that uniquely identifies an entity

`network services` supports the following attributes:

* `nw_service_groups` - (List of Objects) Any number of predefined or custom network service groups to which the rule applies.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `nw_services`- (List of Objects) When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.
      - `id` - (Integer) Identifier that uniquely identifies an entity

`source ip addresses` supports the following attributes:

* `source_countries` (Set of String) The countries of origin of traffic for which the rule is applicable. If not set, the rule is not restricted to specific source countries.
    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `src_ip_groups` - (List of Objects)Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `src_ipv6_groups` - (List of Objects) Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `src_ips` - (Set of String) Source IP addresses or FQDNs to which the rule applies. If not set, the rule is not restricted to a specific source IP address. Each IP entry can be a single IP address, CIDR (e.g., 10.10.33.0/24), or an IP range (e.g., 10.10.33.1-10.10.33.10).

`destinations` supports the following attributes:

* `dest_addresses` (Set of String) Destination IP addresses or FQDNs to which the rule applies. If not set, the rule is not restricted to a specific destination IP address. Each IP entry can be a single IP address, CIDR (e.g., 10.10.33.0/24), or an IP range (e.g., 10.10.33.1-10.10.33.10).

* `dest_countries` (Set of String) Identify destinations based on the location of a server, select Any to apply the rule to all countries or select the countries to which you want to control traffic.
    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `res_categories` (Set of String) URL categories associated with resolved IP addresses to which the rule applies. If not set, the rule is not restricted to a specific URL category.

* `dest_ip_categories` (Set of String)  identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
* `dest_ip_groups`** - (List of Objects) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `threat_categories` (List of Objects) Advanced threat categories to which the rule applies
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `labels` (List of Objects) Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `zpa_app_segments` (List of Objects) The ZPA application segments to which the rule applies
      - `id` - (Integer) Identifier that uniquely identifies an entity
