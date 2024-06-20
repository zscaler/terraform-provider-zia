package zia

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/sandbox/sandbox_settings"
)

func dataSourceSandboxSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSandboxSettingsRead,
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

func dataSourceSandboxSettingsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.sandbox_settings

	resp, err := sandbox_settings.Get(service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("hash_id")
		_ = d.Set("file_hashes_to_be_blocked", resp.FileHashesToBeBlocked)

	} else {
		return fmt.Errorf("couldn't find any file hashes")
	}

	return nil
}
