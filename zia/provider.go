package zia

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("ZIA_USERNAME"),
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("ZIA_PASSWORD"),
				Required:    true,
				Sensitive:   true,
			},
			"api_key": {
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("ZIA_API_KEY"),
				Required:    true,
				Sensitive:   true,
			},
			"zia_cloud": {
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("ZIA_CLOUD"),
				ValidateFunc: validation.StringInSlice([]string{
					"zscaler",
					"zscalerone",
					"zscalertwo",
					"zscalerthree",
					"zscloud",
					"zscalerbeta",
					"zscalergov",
				}, false),
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"zia_admin_users":                                   resourceAdminUsers(),
			"zia_dlp_dictionaries":                              resourceDLPDictionaries(),
			"zia_dlp_notification_templates":                    resourceDLPNotificationTemplates(),
			"zia_dlp_web_rules":                                 resourceDlpWebRules(),
			"zia_firewall_filtering_rule":                       resourceFirewallFilteringRules(),
			"zia_firewall_filtering_destination_groups":         resourceFWIPDestinationGroups(),
			"zia_firewall_filtering_ip_source_groups":           resourceFWIPSourceGroups(),
			"zia_firewall_filtering_network_service":            resourceFWNetworkServices(),
			"zia_firewall_filtering_network_service_groups":     resourceFWNetworkServiceGroups(),
			"zia_firewall_filtering_network_application_groups": resourceFWNetworkApplicationGroups(),
			"zia_traffic_forwarding_gre_tunnel":                 resourceTrafficForwardingGRETunnel(),
			"zia_traffic_forwarding_static_ip":                  resourceTrafficForwardingStaticIP(),
			"zia_traffic_forwarding_vpn_credentials":            resourceTrafficForwardingVPNCredentials(),
			"zia_location_management":                           resourceLocationManagement(),
			"zia_url_categories":                                resourceURLCategories(),
			"zia_url_filtering_rules":                           resourceURLFilteringRules(),
			"zia_user_management":                               resourceUserManagement(),
			"zia_activation_status":                             resourceActivationStatus(),
			"zia_rule_labels":                                   resourceRuleLabels(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"zia_admin_users":                                   dataSourceAdminUsers(),
			"zia_admin_roles":                                   dataSourceAdminRoles(),
			"zia_user_management":                               dataSourceUserManagement(),
			"zia_group_management":                              dataSourceGroupManagement(),
			"zia_department_management":                         dataSourceDepartmentManagement(),
			"zia_firewall_filtering_rule":                       dataSourceFirewallFilteringRule(),
			"zia_firewall_filtering_network_service":            dataSourceFWNetworkServices(),
			"zia_firewall_filtering_network_service_groups":     dataSourceFWNetworkServiceGroups(),
			"zia_firewall_filtering_network_application":        dataSourceFWNetworkApplication(),
			"zia_firewall_filtering_network_application_groups": dataSourceFWNetworkApplicationGroups(),
			"zia_firewall_filtering_ip_source_groups":           dataSourceFWIPSourceGroups(),
			"zia_firewall_filtering_destination_groups":         dataSourceFWIPDestinationGroups(),
			"zia_firewall_filtering_time_window":                dataSourceFWTimeWindow(),
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
			"zia_dlp_dictionaries":                              dataSourceDLPDictionaries(),
			"zia_dlp_notification_templates":                    dataSourceDLPNotificationTemplates(),
			"zia_dlp_web_rules":                                 dataSourceDlpWebRules(),
			"zia_dlp_engines":                                   dataSourceDLPEngines(),
			"zia_device_groups":                                 dataSourceDeviceGroups(),
			"zia_devices":                                       dataSourceDevices(),
			"zia_rule_labels":                                   dataSourceRuleLabels(),
			"zia_activation_status":                             dataSourceActivationStatus(),
		},

		ConfigureFunc: ziaConfigure,
	}
}

func ziaConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing ZIA client")
	config := Config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		APIKey:   d.Get("api_key").(string),
		ZIACloud: d.Get("zia_cloud").(string),
	}

	return config.Client()
}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		return nil, nil
	}
}
