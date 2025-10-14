---
subcategory: "Forwarding Control Policy"
layout: "zscaler"
page_title: "ZIA): forwarding_control_rule"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-forwarding-policy
  API documentation https://help.zscaler.com/zia/forwarding-control-policy#/forwardingRules-get
  Creates and manages forwarding control rule.
---

# zia_forwarding_control_rule (Resource)

* [Official documentation](https://help.zscaler.com/zia/configuring-forwarding-policy)
* [API documentation](https://help.zscaler.com/zia/forwarding-control-policy#/forwardingRules-get)

The **zia_forwarding_control_rule** resource allows the creation and management of ZIA Forwarding Control rules in the Zscaler Internet Access.

⚠️ **WARNING:**  - [PR #373](https://github.com/zscaler/terraform-provider-zia/pull/373) - The resource `zia_forwarding_control_rule` now pauses for 60 seconds before proceeding with the create or update process whenever the `forward_method` attribute is set to `ZPA`. In case of a failure related to resource synchronization, the provider will retry the resource creation or update up to 3 times, waiting 30 seconds between each retry. This behavior ensures that ZIA and ZPA have sufficient time to synchronize and replicate the necessary resource IDs, reducing the risk of transient errors during provisioning.

  **NOTE**: This retry mechanism helps to automatically overcome temporary latency without manual intervention. This behavior does not affect forwarding rules configured with other forward_methods such as `DIRECT`.

## Example Usage - DIRECT Forwarding Method

```hcl
resource "zia_forwarding_control_rule" "this" {
  name               = "FC_DIRECT_RULE"
  description        = "FC_DIRECT_RULE"
  order              = 1
  rank               = 7
  state              = "ENABLED"
  type               = "FORWARDING"
  forward_method     = "DIRECT"
  src_ips            = ["192.168.200.200"]
  dest_addresses     = ["192.168.255.1"]
  dest_ip_categories = ["ZSPROXY_IPS", "CUSTOM_01"]
  dest_countries     = ["CA", "US"]
}
```

## Example Usage - ZPA Forwarding Method

  ⚠️ **WARNING:**: You must use the [ZPA provider](https://registry.terraform.io/providers/zscaler/zpa/latest/docs) in combination with the ZIA Terraform Provider to successfully configure a Forwarding control rule where the `forward_method` is `ZPA`

```hcl
# ZPA Server Group
data "zpa_server_group" "this" {
  name = "Server_Group_IP_Source_Anchoring"
}

# ZPA Application Segment
data "zpa_application_segment" "this" {
  name = "App_Segment_IP_Source_Anchoring"
}

resource "zia_forwarding_control_zpa_gateway" "this" {
    name = "ZPA_GW01"
    description = "ZPA_GW01"
    type = "ZPA"
    zpa_server_group {
      external_id = data.zpa_server_group.this.id
      name = data.zpa_server_group.this.id
    }
    zpa_app_segments {
        external_id = data.zpa_application_segment.this.id
        name = data.zpa_application_segment.this.name
    }
}

resource "zia_forwarding_control_rule" "this" {
  name           = "ZPA_FORWARDING_RULE"
  description    = "ZPA_FORWARDING_RULE"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  type           = "FORWARDING"
  forward_method = "ZPA"
  zpa_gateway {
    id   = zia_forwarding_control_zpa_gateway.this.id
    name = zia_forwarding_control_zpa_gateway.this.name
  }
  zpa_app_segments {
    name        = data.zpa_application_segment.this.name
    external_id = data.zpa_application_segment.this.id
  }
}
```

## Example Usage - PROXYCHAIN Forwarding Method

  ⚠️ **WARNING:**: Creating or retrieving a Proxy Gateway via API is not currently supported; hence, the `id` and `name` for the `proxy_gateway` must be passed manually to the `proxy_gateway` block in the below configuration.

```hcl
resource "zia_forwarding_control_rule" "this" {
  name               = "PROXYCHAIN_FORWARDING_RULE"
  description        = "PROXYCHAIN_FORWARDING_RULE"
  order              = 1
  rank               = 7
  state              = "ENABLED"
  type               = "FORWARDING"
  forward_method     = "PROXYCHAIN"
  src_ips            = ["192.168.200.200"]
  dest_addresses     = ["192.168.255.1"]
  dest_ip_categories = ["ZSPROXY_IPS", "CUSTOM_01"]
  dest_countries     = ["CA", "US"]
  proxy_gateway {
    id   = 2589270
    name = "ProxyGW01"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the Firewall Filtering policy rule
* `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.
* `type` - (string) -  The rule type selected from the available options. Supported Values: ``FORWARDING``
* `forward_method` - (string) - The type of traffic forwarding method selected from the available options.
      - `DIRECT` - If forward_method is `DIRECT` no other attribute is required.
      - `ZPA` - If forward_method is `ZPA` the attributes `zpa_gateway` and `zpa_app_segments` are required.
      - `PROXYCHAIN` - If forward_method is `PROXYCHAIN` the attributes `proxy_gateway` is required.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) - Additional information about the forwarding rule
* `state` - (string) - Indicates whether the forwarding rule is enabled or disabled. Supported values are: `ENABLED` and `DISABLED`.
* `order` - (int) - The order of execution for the forwarding rule order.

`Who, Where and When` supports the following attributes:

* `locations` - (Optional) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (int) Identifier that uniquely identifies an entity
* `location_groups` - (Optional) You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (int) Identifier that uniquely identifies an entity
* `ec_groups` - (list) - Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies
      - `id` - (int) Identifier that uniquely identifies an entity
* `departments` - (list) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (int) Identifier that uniquely identifies an entity
* `groups` - (list) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (int) Identifier that uniquely identifies an entity
* `users` - (list) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (int) Identifier that uniquely identifies an entity

`network services` supports the following attributes:

* `nw_service_groups` - (list) Any number of predefined or custom network service groups to which the rule applies.
* `nw_services`- (list) When not used it applies the rule to all network services or you can select specific network services. The Zscaler firewall has predefined services and you can configure up to `1,024` additional custom services.

`network applications` supports the following attributes:

* `nw_application_groups` - (list) Any number of application groups that you want to control with this rule. The service provides predefined applications that you can group, but not modify
* `nw_applications` - (Optional) When not used it applies the rule to all applications. The service provides predefined applications, which you can group, but not modify.

`source ip addresses` supports the following attributes:

* `src_ip_groups` - (list) Any number of source IP address groups that you want to control with this rule.
      - `id` - (int) Identifier that uniquely identifies an entity
* `src_ipv6_groups` - (list) Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.
      - `id` - (int) Identifier that uniquely identifies an entity
* `src_ips` - (Optional) You can enter individual IP addresses, subnets, or address ranges.

`destinations` supports the following attributes:

* `dest_addresses`** - (list) -  IP addresses and fully qualified domain names (FQDNs), if the domain has multiple destination IP addresses or if its IP addresses may change. For IP addresses, you can enter individual IP addresses, subnets, or address ranges. If adding multiple items, hit Enter after each entry.
* `dest_countries`** - (list) destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries. Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes).
* `res_categories`** - (list) List of destination domain categories to which the rule applies.
* `dest_ip_categories`** - (list) identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
      - `id` - (int) Identifier that uniquely identifies an entity
* `dest_ip_groups`** - (list) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (int) Identifier that uniquely identifies an entity
* `dest_ipv6_groups`** - (list) Destination IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.
      - `id` - (int) Identifier that uniquely identifies an entity

* `app_service_groups` (list) - Application service groups on which this rule is applied
      - `id` - (int) Identifier that uniquely identifies an entity

* `labels` (list) Labels that are applicable to the rule.
      - `id` - (int) Identifier that uniquely identifies an entity

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

* `zpa_application_segments` (set) List of ZPA Application Segments for which this rule is applicable. This field is applicable only for the `ECZPA` forwarding method (used for Zscaler Cloud Connector).
      - `name` - (string) The configured name of the entity
      - `external_id` - (int) Identifier that uniquely identifies an entity

* `zpa_application_segment_groups` (set) List of ZPA Application Segment Groups for which this rule is applicable. This field is applicable only for the `ECZPA` forwarding method (used for Zscaler Cloud Connector).
      - `name` - (string) The configured name of the entity
      - `external_id` - (int) Identifier that uniquely identifies an entity

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_forwarding_control_rule** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import zia_forwarding_control_rule.example <rule_id>
```

or

```shell
terraform import zia_forwarding_control_rule.example <rule_name>
```
