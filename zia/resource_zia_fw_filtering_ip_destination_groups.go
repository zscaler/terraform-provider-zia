package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipdestinationgroups"
)

func resourceFWIPDestinationGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFWIPDestinationGroupsCreate,
		ReadContext:   resourceFWIPDestinationGroupsRead,
		UpdateContext: resourceFWIPDestinationGroupsUpdate,
		DeleteContext: resourceFWIPDestinationGroupsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("group_id", idInt)
				} else {
					resp, err := ipdestinationgroups.GetByName(ctx, service, id)
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
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Additional information about the destination IP group",
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
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
			"countries":     getISOCountryCodes(),
		},
	}
}

func resourceFWIPDestinationGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupType := d.Get("type").(string)
	if groupType == "DSTN_OTHER" {
		ipCategories, ipCategoriesOk := d.GetOk("ip_categories")
		countries, countriesOk := d.GetOk("countries")
		if (!ipCategoriesOk || ipCategories.(*schema.Set).Len() == 0) && (!countriesOk || countries.(*schema.Set).Len() == 0) {
			return diag.FromErr(fmt.Errorf("when 'type' is set to 'DSTN_OTHER', either 'ip_categories' or 'countries' must be set"))
		}
	}

	zClient := meta.(*Client)
	service := zClient.Service

	req := expandIPDestinationGroups(d)
	log.Printf("[INFO] Creating zia ip destination groups\n%+v\n", req)

	resp, err := ipdestinationgroups.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created zia ip destination groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceFWIPDestinationGroupsRead(ctx, d, meta)
}

func resourceFWIPDestinationGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no ip destination groups id is set"))
	}
	resp, err := ipdestinationgroups.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia ip destination groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
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

func resourceFWIPDestinationGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupType := d.Get("type").(string)

	if groupType == "DSTN_OTHER" || d.HasChange("countries") || d.HasChange("ip_categories") {
		ipCategories, ipCategoriesOk := d.GetOk("ip_categories")
		countries, countriesOk := d.GetOk("countries")

		if (!ipCategoriesOk || ipCategories.(*schema.Set).Len() == 0) && (!countriesOk || countries.(*schema.Set).Len() == 0) {
			return diag.FromErr(fmt.Errorf("when 'type' is set to 'DSTN_OTHER', either 'ip_categories' or 'countries' must be set"))
		}
	}

	zClient := meta.(*Client)
	service := zClient.Service
	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("[ERROR] IP destination groups ID not set: %v", id))
	}

	log.Printf("[INFO] Updating ZIA IP destination groups ID: %v", id)
	req := expandIPDestinationGroups(d)

	_, _, err := ipdestinationgroups.Update(ctx, service, id, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceFWIPDestinationGroupsRead(ctx, d, meta)
}

func resourceFWIPDestinationGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] ip destination groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia ip destination groups ID: %v\n", (d.Id()))
	err := DetachRuleIDNameExtensions(
		ctx,
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
		return diag.FromErr(err)
	}
	if _, err := ipdestinationgroups.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia ip destination groups deleted")
	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

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
