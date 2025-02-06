---
subcategory: "SSL Inspection Rules"
layout: "zscaler"
page_title: "ZIA: ssl_inspection_rules"
description: |-
  Creates and manages SSL Inspection Rules.
---

# Resource: zia_ssl_inspection_rules

The **zia_ssl_inspection_rules** resource allows the creation and management of SSL Inspection rules in the Zscaler Internet Access.

**NOTE 1** Zscaler SSL Inspection rules contain default and predefined rules which are placed in their respective orders. These rules `CANNOT` be deleted. When configuring your rules make sure that the `order` attributue value consider these pre-existing rules so that Terraform can place the new rules in the correct position, and drifts can be avoided. i.e If there are 2 pre-existing rules, you should start your rule order at `3` and manage your rule sets from that number onwards. The provider will reorder the rules automatically while ignoring the order of pre-existing rules, as the API will be responsible for moving these rules to their respective positions as API calls are made.

The most common default and predefined rules:

|              Rule Names                      |  Default or Predefined   |   Rule Number Associated |
|:--------------------------------------------:|:------------------------:|:------------------------:|
|-----------------------------|--------------------------|-------------------|
|  `Zscaler Recommended Exemptions`                 |      `Predefined`       |           `Yes`          |
|  `Office 365 One Click`             |      `Predefined`       |           `Yes`          |
|  `Office365 Inspection`             |      `Predefined`       |           `Yes`          |
|  `UCaaS One Click`             |      `Predefined`       |           `Yes`          |
|  `Default SSL Inspection Rule`             |      `Default`       |           `No`          |
|-------------------------|-------------------------|-----------------|

**NOTE 2** Certain attributes on `predefined` rules can still be managed or updated via Terraform such as:

- `description` - (Optional) Enter additional notes or information. The description cannot exceed 10,240 characters.
- `state` - (Optional) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to
- `labels` (list) - Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

**NOTE 3** The import of `predefined` rules is still possible in case you want o have them under the Terraform management; however, remember that these rules cannot be deleted. That means, the provider will fail when executing `terraform destroy`; hence, you must remove the rules you want to delete, and re-run `terraform apply` instead.

## Example Usage - Action - DECRYPT

```hcl

data "zia_group_management" "this" {
    name = "A000"
}

resource "zia_ssl_inspection_rules" "this" {
  name                         = "SSL_Inspection_Rule_Decrypt"
  description                  = "SSL_Inspection_Rule_Decrypt"
  state                        = "ENABLED"
  order                        = 1
  rank                         = 7
  road_warrior_for_kerberos    = true
  cloud_applications           = ["CHATGPT_AI", "ANDI"]
  platforms                    = ["SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX"]

  action {
    type                         = "DECRYPT"
    # show_eun                   = false
    # show_eunatp                = false
    override_default_certificate = false

    ssl_interception_cert {
      id                  = 1
      name                = "Zscaler Intermediate CA Certificate"
      default_certificate = true
    }

    decrypt_sub_actions {
      server_certificates                   = "ALLOW"
      ocsp_check                            = true
      block_ssl_traffic_with_no_sni_enabled = true
      min_client_tls_version                = "CLIENT_TLS_1_0"
      min_server_tls_version                = "SERVER_TLS_1_0"
      block_undecrypt                       = true
      http2_enabled                         = false
    }
  }
  groups {
        id = [ data.zia_group_management.this.id ]
    }
}
```

## Example Usage - Action - DO_NOT_DECRYPT - Bypass Rule (False)

```hcl

data "zia_group_management" "this" {
    name = "A000"
}

resource "zia_ssl_inspection_rules" "this" {
  name                         = "SSL_Rule_Do_Not_Decrypt"
  description                  = "SSL_Rule_Do_Not_Decrypt"
  state                        = "ENABLED"
  order                        = 1
  rank                         = 7
  road_warrior_for_kerberos    = true
  cloud_applications           = ["CHATGPT_AI", "ANDI"]
  platforms                    = ["SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX"]

  action {
    type                                    = "DO_NOT_DECRYPT"
    do_not_decrypt_sub_actions {
      bypass_other_policies                 = false
      server_certificates                   = "ALLOW"
      ocsp_check                            = true
      block_ssl_traffic_with_no_sni_enabled = true
      min_tls_version                       = "SERVER_TLS_1_0"
    }
  }
  groups {
        id = [ data.zia_group_management.this.id ]
    }
}
```

## Example Usage - Action - DO_NOT_DECRYPT - Bypass Rule (True)

```hcl

data "zia_group_management" "this" {
    name = "A000"
}

resource "zia_ssl_inspection_rules" "this" {
  name                         = "SSL_Rule_Bypass_Rule"
  description                  = "SSL_Rule_Bypass_Rule"
  state                        = "ENABLED"
  order                        = 1
  rank                         = 7
  road_warrior_for_kerberos    = true
  cloud_applications           = ["CHATGPT_AI", "ANDI"]
  platforms                    = ["SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX"]

  action {
    type                                    = "DO_NOT_DECRYPT"
    do_not_decrypt_sub_actions {
      bypass_other_policies                 = true
      block_ssl_traffic_with_no_sni_enabled = true
    }
  }
  groups {
        id = [ data.zia_group_management.this.id ]
    }
}
```

## Example Usage - Action - BLOCK

```hcl

data "zia_group_management" "this" {
    name = "A000"
}

resource "zia_ssl_inspection_rules" "this" {
  name                         = "SSL_Rule_BLOCK"
  description                  = "SSL_Rule_BLOCK"
  state                        = "ENABLED"
  order                        = 1
  rank                         = 7
  road_warrior_for_kerberos    = true
  cloud_applications           = ["CHATGPT_AI", "ANDI"]
  platforms                    = ["SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX"]

  action {
    type                                    = "BLOCK"
    ssl_interception_cert {
      id                                    = 1
    }
  }
  groups {
        id = [ data.zia_group_management.this.id ]
    }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (String) Name of the SSL Inspection
* `order` - (String) Unique identifier for the SSL Inspection

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

* `description` (String) -  Enter additional notes or information. The description cannot exceed 10,240 characters.
* `order` (String) -  Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.
* `state` (String) - The state of the rule indicating whether it is enabled or disabled. Supported values: `ENABLED` or `DISABLED`
* `rank` (Integer) - The admin rank specified for the rule based on your assigned admin rank. Admin rank determines the rule order that can be specified for the rule. Admin rank can be configured if it is enabled in the Advanced Settings.
* `access_control` (String) - The access privilege (RBA) for this rule.
* `road_warrior_for_kerberos` (Boolean) - Indicates whether this rule is applied to remote users that use PAC with Kerberos authentication.
* `platforms` (Set of String) -  Zscaler Client Connector device platforms for which this rule is applied. Supported Values: `SCAN_IOS`, `SCAN_ANDROID`, `SCAN_MACOS`, `SCAN_WINDOWS`, `NO_CLIENT_CONNECTOR`, `SCAN_LINUX`
* `cloud_applications` (Set of String) -  The list of URL categories to which the DLP policy rule must be applied. For the complete list of supported file types refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post)
* `url_categories` (Set of String) -  The list of URL categories to which the DLP policy rule must be applied.
* `user_agent_types` (Set of String) -  A list of user agent types the rule applies to.
* `device_trust_levels` (Set of String)  - List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation. Supported values: `ANY`, `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`
* `action` (Block List) - Action taken when the traffic matches policy
* `devices` (Block List) - ID pairs of devices for which the rule is applied
* `device_groups` (Block List) - ID pairs of device groups for which the rule is applied.
* `departments` (Block List) - ID pairs of departments for which the rule is applied.
* `groups` (Block List) - ID pairs of groups for which the rule is applied. If not set, rule is applied for all groups.
* `labels` (Block List) - ID pairs of labels associated with the rule.
* `locations` (Block List) - ID pairs of locations to which the rule is applied. When empty, it implies applying to all locations.
* `location_groups` (Block List) - ID pairs of location groups to which the rule is applied. When empty, it implies applying to all location groups.
* `dest_ip_groups` (Block List) - ID pairs of destination IP address groups for which the rule is applied.
* `source_ip_groups` (Block List) - ID pairs of source IP address groups for which the rule is applied.
* `proxy_gateways` (Block List) - When using ZPA Gateway forwarding, name-ID pairs of ZPA Application Segments for which the rule is applicable.
* `zpa_app_segments` (Block List) - The list of ZPA Application Segments for which this rule is applicable (applicable only for ZPA Gateway forwarding).
* `workload_groups` (Block List) - The list of preconfigured workload groups to which the policy must be applied.
* `time_windows` (Block List) - The time intervals during which the rule applies
* `users` (Block List) - The list of preconfigured workload groups to which the policy must be applied.

### Action Attributes

`action` has the following attributes:

* `type` (String) - The action type for this rule. Possible values: `BLOCK`, `DECRYPT`, or `DO_NOT_DECRYPT`.
* `show_eun` (Boolean) - Whether to show End User Notification (EUN).
* `show_eunatp` (Boolean) - Whether to display the EUN ATP page.
* `override_default_certificate` (Boolean) - Whether to override the default SSL interception certificate.
* `ssl_interception_cert` (Block List) - The SSL interception certificate to be used. If not set it will use the default Zscaler certificate
* `decrypt_sub_actions` (Block List) - Action taken when enabling SSL intercept
* `do_not_decrypt_sub_actions` (Block List) - Action taken when bypassing SSL intercept

### ssl_interception_cert Attributes

`ssl_interception_cert` has the following attributes:

* `id` (Integer) - The unique ID of the SSL interception certificate.
* `name` (String) - The name of the SSL interception certificate.
* `default_certificate` (Boolean) - Indicates if this certificate is the default certificate.

### decrypt_sub_actions Attributes

`decrypt_sub_actions` has the following attributes:

* `server_certificates` (String) - Action to take on server certificates. Valid values might include `ALLOW`, `BLOCK`, or `PASS_THRU`.
* `ocsp_check` (Boolean) - Whether to enable OCSP check.
* `block_ssl_traffic_with_no_sni_enabled` (Boolean) - Whether to block SSL traffic when SNI is not present.
* `min_client_tls_version` (String) - The minimum TLS version allowed on the client side: Supported Values are: `CLIENT_TLS_1_0`, `CLIENT_TLS_1_1`, `CLIENT_TLS_1_2`,  `CLIENT_TLS_1_3`.
* `min_server_tls_version` (String) - The minimum TLS version allowed on the server side: Supported Values are: `SERVER_TLS_1_0`, `SERVER_TLS_1_1`, `SERVER_TLS_1_2`,  `SERVER_TLS_1_3`.
* `block_undecrypt` (Boolean) - Enable to block traffic from servers that use non-standard encryption methods or require mutual TLS authentication.
* `http2_enabled` (Boolean)

### do_not_decrypt_sub_actions Attributes

`do_not_decrypt_sub_actions` has the following attributes:

* `bypass_other_policies` (Boolean) - Whether to bypass other policies when action is set to `DO_NOT_DECRYPT`.
* `server_certificates` (String) - Action to take on server certificates. Valid values might include `ALLOW`, `BLOCK`, or `PASS_THRU`.
* `ocsp_check` (Boolean) - Whether to enable OCSP check.
* `block_ssl_traffic_with_no_sni_enabled` (Boolean) - Whether to block SSL traffic when SNI is not present.
* `min_tls_version` (String) -  The minimum TLS version allowed on the server side: Supported Values are: `SERVER_TLS_1_0`, `SERVER_TLS_1_1`, `SERVER_TLS_1_2`,  `SERVER_TLS_1_3`.

### Devices Attributes

* `id` (Integer) - A unique identifier for the device.

### Device Groups Attributes

* `id` (Integer) - A unique identifier for the device groups.

### Labels Attributes

* `id` (Integer) - A unique identifier for the label.

### Locations Attributes

* `id` (Integer) - A unique identifier for the locations.

### Location Groups Attributes

* `id` (Integer) - A unique identifier for the location groups.

### Departments Attributes

* `id` (Integer) - A unique identifier for the departments.

### Destination IP Groups Attributes

* `id` (Integer) - A unique identifier for the destination ip group.

### Groups Attributes

* `id` (Integer) - A unique identifier for the groups.

### Source IP Groups Attributes

* `id` (Integer) - A unique identifier for the source ip group.

### Users Attributes

* `id` (Integer) - A unique identifier for the users.

### Time Windows Attributes

* `id` (Integer) - A unique identifier for the time window.

### Proxy Gateways Attributes

* `id` (Integer) - A unique identifier assigned to the Application Segment

### ZPA App Segments Attributes

* `id` (Integer) - A unique identifier assigned to the Application Segment

### Workload Groups Attributes

* `id` (Integer) - A unique identifier assigned to the workload group
