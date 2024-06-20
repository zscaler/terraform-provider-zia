package zia

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/security_policy_settings"
)

func dataSourceSecurityPolicySettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecurityPolicySettingsRead,
		Schema: map[string]*schema.Schema{
			"whitelist_urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"blacklist_urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSecurityPolicySettingsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.security_policy_settings

	resp, err := security_policy_settings.GetListUrls(service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("url_id")
		_ = d.Set("whitelist_urls", resp.White)
		_ = d.Set("blacklist_urls", resp.Black)

	} else {
		return fmt.Errorf("couldn't find the activation status")
	}

	return nil
}
