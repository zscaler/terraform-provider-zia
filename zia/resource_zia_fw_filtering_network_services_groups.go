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
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkservicegroups"
)

func resourceFWNetworkServiceGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWNetworkServiceGroupsCreate,
		Read:   resourceFWNetworkServiceGroupsRead,
		Update: resourceFWNetworkServiceGroupsUpdate,
		Delete: resourceFWNetworkServiceGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("group_id", idInt)
				} else {
					resp, err := zClient.networkservicegroups.GetNetworkServiceGroupsByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("group_id", resp.ID)
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
			"group_id": {
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
			"services": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "list of services IDs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		},
	}
}

func resourceFWNetworkServiceGroupsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandNetworkServiceGroups(d)
	log.Printf("[INFO] Creating network service groups\n%+v\n", req)

	resp, err := zClient.networkservicegroups.CreateNetworkServiceGroups(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia network service groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)

	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return resourceFWNetworkServiceGroupsRead(d, m)
}

func resourceFWNetworkServiceGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return fmt.Errorf("no network service groups id is set")
	}
	resp, err := zClient.networkservicegroups.GetNetworkServiceGroups(id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia network service groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting network service groups :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("group_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)

	if err := d.Set("services", flattenServicesSimple(resp.Services)); err != nil {
		return err
	}

	return nil
}

func flattenServicesSimple(list []networkservicegroups.Services) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]int, len(list))
	for i, item := range list {
		ids[i] = item.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func resourceFWNetworkServiceGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] network service groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating network service groups ID: %v\n", id)
	req := expandNetworkServiceGroups(d)
	if _, err := zClient.networkservicegroups.GetNetworkServiceGroups(req.ID); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := zClient.networkservicegroups.UpdateNetworkServiceGroups(id, &req); err != nil {
		return err
	}
	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return resourceFWNetworkServiceGroupsRead(d, m)
}

func resourceFWNetworkServiceGroupsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] network service groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting network service groups ID: %v\n", (d.Id()))
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
	if _, err := zClient.networkservicegroups.DeleteNetworkServiceGroups(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] network service groups deleted")

	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return nil
}

func expandNetworkServiceGroups(d *schema.ResourceData) networkservicegroups.NetworkServiceGroups {
	id, _ := getIntFromResourceData(d, "group_id")
	result := networkservicegroups.NetworkServiceGroups{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Services:    expandServicesSet(d),
	}

	return result
}

func expandServicesSet(d *schema.ResourceData) []networkservicegroups.Services {
	setInterface, ok := d.GetOk("services")
	if ok {
		set := setInterface.(*schema.Set)
		var result []networkservicegroups.Services
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil {
				idSet, ok := itemMap["id"].(*schema.Set)
				if ok {
					for _, id := range idSet.List() {
						result = append(result, networkservicegroups.Services{
							ID: id.(int),
						})
					}
				}
			}
		}
		return result
	}
	return []networkservicegroups.Services{}
}
