---
subcategory: "Forwarding Control Policy"
layout: "zscaler"
page_title: "ZIA): forwarding_control_rule"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-forwarding-policy
  API documentation https://help.zscaler.com/zia/forwarding-control-policy#/forwardingRules-get
  Get information about forwarding control rule.
---

# zia_forwarding_control_rule (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-forwarding-policy)
* [API documentation](https://help.zscaler.com/zia/forwarding-control-policy#/forwardingRules-get)

Use the **zia_forwarding_control_rule** data source to get information about a forwarding control rule which is used to forward selective Zscaler traffic to specific destinations based on your needs.For example, if you want to forward specific web traffic to a third-party proxy service or if you want to forward source IP anchored application traffic to a specific Zscaler Private Access (ZPA) App Connector or internal application traffic through ZIA threat and data protection engines, use forwarding control by configuring appropriate rules.

## Example Usage

```hcl
# ZIA Forwarding Control - ZPA Gateway
data "zia_forwarding_control_rule" "this" {
  name = "FWD_RULE01"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the forwarding rule.
* `id` - (Optional) A unique identifier assigned to the forwarding rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) - Additional information about the forwarding rule
* `type` - (string) -  The rule type selected from the available options
* `forward_method` - (string) - The type of traffic forwarding method selected from the available options.
* `state` - (string) - Indicates whether the forwarding rule is enabled or disabled.
* `order` - (string) - The order of execution for the forwarding rule order.

`Who, Where and When` supports the following attributes:

* `locations` - (Optional) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `location_groups` - (Optional) You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `ec_groups` - (list) - Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `departments` - (list) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `groups` - (list) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `users` - (list) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

`network services` supports the following attributes:

* `nw_service_groups` - (list) Any number of predefined or custom network service groups to which the rule applies.
* `nw_services`- (list) When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.

`network applications` supports the following attributes:

* `nw_application_groups` - (list) Any number of application groups that you want to control with this rule. The service provides predefined applications that you can group, but not modify
* `nw_applications` - (Optional) When not used it applies the rule to all applications. The service provides predefined applications, which you can group, but not modify.

`source ip addresses` supports the following attributes:

* `src_ip_groups` - (list) Any number of source IP address groups that you want to control with this rule.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `src_ips` - (Optional) You can enter individual IP addresses, subnets, or address ranges.

`destinations` supports the following attributes:

* `dest_addresses`** - (list) -  IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges. If adding multiple items, hit Enter after each entry.
* `dest_countries`** - (list) estination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries. Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes).
* `res_categories`** - (list) List of destination domain categories to which the rule applies.
* `dest_ip_categories`** - (list) identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `dest_ip_groups`** - (list) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `app_service_groups` (list) - Application service groups on which this rule is applied
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `app_services` (list) - Application services on which this rule is applied
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `labels` (list) Labels that are applicable to the rule.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity.
* `devices` (list) Name-ID pairs of devices for which the rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (int) Identifier that uniquely identifies an entity

* `device_groups` (list) Name-ID pairs of device groups for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (int) Identifier that uniquely identifies an entity

* `zpa_gateway` (set) The ZPA Gateway for which this rule is applicable. This field is applicable only for the `ZPA` forwarding method.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (string) The configured name of the entity

* `zpa_app_segments` (set) The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the `ZPA` Gateway forwarding method.
      - `name` - (string) The configured name of the entity
      - `external_id` - (int) Identifier that uniquely identifies an entity

* `proxy_gateway` (set) The proxy gateway for which the rule is applicable. This field is applicable only for the `PROXYCHAIN` forwarding method.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (string) The configured name of the entity.

* `dedicated_ip_gateway` (set) The dedicated IP gateway for which the rule is applicable. This field is applicable only for the `ENATDEDIP` forwarding method.
      - `id` - (int) Identifier that uniquely identifies an entity
      - `name` - (string) The configured name of the entity.

* `zpa_application_segments` (set) List of ZPA Application Segments for which this rule is applicable. This field is applicable only for the `ECZPA` forwarding method (used for Zscaler Cloud Connector).
      - `name` - (string) The configured name of the entity
      - `external_id` - (int) Identifier that uniquely identifies an entity

* `zpa_application_segment_groups` (set) List of ZPA Application Segment Groups for which this rule is applicable. This field is applicable only for the `ECZPA` forwarding method (used for Zscaler Cloud Connector).
      - `name` - (string) The configured name of the entity
      - `external_id` - (int) Identifier that uniquely identifies an entity
