package zia

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/user_authentication_settings"
)

func dataSourceAuthSettingsUrls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthSettingsUrlsRead,
		Schema: map[string]*schema.Schema{
			"urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAuthSettingsUrlsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	res, err := user_authentication_settings.Get(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("exempted_urls")
	_ = d.Set("urls", res.URLs)
	return nil
}
