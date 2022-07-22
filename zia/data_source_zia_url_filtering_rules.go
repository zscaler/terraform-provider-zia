package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zia/services/urlfilteringpolicies"
)

func dataSourceURLFilteringRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceURLFilteringRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"url_categories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_windows": {
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
			"rank": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"request_methods": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"end_user_notification_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"override_users": {
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
			"override_groups": {
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
			"block_override": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"time_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
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
			"last_modified_by": {
				Type:     schema.TypeSet,
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
			"enforce_time_validity": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"user_agent_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cbi_profile_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ciparule": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceURLFilteringRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *urlfilteringpolicies.URLFilteringRule
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data url filtering policy id: %d\n", id)
		res, err := zClient.urlfilteringpolicies.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting url filtering policy : %s\n", name)
		res, err := zClient.urlfilteringpolicies.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("order", resp.Order)
		_ = d.Set("protocols", resp.Protocols)
		_ = d.Set("url_categories", resp.URLCategories)
		_ = d.Set("state", resp.State)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("request_methods", resp.RequestMethods)
		_ = d.Set("end_user_notification_url", resp.EndUserNotificationURL)
		_ = d.Set("block_override", resp.BlockOverride)
		_ = d.Set("time_quota", resp.TimeQuota)
		_ = d.Set("size_quota", resp.SizeQuota)
		_ = d.Set("description", resp.Description)
		_ = d.Set("validity_start_time", resp.ValidityStartTime)
		_ = d.Set("validity_end_time", resp.ValidityEndTime)
		_ = d.Set("validity_time_zone_id", resp.ValidityTimeZoneID)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)
		_ = d.Set("enforce_time_validity", resp.EnforceTimeValidity)
		_ = d.Set("action", resp.Action)
		_ = d.Set("ciparule", resp.Ciparule)

		if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
			return err
		}

		if err := d.Set("groups", flattenIDNameExtensions(resp.Groups)); err != nil {
			return err
		}

		if err := d.Set("departments", flattenIDNameExtensions(resp.Departments)); err != nil {
			return err
		}

		if err := d.Set("users", flattenIDNameExtensions(resp.Users)); err != nil {
			return err
		}

		if err := d.Set("time_windows", flattenIDNameExtensions(resp.TimeWindows)); err != nil {
			return err
		}

		if err := d.Set("override_users", flattenIDNameExtensions(resp.OverrideUsers)); err != nil {
			return err
		}

		if err := d.Set("override_groups", flattenIDNameExtensions(resp.OverrideGroups)); err != nil {
			return err
		}

		if err := d.Set("location_groups", flattenIDNameExtensions(resp.LocationGroups)); err != nil {
			return err
		}

		if err := d.Set("device_groups", flattenIDNameExtensions(resp.DeviceGroups)); err != nil {
			return err
		}

		if err := d.Set("devices", flattenIDNameExtensions(resp.Devices)); err != nil {
			return err
		}

		if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
			return err
		}

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return err
		}
		if err := d.Set("device_groups", flattenIDNameExtensions(resp.DeviceGroups)); err != nil {
			return err
		}

		if err := d.Set("devices", flattenIDNameExtensions(resp.Devices)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any url filtering rule with name '%s' or id '%d'", name, id)
	}

	return nil
}
