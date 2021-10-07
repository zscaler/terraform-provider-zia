package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlfilteringpolicies"
)

func dataSourceURLFilteringRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceURLFilteringRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"enforce_time_validity": {
				Type:     schema.TypeBool,
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

		if err := d.Set("locations", flattenURLFilteringLocation(resp.Locations)); err != nil {
			return err
		}

		if err := d.Set("groups", flattenURLFilteringGroups(resp.Groups)); err != nil {
			return err
		}

		if err := d.Set("departments", flattenURLFilteringDepartments(resp.Departments)); err != nil {
			return err
		}

		if err := d.Set("users", flattenURLFilteringUsers(resp.Users)); err != nil {
			return err
		}

		if err := d.Set("time_windows", flattenURLFilteringTimeWindows(resp.TimeWindows)); err != nil {
			return err
		}

		if err := d.Set("override_users", flattenURLFilteringOverrideUsers(resp.OverrideUsers)); err != nil {
			return err
		}

		if err := d.Set("override_groups", flattenURLFilteringOverrideGroups(resp.OverrideGroups)); err != nil {
			return err
		}

		if err := d.Set("location_groups", flattenURLFilteringLocationGroups(resp.LocationGroups)); err != nil {
			return err
		}

		if err := d.Set("labels", flattenURLFilteringLabels(resp.Labels)); err != nil {
			return err
		}

		if err := d.Set("last_modified_by", flattenURLFilteringLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any url filtering rule with name '%s' or id '%d'", name, id)
	}

	return nil
}

func flattenURLFilteringLocation(locations []urlfilteringpolicies.Locations) []interface{} {
	location := make([]interface{}, len(locations))
	for i, val := range locations {
		location[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return location
}

func flattenURLFilteringGroups(groups []urlfilteringpolicies.Groups) []interface{} {
	group := make([]interface{}, len(groups))
	for i, val := range groups {
		group[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return group
}

func flattenURLFilteringDepartments(departments []urlfilteringpolicies.Departments) []interface{} {
	department := make([]interface{}, len(departments))
	for i, val := range departments {
		department[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return department
}

func flattenURLFilteringUsers(users []urlfilteringpolicies.Users) []interface{} {
	user := make([]interface{}, len(users))
	for i, val := range users {
		user[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return user
}

func flattenURLFilteringTimeWindows(timeWindows []urlfilteringpolicies.TimeWindows) []interface{} {
	timeWindow := make([]interface{}, len(timeWindows))
	for i, val := range timeWindows {
		timeWindow[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return timeWindow
}

func flattenURLFilteringOverrideUsers(overrideUsers []urlfilteringpolicies.OverrideUsers) []interface{} {
	override := make([]interface{}, len(overrideUsers))
	for i, val := range overrideUsers {
		override[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return override
}

func flattenURLFilteringOverrideGroups(overrideGroups []urlfilteringpolicies.OverrideGroups) []interface{} {
	override := make([]interface{}, len(overrideGroups))
	for i, val := range overrideGroups {
		override[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return override
}

func flattenURLFilteringLocationGroups(overrideLocationGroups []urlfilteringpolicies.LocationGroups) []interface{} {
	override := make([]interface{}, len(overrideLocationGroups))
	for i, val := range overrideLocationGroups {
		override[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return override
}

func flattenURLFilteringLabels(labels []urlfilteringpolicies.Labels) []interface{} {
	label := make([]interface{}, len(labels))
	for i, val := range labels {
		label[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return label
}

func flattenURLFilteringLastModifiedBy(lastModifiedBy *urlfilteringpolicies.LastModifiedBy) interface{} {
	return []map[string]interface{}{
		{
			"id":         lastModifiedBy.ID,
			"name":       lastModifiedBy.Name,
			"extensions": lastModifiedBy.Extensions,
		},
	}
}
