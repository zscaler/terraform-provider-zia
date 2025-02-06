---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_rule"
description: |-
  Creates and manages ZIA Cloud firewall filtering rule.
---

# Resource: zia_firewall_filtering_rule

The **zia_firewall_filtering_rule** resource allows the creation and management of ZIA Cloud Firewall filtering rules in the Zscaler Internet Access.

**NOTE 1** Zscaler Cloud Firewall contain default and predefined rules which are placed in their respective orders. These rules `CANNOT` be deleted. When configuring your rules make sure that the `order` attributue value consider these pre-existing rules so that Terraform can place the new rules in the correct position, and drifts can be avoided. i.e If there are 2 pre-existing rules, you should start your rule order at `3` and manage your rule sets from that number onwards. The provider will reorder the rules automatically while ignoring the order of pre-existing rules, as the API will be responsible for moving these rules to their respective positions as API calls are made.

The most common default rules are:

The most common default and predefined rules:

|              Rule Names                 |  Default or Predefined   |   Rule Number Associated |
|:---------------------------------------:|:------------------------:|:------------------------:|
|------------------------------|--------------------------|-------------------|
|  `Office 365 One Click Rule`            |      `Predefined`,       |           `Yes`          |
|  `UCaaS One Click Rule`                 |      `Predefined`,       |           `Yes`          |
|  `Block All IPv6`                       |      `Predefined`,       |           `Yes`          |
|  `Block malicious IPs and domains`      |      `Predefined`,       |           `Yes`          |
|  `Default Firewall Filtering Rule`      |      `Default`,          |           `Yes`          |
|-----------------------|-----------------------------|

**NOTE 2** Certain attributes on `predefined` rules can still be managed or updated via Terraform such as:

- `description` - (Optional) Enter additional notes or information. The description cannot exceed 10,240 characters.
- `state` - (Optional) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule.
- `labels` (list) - Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

**NOTE 3** The import of `predefined` rules is still possible in case you want o have them under the Terraform management; however, remember that these rules cannot be deleted. That means, the provider will fail when executing `terraform destroy`; hence, you must remove the rules you want to delete, and re-run `terraform apply` instead.

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

- `name` - (Required) Name of the Firewall Filtering policy rule
- `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.

**NOTE 1** Zscaler Cloud Firewall contain `default` and `predefined` rules which are placed in their respective orders. These rules `CANNOT` be deleted. When configuring your rules make sure that the `order` attributue value consider these pre-existing rules so that Terraform can place the new rules in the correct position, and drifts can be avoided. i.e If there are 2 pre-existing rules, you should start your rule order at `3` and manage your rule sets from that number onwards. The provider will reorder the rules automatically while ignoring the order of pre-existing rules, as the API will be responsible for moving these rules to their respective positions as API calls are made.

The most common default rules are:

- `Office 365 One Click Rule`
- `UCaaS One Click Rule`
- `Block All IPv6`
- `Block malicious IPs and domains`
- `Default Firewall Filtering Rule`

**NOTE 2** Certain attributes on `predefined` rules can still be managed or updated via Terraform such as:

- `description` - (Optional) Enter additional notes or information. The description cannot exceed 10,240 characters.
- `state` - (Optional) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule.
- `labels` (list) - Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

### Optional

- `description` - (String) Enter additional notes or information. The description cannot exceed 10,240 characters.
- `state` - (String) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule.
- `action` - (String) Choose the action of the service when packets match the rule. The following actions are accepted: `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BLOCK_ICMP`, `EVAL_NWAPP`
- `rank` - (Integer) By default, the admin ranking is disabled. To use this feature, you must enable admin rank. The default value is `7`.

`Who, Where and When` supports the following attributes:

- `locations` (list) - You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
- `location_groups` (list) - You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `users` (list) - You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `groups` (list) - You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
- `departments` (list) - Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (Integer) Identifier that uniquely identifies an entity
- `time_windows` (list) - You can manually select up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (Integer) Identifier that uniquely identifies an entity

`network services` supports the following attributes:

- `nw_service_groups` (list) - Any number of predefined or custom network service groups to which the rule applies.
- `nw_services` (list) - When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.

`network applications` -  supports the following attributes:

- `nw_application_groups` (list) - Any number of application groups that you want to control with this rule. The service provides predefined applications that you can group, but not modify
- `nw_applications` (list) - When not used it applies the rule to all applications. The service provides predefined applications, which you can group, but not modify.

- `src_ip_groups` (list) - Any number of source IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity
- `src_ips` (list) - You can enter individual IP addresses, subnets, or address ranges.

- `dest_addresses` (list) - IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges.
      **NOTE**: PLEASE BE AWARE. The API supports ONLY `IPv4` addresses. `IPV6` addresses are not supported.

- `dest_countries` (list) - Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
      **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

- `source_countries` (list) - The list of source countries that must be included or excluded from the rule based on the excludeSrcCountries field value. If no value is set, this field is ignored during policy evaluation and the rule is applied to all source countries.
      **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

- `dest_ip_categories` (list) - identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
      - `id` - (Integer) Identifier that uniquely identifies an entity
- `dest_ip_groups` (list) - Any number of destination IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `app_service_groups` (list) - Application service groups on which this rule is applied
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `app_services`  (list) - Application services on which this rule is applied
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `labels` (list) - Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `workload_groups` (list) - The list of preconfigured workload groups to which the policy must be applied
  - `id` - (Optional) A unique identifier assigned to the workload group
  - `name` - (Optional) The name of the workload group

- `Other Exported Arguments`
  - `enable_full_logging` (Boolean)
  `Aggregate`: The service groups together individual sessions based on { user, rule, network service, network application } and records them periodically.`Full`: The service logs all sessions of the rule individually, except HTTP(S). Only Block rules support full logging. Full logging on all other rules requires the Full Logging license.
  - `predefined` - (Boolean) If set to true, a predefined rule is applied

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
