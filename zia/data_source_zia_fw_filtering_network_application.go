package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkapplications"
)

func dataSourceFWNetworkApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWNetworkApplicationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"deprecated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"parent_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFWNetworkApplicationRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.networkapplications

	id, ok := getStringFromResourceData(d, "id")
	if !ok {
		return fmt.Errorf("network application id is required '%s'", id)
	}

	log.Printf("[INFO] Getting network application group id: %s\n", id)
	resp, err := networkapplications.GetNetworkApplication(service, id, d.Get("locale").(string))
	if err != nil {
		return err
	}

	d.SetId(resp.ID)
	_ = d.Set("id", resp.ID)
	_ = d.Set("deprecated", resp.Deprecated)
	_ = d.Set("parent_category", resp.ParentCategory)
	_ = d.Set("description", resp.Description)

	return nil
}
