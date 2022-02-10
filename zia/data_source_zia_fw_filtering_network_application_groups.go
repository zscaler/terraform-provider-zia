package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkapplications"
)

func dataSourceFWNetworkApplicationGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWNetworkApplicationGroupsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"network_applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFWNetworkApplicationGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *networkapplications.NetworkApplicationGroups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting network application group id: %d\n", id)
		res, err := zClient.networkapplications.GetNetworkApplicationGroups(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting network application group : %s\n", name)
		res, err := zClient.networkapplications.GetNetworkApplicationGroupsByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("network_applications", resp.NetworkApplications)
		_ = d.Set("description", resp.Description)

	} else {
		return fmt.Errorf("couldn't find any network application group with name '%s' or id '%d'", name, id)
	}

	return nil
}
