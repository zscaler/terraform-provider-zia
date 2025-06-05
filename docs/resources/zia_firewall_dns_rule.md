---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_dns_rule"
description: |-
  Creates and manages ZIA Cloud firewall DNS rule.
---

# Resource: zia_firewall_dns_rule

The **zia_firewall_dns_rule** resource allows the creation and management of ZIA Cloud Firewall DNS rules in the Zscaler Internet Access.

**NOTE 1** Zscaler Cloud Firewall DNS Rules contain default and predefined rules which are placed in their respective orders. These rules `CANNOT` be deleted. When configuring your rules make sure that the `order` attributue value consider these pre-existing rules so that Terraform can place the new rules in the correct position, and drifts can be avoided. i.e If there are 2 pre-existing rules, you should start your rule order at `3` and manage your rule sets from that number onwards. The provider will reorder the rules automatically while ignoring the order of pre-existing rules, as the API will be responsible for moving these rules to their respective positions as API calls are made.

The most common default and predefined rules:

|              Rule Names                      |  Default or Predefined   |   Rule Number Associated |
|:--------------------------------------------:|:------------------------:|:------------------------:|
|-----------------------------|--------------------------|-------------------|
|  `Office 365 One Click Rule`                 |      `Predefined`       |           `Yes`          |
|  `ZPA Resolver for Road Warrior`             |      `Predefined`       |           `Yes`          |
|  `Critical risk DNS categories`              |      `Predefined`       |           `Yes`          |
|  `Critical risk DNS tunnels`                 |      `Predefined`       |           `Yes`          |
|  `High risk DNS categories`                  |      `Predefined`       |           `Yes`          |
|  `High risk DNS tunnels`                     |      `Predefined`       |           `Yes`          |
|  `Risky DNS categories`                      |      `Predefined`       |           `Yes`          |
|  `Risky DNS Risky DNS tunnels`               |      `Predefined`       |           `Yes`          |
|  `Unknown DNS Traffic`                       |      `Predefined`       |           `Yes`          |
|  `Default Firewall DNS Rule`                 |      `Predefined`       |           `Yes`          |
|  `ZPA Resolver for Locations`                |        `Default`        |           `No`           |
|  `Fallback ZPA Resolver for Locations`       |        `Default`        |           `No`           |
|  `Fallback ZPA Resolver for Road Warrior`    |        `Default`        |           `No`           |
|-------------------------|-------------------------|-----------------|

**NOTE 2** Certain attributes on `predefined` rules can still be managed or updated via Terraform such as:

- `description` - (Optional) Enter additional notes or information. The description cannot exceed 10,240 characters.
- `state` - (Optional) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to
- `labels` (list) - Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

**NOTE 3** The import of `predefined` rules is still possible in case you want o have them under the Terraform management; however, remember that these rules cannot be deleted. That means, the provider will fail when executing `terraform destroy`; hence, you must remove the rules you want to delete, and re-run `terraform apply` instead.

## Example Usage - Create Firewall DNS Rules - Redirect Action

```hcl
data "zia_department_management" "engineering" {
 name = "Engineering"
}

data "zia_group_management" "normal_internet" {
    name = "Normal_Internet"
}

data "zia_firewall_filtering_time_window" "work_hours" {
    name = "Work hours"
}

resource "zia_firewall_dns_rule" "this" {
    name = "Example_DNS_Rule01"
    description = "Example_DNS_Rule01"
    action = "REDIR_REQ"
    state = "ENABLED"
    order = 10
    rank = 7
    redirect_ip = "8.8.8.8"
    dest_countries = ["CA", "US"]
    source_countries = ["CA", "US"]
    protocols = ["ANY_RULE"]
    departments {
        id = [ data.zia_department_management.engineering.id ]
    }
    groups {
        id = [ data.zia_group_management.normal_internet.id ]
    }
    time_windows {
        id = [ data.zia_firewall_filtering_time_window.work_hours.id ]
    }
}
```

## Example Usage - Create Firewall DNS Rules - Redirect Request DOH

```hcl
resource "zia_firewall_dns_rule" "this2" {
    name = "Example_DNS_Rule02"
    description = "Example_DNS_Rule02"
    action = "REDIR_REQ_DOH"
    state = "ENABLED"
    order = 12
    rank = 7
    dest_countries = ["CA", "US"]
    source_countries = ["CA", "US"]
    protocols = ["ANY_RULE"]
    dns_gateway {
      id = 18207342
      name = "DNS_GW01"
    }
}
```

## Example Usage - Create Firewall DNS Rules - Redirect TCP Request

resource "zia_firewall_dns_rule" "this3" {
    name = "Example_DNS_Rule03"
    description = "Example_DNS_Rule03"
    action = "REDIR_REQ_TCP"
    state = "ENABLED"
    order = 13
    rank = 7
    dest_countries = ["CA", "US"]
    source_countries = ["CA", "US"]
    protocols = ["ANY_RULE"]
    dns_gateway {
      id = 18207342
      name = "DNS_GW01"
    }
}

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the Firewall Filtering policy rule
- `id` - (Optional) Unique identifier for the Firewall Filtering policy rule

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `description` - (String) Enter additional notes or information. The description cannot exceed 10,240 characters.
- `order` - (Integer) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.
- `state` - (String) An enabled rule is actively enforced. A disabled rule is not actively enforced but does not lose its place in the Rule Order. The service skips it and moves to the next rule.
- `action` - (String) The action configured for the rule that must take place if the traffic matches the rule criteria, such as allowing or blocking the traffic or bypassing the rule. The following actions are accepted: `ALLOW`, `BLOCK`, `REDIR_REQ`, `REDIR_RES`, `REDIR_ZPA`, `REDIR_REQ_DOH`, `REDIR_REQ_KEEP_SENDER`, `REDIR_REQ_TCP`, `REDIR_REQ_UDP`, `BLOCK_WITH_RESPONSE`
- `rank` - (Integer) By default, the admin ranking is disabled. To use this feature, you must enable admin rank. The default value is `7`.
- `access_control` - (String) The admin’s access privilege to this rule based on the assigned role
- `redirect_ip` - (String) The IP address to which the traffic will be redirected to when the DNAT rule is triggered. If not set, no redirection is done to specific IP addresses. Only supported when the `action` is `REDIR_REQ`

- `dns_rule_request_types` - (Set of Strings) DNS request types to which the rule applies. Supportedn values are:
`A`, `NS`, `MD`, `MF`, `CNAME`, `SOA`, `MB`, `MG`, `MR`, `NULL`, `WKS`, `PTR`, `HINFO`, `MINFO`, `MX`, `TXT`, `RP`, `AFSDB`,
`X25`, `ISDN`, `RT`, `NSAP`, `NSAP_PTR`, `SIG`, `KEY`, `PX`, `GPOS`, `AAAA`, `LOC`, `NXT`, `EID`, `NIMLOC`, `SRV`, `ATMA`,
`NAPTR`, `KX`, `CERT`, `A6`, `DNAME`, `SINK`, `OPT`, `APL`, `DS`, `SSHFP`, `PSECKEF`, `RRSIG`, `NSEC`, `DNSKEY`,
`DHCID`, `NSEC3`, `NSEC3PARAM`, `TLSA`, `HIP`, `NINFO`, `RKEY`, `TALINK`, `CDS`, `CDNSKEY`, `OPENPGPKEY`, `CSYNC`,
`ZONEMD`, `SVCB`, `HTTPS`,

- `protocols` - (Set of Strings) The protocols to which the rules applies. Supported Values: `ANY_RULE`, `SMRULEF_CASCADING_ALLOWED`, `TCP_RULE`, `UDP_RULE`, `DOHTTPS_RULE`

- `applications` - (Set of Strings) DNS tunnels and network applications to which the rule applies. To retrieve the available list of DNS tunnels applications use the data source: `zia_cloud_applications` with the `app_class` value `DNS_OVER_HTTPS`. See example:

```hcl
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_policy"
  app_class = ["DNS_OVER_HTTPS"]
}
```

- `block_response_code` - (String) Specifies the DNS response code to be sent to the client when the action is configured to block and send response code. Supported values are: `ANY`, `NONE`, `FORMERR`, `SERVFAIL`, `NXDOMAIN`, `NOTIMP`, `REFUSED`, `YXDOMAIN`, `YXRRSET`, `NXRRSET`, `NOTAUTH`, `NOTZONE`, `BADVERS`, `BADKEY`, `BADTIME`, `BADMODE`, `BADNAME`, `BADALG`, `BADTRUNC`, `UNSUPPORTED`, `BYPASS`, `INT_ERROR`, `SRV_TIMEOUT`, `EMPTY_RESP`,
`REQ_BLOCKED`, `ADMIN_DROP`, `WCDN_TIMEOUT`, `IPS_BLOCK`, `FQDN_RESOLV_FAIL`

- `capture_pcap` - (Boolean) Value that indicates whether packet capture (PCAP) is enabled or not
- `predefined` - (Boolean) A Boolean field that indicates that the rule is predefined by using a true value
- `default_rule` - (Boolean) Value that indicates whether the rule is the Default Cloud DNS Rule or not

`Devices`

- `devices` - (List of Objects) Devices to which the rule applies. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `device_groups` - (List of Objects) Device groups to which the rule applies. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
      - `id` - (Integer) Identifier that uniquely identifies an entity

`Who, Where and When` supports the following attributes:

- `locations` - (List of Objects) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `location_groups` - (List of Objects)You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `users` - (List of Objects) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `groups` - (List of Objects) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `departments` - (List of Objects) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `time_windows` - (List of Objects) You can manually select up to `1` time intervals. When not used it implies `always` to apply the rule to all time intervals.
      - `id` - (Integer) Identifier that uniquely identifies an entity

`source ip addresses` supports the following attributes:

- `source_countries` (Set of String) The countries of origin of traffic for which the rule is applicable. If not set, the rule is not restricted to specific source countries.
    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

- `src_ip_groups` - (List of Objects)Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `src_ipv6_groups` - (List of Objects) Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `src_ips` - (Set of String) Source IP addresses or FQDNs to which the rule applies. If not set, the rule is not restricted to a specific source IP address. Each IP entry can be a single IP address, CIDR (e.g., 10.10.33.0/24), or an IP range (e.g., 10.10.33.1-10.10.33.10).

`destinations` supports the following attributes:

- `dest_addresses` (Set of String) Destination IP addresses or FQDNs to which the rule applies. If not set, the rule is not restricted to a specific destination IP address. Each IP entry can be a single IP address, CIDR (e.g., 10.10.33.0/24), or an IP range (e.g., 10.10.33.1-10.10.33.10).

- `dest_countries` (Set of String) Identify destinations based on the location of a server, select Any to apply the rule to all countries or select the countries to which you want to control traffic.
    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

- `res_categories` (Set of String) URL categories associated with resolved IP addresses to which the rule applies. If not set, the rule is not restricted to a specific URL category.

- `dest_ip_categories` (Set of String)  identify destinations based on the URL category of the domain, select Any to apply the rule to all categories or select the specific categories you want to control.
- `dest_ip_groups`** - (List of Objects) Any number of destination IP address groups that you want to control with this rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `labels` (List of Objects) Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `application_groups` (List of Objects) DNS application groups to which the rule applies
      - `id` - (Integer) Identifier that uniquely identifies an entity

- `edns_ecs_object` (List of Objects) The EDNS ECS object which resolves DNS request. Only one object is supported.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (Integer) The configured name of the entity

- `dns_gateway` (Set of Objects) The DNS gateway used to redirect traffic, specified when the rule action is to redirect DNS request to an external DNS service. Only one DNS Gateway is supported.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (Integer) The configured name of the entity

- `zpa_ip_group` (Set of Objects) The ZPA IP pool specified when the rule action is to resolve domain names of ZPA applications to an ephemeral IP address from a preconfigured IP pool. Only one object is supported.
      - `id` - (Integer) Identifier that uniquely identifies an entity
      - `name` - (Integer) The configured name of the entity
