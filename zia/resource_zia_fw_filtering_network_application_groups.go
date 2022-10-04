package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/networkapplications"
)

func resourceFWNetworkApplicationGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWNetworkApplicationGroupsCreate,
		Read:   resourceFWNetworkApplicationGroupsRead,
		Update: resourceFWNetworkApplicationGroupsUpdate,
		Delete: resourceFWNetworkApplicationGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("app_id", idInt)
				} else {
					resp, err := zClient.networkapplications.GetNetworkApplicationGroupsByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("app_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_applications": {
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

func resourceFWNetworkApplicationGroupsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandNetworkApplicationGroups(d)
	log.Printf("[INFO] Creating network application groups\n%+v\n", req)

	resp, err := zClient.networkapplications.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia network application groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("app_id", resp.ID)
	return resourceFWNetworkApplicationGroupsRead(d, m)
}

func resourceFWNetworkApplicationGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "app_id")
	if !ok {
		return fmt.Errorf("no network application groups id is set")
	}
	resp, err := zClient.networkapplications.GetNetworkApplicationGroups(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia network application groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting network application groups :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("app_id", resp.ID)
	_ = d.Set("network_applications", resp.NetworkApplications)
	_ = d.Set("description", resp.Description)

	return nil
}

func resourceFWNetworkApplicationGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "app_id")
	if !ok {
		log.Printf("[ERROR] network application groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating network application groups ID: %v\n", id)
	req := expandNetworkApplicationGroups(d)

	if _, _, err := zClient.networkapplications.Update(id, &req); err != nil {
		return err
	}

	return resourceFWNetworkApplicationGroupsRead(d, m)
}

func resourceFWNetworkApplicationGroupsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "app_id")
	if !ok {
		log.Printf("[ERROR] network application groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting network application groups ID: %v\n", (d.Id()))

	if _, err := zClient.networkapplications.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] network application groups deleted")
	return nil
}

func expandNetworkApplicationGroups(d *schema.ResourceData) networkapplications.NetworkApplicationGroups {
	id, _ := getIntFromResourceData(d, "app_id")
	result := networkapplications.NetworkApplicationGroups{
		ID:                  id,
		Name:                d.Get("name").(string),
		NetworkApplications: SetToStringList(d, "network_applications"),
		Description:         d.Get("description").(string),
	}

	return result
}
