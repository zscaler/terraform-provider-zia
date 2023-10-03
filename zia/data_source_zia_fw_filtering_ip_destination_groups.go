package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/ipdestinationgroups"
)

func dataSourceFWIPDestinationGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWIPDestinationGroupsRead,
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_categories": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"countries": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceFWIPDestinationGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *ipdestinationgroups.IPDestinationGroups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for ip destination groups id: %d\n", id)
		res, err := zClient.ipdestinationgroups.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for ip destination groups : %s\n", name)
		res, err := zClient.ipdestinationgroups.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("type", resp.Type)
		_ = d.Set("addresses", resp.Addresses)
		_ = d.Set("description", resp.Description)
		_ = d.Set("ip_categories", resp.IPCategories)
		_ = d.Set("countries", resp.Countries)

	} else {
		return fmt.Errorf("couldn't find any ip destination groups with name '%s' or id '%d'", name, id)
	}

	return nil
}
