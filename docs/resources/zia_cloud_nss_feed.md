---
subcategory: "Cloud Nanolog Streaming Service (NSS)"
layout: "zscaler"
page_title: "ZIA: cloud_nss_feed"
description: |-
  Official documentation https://help.zscaler.com/zia/about-nss-feeds
  API documentation https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get
  Manages cloud NSS feeds in the ZIA Admin Portal
---

# zia_cloud_nss_feed (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-nss-feeds)
* [API documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

Use the **zia_cloud_nss_feed** resource to create, update, and delete cloud NSS feeds in the ZIA Admin Portal

## Example Usage - NSS Splunk Feed

```hcl
resource "zia_cloud_nss_feed" "this" {
  name          = "Splunk_Feed"
  feed_status   = "ENABLED"
  nss_log_type  = "WEBLOG"
  nss_feed_type = "JSON"
  time_zone     = "GMT"
  custom_escaped_character = [
    "ASCII_34",
    "ASCII_44",
    "ASCII_92"
  ]
  eps_rate_limit     = 0
  max_batch_size     = 512
  json_array_toggle  = true
  siem_type          = "SPLUNK"
  connection_url     = "http://3.87.81.187:8088/services/collector?auto_extract_timestamp=true"
  connection_headers = ["Authorization:Splunk xxxxx-xxx-xxxx-xxxx-xxxxxxxx"]
  nss_type           = "NSS_FOR_WEB"

  feed_output_format = "\\{ \"sourcetype\" : \"zscalernss-web\", \"event\" : \\{\"datetime\":\"%d{yy}-%02d{mth}-%02d{dd} %02d{hh}:%02d{mm}:%02d{ss}\",\"reason\":\"%s{reason}\",\"event_id\":\"%d{recordid}\",\"protocol\":\"%s{proto}\",\"action\":\"%s{action}\",\"transactionsize\":\"%d{totalsize}\",\"responsesize\":\"%d{respsize}\",\"requestsize\":\"%d{reqsize}\",\"urlcategory\":\"%s{urlcat}\",\"serverip\":\"%s{sip}\",\"requestmethod\":\"%s{reqmethod}\",\"refererURL\":\"%s{ereferer}\",\"useragent\":\"%s{eua}\",\"product\":\"NSS\",\"location\":\"%s{elocation}\",\"ClientIP\":\"%s{cip}\",\"status\":\"%s{respcode}\",\"user\":\"%s{elogin}\",\"url\":\"%s{eurl}\",\"vendor\":\"Zscaler\",\"hostname\":\"%s{ehost}\",\"clientpublicIP\":\"%s{cintip}\",\"threatcategory\":\"%s{malwarecat}\",\"threatname\":\"%s{threatname}\",\"filetype\":\"%s{filetype}\",\"appname\":\"%s{appname}\",\"app_status\":\"%s{app_status}\",\"pagerisk\":\"%d{riskscore}\",\"threatseverity\":\"%s{threatseverity}\",\"department\":\"%s{edepartment}\",\"urlsupercategory\":\"%s{urlsupercat}\",\"appclass\":\"%s{appclass}\",\"dlpengine\":\"%s{dlpeng}\",\"urlclass\":\"%s{urlclass}\",\"threatclass\":\"%s{malwareclass}\",\"dlpdictionaries\":\"%s{dlpdict}\",\"fileclass\":\"%s{fileclass}\",\"bwthrottle\":\"%s{bwthrottle}\",\"contenttype\":\"%s{contenttype}\",\"unscannabletype\":\"%s{unscannabletype}\",\"deviceowner\":\"%s{deviceowner}\",\"devicehostname\":\"%s{devicehostname}\",\"keyprotectiontype\":\"%s{keyprotectiontype}\"\\}\\}\n"
    departments {
        id = [ 4451590 ]
    }
    users {
        id = [ 6438644 ]
    }
    url_categories {
        id = [ data.zia_url_categories.this.val ]
    }
}
```

## Example Usage - NSS Splunk Feed - Google Chronicle

```hcl
resource "zia_cloud_nss_feed" "this" {
  name          = "Google_Firewall_Terraform"
  feed_status   = "ENABLED"
  nss_log_type  = "FWLOG"
  nss_feed_type = "JSON"
  time_zone     = "GMT_06_00_BANGLADESH_CENTRAL_ASIA_GMT_06_00"
  custom_escaped_character = [
    "ASCII_44",
    "ASCII_92",
    "ASCII_34"
  ]
  # eps_rate_limit     = 0
  max_batch_size     = 512
  json_array_toggle  = true
  siem_type          = "SPLUNK"
  connection_url     = "https://us-chronicle.googleapis.com/v1alpha/projects/xxxxxxx/locations/us/instances/xxxxxxxxxxx/feeds/xxxx-xxxx-xxxx-xxxx-xxxxxxxx:importPushLogs"
  connection_headers = [
    "X-goog-api-key: <Your API Key>",
    "X-Webhook-Access-Key:<Your Webhook API Key>"
    ]
  nss_type           = "NSS_FOR_FIREWALL"
  firewall_logging_mode = "ALL"
  feed_output_format = "\\{ \"sourcetype\" : \"zscalernss-fw\", \"event\" :\\{\"datetime\":\"%s{time}\",\"user\":\"%s{elogin}\",\"department\":\"%s{dept}\",\"locationname\":\"%s{location}\",\"cdport\":\"%d{cdport}\",\"csport\":\"%d{csport}\",\"sdport\":\"%d{sdport}\",\"ssport\":\"%d{ssport}\",\"csip\":\"%s{csip}\",\"cdip\":\"%s{cdip}\",\"ssip\":\"%s{ssip}\",\"sdip\":\"%s{sdip}\",\"tsip\":\"%s{tsip}\",\"tunsport\":\"%d{tsport}\",\"tuntype\":\"%s{ttype}\",\"action\":\"%s{action}\",\"dnat\":\"%s{dnat}\",\"stateful\":\"%s{stateful}\",\"aggregate\":\"%s{aggregate}\",\"nwsvc\":\"%s{nwsvc}\",\"nwapp\":\"%s{nwapp}\",\"proto\":\"%s{ipproto}\",\"ipcat\":\"%s{ipcat}\",\"destcountry\":\"%s{destcountry}\",\"avgduration\":\"%d{avgduration}\",\"rulelabel\":\"%s{erulelabel}\",\"inbytes\":\"%ld{inbytes}\",\"outbytes\":\"%ld{outbytes}\",\"duration\":\"%d{duration}\",\"durationms\":\"%d{durationms}\",\"numsessions\":\"%d{numsessions}\",\"ipsrulelabel\":\"%s{ipsrulelabel}\",\"threatcat\":\"%s{threatcat}\",\"threatname\":\"%s{ethreatname}\",\"deviceowner\":\"%s{deviceowner}\",\"devicehostname\":\"%s{devicehostname}\",\"threat_score\":\"%d{threat_score}\",\"threat_severity\":\"%s{threat_severity}\"\\}\\}\n"

    departments {
        id = [ 4451590 ]
    }
    users {
        id = [ 6438644 ]
    }
    url_categories {
        id = [ data.zia_url_categories.this.val ]
    }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The name of the cloud NSS feed

### Optional

* `feed_status` - (Optional) The status of the feed

* `nss_log_type` - (Optional) The type of NSS logs that are streamed (e.g. Web, Firewall, DNS, Alert, etc.). Supported values: `ADMIN_AUDIT`, `WEBLOG`, `ALERT`, `FWLOG`, `DNSLOG`, `MULTIFEEDLOG`, `CASB_FILELOG`, `CASB_MAILLOG`, `ECLOG`, `EC_DNSLOG`, `CASB_ITSM`, `CASB_CRM`, `CASB_CODE_REPO`, `CASB_COLLAB`, `CASB_PCS`, `USER_ACT_REP`, `USER_COUNT_ALERT`, `USER_IMP_TRAVEL_ALERT`, `ENDPOINT_DLP`, `EC_EVENTLOG`, `EMAIL_DLP`

* `nss_feed_type` - (Optional) NSS feed format type (e.g. CSV, syslog, Splunk Common Information Model (CIM), etc.). Supported values: `QRADAR`, `SYSLOG`, `CSV`, `TAB_SEPARATED`, `CUSTOM`, `SPLUNK_CIM`, `NAME_VALUE_PAIRS`, `RSA_SECURITY`, `ARCSIGHT_CEF`, `SYMANTEC_MSS`, `LOGRHYTHM`, `ZBRIDGE`, `MCAS`, `JSON`, `ZFAB_AGENT`

* `feed_output_format` - (Optional) Output format used for the feed

* `time_zone` - (Optional) Specifies the time zone that must be used in the output file. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get) for a list of supported time zones.

* `custom_escaped_character` - (Optional) Characters that need to be encoded using hex when they appear in URL, Host, or Referrer

* `eps_rate_limit` - (Optional) Event per second limit

* `json_array_toggle` - (Optional) A Boolean value indicating whether streaming of logs in JSON array format (e.g., [{JSON1},{JSON2}]) is enabled or disabled for the JSON feed output type

* `siem_type` - (Optional) Cloud NSS SIEM type. Supported values: `SPLUNK`, `SUMO_LOGIC`, `DEVO`, `OTHER`, `AZURE_SENTINEL`, `S3`

* `max_batch_size` - (Optional) The maximum batch size in KB

* `connection_url` - (Optional) The HTTPS URL of the SIEM log collection API endpoint

* `authentication_token` - (Optional) The authentication token value

* `connection_headers` - (Optional) The HTTP Connection headers

* `base64_encoded_certificate` - (Optional) Base64-encoded certificate

* `nss_type` - (Optional) NSS type. Supported values: `NONE`, `SOFTWARE_AA_FLAG`, `NSS_FOR_WEB`, `NSS_FOR_FIREWALL`, `VZEN`, `VZEN_SME`, `VZEN_SMLB`, `PINNED_NSS`, `MD5_CAPABLE`, `ADP`, `ZIRSVR`, `NSS_FOR_ZPA`

* `client_id` - (Optional) Client ID applicable when SIEM type is set to S3 or Azure Sentinel

* `client_secret` - (Optional) Client secret applicable when SIEM type is set to S3 or Azure Sentinel

* `authentication_url` - (Optional) Authentication URL applicable when SIEM type is set to Azure Sentinel

* `grant_type` - (Optional) Grant type applicable when SIEM type is set to Azure Sentinel

* `scope` - (Optional) Scope applicable when SIEM type is set to Azure Sentinel

* `oauth_authentication` - (Optional) A Boolean value indicating whether OAuth 2.0 authentication is enabled or not

* `server_ips` - (Optional) Filter to limit the logs based on the server's IPv4 addresses

* `client_ips` - (Optional) Filter to limit the logs based on a client's public IPv4 addresses

* `domains` - (Optional) Filter to limit the logs to sessions associated with specific domains

* `dns_request_types` - (Optional) DNS request types included in the feed. Supported values: `ANY`, `NONE`, `DNSREQ_A`, `DNSREQ_NS`, `DNSREQ_CNAME`, `DNSREQ_SOA`, `DNSREQ_WKS`, `DNSREQ_PTR`, `DNSREQ_HINFO`, `DNSREQ_MINFO`, `DNSREQ_MX`, `DNSREQ_TXT`, `DNSREQ_AAAA`, `DNSREQ_ISDN`, `DNSREQ_LOC`, `DNSREQ_RP`, `DNSREQ_RT`, `DNSREQ_MR`, `DNSREQ_MG`, `DNSREQ_MB`, `DNSREQ_AFSDB`, `DNSREQ_HIP`, `DNSREQ_SRV`, `DNSREQ_DS`, `DNSREQ_NAPTR`, `DNSREQ_NSEC`, `DNSREQ_DNSKEY`, `DNSREQ_HTTPS`, `DNSREQ_UNKNOWN`

* `dns_response_types` - (Optional) DNS response types filter. Supported values: `ANY`, `DNSRES_ZSCODE`, `DNSRES_CNAME`, `DNSRES_IPV6`, `DNSRES_SRV_CODE`, `DNSRES_IPV4`

* `dns_responses` - (Optional) DNS responses filter

* `durations` - (Optional) Filter based on time durations

* `dns_actions` - (Optional) DNS Control policy action filter

* `firewall_logging_mode` - (Optional) Filter based on the Firewall Filtering policy logging mode. Supported values: `SESSION`, `AGGREGATE`, `ALL`

* `client_source_ips` - (Optional) Client source IPs configured for NSS feed

* `firewall_actions` - (Optional) Firewall actions included in the NSS feed. Supported values: `BLOCK`, `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BLOCK_ICMP`, `COUNTRY_BLOCK`, `IPS_BLOCK_DROP`, `IPS_BLOCK_RESET`, `ALLOW_INSUFFICIENT_APPDATA`, `BLOCK_ABUSE_DROP`, `INT_ERR_DROP`, `CFG_BYPASSED`, `CFG_TIMEDOUT`

* `countries` - (Optional) Countries filter in the Firewall policy. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `server_source_ports` - (Optional) Firewall log filter based on the traffic destination name

* `client_source_ports` - (Optional) Firewall log filter based on a client's source ports

* `action_filter` - (Optional) Policy action filter. Supported values: `ALLOWED`, `BLOCKED`

* `email_dlp_policy_action` - (Optional) Action filter for Email DLP log type. Supported values: `ALLOW`, `CUSTOMHEADERINSERTION`, `BLOCK`

* `direction` - (Optional) Traffic direction filter specifying inbound or outbound. Supported values: `INBOUND`, `OUTBOUND`

* `event` - (Optional) CASB event filter. Supported values: `SCAN`, `VIOLATION`, `INCIDENT`

* `policy_reasons` - (Optional) Policy reason filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `protocol_types` - (Optional) Protocol types filter. Supported values: `TUNNEL`, `SSL`, `HTTP`, `HTTPS`, `FTP`, `FTPOVERHTTP`, `HTTP_PROXY`, `TUNNEL_SSL`, `DNSOVERHTTPS`, `WEBSOCKET`, `WEBSOCKET_SSL`

* `user_agents` - (Optional) Predefined user agents filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `request_methods` - (Optional) Request methods filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `casb_severity` - (Optional) Zscaler's Cloud Access Security Broker (CASB) severity filter. Supported values: `RULE_SEVERITY_HIGH`, `RULE_SEVERITY_MEDIUM`, `RULE_SEVERITY_LOW`, `RULE_SEVERITY_INFO`

* `casb_policy_types` - (Optional) CASB policy type filter. Supported values: `MALWARE`, `DLP`, `ALL_INCIDENT`

* `casb_applications` - (Optional) CASB application filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `casb_action` - (Optional) CASB policy action filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `url_super_categories` - (Optional) URL supercategory filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `web_applications` - (Optional) Filter to include specific cloud applications in the logs. By default, all cloud applications are included in the logs. To obtain the list of cloud applications that can be specified in this attribute, use the GET /cloudApplications/lite request. To retrieve the list of cloud applications, use the data source: `zia_cloud_applications`

* `web_applications_exclude` - (Optional) Filter to exclude specific cloud applications from the logs. By default, no cloud applications is excluded from the logs. To obtain the list of cloud applications that can be specified in this attribute, use the GET /cloudApplications/lite request. To retrieve the list of cloud applications, use the data source: `zia_cloud_applications`

* `web_application_classes` - (Optional) Cloud application categories Filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `malware_names` - (Optional) Filter based on malware names. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `malware_classes` - (Optional) Malware category filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `url_classes` - (Optional) URL category filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `advanced_threats` - (Optional) Advanced threats filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `response_codes` - (Optional) Response codes filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `nw_applications` - (Optional) Filter to include specific network applications in the logs. By default, all network applications are included in the logs. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `nw_applications_exclude` - (Optional) Filter to include specific network applications in the logs. By default, no network application is excluded from the logs. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `nat_actions` - (Optional) NAT Control policy actions filter. Supported values: `NONE`, `DNAT`

* `traffic_forwards` - (Optional) Filter based on the firewall traffic forwarding method. Supported values: `ANY`, `NONE`, `PBF`, `GRE`, `IPSEC`, `Z_APP`, `ZAPP_GRE`, `ZAPP_IPSEC`, `EC`, `MTGRE`, `ZAPP_DIRECT`, `CCA`, `MTUN_PROXY`, `MTUN_CBI`

* `web_traffic_forwards` - (Optional) Filter based on the web traffic forwarding method. Supported values: `ANY`, `NONE`, `PBF`, `GRE`, `IPSEC`, `Z_APP`, `ZAPP_GRE`, `ZAPP_IPSEC`, `EC`, `MTGRE`, `ZAPP_DIRECT`, `CCA`, `MTUN_PROXY`, `MTUN_CBI`

* `tunnel_types` - (Optional) Tunnel type filter. Supported values: `GRE`, `IPSEC_IKEV1`, `IPSEC_IKEV2`, `SVPN`, `EXTRANET`, `ZUB`, `ZCB`

* `alerts` - (Optional) Alert filter. Supported values: `CRITICAL`, `WARN`

* `object_type` - (Optional) CRM object type filter

* `activity` - (Optional) CASB activity filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `object_type1` - (Optional) CASB activity object type filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `object_type2` - (Optional) CASB activity object type filter if applicable. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `end_point_dlp_log_type` - (Optional) Endpoint DLP log type filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `email_dlp_log_type` - (Optional) Email DLP record type filter. Supported values: `EPDLP_SCAN_AGGREGATE`, `EPDLP_SENSITIVE_ACTIVITY`, `EPDLP_DLP_INCIDENT`

* `file_type_super_categories` - (Optional) Filter based on the category of file type in download. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `file_type_categories` - (Optional) Filter based on the file type in download. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `casb_file_type_super_categories` - (Optional) Endpoint DLP file type category filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `casb_file_type` - (Optional) Endpoint DLP file type filter. See the [Cloud Nanolog Streaming Service (NSS) documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get)

* `file_sizes` - (Optional) File size filter

* `request_sizes` - (Optional) Request size filter

* `response_sizes` - (Optional) Request size filter

* `transaction_sizes` - (Optional) Transaction size filter

* `in_bound_bytes` - (Optional) Filter based on inbound bytes

* `out_bound_bytes` - (Optional) Filter based on outbound bytes

* `download_time` - (Optional) Download time filter

* `scan_time` - (Optional) Scan time filter

* `server_source_ips` - (Optional) Filter based on the server's source IPv4 addresses in Firewall policy

* `server_destination_ips` - (Optional) Filter based on the server's destination IPv4 addresses in Firewall policy

* `tunnel_ips` - (Optional) Filter based on tunnel IPv4 addresses in Firewall policy

* `internal_ips` - (Optional) Filter based on internal IPv4 addresses

* `tunnel_source_ips` - (Optional) Source IPv4 addresses of tunnels

* `tunnel_dest_ips` - (Optional) Destination IPv4 addresses of tunnels

* `client_destination_ips` - (Optional) Client's destination IPv4 addresses in Firewall policy

* `audit_log_type` - (Optional) Audit log type filter

* `project_name` - (Optional) Repository project name filter

* `repo_name` - (Optional) Repository name filter

* `object_name` - (Optional) CRM object name filter

* `channel_name` - (Optional) Collaboration channel name filter

* `file_source` - (Optional) Filter based on the file source

* `file_name` - (Optional) Filter based on the file name

* `session_counts` - (Optional) Firewall logs filter based on the number of sessions

* `adv_user_agents` - (Optional) Filter based on custom user agent strings

* `referer_urls` - (Optional) Referrer URL filter

* `host_names` - (Optional) Filter to limit the logs based on specific hostnames

* `full_urls` - (Optional) Filter to limit the logs based on specific full URLs

* `threat_names` - (Optional) Filter based on threat names

* `page_risk_indexes` - (Optional) Page Risk Index filter

* `client_destination_ports` - (Optional) Firewall logs filter based on a client's destination

* `tunnel_source_port` - (Optional) Filter based on the tunnel source port

## Block Attributes

### `casb_tenant` CASB tenant filter

* `id` - (Set of Integer) The unique identifiers for the CASB tenant.

### `locations` Location filter

* `id` - (Set of Integer) The unique identifiers for the locations.

### `location_groups` A filter based on location groups

* `id` - (Set of Integer) The unique identifiers for the location groups.

### `departments` Department filter

* `id` - (Set of Integer) The unique identifiers for the departments.

### `users` User filter

* `id` - (Set of Integer) The unique identifiers for the users.

### `sender_name` Filter based on sender or owner name

* `id` - (Set of Integer) The unique identifiers for the sender names.

### `buckets` Filter based on public cloud storage buckets

* `id` - (Set of Integer) The unique identifiers for the buckets.

### `external_owners` Filter logs associated with file owners (inside or outside your organization) who are not provisioned to ZIA services

* `id` - (Set of Integer) The unique identifiers for the external owners.

### `external_collaborators`

* `id` - (Set of Integer) The unique identifiers for the external collaborators.

### `internal_collaborators` Filter logs to specific recipients within your organization

* `id` - (Set of Integer) The unique identifiers for the internal collaborators.

### `itsm_object_type` ITSM object type filter

* `id` - (Set of Integer) The unique identifiers for the ITSM object types.

### `dlp_engines` DLP engine filter

* `id` - (Set of Integer) The unique identifiers for the DLP engines.
  ~> **NOTE** When associating a DLP Engine, you can use the `zia_dlp_engines` resource or data source.

### `dlp_dictionaries` DLP dictionary filter

* `id` - (Set of Integer) The unique identifiers for the DLP dictionaries.
  ~> **NOTE** When associating a DLP Dictionary, you can use the `zia_dlp_dictionaries` resource or data source.

### `url_categories` -  URL category filter

* `id` - (Set of Integer) Identifier that uniquely identifies an entity
  ~> **NOTE** When associating a URL category, you can use the `zia_url_categories` resource or data source; however, you must export the attribute `val`

### `vpn_credentials` Filter based on specific VPN credentials

* `id` - (Set of Integer) The unique identifiers for the VPN credentials.

### `rules` Policy rules filter (e.g., Firewall Filtering or DNS Control rule filter)

* `id` - (Set of Integer) The unique identifiers for the rules.

### `nw_services` Network services filter

* `id` - (Set of Integer) The unique identifiers for the network services.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Cloud NSS feeds can be imported using the `id` or `name`, e.g.

```shell
terraform import zia_cloud_nss_feed.example 123456789
```

or

```shell
terraform import zia_cloud_nss_feed.example "Splunk_Audit_Feed_Terraform"
```
