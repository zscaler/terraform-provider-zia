package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/groups"
)

func dataSourceGroupManagement() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupManagementRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
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

func dataSourceGroupManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *groups.Groups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for user id: %d\n", id)
		res, err := groups.GetGroups(ctx, service, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting group by ID %d: %s", id, err))
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for user : %s\n", name)
		res, err := groups.GetGroupByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting group by name %s: %s", name, err))
		}
		resp = res
		log.Printf("[DEBUG] Group received: %+v", resp) // Log the received group for debugging
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		err := d.Set("name", resp.Name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error setting name: %s", err))
		}
		err = d.Set("idp_id", resp.IdpID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error setting idp_id: %s", err))
		}
		err = d.Set("comments", resp.Comments)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error setting comments: %s", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any user with name '%s' or id '%d'", name, id))
	}

	return nil
}
