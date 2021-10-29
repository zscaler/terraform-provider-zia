---
subcategory: "Firewall Filtering Policy"
layout: "zia"
page_title: "ZIA: url_filtering_rules"
description: |-
        Gets all rules in the Firewall Filtering policy.
---
# zia_url_filtering_rules (Data Source)

The **zia_url_filtering_rules** - data source retrieves information about a Firewall Filtering policy rule information for the specified `ID` or `Name`.

```hcl
data "zia_firewall_filtering_rule" "example" {
    name = "Office 365 One Click Rule"
}

output "zia_firewall_filtering_rule" {
  value = data.zia_firewall_filtering_rule.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (String) Name of the Firewall Filtering policy rule
* `id` - (String) Unique identifier for the Firewall Filtering policy rule

### Read-Only

* `order` - (Number) Rule order number of the Firewall Filtering policy rule
* `rank` - (Number) Admin rank of the Firewall Filtering policy rule
* `action` - (String) The action the Firewall Filtering policy rule takes when packets match the rule
* `state` - (String) Determines whether the Firewall Filtering policy rule is enabled or disabled
* `description` - (String) Additional information about the rule
* `last_modified_time` - (Number) Timestamp when the rule was last modified. Ignored if the request is POST or PUT. For GET, ignored if or the rule is current version.
* `src_ips` - (String) User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.
* `dest_addresses` - (String) List of destination IP addresses to which this rule will be applied. CIDR notation can be used for destination IP addresses. If not set, the rule is not restricted to a specific destination addresses unless specified by destCountries, destIpGroups or destIpCategories.
* `dest_ip_categories` - (String) IP address categories of destination for which the DNAT rule is applicable. If not set, the rule is not restricted to specific destination IP categories.
* `dest_countries` - (String) Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
* `nw_applications` - (String) User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.
* `nw_applications` - (String) User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.
* `default_rule` - (Boolean) If set to true, the default rule is applied
* `predefined` - (Boolean) If set to true, a predefined rule is applied

`locations` - (List of Object) The locations to which the Firewall Filtering policy rule applies

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`location_groups` - (List of Object) The location groups to which the Firewall Filtering policy rule applies

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`departments` - (List of Object) The departments to which the Firewall Filtering policy rule applies

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`groups` - (List of Object) The groups to which the Firewall Filtering policy rule applies

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`users` - (List of Object) The users to which the Firewall Filtering policy rule applies

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`time_windows` - (List of Object) The time interval in which the Firewall Filtering policy rule applies

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`labels` -(List of Object) Labels that are applicable to the rule.

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`last_modified_by`

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`src_ip_groups`

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`dest_ip_groups`

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`nw_services`

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`nw_service_groups`

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`nw_application_groups`

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`app_services`

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)

`app_service_groups`

* `id` -(Number)
* `name` -(String)
* `extensions` -(Map of String)

`labels`

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` -(String) The configured name of the entity
* `extensions` - (Map of String)
