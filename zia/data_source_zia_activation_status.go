package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/activation"
)

func dataSourceActivationStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceActivationStatusRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceActivationStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := activation.GetActivationStatus(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("activation")
		_ = d.Set("status", resp.Status)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find the activation status"))
	}

	return nil
}
