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
				DefaultFunc: schema.EnvDefaultFunc("ZIA_USERNAME", nil),
				Optional:    true,
			},
			"password": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ZIA_PASSWORD", nil),
				Optional:    true,
				Sensitive:   true,
			},
			"api_key": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ZIA_API_KEY", nil),
				Optional:    true,
				Sensitive:   true,
			},
			"zia_base_url": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ZIA_BASE_URL", nil),
				Optional:    true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{},

		DataSourcesMap: map[string]*schema.Resource{
			//"zia_admin_role_mgmt": dataSourceAdminRoleMgmt(),
			//"zia_dlp_dictionary":  dataSourceDLPDictionary(),
			"zia_user_management":              dataSourceUserManagement(),
			"zia_public_node_virtual_address":  dataSourcePublicNodeVirtualAddress(),
			"zia_gre_virtual_ip_address_list":  dataSourceGreVirtualIPAddressesList(),
			"zia_traffic_forwarding_static_ip": dataSourceTrafficForwardingStaticIP(),
			"zia_activation_status":            dataSourceActivationStatus(),
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
