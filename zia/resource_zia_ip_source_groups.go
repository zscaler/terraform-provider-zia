package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipsourcegroups"
)

func resourceIPSourceGroups() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIPSourceGroupsCreate,
		Read:     resourceIPSourceGroupsRead,
		Update:   resourceIPSourceGroupsUpdate,
		Delete:   resourceIPSourceGroupsDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIPSourceGroupsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandIPSourceGroups(d)
	log.Printf("[INFO] Creating zia ip destisourcenation groups\n%+v\n", req)

	resp, err := zClient.ipsourcegroups.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia ip source groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))

	return resourceIPSourceGroupsRead(d, m)
}

func resourceIPSourceGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "id")
	if !ok {
		return fmt.Errorf("no ip source groups id is set")
	}
	resp, err := zClient.ipsourcegroups.Get(id)

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing zia ip source groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting zia ip source groups:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("ip_addresses", resp.IPAddresses)
	_ = d.Set("description", resp.Description)

	return nil
}

func resourceIPSourceGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating zia ip source groupsID: %v\n", id)
	req := expandIPSourceGroups(d)

	if _, err := zClient.ipsourcegroups.Update(id, &req); err != nil {
		return err
	}

	return resourceIPSourceGroupsRead(d, m)
}

func resourceIPSourceGroupsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	// Need to pass the ID (int) of the resource for deletion
	log.Printf("[INFO] Deleting zia ip source groups ID: %v\n", (d.Id()))

	if err := zClient.ipsourcegroups.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] zia ip source groups deleted")
	return nil
}

func expandIPSourceGroups(d *schema.ResourceData) ipsourcegroups.IPSourceGroups {
	return ipsourcegroups.IPSourceGroups{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		IPAddresses: ListToStringSlice(d.Get("ip_addresses").([]interface{})),
	}
}
