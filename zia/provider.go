package zia

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"zia_base_url": {
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("ZIA_BASE_URL"),
				Required:    true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"zia_admin_users":                        resourceAdminUsers(),
			"zia_dlp_dictionaries":                   resourceDLPDictionaries(),
			"zia_firewall_filtering_rule":            resourceFirewallFilteringRules(),
			"zia_ip_destination_groups":              resourceIPDestinationGroups(),
			"zia_ip_source_groups":                   resourceIPSourceGroups(),
			"zia_traffic_forwarding_gre_tunnel":      resourceTrafficForwardingGRETunnel(),
			"zia_traffic_forwarding_static_ip":       resourceTrafficForwardingStaticIP(),
			"zia_traffic_forwarding_vpn_credentials": resourceTrafficForwardingVPNCredentials(),
			"zia_location_management":                resourceLocationManagement(),
			"zia_url_filtering_rules":                resourceURLFilteringRules(),
			"zia_activation_status":                  resourceActivationStatus(),
			//"zia_url_categories":                     resourceURLCategories(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"zia_admin_users":                         dataSourceAdminUsers(),
			"zia_admin_roles":                         dataSourceAdminRoles(),
			"zia_user_management":                     dataSourceUserManagement(),
			"zia_group_management":                    dataSourceGroupManagement(),
			"zia_department_management":               dataSourceDepartmentManagement(),
			"zia_firewall_filtering_rule":             dataSourceFirewallFilteringRule(),
			"zia_ip_destination_groups":               dataSourceIPDestinationGroups(),
			"zia_network_application_groups_lite":     dataSourceNetworkApplicationGroupsLite(),
			"zia_network_application":                 dataSourceNetworkApplication(),
			"zia_network_service_groups":              dataSourceNetworkServiceGroups(),
			"zia_network_services":                    dataSourceNetworkServices(),
			"zia_url_categories":                      dataSourceURLCategories(),
			"zia_url_filtering_rules":                 dataSourceURLFilteringRules(),
			"zia_ip_source_groups":                    dataSourceIPSourceGroups(),
			"zia_traffic_forwarding_public_node_vips": dataSourceTrafficForwardingPublicNodeVIPs(),
			"zia_gre_virtual_ip_address_list":         dataSourceGreVirtualIPAddressesList(),
			"zia_gre_internal_ip_ranges":              dataSourceTrafficForwardingGreInternalIPRanges(),
			"zia_traffic_forwarding_static_ip":        dataSourceTrafficForwardingStaticIP(),
			"zia_traffic_forwarding_gre_tunnel":       dataSourceTrafficForwardingGreTunnels(),
			"zia_location_management":                 dataSourceLocationManagement(),
			"zia_traffic_forwarding_vpn_credentials":  dataSourceTrafficForwardingVPNCredentials(),
			"zia_dlp_dictionaries":                    dataSourceDLPDictionaries(),
			"zia_activation_status":                   dataSourceActivationStatus(),
		},

		ConfigureFunc: ziaConfigure,
	}
}

func ziaConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing ZIA client")
	config := Config{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		APIKey:     d.Get("api_key").(string),
		ZIABaseURL: d.Get("zia_base_url").(string),
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
