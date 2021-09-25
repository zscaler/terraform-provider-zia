package zscaler

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
			"zia_url": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ZIA_URL", nil),
				Optional:    true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{},

		DataSourcesMap: map[string]*schema.Resource{
			"zia_admin_role_mgmt": dataSourceAdminRoleMgmt(),
			"zia_dlp_dictionary":  dataSourceDLPDictionary(),
		},

		ConfigureFunc: zscalerConfigure,
	}
}

func zscalerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing ZIA client")
	config := Config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		APIKey:   d.Get("api_key").(string),
		ZIAUrl:   d.Get("zia_url").(string),
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
