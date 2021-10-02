package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
)

func resourceFirewallFilteringRules() *schema.Resource {
	return &schema.Resource{
		Create:   resourceFirewallFilteringRulesCreate,
		Read:     resourceFirewallFilteringRulesRead,
		Update:   resourceFirewallFilteringRulesUpdate,
		Delete:   resourceFirewallFilteringRulesDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"up_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dn_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"country": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tz": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ports": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_credentials": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CN",
								"IP",
								"UFQDN",
								"XAUTH",
							}, false),
						},
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pre_shared_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"auth_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssl_scan_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"zapp_ssl_scan_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"xff_forward_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"surrogate_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"idle_time_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"display_time_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MINUTE",
					"HOUR",
					"DAY",
				}, false),
			},
			"surrogate_ip_enforced_for_known_browsers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"surrogate_refresh_time_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"surrogate_refresh_time_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MINUTE",
					"HOUR",
					"DAY",
				}, false),
			},
			"ofw_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ips_control": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"caution_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aup_block_internet_until_accepted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aup_force_ssl_inspection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aup_timeout_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"profile": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"CORPORATE",
					"SERVER",
					"GUESTWIFI",
					"IOT",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
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

	return resourceFirewallFilteringRulesRead(d, m)
}

func resourceFirewallFilteringRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.filteringrules.Get(d.Id())

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing firewall filtering rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting firewall filtering rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
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

	if err := d.Set("locations", flattenLocations(resp.Locations)); err != nil {
		return err
	}

	if err := d.Set("location_groups", flattenLocationGroups(resp.LocationsGroups)); err != nil {
		return err
	}

	if err := d.Set("departments", flattenDepartments(resp.Departments)); err != nil {
		return err
	}

	if err := d.Set("groups", flattenGroups(resp.Groups)); err != nil {
		return err
	}

	if err := d.Set("users", flattenUsers(resp.Users)); err != nil {
		return err
	}

	if err := d.Set("time_windows", flattenTimeWindows(resp.TimeWindows)); err != nil {
		return err
	}

	if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
		return err
	}

	if err := d.Set("src_ip_groups", flattenSrcIPGroups(resp.SrcIpGroups)); err != nil {
		return err
	}

	if err := d.Set("dest_ip_groups", flattenDestIPGroups(resp.DestIpGroups)); err != nil {
		return err
	}

	if err := d.Set("nw_services", flattenNWServices(resp.NwServices)); err != nil {
		return err
	}

	if err := d.Set("nw_service_groups", flattenNWServiceGroups(resp.NwServiceGroups)); err != nil {
		return err
	}

	if err := d.Set("nw_application_groups", flattenNWApplicationGroups(resp.NwApplicationGroups)); err != nil {
		return err
	}

	if err := d.Set("app_services", flattenAppServices(resp.AppServices)); err != nil {
		return err
	}

	if err := d.Set("app_services", flattenAppServiceGroups(resp.AppServiceGroups)); err != nil {
		return err
	}

	if err := d.Set("labels", flattenLabels(resp.Labels)); err != nil {
		return err
	}

	return nil
}

func resourceFirewallFilteringRulesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating firewall filtering rule ID: %v\n", id)
	req := expandFirewallFilteringRules(d)

	if _, err := zClient.filteringrules.Update(id, &req); err != nil {
		return err
	}

	return resourceFirewallFilteringRulesRead(d, m)
}

func resourceFirewallFilteringRulesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	// Need to pass the ID (int) of the resource for deletion
	log.Printf("[INFO] Deleting firewall filtering rule ID: %v\n", (d.Id()))

	if _, err := zClient.filteringrules.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] firewall filtering rule deleted")
	return nil
}
