package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudnss/cloudnss"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

func resourceCloudNSSFeed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudNSSFeedCreate,
		ReadContext:   resourceCloudNSSFeedRead,
		UpdateContext: resourceCloudNSSFeedUpdate,
		DeleteContext: resourceCloudNSSFeedDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("nss_id", idInt)
				} else {
					resp, err := cloudnss.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("nss_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nss_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the cloud NSS feed",
			},
			"feed_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the feed",
			},
			"nss_log_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of NSS logs that are streamed (e.g. Web, Firewall, DNS, Alert, etc.)",
				ValidateFunc: validation.StringInSlice([]string{
					"ADMIN_AUDIT",
					"WEBLOG",
					"ALERT",
					"FWLOG",
					"DNSLOG",
					"MULTIFEEDLOG",
					"CASB_FILELOG",
					"CASB_MAILLOG",
					"ECLOG",
					"EC_DNSLOG",
					"CASB_ITSM",
					"CASB_CRM",
					"CASB_CODE_REPO",
					"CASB_COLLAB",
					"CASB_PCS",
					"USER_ACT_REP",
					"USER_COUNT_ALERT",
					"USER_IMP_TRAVEL_ALERT",
					"ENDPOINT_DLP",
					"EC_EVENTLOG",
					"EMAIL_DLP",
				}, false),
			},
			"nss_feed_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "NSS feed format type (e.g. CSV, syslog, Splunk Common Information Model (CIM), etc.",
				ValidateFunc: validation.StringInSlice([]string{
					"QRADAR",
					"SYSLOG",
					"CSV",
					"TAB_SEPARATED",
					"CUSTOM",
					"SPLUNK_CIM",
					"NAME_VALUE_PAIRS",
					"RSA_SECURITY",
					"ARCSIGHT_CEF",
					"SYMANTEC_MSS",
					"LOGRHYTHM",
					"ZBRIDGE",
					"MCAS",
					"JSON",
					"ZFAB_AGENT",
				}, false),
			},
			"feed_output_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Output format used for the feed",
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the time zone that must be used in the output file
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"custom_escaped_character": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Characters that need to be encoded using hex when they appear in URL, Host, or Referrer",
			},
			"eps_rate_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Event per second limit",
			},
			"json_array_toggle": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value indicating whether streaming of logs in JSON array format (e.g., [{JSON1},{JSON2}]) is enabled or disabled for the JSON feed output type",
			},
			"siem_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud NSS SIEM type",
				ValidateFunc: validation.StringInSlice([]string{
					"SPLUNK",
					"SUMO_LOGIC",
					"DEVO",
					"OTHER",
					"AZURE_SENTINEL",
					"S3",
				}, false),
			},
			"max_batch_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum batch size in KB",
			},
			"connection_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The HTTPS URL of the SIEM log collection API endpoint",
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"authentication_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The authentication token value",
			},
			"connection_headers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The HTTP Connection headers",
			},
			"base64_encoded_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Base64-encoded certificate",
			},
			"nss_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "NSS type",
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"SOFTWARE_AA_FLAG",
					"NSS_FOR_WEB",
					"NSS_FOR_FIREWALL",
					"VZEN",
					"VZEN_SME",
					"VZEN_SMLB",
					"PINNED_NSS",
					"MD5_CAPABLE",
					"ADP",
					"ZIRSVR",
					"NSS_FOR_ZPA",
				}, false),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client ID applicable when SIEM type is set to S3 or Azure Sentinel",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client secret applicable when SIEM type is set to S3 or Azure Sentinel",
			},
			"authentication_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Authentication URL applicable when SIEM type is set to Azure Sentinel",
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"grant_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Grant type applicable when SIEM type is set to Azure Sentinel",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Scope applicable when SIEM type is set to Azure Sentinel",
			},
			"oauth_authentication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value indicating whether OAuth 2.0 authentication is enabled or not",
			},
			"server_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs based on the server's IPv4 addresses",
			},
			"client_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs based on a client's public IPv4 addresses",
			},
			"domains": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs to sessions associated with specific domains",
			},
			"dns_request_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `DNS request types included in the feed
				Supported Values: ANY, NONE, DNSREQ_A, DNSREQ_NS, DNSREQ_CNAME, DNSREQ_SOA, DNSREQ_WKS,
				DNSREQ_PTR, DNSREQ_HINFO, DNSREQ_MINFO, DNSREQ_MX, DNSREQ_TXT, DNSREQ_AAAA,
				DNSREQ_ISDN, DNSREQ_LOC, DNSREQ_RP, DNSREQ_RT, DNSREQ_MR, DNSREQ_MG,
				DNSREQ_MB, DNSREQ_AFSDB, DNSREQ_HIP, DNSREQ_SRV, DNSREQ_DS, DNSREQ_NAPTR,
				DNSREQ_NSEC, DNSREQ_DNSKEY, DNSREQ_HTTPS, DNSREQ_UNKNOWN`,
			},
			"dns_response_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `DNS response types filter
				Supported Values: ANY, DNSRES_ZSCODE, DNSRES_CNAME, DNSRES_IPV6, DNSRES_SRV_CODE, DNSRES_IPV4`,
			},
			"dns_responses": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "DNS responses filter",
			},
			"durations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on time durations",
			},
			"dns_actions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "DNS Control policy action filter",
			},
			"firewall_logging_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter based on the Firewall Filtering policy logging mode",
				ValidateFunc: validation.StringInSlice([]string{
					"SESSION",
					"AGGREGATE",
					"ALL",
				}, false),
			},
			"client_source_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Client source IPs configured for NSS feed.",
			},
			"firewall_actions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Firewall actions included in the NSS feed
				Supported Values: BLOCK,ALLOW,BLOCK_DROP,BLOCK_RESET,BLOCK_ICMP,COUNTRY_BLOCK
				IPS_BLOCK_DROP,IPS_BLOCK_RESET,ALLOW_INSUFFICIENT_APPDATA,
				BLOCK_ABUSE_DROP,INT_ERR_DROP,CFG_BYPASSED,CFG_TIMEDOUT`,
			},
			"countries": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Countries filter in the Firewall policy
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"server_source_ports": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Firewall log filter based on the traffic destination name",
			},
			"client_source_ports": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Firewall log filter based on a client's source ports",
			},
			"action_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy action filter",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOWED",
					"BLOCKED",
				}, false),
			},
			"email_dlp_policy_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Action filter for Email DLP log type",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"CUSTOMHEADERINSERTION",
					"BLOCK",
				}, false),
			},
			"direction": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Traffic direction filter specifying inbound or outbound",
				ValidateFunc: validation.StringInSlice([]string{
					"INBOUND",
					"OUTBOUND",
				}, false),
			},
			"event": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CASB event filter",
				ValidateFunc: validation.StringInSlice([]string{
					"SCAN",
					"VIOLATION",
					"INCIDENT",
				}, false),
			},
			"policy_reasons": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Policy reason filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"protocol_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Protocol types filter
				Supported Values: TUNNEL, SSL, HTTP, HTTPS, FTP, FTPOVERHTTP, HTTP_PROXY, TUNNEL_SSL, DNSOVERHTTPS, WEBSOCKET, WEBSOCKET_SSL`,
			},
			"user_agents": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Predefined user agents filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"request_methods": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Request methods filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"casb_severity": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Zscaler's Cloud Access Security Broker (CASB) severity filter
				Supported Values: RULE_SEVERITY_HIGH, RULE_SEVERITY_MEDIUM, RULE_SEVERITY_LOW, RULE_SEVERITY_INFO`,
			},
			"casb_policy_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `CASB policy type filter
				Supported Values: MALWARE, DLP, ALL_INCIDENT`,
			},
			"casb_applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `CASB application filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"casb_action": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `CASB policy action filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"url_super_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `URL supercategory filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"web_applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter to include specific cloud applications in the logs.
				By default, all cloud applications are included in the logs.
				To obtain the list of cloud applications that can be specified in this attribute, use the GET /cloudApplications/lite request
				To retrieve the list of cloud applications, use the data source: zia_cloud_applications`,
			},
			"web_applications_exclude": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter to exclude specific cloud applications from the logs.
				By default, no cloud applications is excluded from the logs.
				To obtain the list of cloud applications that can be specified in this attribute, use the GET /cloudApplications/lite request.
				To retrieve the list of cloud applications, use the data source: zia_cloud_applications`,
			},
			"web_application_classes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Cloud application categories Filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"malware_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter based on malware names
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"malware_classes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Malware category filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"url_classes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `URL category filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"advanced_threats": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Advanced threats filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"response_codes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Response codes filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"nw_applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter to include specific network applications in the logs.
				By default, all network applications are included in the logs
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"nw_applications_exclude": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter to include specific network applications in the logs.
				By default, no network application is excluded from the logs
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"nat_actions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `NAT Control policy actions filter
				Supported Values: NONE, DNAT`,
			},
			"traffic_forwards": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter based on the firewall traffic forwarding method
				Supported Values: ANY, NONE, PBF, GRE, IPSEC, Z_APP, ZAPP_GRE, ZAPP_IPSEC, EC, MTGRE, ZAPP_DIRECT, CCA, MTUN_PROXY, MTUN_CBI`,
			},
			"web_traffic_forwards": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter based on the web traffic forwarding method
				Supported Values: ANY, NONE, PBF, GRE, IPSEC, Z_APP, ZAPP_GRE, ZAPP_IPSEC, EC, MTGRE, ZAPP_DIRECT, CCA, MTUN_PROXY, MTUN_CBI`,
			},
			"tunnel_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Tunnel type filter
				Supported Values: GRE, IPSEC_IKEV1, IPSEC_IKEV2, SVPN, EXTRANET, ZUB, ZCB`,
			},
			"alerts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Alert filter
				Supported Values: CRITICAL, WARN`,
			},
			"object_type": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CRM object type filter",
			},
			"activity": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `CASB activity filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"object_type1": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `CASB activity object type filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"object_type2": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `CASB activity object type filter if applicable
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"end_point_dlp_log_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Endpoint DLP log type filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"email_dlp_log_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Email DLP record type filter
				Supported Values: EPDLP_SCAN_AGGREGATE, EPDLP_SENSITIVE_ACTIVITY, EPDLP_DLP_INCIDENT`,
			},
			"file_type_super_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter based on the category of file type in download
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"file_type_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter based on the file type in download
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"casb_file_type_super_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Endpoint DLP file type category filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"casb_file_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Endpoint DLP file type filter
				See the [Cloud Nanolog Streaming Service (NSS) documentation
				https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssFeeds-get`,
			},
			"file_sizes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "File size filter",
			},
			"request_sizes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request size filter",
			},
			"response_sizes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request size filter",
			},
			"transaction_sizes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Transaction size filter",
			},
			"in_bound_bytes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on inbound bytes",
			},
			"out_bound_bytes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on outbound bytes",
			},
			"download_time": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Download time filter",
			},
			"scan_time": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Scan time filter",
			},
			"server_source_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the server's source IPv4 addresses in Firewall policy",
			},
			"server_destination_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the server's destination IPv4 addresses in Firewall policy",
			},
			"tunnel_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on tunnel IPv4 addresses in Firewall policy",
			},
			"internal_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on internal IPv4 addresses",
			},
			"tunnel_source_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Source IPv4 addresses of tunnels",
			},
			"tunnel_dest_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Destination IPv4 addresses of tunnels",
			},
			"client_destination_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Client's destination IPv4 addresses in Firewall policy",
			},
			"audit_log_type": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Audit log type filter",
			},
			"project_name": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Repository project name filter",
			},
			"repo_name": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Repository name filter",
			},
			"object_name": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CRM object name filter",
			},
			"channel_name": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Collaboration channel name filter",
			},
			"file_source": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the file source",
			},
			"file_name": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the file name",
			},
			"session_counts": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Firewall logs filter based on the number of sessions",
			},
			"adv_user_agents": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on custom user agent strings",
			},
			"referer_urls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Referrer URL filter",
			},
			"host_names": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs based on specific hostnames",
			},
			"full_urls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs based on specific full URLs",
			},
			"threat_names": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on threat names",
			},
			"page_risk_indexes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Page Risk Index filter",
			},
			"client_destination_ports": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Firewall logs filter based on a client's destination",
			},
			"tunnel_source_port": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the tunnel source port",
			},
			"casb_tenant":            setIDsSchemaTypeCustom(nil, "CASB tenant filter"),
			"locations":              setIDsSchemaTypeCustom(nil, "Location filter"),
			"location_groups":        setIDsSchemaTypeCustom(nil, "A filter based on location groups"),
			"departments":            setIDsSchemaTypeCustom(nil, "Departments filter"),
			"users":                  setIDsSchemaTypeCustom(nil, "Users filter"),
			"sender_name":            setIDsSchemaTypeCustom(nil, "Filter based on sender or owner name"),
			"buckets":                setIDsSchemaTypeCustom(nil, "Filter based on public cloud storage buckets"),
			"external_owners":        setIDsSchemaTypeCustom(nil, "Filter logs associated with file owners (inside or outside your organization) who are not provisioned to ZIA services"),
			"external_collaborators": setIDsSchemaTypeCustom(nil, "Filter logs to specific recipients outside your organization"),
			"internal_collaborators": setIDsSchemaTypeCustom(nil, "Filter logs to specific recipients within your organization"),
			"itsm_object_type":       setIDsSchemaTypeCustom(nil, "ITSM object type filter"),
			"dlp_engines":            setIDsSchemaTypeCustom(nil, "DLP engine filter"),
			"dlp_dictionaries":       setIDsSchemaTypeCustom(nil, "DLP dictionary filter"),
			"url_categories":         setIDsSchemaTypeCustom(nil, "URL category filter"),
			"vpn_credentials":        setIDsSchemaTypeCustom(nil, "Filter based on specific VPN credentials"),
			"rules":                  setIDsSchemaTypeCustom(nil, "Policy rules filter (e.g., Firewall Filtering or DNS Control rule filter)"),
			"nw_services":            setIDsSchemaTypeCustom(nil, "Firewall network services filter"),
		},
	}
}

func resourceCloudNSSFeedCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandCloudNSSFeed(d)
	log.Printf("[INFO] Creating ZIA cloud nss feeds\n%+v\n", req)

	resp, err := cloudnss.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA cloud nss feeds request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("nss_id", resp.ID)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceCloudNSSFeedRead(ctx, d, meta)
}

func resourceCloudNSSFeedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "nss_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no cloud nss feeds id is set"))
	}
	resp, err := cloudnss.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia cloud nss feeds %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia cloud nss feeds:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("nss_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("feed_status", resp.FeedStatus)
	_ = d.Set("nss_log_type", resp.NssLogType)
	_ = d.Set("nss_feed_type", resp.NssFeedType)
	_ = d.Set("feed_output_format", resp.FeedOutputFormat)
	_ = d.Set("time_zone", resp.TimeZone)
	_ = d.Set("custom_escaped_character", resp.CustomEscapedCharacter)
	_ = d.Set("eps_rate_limit", resp.EpsRateLimit)
	_ = d.Set("json_array_toggle", resp.JsonArrayToggle)
	_ = d.Set("siem_type", resp.SiemType)
	_ = d.Set("max_batch_size", resp.MaxBatchSize)
	_ = d.Set("connection_url", resp.ConnectionURL)
	_ = d.Set("authentication_token", resp.AuthenticationToken)
	_ = d.Set("connection_headers", resp.ConnectionHeaders)
	_ = d.Set("base64_encoded_certificate", resp.Base64EncodedCertificate)
	_ = d.Set("nss_type", resp.NssType)
	_ = d.Set("client_id", resp.ClientID)
	_ = d.Set("client_secret", resp.ClientSecret)
	_ = d.Set("authentication_url", resp.AuthenticationUrl)
	_ = d.Set("grant_type", resp.GrantType)
	_ = d.Set("scope", resp.Scope)
	_ = d.Set("oauth_authentication", resp.OauthAuthentication)
	_ = d.Set("server_ips", resp.ServerIps)
	_ = d.Set("client_ips", resp.ClientIps)
	_ = d.Set("domains", resp.Domains)
	_ = d.Set("dns_request_types", resp.DNSRequestTypes)
	_ = d.Set("dns_response_types", resp.DNSResponseTypes)
	_ = d.Set("dns_responses", resp.DNSResponses)
	_ = d.Set("durations", resp.Durations)
	_ = d.Set("dns_actions", resp.DNSActions)
	_ = d.Set("firewall_logging_mode", resp.FirewallLoggingMode)
	_ = d.Set("client_source_ips", resp.ClientSourceIps)
	_ = d.Set("firewall_actions", resp.FirewallActions)
	_ = d.Set("countries", resp.Countries)
	_ = d.Set("server_source_ports", resp.ServerSourcePorts)
	_ = d.Set("client_source_ports", resp.ClientSourcePorts)
	_ = d.Set("action_filter", resp.ActionFilter)
	_ = d.Set("email_dlp_policy_action", resp.EmailDlpPolicyAction)
	_ = d.Set("direction", resp.Direction)
	_ = d.Set("event", resp.Event)
	_ = d.Set("policy_reasons", resp.PolicyReasons)
	_ = d.Set("protocol_types", resp.ProtocolTypes)
	_ = d.Set("user_agents", resp.UserAgents)
	_ = d.Set("request_methods", resp.RequestMethods)
	_ = d.Set("casb_severity", resp.CasbSeverity)
	_ = d.Set("casb_policy_types", resp.CasbPolicyTypes)
	_ = d.Set("casb_applications", resp.CasbApplications)
	_ = d.Set("casb_action", resp.CasbAction)
	_ = d.Set("url_super_categories", resp.URLSuperCategories)
	_ = d.Set("web_applications", resp.WebApplications)
	_ = d.Set("web_application_classes", resp.WebApplicationClasses)
	_ = d.Set("malware_names", resp.MalwareNames)
	_ = d.Set("malware_classes", resp.MalwareClasses)
	_ = d.Set("url_classes", resp.URLClasses)
	_ = d.Set("advanced_threats", resp.AdvancedThreats)
	_ = d.Set("response_codes", resp.ResponseCodes)
	_ = d.Set("nw_applications", resp.NwApplications)
	_ = d.Set("nat_actions", resp.NatActions)
	_ = d.Set("traffic_forwards", resp.TrafficForwards)
	_ = d.Set("web_traffic_forwards", resp.WebTrafficForwards)
	_ = d.Set("tunnel_types", resp.TunnelTypes)
	_ = d.Set("alerts", resp.Alerts)
	_ = d.Set("object_type", resp.ObjectType)
	_ = d.Set("activity", resp.Activity)
	_ = d.Set("object_type1", resp.ObjectType1)
	_ = d.Set("object_type2", resp.ObjectType2)
	_ = d.Set("end_point_dlp_log_type", resp.EndPointDLPLogType)
	_ = d.Set("email_dlp_log_type", resp.EmailDLPLogType)
	_ = d.Set("file_type_super_categories", resp.FileTypeSuperCategories)
	_ = d.Set("file_type_categories", resp.FileTypeCategories)
	_ = d.Set("casb_file_type_super_categories", resp.CasbFileTypeSuperCategories)
	_ = d.Set("casb_file_type", resp.CasbFileType)
	_ = d.Set("file_sizes", resp.FileSizes)
	_ = d.Set("request_sizes", resp.RequestSizes)
	_ = d.Set("response_sizes", resp.ResponseSizes)
	_ = d.Set("transaction_sizes", resp.TransactionSizes)
	_ = d.Set("in_bound_bytes", resp.InBoundBytes)
	_ = d.Set("out_bound_bytes", resp.OutBoundBytes)
	_ = d.Set("download_time", resp.DownloadTime)
	_ = d.Set("scan_time", resp.ScanTime)
	_ = d.Set("server_source_ips", resp.ServerSourceIps)
	_ = d.Set("server_destination_ips", resp.ServerDestinationIps)
	_ = d.Set("tunnel_ips", resp.TunnelIps)
	_ = d.Set("internal_ips", resp.InternalIps)
	_ = d.Set("tunnel_source_ips", resp.TunnelSourceIps)
	_ = d.Set("tunnel_dest_ips", resp.TunnelDestIps)
	_ = d.Set("client_destination_ips", resp.ClientDestinationIps)
	_ = d.Set("audit_log_type", resp.AuditLogType)
	_ = d.Set("project_name", resp.ProjectName)
	_ = d.Set("repo_name", resp.RepoName)
	_ = d.Set("object_name", resp.ObjectName)
	_ = d.Set("channel_name", resp.ChannelName)
	_ = d.Set("file_source", resp.FileSource)
	_ = d.Set("file_name", resp.FileName)
	_ = d.Set("session_counts", resp.SessionCounts)
	_ = d.Set("adv_user_agents", resp.AdvUserAgents)
	_ = d.Set("referer_urls", resp.RefererUrls)
	_ = d.Set("host_names", resp.HostNames)
	_ = d.Set("full_urls", resp.FullUrls)
	_ = d.Set("threat_names", resp.ThreatNames)
	_ = d.Set("page_risk_indexes", resp.PageRiskIndexes)
	_ = d.Set("client_destination_ports", resp.ClientDestinationPorts)
	_ = d.Set("tunnel_source_port", resp.TunnelSourcePort)

	if err := d.Set("casb_tenant", flattenCommonNSSIDs(resp.CasbTenant)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("locations", flattenCommonNSSIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenCommonNSSIDs(resp.LocationGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("departments", flattenCommonNSSIDs(resp.Departments)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("users", flattenCommonNSSIDs(resp.Users)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("sender_name", flattenCommonNSSIDs(resp.SenderName)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("buckets", flattenCommonNSSIDs(resp.Buckets)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vpn_credentials", flattenCommonNSSIDs(resp.VPNCredentials)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("external_owners", flattenIDExtensionsListIDs(resp.ExternalOwners)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("external_collaborators", flattenIDExtensionsListIDs(resp.ExternalCollaborators)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("internal_collaborators", flattenIDExtensionsListIDs(resp.InternalCollaborators)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("itsm_object_type", flattenIDExtensionsListIDs(resp.ItsmObjectType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url_categories", flattenIDExtensionsListIDs(resp.URLCategories)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dlp_engines", flattenIDExtensionsListIDs(resp.DLPEngines)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dlp_dictionaries", flattenIDExtensionsListIDs(resp.DLPDictionaries)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rules", flattenIDExtensionsListIDs(resp.Rules)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nw_services", flattenIDExtensionsListIDs(resp.NwServices)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceCloudNSSFeedUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "nss_id")
	if !ok {
		log.Printf("[ERROR] cloud nss feed ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia cloud nss feed ID: %v\n", id)
	req := expandCloudNSSFeed(d)
	if _, err := cloudnss.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := cloudnss.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceCloudNSSFeedRead(ctx, d, meta)
}

func resourceCloudNSSFeedDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "nss_id")
	if !ok {
		log.Printf("[ERROR] cloud nss feed ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia cloud nss feed ID: %v\n", (d.Id()))

	if _, err := cloudnss.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia cloud nss feed deleted")

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandCloudNSSFeed(d *schema.ResourceData) cloudnss.NSSFeed {
	id, _ := getIntFromResourceData(d, "nss_id")
	result := cloudnss.NSSFeed{
		ID:                          id,
		Name:                        d.Get("name").(string),
		FeedStatus:                  d.Get("feed_status").(string),
		NssLogType:                  d.Get("nss_log_type").(string),
		NssFeedType:                 d.Get("nss_feed_type").(string),
		FeedOutputFormat:            d.Get("feed_output_format").(string),
		TimeZone:                    d.Get("time_zone").(string),
		EpsRateLimit:                d.Get("eps_rate_limit").(int),
		JsonArrayToggle:             d.Get("json_array_toggle").(bool),
		SiemType:                    d.Get("siem_type").(string),
		MaxBatchSize:                d.Get("max_batch_size").(int),
		ConnectionURL:               d.Get("connection_url").(string),
		AuthenticationToken:         d.Get("authentication_token").(string),
		Base64EncodedCertificate:    d.Get("base64_encoded_certificate").(string),
		NssType:                     d.Get("nss_type").(string),
		ClientID:                    d.Get("client_id").(string),
		ClientSecret:                d.Get("client_secret").(string),
		AuthenticationUrl:           d.Get("authentication_url").(string),
		GrantType:                   d.Get("grant_type").(string),
		Scope:                       d.Get("scope").(string),
		OauthAuthentication:         d.Get("oauth_authentication").(bool),
		FirewallLoggingMode:         d.Get("firewall_logging_mode").(string),
		ActionFilter:                d.Get("action_filter").(string),
		EmailDlpPolicyAction:        d.Get("email_dlp_policy_action").(string),
		Direction:                   d.Get("direction").(string),
		Event:                       d.Get("event").(string),
		CustomEscapedCharacter:      SetToStringList(d, "custom_escaped_character"),
		ConnectionHeaders:           SetToStringList(d, "connection_headers"),
		ServerIps:                   SetToStringList(d, "server_ips"),
		ClientIps:                   SetToStringList(d, "client_ips"),
		Domains:                     SetToStringList(d, "domains"),
		DNSRequestTypes:             SetToStringList(d, "dns_request_types"),
		DNSResponseTypes:            SetToStringList(d, "dns_response_types"),
		DNSResponses:                SetToStringList(d, "dns_responses"),
		Durations:                   SetToStringList(d, "durations"),
		DNSActions:                  SetToStringList(d, "dns_actions"),
		ClientSourceIps:             SetToStringList(d, "client_source_ips"),
		FirewallActions:             SetToStringList(d, "firewall_actions"),
		Countries:                   SetToStringList(d, "countries"),
		ServerSourcePorts:           SetToStringList(d, "server_source_ports"),
		ClientSourcePorts:           SetToStringList(d, "client_source_ports"),
		PolicyReasons:               SetToStringList(d, "policy_reasons"),
		ProtocolTypes:               SetToStringList(d, "protocol_types"),
		UserAgents:                  SetToStringList(d, "user_agents"),
		RequestMethods:              SetToStringList(d, "request_methods"),
		CasbSeverity:                SetToStringList(d, "casb_severity"),
		CasbPolicyTypes:             SetToStringList(d, "casb_policy_types"),
		CasbApplications:            SetToStringList(d, "casb_applications"),
		CasbAction:                  SetToStringList(d, "casb_action"),
		URLSuperCategories:          SetToStringList(d, "url_super_categories"),
		WebApplications:             SetToStringList(d, "web_applications"),
		WebApplicationClasses:       SetToStringList(d, "web_application_classes"),
		MalwareNames:                SetToStringList(d, "malware_names"),
		MalwareClasses:              SetToStringList(d, "malware_classes"),
		URLClasses:                  SetToStringList(d, "url_classes"),
		AdvancedThreats:             SetToStringList(d, "advanced_threats"),
		ResponseCodes:               SetToStringList(d, "response_codes"),
		NwApplications:              SetToStringList(d, "nw_applications"),
		NatActions:                  SetToStringList(d, "nat_actions"),
		TrafficForwards:             SetToStringList(d, "traffic_forwards"),
		WebTrafficForwards:          SetToStringList(d, "web_traffic_forwards"),
		TunnelTypes:                 SetToStringList(d, "tunnel_types"),
		Alerts:                      SetToStringList(d, "alerts"),
		ObjectType:                  SetToStringList(d, "object_type"),
		Activity:                    SetToStringList(d, "activity"),
		ObjectType1:                 SetToStringList(d, "object_type1"),
		ObjectType2:                 SetToStringList(d, "object_type2"),
		EndPointDLPLogType:          SetToStringList(d, "end_point_dlp_log_type"),
		EmailDLPLogType:             SetToStringList(d, "email_dlp_log_type"),
		FileTypeSuperCategories:     SetToStringList(d, "file_type_super_categories"),
		FileTypeCategories:          SetToStringList(d, "file_type_categories"),
		CasbFileTypeSuperCategories: SetToStringList(d, "casb_file_type_super_categories"),
		CasbFileType:                SetToStringList(d, "casb_file_type"),
		FileSizes:                   SetToStringList(d, "file_sizes"),
		RequestSizes:                SetToStringList(d, "request_sizes"),
		ResponseSizes:               SetToStringList(d, "response_sizes"),
		TransactionSizes:            SetToStringList(d, "transaction_sizes"),
		InBoundBytes:                SetToStringList(d, "in_bound_bytes"),
		OutBoundBytes:               SetToStringList(d, "out_bound_bytes"),
		DownloadTime:                SetToStringList(d, "download_time"),
		ScanTime:                    SetToStringList(d, "scan_time"),
		ServerSourceIps:             SetToStringList(d, "server_source_ips"),
		ServerDestinationIps:        SetToStringList(d, "server_destination_ips"),
		TunnelIps:                   SetToStringList(d, "tunnel_ips"),
		InternalIps:                 SetToStringList(d, "internal_ips"),
		TunnelSourceIps:             SetToStringList(d, "tunnel_source_ips"),
		TunnelDestIps:               SetToStringList(d, "tunnel_dest_ips"),
		ClientDestinationIps:        SetToStringList(d, "client_destination_ips"),
		AuditLogType:                SetToStringList(d, "audit_log_type"),
		ProjectName:                 SetToStringList(d, "project_name"),
		RepoName:                    SetToStringList(d, "repo_name"),
		ObjectName:                  SetToStringList(d, "object_name"),
		ChannelName:                 SetToStringList(d, "channel_name"),
		FileSource:                  SetToStringList(d, "file_source"),
		FileName:                    SetToStringList(d, "file_name"),
		SessionCounts:               SetToStringList(d, "session_counts"),
		AdvUserAgents:               SetToStringList(d, "adv_user_agents"),
		RefererUrls:                 SetToStringList(d, "referer_urls"),
		HostNames:                   SetToStringList(d, "host_names"),
		FullUrls:                    SetToStringList(d, "full_urls"),
		ThreatNames:                 SetToStringList(d, "threat_names"),
		PageRiskIndexes:             SetToStringList(d, "page_risk_indexes"),
		ClientDestinationPorts:      SetToStringList(d, "client_destination_ports"),
		TunnelSourcePort:            SetToStringList(d, "tunnel_source_port"),
		CasbTenant:                  expandCloudNSSIDSet(d, "casb_tenant"),
		Locations:                   expandCloudNSSIDSet(d, "locations"),
		LocationGroups:              expandCloudNSSIDSet(d, "location_groups"),
		Departments:                 expandCloudNSSIDSet(d, "departments"),
		Users:                       expandCloudNSSIDSet(d, "users"),
		SenderName:                  expandCloudNSSIDSet(d, "sender_name"),
		Buckets:                     expandCloudNSSIDSet(d, "buckets"),
		VPNCredentials:              expandCloudNSSIDSet(d, "vpn_credentials"),
		DLPEngines:                  expandIDNameExtensionsSet(d, "dlp_engines"),
		DLPDictionaries:             expandIDNameExtensionsSet(d, "dlp_dictionaries"),
		ExternalOwners:              expandIDNameExtensionsSet(d, "external_owners"),
		ExternalCollaborators:       expandIDNameExtensionsSet(d, "external_collaborators"),
		InternalCollaborators:       expandIDNameExtensionsSet(d, "internal_collaborators"),
		ItsmObjectType:              expandIDNameExtensionsSet(d, "itsm_object_type"),
		NwServices:                  expandIDNameExtensionsSet(d, "nw_services"),
		URLCategories:               expandIDNameExtensionsSet(d, "url_categories"),
		Rules:                       expandIDNameExtensionsSet(d, "rules"),
	}
	return result
}

func expandCloudNSSIDSet(d *schema.ResourceData, key string) []common.CommonNSS {
	setInterface, ok := d.GetOk(key)
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.CommonNSS
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil && itemMap["id"] != nil {
				set := itemMap["id"].(*schema.Set)
				for _, id := range set.List() {
					result = append(result, common.CommonNSS{
						ID: id.(int),
					})
				}
			}
		}
		return result
	}
	return []common.CommonNSS{}
}

func flattenCommonNSSIDs(list []common.CommonNSS) []interface{} {
	if len(list) == 0 {
		// Return an empty slice instead of nil
		return []interface{}{}
	}

	ids := []int{}
	for _, item := range list {
		if item.ID == 0 && item.Name == "" {
			continue
		}
		ids = append(ids, item.ID)
	}

	if len(ids) == 0 {
		// Again return []interface{}{} instead of nil
		return []interface{}{}
	}

	// The rest remains the same
	return []interface{}{
		map[string]interface{}{
			"id": ids,
		},
	}
}
