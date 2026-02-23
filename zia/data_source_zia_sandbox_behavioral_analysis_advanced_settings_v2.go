package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_settings"
)

func dataSourceSandboxSettingsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSandboxSettingsV2Read,
		Schema: map[string]*schema.Schema{
			"md5_hash_value_list": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A custom list of MD5 hash values with metadata for sandbox blocking.",
				Set:         md5HashValueListHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MD5 hash identifier for the entry.",
						},
						"url_comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A comment describing the MD5 hash entry.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the entry (e.g. CUSTOM_FILEHASH_ALLOW, CUSTOM_FILEHASH_DENY, MALWARE).",
						},
					},
				},
			},
		},
	}
}

func dataSourceSandboxSettingsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := sandbox_settings.Getv2(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error fetching sandbox behavioral analysis advanced settings: %s", err))
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("no sandbox behavioral analysis settings found"))
	}

	d.SetId("sandbox_settings")

	if len(resp.Md5HashValueList) > 0 {
		if err := d.Set("md5_hash_value_list", flattenMd5HashValueList(resp.Md5HashValueList)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting md5_hash_value_list: %s", err))
		}
	} else {
		if err := d.Set("md5_hash_value_list", []interface{}{}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting md5_hash_value_list: %s", err))
		}
	}

	return nil
}
