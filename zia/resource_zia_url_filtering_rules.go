package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlfilteringpolicies"
)

func resourceURLFilteringRules() *schema.Resource {
	return &schema.Resource{
		Create:   resourceURLFilteringRulesCreate,
		Read:     resourceURLFilteringRulesRead,
		Update:   resourceURLFilteringRulesUpdate,
		Delete:   resourceURLFilteringRulesDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocols": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"locations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Name-ID pairs of locations for which rule must be applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Name-ID pairs of groups for which rule must be applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"departments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Name-ID pairs of departments for which rule must be applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"users": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Name-ID pairs of users for which rule must be applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"url_categories": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"time_windows": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Name-ID pairs of time interval during which rule must be enforced.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Admin rank of the admin who creates this rule",
			},
			"request_methods": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request method for which the rule must be applied. If not set, rule will be applied to all methods",
			},
			"end_user_notification_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"override_users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"override_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"block_override": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"time_quota": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"size_quota": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"validity_start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"validity_end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"validity_time_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_modified_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"last_modified_by": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
			"enforce_time_validity": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"NONE",
					"BLOCK",
					"CAUTION",
					"ALLOW",
					"ICAP_RESPONSE",
				}, false),
			},
			"ciparule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the CIPA Compliance rule is enabled",
			},
		},
	}
}

func resourceURLFilteringRulesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandURLFilteringRules(d)
	log.Printf("[INFO] Creating url filtering rule\n%+v\n", req)

	resp, _, err := zClient.urlfilteringpolicies.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia url filtering rule request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))

	return resourceURLFilteringRulesRead(d, m)
}

func resourceURLFilteringRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "id")
	if !ok {
		return fmt.Errorf("no url filtering rule id is set")
	}
	resp, err := zClient.urlfilteringpolicies.Get(id)

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing zia url filtering rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting url category :\n%+v\n", resp)
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

	return nil
}

func resourceURLFilteringRulesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "id")
	if !ok {
		log.Printf("[ERROR] url filtering rule ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating url filtering rule ID: %v\n", id)
	req := expandURLFilteringRules(d)

	if _, _, err := zClient.urlfilteringpolicies.Update(id, &req); err != nil {
		return err
	}

	return resourceURLFilteringRulesRead(d, m)
}

func resourceURLFilteringRulesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "id")
	if !ok {
		log.Printf("[ERROR] url filtering rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting url filtering rule ID: %v\n", id)

	if _, err := zClient.urlfilteringpolicies.Delete(id); err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[INFO] url filtering rule deleted")
	return nil
}

func expandURLFilteringRules(d *schema.ResourceData) urlfilteringpolicies.URLFilteringRule {
	id, _ := getIntFromResourceData(d, "id")
	result := urlfilteringpolicies.URLFilteringRule{
		ID:                     id,
		Name:                   d.Get("name").(string),
		Order:                  d.Get("order").(int),
		Protocols:              ListToStringSlice(d.Get("protocols").([]interface{})),
		URLCategories:          ListToStringSlice(d.Get("url_categories").([]interface{})),
		State:                  d.Get("state").(string),
		Rank:                   d.Get("rank").(int),
		RequestMethods:         ListToStringSlice(d.Get("request_methods").([]interface{})),
		EndUserNotificationURL: d.Get("end_user_notification_url").(string),
		BlockOverride:          d.Get("block_override").(bool),
		TimeQuota:              d.Get("time_quota").(int),
		SizeQuota:              d.Get("size_quota").(int),
		Description:            d.Get("description").(string),
		ValidityStartTime:      d.Get("validity_start_time").(int),
		ValidityEndTime:        d.Get("validity_end_time").(int),
		ValidityTimeZoneID:     d.Get("validity_time_zone_id").(string),
		LastModifiedTime:       d.Get("last_modified_time").(int),
		EnforceTimeValidity:    d.Get("enforce_time_validity").(bool),
		Action:                 d.Get("action").(string),
		Ciparule:               d.Get("ciparule").(bool),
	}
	locations := expandURLFilteringLocations(d)
	if locations != nil {
		result.Locations = locations
	}
	groups := expandURLFilteringGroups(d)
	if groups != nil {
		result.Groups = groups
	}
	departments := expandURLFilteringDepartments(d)
	if departments != nil {
		result.Departments = departments
	}
	users := expandURLFilteringUsers(d)
	if users != nil {
		result.Users = users
	}
	timeWindows := expandURLFilteringTimeWindows(d)
	if timeWindows != nil {
		result.TimeWindows = timeWindows
	}
	overrideUsers := expandURLFilteringOverrideUsers(d)
	if overrideUsers != nil {
		result.OverrideUsers = overrideUsers
	}
	overrideGroups := expandURLFilteringOverrideGroups(d)
	if overrideGroups != nil {
		result.OverrideGroups = overrideGroups
	}
	locationGroups := expandURLFilteringLocationGroups(d)
	if locationGroups != nil {
		result.LocationGroups = locationGroups
	}
	labels := expandURLFilteringLabels(d)
	if labels != nil {
		result.Labels = labels
	}
	lastModifiedBy := expandURLFilteringLastModifiedBy(d)
	if lastModifiedBy != nil {
		result.LastModifiedBy = lastModifiedBy
	}
	return result
}

func expandURLFilteringLocations(d *schema.ResourceData) []urlfilteringpolicies.Locations {
	var locations []urlfilteringpolicies.Locations
	if locationsInterface, ok := d.GetOk("locations"); ok {
		location := locationsInterface.([]interface{})
		locations = make([]urlfilteringpolicies.Locations, len(location))
		for i, val := range location {
			locationItem := val.(map[string]interface{})
			locations[i] = urlfilteringpolicies.Locations{
				ID:         locationItem["id"].(int),
				Name:       locationItem["name"].(string),
				Extensions: locationItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return locations
}

func expandURLFilteringGroups(d *schema.ResourceData) []urlfilteringpolicies.Groups {
	var groups []urlfilteringpolicies.Groups
	if groupsInterface, ok := d.GetOk("groups"); ok {
		group := groupsInterface.([]interface{})
		groups = make([]urlfilteringpolicies.Groups, len(group))
		for i, val := range group {
			groupItem := val.(map[string]interface{})
			groups[i] = urlfilteringpolicies.Groups{
				ID:         groupItem["id"].(int),
				Name:       groupItem["name"].(string),
				Extensions: groupItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return groups
}

func expandURLFilteringDepartments(d *schema.ResourceData) []urlfilteringpolicies.Departments {
	var departments []urlfilteringpolicies.Departments
	if departmentsInterface, ok := d.GetOk("departments"); ok {
		department := departmentsInterface.([]interface{})
		departments = make([]urlfilteringpolicies.Departments, len(department))
		for i, val := range department {
			departmentItem := val.(map[string]interface{})
			departments[i] = urlfilteringpolicies.Departments{
				ID:         departmentItem["id"].(int),
				Name:       departmentItem["name"].(string),
				Extensions: departmentItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return departments
}

func expandURLFilteringUsers(d *schema.ResourceData) []urlfilteringpolicies.Users {
	var users []urlfilteringpolicies.Users
	if usersInterface, ok := d.GetOk("users"); ok {
		user := usersInterface.([]interface{})
		users = make([]urlfilteringpolicies.Users, len(user))
		for i, val := range user {
			userItem := val.(map[string]interface{})
			users[i] = urlfilteringpolicies.Users{
				ID:         userItem["id"].(int),
				Name:       userItem["name"].(string),
				Extensions: userItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return users
}

func expandURLFilteringTimeWindows(d *schema.ResourceData) []urlfilteringpolicies.TimeWindows {
	var timewindows []urlfilteringpolicies.TimeWindows
	if timewindowsInterface, ok := d.GetOk("time_windows"); ok {
		timewindow := timewindowsInterface.([]interface{})
		timewindows = make([]urlfilteringpolicies.TimeWindows, len(timewindow))
		for i, val := range timewindow {
			timewindowsItem := val.(map[string]interface{})
			timewindows[i] = urlfilteringpolicies.TimeWindows{
				ID:         timewindowsItem["id"].(int),
				Name:       timewindowsItem["name"].(string),
				Extensions: timewindowsItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return timewindows
}

func expandURLFilteringOverrideUsers(d *schema.ResourceData) []urlfilteringpolicies.OverrideUsers {
	var overrideUsers []urlfilteringpolicies.OverrideUsers
	if overrideUsersInterface, ok := d.GetOk("override_users"); ok {
		overrideUser := overrideUsersInterface.([]interface{})
		overrideUsers = make([]urlfilteringpolicies.OverrideUsers, len(overrideUser))
		for i, val := range overrideUser {
			overrideUserItem := val.(map[string]interface{})
			overrideUsers[i] = urlfilteringpolicies.OverrideUsers{
				ID:         overrideUserItem["id"].(int),
				Name:       overrideUserItem["name"].(string),
				Extensions: overrideUserItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return overrideUsers
}

func expandURLFilteringOverrideGroups(d *schema.ResourceData) []urlfilteringpolicies.OverrideGroups {
	var overrideGroups []urlfilteringpolicies.OverrideGroups
	if overrideGroupsInterface, ok := d.GetOk("override_groups"); ok {
		overrideGroup := overrideGroupsInterface.([]interface{})
		overrideGroups = make([]urlfilteringpolicies.OverrideGroups, len(overrideGroup))
		for i, val := range overrideGroup {
			overrideGroupItem := val.(map[string]interface{})
			overrideGroups[i] = urlfilteringpolicies.OverrideGroups{
				ID:         overrideGroupItem["id"].(int),
				Name:       overrideGroupItem["name"].(string),
				Extensions: overrideGroupItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return overrideGroups
}

func expandURLFilteringLocationGroups(d *schema.ResourceData) []urlfilteringpolicies.LocationGroups {
	var locationGroups []urlfilteringpolicies.LocationGroups
	if locationGroupsInterface, ok := d.GetOk("location_groups"); ok {
		locationGroup := locationGroupsInterface.([]interface{})
		locationGroups = make([]urlfilteringpolicies.LocationGroups, len(locationGroup))
		for i, val := range locationGroup {
			locationGroupItem := val.(map[string]interface{})
			locationGroups[i] = urlfilteringpolicies.LocationGroups{
				ID:         locationGroupItem["id"].(int),
				Name:       locationGroupItem["name"].(string),
				Extensions: locationGroupItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return locationGroups
}

func expandURLFilteringLabels(d *schema.ResourceData) []urlfilteringpolicies.Labels {
	var labels []urlfilteringpolicies.Labels
	if labelsInterface, ok := d.GetOk("labels"); ok {
		label := labelsInterface.([]interface{})
		labels = make([]urlfilteringpolicies.Labels, len(label))
		for i, val := range label {
			labelItem := val.(map[string]interface{})
			labels[i] = urlfilteringpolicies.Labels{
				ID:         labelItem["id"].(int),
				Name:       labelItem["name"].(string),
				Extensions: labelItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return labels
}

func expandURLFilteringLastModifiedBy(d *schema.ResourceData) *urlfilteringpolicies.LastModifiedBy {
	lastModifiedByObj, ok := d.GetOk("last_modified_by")
	if !ok {
		return nil
	}
	lastMofiedBy, ok := lastModifiedByObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(lastMofiedBy.List()) > 0 {
		lastModifiedByObj := lastMofiedBy.List()[0]
		lastMofied, ok := lastModifiedByObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &urlfilteringpolicies.LastModifiedBy{
			ID:         lastMofied["id"].(int),
			Name:       lastMofied["name"].(string),
			Extensions: lastMofied["extensions"].(map[string]interface{}),
		}
	}
	return nil
}
