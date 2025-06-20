---
subcategory: "SSL Inspection Rules"
layout: "zscaler"
page_title: "ZIA: ssl_inspection_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/about-ssl-inspection-policy
  API documentation https://help.zscaler.com/zia/ssl-inspection-policy#/sslInspectionRules-get
  Retrieves the list of all SSL Inspection rules configured in the ZIA Admin Portal.
---

# zia_ssl_inspection_rules (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-ssl-inspection-policy)
* [API documentation](https://help.zscaler.com/zia/ssl-inspection-policy#/sslInspectionRules-get)

Use the **zia_ssl_inspection_rules** data source to get information about a ssl inspection rule in the Zscaler Internet Access.

## Example Usage

```hcl
# ZIA SSL Inspection by name
data "zia_ssl_inspection_rules" "this" {
    name = "SSL_Inspection_Rule01"
}
```

```hcl
# ZIA SSL Inspection by ID
data "zia_ssl_inspection_rules" "this" {
    id = "12365478"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the SSL Inspection
* `id` - (Optional) Unique identifier for the SSL Inspection

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` -  Enter additional notes or information. The description cannot exceed 10,240 characters.
* `order` -  Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.
* `state` - The state of the rule indicating whether it is enabled or disabled. Supported values: `ENABLED` or `DISABLED`
* `rank` - The admin rank specified for the rule based on your assigned admin rank. Admin rank determines the rule order that can be specified for the rule. Admin rank can be configured if it is enabled in the Advanced Settings.
* `access_control` - The access privilege (RBA) for this rule.
* `road_warrior_for_kerberos` - Indicates whether this rule is applied to remote users that use PAC with Kerberos authentication.
* `platforms` -  Zscaler Client Connector device platforms for which this rule is applied. Supported Values: `SCAN_IOS`, `SCAN_ANDROID`, `SCAN_MACOS`, `SCAN_WINDOWS`, `NO_CLIENT_CONNECTOR`, `SCAN_LINUX`
* `cloud_applications` -  The list of URL categories to which the DLP policy rule must be applied.
* `url_categories` -  The list of URL categories to which the DLP policy rule must be applied.
* `user_agent_types` -  A list of user agent types the rule applies to.
* `device_trust_levels` -  Lists device trust levels for which the rule must be applied (for devices managed using Zscaler Client Connector).
* `action` - Action taken when the traffic matches policy
* `devices` - ID pairs of devices for which the rule is applied
* `device_groups` - ID pairs of device groups for which the rule is applied.
* `departments` - ID pairs of departments for which the rule is applied.
* `groups` - ID pairs of groups for which the rule is applied. If not set, rule is applied for all groups.
* `labels` - ID pairs of labels associated with the rule.
* `locations` - ID pairs of locations to which the rule is applied. When empty, it implies applying to all locations.
* `location_groups` - ID pairs of location groups to which the rule is applied. When empty, it implies applying to all location groups.
* `dest_ip_groups` - ID pairs of destination IP address groups for which the rule is applied.
* `source_ip_groups` - ID pairs of source IP address groups for which the rule is applied.
* `proxy_gateways` - When using ZPA Gateway forwarding, name-ID pairs of ZPA Application Segments for which the rule is applicable.
* `zpa_app_segments` - The list of ZPA Application Segments for which this rule is applicable (applicable only for ZPA Gateway forwarding).
* `workload_groups` - The list of preconfigured workload groups to which the policy must be applied.
* `time_windows` - The time intervals during which the rule applies
* `users` - The list of preconfigured workload groups to which the policy must be applied.

### Action Attributes

`action` has the following attributes:

* `type` - The action type for this rule. Possible values: `BLOCK`, `DECRYPT`, or `DO_NOT_DECRYPT`.
* `show_eun` - Whether to show End User Notification (EUN).
* `show_eunatp` - Whether to display the EUN ATP page.
* `override_default_certificate` - Whether to override the default SSL interception certificate.
* `ssl_interception_cert` - Action taken when enabling SSL intercept
* `do_not_decrypt_sub_actions` - Action taken when bypassing SSL intercept

### ssl_interception_cert Attributes

`ssl_interception_cert` has the following attributes:

* `id` - The unique ID of the SSL interception certificate.
* `name` - The name of the SSL interception certificate.
* `default_certificate` - Indicates if this certificate is the default certificate.

### do_not_decrypt_sub_actions Attributes

`do_not_decrypt_sub_actions` has the following attributes:

* `bypass_other_policies` - Whether to bypass other policies when action is set to `DO_NOT_DECRYPT`.
* `server_certificates` - Action to take on server certificates. Valid values might include `ALLOW`, `BLOCK`, or `PASS_THRU`.
* `ocsp_check` - Whether to enable OCSP check.
* `block_ssl_traffic_with_no_sni_enabled` - Whether to block SSL traffic when SNI is not present.
* `min_tls_version` - The minimum TLS version allowed when action is `DO_NOT_DECRYPT`.

### Devices Attributes

* `id` - A unique identifier for the device.
* `name` - The name of the device.
* `extensions` - Additional information about the device.

### Device Groups Attributes

* `id` - A unique identifier for the device groups.
* `name` - The name of the device groups.
* `extensions` - Additional information about the device groups.

### Labels Attributes

* `id` - A unique identifier for the label.
* `name` - The name of the label.
* `extensions` - Additional information about the label.

### Locations Attributes

* `id` - A unique identifier for the locations.
* `name` - The name of the locations.
* `extensions` - Additional information about the locations.

### Location Groups Attributes

* `id` - A unique identifier for the location groups.
* `name` - The name of the location groups.
* `extensions` - Additional information about the location groups.

### Departments Attributes

* `id` - A unique identifier for the departments.
* `name` - The name of the departments.
* `extensions` - Additional information about the departments.

### Destination IP Groups Attributes

* `id` - A unique identifier for the destination ip group.
* `name` - The name of the destination ip group.
* `extensions` - Additional information about the destination ip group.

### Groups Attributes

* `id` - A unique identifier for the groups.
* `name` - The name of the groups.
* `extensions` - Additional information about the groups.

### Source IP Groups Attributes

* `id` - A unique identifier for the source ip group.
* `name` - The name of the source ip group.
* `extensions` - Additional information about the source ip group.

### Users Attributes

* `id` - A unique identifier for the users.
* `name` - The name of the users.
* `extensions` - Additional information about the users.

### Time Windows Attributes

* `id` - A unique identifier for the time window.
* `name` - The name of the time window.
* `extensions` - Additional information about the time window.

### Proxy Gateways Attributes

* `id` - A unique identifier assigned to the Application Segment
* `name` - The name of the Application Segment
* `external_id` - Indicates the external ID. Applicable only when this reference is of an external entity.

### ZPA App Segments Attributes

* `id` - A unique identifier assigned to the Application Segment
* `name` - The name of the Application Segment
* `external_id` - Indicates the external ID. Applicable only when this reference is of an external entity.

### Workload Groups Attributes

* `id` - A unique identifier assigned to the workload group
* `name` - The name of the workload group
* `description` - The description of the workload group
* `expression` - The expression used within the workload group.
* `last_modified_time` - Timestamp when the workload group was last modified.
* `last_modified_by` - A nested block with details about who last modified the workload group.
* `expression_json` - A nested block describing the JSON expression for the workload group.

### Workload Groups Expression JSON Attributes

* `expression_containers` - Contains one or more tag types (and associated tags) combined using logical operators within a workload group

### Workload Groups Expression Containers Attributes

* `tag_type` - The tag type selected from a predefined list
* `operator` - The operator (either AND or OR) used to create logical relationships among tag types
* `tag_container` - Contains one or more tags and the logical operator used to combine the tags within a tag type
