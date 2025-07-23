---
subcategory: "Bandwidth Control"
layout: "zscaler"
page_title: "ZIA): bandwidth_control_rule"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-rules-bandwidth-control-policy
  API documentation https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthControlRules-get
  Retrieves all the rules in the Bandwidth Control policy.

---
# zia_bandwidth_control_rule (Data Source)

* [Official documentation](https://help.zscaler.com/zia/adding-rules-bandwidth-control-policy)
* [API documentation](https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthControlRules-get)

Use the **zia_bandwidth_control_rule** Retrieves all the rules in the Bandwidth Control policy.

**NOTE**: Bandwidth control rule resource is only supported via Zscaler OneAPI.

## Example Usage - By Name

```hcl

data "zia_bandwidth_control_rule" "this" {
    name = "Streaming Media Bandwidth"
}
```

## Example Usage - By ID

```hcl

data "zia_bandwidth_control_rule" "this" {
  id = 154658
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) System-generated identifier for bandwidth control rule
* `name` - (Optional) Rule name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) Additional information about the rule
* `rank` - (int) Admin rank of the Bandwidth Control policy rule
* `state` - (string) Administrative state of the rule.
* `protocols` - (List of string) Protocol criteria. Supported values: `WEBSOCKETSSL_RULE`, `WEBSOCKET_RULE`, `DOHTTPS_RULE`, `TUNNELSSL_RULE`, `HTTP_PROXY`, `FOHTTP_RULE`, `FTP_RULE`, `HTTPS_RULE`, `HTTP_RULE`, `SSL_RULE`, `SSL_RULE`, `TUNNEL_RULE`

* `min_bandwidth` - (int) The minimum percentage of a location's bandwidth you want to be guaranteed for each selected bandwidth class. This percentage includes bandwidth for uploads and downloads.
* `max_bandwidth` - (int) The maximum percentage of a location's bandwidth to be guaranteed for each selected bandwidth class. This percentage includes bandwidth for uploads and downloads.

### Block Attributes

Each of the following blocks supports nested attributes:

#### `labels` - Labels associated with the rule

* `id` - (int) A unique identifier for the label.

#### `bandwidth_classes` - The bandwidth classes to which you want to apply this rule

* `id` - (int) A unique identifier for an entity

#### `time_windows` - The time interval in which the Bandwidth Control policy rule applies

* `id` - (int) A unique identifier for an entity

#### `locations` - The Name-ID pairs of locations to which the DLP policy rule must be applied. Maximum of up to `8` locations. When not used it implies `Any` to apply the rule to all locations

* `id` - (int) Identifier that uniquely identifies an entity

#### `location_groups` - The Name-ID pairs of locations groups to which the DLP policy rule must be applied. Maximum of up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups

* `id` - (Number) Identifier that uniquely identifies an entity
