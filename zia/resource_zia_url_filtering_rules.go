package zia

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlfilteringpolicies"
)

type listrules struct {
	orders map[int]int
	sync.Mutex
}

var rules = listrules{
	orders: make(map[int]int),
}

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
			"rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"protocols": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"url_categories": {
				Type:     schema.TypeSet,
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
			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     7,
				Description: "Admin rank of the admin who creates this rule",
			},
			"request_methods": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request method for which the rule must be applied. If not set, rule will be applied to all methods",
			},
			"end_user_notification_url": {
				Type:     schema.TypeString,
				Optional: true,
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
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"user_agent_types": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"OPERA",
					"FIREFOX",
					"MSIE",
					"MSEDGE",
					"CHROME",
					"SAFARI",
					"OTHER",
				}, false),
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
					"ISOLATE",
					"ICAP_RESPONSE",
				}, false),
			},
			"ciparule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the CIPA Compliance rule is enabled",
			},
			"locations":       listIDsSchemaType("Name-ID pairs of locations for which rule must be applied"),
			"groups":          listIDsSchemaType("Name-ID pairs of groups for which rule must be applied"),
			"departments":     listIDsSchemaType("Name-ID pairs of departments for which rule must be applied"),
			"users":           listIDsSchemaType("Name-ID pairs of users for which rule must be applied"),
			"time_windows":    listIDsSchemaType("Name-ID pairs of time interval during which rule must be enforced."),
			"override_users":  listIDsSchemaType("list of override users"),
			"override_groups": listIDsSchemaType("list of override groups"),
			"location_groups": listIDsSchemaType("list of locations groups"),
			"labels":          listIDsSchemaType("list of labels"),
		},
	}
}

func resourceURLFilteringRulesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandURLFilteringRules(d)
	log.Printf("[INFO] Creating url filtering rule\n%+v\n", req)
	orderObj, orderIsSet := d.GetOk("order")
	if orderIsSet {
		// always set it to 1, and let the re-ordering happen after ( because having an invalid order will cause a bad request)
		req.Order = 1
	}
	resp, err := zClient.urlfilteringpolicies.Create(&req)
	if err != nil {
		return err
	}
	if orderIsSet {
		req.Order = orderObj.(int)
		go reorder(req.Order, resp.ID, zClient)
	}
	log.Printf("[INFO] Created zia url filtering rule request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_id", resp.ID)

	return resourceURLFilteringRulesRead(d, m)
}

func resourceURLFilteringRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return fmt.Errorf("no url filtering rule id is set")
	}
	resp, err := zClient.urlfilteringpolicies.Get(id)

	if err != nil {
		if obj, ok := err.(*client.ErrorResponse); ok && obj.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia url filtering rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting url category :\n%+v\n", resp)
	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("protocols", resp.Protocols)
	_ = d.Set("url_categories", resp.URLCategories)
	_ = d.Set("state", resp.State)
	_ = d.Set("user_agent_types", resp.UserAgentTypes)
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

	if err := d.Set("locations", flattenIDs(resp.Locations)); err != nil {
		return err
	}

	if err := d.Set("groups", flattenIDs(resp.Groups)); err != nil {
		return err
	}

	if err := d.Set("departments", flattenIDs(resp.Departments)); err != nil {
		return err
	}

	if err := d.Set("users", flattenIDs(resp.Users)); err != nil {
		return err
	}

	if err := d.Set("time_windows", flattenIDs(resp.TimeWindows)); err != nil {
		return err
	}

	if err := d.Set("override_users", flattenIDs(resp.OverrideUsers)); err != nil {
		return err
	}

	if err := d.Set("override_groups", flattenIDs(resp.OverrideGroups)); err != nil {
		return err
	}

	if err := d.Set("location_groups", flattenIDs(resp.LocationGroups)); err != nil {
		return err
	}

	if err := d.Set("labels", flattenIDs(resp.Labels)); err != nil {
		return err
	}

	if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
		return err
	}

	return nil
}

func resourceURLFilteringRulesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] url filtering rule ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating url filtering rule ID: %v\n", id)
	req := expandURLFilteringRules(d)
	if d.HasChange("order") {
		_, orderIsSet := d.GetOk("order")
		if orderIsSet {
			go reorder(req.Order, req.ID, zClient)
		}
		req.Order = 1
	}
	if _, _, err := zClient.urlfilteringpolicies.Update(id, &req); err != nil {
		return err
	}

	return resourceURLFilteringRulesRead(d, m)
}

func resourceURLFilteringRulesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] url filtering rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting url filtering rule ID: %v\n", (d.Id()))

	if _, err := zClient.urlfilteringpolicies.Delete(id); err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[INFO] url filtering rule deleted")
	return nil
}

func expandURLFilteringRules(d *schema.ResourceData) urlfilteringpolicies.URLFilteringRule {
	id, _ := getIntFromResourceData(d, "rule_id")
	result := urlfilteringpolicies.URLFilteringRule{
		ID:                     id,
		Name:                   d.Get("name").(string),
		Order:                  d.Get("order").(int),
		Protocols:              SetToStringList(d, "protocols"),
		URLCategories:          SetToStringList(d, "url_categories"),
		State:                  d.Get("state").(string),
		UserAgentTypes:         SetToStringList(d, "user_agent_types"),
		Rank:                   d.Get("rank").(int),
		RequestMethods:         SetToStringList(d, "request_methods"),
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
	locations := expandIDNameExtensionsSet(d, "locations")
	if locations != nil {
		result.Locations = locations
	}
	groups := expandIDNameExtensionsSet(d, "groups")
	if groups != nil {
		result.Groups = groups
	}
	departments := expandIDNameExtensionsSet(d, "departments")
	if departments != nil {
		result.Departments = departments
	}
	users := expandIDNameExtensionsSet(d, "users")
	if users != nil {
		result.Users = users
	}
	timeWindows := expandIDNameExtensionsSet(d, "time_windows")
	if timeWindows != nil {
		result.TimeWindows = timeWindows
	}
	overrideUsers := expandIDNameExtensionsSet(d, "override_users")
	if overrideUsers != nil {
		result.OverrideUsers = overrideUsers
	}
	overrideGroups := expandIDNameExtensionsSet(d, "override_groups")
	if overrideGroups != nil {
		result.OverrideGroups = overrideGroups
	}
	locationGroups := expandIDNameExtensionsSet(d, "location_groups")
	if locationGroups != nil {
		result.LocationGroups = locationGroups
	}
	labels := expandIDNameExtensionsSet(d, "labels")
	if labels != nil {
		result.Labels = labels
	}
	lastModifiedBy := expandIDNameExtensions(d, "last_modified_by")
	if lastModifiedBy != nil {
		result.LastModifiedBy = lastModifiedBy
	}
	return result
}

func reorder(order, id int, zClient *Client) {
	defer reorderAll(zClient)
	rules.Lock()
	rules.orders[id] = order
	rules.Unlock()
}

// we keep calling reordering endpoint to reorder all rules after new rule was added
// because the reorder endpoint shifts all order up to replac the new order.
func reorderAll(zClient *Client) {
	rules.Lock()
	defer rules.Unlock()
	count := zClient.urlfilteringpolicies.RulesCount()
	for k, v := range rules.orders {
		// the only valid order you can set is 0,count
		if v <= count {
			_, err := zClient.urlfilteringpolicies.Reorder(k, v)
			if err != nil {
				log.Printf("[ERROR] couldn't reorder the url filtering policy, the order may not have taken place: %v\n", err)
			}
		}
	}
}
