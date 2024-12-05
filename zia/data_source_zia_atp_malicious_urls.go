package zia

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/advancedthreatsettings"
)

func dataSourceATPMaliciousUrls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceATPMaliciousUrlsRead,
		Schema: map[string]*schema.Schema{
			"malicious_urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceATPMaliciousUrlsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	res, err := advancedthreatsettings.GetMaliciousURLs(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("url_list")
	_ = d.Set("malicious_urls", res.MaliciousUrls)
	return nil
}
