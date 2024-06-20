package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/ipsourcegroups"
)

func dataSourceFWIPSourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWIPSourceGroupsRead,
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
			"ip_addresses": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFWIPSourceGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ipsourcegroups

	var resp *ipsourcegroups.IPSourceGroups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting ip source group id: %d\n", id)
		res, err := ipsourcegroups.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting ip source group : %s\n", name)
		res, err := ipsourcegroups.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("ip_addresses", resp.IPAddresses)

	} else {
		return fmt.Errorf("couldn't find any ip source group with name '%s' or id '%d'", name, id)
	}

	return nil
}
