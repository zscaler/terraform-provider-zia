package zia

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/urlfilteringpolicies"
)

/*
type listrules struct {
	orders map[int]int
	sync.Mutex
}

var rules = listrules{
	orders: make(map[int]int),
}
*/

func resourceURLFilteringRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceURLFilteringRulesCreate,
		Read:   resourceURLFilteringRulesRead,
		Update: resourceURLFilteringRulesUpdate,
		Delete: resourceURLFilteringRulesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("rule_id", idInt)
				} else {
					resp, err := zClient.urlfilteringpolicies.GetByName(id)
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL Filtering Rule ID",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "URL Filtering Rule ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule Name",
			},
			"order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Order of execution of rule with respect to other URL Filtering rules",
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(1, 7),
				Description:  "Admin rank of the admin who creates this rule",
			},
			"end_user_notification_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of end user notification page to be displayed when the rule is matched. Not applicable if either 'overrideUsers' or 'overrideGroups' is specified.",
			},
			"block_override": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"time_quota": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(15, 600),
				Description:  "Time quota in minutes, after which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to 'BLOCK', this field is not applicable.",
			},
			"size_quota": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(10, 100000),
				Description:  "Size quota in KB beyond which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to 'BLOCK', this field is not applicable.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
				Description:  "Additional information about the URL Filtering rule",
			},
			"validity_start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If enforceTimeValidity is set to true, the URL Filtering rule will be valid starting on this date and time.",
			},
			"validity_end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If enforceTimeValidity is set to true, the URL Filtering rule will cease to be valid on this end date and time.",
			},
			"validity_time_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If enforceTimeValidity is set to true, the URL Filtering rule date and time will be valid based on this time zone ID.",
			},
			"enforce_time_validity": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enforce a set a validity time period for the URL Filtering rule.",
			},
			"user_agent_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Action taken when traffic matches rule criteria",
				ValidateFunc: validation.StringInSlice([]string{
					"BLOCK",
					"CAUTION",
					"ALLOW",
					"ISOLATE",
				}, false),
			},
			"ciparule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, the CIPA Compliance rule is enabled",
			},
			"locations":           listIDsSchemaTypeCustom(8, "Name-ID pairs of locations for which rule must be applied"),
			"groups":              listIDsSchemaTypeCustom(8, "Name-ID pairs of groups for which rule must be applied"),
			"departments":         listIDsSchemaTypeCustom(8, "Name-ID pairs of departments for which rule must be applied"),
			"users":               listIDsSchemaTypeCustom(4, "Name-ID pairs of users for which rule must be applied"),
			"time_windows":        listIDsSchemaType("Name-ID pairs of time interval during which rule must be enforced."),
			"override_users":      listIDsSchemaType("Name-ID pairs of users for which this rule can be overridden."),
			"override_groups":     listIDsSchemaTypeCustom(8, "Name-ID pairs of groups for which this rule can be overridden."),
			"device_groups":       listIDsSchemaType("This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":             listIDsSchemaType("Name-ID pairs of devices for which rule must be applied."),
			"location_groups":     listIDsSchemaTypeCustom(32, "Name-ID pairs of the location groups to which the rule must be applied."),
			"labels":              listIDsSchemaType("The URL Filtering rule's label."),
			"device_trust_levels": getDeviceTrustLevels(),
			"url_categories":      getURLCategories(),
			"request_methods":     getURLRequestMethods(),
			"protocols":           getURLProtocols(),
		},
	}
}

func currentOrderVsRankWording(zClient *Client) string {
	list, err := zClient.urlfilteringpolicies.GetAll()
	if err != nil {
		return ""
	}
	result := ""
	for i, r := range list {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("Rank %d VS Order %d", r.Rank, r.Order)

	}
	return result
}

func resourceURLFilteringRulesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandURLFilteringRules(d)
	log.Printf("[INFO] Creating url filtering rule\n%+v\n", req)
	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate)-time.Minute, func() *resource.RetryError {
		order := req.Order
		list, _ := zClient.urlfilteringpolicies.GetAll()
		req.Order = len(list) + 1
		resp, err := zClient.urlfilteringpolicies.Create(&req)
		if err != nil {
			// check is err matches regex "Rule with rank [0-9]+ is not allowed at order [0-9]+
			// if so, retry
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") && reg.MatchString(err.Error()) {
				return resource.NonRetryableError(fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, order, req.Rank, currentOrderVsRankWording(zClient), err))
			}
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Creating url filtering rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
				time.Sleep(time.Second * time.Duration(order+1))
				return resource.RetryableError(errors.New("expected resource to be created but was not"))
			}
			return resource.NonRetryableError(fmt.Errorf("error creating resource: %s", err))
		}
		log.Printf("[INFO] Created url filtering rule request. ID: %v\n", resp)
		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		err = resourceURLFilteringRulesRead(d, m)
		if err != nil {
			return resource.NonRetryableError(err)
		} else {
			reorder(order, resp.ID, "url_filtering_rules", func() (int, error) {
				list, err := zClient.urlfilteringpolicies.GetAll()
				return len(list), err

			}, func(id, order int) error {
				rule, err := zClient.urlfilteringpolicies.Get(id)
				if err != nil {
					return err
				}
				rule.Order = order
				_, _, err = zClient.urlfilteringpolicies.Update(id, rule)
				return err
			})
			return nil
		}
	})
}

func resourceURLFilteringRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return fmt.Errorf("no url filtering rule id is set")
	}
	resp, err := zClient.urlfilteringpolicies.Get(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
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
	_ = d.Set("description", resp.Description)
	_ = d.Set("protocols", resp.Protocols)
	if len(resp.URLCategories) == 0 {
		_ = d.Set("url_categories", []string{"ANY"})
	} else {
		_ = d.Set("url_categories", resp.URLCategories)
	}
	_ = d.Set("state", resp.State)
	_ = d.Set("user_agent_types", resp.UserAgentTypes)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)
	_ = d.Set("request_methods", resp.RequestMethods)
	_ = d.Set("end_user_notification_url", resp.EndUserNotificationURL)
	_ = d.Set("block_override", resp.BlockOverride)
	_ = d.Set("time_quota", resp.TimeQuota)
	_ = d.Set("size_quota", resp.SizeQuota)
	_ = d.Set("validity_start_time", resp.ValidityStartTime)
	_ = d.Set("validity_end_time", resp.ValidityEndTime)
	_ = d.Set("validity_time_zone_id", resp.ValidityTimeZoneID)
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

	if err := d.Set("device_groups", flattenIDs(resp.DeviceGroups)); err != nil {
		return err
	}

	if err := d.Set("devices", flattenIDs(resp.Devices)); err != nil {
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
	if _, err := zClient.urlfilteringpolicies.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate)-time.Minute, func() *resource.RetryError {
		_, _, err := zClient.urlfilteringpolicies.Update(id, &req)
		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				time.Sleep(time.Second * time.Duration(req.Order+1))
				log.Printf("[INFO] Updating url filtering rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				return resource.RetryableError(errors.New("expected resource to be updated but was not"))
			}
			return resource.NonRetryableError(fmt.Errorf("error updating resource: %s", err))
		}

		err = resourceURLFilteringRulesRead(d, m)
		if err != nil {
			return resource.NonRetryableError(err)
		} else {
			reorder(req.Order, req.ID, "url_filtering_rules", func() (int, error) {
				list, err := zClient.urlfilteringpolicies.GetAll()
				return len(list), err

			}, func(id, order int) error {
				rule, err := zClient.urlfilteringpolicies.Get(id)
				if err != nil {
					return err
				}
				rule.Order = order
				_, _, err = zClient.urlfilteringpolicies.Update(id, rule)
				return err
			})
			return nil
		}
	})
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
		Description:            d.Get("description").(string),
		Order:                  d.Get("order").(int),
		Protocols:              SetToStringList(d, "protocols"),
		URLCategories:          SetToStringList(d, "url_categories"),
		DeviceTrustLevels:      SetToStringList(d, "device_trust_levels"),
		RequestMethods:         SetToStringList(d, "request_methods"),
		UserAgentTypes:         SetToStringList(d, "user_agent_types"),
		State:                  d.Get("state").(string),
		Rank:                   d.Get("rank").(int),
		EndUserNotificationURL: d.Get("end_user_notification_url").(string),
		BlockOverride:          d.Get("block_override").(bool),
		TimeQuota:              d.Get("time_quota").(int),
		SizeQuota:              d.Get("size_quota").(int),
		ValidityStartTime:      d.Get("validity_start_time").(int),
		ValidityEndTime:        d.Get("validity_end_time").(int),
		ValidityTimeZoneID:     d.Get("validity_time_zone_id").(string),
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
	deviceGroups := expandIDNameExtensionsSet(d, "device_groups")
	if deviceGroups != nil {
		result.DeviceGroups = deviceGroups
	}
	devices := expandIDNameExtensionsSet(d, "devices")
	if devices != nil {
		result.Devices = devices
	}
	return result
}

/*

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
*/
