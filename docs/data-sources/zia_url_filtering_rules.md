---
subcategory: "URL Filtering Rule"
layout: "zscaler"
page_title: "ZIA: url_filtering_rules"
description: |-
    Official documentation https://help.zscaler.com/zia/about-url-filtering
    API documentation https://help.zscaler.com/zia/url-filtering-policy#/urlFilteringRules-get
    Gets a list of all of URL Filtering Policy rules.
---

# zia_url_filtering_rules (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-url-filtering)
* [API documentation](https://help.zscaler.com/zia/url-filtering-policy#/urlFilteringRules-post)

Use the **zia_url_filtering_rules** data source to get information about a URL filtering rule information for the specified `Name`.

```hcl
# URL filtering rule
data "zia_url_filtering_rules" "example"{
    name = "Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (String) Name of the URL Filtering policy rule
* `id` - (String) URL Filtering Rule ID

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `order` - (Number) Order of execution of rule with respect to other URL Filtering rules
* `protocols` - (List of Object) Protocol criteria. Supported values: `SMRULEF_ZPA_BROKERS_RULE`, `ANY_RULE`, `TCP_RULE`, `UDP_RULE`, `DOHTTPS_RULE`, `TUNNELSSL_RULE`, `HTTP_PROXY`, `FOHTTP_RULE`, `FTP_RULE`, `HTTPS_RULE`, `HTTP_RULE`, `SSL_RULE`, `TUNNEL_RULE`, `WEBSOCKETSSL_RULE`, `WEBSOCKET_RULE`
* `state` - (String) Rule State
* `rank` - (String) Admin rank of the admin who creates this rule
* `request_methods` - (String) Request method for which the rule must be applied. If not set, rule will be applied to all methods
* `end_user_notification_url` - (String) URL of end user notification page to be displayed when the rule is matched. Not applicable if either 'overrideUsers' or 'overrideGroups' is specified.
* `block_override` - (String) When set to true, a `BLOCK` action triggered by the rule could be overridden. If true and both overrideGroup and overrideUsers are not set, the `BLOCK` triggered by this rule could be overridden for any users. If block_override is not set, `BLOCK` action cannot be overridden.
* `time_quota` - (String) Time quota in minutes, after which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to `BLOCK`, this field is not applicable.
* `size_quota` - (String) Size quota in KB beyond which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to `BLOCK`, this field is not applicable.
* `description` - (String) Additional information about the rule
* `validity_start_time` - (Number) If enforceTimeValidity is set to true, the URL Filtering rule will be valid starting on this date and time.
* `validity_end_time` - (Number) If enforceTimeValidity is set to true, the URL Filtering rule will cease to be valid on this end date and time.
* `validity_time_zone_id` - (Number) If enforceTimeValidity is set to true, the URL Filtering rule date and time will be valid based on this time zone ID.
* `last_modified_time` - (Number) When the rule was last modified
* `enforce_time_validity` - (String) Enforce a set a validity time period for the URL Filtering rule.
* `action` - (String) Action taken when traffic matches rule criteria. Supported values: `ANY`, `NONE`, `BLOCK`, `CAUTION`, `ALLOW`, `ICAP_RESPONSE`
* `source_countries`** - (List of String) Identify destinations based on the location of a server. Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``
* `device_trust_levels` - (List) List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation. Supported values: `ANY`, `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`

* `user_risk_score_levels` (List) - Indicates the user risk score level selectedd for the DLP rule violation: Returned values are: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

* `user_agent_types` (List) - User Agent types on which this rule will be applied: Returned values are: `CHROME`, `FIREFOX`, `MSIE`, `MSEDGE`,   `MSCHREDGE`, `OPERA`, `OTHER`, `SAFARI`

* `cbi_profile` - (List) The cloud browser isolation profile to which the ISOLATE action is applied in the URL Filtering Policy rules. This block is required when the attribute `action` is set to `ISOLATE`
  * `id` - (String) The universally unique identifier (UUID) for the browser isolation profile
  * `name` - (String) Name of the browser isolation profile
  * `url` - (String) The browser isolation profile URL

* `cipa_rule` - (String) If set to true, the CIPA Compliance rule is enabled
* `url_categories` - (String) List of URL categories for which rule must be applied

* `locations` - (List of Object) The locations to which the Firewall Filtering policy rule applies
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `groups` - (List of Object) The groups to which the Firewall Filtering policy rule applies
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `departments` - (List of Object) The departments to which the Firewall Filtering policy rule applies
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `users` - (List of Object) The users to which the Firewall Filtering policy rule applies
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `time_windows` - (List of Object) The time interval in which the Firewall Filtering policy rule applies
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `override_users` - (List of Object) Name-ID pairs of users for which this rule can be overridden. Applicable only if blockOverride is set to `true`, action is `BLOCK` and overrideGroups is not set.If this overrideUsers is not set, `BLOCK` action can be overridden for any user.
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `override_groups` - (List of Object) Name-ID pairs of users for which this rule can be overridden. Applicable only if blockOverride is set to `true`, action is `BLOCK` and overrideGroups is not set.If this overrideUsers is not set, `BLOCK` action can be overridden for any group.
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `location_groups` - (List of Object) The location groups to which the Firewall Filtering policy rule applies
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `labels`
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `last_modified_by`
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)

* `workload_groups` (List) The list of preconfigured workload groups to which the policy must be applied
  * `id` - (Number) A unique identifier assigned to the workload group
  * `name` - (String) The name of the workload group
