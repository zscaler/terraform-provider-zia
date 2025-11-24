---
subcategory: "Traffic Capture Policy"
layout: "zscaler"
page_title: "ZIA: traffic_capture_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/about-traffic-capture-policy
  API documentation https://help.zscaler.com/zia/traffic-capture-policy#/trafficCaptureRules-get
  Get information about traffic capture rules
---

# zia_traffic_capture_rules (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-traffic-capture-policy)
* [API documentation](https://help.zscaler.com/zia/traffic-capture-policy#/trafficCaptureRules-get)

The **zia_traffic_capture_rules** resource allows the creation and management of ZIA traffic capture rules in the Zscaler Internet Access.

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

resource "zia_traffic_capture_rules" "example" {
    name                = "Example Traffic Capture Rule"
    description         = "Example traffic capture rule for engineering department"
    action              = "ALLOW"
    state               = "ENABLED"
    order               = 1
    enable_full_logging = true
    txn_size_limit      = "UNLIMITED"
    txn_sampling        = "HUNDRED_PERCENT"
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

* `name` - (Required) Name of the Traffic Capture policy rule
* `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.

**NOTE 1** Zscaler Traffic Capture contains `default` and `predefined` rules which are placed in their respective orders. These rules `CANNOT` be deleted. When configuring your rules make sure that the `order` attribute value considers these pre-existing rules so that Terraform can place the new rules in the correct position, and drifts can be avoided. i.e If there are 2 pre-existing rules, you should start your rule order at `3` and manage your rule sets from that number onwards. The provider will reorder the rules automatically while ignoring the order of pre-existing rules, as the API will be responsible for moving these rules to their respective positions as API calls are made.

**NOTE 2** Certain attributes on `predefined` rules can still be managed or updated via Terraform such as:

* `description` - (Optional) Enter additional notes or information. The description cannot exceed 10,240 characters.
* `state` - (Optional) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule. Supported Values: `ENABLED`, `DISABLED`

* `order` - (Integer) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.

* `labels` (list) - Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

**NOTE 3** The following attributes on `predefined` rules **cannot** be updated:

* `name` - Predefined rule names are fixed and cannot be changed

* `action` - (String) The action configured for the rule that must take place if the traffic matches the rule criteria. The following actions are accepted: `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BLOCK_ICMP`, `EVAL_NWAPP`

* Most other attributes that define the rule's behavior

### Optional

* `description` - (String) Additional information about the rule. The description cannot exceed 10,240 characters.
* `state` - (String) Determines whether the Traffic Capture policy rule is enabled or disabled. An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule. Supported Values: `ENABLED`, `DISABLED`
* `action` - (String) The action the Traffic Capture policy rule takes when packets match the rule. The following actions are accepted: `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BLOCK_ICMP`, `EVAL_NWAPP`
* `rank` - (Integer) Admin rank of the Traffic Capture policy rule. By default, the admin ranking is disabled. To use this feature, you must enable admin rank in UI first. The default value is `7`. Valid range is 0-7. Visit to learn more [About Admin Rank](https://help.zscaler.com/zia/about-admin-rank)
* `txn_size_limit` - (String) The maximum size of traffic to capture per connection. Supported values: `NONE`, `UNLIMITED`, `THIRTY_TWO_KB`, `TWO_FIFTY_SIX_KB`, `TWO_MB`, `FOUR_MB`, `THIRTY_TWO_MB`, `SIXTY_FOUR_MB`
* `txn_sampling` - (String) The percentage of connections sampled for capturing each time the rule is triggered. Supported values: `NONE`, `ONE_PERCENT`, `TWO_PERCENT`, `FIVE_PERCENT`, `TEN_PERCENT`, `TWENTY_FIVE_PERCENT`, `HUNDRED_PERCENT`

`Who, Where and When` supports the following attributes:

* `locations` (list) - You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `location_groups` (list) - You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `users` (list) - You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `groups` (list) - You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `departments` (list) - Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `devices` (list) - Specifies devices that are managed using Zscaler Client Connector.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `device_groups` (list) - This field is applicable for devices that are managed using Zscaler Client Connector.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `device_trust_levels` - (Optional) List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation. Supported values: `ANY`, `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`

* `time_windows` (list) - The time interval in which the Traffic Capture policy rule applies. You can manually select up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (Integer) Identifier that uniquely identifies an entity

`Network Services` supports the following attributes:

* `nw_service_groups` (list) - Any number of predefined or custom network service groups to which the rule applies.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `nw_services` (list) - User-defined network services on which the rule is applied. When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.
      - `id` - (Integer) Identifier that uniquely identifies an entity

`Network Applications` supports the following attributes:

* `nw_application_groups` (list) - Any number of application groups that you want to control with this rule. The service provides predefined applications that you can group, but not modify.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `nw_applications` (list) - User-defined network service applications on which the rule is applied. When not used it applies the rule to all applications. The service provides predefined applications, which you can group, but not modify. To retrieve the list of cloud applications, use the data source: `zia_cloud_applications`

* `src_ip_groups` (list) - Any number of source IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `src_ips` (list) - You can enter individual IP addresses, subnets, or address ranges.

* `dest_addresses` (list) - IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges.
      **NOTE**: PLEASE BE AWARE. The API supports ONLY `IPv4` addresses. `IPV6` addresses are not supported.

* `dest_countries` (list) - Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
      **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `source_countries` (list) - The list of source countries that must be included or excluded from the rule based on the excludeSrcCountries field value. If no value is set, this field is ignored during policy evaluation and the rule is applied to all source countries.
      **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `dest_ip_categories` (list) - IP address categories of destination for which the rule is applicable. Select Any to apply the rule to all categories or select the specific categories you want to control.
* `dest_ip_groups` (list) - User-defined destination IP address groups on which the rule is applied. Any number of destination IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `app_service_groups` (list) - Application service groups on which this rule is applied.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `labels` (list) - Labels that are applicable to the rule. You can manually select up to `1` label.
      - `id` - (Integer) Identifier that uniquely identifies an entity

### Computed / Exported Attributes

* `id` - (String) The unique identifier for the Traffic Capture rule
* `rule_id` - (Integer) The unique identifier for the Traffic Capture rule (same as `id` but as integer)
* `enable_full_logging` - (Boolean) A Boolean value that indicates whether full logging is enabled. A `true` value indicates that full logging is enabled, whereas a `false` value indicates that aggregate logging is enabled. Aggregate logging groups together individual sessions based on { user, rule, network service, network application } and records them periodically. Full logging logs all sessions of the rule individually, except HTTP(S).
* `default_rule` - (Boolean) If set to true, the default rule is applied
* `predefined` - (Boolean) If set to true, a predefined rule is applied

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_traffic_capture_rules** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import zia_traffic_capture_rules.example <rule_id>
```

or

```shell
terraform import zia_traffic_capture_rules.example <rule_name>
```
