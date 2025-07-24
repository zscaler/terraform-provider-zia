package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudnss/cloudnss"
)

func dataSourceCloudNSSFeed() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudNSSFeedRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the nss server",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the cloud NSS feed",
			},
			"feed_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the feed",
			},
			"nss_log_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of NSS logs that are streamed (e.g. Web, Firewall, DNS, Alert, etc.)",
			},
			"nss_feed_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NSS feed format type (e.g. CSV, syslog, Splunk Common Information Model (CIM), etc.",
			},
			"feed_output_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output format used for the feed",
			},
			"user_obfuscation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies whether user obfuscation is enabled or disabled",
			},
			"time_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the time zone that must be used in the output file",
			},
			"custom_escaped_character": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Characters that need to be encoded using hex when they appear in URL, Host, or Referrer",
			},
			"eps_rate_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Event per second limit",
			},
			"json_array_toggle": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether streaming of logs in JSON array format (e.g., [{JSON1},{JSON2}]) is enabled or disabled for the JSON feed output type",
			},
			"siem_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud NSS SIEM type",
			},
			"max_batch_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum batch size in KB",
			},
			"connection_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The HTTPS URL of the SIEM log collection API endpoint",
			},
			"authentication_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authentication token value",
			},
			"connection_headers": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The HTTP Connection headers",
			},
			"last_success_full_test": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timestamp of the last successful test. Value is in Unix time.",
			},
			"test_connectivity_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The code from the last test",
			},
			"base64_encoded_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Base64-encoded certificate",
			},
			"nss_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NSS type",
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Client ID applicable when SIEM type is set to S3 or Azure Sentinel",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Client secret applicable when SIEM type is set to S3 or Azure Sentinel",
			},
			"authentication_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Authentication URL applicable when SIEM type is set to Azure Sentinel",
			},
			"grant_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Grant type applicable when SIEM type is set to Azure Sentinel",
			},
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scope applicable when SIEM type is set to Azure Sentinel",
			},
			"oauth_authentication": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether OAuth 2.0 authentication is enabled or not",
			},
			"server_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs based on the server's IPv4 addresses",
			},
			"client_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs based on a client's public IPv4 addresses",
			},
			"domains": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs to sessions associated with specific domains",
			},
			"dns_request_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "DNS request types included in the feed",
			},
			"dns_response_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "DNS response types filter",
			},
			"dns_responses": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "DNS responses filter",
			},
			"durations": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on time durations",
			},
			"dns_actions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "DNS Control policy action filter",
			},
			"firewall_logging_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filter based on the Firewall Filtering policy logging mode",
			},
			"client_source_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Client source IPs configured for NSS feed.",
			},
			"firewall_actions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Firewall actions included in the NSS feed.",
			},
			"countries": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Countries filter in the Firewall policy",
			},
			"server_source_ports": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Firewall log filter based on the traffic destination name",
			},
			"client_source_ports": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Firewall log filter based on a client's source ports",
			},
			"action_filter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy action filter",
			},
			"email_dlp_policy_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action filter for Email DLP log type",
			},
			"direction": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Traffic direction filter specifying inbound or outbound",
			},
			"event": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CASB event filter",
			},
			"policy_reasons": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Policy reason filter",
			},
			"protocol_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Protocol types filter",
			},
			"user_agents": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Predefined user agents filter",
			},
			"request_methods": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request methods filter",
			},
			"casb_severity": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Zscaler's Cloud Access Security Broker (CASB) severity filter",
			},
			"casb_policy_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CASB policy type filter",
			},
			"casb_applications": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CASB application filter",
			},
			"casb_action": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CASB policy action filter",
			},
			"url_super_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL supercategory filter",
			},
			"web_applications": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter to include specific cloud applications in the logs.
				By default, all cloud applications are included in the logs.
				To obtain the list of cloud applications that can be specified in this attribute, use the GET /cloudApplications/lite request.`,
			},
			"web_applications_exclude": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter to exclude specific cloud applications from the logs.
				By default, no cloud applications is excluded from the logs.
				To obtain the list of cloud applications that can be specified in this attribute, use the GET /cloudApplications/lite request.`,
			},
			"web_application_classes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud application categories Filter",
			},
			"malware_names": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on malware names",
			},
			"malware_classes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Malware category filter",
			},
			"url_classes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL category filter",
			},
			"advanced_threats": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Advanced threats filter",
			},
			"response_codes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Advanced threats filter",
			},
			"nw_applications": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter to include specific network applications in the logs.
				By default, all network applications are included in the logs`,
			},
			"nw_applications_exclude": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Filter to include specific network applications in the logs.
				By default, no network application is excluded from the logs`,
			},
			"nat_actions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "NAT Control policy actions filter",
			},
			"traffic_forwards": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the firewall traffic forwarding method",
			},
			"web_traffic_forwards": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the web traffic forwarding method",
			},
			"tunnel_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tunnel type filter",
			},
			"alerts": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Alert filter",
			},
			"object_type": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CRM object type filter",
			},
			"activity": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CASB activity filter",
			},
			"object_type1": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CASB activity object type filter",
			},
			"object_type2": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CASB activity object type filter if applicable",
			},
			"end_point_dlp_log_type": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Endpoint DLP log type filter",
			},
			"email_dlp_log_type": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Email DLP record type filter",
			},
			"file_type_super_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the category of file type in download",
			},
			"file_type_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the file type in download",
			},
			"casb_file_type_super_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Endpoint DLP file type category filer",
			},
			"file_sizes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "File size filter",
			},
			"request_sizes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request size filter",
			},
			"response_sizes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request size filter",
			},
			"transaction_sizes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Transaction size filter",
			},
			"in_bound_bytes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on inbound bytes",
			},
			"out_bound_bytes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on outbound bytes",
			},
			"download_time": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Download time filter",
			},
			"scan_time": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Scan time filter",
			},
			"server_source_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the server's source IPv4 addresses in Firewall policy",
			},
			"server_destination_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the server's destination IPv4 addresses in Firewall policy",
			},
			"tunnel_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on tunnel IPv4 addresses in Firewall policy",
			},
			"internal_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on internal IPv4 addresses",
			},
			"tunnel_source_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Source IPv4 addresses of tunnels",
			},
			"tunnel_dest_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Destination IPv4 addresses of tunnels",
			},
			"client_destination_ips": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Client's destination IPv4 addresses in Firewall policy",
			},
			"audit_log_type": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Audit log type filter",
			},
			"project_name": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Repository project name filter",
			},
			"repo_name": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Repository name filter",
			},
			"object_name": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CRM object name filter",
			},
			"channel_name": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Collaboration channel name filter",
			},
			"file_source": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the file source",
			},
			"file_name": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the file name",
			},
			"session_counts": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Firewall logs filter based on the number of sessions",
			},
			"adv_user_agents": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on custom user agent strings",
			},
			"referer_urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Referrer URL filter",
			},
			"host_names": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs based on specific hostnames",
			},
			"full_urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter to limit the logs based on specific full URLs",
			},
			"threat_names": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on threat names",
			},
			"page_risk_indexes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Page Risk Index filter",
			},
			"client_destination_ports": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Firewall logs filter based on a client's destination",
			},
			"tunnel_source_port": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based on the tunnel source port",
			},
			"external_owners": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"external_collaborators": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"internal_collaborators": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"itsm_object_type": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"url_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"dlp_engines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"dlp_dictionaries": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"nw_services": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Location filter",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"getl_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"location_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A filter based on location groups",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"getl_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"casb_tenant": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CASB tenant filter",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"getl_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CASB tenant filter",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"getl_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"departments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CASB tenant filter",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"getl_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"sender_name": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CASB tenant filter",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"getl_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"buckets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CASB tenant filter",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"getl_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"vpn_credentials": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CASB tenant filter",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"getl_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudNSSFeedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *cloudnss.NSSFeed
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for NSS Server id: %d\n", id)
		res, err := cloudnss.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for NSS Server name: %s\n", name)
		res, err := cloudnss.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("feed_status", resp.FeedStatus)
		_ = d.Set("nss_log_type", resp.NssLogType)
		_ = d.Set("nss_feed_type", resp.NssFeedType)
		_ = d.Set("feed_output_format", resp.FeedOutputFormat)
		_ = d.Set("user_obfuscation", resp.UserObfuscation)
		_ = d.Set("time_zone", resp.TimeZone)
		_ = d.Set("custom_escaped_character", resp.CustomEscapedCharacter)
		_ = d.Set("eps_rate_limit", resp.EpsRateLimit)
		_ = d.Set("json_array_toggle", resp.JsonArrayToggle)
		_ = d.Set("siem_type", resp.SiemType)
		_ = d.Set("max_batch_size", resp.MaxBatchSize)
		_ = d.Set("connection_url", resp.ConnectionURL)
		_ = d.Set("authentication_token", resp.AuthenticationToken)
		_ = d.Set("connection_headers", resp.ConnectionHeaders)
		_ = d.Set("last_success_full_test", resp.LastSuccessFullTest)
		_ = d.Set("test_connectivity_code", resp.TestConnectivityCode)
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

		// Handle all TypeList attributes using flattenCommonNSS
		if err := d.Set("casb_tenant", flattenCommonNSS(resp.CasbTenant)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("users", flattenCommonNSS(resp.Users)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("departments", flattenCommonNSS(resp.Departments)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("sender_name", flattenCommonNSS(resp.SenderName)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("buckets", flattenCommonNSS(resp.Buckets)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("vpn_credentials", flattenCommonNSS(resp.VPNCredentials)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("locations", flattenCommonNSS(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("location_groups", flattenCommonNSS(resp.LocationGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("external_owners", flattenIDExtensions(resp.ExternalOwners)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("external_collaborators", flattenIDExtensions(resp.ExternalCollaborators)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("internal_collaborators", flattenIDExtensions(resp.InternalCollaborators)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("itsm_object_type", flattenIDExtensions(resp.ItsmObjectType)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("url_categories", flattenIDExtensions(resp.URLCategories)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dlp_engines", flattenIDExtensions(resp.DLPEngines)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dlp_dictionaries", flattenIDExtensions(resp.DLPDictionaries)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("rules", flattenIDExtensions(resp.Rules)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("nw_services", flattenIDExtensions(resp.NwServices)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any Cloud NSS name '%s' or id '%d'", name, id))
	}

	return nil
}
