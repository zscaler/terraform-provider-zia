---
subcategory: "Firewall Filtering Rule"
layout: "zia"
page_title: "ZIA: firewall_filtering_rule"
description: |-
  Retrieve ZIA firewall filtering rule.
  
---
# zia_firewall_filtering_rule (Data Source)

The **zia_firewall_filtering_rule** data source provides details about a specific cloud firewall rule available in the Zscaler Internet Access cloud firewall.

## Example Usage

```hcl
# ZIA Firewall Filtering Rule
data "zia_firewall_filtering_rule" "example" {
    name = "Office 365 One Click Rule"
}

output "zia_firewall_filtering_rule" {
  value = data.zia_firewall_filtering_rule.example
}
```

## Argument Reference

The following arguments are supported:

`Who, Where and When` supports the following attributes:

* `name` - (Required) The firewall automatically creates a Rule Name, which you can change. The maximum length is 31 characters.
* `description` - (Optional) Enter additional notes or information. The description cannot exceed 10,240 characters.
* `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.
* `state` - (Optional) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule.
* `action` - (Optional) Choose the action of the service when packets match the rule. The following actions are accepted: `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BLOCK_ICMP`, `EVAL_NWAPP`
* `rank` - (Optional) By default, the admin ranking is disabled. To use this feature, you must enable admin rank. The default value is `7`.

`Who, Where and When` supports the following attributes:

* `users` - (Optional) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
* `groups` - (Optional) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
* `departments` - (Optional) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
* `locations` - (Optional) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
* `location_groups` - (Optional) You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
* `time_windows` - (Optional) You can manually select up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.

`network services` supports the following attributes:

* `nw_service_groups` - (Optional) Any number of predefined or custom network service groups to which the rule applies.
* `nw_services`- (Optional) When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.

`network applications` supports the following attributes:

* `nw_application_groups` - (Optional) Any number of application groups that you want to control with this rule. The service provides predefined applications that you can group, but not modify
* `nw_applications` - (Optional) When not used it applies the rule to all applications. The service provides predefined applications, which you can group, but not modify.

`source ip addresses` supports the following attributes:

* `src_ip_groups` - (Optional) Any number of source IP address groups that you want to control with this rule.
* `src_ips` - (Optional) You can enter individual IP addresses, subnets, or address ranges.

`destinations` supports the following attributes:

* `dest_addresses`** - (Optional) -  IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges. If adding multiple items, hit Enter after each entry.
* `dest_countries`** - (Optional) Identify destinations based on the location of a server, select Any to apply the rule to all countries or select the countries to which you want to control traffic.
* `dest_ip_categories`** - (Optional) identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
* `dest_ip_groups`** - (Optional) Any number of destination IP address groups that you want to control with this rule.

`Other Supported Arguments` supports the following attributes:

* `id` - (Number) The ID of this resource.
* `last_modified_time` - (Number)
* `access_control` - (String)
* `enable_full_logging` (Boolean)
* `default_rule` - (Boolean)
* `predefined` - (Boolean)

`app_service_groups` - Application service groups on which this rule is applied

* `name` (String)
* `id` (Number)
* `extensions`(Map of String)

* `app_services` - Application services on which this rule is applied

* `name` (String)
* `id` (Number)
* `extensions`(Map of String)

`labels` Labels that are applicable to the rule.

* `name` (String)
* `id` (Number)
* `extensions`(Map of String)
