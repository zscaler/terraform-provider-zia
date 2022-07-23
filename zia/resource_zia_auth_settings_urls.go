package zia

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zia/services/user_authentication_settings"
)

func resourceAuthSettingsUrls() *schema.Resource {
	return &schema.Resource{
		Read:   resourceAuthSettingsUrlsRead,
		Create: resourceAuthSettingsUrlsCreate,
		Update: resourceAuthSettingsUrlsUpdate,
		Delete: resourceAuthSettingsUrlsDelete,
		Schema: map[string]*schema.Schema{
			"urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				MaxItems: 25000,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAuthSettingsUrlsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	res, err := zClient.user_authentication_settings.Get()
	if err != nil {
		return err
	}
	d.SetId("exempted_urls")
	_ = d.Set("urls", res.URLs)
	return nil
}

func expandAuthSettingsUrls(d *schema.ResourceData) user_authentication_settings.ExemptedUrls {
	return user_authentication_settings.ExemptedUrls{
		URLs: SetToStringList(d, "urls"),
	}
}
func resourceAuthSettingsUrlsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	urls := expandAuthSettingsUrls(d)
	_, err := zClient.user_authentication_settings.Update(urls)
	if err != nil {
		return err
	}
	d.SetId("exempted_urls")
	return resourceAuthSettingsUrlsRead(d, m)
}

func resourceAuthSettingsUrlsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	urls := expandAuthSettingsUrls(d)

	_, err := zClient.user_authentication_settings.Update(urls)
	if err != nil {
		return err
	}
	return resourceAuthSettingsUrlsRead(d, m)
}

func resourceAuthSettingsUrlsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
