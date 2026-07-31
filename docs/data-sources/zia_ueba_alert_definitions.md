---
subcategory: "Security & UEBA Alerts"
layout: "zscaler"
page_title: "ZIA: ueba_alert_definitions"
description: |-
  Official documentation https://help.zscaler.com/zia/about-alerts
  API documentation https://help.zscaler.com/legacy-apis/security-ueba-alerts#/alertDefinitions-get
  Get information about a ZIA Security & UEBA alert definition.
---

# zia_ueba_alert_definitions (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-alerts)
* [API documentation](https://help.zscaler.com/legacy-apis/security-ueba-alerts#/alertDefinitions-get)

Use the **zia_ueba_alert_definitions** data source to get information about a Security & UEBA alert definition. The alert definition can be looked up by its system-generated `id` or by its `alert_name`.

## Example Usage - By Alert Name

```hcl
data "zia_ueba_alert_definitions" "this" {
  alert_name = "OUTGOING_VIRUSES"
}
```

## Example Usage - By ID

```hcl
data "zia_ueba_alert_definitions" "this" {
  id = 123456
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, Integer) The system-generated identifier of the alert definition.
* `alert_name` - (Optional, String) The alert name that identifies the threat or event type the alert is triggered for.

~> **NOTE** One of `id` or `alert_name` must be provided.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - (String) The status of the alert rule.
* `occurrence` - (String) Specifies the occurrence of an ongoing alert for a specific threat or event type.
* `traffic_change_percent` - (Integer) Specifies the percentage change in traffic.
* `interval` - (String) The time span within which an event's occurrence triggers an alert.
* `scope` - (String) Specifies if the alert is triggered for a user, location, department, or organization.
* `severity` - (String) The threat severity based on which the alert is triggered.
* `comments` - (String) Additional information about the triggered alert.
* `entity` - (List) An immutable reference to the entity the alert is scoped to.
  * `id` - (Integer) A unique identifier for the entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional information about the entity.
