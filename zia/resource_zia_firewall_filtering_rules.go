package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/terraform-provider-zia/gozscaler/client"
	"github.com/zscaler/terraform-provider-zia/gozscaler/firewallpolicies/filteringrules"
)

func resourceFirewallFilteringRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallFilteringRulesCreate,
		Read:   resourceFirewallFilteringRulesRead,
		Update: resourceFirewallFilteringRulesUpdate,
		Delete: resourceFirewallFilteringRulesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("rule_id", id)
				} else {
					resp, err := zClient.filteringrules.GetByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("rule_id", resp.ID)
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
			"rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Firewall Filtering policy rule",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Rule order number of the Firewall Filtering policy rule",
			},
			"rank": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 7),
				Description:  "Admin rank of the Firewall Filtering policy rule",
			},
			"access_control": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_full_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The action the Firewall Filtering policy rule takes when packets match the rule",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"BLOCK_DROP",
					"BLOCK_RESET",
					"BLOCK_ICMP",
					"EVAL_NWAPP",
				}, false),
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether the Firewall Filtering policy rule is enabled or disabled",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the rule",
			},
			"last_modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"last_modified_by": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"src_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.",
			},
			"dest_addresses": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_ip_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_rule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, the default rule is applied",
			},
			"predefined": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, a predefined rule is applied",
			},
			"locations":             listIDsSchemaTypeCustom(8, "list of locations for which rule must be applied"),
			"location_groups":       listIDsSchemaTypeCustom(32, "list of locations groups"),
			"users":                 listIDsSchemaTypeCustom(4, "list of users for which rule must be applied"),
			"groups":                listIDsSchemaTypeCustom(8, "list of groups for which rule must be applied"),
			"departments":           listIDsSchemaType("list of departments for which rule must be applied"),
			"time_windows":          listIDsSchemaType("list of time interval during which rule must be enforced."),
			"labels":                listIDsSchemaType("list of Labels that are applicable to the rule."),
			"src_ip_groups":         listIDsSchemaType("list of src ip groups"),
			"dest_ip_groups":        listIDsSchemaType("list of dest ip groups"),
			"app_service_groups":    listIDsSchemaType("list of app service groups"),
			"app_services":          listIDsSchemaType("list of app services"),
			"nw_application_groups": listIDsSchemaType("list of nw application groups"),
			"nw_service_groups":     listIDsSchemaType("list of nw service groups"),
			"nw_services":           listIDsSchemaTypeCustom(1024, "list of nw services"),
			"dest_countries":        getCloudFirewallDstCountries(),
			"nw_applications":       getCloudFirewallNwApplications(),
		},
	}
}

func resourceFirewallFilteringRulesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandFirewallFilteringRules(d)
	log.Printf("[INFO] Creating zia firewall filtering rule\n%+v\n", req)

	resp, err := zClient.filteringrules.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia firewall filtering rule request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_id", resp.ID)

	return resourceFirewallFilteringRulesRead(d, m)
}

func resourceFirewallFilteringRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return fmt.Errorf("no zia firewall filtering rule id is set")
	}
	resp, err := zClient.filteringrules.Get(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing firewall filtering rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting firewall filtering rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("access_control", resp.AccessControl)
	_ = d.Set("enable_full_logging", resp.EnableFullLogging)
	_ = d.Set("action", resp.Action)
	_ = d.Set("state", resp.State)
	_ = d.Set("description", resp.Description)
	_ = d.Set("last_modified_time", resp.LastModifiedTime)
	_ = d.Set("src_ips", resp.SrcIps)
	_ = d.Set("dest_addresses", resp.DestAddresses)
	_ = d.Set("dest_ip_categories", resp.DestIpCategories)
	_ = d.Set("dest_countries", resp.DestCountries)
	_ = d.Set("nw_applications", resp.NwApplications)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("predefined", resp.Predefined)

	if err := d.Set("locations", flattenIDs(resp.Locations)); err != nil {
		return err
	}

	if err := d.Set("location_groups", flattenIDs(resp.LocationsGroups)); err != nil {
		return err
	}

	if err := d.Set("departments", flattenIDs(resp.Departments)); err != nil {
		return err
	}

	if err := d.Set("groups", flattenIDs(resp.Groups)); err != nil {
		return err
	}

	if err := d.Set("users", flattenIDs(resp.Users)); err != nil {
		return err
	}

	if err := d.Set("time_windows", flattenIDs(resp.TimeWindows)); err != nil {
		return err
	}

	if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
		return err
	}

	if err := d.Set("src_ip_groups", flattenIDs(resp.SrcIpGroups)); err != nil {
		return err
	}

	if err := d.Set("dest_ip_groups", flattenIDs(resp.DestIpGroups)); err != nil {
		return err
	}

	if err := d.Set("nw_services", flattenIDs(resp.NwServices)); err != nil {
		return err
	}

	if err := d.Set("nw_service_groups", flattenIDs(resp.NwServiceGroups)); err != nil {
		return err
	}

	if err := d.Set("nw_application_groups", flattenIDs(resp.NwApplicationGroups)); err != nil {
		return err
	}

	if err := d.Set("app_services", flattenIDs(resp.AppServices)); err != nil {
		return err
	}

	if err := d.Set("labels", flattenIDs(resp.Labels)); err != nil {
		return err
	}
	if err := d.Set("app_service_groups", flattenIDs(resp.AppServiceGroups)); err != nil {
		return err
	}

	return nil
}

func resourceFirewallFilteringRulesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall filteringrule ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating firewall filtering rule ID: %v\n", id)
	req := expandFirewallFilteringRules(d)

	if _, err := zClient.filteringrules.Update(id, &req); err != nil {
		return err
	}

	return resourceFirewallFilteringRulesRead(d, m)
}

func resourceFirewallFilteringRulesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall filtering rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting firewall filtering rule ID: %v\n", (d.Id()))

	if _, err := zClient.filteringrules.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] firewall filtering rule deleted")
	return nil
}

func expandFirewallFilteringRules(d *schema.ResourceData) filteringrules.FirewallFilteringRules {
	id, _ := getIntFromResourceData(d, "rule_id")
	result := filteringrules.FirewallFilteringRules{
		ID:                  id,
		Name:                d.Get("name").(string),
		Order:               d.Get("order").(int),
		Rank:                d.Get("rank").(int),
		Action:              d.Get("action").(string),
		State:               d.Get("state").(string),
		Description:         d.Get("description").(string),
		LastModifiedTime:    d.Get("last_modified_time").(int),
		SrcIps:              SetToStringList(d, "src_ips"),
		DestAddresses:       SetToStringList(d, "dest_addresses"),
		DestIpCategories:    SetToStringList(d, "dest_ip_categories"),
		DestCountries:       SetToStringList(d, "dest_countries"),
		NwApplications:      SetToStringList(d, "nw_applications"),
		DefaultRule:         d.Get("default_rule").(bool),
		Predefined:          d.Get("predefined").(bool),
		LastModifiedBy:      expandIDNameExtensions(d, "last_modified_by"),
		Locations:           expandIDNameExtensionsSet(d, "locations"),
		LocationsGroups:     expandIDNameExtensionsSet(d, "location_groups"),
		Departments:         expandIDNameExtensionsSet(d, "departments"),
		Groups:              expandIDNameExtensionsSet(d, "groups"),
		Users:               expandIDNameExtensionsSet(d, "users"),
		TimeWindows:         expandIDNameExtensionsSet(d, "time_windows"),
		SrcIpGroups:         expandIDNameExtensionsSet(d, "src_ip_groups"),
		DestIpGroups:        expandIDNameExtensionsSet(d, "dest_ip_groups"),
		NwServices:          expandIDNameExtensionsSet(d, "nw_services"),
		NwServiceGroups:     expandIDNameExtensionsSet(d, "nw_service_groups"),
		NwApplicationGroups: expandIDNameExtensionsSet(d, "nw_application_groups"),
		AppServices:         expandIDNameExtensionsSet(d, "app_services"),
		AppServiceGroups:    expandIDNameExtensionsSet(d, "app_service_groups"),
		Labels:              expandIDNameExtensionsSet(d, "labels"),
	}
	return result
}
