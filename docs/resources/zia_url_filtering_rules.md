---
subcategory: "URL Filtering Rule"
layout: "zscaler"
page_title: "ZIA:  url_filtering_rules"
description: |-
       Creates and manages a URL Filtering Policy rule.
---
# Resource: zia_url_filtering_rules

The **zia_url_filtering_rules** resource creates and manages a URL filtering rules within the Zscaler Internet Access cloud.

## Example Usage - BLOCK ACTION

```hcl
resource "zia_url_filtering_rules" "this" {
    name                = "Example"
    description         = "Example"
    state               = "ENABLED"
    action              = "BLOCK"
    order               = 1
    url_categories      = ["ANY"]
    device_trust_levels = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    protocols           = ["ANY_RULE"]
    request_methods     = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}
```

## Example Usage - ISOLATE ACTION

⚠️ **WARNING 1:**: Creating a URL Filtering rule with the action of `ISOLATE` requires the Cloud Browser Isolation subscription. To learn more, contact Zscaler Support or your local account team.

```hcl
data "zia_cloud_browser_isolation_profile" "this" {
    name = "BD_SA_Profile1_ZIA"
}

resource "zia_url_filtering_rules" "this" {
    name                = "Example"
    description         = "Example"
    state               = "ENABLED"
    action              = "ISOLATE"
    order               = 1
    url_categories      = ["ANY"]
    device_trust_levels = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    protocols           = [ "HTTPS_RULE", "HTTP_RULE" ]
    request_methods     = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE" ]
    cbi_profile {
        id = data.zia_cloud_browser_isolation_profile.this.id
        name = data.zia_cloud_browser_isolation_profile.this.name
        url = data.zia_cloud_browser_isolation_profile.this.url
    }
    user_agent_types = [ "OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE" ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Firewall Filtering policy rule
* `order` - (Required) Order of execution of rule with respect to other URL Filtering rules
* `protocols` - (List of Object) Protocol criteria. Supported values: `SMRULEF_ZPA_BROKERS_RULE`, `ANY_RULE`, `TCP_RULE`, `UDP_RULE`, `DOHTTPS_RULE`, `TUNNELSSL_RULE`, `HTTP_PROXY`, `FOHTTP_RULE`, `FTP_RULE`, `HTTPS_RULE`, `HTTP_RULE`, `SSL_RULE`, `TUNNEL_RULE`, `WEBSOCKETSSL_RULE`, `WEBSOCKET_RULE`,

### Optional

* `request_methods` - (Optional) Request method for which the rule must be applied. If not set, rule will be applied to all methods
* `rank` - (Optional) Admin rank of the admin who creates this rule
* `state` - (Optional) Rule State
* `end_user_notification_url` - (Optional) URL of end user notification page to be displayed when the rule is matched. Not applicable if either 'overrideUsers' or 'overrideGroups' is specified.
* `block_override` - (Optional) When set to true, a `BLOCK` action triggered by the rule could be overridden. If true and both overrideGroup and overrideUsers are not set, the `BLOCK` triggered by this rule could be overridden for any users. If block_override is not set, `BLOCK` action cannot be overridden.
* `time_quota` - (Optional) Time quota in minutes, after which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to `BLOCK`, this field is not applicable.
* `size_quota` - (Optional) Size quota in KB beyond which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to `BLOCK`, this field is not applicable.
* `description` - (Optional) Additional information about the rule
* `validity_start_time` - (Optional) If enforceTimeValidity is set to true, the URL Filtering rule will be valid starting on this date and time.
* `validity_end_time` - (Optional) If `enforceTimeValidity` is set to true, the URL Filtering rule will cease to be valid on this end date and time.
* `validity_time_zone_id` - (Optional) If `enforceTimeValidity` is set to true, the URL Filtering rule date and time will be valid based on this time zone ID.
* `last_modified_time` - (Optional) When the rule was last modified
* `enforce_time_validity` - (Optional) Enforce a set a validity time period for the URL Filtering rule.
* `action` - (Optional) Action taken when traffic matches rule criteria. Supported values: `ANY`, `NONE`, `BLOCK`, `CAUTION`, `ALLOW`, `ICAP_RESPONSE`
* `device_trust_levels` - (Optional) List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation. Supported values: `ANY`, `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`

* `user_risk_score_levels` (Optional) - Indicates the user risk score level selectedd for the DLP rule violation: Returned values are: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

* `cbi_profile` - (Optional) The cloud browser isolation profile to which the ISOLATE action is applied in the URL Filtering Policy rules. This block is required when the attribute `action` is set to `ISOLATE`
  * `id` - (Optional) The universally unique identifier (UUID) for the browser isolation profile
  * `name` - (Optional) Name of the browser isolation profile
  * `url` - (Optional) The browser isolation profile URL

* `cipa_rule` - (Optional) If set to true, the CIPA Compliance rule is enabled
* `url_categories` - (Optional) List of URL categories for which rule must be applied

* `locations` - (List of Object) The locations to which the Firewall Filtering policy rule applies
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `groups` - (List of Object) The groups to which the Firewall Filtering policy rule applies
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `departments` - (List of Object) The departments to which the Firewall Filtering policy rule applies
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `users` - (List of Object) The users to which the Firewall Filtering policy rule applies
  * `id` - (Number) Identifier that uniquely identifies an entity

* `time_windows` - (List of Object) The time interval in which the Firewall Filtering policy rule applies
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `override_users` - (List of Object) Name-ID pairs of users for which this rule can be overridden. Applicable only if blockOverride is set to `true`, action is `BLOCK` and overrideGroups is not set.If this overrideUsers is not set, `BLOCK` action can be overridden for any user.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `override_groups` - (List of Object) Name-ID pairs of users for which this rule can be overridden. Applicable only if blockOverride is set to `true`, action is `BLOCK` and overrideGroups is not set.If this overrideUsers is not set, `BLOCK` action can be overridden for any group.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `location_groups` - (List of Object) The location groups to which the Firewall Filtering policy rule applies
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `labels`
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `workload_groups` (Optional) The list of preconfigured workload groups to which the policy must be applied
  * `id` - (Optional) A unique identifier assigned to the workload group
  * `name` - (Optional) The name of the workload group
