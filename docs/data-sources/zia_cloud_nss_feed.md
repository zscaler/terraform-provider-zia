---
subcategory: "Cloud Nanolog Streaming Service (NSS)"
layout: "zscaler"
page_title: "ZIA: cloud_nss_feed"
description: |-
  Official documentation https://help.zscaler.com/zia/about-nss-feeds
  API documentation https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get
  Retrieves the cloud NSS feeds configured in the ZIA Admin Portal
---

# zia_cloud_nss_feed (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-nss-feeds)
* [API documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

Use the **zia_cloud_nss_feed** data source to get information about cloud NSS feeds configured in the ZIA Admin Portal

## Example Usage

```hcl
data "zia_cloud_nss_feed" "this" {
  name = "Google FW"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The unique identifier for the nss server
* `name` - (Optional) The name of the cloud NSS feed

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `feed_status` - (string) The status of the feed
* `nss_log_type` - (string) The type of NSS logs that are streamed (e.g. Web, Firewall, DNS, Alert, etc.)
* `nss_feed_type` - (string) NSS feed format type (e.g. CSV, syslog, Splunk Common Information Model (CIM), etc.)
* `feed_output_format` - (string) Output format used for the feed
* `user_obfuscation` - (string) Specifies whether user obfuscation is enabled or disabled
* `time_zone` - (string) Specifies the time zone that must be used in the output file
* `custom_escaped_character` - (Set of String) Characters that need to be encoded using hex when they appear in URL, Host, or Referrer
* `eps_rate_limit` - (int) Event per second limit
* `json_array_toggle` - (bool) A Boolean value indicating whether streaming of logs in JSON array format (e.g., [{JSON1},{JSON2}]) is enabled or disabled for the JSON feed output type
* `siem_type` - (string) Cloud NSS SIEM type
* `max_batch_size` - (int) The maximum batch size in KB
* `connection_url` - (string) The HTTPS URL of the SIEM log collection API endpoint
* `authentication_token` - (string) The authentication token value
* `connection_headers` - (Set of String) The HTTP Connection headers
* `last_success_full_test` - (int) The timestamp of the last successful test. Value is in Unix time.
* `test_connectivity_code` - (int) The code from the last test
* `base64_encoded_certificate` - (string) Base64-encoded certificate
* `nss_type` - (string) NSS type
* `client_id` - (string) Client ID applicable when SIEM type is set to S3 or Azure Sentinel
* `client_secret` - (string) Client secret applicable when SIEM type is set to S3 or Azure Sentinel
* `authentication_url` - (string) Authentication URL applicable when SIEM type is set to Azure Sentinel
* `grant_type` - (string) Grant type applicable when SIEM type is set to Azure Sentinel
* `scope` - (string) Scope applicable when SIEM type is set to Azure Sentinel
* `oauth_authentication` - (bool) A Boolean value indicating whether OAuth 2.0 authentication is enabled or not
* `server_ips` - (Set of String) Filter to limit the logs based on the server's IPv4 addresses
* `client_ips` - (Set of String) Filter to limit the logs based on a client's public IPv4 addresses
* `domains` - (Set of String) Filter to limit the logs to sessions associated with specific domains
* `dns_request_types` - (Set of String) DNS request types included in the feed
* `dns_response_types` - (Set of String) DNS response types filter
* `dns_responses` - (Set of String) DNS responses filter
* `durations` - (Set of String) Filter based on time durations
* `dns_actions` - (Set of String) DNS Control policy action filter
* `firewall_logging_mode` - (string) Filter based on the Firewall Filtering policy logging mode
* `client_source_ips` - (Set of String) Client source IPs configured for NSS feed.
* `firewall_actions` - (Set of String) Firewall actions included in the NSS feed.
* `countries` - (Set of String) Countries filter in the Firewall policy
* `server_source_ports` - (Set of String) Firewall log filter based on the traffic destination name
* `client_source_ports` - (Set of String) Firewall log filter based on a client's source ports
* `action_filter` - (string) Policy action filter
* `email_dlp_policy_action` - (string) Action filter for Email DLP log type
* `direction` - (string) Traffic direction filter specifying inbound or outbound
* `event` - (string) CASB event filter
* `policy_reasons` - (Set of String) Policy reason filter
* `protocol_types` - (Set of String) Protocol types filter
* `user_agents` - (Set of String) Predefined user agents filter
* `request_methods` - (Set of String) Request methods filter
* `casb_severity` - (Set of String) Zscaler's Cloud Access Security Broker (CASB) severity filter
* `casb_policy_types` - (Set of String) CASB policy type filter
* `casb_applications` - (Set of String) CASB application filter
* `casb_action` - (Set of String) CASB policy action filter
* `url_super_categories` - (Set of String) URL supercategory filter
* `web_applications` - (Set of String) Filter to include specific cloud applications in the logs. By default, all cloud applications are included in the logs. To obtain the list of cloud applications that can be specified in this attribute, use the GET /cloudApplications/lite request.
* `web_applications_exclude` - (Set of String) Filter to exclude specific cloud applications from the logs. By default, no cloud applications is excluded from the logs. To obtain the list of cloud applications that can be specified in this attribute, use the GET /cloudApplications/lite request.
* `web_application_classes` - (Set of String) Cloud application categories Filter
* `malware_names` - (Set of String) Filter based on malware names
* `malware_classes` - (Set of String) Malware category filter
* `url_classes` - (Set of String) URL category filter
* `advanced_threats` - (Set of String) Advanced threats filter
* `response_codes` - (Set of String) Advanced threats filter
* `nw_applications` - (Set of String) Filter to include specific network applications in the logs. By default, all network applications are included in the logs
* `nw_applications_exclude` - (Set of String) Filter to include specific network applications in the logs. By default, no network application is excluded from the logs
* `nat_actions` - (Set of String) NAT Control policy actions filter
* `traffic_forwards` - (Set of String) Filter based on the firewall traffic forwarding method
* `web_traffic_forwards` - (Set of String) Filter based on the web traffic forwarding method
* `tunnel_types` - (Set of String) Tunnel type filter
* `alerts` - (Set of String) Alert filter
* `object_type` - (Set of String) CRM object type filter
* `activity` - (Set of String) CASB activity filter
* `object_type1` - (Set of String) CASB activity object type filter
* `object_type2` - (Set of String) CASB activity object type filter if applicable
* `end_point_dlp_log_type` - (Set of String) Endpoint DLP log type filter
* `email_dlp_log_type` - (Set of String) Email DLP record type filter
* `file_type_super_categories` - (Set of String) Filter based on the category of file type in download
* `file_type_categories` - (Set of String) Filter based on the file type in download
* `casb_file_type_super_categories` - (Set of String) Endpoint DLP file type category filer
* `file_sizes` - (Set of String) File size filter
* `request_sizes` - (Set of String) Request size filter
* `response_sizes` - (Set of String) Request size filter
* `transaction_sizes` - (Set of String) Transaction size filter
* `in_bound_bytes` - (Set of String) Filter based on inbound bytes
* `out_bound_bytes` - (Set of String) Filter based on outbound bytes
* `download_time` - (Set of String) Download time filter
* `scan_time` - (Set of String) Scan time filter
* `server_source_ips` - (Set of String) Filter based on the server's source IPv4 addresses in Firewall policy
* `server_destination_ips` - (Set of String) Filter based on the server's destination IPv4 addresses in Firewall policy
* `tunnel_ips` - (Set of String) Filter based on tunnel IPv4 addresses in Firewall policy
* `internal_ips` - (Set of String) Filter based on internal IPv4 addresses
* `tunnel_source_ips` - (Set of String) Source IPv4 addresses of tunnels
* `tunnel_dest_ips` - (Set of String) Destination IPv4 addresses of tunnels
* `client_destination_ips` - (Set of String) Client's destination IPv4 addresses in Firewall policy
* `audit_log_type` - (Set of String) Audit log type filter
* `project_name` - (Set of String) Repository project name filter
* `repo_name` - (Set of String) Repository name filter
* `object_name` - (Set of String) CRM object name filter
* `channel_name` - (Set of String) Collaboration channel name filter
* `file_source` - (Set of String) Filter based on the file source
* `file_name` - (Set of String) Filter based on the file name
* `session_counts` - (Set of String) Firewall logs filter based on the number of sessions
* `adv_user_agents` - (Set of String) Filter based on custom user agent strings
* `referer_urls` - (Set of String) Referrer URL filter
* `host_names` - (Set of String) Filter to limit the logs based on specific hostnames
* `full_urls` - (Set of String) Filter to limit the logs based on specific full URLs
* `threat_names` - (Set of String) Filter based on threat names
* `page_risk_indexes` - (Set of String) Page Risk Index filter
* `client_destination_ports` - (Set of String) Firewall logs filter based on a client's destination
* `tunnel_source_port` - (Set of String) Filter based on the tunnel source port

### Block Attributes

Each of the following blocks supports nested attributes:

#### `external_owners`

* `id` - (int) Identifier that uniquely identifies an entity
* `name` - (string) Identifier that uniquely identifies an entity
* `extensions` - (Map of String) Optional metadata for the entity

#### `external_collaborators`

* `id` - (int) Identifier that uniquely identifies an entity
* `name` - (string) Identifier that uniquely identifies an entity
* `extensions` - (Map of String) Optional metadata for the entity

#### `internal_collaborators`

* `id` - (int) Identifier that uniquely identifies an entity
* `name` - (string) Identifier that uniquely identifies an entity
* `extensions` - (Map of String) Optional metadata for the entity

#### `itsm_object_type`

* `id` - (int) Identifier that uniquely identifies an entity
* `name` - (string) Identifier that uniquely identifies an entity
* `extensions` - (Map of String) Optional metadata for the entity

#### `url_categories`

* `id` - (int) Identifier that uniquely identifies an entity
* `name` - (string) Identifier that uniquely identifies an entity
* `extensions` - (Map of String) Optional metadata for the entity

#### `dlp_engines`

* `id` - (int) Identifier that uniquely identifies an entity
* `name` - (string) Identifier that uniquely identifies an entity
* `extensions` - (Map of String) Optional metadata for the entity

#### `dlp_dictionaries`

* `id` - (int) Identifier that uniquely identifies an entity
* `name` - (string) Identifier that uniquely identifies an entity
* `extensions` - (Map of String) Optional metadata for the entity

#### `rules`

* `id` - (int) Identifier that uniquely identifies an entity
* `name` - (string) Identifier that uniquely identifies an entity
* `extensions` - (Map of String) Optional metadata for the entity

#### `nw_services`

* `id` - (int) Identifier that uniquely identifies an entity
* `name` - (string) Identifier that uniquely identifies an entity
* `extensions` - (Map of String) Optional metadata for the entity

#### `locations`

* `id` - (int) A unique identifier for the location
* `pid` - (int) Parent identifier for the location
* `name` - (string) The configured name of the location
* `description` - (string) Description of the location
* `deleted` - (bool) Indicates if the location is deleted
* `getl_id` - (int) GETL identifier for the location

#### `location_groups`

* `id` - (int) A unique identifier for the location group
* `pid` - (int) Parent identifier for the location group
* `name` - (string) The configured name of the location group
* `description` - (string) Description of the location group
* `deleted` - (bool) Indicates if the location group is deleted
* `getl_id` - (int) GETL identifier for the location group

#### `casb_tenant`

* `id` - (int) A unique identifier for the CASB tenant
* `pid` - (int) Parent identifier for the CASB tenant
* `name` - (string) The configured name of the CASB tenant
* `description` - (string) Description of the CASB tenant
* `deleted` - (bool) Indicates if the CASB tenant is deleted
* `getl_id` - (int) GETL identifier for the CASB tenant

#### `users`

* `id` - (int) A unique identifier for the user
* `pid` - (int) Parent identifier for the user
* `name` - (string) The configured name of the user
* `description` - (string) Description of the user
* `deleted` - (bool) Indicates if the user is deleted
* `getl_id` - (int) GETL identifier for the user

#### `departments`

* `id` - (int) A unique identifier for the department
* `pid` - (int) Parent identifier for the department
* `name` - (string) The configured name of the department
* `description` - (string) Description of the department
* `deleted` - (bool) Indicates if the department is deleted
* `getl_id` - (int) GETL identifier for the department

#### `sender_name`

* `id` - (int) A unique identifier for the sender
* `pid` - (int) Parent identifier for the sender
* `name` - (string) The configured name of the sender
* `description` - (string) Description of the sender
* `deleted` - (bool) Indicates if the sender is deleted
* `getl_id` - (int) GETL identifier for the sender

#### `buckets`

* `id` - (int) A unique identifier for the bucket
* `pid` - (int) Parent identifier for the bucket
* `name` - (string) The configured name of the bucket
* `description` - (string) Description of the bucket
* `deleted` - (bool) Indicates if the bucket is deleted
* `getl_id` - (int) GETL identifier for the bucket

#### `vpn_credentials`

* `id` - (int) A unique identifier for the VPN credential
* `pid` - (int) Parent identifier for the VPN credential
* `name` - (string) The configured name of the VPN credential
* `description` - (string) Description of the VPN credential
* `deleted` - (bool) Indicates if the VPN credential is deleted
* `getl_id` - (int) GETL identifier for the VPN credential
