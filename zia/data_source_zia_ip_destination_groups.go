package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/usermanagement"
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
			"idp_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGroupManagementRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *usermanagement.Group
	idObj, idSet := d.GetOk("id")
	id, idIsInt := idObj.(int)
	if idSet && idIsInt && id > 0 {
		log.Printf("[INFO] Getting data for user id: %d\n", id)
		res, err := zClient.usermanagement.GetGroups(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for user : %s\n", name)
		res, err := zClient.usermanagement.GetGroupByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("idp_id", resp.IdpID)
		_ = d.Set("comments", resp.Comments)

	} else {
		return fmt.Errorf("couldn't find any user with name '%s' or id '%d'", name, id)
	}

	return nil
}
