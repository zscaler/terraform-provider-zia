---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_rule"
description: |-
  Creates and manages ZIA Cloud firewall filtering rule.
---

# Resource: zia_firewall_filtering_rule

The **zia_firewall_filtering_rule** resource allows the creation and management of ZIA Cloud Firewall filtering rules in the Zscaler Internet Access.

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

resource "zia_firewall_filtering_rule" "example" {
    name                = "Example"
    description         = "Example"
    action              = "ALLOW"
    state               = "ENABLED"
    order               = 1
    enable_full_logging = true
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

* `name` - (Required) Name of the Firewall Filtering policy rule
* `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.

### Optional

* `description` - (Optional) Enter additional notes or information. The description cannot exceed 10,240 characters.
* `state` - (Optional) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule.
* `action` - (Optional) Choose the action of the service when packets match the rule. The following actions are accepted: `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BLOCK_ICMP`, `EVAL_NWAPP`
* `rank` - (Optional) By default, the admin ranking is disabled. To use this feature, you must enable admin rank. The default value is `7`.

`Who, Where and When` supports the following attributes:

* `locations` - (Optional) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `location_groups` - (Optional) You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `users` - (Optional) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `groups` - (Optional) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `departments` - (Optional) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `time_windows` - (Optional) You can manually select up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (Integer) Identifier that uniquely identifies an entity

`network services` supports the following attributes:

* `nw_service_groups` - (Optional) Any number of predefined or custom network service groups to which the rule applies.
* `nw_services`- (Optional) When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.

`network applications` supports the following attributes:

* `nw_application_groups` - (Optional) Any number of application groups that you want to control with this rule. The service provides predefined applications that you can group, but not modify
* `nw_applications` - (Optional) When not used it applies the rule to all applications. The service provides predefined applications, which you can group, but not modify.

* `src_ip_groups` - (Optional) Any number of source IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `src_ips` - (Optional) You can enter individual IP addresses, subnets, or address ranges.

* `dest_addresses`** - (Optional) -  IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges.
      **NOTE**: PLEASE BE AWARE. The API supports ONLY `IPv4` addresses. `IPV6` addresses are not supported.

* `dest_countries`** - (Optional) Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
      **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `source_countries`** - (Optional) The list of source countries that must be included or excluded from the rule based on the excludeSrcCountries field value. If no value is set, this field is ignored during policy evaluation and the rule is applied to all source countries.
      **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `dest_ip_categories`** - (Optional) identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `dest_ip_groups`** - (Optional) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `app_service_groups` - Application service groups on which this rule is applied
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `app_services` - Application services on which this rule is applied
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `labels` Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `workload_groups` (Optional) The list of preconfigured workload groups to which the policy must be applied
  * `id` - (Optional) A unique identifier assigned to the workload group
  * `name` - (Optional) The name of the workload group

* `Other Exported Arguments`
  * `enable_full_logging` (Boolean)
  * `predefined` - (Boolean) If set to true, a predefined rule is applied

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_firewall_filtering_rule** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import zia_firewall_filtering_rule.example <rule_id>
```

or

```shell
terraform import zia_firewall_filtering_rule.example <rule_name>
```
