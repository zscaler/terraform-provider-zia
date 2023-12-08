package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkapplicationgroups"
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
					resp, err := zClient.networkapplicationgroups.GetNetworkApplicationGroupsByName(id)
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
			"network_applications": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceFWNetworkApplicationGroupsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandNetworkApplicationGroups(d)
	log.Printf("[INFO] Creating network application groups\n%+v\n", req)

	resp, err := zClient.networkapplicationgroups.Create(&req)
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
	resp, err := zClient.networkapplicationgroups.GetNetworkApplicationGroups(id)
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
	if _, err := zClient.networkapplicationgroups.GetNetworkApplicationGroups(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := zClient.networkapplicationgroups.Update(id, &req); err != nil {
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
	err := DetachRuleIDNameExtensions(
		zClient,
		id,
		"NwApplicationGroups",
		func(r *filteringrules.FirewallFilteringRules) []common.IDNameExtensions {
			return r.NwApplicationGroups
		},
		func(r *filteringrules.FirewallFilteringRules, ids []common.IDNameExtensions) {
			r.NwApplicationGroups = ids
		},
	)
	if err != nil {
		return err
	}
	if _, err := zClient.networkapplicationgroups.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] network application groups deleted")
	return nil
}

func expandNetworkApplicationGroups(d *schema.ResourceData) networkapplicationgroups.NetworkApplicationGroups {
	id, _ := getIntFromResourceData(d, "app_id")
	result := networkapplicationgroups.NetworkApplicationGroups{
		ID:                  id,
		Name:                d.Get("name").(string),
		NetworkApplications: SetToStringList(d, "network_applications"),
		Description:         d.Get("description").(string),
	}

	return result
}
