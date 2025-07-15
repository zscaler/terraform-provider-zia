package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ZIAProvider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "zpa client id",
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "zpa client secret",
				ConflictsWith: []string{"private_key"},
			},
			"private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "zpa private key",
				ConflictsWith: []string{"client_secret"},
			},
			"vanity_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Zscaler Vanity Domain",
			},
			"zscaler_cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Zscaler Cloud Name",
			},
			"sandbox_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Zscaler Sandbox Token",
			},
			"sandbox_cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Zscaler Sandbox Cloud",
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"zia_cloud": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"zscaler",
					"zscalerone",
					"zscalertwo",
					"zscalerthree",
					"zscloud",
					"zscalerbeta",
					"zscalergov",
					"zscalerten",
					"zspreview",
				}, false),
				Optional: true,
			},
			"use_legacy_client": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
			},
			"http_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alternate HTTP proxy of scheme://hostname or scheme://hostname:port format",
			},
			"max_retries": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intAtMost(100),
				Description:      "maximum number of retries to attempt before erroring out.",
			},
			"parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible. Take note of https://help.zscaler.com/oneapi/understanding-rate-limiting.",
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 300),
				Description:      "Timeout for single request (in seconds) which is made to Zscaler, the default is `0` (means no limit is set). The maximum value can be `300`.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"zia_admin_users":                        resourceAdminUsers(),
			"zia_admin_roles":                        resourceAdminRoles(),
			"zia_browser_control_policy":             resourceBrowserControlPolicy(),
			"zia_bandwidth_classes":                  resourceBandwdithClasses(),
			"zia_bandwidth_classes_web_conferencing": resourceBandwdithClassesWebConferencing(),
			"zia_bandwidth_classes_file_size":        resourceBandwdithClassesFileSize(),
			"zia_dlp_dictionaries":                   resourceDLPDictionaries(),
			"zia_dlp_engines":                        resourceDLPEngines(),
			"zia_dlp_notification_templates":         resourceDLPNotificationTemplates(),
			"zia_dlp_web_rules":                      resourceDlpWebRules(),
			"zia_firewall_filtering_rule":            resourceFirewallFilteringRules(),
			"zia_firewall_ips_rule":                  resourceFirewallIPSRules(),
			"zia_firewall_dns_rule":                  resourceFirewallDNSRules(),
			"zia_cloud_app_control_rule":             resourceCloudAppControlRules(),
			"zia_casb_dlp_rules":                     resourceCasbDlpRules(),
			"zia_casb_malware_rules":                 resourceCasbMalwareRules(),
			"zia_risk_profiles":                      resourceRiskProfiles(),
			"zia_cloud_application_instance":         resourceCloudApplicationInstance(),
			//"zia_tenant_restriction_profile":                    resourceTenantRestrictionProfile(),
			"zia_firewall_filtering_destination_groups":         resourceFWIPDestinationGroups(),
			"zia_firewall_filtering_ip_source_groups":           resourceFWIPSourceGroups(),
			"zia_firewall_filtering_network_service":            resourceFWNetworkServices(),
			"zia_firewall_filtering_network_service_groups":     resourceFWNetworkServiceGroups(),
			"zia_firewall_filtering_network_application_groups": resourceFWNetworkApplicationGroups(),
			"zia_forwarding_control_rule":                       resourceForwardingControlRule(),
			"zia_nat_control_rules":                             resourceNatControlRules(),
			"zia_traffic_forwarding_gre_tunnel":                 resourceTrafficForwardingGRETunnel(),
			"zia_traffic_forwarding_static_ip":                  resourceTrafficForwardingStaticIP(),
			"zia_traffic_forwarding_vpn_credentials":            resourceTrafficForwardingVPNCredentials(),
			"zia_forwarding_control_zpa_gateway":                resourceForwardingControlZPAGateway(),
			"zia_location_management":                           resourceLocationManagement(),
			"zia_url_categories":                                resourceURLCategories(),
			"zia_url_filtering_rules":                           resourceURLFilteringRules(),
			"zia_file_type_control_rules":                       resourceFileTypeControlRules(),
			"zia_user_management":                               resourceUserManagement(),
			"zia_activation_status":                             resourceActivationStatus(),
			"zia_rule_labels":                                   resourceRuleLabels(),
			"zia_auth_settings_urls":                            resourceAuthSettingsUrls(),
			"zia_security_settings":                             resourceSecurityPolicySettings(),
			"zia_sandbox_behavioral_analysis":                   resourceSandboxSettings(),
			"zia_sandbox_file_submission":                       resourceSandboxSubmission(),
			"zia_sandbox_rules":                                 resourceSandboxRules(),
			"zia_ssl_inspection_rules":                          resourceSSLInspectionRules(),
			"zia_advanced_threat_settings":                      resourceAdvancedThreatSettings(),
			"zia_atp_malicious_urls":                            resourceATPMaliciousUrls(),
			"zia_atp_security_exceptions":                       resourceATPSecurityExceptions(),
			"zia_advanced_settings":                             resourceAdvancedSettings(),
			"zia_atp_malware_inspection":                        resourceATPMalwareInspection(),
			"zia_atp_malware_protocols":                         resourceATPMalwareProtocols(),
			"zia_atp_malware_settings":                          resourceATPMalwareSettings(),
			"zia_atp_malware_policy":                            resourceATPMalwarePolicy(),
			"zia_url_filtering_and_cloud_app_settings":          resourceURLFilteringCloludAppSettings(),
			"zia_end_user_notification":                         resourceEndUserNotification(),
			"zia_nss_server":                                    resourceNSSServer(),
			"zia_subscription_alert":                            resourceSubscriptionAlerts(),
			"zia_forwarding_control_proxies":                    resourceForwardingControlProxies(),
			"zia_ftp_control_policy":                            resourceFTPControlPolicy(),
			"zia_mobile_malware_protection_policy":              resourceMobileMalwareProtectionPolicy(),
			"zia_virtual_service_edge_cluster":                  resourceVZENCluster(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"zia_admin_users":                                   dataSourceAdminUsers(),
			"zia_admin_roles":                                   dataSourceAdminRoles(),
			"zia_user_management":                               dataSourceUserManagement(),
			"zia_group_management":                              dataSourceGroupManagement(),
			"zia_department_management":                         dataSourceDepartmentManagement(),
			"zia_browser_control_policy":                        dataSourceBrowserControlPolicy(),
			"zia_bandwidth_classes":                             dataBandwdithClasses(),
			"zia_bandwidth_control_rule":                        dataBandwdithControlRules(),
			"zia_cloud_applications":                            dataSourceCloudApplications(),
			"zia_cloud_app_control_rule":                        dataSourceCloudAppControlRules(),
			"zia_file_type_control_rules":                       dataSourceFileTypeControlRules(),
			"zia_firewall_filtering_rule":                       dataSourceFirewallFilteringRule(),
			"zia_firewall_filtering_network_service":            dataSourceFWNetworkServices(),
			"zia_firewall_filtering_network_service_groups":     dataSourceFWNetworkServiceGroups(),
			"zia_firewall_filtering_network_application":        dataSourceFWNetworkApplication(),
			"zia_firewall_filtering_network_application_groups": dataSourceFWNetworkApplicationGroups(),
			"zia_firewall_filtering_application_services":       dataSourceFWApplicationServicesLite(),
			"zia_firewall_filtering_application_services_group": dataSourceFWApplicationServicesGroupLite(),
			"zia_firewall_filtering_ip_source_groups":           dataSourceFWIPSourceGroups(),
			"zia_firewall_filtering_destination_groups":         dataSourceFWIPDestinationGroups(),
			"zia_firewall_filtering_time_window":                dataSourceFWTimeWindow(),
			"zia_firewall_ips_rule":                             dataSourceFirewallIPSRules(),
			"zia_firewall_dns_rule":                             dataSourceFirewallDNSRules(),
			"zia_forwarding_control_rule":                       dataSourceForwardingControlRule(),
			"zia_nat_control_rules":                             dataSourceNatControlRules(),
			"zia_url_categories":                                dataSourceURLCategories(),
			"zia_url_filtering_rules":                           dataSourceURLFilteringRules(),
			"zia_traffic_forwarding_public_node_vips":           dataSourceTrafficForwardingPublicNodeVIPs(),
			"zia_traffic_forwarding_vpn_credentials":            dataSourceTrafficForwardingVPNCredentials(),
			"zia_traffic_forwarding_gre_vip_recommended_list":   dataSourceTrafficForwardingGreVipRecommendedList(),
			"zia_traffic_forwarding_static_ip":                  dataSourceTrafficForwardingStaticIP(),
			"zia_traffic_forwarding_gre_tunnel":                 dataSourceTrafficForwardingGreTunnels(),
			"zia_traffic_forwarding_gre_tunnel_info":            dataSourceTrafficForwardingIPGreTunnelInfo(),
			"zia_gre_internal_ip_range_list":                    dataSourceTrafficForwardingGreInternalIPRangeList(),
			"zia_location_management":                           dataSourceLocationManagement(),
			"zia_location_groups":                               dataSourceLocationGroup(),
			"zia_location_lite":                                 dataSourceLocationLite(),
			"zia_dlp_dictionaries":                              dataSourceDLPDictionaries(),
			"zia_dlp_dictionary_predefined_identifiers":         dataSourceDLPDictionaryPredefinedIdentifiers(),
			"zia_dlp_engines":                                   dataSourceDLPEngines(),
			"zia_dlp_icap_servers":                              dataSourceDLPICAPServers(),
			"zia_dlp_edm_schema":                                dataSourceDLPEDMSchema(),
			"zia_dlp_idm_profiles":                              dataSourceDLPIDMProfiles(),
			"zia_dlp_idm_profile_lite":                          dataSourceDLPIDMProfileLite(),
			"zia_dlp_incident_receiver_servers":                 dataSourceDLPIncidentReceiverServers(),
			"zia_dlp_notification_templates":                    dataSourceDLPNotificationTemplates(),
			"zia_dlp_web_rules":                                 dataSourceDlpWebRules(),
			"zia_domain_profiles":                               dataSourceDomainProfiles(),
			"zia_casb_email_label":                              dataSourceCasbEmailLabel(),
			"zia_casb_dlp_rules":                                dataSourceCasbDlpRules(),
			"zia_casb_malware_rules":                            dataSourceCasbMalwareRules(),
			"zia_casb_tenant":                                   dataSourceCasbTenant(),
			"zia_casb_tombstone_template":                       dataSourceCasbTombstoneTemplate(),
			"zia_risk_profiles":                                 dataSourceRiskProfiles(),
			"zia_cloud_application_instance":                    dataSourceCloudApplicationInstance(),
			"zia_tenant_restriction_profile":                    dataSourceTenantRestrictionProfile(),
			"zia_device_groups":                                 dataSourceDeviceGroups(),
			"zia_devices":                                       dataSourceDevices(),
			"zia_rule_labels":                                   dataSourceRuleLabels(),
			"zia_activation_status":                             dataSourceActivationStatus(),
			"zia_auth_settings_urls":                            dataSourceAuthSettingsUrls(),
			"zia_security_settings":                             dataSourceSecurityPolicySettings(),
			"zia_sandbox_behavioral_analysis":                   dataSourceSandboxSettings(),
			"zia_sandbox_report":                                dataSourceSandboxReport(),
			"zia_sandbox_rules":                                 dataSourceSandboxRules(),
			"zia_ssl_inspection_rules":                          dataSourceSSLInspectionRules(),
			"zia_forwarding_control_zpa_gateway":                dataSourceForwardingControlZPAGateway(),
			"zia_forwarding_control_proxy_gateway":              dataSourceForwardingControlProxyGateway(),
			"zia_cloud_browser_isolation_profile":               dataSourceCBIProfile(),
			"zia_workload_groups":                               dataSourceWorkloadGroup(),
			"zia_advanced_threat_settings":                      dataSourceAdvancedThreatSettings(),
			"zia_atp_malicious_urls":                            dataSourceATPMaliciousUrls(),
			"zia_atp_security_exceptions":                       dataSourceATPSecurityExceptions(),
			"zia_advanced_settings":                             dataSourceAdvancedSettings(),
			"zia_atp_malware_inspection":                        dataSourceATPMalwareInspection(),
			"zia_atp_malware_protocols":                         dataSourceATPMalwareProtocols(),
			"zia_atp_malware_settings":                          dataSourceATPMalwareSettings(),
			"zia_atp_malware_policy":                            dataSourceATPMalwarePolicy(),
			"zia_url_filtering_and_cloud_app_settings":          dataSourceURLFilteringCloludAppSettings(),
			"zia_end_user_notification":                         dataSourceEndUserNotification(),
			"zia_nss_server":                                    dataSourceNSSServer(),
			"zia_subscription_alert":                            dataSourceSubscriptionAlerts(),
			"zia_forwarding_control_proxies":                    dataSourceForwardingControlProxies(),
			"zia_ftp_control_policy":                            dataSourceFTPControlPolicy(),
			"zia_mobile_malware_protection_policy":              dataSourceMobileMalwareProtectionPolicy(),
			"zia_virtual_service_edge_cluster":                  dataSourceVZENCluster(),
		},
	}

	p.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		r, err := providerConfigure(d, terraformVersion)
		if err != nil {
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "failed configuring the provider",
					Detail:        fmt.Sprintf("error:%v", err),
					AttributePath: cty.Path{},
				},
			}
		}
		return r, nil
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing Zscaler client")

	// Create configuration from schema
	config := NewConfig(d)
	config.TerraformVersion = terraformVersion

	// Load the correct SDK client (prioritizing V3)
	if diags := config.loadClients(); diags.HasError() {
		return nil, diags
	}

	// Return the configured client
	client, err := config.Client()
	if err != nil {
		return nil, diag.Errorf("failed to configure Zscaler client: %v", err)
	}

	// Initialize the global semaphore based on the configured parallelism
	if config.parallelism > 0 {
		apiSemaphore = make(chan struct{}, config.parallelism)
	} else {
		apiSemaphore = make(chan struct{}, 1)
	}

	return client, nil
}

func resourceFuncNoOp(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
