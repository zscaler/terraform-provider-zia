---
subcategory: "Security & UEBA Alerts"
layout: "zscaler"
page_title: "ZIA: ueba_alert_definitions"
description: |-
  Official documentation https://help.zscaler.com/zia/about-alerts
  API documentation https://help.zscaler.com/legacy-apis/security-ueba-alerts#/alertDefinitions-post
  Creates and manages a ZIA Security & UEBA alert definition.
---

# zia_ueba_alert_definitions (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-alerts)
* [API documentation](https://help.zscaler.com/legacy-apis/security-ueba-alerts#/alertDefinitions-post)

Use the **zia_ueba_alert_definitions** resource to create and manage Security & UEBA alert rules. An alert definition specifies the threat or event type to monitor, the severity that triggers the alert, and the scope (user, location, department, or organization) it applies to.

## Example Usage - Organization Scope

```hcl
resource "zia_ueba_alert_definitions" "this" {
  alert_name = "OUTGOING_VIRUSES"
  status     = "ENABLED"
  scope      = "ORGANIZATION"
  occurrence = "OCCURRENCE_1"
  interval   = "INTERVAL_5_MINUTES"
  severity   = "CRITICAL"
  comments   = "Alert on outgoing virus detections"
}
```

## Example Usage - Traffic-Based Alert

```hcl
resource "zia_ueba_alert_definitions" "traffic" {
  alert_name             = "TRAFFIC_INCREASE"
  status                 = "ENABLED"
  scope                  = "ORGANIZATION"
  interval               = "INTERVAL_1_HOUR"
  traffic_change_percent = 50
  severity               = "MAJOR"
  comments               = "Alert when traffic increases by 50%"
}
```

## Example Usage - Scoped to an Entity

```hcl
data "zia_department_management" "this" {
  name = "A000"
}

resource "zia_ueba_alert_definitions" "user_scope" {
  alert_name = "PHISHING"
  status     = "ENABLED"
  scope      = "USER"
  occurrence = "OCCURRENCE_5"
  interval   = "INTERVAL_15_MINUTES"
  severity   = "CRITICAL"

  entity {
    id = data.zia_department_management.this.id
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `alert_name` - (String) The alert name that identifies the threat or event type the alert is triggered for. Supported values are: `LDAP_SUCCESS`, `LDAP_FAILURE`, `LDAP_CONNECTION_DOWN`, `AUTH_BRIDGE_DOWN`, `ADP_SCHEDULE_UPDATE_FAILURE`, `IDM_SCHEDULE_UPDATE_FAILURE`, `EDM_SCHEMA_INDEXING_FAILURE`, `OUTGOING_VIRUSES`, `OUTGOING_SPYWARE`, `OUTGOING_MALWARE`, `OUTGOING_UNSCANNABLE_FILES`, `INCOMING_VIRUSES`, `INCOMING_SPYWARE`, `INCOMING_MALWARE`, `INCOMING_UNSCANNABLE_FILES`, `BOTNET`, `MALICIOUS_CONTENT`, `PHISHING`, `URL_FILTERING_BLOCKED_SITES`, `STREAMING_UPLOAD`, `STREAMING_VIEW_LISTEN`, `SOCIAL_NETWORK_POST`, `CHAT_FILE_TRANSFER`, `WEBMAIL_FILE_ATTACHMENT`, `HIPAA_VIOLATION`, `PCI_VIOLATION`, `GLBA_VIOLATION`, `CUSTOM_ENGINE_VIOLATION`, `POLICY_VIOLATION`, `BA_ADWARE`, `BA_MALWARE`, `BA_ANONYMIZER`, `ADVANCED_SECURITY`, `PEER_TO_PEER`, `UNAUTH_COMM`, `CROSS_SITE_SCRIPTING`, `BROWSER_EXPLOIT`, `SUSPICIOUS_DESTINATION`, `ADWARE_SPYWARE`, `WEB_SPAM`, `PAGE_RISK`, `ADSPYWARE_SITES`, `CRYPTOMINING`, `TRAFFIC_DECREASE`, `TRAFFIC_INCREASE`, `DGA_DOMAINS`, `BA_PATIENT0`.

### Optional

* `status` - (String) The status of the alert rule. Supported values are `ENABLED` and `DISABLED`.
* `occurrence` - (String) Specifies the occurrence of an ongoing alert for a specific threat or event type. Supported values are `OCCURRENCE_1`, `OCCURRENCE_5`, `OCCURRENCE_10`, `OCCURRENCE_100`, and `OCCURRENCE_1000`.
* `traffic_change_percent` - (Integer) Specifies the percentage change in traffic. Applicable to traffic-based alerts such as `TRAFFIC_INCREASE` and `TRAFFIC_DECREASE`.
* `interval` - (String) The time span within which an event's occurrence triggers an alert. Supported values are `INTERVAL_5_MINUTES`, `INTERVAL_15_MINUTES`, `INTERVAL_30_MINUTES`, `INTERVAL_1_HOUR`, and `INTERVAL_1_DAY`.
* `scope` - (String) Specifies if the alert is triggered for a user, location, department, or organization. Supported values are `USER`, `LOCATION`, `DEPARTMENT`, and `ORGANIZATION`.
* `severity` - (String) The threat severity based on which the alert is triggered. Supported values are `CRITICAL`, `MAJOR`, `MINOR`, `INFO`, and `DEBUG`.
* `comments` - (String) Additional information about the triggered alert.
* `entity` - (Block) An immutable reference to the entity (user, location, or department) the alert is scoped to. Applicable ONLY when `scope` is NOT `ORGANIZATION`. Use the data sources `zia_department_management`, `zia_user_management`, `zia_location_management`
  * `id` - (Integer) A unique identifier for the entity.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The Terraform identifier of the alert definition.
* `alert_definition_id` - (Integer) The system-generated identifier of the alert definition.

## Import

Alert definitions can be imported by using `<ID>` or `<ALERT_NAME>` as the import ID.

For example:

```shell
terraform import zia_ueba_alert_definitions.this 123456
```

```shell
terraform import zia_ueba_alert_definitions.this OUTGOING_VIRUSES
```
