package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipdestinationgroups"
)

func dataSourceIPDestinationGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIPDestinationGroupsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extensions": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIPDestinationGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *ipdestinationgroups.IPDestinationGroupsLite
	idObj, idSet := d.GetOk("id")
	id, idIsInt := idObj.(int)
	if idSet && idIsInt && id > 0 {
		log.Printf("[INFO] Getting data for ip destination groups id: %d\n", id)
		res, err := zClient.ipdestinationgroups.GetIPDestinationGroupsLite(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for ip destination groups : %s\n", name)
		res, err := zClient.ipdestinationgroups.GetIPDestinationGroupsLiteByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("extensions", resp.Extensions)

	} else {
		return fmt.Errorf("couldn't find any ip destination groups with name '%s' or id '%d'", name, id)
	}

	return nil
}
