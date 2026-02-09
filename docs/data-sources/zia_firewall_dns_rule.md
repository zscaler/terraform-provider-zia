---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA): firewall_dns_rule"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-dns-control-policy
  API documentation https://help.zscaler.com/zia/dns-control-policy#/firewallDnsRules-post
  Get information about firewall DNS Control policy rule.
---

# zia_firewall_dns_rule (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-dns-control-policy)
* [API documentation](https://help.zscaler.com/zia/dns-control-policy#/firewallDnsRules-post)

Use the **zia_firewall_dns_rule** data source to get information about a cloud firewall DNS rule available in the Zscaler Internet Access.

## Example Usage

```hcl
# ZIA Firewall DNS Rule by name
data "zia_firewall_dns_rule" "this" {
    name = "Default Cloud IPS Rule"
}
```

```hcl
# ZIA Firewall DNS Rule by ID
data "zia_firewall_dns_rule" "this" {
    id = "12365478"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Firewall Filtering policy rule
* `id` - (Optional) Unique identifier for the Firewall Filtering policy rule

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Enter additional notes or information. The description cannot exceed 10,240 characters.
* `order` - (Integer) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.
* `state` - (String) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule.
* `action` - (String) The action configured for the rule that must take place if the traffic matches the rule criteria, such as allowing or blocking the traffic or bypassing the rule. The following actions are accepted: `ALLOW`, `BLOCK`, `REDIR_REQ`, `REDIR_RES`, `REDIR_ZPA`, `REDIR_REQ_DOH`, `REDIR_REQ_KEEP_SENDER`, `REDIR_REQ_TCP`, `REDIR_REQ_UDP`, `BLOCK_WITH_RESPONSE`
* `rank` - (Integer) By default, the admin ranking is disabled. To use this feature, you must enable admin rank. The default value is `7`.
* `access_control` - (String) The admin’s access privilege to this rule based on the assigned role
* `is_web_eun_enabled` - (Boolean) If set to true, Web EUN is enabled for the rule
* `default_dns_rule_name_used` - (Boolean) If set to true, the default DNS rule name is used for the rule
* `redirect_ip` - (String) The IP address to which the traffic will be redirected to when the DNAT rule is triggered. If not set, no redirection is done to specific IP addresses. Only supported when the `action` is `REDIR_REQ`
* `dns_rule_request_types` - (Set of Strings) DNS request types to which the rule applies. Supportedn values are:
`A`, `NS`, `MD`, `MF`, `CNAME`, `SOA`, `MB`, `MG`, `MR`, `NULL`, `WKS`, `PTR`, `HINFO`, `MINFO`, `MX`, `TXT`, `RP`, `AFSDB`,
`X25`, `ISDN`, `RT`, `NSAP`, `NSAP_PTR`, `SIG`, `KEY`, `PX`, `GPOS`, `AAAA`, `LOC`, `NXT`, `EID`, `NIMLOC`, `SRV`, `ATMA`,
`NAPTR`, `KX`, `CERT`, `A6`, `DNAME`, `SINK`, `OPT`, `APL`, `DS`, `SSHFP`, `PSECKEF`, `RRSIG`, `NSEC`, `DNSKEY`,
`DHCID`, `NSEC3`, `NSEC3PARAM`, `TLSA`, `HIP`, `NINFO`, `RKEY`, `TALINK`, `CDS`, `CDNSKEY`, `OPENPGPKEY`, `CSYNC`,
`ZONEMD`, `SVCB`, `HTTPS`,

* `protocols` - (Set of Strings) The protocols to which the rules applies. Supported Values: `ANY_RULE`, `SMRULEF_CASCADING_ALLOWED`, `TCP_RULE`, `UDP_RULE`, `DOHTTPS_RULE`

* `applications` - (Set of Strings) DNS tunnels and network applications to which the rule applies. To retrieve the available list of DNS tunnels applications use the data source: `zia_cloud_applications` with the `app_class` value `DNS_OVER_HTTPS`. See example:

```hcl
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_policy"
  app_class = ["DNS_OVER_HTTPS"]
}
```

* `block_response_code` - (String) Specifies the DNS response code to be sent to the client when the action is configured to block and send response code. Supported values are: `ANY`, `NONE`, `FORMERR`, `SERVFAIL`, `NXDOMAIN`, `NOTIMP`, `REFUSED`, `YXDOMAIN`, `YXRRSET`, `NXRRSET`, `NOTAUTH`, `NOTZONE`, `BADVERS`, `BADKEY`, `BADTIME`, `BADMODE`, `BADNAME`, `BADALG`, `BADTRUNC`, `UNSUPPORTED`, `BYPASS`, `INT_ERROR`, `SRV_TIMEOUT`, `EMPTY_RESP`,
`REQ_BLOCKED`, `ADMIN_DROP`, `WCDN_TIMEOUT`, `IPS_BLOCK`, `FQDN_RESOLV_FAIL`

* `capture_pcap` - (Boolean) Value that indicates whether packet capture (PCAP) is enabled or not
* `predefined` - (Boolean) A Boolean field that indicates that the rule is predefined by using a true value
* `default_rule` - (Boolean) Value that indicates whether the rule is the Default Cloud DNS Rule or not
* `is_web_eun_enabled` - (Boolean) A Boolean value that indicates whether Enhanced User Notification (EUN) is enabled for the rule.
* `default_dns_rule_name_used` - (Boolean) A Boolean value that indicates whether the default DNS rule name is used for the rule.

`Devices`

* `devices` - (List of Objects) Devices to which the rule applies. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `device_groups` - (List of Objects) Device groups to which the rule applies. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

`Who, Where and When` supports the following attributes:

* `locations` - (List of Objects) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `location_groups` - (List of Objects)You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `users` - (List of Objects) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `groups` - (List of Objects) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
* `departments` - (List of Objects) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `time_windows` - (List of Objects) You can manually select up to `1` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

`source ip addresses` supports the following attributes:

* `source_countries` (Set of String) The countries of origin of traffic for which the rule is applicable. If not set, the rule is not restricted to specific source countries.
    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `src_ip_groups` - (List of Objects)Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `src_ipv6_groups` - (List of Objects) Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `src_ips` - (Set of String) Source IP addresses or FQDNs to which the rule applies. If not set, the rule is not restricted to a specific source IP address. Each IP entry can be a single IP address, CIDR (e.g., 10.10.33.0/24), or an IP range (e.g., 10.10.33.1-10.10.33.10).

`destinations` supports the following attributes:

* `dest_addresses` (Set of String) Destination IP addresses or FQDNs to which the rule applies. If not set, the rule is not restricted to a specific destination IP address. Each IP entry can be a single IP address, CIDR (e.g., 10.10.33.0/24), or an IP range (e.g., 10.10.33.1-10.10.33.10).

* `dest_countries` (Set of String) Identify destinations based on the location of a server, select Any to apply the rule to all countries or select the countries to which you want to control traffic.
    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `res_categories` (Set of String) URL categories associated with resolved IP addresses to which the rule applies. If not set, the rule is not restricted to a specific URL category.

* `dest_ip_categories` (Set of String)  identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
* `dest_ip_groups`** - (List of Objects) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `labels` (List of Objects) Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity

* `application_groups` (List of Objects) DNS application groups to which the rule applies
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (string) The configured name of the entity

* `edns_ecs_object` (List of Objects) The EDNS ECS object which resolves DNS request
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (string) The configured name of the entity

* `dns_gateway` (Set) The DNS gateway used to redirect traffic, specified when the rule action is to redirect DNS request to an external DNS service.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (string) The configured name of the entity

* `zpa_ip_group` (Set) The ZPA IP pool specified when the rule action is to resolve domain names of ZPA applications to an ephemeral IP address from a preconfigured IP pool
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (string) The configured name of the entity
