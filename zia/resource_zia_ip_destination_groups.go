package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipdestinationgroups"
)

func resourceIPDestinationGroups() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIPDestinationGroupsCreate,
		Read:     resourceIPDestinationGroupsRead,
		Update:   resourceIPDestinationGroupsUpdate,
		Delete:   resourceIPDestinationGroupsDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_categories": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"countries": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceIPDestinationGroupsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandIPDestinationGroups(d)
	log.Printf("[INFO] Creating zia ip destination groups\n%+v\n", req)

	resp, err := zClient.ipdestinationgroups.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia ip destination groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))

	return resourceIPDestinationGroupsRead(d, m)
}

func resourceIPDestinationGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.ipdestinationgroups.Get(d.Id())

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing zia ip destination groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting zia ip destination groups:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("type", resp.Type)
	_ = d.Set("addresses", resp.Addresses)
	_ = d.Set("description", resp.Description)
	_ = d.Set("ip_categories", resp.IPCategories)
	_ = d.Set("countries", resp.Countries)

	return nil
}

func resourceIPDestinationGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating zia ip destination groupsID: %v\n", id)
	req := expandIPDestinationGroups(d)

	if _, err := zClient.ipdestinationgroups.Update(id, &req); err != nil {
		return err
	}

	return resourceIPDestinationGroupsRead(d, m)
}

func resourceIPDestinationGroupsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	// Need to pass the ID (int) of the resource for deletion
	log.Printf("[INFO] Deleting zia ip destination groups ID: %v\n", (d.Id()))

	if err := zClient.ipdestinationgroups.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] zia ip destination groups deleted")
	return nil
}

func expandIPDestinationGroups(d *schema.ResourceData) ipdestinationgroups.IPDestinationGroups {
	return ipdestinationgroups.IPDestinationGroups{
		Name:         d.Get("name").(string),
		Type:         d.Get("type").(string),
		Description:  d.Get("description").(string),
		Addresses:    d.Get("addresses").([]string),
		IPCategories: d.Get("ip_categories").([]string),
		Countries:    d.Get("countries").([]string),
	}
}
