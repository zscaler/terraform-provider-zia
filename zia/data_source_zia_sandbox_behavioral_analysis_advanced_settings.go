package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_settings"
)

func dataSourceSandboxSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSandboxSettingsRead,
		Schema: map[string]*schema.Schema{
			"file_hashes_to_be_blocked": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSandboxSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := sandbox_settings.Get(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("hash_id")
		_ = d.Set("file_hashes_to_be_blocked", resp.FileHashesToBeBlocked)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any file hashes"))
	}

	return nil
}
