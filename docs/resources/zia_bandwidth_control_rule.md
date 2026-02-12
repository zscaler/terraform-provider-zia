---
subcategory: "Bandwidth Control"
layout: "zscaler"
page_title: "ZIA): bandwidth_control_rule"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-rules-bandwidth-control-policy
  API documentation https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthControlRules-post
  Adds a new Bandwidth Control policy rule

---
# zia_bandwidth_control_rule (Resource)

* [Official documentation](https://help.zscaler.com/zia/adding-rules-bandwidth-control-policy)
* [API documentation](https://help.zscaler.com/zia/bandwidth-control-classes#/)

Use the **zia_bandwidth_control_rule** resource allows the creation and management of ZIA Bandwidth Control Rules in the Zscaler Internet Access cloud or via the API.

**NOTE**: Bandwidth control rule resource is only supported via Zscaler OneAPI.

## Example Usage - By Name

```hcl
data "zia_bandwidth_classes_web_conferencing" "this" {
    name = "BANDWIDTH_CAT_WEBCONF"
}

resource "zia_bandwidth_control_rule" "this" {
    name = "Streaming Media Bandwidth"
    description = "Streaming Media Bandwidth"
    state = "ENABLED"
    order = 1
    rank = 7
    min_bandwidth = 5
    max_bandwidth = 100
    protocols = ["ANY_RULE"]
    bandwidth_classes  {
        id = [data.zia_bandwidth_classes_web_conferencing.this.id]
    }
    labels  {
        id = [1503197]
    }
    location_groups {
        id = [8061255]
    }
    time_windows {
        id = [483]
    }
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Optional) Rule name.
* `order` - (Int) Rule order. Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the rule order reflects this rule's place in the order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (Optional) Additional information about the rule
* `rank` - (Optional) Admin rank of the Bandwidth Control policy rule
* `state` - (Optional) Administrative state of the rule.
* `protocols` - (List of Object) Protocol criteria. Supported values: `WEBSOCKETSSL_RULE`, `WEBSOCKET_RULE`, `DOHTTPS_RULE`, `TUNNELSSL_RULE`, `HTTP_PROXY`, `FOHTTP_RULE`, `FTP_RULE`, `HTTPS_RULE`, `HTTP_RULE`, `SSL_RULE`, `SSL_RULE`, `TUNNEL_RULE`

* `min_bandwidth` - (Optional) The minimum percentage of a location's bandwidth you want to be guaranteed for each selected bandwidth class. This percentage includes bandwidth for uploads and downloads.
* `max_bandwidth` - (Optional) The maximum percentage of a location's bandwidth to be guaranteed for each selected bandwidth class. This percentage includes bandwidth for uploads and downloads.

### Block Attributes

Each of the following blocks supports nested attributes:

#### `labels` - Labels associated with the rule

* `id` - (int) A unique identifier for the label.

#### `bandwidth_classes` - The bandwidth classes to which you want to apply this rule

* `id` - (int) A unique identifier for an entity

#### `time_windows` - The time interval in which the Bandwidth Control policy rule applies

* `id` - (int) A unique identifier for an entity

#### `locations` - The Name-ID pairs of locations to which the DLP policy rule must be applied. Maximum of up to `32` locations. When not used it implies `Any` to apply the rule to all locations

* `id` - (int) Identifier that uniquely identifies an entity

#### `location_groups` - The Name-ID pairs of locations groups to which the DLP policy rule must be applied. Maximum of up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups

* `id` - (Number) Identifier that uniquely identifies an entity

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_bandwidth_control_rule** can be imported by using `<RULE_ID>` or `<RULE_NAME>` as the import ID.

For example:

```shell
terraform import zia_bandwidth_control_rule.this <rule_id>
```

```shell
terraform import zia_bandwidth_control_rule.this <"rule_name">
```
