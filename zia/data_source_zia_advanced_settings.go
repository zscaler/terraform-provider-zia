package zia

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/advanced_settings"
)

func dataSourceAdvancedSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedSettingsRead,
		Schema: map[string]*schema.Schema{
			"auth_bypass_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are exempted from cookie authentication",
			},
			"auth_bypass_urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom URLs that are exempted from cookie authentication for users",
			},
			"basic_bypass_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are exempted from Basic authentication",
			},
			"digest_auth_bypass_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are exempted from Digest authentication",
			},
			"dns_resolution_on_transparent_proxy_exempt_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are excluded from DNS optimization on transparent proxy mode",
			},
			"dns_resolution_on_transparent_proxy_ipv6_exempt_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are excluded from DNS optimization for IPv6 addresses on transparent proxy mode",
			},
			"dns_resolution_on_transparent_proxy_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications to which DNS optimization on transparent proxy mode applies",
			},
			"dns_resolution_on_transparent_proxy_ipv6_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications to which DNS optimization for IPv6 addresses on transparent proxy mode applies",
			},
			"block_domain_fronting_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Applications which are subjected to Domain Fronting",
			},
			"prefer_sni_over_conn_host_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Applications that are exempted from the preferSniOverConnHost setting (i.e., prefer SSL/TLS client hello SNI for DNS resolution instead of the CONNECT host for forward proxy connections)",
			},
			"dns_resolution_on_transparent_proxy_exempt_url_categories": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dns_resolution_on_transparent_proxy_ipv6_exempt_url_categories": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dns_resolution_on_transparent_proxy_url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL categories to which DNS optimization on transparent proxy mode applies",
			},
			"dns_resolution_on_transparent_proxy_ipv6_url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "IPv6 URL categories to which DNS optimization on transparent proxy mode applies",
			},
			"auth_bypass_url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL categories that are exempted from cookie authentication",
			},
			"domain_fronting_bypass_url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL categories that are exempted from domain fronting",
			},
			"kerberos_bypass_url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL categories that are exempted from Kerberos authentication",
			},
			"basic_bypass_url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL categories that are exempted from Basic authentication",
			},
			"http_range_header_remove_url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL categories for which HTTP range headers must be removed",
			},
			"digest_auth_bypass_url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL categories that are exempted from Digest authentication",
			},
			"sni_dns_optimization_bypass_url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL categories that are excluded from the preferSniOverConnHost setting (i.e., prefer SSL/TLS client hello SNI for DNS resolution instead of the CONNECT host for forward proxy connections)",
			},
			"kerberos_bypass_urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom URLs that are exempted from Kerberos authentication",
			},
			"kerberos_bypass_apps": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are exempted from Kerberos authentication",
			},
			"digest_auth_bypass_urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom URLs that are exempted from Digest authentication. Cloud applications that are exempted from Digest authentication",
			},
			"dns_resolution_on_transparent_proxy_exempt_urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URLs that are excluded from DNS optimization on transparent proxy mode",
			},
			"dns_resolution_on_transparent_proxy_urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URLs to which DNS optimization on transparent proxy mode applies",
			},
			"enable_dns_resolution_on_transparent_proxy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether DNS optimization is enabled or disabled for Z-Tunnel 2.0 and transparent proxy mode traffic (e.g., traffic via GRE or IPSec tunnels without a PAC file).",
			},
			"enable_ipv6_dns_resolution_on_transparent_proxy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether DNS optimization is enabled or disabled for IPv6 connections to dual-stack or IPv6-only destinations sent via Z-Tunnel 2.0 and transparent proxy proxy mode (e.g., traffic via GRE or IPSec tunnels without a PAC file).",
			},
			"enable_ipv6_dns_optimization_on_all_transparent_proxy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_evaluate_policy_on_global_ssl_bypass": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable DNS optimization for all IPv6 transparent proxy traffic",
			},
			"enable_office365": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether Microsoft Office 365 One Click Configuration is enabled or not",
			},
			"log_internal_ip": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether to log internal IP address present in X-Forwarded-For (XFF) proxy header or not",
			},
			"enforce_surrogate_ip_for_windows_app": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enforce Surrogate IP authentication for Windows app traffic",
			},
			"track_http_tunnel_on_http_ports": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether to apply configured policies on tunneled HTTP traffic sent via a CONNECT method request on port 80",
			},
			"block_http_tunnel_on_non_http_ports": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether HTTP CONNECT method requests to non-standard ports are allowed or not (i.e., requests directed to ports other than the standard HTTP/S ports 80 and 443)",
			},
			"block_domain_fronting_on_host_header": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether to block or allow HTTP/S transactions in which the FQDN of the request URL is different than the FQDN of the request's host header",
			},
			"zscaler_client_connector_1_and_pac_road_warrior_in_firewall": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether to apply the Firewall rules configured without a specified location criteria (or with the Road Warrior location) to remote user traffic forwarded via Z-Tunnel 1.0 or PAC files",
			},
			"cascade_url_filtering": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether to apply the URL Filtering policy even when the Cloud App Control policy already allows a transaction explicitly",
			},
			"enable_policy_for_unauthenticated_traffic": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether policies that include user and department criteria can be configured and applied for unauthenticated traffic",
			},
			"block_non_compliant_http_request_on_http_ports": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether to allow or block traffic that is not compliant with RFC HTTP protocol standards",
			},
			"enable_admin_rank_access": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether ranks are enabled for admins to allow admin ranks in policy configuration and management",
			},
			"http2_nonbrowser_traffic_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether or not HTTP/2 should be the default web protocol for accessing various applications at your organizational level",
			},
			"ecs_for_all_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether or not to include the ECS option in all DNS queries, originating from all locations and remote users.",
			},
			"dynamic_user_risk_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether to dynamically update user risk score by tracking risky user activities in real time",
			},
			"block_connect_host_sni_mismatch": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether CONNECT host and SNI mismatch (i.e., CONNECT host doesn't match the SSL/TLS client hello SNI) is blocked or not",
			},
			"prefer_sni_over_conn_host": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether or not to use the SSL/TLS client hello SNI for DNS resolution instead of the CONNECT host for forward proxy connections",
			},
			"sipa_xff_header_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether or not to insert XFF header to all traffic forwarded from ZIA to ZPA, including source IP-anchored and ZIA-inspected ZPA application traffic.",
			},
			"block_non_http_on_http_port_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Value indicating whether non-HTTP Traffic on HTTP/S ports are allowed or blocked",
			},
			"ui_session_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the login session timeout for admins accessing the ZIA Admin Portal",
			},
			// "ecs_object": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			// 			"name": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			// 			"external_id": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			// 		},
			// 	},
			// },

		},
	}
}

func dataSourceAdvancedSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Fetch data from the API
	res, err := advanced_settings.GetAdvancedSettings(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set ID for the data source
	d.SetId("advanced_settings")

	// Map values to the Terraform schema
	_ = d.Set("auth_bypass_urls", res.AuthBypassUrls)
	_ = d.Set("kerberos_bypass_urls", res.KerberosBypassUrls)
	_ = d.Set("digest_auth_bypass_urls", res.DigestAuthBypassUrls)
	_ = d.Set("dns_resolution_on_transparent_proxy_exempt_urls", res.DnsResolutionOnTransparentProxyExemptUrls)
	_ = d.Set("dns_resolution_on_transparent_proxy_urls", res.DnsResolutionOnTransparentProxyUrls)
	_ = d.Set("enable_dns_resolution_on_transparent_proxy", res.EnableDnsResolutionOnTransparentProxy)
	_ = d.Set("enable_ipv6_dns_resolution_on_transparent_proxy", res.EnableIPv6DnsResolutionOnTransparentProxy)
	_ = d.Set("enable_ipv6_dns_optimization_on_all_transparent_proxy", res.EnableIPv6DnsOptimizationOnAllTransparentProxy)
	_ = d.Set("enable_evaluate_policy_on_global_ssl_bypass", res.EnableEvaluatePolicyOnGlobalSSLBypass)
	_ = d.Set("enable_office365", res.EnableOffice365)
	_ = d.Set("log_internal_ip", res.LogInternalIp)
	_ = d.Set("enforce_surrogate_ip_for_windows_app", res.EnforceSurrogateIpForWindowsApp)
	_ = d.Set("track_http_tunnel_on_http_ports", res.TrackHttpTunnelOnHttpPorts)
	_ = d.Set("block_http_tunnel_on_non_http_ports", res.BlockHttpTunnelOnNonHttpPorts)
	_ = d.Set("block_domain_fronting_on_host_header", res.BlockDomainFrontingOnHostHeader)
	_ = d.Set("zscaler_client_connector_1_and_pac_road_warrior_in_firewall", res.ZscalerClientConnector1AndPacRoadWarriorInFirewall)
	_ = d.Set("cascade_url_filtering", res.CascadeUrlFiltering)
	_ = d.Set("enable_policy_for_unauthenticated_traffic", res.EnablePolicyForUnauthenticatedTraffic)
	_ = d.Set("block_non_compliant_http_request_on_http_ports", res.BlockNonCompliantHttpRequestOnHttpPorts)
	_ = d.Set("enable_admin_rank_access", res.EnableAdminRankAccess)
	_ = d.Set("http2_nonbrowser_traffic_enabled", res.Http2NonbrowserTrafficEnabled)
	_ = d.Set("ecs_for_all_enabled", res.EcsForAllEnabled)
	_ = d.Set("dynamic_user_risk_enabled", res.DynamicUserRiskEnabled)
	_ = d.Set("block_connect_host_sni_mismatch", res.BlockConnectHostSniMismatch)
	_ = d.Set("prefer_sni_over_conn_host", res.PreferSniOverConnHost)
	_ = d.Set("sipa_xff_header_enabled", res.SipaXffHeaderEnabled)
	_ = d.Set("block_non_http_on_http_port_enabled", res.BlockNonHttpOnHttpPortEnabled)
	_ = d.Set("ui_session_timeout", res.UISessionTimeout)
	_ = d.Set("auth_bypass_apps", res.AuthBypassApps)
	_ = d.Set("kerberos_bypass_apps", res.KerberosBypassApps)
	_ = d.Set("basic_bypass_apps", res.BasicBypassApps)
	_ = d.Set("digest_auth_bypass_apps", res.DigestAuthBypassApps)
	_ = d.Set("dns_resolution_on_transparent_proxy_exempt_apps", res.DnsResolutionOnTransparentProxyExemptApps)
	_ = d.Set("dns_resolution_on_transparent_proxy_ipv6_exempt_apps", res.DnsResolutionOnTransparentProxyIPv6ExemptApps)
	_ = d.Set("dns_resolution_on_transparent_proxy_apps", res.DnsResolutionOnTransparentProxyApps)
	_ = d.Set("dns_resolution_on_transparent_proxy_ipv6_apps", res.DnsResolutionOnTransparentProxyIPv6Apps)
	_ = d.Set("block_domain_fronting_apps", res.BlockDomainFrontingApps)
	_ = d.Set("prefer_sni_over_conn_host_apps", res.PreferSniOverConnHostApps)
	_ = d.Set("dns_resolution_on_transparent_proxy_exempt_url_categories", res.DnsResolutionOnTransparentProxyExemptUrlCategories)
	_ = d.Set("dns_resolution_on_transparent_proxy_ipv6_exempt_url_categories", res.DnsResolutionOnTransparentProxyIPv6ExemptUrlCategories)
	_ = d.Set("dns_resolution_on_transparent_proxy_url_categories", res.DnsResolutionOnTransparentProxyUrlCategories)
	_ = d.Set("dns_resolution_on_transparent_proxy_ipv6_url_categories", res.DnsResolutionOnTransparentProxyIPv6UrlCategories)
	_ = d.Set("auth_bypass_url_categories", res.AuthBypassUrlCategories)
	_ = d.Set("domain_fronting_bypass_url_categories", res.DomainFrontingBypassUrlCategories)
	_ = d.Set("kerberos_bypass_url_categories", res.KerberosBypassUrlCategories)
	_ = d.Set("basic_bypass_url_categories", res.BasicBypassUrlCategories)
	_ = d.Set("http_range_header_remove_url_categories", res.HttpRangeHeaderRemoveUrlCategories)
	_ = d.Set("digest_auth_bypass_url_categories", res.DigestAuthBypassUrlCategories)
	_ = d.Set("sni_dns_optimization_bypass_url_categories", res.SniDnsOptimizationBypassUrlCategories)

	return nil
}
