package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkservices"
)

func dataSourceNetworkServiceGroupsLite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkServiceGroupsLiteRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkServiceGroupsLiteRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *networkservices.NetworkServiceGroups
	idObj, idSet := d.GetOk("id")
	id, idIsInt := idObj.(int)
	if idSet && idIsInt && id > 0 {
		log.Printf("[INFO] Getting network service group id: %d\n", id)
		res, err := zClient.networkservices.GetNetworkServiceGroupsLite(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting network service group : %s\n", name)
		res, err := zClient.networkservices.GetNetworkServiceGroupsLiteByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)

	} else {
		return fmt.Errorf("couldn't find any network service group with name '%s' or id '%d'", name, id)
	}

	return nil
}
