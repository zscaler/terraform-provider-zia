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
			"skip_credentials_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Skip credentials validation and API client initialization entirely. " +
					"Intended for configurations where every zia_* resource and data source is conditionally disabled " +
					"(e.g., count = 0) and no API call will ever be made, such as multi-environment deployments where " +
					"Zscaler is not present in every environment. Any resource or data source that does attempt an API " +
					"call will fail with an explanatory error. Can also be sourced from the " +
					"ZSCALER_SKIP_CREDENTIALS_VALIDATION environment variable.",
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
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "This attribute no longer has any effect and will be removed in a future major release. Remove it from the provider block. API rate limits are handled automatically: the provider honours the Retry-After header returned on a 429 response and retries transparently.",
				Description: "Deprecated and ignored. Previously limited the number of concurrent API requests. " +
					"Rate limiting is now handled automatically and this attribute has no effect.",
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 300),
				Description:      "Timeout for single request (in seconds) which is made to Zscaler, the default is `0` (means no limit is set). The maximum value can be `300`.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"zia_admin_users":                                   resourceAdminUsers(),
			"zia_admin_roles":                                   resourceAdminRoles(),
			"zia_bandwidth_control_rule":                        resourceBandwdithControlRules(),
			"zia_browser_control_policy":                        resourceBrowserControlPolicy(),
			"zia_bandwidth_classes":                             resourceBandwdithClasses(),
			"zia_bandwidth_classes_web_conferencing":            resourceBandwdithClassesWebConferencing(),
			"zia_bandwidth_classes_file_size":                   resourceBandwdithClassesFileSize(),
			"zia_dlp_dictionaries":                              resourceDLPDictionaries(),
			"zia_dlp_engines":                                   resourceDLPEngines(),
			"zia_dlp_notification_templates":                    resourceDLPNotificationTemplates(),
			"zia_dlp_web_rules":                                 resourceDlpWebRules(),
			"zia_endpoint_dlp_rules":                            resourceEndpointDLPRules(),
			"zia_endpoint_dlp_sub_rules":                        resourceEndpointDLPSubRules(),
			"zia_endpoint_dlp_resource":                         resourceEndpointDLPResource(),
			"zia_endpoint_dlp_resource_group":                   resourceEndpointDLPResourceGroup(),
			"zia_endpoint_dlp_application_group":                resourceEndpointDLPApplicationGroup(),
			"zia_endpoint_dlp_custom_apps":                      resourceEndpointDLPCustomApps(),
			"zia_outbound_email_dlp":                            resourceOutboundEmailDLP(),
			"zia_dlp_global_options":                            resourceDLPGlobalOptions(),
			"zia_firewall_filtering_rule":                       resourceFirewallFilteringRules(),
			"zia_firewall_ips_rule":                             resourceFirewallIPSRules(),
			"zia_firewall_dns_rule":                             resourceFirewallDNSRules(),
			"zia_cloud_app_control_rule":                        resourceCloudAppControlRules(),
			"zia_casb_dlp_rules":                                resourceCasbDlpRules(),
			"zia_casb_malware_rules":                            resourceCasbMalwareRules(),
			"zia_risk_profiles":                                 resourceRiskProfiles(),
			"zia_cloud_application_instance":                    resourceCloudApplicationInstance(),
			"zia_tenant_restriction_profile":                    resourceTenantRestrictionProfile(),
			"zia_firewall_filtering_destination_groups":         resourceFWIPDestinationGroups(),
			"zia_firewall_filtering_ip_source_groups":           resourceFWIPSourceGroups(),
			"zia_firewall_filtering_network_service":            resourceFWNetworkServices(),
			"zia_firewall_filtering_network_service_groups":     resourceFWNetworkServiceGroups(),
			"zia_firewall_filtering_network_application_groups": resourceFWNetworkApplicationGroups(),
			"zia_dns_application_groups":                        resourceDNSApplicationGroups(),
			"zia_forwarding_control_rule":                       resourceForwardingControlRule(),
			"zia_nat_control_rules":                             resourceNatControlRules(),
			"zia_traffic_forwarding_gre_tunnel":                 resourceTrafficForwardingGRETunnel(),
			"zia_traffic_forwarding_static_ip":                  resourceTrafficForwardingStaticIP(),
			"zia_traffic_forwarding_vpn_credentials":            resourceTrafficForwardingVPNCredentials(),
			"zia_forwarding_control_zpa_gateway":                resourceForwardingControlZPAGateway(),
			"zia_location_management":                           resourceLocationManagement(),
			"zia_url_categories":                                resourceURLCategories(),
			"zia_url_categories_predefined":                     resourceURLCategoriesPredefined(),
			"zia_url_filtering_rules":                           resourceURLFilteringRules(),
			"zia_file_type_control_rules":                       resourceFileTypeControlRules(),
			"zia_custom_file_types":                             resourceCustomFileTypes(),
			"zia_user_management":                               resourceUserManagement(),
			"zia_activation_status":                             resourceActivationStatus(),
			"zia_rule_labels":                                   resourceRuleLabels(),
			"zia_pac_files":                                     resourcePacFiles(),
			"zia_auth_settings_urls":                            resourceAuthSettingsUrls(),
			"zia_security_settings":                             resourceSecurityPolicySettings(),
			"zia_sandbox_behavioral_analysis":                   resourceSandboxSettings(),
			"zia_sandbox_behavioral_analysis_v2":                resourceSandboxSettingsV2(),
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
			"zia_cloud_nss_feed":                                resourceCloudNSSFeed(),
			"zia_nss_server":                                    resourceNSSServer(),
			"zia_subscription_alert":                            resourceSubscriptionAlerts(),
			"zia_forwarding_control_proxies":                    resourceForwardingControlProxies(),
			"zia_ftp_control_policy":                            resourceFTPControlPolicy(),
			"zia_mobile_malware_protection_policy":              resourceMobileMalwareProtectionPolicy(),
			"zia_virtual_service_edge_cluster":                  resourceVZENCluster(),
			"zia_virtual_service_edge_node":                     resourceVZENNode(),
			"zia_workload_groups":                               resourceWorkloadGroups(),
			"zia_traffic_capture_rules":                         resourceTrafficCaptureRules(),
			"zia_sub_cloud":                                     resourceSubCloud(),
			"zia_extranet":                                      resourceExtranet(),
			"zia_dc_exclusions":                                 resourceDCExclusions(),
			"zia_email_profile":                                 resourceEmailProfile(),
			"zia_ips_signature_rules":                           resourceIPSSignatureRules(),
			"zia_http_header_action_profile":                    resourceHttpHeaderActionProfile(),
			"zia_http_header_profile":                           resourceHttpHeaderProfile(),
			"zia_ueba_alert_definitions":                        resourceUEBAAlertDefinitions(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"zia_admin_users":                                   dataSourceAdminUsers(),
			"zia_admin_roles":                                   dataSourceAdminRoles(),
			"zia_user_management":                               dataSourceUserManagement(),
			"zia_group_management":                              dataSourceGroupManagement(),
			"zia_department_management":                         dataSourceDepartmentManagement(),
			"zia_browser_control_policy":                        dataSourceBrowserControlPolicy(),
			"zia_bandwidth_classes":                             dataSourceBandwdithClasses(),
			"zia_bandwidth_control_rule":                        dataSourceBandwdithControlRules(),
			"zia_cloud_applications":                            dataSourceCloudApplications(),
			"zia_cloud_app_control_rule":                        dataSourceCloudAppControlRules(),
			"zia_cloud_app_control_rule_actions":                dataSourceCloudAppControlRuleActions(),
			"zia_file_type_control_rules":                       dataSourceFileTypeControlRules(),
			"zia_custom_file_types":                             dataSourceCustomFileTypes(),
			"zia_file_type_categories":                          dataSourceFileTypeCategories(),
			"zia_firewall_filtering_rule":                       dataSourceFirewallFilteringRule(),
			"zia_firewall_filtering_network_service":            dataSourceFWNetworkServices(),
			"zia_firewall_filtering_network_service_groups":     dataSourceFWNetworkServiceGroups(),
			"zia_firewall_filtering_network_application":        dataSourceFWNetworkApplication(),
			"zia_firewall_filtering_network_application_groups": dataSourceFWNetworkApplicationGroups(),
			"zia_firewall_filtering_application_services":       dataSourceFWApplicationServicesLite(),
			"zia_firewall_filtering_application_services_group": dataSourceFWApplicationServicesGroupLite(),
			"zia_firewall_filtering_ip_source_groups":           dataSourceFWIPSourceGroups(),
			"zia_dns_application_groups":                        dataSourceDNSApplicationGroups(),
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
			"zia_endpoint_dlp_rules":                            dataSourceEndpointDLPRules(),
			"zia_outbound_email_dlp":                            dataSourceOutboundEmailDLP(),
			"zia_endpoint_dlp_custom_apps":                      dataSourceEndpointDLPCustomApps(),
			"zia_endpoint_dlp_application":                      dataSourceEndpointDLPApplication(),
			"zia_eun_template_product":                          dataSourceEUNTemplateProduct(),
			"zia_eun_user_confirmation_template_product":        dataSourceEUNUserConfirmationTemplateProduct(),
			"zia_dlp_endpoint_resource_channels":                dataSourceDLPEndpointResourceChannels(),
			"zia_dlp_endpoint_resource_group_tag":               dataSourceDLPEndpointResourceGroupTag(),
			"zia_dlp_cloud_to_cloud_ir":                         dataSourceDLPCloudToCloudIR(),
			"zia_dlp_global_options":                            dataSourceDLPGlobalOptions(),
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
			"zia_pac_files":                                     dataSourcePacFiles(),
			"zia_activation_status":                             dataSourceActivationStatus(),
			"zia_auth_settings_urls":                            dataSourceAuthSettingsUrls(),
			"zia_security_settings":                             dataSourceSecurityPolicySettings(),
			"zia_sandbox_behavioral_analysis":                   dataSourceSandboxSettings(),
			"zia_sandbox_behavioral_analysis_v2":                dataSourceSandboxSettingsV2(),
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
			"zia_cloud_nss_feed":                                dataSourceCloudNSSFeed(),
			"zia_nss_server":                                    dataSourceNSSServer(),
			"zia_subscription_alert":                            dataSourceSubscriptionAlerts(),
			"zia_forwarding_control_proxies":                    dataSourceForwardingControlProxies(),
			"zia_dedicated_ip_proxy":                            dataSourceDedicatedIPProxy(),
			"zia_ftp_control_policy":                            dataSourceFTPControlPolicy(),
			"zia_mobile_malware_protection_policy":              dataSourceMobileMalwareProtectionPolicy(),
			"zia_virtual_service_edge_cluster":                  dataSourceVZENCluster(),
			"zia_virtual_service_edge_node":                     dataSourceVZENNode(),
			"zia_traffic_capture_rules":                         dataSourceTrafficCaptureRules(),
			"zia_sub_cloud":                                     dataSourceSubCloud(),
			"zia_datacenters":                                   dataSourceDatacenters(),
			"zia_extranet":                                      dataSourceExtranet(),
			"zia_dc_exclusions":                                 dataSourceDCExclusions(),
			"zia_email_profile":                                 dataSourceEmailProfile(),
			"zia_ips_signature_rules":                           dataSourceIPSSignatureRules(),
			"zia_ips_categories":                                dataSourceIPSCategories(),
			"zia_supported_browser_version":                     dataSourceSupportedBrowserVersion(),
			"zia_adaptive_access_profile":                       dataSourceAdaptiveAccessProfile(),
			"zia_http_header_action_profile":                    dataSourceHttpHeaderActionProfile(),
			"zia_http_header_profile":                           dataSourceHttpHeaderProfile(),
			"zia_ueba_alert_definitions":                        dataSourceUEBAAlertDefinitions(),
		},
	}

	// Guard every resource and data source against the inert client returned
	// when skip_credentials_validation is enabled, so an accidental API call
	// yields a descriptive error instead of a nil-pointer panic.
	for _, r := range p.ResourcesMap {
		guardResourceAgainstInertClient(r)
	}
	for _, ds := range p.DataSourcesMap {
		guardResourceAgainstInertClient(ds)
	}

	p.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		r, diags := providerConfigure(d, terraformVersion)
		if diags.HasError() {
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "failed configuring the provider",
					Detail:        fmt.Sprintf("error:%v", diags),
					AttributePath: cty.Path{},
				},
			}
		}
		// Pass through non-error diagnostics (e.g. the
		// skip_credentials_validation warning).
		return r, diags
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing Zscaler client")

	// Create configuration from schema
	config := NewConfig(d)
	config.TerraformVersion = terraformVersion

	// Skip mode: never construct the SDK client. The OneAPI SDK authenticates
	// inside its constructor, so building a client without valid credentials
	// would fail at configure time even when no resource will ever call the
	// API. Return an inert client instead; the CRUD guards installed in
	// ZIAProvider turn any actual API use into a descriptive error.
	if config.skipCredentialsValidation {
		log.Printf("[WARN] skip_credentials_validation is enabled; ZIA API client was not initialized")
		return &Client{skipCredentialsValidation: true}, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "ZIA credentials were not validated",
				Detail: "skip_credentials_validation is enabled, so the ZIA API client was not initialized. " +
					"Any zia_* resource or data source that attempts an API call will fail. " +
					"This mode is intended for configurations where all ZIA resources are conditionally disabled (e.g., count = 0).",
			},
		}
	}

	// Load the correct SDK client (prioritizing V3)
	if diags := config.loadClients(); diags.HasError() {
		return nil, diags
	}

	// Return the configured client
	client, err := config.Client()
	if err != nil {
		return nil, diag.Errorf("failed to configure Zscaler client: %v", err)
	}

	return client, nil
}

// inertClientDiag is the error returned when a resource or data source is
// evaluated while the provider is in skip_credentials_validation mode.
func inertClientDiag() diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ZIA provider was configured with skip_credentials_validation",
			Detail: "This resource or data source attempted a ZIA API call, but the provider was configured with " +
				"skip_credentials_validation = true (or ZSCALER_SKIP_CREDENTIALS_VALIDATION=true), so no API client exists. " +
				"Either provide valid credentials and remove skip_credentials_validation, or ensure every zia_* resource " +
				"and data source is disabled (e.g., count = 0) in this configuration.",
		},
	}
}

// isInertClient reports whether the provider meta is the inert client created
// in skip_credentials_validation mode.
func isInertClient(meta interface{}) bool {
	client, ok := meta.(*Client)
	return ok && client.skipCredentialsValidation
}

// guardResourceAgainstInertClient wraps a resource's CRUD and import
// functions so that, in skip_credentials_validation mode, they return a
// descriptive error instead of dereferencing the nil SDK service.
func guardResourceAgainstInertClient(r *schema.Resource) {
	wrap := func(f func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
		if f == nil {
			return nil
		}
		return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			if isInertClient(meta) {
				return inertClientDiag()
			}
			return f(ctx, d, meta)
		}
	}

	r.CreateContext = wrap(r.CreateContext)
	r.ReadContext = wrap(r.ReadContext)
	r.UpdateContext = wrap(r.UpdateContext)
	r.DeleteContext = wrap(r.DeleteContext)

	if r.Importer != nil && r.Importer.StateContext != nil {
		importer := r.Importer.StateContext
		r.Importer.StateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			if isInertClient(meta) {
				return nil, fmt.Errorf("cannot import: the ZIA provider was configured with skip_credentials_validation, so no API client exists")
			}
			return importer(ctx, d, meta)
		}
	}
}

func resourceFuncNoOp(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
