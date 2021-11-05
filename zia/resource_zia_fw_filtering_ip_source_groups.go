package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipsourcegroups"
)

func resourceFWIPSourceGroups() *schema.Resource {
	return &schema.Resource{
		Create:   resourceFWIPSourceGroupsCreate,
		Read:     resourceFWIPSourceGroupsRead,
		Update:   resourceFWIPSourceGroupsUpdate,
		Delete:   resourceFWIPSourceGroupsDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_source_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_addresses": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
		},
	}
}

func resourceFWIPSourceGroupsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandFWIPSourceGroups(d)
	log.Printf("[INFO] Creating zia ip source groups\n%+v\n", req)

	resp, err := zClient.ipsourcegroups.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia ip source groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("ip_source_group_id", resp.ID)
	return resourceFWIPSourceGroupsRead(d, m)
}

func resourceFWIPSourceGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "ip_source_group_id")
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
	_ = d.Set("ip_source_group_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("ip_addresses", resp.IPAddresses)
	_ = d.Set("description", resp.Description)

	return nil
}

func resourceFWIPSourceGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "ip_source_group_id")
	if !ok {
		log.Printf("[ERROR] ip source groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia ip source groups ID: %v\n", id)
	req := expandFWIPSourceGroups(d)

	if _, err := zClient.ipsourcegroups.Update(id, &req); err != nil {
		return err
	}

	return resourceFWIPSourceGroupsRead(d, m)
}

func resourceFWIPSourceGroupsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "ip_source_group_id")
	if !ok {
		log.Printf("[ERROR] ip source groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia ip source groups ID: %v\n", (d.Id()))

	if _, err := zClient.ipsourcegroups.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] zia ip source groups deleted")
	return nil
}

func expandFWIPSourceGroups(d *schema.ResourceData) ipsourcegroups.IPSourceGroups {
	return ipsourcegroups.IPSourceGroups{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		IPAddresses: SetToStringList(d, "ip_addresses"),
	}
}
