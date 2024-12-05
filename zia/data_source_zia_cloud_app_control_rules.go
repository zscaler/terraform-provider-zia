package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudappcontrol"
)

func dataSourceCloudAppControlRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudAppControlRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier for the device.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier for the device.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier for the device.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"actions": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"time_quota": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"size_quota": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"cascading_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"access_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"applications": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"number_of_applications": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"eun_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"eun_template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"browser_eun_template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"predefined": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"device_trust_levels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"validity_start_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"validity_end_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"validity_time_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enforce_time_validity": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"user_agent_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cbi_profile": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The universally unique identifier (UUID) for the browser isolation profile",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the browser isolation profile",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The browser isolation profile URL",
						},
						"default_profile": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The browser isolation profile URL",
						},
						"sandbox_mode": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The browser isolation profile URL",
						},
					},
				},
			},
			"devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"device_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"location_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"departments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudAppControlRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	ruleType, ok := d.Get("type").(string)
	if !ok || ruleType == "" {
		return diag.FromErr(fmt.Errorf("type must be specified"))
	}

	var resp *cloudappcontrol.WebApplicationRules
	id, idOk := getIntFromResourceData(d, "id")
	if idOk {
		log.Printf("[INFO] Getting data for cloud app control rule id: %d and type: %s\n", id, ruleType)
		res, err := cloudappcontrol.GetByRuleID(ctx, service, ruleType, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, nameOk := d.Get("name").(string)
	if resp == nil && nameOk && name != "" {
		log.Printf("[INFO] Getting data for cloud app control rule: %s and type: %s\n", name, ruleType)
		res, err := cloudappcontrol.GetByRuleType(ctx, service, ruleType)
		if err != nil {
			return diag.FromErr(err)
		}

		// Look for the rule with the specified name
		for _, rule := range res {
			if rule.Name == name {
				resp = &rule
				break
			}
		}
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("order", resp.Order)
		_ = d.Set("actions", resp.Actions)
		_ = d.Set("state", resp.State)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("type", resp.Type)
		_ = d.Set("time_quota", resp.TimeQuota)
		_ = d.Set("size_quota", resp.SizeQuota)
		_ = d.Set("cascading_enabled", resp.CascadingEnabled)
		_ = d.Set("access_control", resp.AccessControl)
		_ = d.Set("applications", resp.Applications)
		_ = d.Set("number_of_applications", resp.NumberOfApplications)
		_ = d.Set("eun_enabled", resp.EunEnabled)
		_ = d.Set("eun_template_id", resp.EunTemplateID)
		_ = d.Set("browser_eun_template_id", resp.BrowserEunTemplateID)
		_ = d.Set("predefined", resp.Predefined)
		_ = d.Set("validity_start_time", resp.ValidityStartTime)
		_ = d.Set("validity_end_time", resp.ValidityEndTime)
		_ = d.Set("validity_time_zone_id", resp.ValidityTimeZoneID)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)
		_ = d.Set("enforce_time_validity", resp.EnforceTimeValidity)
		_ = d.Set("user_agent_types", resp.UserAgentTypes)
		if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("groups", flattenIDNameExtensions(resp.Groups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("departments", flattenIDNameExtensions(resp.Departments)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("users", flattenIDNameExtensions(resp.Users)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("location_groups", flattenIDNameExtensions(resp.LocationGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("device_groups", flattenIDNameExtensions(resp.DeviceGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("devices", flattenIDNameExtensions(resp.Devices)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}

		if resp.CBIProfile.ID != "" {
			if err := d.Set("cbi_profile", flattenCloudAppControlCBIProfile(&resp.CBIProfile)); err != nil {
				return diag.FromErr(err)
			}
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any cloud application rule with name '%s' or id '%d'", name, id))
	}

	return nil
}

func flattenCloudAppControlCBIProfile(cbiProfile *cloudappcontrol.CBIProfile) map[string]interface{} {
	if cbiProfile == nil {
		return nil
	}

	return map[string]interface{}{
		"id":              cbiProfile.ID,
		"name":            cbiProfile.Name,
		"url":             cbiProfile.URL,
		"profile_seq":     cbiProfile.ProfileSeq,
		"default_profile": cbiProfile.DefaultProfile,
		"sandbox_mode":    cbiProfile.SandboxMode,
	}
}
