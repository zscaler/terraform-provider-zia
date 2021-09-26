package zia

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceActivationStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceActivationStatusRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceActivationStatusRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.activation.GetActivationStatus()
	if err != nil {
		return nil
	}

	_ = d.Set("status", resp.Status)

	return nil
}
