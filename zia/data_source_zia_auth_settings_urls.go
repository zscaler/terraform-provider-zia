package zia

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAuthSettingsUrls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAuthSettingsUrlsRead,
		Schema: map[string]*schema.Schema{
			"urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAuthSettingsUrlsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	res, err := zClient.user_authentication_settings.Get()
	if err != nil {
		return err
	}
	d.SetId("exempted_urls")
	_ = d.Set("urls", res.URLs)
	return nil
}
