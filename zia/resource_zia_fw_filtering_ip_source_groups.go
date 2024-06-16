package zia

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/ipsourcegroups"
)

func resourceFWIPSourceGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWIPSourceGroupsCreate,
		Read:   resourceFWIPSourceGroupsRead,
		Update: resourceFWIPSourceGroupsUpdate,
		Delete: resourceFWIPSourceGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)
				service := zClient.ipsourcegroups

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("group_id", idInt)
				} else {
					resp, err := ipsourcegroups.GetByName(service, id)
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
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
			"ip_addresses": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceFWIPSourceGroupsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ipsourcegroups

	req := expandFWIPSourceGroups(d)
	log.Printf("[INFO] Creating zia ip source groups\n%+v\n", req)

	resp, err := ipsourcegroups.Create(service, &req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia ip source groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceFWIPSourceGroupsRead(d, m)
}

func resourceFWIPSourceGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ipsourcegroups

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return fmt.Errorf("no ip source groups id is set")
	}
	resp, err := ipsourcegroups.Get(service, id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia ip source groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting zia ip source groups:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("group_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("ip_addresses", resp.IPAddresses)
	_ = d.Set("description", resp.Description)

	return nil
}

func resourceFWIPSourceGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ipsourcegroups

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] ip source groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia ip source groups ID: %v\n", id)
	req := expandFWIPSourceGroups(d)
	if _, err := ipsourcegroups.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := ipsourcegroups.Update(service, id, &req); err != nil {
		return err
	}
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceFWIPSourceGroupsRead(d, m)
}

func resourceFWIPSourceGroupsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.ipsourcegroups

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] ip source groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia ip source groups ID: %v\n", (d.Id()))
	err := DetachRuleIDNameExtensions(
		zClient,
		id,
		"FWIPSourceGroups",
		func(r *filteringrules.FirewallFilteringRules) []common.IDNameExtensions {
			return r.SrcIpGroups
		},
		func(r *filteringrules.FirewallFilteringRules, ids []common.IDNameExtensions) {
			r.SrcIpGroups = ids
		},
	)
	if err != nil {
		return err
	}
	if _, err := ipsourcegroups.Delete(service, id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] zia ip source groups deleted")
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandFWIPSourceGroups(d *schema.ResourceData) ipsourcegroups.IPSourceGroups {
	return ipsourcegroups.IPSourceGroups{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		IPAddresses: SetToStringList(d, "ip_addresses"),
	}
}
