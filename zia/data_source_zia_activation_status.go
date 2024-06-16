package zia

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/activation"
)

func dataSourceActivationStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceActivationStatusRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceActivationStatusRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.activation

	resp, err := activation.GetActivationStatus(service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("activation")
		_ = d.Set("status", resp.Status)

	} else {
		return fmt.Errorf("couldn't find the activation status")
	}

	return nil
}
