package zia

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/ipdestinationgroups"
)

func resourceFWIPDestinationGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWIPDestinationGroupsCreate,
		Read:   resourceFWIPDestinationGroupsRead,
		Update: resourceFWIPDestinationGroupsUpdate,
		Delete: resourceFWIPDestinationGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("group_id", idInt)
				} else {
					resp, err := zClient.ipdestinationgroups.GetByName(id)
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifer for the destination IP group",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifer for the destination IP group",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Destination IP group name",
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Additional information about the destination IP group",
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Destination IP group type (i.e., the group can contain destination IP addresses or FQDNs)",
				ValidateFunc: validation.StringInSlice([]string{
					"DSTN_IP",
					"DSTN_FQDN",
					"DSTN_DOMAIN",
					"DSTN_OTHER",
				}, false),
			},
			"addresses": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Destination IP addresses within the group",
			},
			"ip_categories": getURLCategories(),
			"countries":     getDestinationCountries(),
		},
	}
}

func resourceFWIPDestinationGroupsCreate(d *schema.ResourceData, m interface{}) error {
	groupType := d.Get("type").(string)
	if groupType == "DSTN_OTHER" {
		ipCategories, ipCategoriesOk := d.GetOk("ip_categories")
		countries, countriesOk := d.GetOk("countries")
		if (!ipCategoriesOk || ipCategories.(*schema.Set).Len() == 0) && (!countriesOk || countries.(*schema.Set).Len() == 0) {
			return fmt.Errorf("when 'type' is set to 'DSTN_OTHER', either 'ip_categories' or 'countries' must be set")
		}
	}

	zClient := m.(*Client)
	req := expandIPDestinationGroups(d)
	log.Printf("[INFO] Creating zia ip destination groups\n%+v\n", req)

	resp, err := zClient.ipdestinationgroups.Create(&req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created zia ip destination groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)
	return resourceFWIPDestinationGroupsRead(d, m)
}

func resourceFWIPDestinationGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return fmt.Errorf("no ip destination groups id is set")
	}
	resp, err := zClient.ipdestinationgroups.Get(id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia ip destination groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	processedCountries := make([]string, len(resp.Countries))
	for i, country := range resp.Countries {
		processedCountries[i] = strings.TrimPrefix(country, "COUNTRY_")
	}

	log.Printf("[INFO] Getting zia ip destination groups:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("group_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("type", resp.Type)
	_ = d.Set("addresses", resp.Addresses)
	_ = d.Set("description", resp.Description)
	_ = d.Set("ip_categories", resp.IPCategories)
	_ = d.Set("countries", processedCountries)

	return nil
}

func resourceFWIPDestinationGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	groupType := d.Get("type").(string)

	if groupType == "DSTN_OTHER" || d.HasChange("countries") || d.HasChange("ip_categories") {
		ipCategories, ipCategoriesOk := d.GetOk("ip_categories")
		countries, countriesOk := d.GetOk("countries")

		if (!ipCategoriesOk || ipCategories.(*schema.Set).Len() == 0) && (!countriesOk || countries.(*schema.Set).Len() == 0) {
			return fmt.Errorf("when 'type' is set to 'DSTN_OTHER', either 'ip_categories' or 'countries' must be set")
		}
	}

	zClient := m.(*Client)
	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return fmt.Errorf("[ERROR] IP destination groups ID not set: %v", id)
	}

	log.Printf("[INFO] Updating ZIA IP destination groups ID: %v", id)
	req := expandIPDestinationGroups(d)

	_, _, err := zClient.ipdestinationgroups.Update(id, &req)
	if err != nil {
		return err
	}

	return resourceFWIPDestinationGroupsRead(d, m)
}

func resourceFWIPDestinationGroupsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] ip destination groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia ip destination groups ID: %v\n", (d.Id()))
	err := DetachRuleIDNameExtensions(
		zClient,
		id,
		"DestIpGroups",
		func(r *filteringrules.FirewallFilteringRules) []common.IDNameExtensions {
			return r.DestIpGroups
		},
		func(r *filteringrules.FirewallFilteringRules, ids []common.IDNameExtensions) {
			r.DestIpGroups = ids
		},
	)
	if err != nil {
		return err
	}
	if _, err := zClient.ipdestinationgroups.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] zia ip destination groups deleted")
	return nil
}

func expandIPDestinationGroups(d *schema.ResourceData) ipdestinationgroups.IPDestinationGroups {
	name := d.Get("name").(string)
	groupType := d.Get("type").(string)
	description := d.Get("description").(string)
	addresses := SetToStringList(d, "addresses")
	ipCategories := SetToStringList(d, "ip_categories")

	// Process countries to prepend "COUNTRY_"
	rawCountries := SetToStringList(d, "countries")
	countries := make([]string, len(rawCountries))
	for i, code := range rawCountries {
		countries[i] = "COUNTRY_" + code
	}

	return ipdestinationgroups.IPDestinationGroups{
		Name:         name,
		Type:         groupType,
		Description:  description,
		Addresses:    addresses,
		IPCategories: ipCategories,
		Countries:    countries,
	}
}
