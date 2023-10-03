package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/usermanagement/groups"
)

func dataSourceGroupManagement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupManagementRead,
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

	var resp *groups.Groups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for user id: %d\n", id)
		res, err := zClient.groups.GetGroups(id)
		if err != nil {
			return fmt.Errorf("error getting group by ID %d: %s", id, err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for user : %s\n", name)
		res, err := zClient.groups.GetGroupByName(name)
		if err != nil {
			return fmt.Errorf("error getting group by name %s: %s", name, err)
		}
		resp = res
		log.Printf("[DEBUG] Group received: %+v", resp) // Log the received group for debugging
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		err := d.Set("name", resp.Name)
		if err != nil {
			return fmt.Errorf("error setting name: %s", err)
		}
		err = d.Set("idp_id", resp.IdpID)
		if err != nil {
			return fmt.Errorf("error setting idp_id: %s", err)
		}
		err = d.Set("comments", resp.Comments)
		if err != nil {
			return fmt.Errorf("error setting comments: %s", err)
		}
	} else {
		return fmt.Errorf("couldn't find any user with name '%s' or id '%d'", name, id)
	}

	return nil
}
