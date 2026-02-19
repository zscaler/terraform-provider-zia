---
subcategory: "URL Filtering Rules"
layout: "zscaler"
page_title: "ZIA:  url_filtering_rules"
description: |-
    Official documentation https://help.zscaler.com/zia/about-url-filtering
    API documentation https://help.zscaler.com/zia/url-filtering-policy#/urlFilteringRules-get
    Creates and manages a URL Filtering Policy rule.
---

# zia_url_filtering_rules (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-url-filtering)
* [API documentation](https://help.zscaler.com/zia/url-filtering-policy#/urlFilteringRules-post)

The **zia_url_filtering_rules** resource creates and manages a URL filtering rules within the Zscaler Internet Access cloud.

## Example Usage - ALLOW ACTION

```hcl
resource "zia_url_filtering_rules" "this" {
    name                  = "Example"
    description           = "Example"
    state                 = "ENABLED"
    action                = "ALLOW"
    order                 = 1
    enforce_time_validity = true
    validity_start_time   = "Mon, 17 Jun 2024 23:30:00 UTC"
    validity_end_time     = "Tue, 17 Jun 2025 23:00:00 UTC"
    validity_time_zone_id = "US/Pacific"
    time_quota            = 15
    size_quota            = 10
    url_categories        = ["ANY"]
    device_trust_levels   = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    protocols             = ["ANY_RULE"]
    request_methods       = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
    user_agent_types      = ["OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"]
}
```

## Example Usage - BLOCK ACTION

```hcl
resource "zia_url_filtering_rules" "this" {
    name                  = "Example"
    description           = "Example"
    state                 = "ENABLED"
    action                = "BLOCK"
    order                 = 1
    enforce_time_validity = true
    validity_start_time   = "Mon, 17 Jun 2024 23:30:00 UTC"
    validity_end_time     = "Tue, 17 Jun 2025 23:00:00 UTC"
    validity_time_zone_id = "US/Pacific"
    time_quota = 15
    size_quota = 10
    url_categories        = ["ANY"]
    device_trust_levels   = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    protocols             = ["ANY_RULE"]
    request_methods       = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
    user_agent_types      = ["OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"]
    block_override        = true
    override_users {
      id = [ 45513075 ]
    }
    override_groups {
      id = [ 76662385 ]
    }
}
```

## Example Usage - CAUTION ACTION

```hcl
resource "zia_url_filtering_rules" "this" {
    name                  = "Example"
    description           = "Example"
    state                 = "ENABLED"
    action                = "CAUTION"
    order                 = 1
    enforce_time_validity = true
    validity_start_time   = "Mon, 17 Jun 2024 23:30:00 UTC"
    validity_end_time     = "Tue, 17 Jun 2025 23:00:00 UTC"
    validity_time_zone_id = "US/Pacific"
    time_quota            = 15
    size_quota            = 10
    url_categories        = ["ANY"]
    device_trust_levels   = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    protocols             = ["ANY_RULE"]
    request_methods       = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
    user_agent_types      = ["OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"]
    end_user_notification_url = "https://caution.acme.com"
}
```

## Example Usage - ISOLATE ACTION

⚠️ **WARNING 1:**: Creating a URL Filtering rule with the action of `ISOLATE` requires the Cloud Browser Isolation subscription. To learn more, contact Zscaler Support or your local account team.

```hcl
data "zia_cloud_browser_isolation_profile" "this" {
    name = "BD_SA_Profile1_ZIA"
}

resource "zia_url_filtering_rules" "this" {
    name                  = "Example"
    description           = "Example"
    state                 = "ENABLED"
    action                = "ISOLATE"
    order                 = 1
    enforce_time_validity = true
    validity_start_time   = "Mon, 17 Jun 2024 23:30:00 UTC"
    validity_end_time     = "Tue, 17 Jun 2025 23:00:00 UTC"
    validity_time_zone_id = "US/Pacific"
    time_quota            = 15
    size_quota            = 10
    url_categories        = ["ANY"]
    device_trust_levels   = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    protocols             = [ "HTTPS_RULE", "HTTP_RULE" ]
    request_methods       = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE" ]
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
Supported values: `OPTIONS`, `GET`, `HEAD`, `POST`, `PUT`, `DELETE`, `TRACE`, `CONNECT`, `OTHER`, `PROPFIND`, `PROPPATCH`, `MOVE`, `MKCOL`, `LOCK`, `COPY`, `UNLOCK`, `PATCH`,, `HTTPS_RULE`, `HTTP_RULE`, `SSL_RULE`, `TUNNEL_RULE`, `WEBSOCKETSSL_RULE`, `WEBSOCKET_RULE`,
* `rank` - (Optional) Admin rank of the admin who creates this rule
* `state` - (Optional) Rule State
* `end_user_notification_url` - (Optional) URL of end user notification page to be displayed when the rule is matched. Not applicable if either 'overrideUsers' or 'overrideGroups' is specified.
* `block_override` - (Optional) When set to true, a `BLOCK` action triggered by the rule could be overridden. If true and both overrideGroup and overrideUsers are not set, the `BLOCK` triggered by this rule could be overridden for any users. If block_override is not set, `BLOCK` action cannot be overridden.
* `time_quota` - (Optional) Time quota in minutes, after which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to `BLOCK`, this field is not applicable.
* `size_quota` - (Optional) Size quota in MB beyond which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to `BLOCK`, this field is not applicable.
* `description` - (Optional) Additional information about the rule
* `validity_start_time` - (Optional) If enforce_time_validity is set to true, the URL Filtering rule will be valid starting on this date and time. The date and time must be provided in `RFC1123` format i.e `Sun, 16 Jun 2024 15:04:05 UTC`
* `validity_end_time` - (Optional) If `enforce_time_validity` is set to true, the URL Filtering rule will cease to be valid on this end date and time. The date and time must be provided in `RFC1123` format i.e `Sun, 16 Jun 2024 15:04:05 UTC`

  **NOTE** Notice that according to RFC1123 the day must be provided as a double digit value for `validity_start_time` and `validity_end_time` i.e `01`, `02` etc.

* `validity_time_zone_id` - (Optional) If `enforce_time_validity` is set to true, the URL Filtering rule date and time will be valid based on this time zone ID. The attribute is validated against the official [IANA List](https://nodatime.org/TimeZones)

* `enforce_time_validity` - (Optional) Enforce a set a validity time period for the URL Filtering rule.
* `action` - (Optional) Action taken when traffic matches rule criteria. Supported values: `BLOCK`, `CAUTION`, `ALLOW`, `ICAP_RESPONSE`

* `devices` (list) - Specifies devices that are managed using Zscaler Client Connector.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `device_groups` (list) - This field is applicable for devices that are managed using Zscaler Client Connector.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `device_trust_levels` - (Optional) List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation. Supported values: `ANY`, `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`

* `user_risk_score_levels` (Optional) - Indicates the user risk score level selectedd for the DLP rule violation: Returned values are: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

* `user_agent_types` (Optional) - User Agent types on which this rule will be applied: Returned values are: `CHROME`, `FIREFOX`, `MSIE`, `MSEDGE`,   `MSCHREDGE`, `OPERA`, `OTHER`, `SAFARI`

* `cbi_profile` - (Optional) The cloud browser isolation profile to which the ISOLATE action is applied in the URL Filtering Policy rules. This block is required when the attribute `action` is set to `ISOLATE`
  * `id` - (Optional) The universally unique identifier (UUID) for the browser isolation profile
  * `name` - (Optional) Name of the browser isolation profile
  * `url` - (Optional) The browser isolation profile URL

* `cipa_rule` - (Optional) If set to true, the CIPA Compliance rule is enabled

* `url_categories` - (List of Strings) The list of URL categories to which the URL Filtering rule must be applied. See the [URL Categories API](https://help.zscaler.com/zia/url-categories#/urlCategories-get) for the list of available categories or use the data source `zia_url_categories` to retrieve the list of URL categories.

* `source_countries`** - (List of String) Identify destinations based on the location of a server. Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

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

* `location_groups` - (List of Object) The location groups to which the URL Filtering policy rule applies
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `labels` - (List of Object) The rule labels to which the URL Filtering policy rule applies
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `source_ip_groups` - (List of Object) The source ip groups to which the URL Filtering policy rule applies
  * `id` - (Optional) Source IP address groups for which the rule is applicable.

* `workload_groups` (Optional) The list of preconfigured workload groups to which the policy must be applied
  * `id` - (Optional) A unique identifier assigned to the workload group
  * `name` - (Optional) The name of the workload group

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_url_filtering_rules** can be imported by using `<RULE_ID>` or `<RULE_NAME>` as the import ID.

For example:

```shell
terraform import zia_url_filtering_rules.example <rule_id>
```

or

```shell
terraform import zia_url_filtering_rules.example <rule_name>
```
