---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA): firewall_filtering_rule"
description: |-
  Get information about firewall filtering rule.

---
# Data Source: zia_firewall_filtering_rule

Use the **zia_firewall_filtering_rule** data source to get information about a cloud firewall rule available in the Zscaler Internet Access cloud firewall.

## Example Usage

```hcl
# ZIA Firewall Filtering Rule
data "zia_firewall_filtering_rule" "example" {
    name = "Office 365 One Click Rule"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Firewall Filtering policy rule
* `id` - (Optional) Unique identifier for the Firewall Filtering policy rule

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (Optional) Enter additional notes or information. The description cannot exceed 10,240 characters.
* `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.
* `state` - (Optional) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule.
* `action` - (Optional) Choose the action of the service when packets match the rule. The following actions are accepted: `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BLOCK_ICMP`, `EVAL_NWAPP`
* `rank` - (Optional) By default, the admin ranking is disabled. To use this feature, you must enable admin rank. The default value is `7`.

`Who, Where and When` supports the following attributes:

* `locations` - (Optional) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)
* `location_groups` - (Optional) You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)
* `users` - (Optional) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)
* `groups` - (Optional) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)
* `departments` - (Optional) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `time_windows` - (Optional) You can manually select up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

`network services` supports the following attributes:

* `nw_service_groups` - (Optional) Any number of predefined or custom network service groups to which the rule applies.
* `nw_services`- (Optional) When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.

`network applications` supports the following attributes:

* `nw_application_groups` - (Optional) Any number of application groups that you want to control with this rule. The service provides predefined applications that you can group, but not modify
* `nw_applications` - (Optional) When not used it applies the rule to all applications. The service provides predefined applications, which you can group, but not modify.

`source ip addresses` supports the following attributes:

* `src_ip_groups` - (Optional) Any number of source IP address groups that you want to control with this rule.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)
* `src_ips` - (Optional) You can enter individual IP addresses, subnets, or address ranges.

`destinations` supports the following attributes:

* `dest_addresses`** - (Optional) -  IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges. If adding multiple items, hit Enter after each entry.
* `dest_countries`** - (Optional) Identify destinations based on the location of a server, select Any to apply the rule to all countries or select the countries to which you want to control traffic.
* `dest_ip_categories`** - (Optional) identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)
* `dest_ip_groups`** - (Optional) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `app_service_groups` - Application service groups on which this rule is applied
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `app_services` - Application services on which this rule is applied
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `labels` Labels that are applicable to the rule.
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `workload_groups` (List) The list of preconfigured workload groups to which the policy must be applied
  * `id` - (Number) A unique identifier assigned to the workload group
  * `name` - (String) The name of the workload group

* `Other Exported Arguments`
  * `id` - (Number) The ID of this resource.
  * `last_modified_time` - (Number)
  * `access_control` - (String)
  * `enable_full_logging` (Boolean)
  * `default_rule` - (Boolean)
  * `predefined` - (Boolean)
