package zia

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	resp, err := zClient.activation.GetActivationStatus()
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
