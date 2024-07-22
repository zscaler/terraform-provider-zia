package zia

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/cloudappcontrol"
)

func resourceCloudAppControlRules() *schema.Resource {
	return &schema.Resource{
		Create:        resourceCloudAppControlRulesCreate,
		Read:          resourceCloudAppControlRulesRead,
		Update:        resourceCloudAppControlRulesUpdate,
		Delete:        resourceCloudAppControlRulesDelete,
		CustomizeDiff: validateActionsCustomizeDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)
				service := zClient.cloudappcontrol

				id := d.Id()
				var ruleType, identifier string

				// Check if id contains a colon to split rule type and identifier
				if strings.Contains(id, ":") {
					parts := strings.SplitN(id, ":", 2)
					ruleType = parts[0]
					identifier = parts[1]
				} else {
					// If no colon, treat entire id as the identifier (assuming it's the rule ID for now)
					return nil, fmt.Errorf("invalid import format: expected 'rule_type:rule_id' or 'rule_type:rule_name'")
				}

				// Check if the identifier is a rule ID
				idInt, parseIDErr := strconv.Atoi(identifier)
				if parseIDErr == nil {
					// If identifier is an ID
					resp, err := cloudappcontrol.GetByRuleID(service, ruleType, idInt)
					if err != nil {
						return nil, err
					}
					d.SetId(strconv.Itoa(resp.ID))
					_ = d.Set("rule_id", resp.ID)
					_ = d.Set("type", ruleType)
				} else {
					// If identifier is a name
					resources, err := cloudappcontrol.GetByRuleType(service, ruleType)
					if err != nil {
						return nil, err
					}
					for _, r := range resources {
						if r.Name == identifier {
							d.SetId(strconv.Itoa(r.ID))
							_ = d.Set("rule_id", r.ID)
							_ = d.Set("type", ruleType)
							break
						}
					}
					if d.Id() == "" {
						return nil, fmt.Errorf("couldn't find any cloud application rule with name '%s'", identifier)
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier assigned to the forwarding rule",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "A unique identifier assigned to the forwarding rule",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the forwarding rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the forwarding rule",
			},
			"order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The order of execution for the forwarding rule order",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Determines whether the Firewall Filtering policy rule is enabled or disabled",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Admin rank assigned to the forwarding rule",
			},
			"actions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Actions allowed for the specified type.",
			},
			"applications": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
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
			"enforce_time_validity": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enforce a set a validity time period for the URL Filtering rule.",
			},
			"validity_start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If enforceTimeValidity is set to true, the URL Filtering rule is valid starting on this date and time.",
			},
			"validity_end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If enforceTimeValidity is set to true, the URL Filtering rule ceases to be valid on this end date and time.",
			},
			"validity_time_zone_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateTimeZone,
				Description:  "If enforceTimeValidity is set to true, the URL Filtering rule date and time is valid based on this time zone ID. Use IANA Format TimeZone.",
			},
			"cascading_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The unique identifier for the device.",
			},
			"cbi_profile": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"device_groups":          setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":                setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"time_windows":           setIDsSchemaTypeCustom(nil, "Name-ID pairs of time interval during which rule must be enforced."),
			"labels":                 setIDsSchemaTypeCustom(nil, "The URL Filtering rule's label."),
			"locations":              setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of locations for which rule must be applied"),
			"location_groups":        setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the location groups to which the rule must be applied."),
			"groups":                 setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of groups for which rule must be applied"),
			"departments":            setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of departments for which rule must be applied"),
			"users":                  setIDsSchemaTypeCustom(intPtr(4), "Name-ID pairs of users for which rule must be applied"),
			"device_trust_levels":    getDeviceTrustLevels(),
			"user_risk_score_levels": getUserRiskScoreLevels(),
			"user_agent_types":       getUserAgentTypes(),
			"type":                   getAppControlType(),
		},
	}
}

func resourceCloudAppControlRulesCreate(d *schema.ResourceData, m interface{}) error {

	zClient := m.(*Client)
	service := zClient.cloudappcontrol

	req := expandCloudAppControlRules(d)
	log.Printf("[INFO] Creating zia cloud app control rule\n%+v\n", req)

	resp, err := cloudappcontrol.Create(service, req.Type, &req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created zia cloud app control request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_id", resp.ID)
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceCloudAppControlRulesRead(d, m)
}

func resourceCloudAppControlRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.cloudappcontrol

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return fmt.Errorf("no zia cloud app control rule id is set")
	}
	ruleType, ok := d.Get("type").(string)
	if !ok || ruleType == "" {
		return fmt.Errorf("no rule type is set")
	}
	resp, err := cloudappcontrol.GetByRuleID(service, ruleType, id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing cloud app control rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting cloud app control rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("order", resp.Order)
	_ = d.Set("actions", resp.Actions)
	_ = d.Set("state", resp.State)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("type", resp.Type)
	_ = d.Set("time_quota", resp.TimeQuota)
	_ = d.Set("cascading_enabled", resp.CascadingEnabled)
	_ = d.Set("applications", resp.Applications)
	_ = d.Set("enforce_time_validity", resp.EnforceTimeValidity)
	_ = d.Set("validity_start_time", resp.ValidityStartTime)
	_ = d.Set("validity_end_time", resp.ValidityEndTime)
	_ = d.Set("validity_time_zone_id", resp.ValidityTimeZoneID)
	_ = d.Set("user_agent_types", resp.UserAgentTypes)
	_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)
	_ = d.Set("user_risk_score_levels", resp.UserRiskScoreLevels)

	// Convert size_quota from KB back to MB
	sizeQuotaMB := resp.SizeQuota / 1024
	_ = d.Set("size_quota", sizeQuotaMB)

	// Set validity_start_time only if it is not the default value
	if resp.ValidityStartTime != 0 {
		_ = d.Set("validity_start_time", time.Unix(int64(resp.ValidityStartTime), 0).UTC().Format(time.RFC1123))
	} else {
		_ = d.Set("validity_start_time", nil)
	}

	// Set validity_end_time only if it is not the default value
	if resp.ValidityEndTime != 0 {
		_ = d.Set("validity_end_time", time.Unix(int64(resp.ValidityEndTime), 0).UTC().Format(time.RFC1123))
	} else {
		_ = d.Set("validity_end_time", nil)
	}

	_ = d.Set("validity_time_zone_id", resp.ValidityTimeZoneID)
	_ = d.Set("enforce_time_validity", resp.EnforceTimeValidity)

	// Update the cbi_profile block in the state
	if resp.CBIProfile.ID != "" {
		if err := d.Set("cbi_profile", flattenCloudAppControlCBIProfileSimple(&resp.CBIProfile)); err != nil {
			return err
		}
	}

	if err := d.Set("locations", flattenIDs(resp.Locations)); err != nil {
		return err
	}

	if err := d.Set("location_groups", flattenIDs(resp.LocationGroups)); err != nil {
		return err
	}

	if err := d.Set("groups", flattenIDs(resp.Groups)); err != nil {
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

	if err := d.Set("device_groups", flattenIDs(resp.DeviceGroups)); err != nil {
		return err
	}

	if err := d.Set("devices", flattenIDs(resp.Devices)); err != nil {
		return err
	}

	if err := d.Set("labels", flattenIDs(resp.Labels)); err != nil {
		return err
	}

	if err := d.Set("time_windows", flattenIDs(resp.TimeWindows)); err != nil {
		return err
	}

	return nil
}

func resourceCloudAppControlRulesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.cloudappcontrol

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] cloud application control rule ID not set: %v\n", id)
	}

	ruleType, ok := d.Get("type").(string)
	if !ok || ruleType == "" {
		return fmt.Errorf("no rule type is set")
	}

	log.Printf("[INFO] Updating zia cloud application control rule ID: %v\n", id)
	req := expandCloudAppControlRules(d)

	if _, err := cloudappcontrol.GetByRuleID(service, ruleType, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := cloudappcontrol.Update(service, ruleType, id, &req); err != nil {
		return err
	}
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceCloudAppControlRulesRead(d, m)
}

func resourceCloudAppControlRulesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.cloudappcontrol

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] cloud application control rule not set: %v\n", id)
	}
	ruleType, ok := d.Get("type").(string)
	if !ok || ruleType == "" {
		return fmt.Errorf("no rule type is set")
	}
	log.Printf("[INFO] Deleting cloud application control rule ID: %v\n", (d.Id()))

	if _, err := cloudappcontrol.Delete(service, ruleType, id); err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[INFO] cloud application control rule deleted")
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandCloudAppControlRules(d *schema.ResourceData) cloudappcontrol.WebApplicationRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	validityStartTimeStr := d.Get("validity_start_time").(string)
	validityEndTimeStr := d.Get("validity_end_time").(string)

	var validityStartTime int
	var validityEndTime int
	var err error

	if validityStartTimeStr != "" {
		log.Printf("[INFO] Converting validity_start_time: %s", validityStartTimeStr)
		validityStartTime, err = ConvertRFC1123ToEpoch(validityStartTimeStr)
		if err != nil {
			log.Printf("[ERROR] Invalid validity_start_time: %v", err)
			validityStartTime = 0
		} else {
			log.Printf("[INFO] Converted validity_start_time: %d", validityStartTime)
		}
	} else {
		log.Printf("[INFO] validity_start_time is empty")
	}

	if validityEndTimeStr != "" {
		log.Printf("[INFO] Converting validity_end_time: %s", validityEndTimeStr)
		validityEndTime, err = ConvertRFC1123ToEpoch(validityEndTimeStr)
		if err != nil {
			log.Printf("[ERROR] Invalid validity_end_time: %v", err)
			validityEndTime = 0
		} else {
			log.Printf("[INFO] Converted validity_end_time: %d", validityEndTime)
		}
	} else {
		log.Printf("[INFO] validity_end_time is empty")
	}

	sizeQuotaMB := d.Get("size_quota").(int)
	sizeQuotaKB, err := convertAndValidateSizeQuota(sizeQuotaMB)
	if err != nil {
		log.Printf("[ERROR] Invalid size_quota: %v", err)
	}

	result := cloudappcontrol.WebApplicationRules{
		ID:                  id,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Type:                d.Get("type").(string),
		Order:               d.Get("order").(int),
		State:               d.Get("state").(string),
		Rank:                d.Get("rank").(int),
		TimeQuota:           d.Get("time_quota").(int),
		SizeQuota:           sizeQuotaKB,
		ValidityStartTime:   validityStartTime,
		ValidityEndTime:     validityEndTime,
		ValidityTimeZoneID:  d.Get("validity_time_zone_id").(string),
		EnforceTimeValidity: d.Get("enforce_time_validity").(bool),
		Actions:             SetToStringList(d, "actions"),
		Applications:        SetToStringList(d, "applications"),
		UserRiskScoreLevels: SetToStringList(d, "user_risk_score_levels"),
		DeviceTrustLevels:   SetToStringList(d, "device_trust_levels"),
		UserAgentTypes:      SetToStringList(d, "user_agent_types"),
		Locations:           expandIDNameExtensionsSet(d, "locations"),
		Groups:              expandIDNameExtensionsSet(d, "groups"),
		Departments:         expandIDNameExtensionsSet(d, "departments"),
		Users:               expandIDNameExtensionsSet(d, "users"),
		TimeWindows:         expandIDNameExtensionsSet(d, "time_windows"),
		LocationGroups:      expandIDNameExtensionsSet(d, "location_groups"),
		Labels:              expandIDNameExtensionsSet(d, "labels"),
		DeviceGroups:        expandIDNameExtensionsSet(d, "device_groups"),
		Devices:             expandIDNameExtensionsSet(d, "devices"),
		CBIProfile:          expandCloudAppControlCBIProfile(d),
	}

	return result
}

func expandCloudAppControlCBIProfile(d *schema.ResourceData) cloudappcontrol.CBIProfile {
	if v, ok := d.GetOk("cbi_profile"); ok {
		cbiProfileList := v.([]interface{})
		if len(cbiProfileList) > 0 {
			cbiProfileData := cbiProfileList[0].(map[string]interface{})
			return cloudappcontrol.CBIProfile{
				ID:   cbiProfileData["id"].(string),
				Name: cbiProfileData["name"].(string),
				URL:  cbiProfileData["url"].(string),
			}
		}
	}
	return cloudappcontrol.CBIProfile{}
}

func flattenCloudAppControlCBIProfileSimple(cbiProfile *cloudappcontrol.CBIProfile) []interface{} {
	if cbiProfile == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"id":   cbiProfile.ID,
			"name": cbiProfile.Name,
			"url":  cbiProfile.URL,
		},
	}
}
