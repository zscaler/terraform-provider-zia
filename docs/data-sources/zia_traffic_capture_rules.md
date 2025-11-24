---
subcategory: "Traffic Capture Policy"
layout: "zscaler"
page_title: "ZIA: traffic_capture_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/about-traffic-capture-policy
  API documentation https://help.zscaler.com/zia/traffic-capture-policy#/trafficCaptureRules-get
  Get information about traffic capture rules
---

# zia_traffic_capture_rules (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-traffic-capture-policy)
* [API documentation](https://help.zscaler.com/zia/traffic-capture-policy#/trafficCaptureRules-get)

Use the **zia_traffic_capture_rules** data source to get information about a traffic capture rules available in the Zscaler Internet Access cloud firewall.

## Example Usage - Retrieve by Name

```hcl
data "zia_traffic_capture_rules" "example" {
    name = "Capture_Rule01"
}
```

## Example Usage - Retrieve by ID

```hcl
data "zia_traffic_capture_rules" "example" {
    id = 1254674585
}
```

## Argument Reference

The following arguments are supported:

### Required

At least one of the following must be provided:

* `id` - (Integer) Unique identifier for the Traffic Capture policy rule
* `name` - (String) Name of the Traffic Capture policy rule

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Additional information about the rule. Cannot exceed 10,240 characters.
* `order` - (Integer) Rule order number. Policy rules are evaluated in ascending numerical order.
* `rank` - (Integer) Admin rank of the Traffic Capture policy rule. Default value is `7`.
* `state` - (String) Rule state. An enabled rule is actively enforced. Values: `ENABLED`, `DISABLED`
* `action` - (String) The action the rule takes when packets match. Values: `CAPTURE`, `SKIP`
* `access_control` - (String) The admin's access privilege to this rule based on the assigned role
* `default_rule` - (Boolean) Indicates if this is a default rule
* `predefined` - (Boolean) Indicates if this is a predefined rule
* `txn_size_limit` - (String) The maximum size of traffic to capture per connection. Supported Values: `NONE`, `UNLIMITED`, `THIRTY_TWO_KB`, `TWO_FIFTY_SIX_KB`, `TWO_MB`, `FOUR_MB`, `THIRTY_TWO_MB`, `SIXTY_FOUR_MB`

* `txn_sampling` - (String) The percentage of connections sampled for capturing each time the rule is triggered. Supported Values: `NONE`, `ONE_PERCENT`, `TWO_PERCENT`, `FIVE_PERCENT`, `TEN_PERCENT`, `TWENTY_FIVE_PERCENT`, `HUNDRED_PERCENT`

* `last_modified_time` - (Integer) Timestamp when the rule was last modified

### Last Modified By

* `last_modified_by` - (List) User who last modified the rule
  * `id` - (Integer) Identifier of the user
  * `name` - (String) Name of the user
  * `extensions` - (Map of String) Additional user extensions

### Who, Where and When

* `locations` - (List) Locations for which the rule applies. You can manually select up to `8` locations.
  * `id` - (Integer) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String) Additional location extensions

* `location_groups` - (List) Location groups for which the rule applies. You can manually select up to `32` location groups.
  * `id` - (Integer) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String) Additional location group extensions

* `users` - (List) Users for which the rule applies. You can manually select up to `4` general and/or special users.
  * `id` - (Integer) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String) Additional user extensions

* `groups` - (List) Groups for which the rule applies. You can manually select up to `8` groups.
  * `id` - (Integer) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String) Additional group extensions

* `departments` - (List) Departments for which the rule applies. You can apply to any number of departments.
  * `id` - (Integer) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String) Additional department extensions

* `time_windows` - (List) Time intervals in which the rule applies. You can manually select up to `2` time intervals.
  * `id` - (Integer) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String) Additional time window extensions

### Network Services

* `nw_services` - (List) Network services for which the rule applies. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.
  * `id` - (Integer) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String) Additional network service extensions

* `nw_service_groups` - (List) Network service groups for which the rule applies.
  * `id` - (Integer) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String) Additional network service group extensions

### Network Applications

* `nw_applications` - (List of String) Network service applications. The service provides predefined applications.

* `nw_application_groups` - (List) Network application groups for which the rule applies.
  * `id` - (Integer) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String) Additional network application group extensions

### Source IP

* `src_ip_groups` - (Optional) Any number of source IP address groups that you want to control with this rule.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `src_ips` - (Optional) You can enter individual IP addresses, subnets, or address ranges.

* `source_countries`** - (List of String) Identify destinations based on the location of a server. Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `exclude_src_countries` - (Boolean) Indicates whether the countries specified in the sourceCountries field are included or excluded from the rule. A true value denotes that the specified source countries are excluded from the rule. A false value denotes that the rule is applied to the source countries if there is a match.

### Destination IP

* `dest_ip_groups`** - (Optional) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `dest_addresses`** - (List of String) -  IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges. If adding multiple items, hit Enter after each entry.
* `dest_countries`** - (List of String) Identify destinations based on the location of a server. Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

### Workload Group

* `workload_groups` (List) The list of preconfigured workload groups to which the policy must be applied
  * `id` - (Number) A unique identifier assigned to the workload group
  * `name` - (String) The name of the workload group
