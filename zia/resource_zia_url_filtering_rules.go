package zia

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/urlfilteringpolicies"
)

var (
	urlFilteringLock          sync.Mutex
	urlFilteringStartingOrder int
)

func resourceURLFilteringRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceURLFilteringRulesCreate,
		Read:   resourceURLFilteringRulesRead,
		Update: resourceURLFilteringRulesUpdate,
		Delete: resourceURLFilteringRulesDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			action := d.Get("action").(string)

			switch action {
			case "ISOLATE":
				// Validation for ISOLATE action
				cbiProfile, ok := d.GetOk("cbi_profile")
				if !ok || len(cbiProfile.([]interface{})) == 0 {
					return errors.New("cbi_profile attribute is required when action is ISOLATE")
				}

				cbiProfileMap := cbiProfile.([]interface{})[0].(map[string]interface{})
				if cbiProfileMap["id"] == "" && cbiProfileMap["name"] == "" && cbiProfileMap["url"] == "" {
					return errors.New("cbi_profile attribute is required when action is ISOLATE")
				}

				userAgentTypes := d.Get("user_agent_types").(*schema.Set).List()
				for _, userAgent := range userAgentTypes {
					if userAgent.(string) == "OTHER" {
						return errors.New("user_agent_types should not contain 'OTHER' when action is ISOLATE. Valid options are: FIREFOX, MSIE, MSEDGE, CHROME, SAFARI, MSCHREDGE")
					}
				}

				validProtocols := map[string]bool{"HTTPS_RULE": true, "HTTP_RULE": true}
				protocols := d.Get("protocols").(*schema.Set).List()
				for _, protocol := range protocols {
					if !validProtocols[strings.ToUpper(protocol.(string))] {
						return errors.New("when action is ISOLATE, valid options for protocols are: HTTP and/or HTTPS")
					}
				}

			case "CAUTION":
				// Validation for CAUTION action
				validMethods := map[string]bool{"CONNECT": true, "GET": true, "HEAD": true}
				requestMethods, ok := d.GetOk("request_methods")
				if !ok {
					return errors.New("request_methods attribute is required when action is CAUTION")
				}
				requestMethodsList := requestMethods.(*schema.Set).List()
				for _, method := range requestMethodsList {
					if !validMethods[strings.ToUpper(method.(string))] {
						return errors.New("'CAUTION' action is allowed only for CONNECT/GET/HEAD request methods")
					}
				}

				// Ensure BLOCK specific attributes are not set
				if blockOverride, blockOverrideOk := d.GetOk("block_override"); blockOverrideOk && blockOverride.(bool) {
					return errors.New("block_override can only be set when action is BLOCK")
				}

				if _, overrideUsersOk := d.GetOk("override_users"); overrideUsersOk {
					return errors.New("override_users can only be set when action is BLOCK and block_override is true")
				}

				if _, overrideGroupsOk := d.GetOk("override_groups"); overrideGroupsOk {
					return errors.New("override_groups can only be set when action is BLOCK and block_override is true")
				}

			case "BLOCK":
				// Validation for BLOCK action
				blockOverride, blockOverrideOk := d.GetOk("block_override")
				if blockOverrideOk && blockOverride.(bool) {
					overrideUsers, overrideUsersOk := d.GetOk("override_users")
					if !overrideUsersOk || len(overrideUsers.(*schema.Set).List()) == 0 {
						return errors.New("override_users must be set when block_override is true")
					}

					overrideGroups, overrideGroupsOk := d.GetOk("override_groups")
					if !overrideGroupsOk || len(overrideGroups.(*schema.Set).List()) == 0 {
						return errors.New("override_groups must be set when block_override is true")
					}
				} else {
					if _, overrideUsersOk := d.GetOk("override_users"); overrideUsersOk {
						return errors.New("override_users can only be set when block_override is true")
					}

					if _, overrideGroupsOk := d.GetOk("override_groups"); overrideGroupsOk {
						return errors.New("override_groups can only be set when block_override is true")
					}
				}
			}
			// Validate enforce_time_validity, validity_start_time, and validity_end_time
			if enforceTimeValidity, ok := d.GetOk("enforce_time_validity"); ok && enforceTimeValidity.(bool) {
				if _, ok := d.GetOk("validity_start_time"); !ok {
					return errors.New("validity_start_time must be set when enforce_time_validity is true")
				}
				if _, ok := d.GetOk("validity_end_time"); !ok {
					return errors.New("validity_end_time must be set when enforce_time_validity is true")
				}
			}

			return nil
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)
				service := zClient.urlfilteringpolicies

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("rule_id", idInt)
				} else {
					resp, err := urlfilteringpolicies.GetByName(service, id)
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
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
				Description:  "Additional information about the URL Filtering rule",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
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
				ValidateFunc: validation.IntBetween(0, 7),
				Description:  "Admin rank of the Firewall Filtering policy rule",
			},
			"end_user_notification_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of end user notification page to be displayed when the rule is matched. Not applicable if either 'overrideUsers' or 'overrideGroups' is specified.",
			},
			"block_override": {
				Type:     schema.TypeBool,
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "If enforceTimeValidity is set to true, the URL Filtering rule is valid starting on this date and time.",
				ValidateFunc: validateTimesNotInPast,
			},
			"validity_end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "If enforceTimeValidity is set to true, the URL Filtering rule ceases to be valid on this end date and time.",
				ValidateFunc: validateTimesNotInPast,
			},
			"validity_time_zone_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateTimeZone,
				Description:  "If enforceTimeValidity is set to true, the URL Filtering rule date and time is valid based on this time zone ID. Use IANA Format TimeZone.",
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
				Description: "If set to true, the CIPA Compliance rule is enabled",
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
			"locations":              setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of locations for which rule must be applied"),
			"groups":                 setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of groups for which rule must be applied"),
			"departments":            setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of departments for which rule must be applied"),
			"users":                  setIDsSchemaTypeCustom(intPtr(4), "Name-ID pairs of users for which rule must be applied"),
			"time_windows":           setIDsSchemaTypeCustom(nil, "Name-ID pairs of time interval during which rule must be enforced."),
			"override_users":         setIDsSchemaTypeCustom(intPtr(4), "Name-ID pairs of users for which this rule can be overridden."),
			"override_groups":        setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of groups for which this rule can be overridden."),
			"device_groups":          setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":                setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"location_groups":        setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the location groups to which the rule must be applied."),
			"labels":                 setIDsSchemaTypeCustom(nil, "The URL Filtering rule's label."),
			"source_ip_groups":       setIDsSchemaTypeCustom(nil, "list of source ip groups"),
			"workload_groups":        setIdNameSchemaCustom(255, "The list of preconfigured workload groups to which the policy must be applied"),
			"device_trust_levels":    getDeviceTrustLevels(),
			"user_risk_score_levels": getUserRiskScoreLevels(),
			"url_categories":         getURLCategories(),
			"request_methods":        getURLRequestMethods(),
			"protocols":              getURLProtocols(),
			"user_agent_types":       getUserAgentTypes(),
		},
	}
}

func currentOrderVsRankWording(zClient *Client) string {
	service := zClient.urlfilteringpolicies

	list, err := urlfilteringpolicies.GetAll(service)
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
	service := zClient.urlfilteringpolicies

	req := expandURLFilteringRules(d)
	log.Printf("[INFO] Creating url filtering rule\n%+v\n", req)

	// Validate URL Filtering Actions
	if err := validateURLFilteringActions(req); err != nil {
		return err
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		urlFilteringLock.Lock()
		if urlFilteringStartingOrder == 0 {
			list, _ := urlfilteringpolicies.GetAll(service)
			for _, r := range list {
				if r.Order > urlFilteringStartingOrder {
					urlFilteringStartingOrder = r.Order
				}
			}
			if urlFilteringStartingOrder == 0 {
				urlFilteringStartingOrder = 1
			}
		}
		urlFilteringLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = urlFilteringStartingOrder
		resp, err := urlfilteringpolicies.Create(service, &req)
		if err != nil {
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") && reg.MatchString(err.Error()) {
				return fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, order, req.Rank, currentOrderVsRankWording(zClient), err)
			}

			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Creating url filtering rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return fmt.Errorf("error creating resource: %s", err)
		}

		log.Printf("[INFO] Created url filtering rule request. took:%s, without locking:%s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorder(order, resp.ID, "url_filtering_rules", func() (int, error) {
			list, err := urlfilteringpolicies.GetAll(service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := urlfilteringpolicies.Get(service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, _, err = urlfilteringpolicies.Update(service, id, rule)
			return err
		})

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		err = resourceURLFilteringRulesRead(d, m)
		if err != nil {
			if time.Since(start) < timeout {
				time.Sleep(5 * time.Second) // Wait before retrying
				continue
			}
			return err
		}
		markOrderRuleAsDone(resp.ID, "url_filtering_rules")
		break
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

	return nil
}

func resourceURLFilteringRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.urlfilteringpolicies

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return fmt.Errorf("no url filtering rule id is set")
	}
	resp, err := urlfilteringpolicies.Get(service, id)
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
	_ = d.Set("order", resp.Order)
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
	_ = d.Set("user_risk_score_levels", resp.UserRiskScoreLevels)
	_ = d.Set("end_user_notification_url", resp.EndUserNotificationURL)
	_ = d.Set("block_override", resp.BlockOverride)
	_ = d.Set("time_quota", resp.TimeQuota)
	// Convert size_quota from KB back to MB
	sizeQuotaMB := resp.SizeQuota / 1024
	_ = d.Set("size_quota", sizeQuotaMB)
	_ = d.Set("request_methods", resp.RequestMethods)

	// Convert epoch time back to RFC1123 in UTC
	_ = d.Set("validity_start_time", time.Unix(int64(resp.ValidityStartTime), 0).UTC().Format(time.RFC1123))
	_ = d.Set("validity_end_time", time.Unix(int64(resp.ValidityEndTime), 0).UTC().Format(time.RFC1123))

	_ = d.Set("validity_time_zone_id", resp.ValidityTimeZoneID)
	_ = d.Set("enforce_time_validity", resp.EnforceTimeValidity)
	_ = d.Set("action", resp.Action)
	_ = d.Set("ciparule", resp.Ciparule)

	// Update the cbi_profile block in the state
	if resp.CBIProfile.ID != "" {
		if err := d.Set("cbi_profile", flattenCBIProfileSimple(&resp.CBIProfile)); err != nil {
			return err
		}
	}

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
	if err := d.Set("source_ip_groups", flattenIDs(resp.SourceIPGroups)); err != nil {
		return err
	}
	if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
		return fmt.Errorf("error setting workload_groups: %s", err)
	}
	return nil
}

func resourceURLFilteringRulesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.urlfilteringpolicies

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] url filtering rule ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating url filtering rule ID: %v\n", id)
	req := expandURLFilteringRules(d)

	// Validate URL Filtering Actions
	if err := validateURLFilteringActions(req); err != nil {
		return err
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, _, err := urlfilteringpolicies.Update(service, id, &req)
		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating url filtering rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(5 * time.Second) // Wait before retrying
					continue
				}
			}
			return fmt.Errorf("error updating resource: %s", err)
		}

		reorder(req.Order, req.ID, "url_filtering_rules", func() (int, error) {
			list, err := urlfilteringpolicies.GetAll(service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := urlfilteringpolicies.Get(service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, _, err = urlfilteringpolicies.Update(service, id, rule)
			return err
		})

		err = resourceURLFilteringRulesRead(d, m)
		if err != nil {
			if time.Since(start) < timeout {
				time.Sleep(5 * time.Second) // Wait before retrying
				continue
			}
			return err
		}
		markOrderRuleAsDone(req.ID, "url_filtering_rules")
		break
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

	return nil
}

func resourceURLFilteringRulesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.urlfilteringpolicies

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] url filtering rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting url filtering rule ID: %v\n", (d.Id()))

	if _, err := urlfilteringpolicies.Delete(service, id); err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[INFO] url filtering rule deleted")
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

func expandURLFilteringRules(d *schema.ResourceData) urlfilteringpolicies.URLFilteringRule {
	id, _ := getIntFromResourceData(d, "rule_id")

	validityStartTimeStr := d.Get("validity_start_time").(string)
	validityEndTimeStr := d.Get("validity_end_time").(string)

	validityStartTime, err := ConvertRFC1123ToEpoch(validityStartTimeStr)
	if err != nil {
		log.Printf("[ERROR] Invalid validity_start_time: %v", err)
	}
	validityEndTime, err := ConvertRFC1123ToEpoch(validityEndTimeStr)
	if err != nil {
		log.Printf("[ERROR] Invalid validity_end_time: %v", err)
	}

	sizeQuotaMB := d.Get("size_quota").(int)
	sizeQuotaKB, err := convertAndValidateSizeQuota(sizeQuotaMB)
	if err != nil {
		log.Printf("[ERROR] Invalid size_quota: %v", err)
	}

	result := urlfilteringpolicies.URLFilteringRule{
		ID:                     id,
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		Order:                  d.Get("order").(int),
		Protocols:              SetToStringList(d, "protocols"),
		URLCategories:          SetToStringList(d, "url_categories"),
		UserRiskScoreLevels:    SetToStringList(d, "user_risk_score_levels"),
		DeviceTrustLevels:      SetToStringList(d, "device_trust_levels"),
		RequestMethods:         SetToStringList(d, "request_methods"),
		UserAgentTypes:         SetToStringList(d, "user_agent_types"),
		State:                  d.Get("state").(string),
		Rank:                   d.Get("rank").(int),
		EndUserNotificationURL: d.Get("end_user_notification_url").(string),
		BlockOverride:          d.Get("block_override").(bool),
		TimeQuota:              d.Get("time_quota").(int),
		SizeQuota:              sizeQuotaKB,
		ValidityStartTime:      validityStartTime,
		ValidityEndTime:        validityEndTime,
		ValidityTimeZoneID:     d.Get("validity_time_zone_id").(string),
		EnforceTimeValidity:    d.Get("enforce_time_validity").(bool),
		Action:                 d.Get("action").(string),
		Ciparule:               d.Get("ciparule").(bool),
		Locations:              expandIDNameExtensionsSet(d, "locations"),
		Groups:                 expandIDNameExtensionsSet(d, "groups"),
		Departments:            expandIDNameExtensionsSet(d, "departments"),
		Users:                  expandIDNameExtensionsSet(d, "users"),
		TimeWindows:            expandIDNameExtensionsSet(d, "time_windows"),
		OverrideUsers:          expandIDNameExtensionsSet(d, "override_users"),
		OverrideGroups:         expandIDNameExtensionsSet(d, "override_groups"),
		LocationGroups:         expandIDNameExtensionsSet(d, "location_groups"),
		Labels:                 expandIDNameExtensionsSet(d, "labels"),
		DeviceGroups:           expandIDNameExtensionsSet(d, "device_groups"),
		Devices:                expandIDNameExtensionsSet(d, "devices"),
		SourceIPGroups:         expandIDNameExtensionsSet(d, "source_ip_groups"),
		WorkloadGroups:         expandWorkloadGroups(d, "workload_groups"),
		CBIProfile:             expandCBIProfile(d),
	}

	return result
}

func expandCBIProfile(d *schema.ResourceData) urlfilteringpolicies.CBIProfile {
	if v, ok := d.GetOk("cbi_profile"); ok {
		cbiProfileList := v.([]interface{})
		if len(cbiProfileList) > 0 {
			cbiProfileData := cbiProfileList[0].(map[string]interface{})
			return urlfilteringpolicies.CBIProfile{
				ID:   cbiProfileData["id"].(string),
				Name: cbiProfileData["name"].(string),
				URL:  cbiProfileData["url"].(string),
			}
		}
	}
	return urlfilteringpolicies.CBIProfile{}
}

func flattenCBIProfileSimple(cbiProfile *urlfilteringpolicies.CBIProfile) []interface{} {
	if cbiProfile == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"id":   cbiProfile.ID,
			"name": cbiProfile.Name,
			"url":  cbiProfile.URL,
			// "profile_seq": cbiProfile.ProfileSeq,
		},
	}
}

func validateURLFilteringActions(rule urlfilteringpolicies.URLFilteringRule) error {
	switch rule.Action {
	case "ISOLATE":
		// Validation 1: Check if any field in CBIProfile is set
		if rule.CBIProfile.ID == "" && rule.CBIProfile.Name == "" && rule.CBIProfile.URL == "" {
			return errors.New("cbi_profile attribute is required when action is ISOLATE")
		}

		// Validation 2: Check user_agent_types does not contain "OTHER"
		for _, userAgent := range rule.UserAgentTypes {
			if userAgent == "OTHER" {
				return errors.New("user_agent_types should not contain 'OTHER' when action is ISOLATE. Valid options are: FIREFOX, MSIE, MSEDGE, CHROME, SAFARI, MSCHREDGE")
			}
		}

		// Validation 3: Check Protocols should be HTTP or HTTPS
		validProtocols := map[string]bool{"HTTPS_RULE": true, "HTTP_RULE": true}
		for _, protocol := range rule.Protocols {
			if !validProtocols[strings.ToUpper(protocol)] {
				return errors.New("when action is ISOLATE, valid options for protocols are: HTTP and/or HTTPS")
			}
		}

	case "CAUTION":
		// Validation 4: Ensure request_methods only contain CONNECT, GET, HEAD
		validMethods := map[string]bool{"CONNECT": true, "GET": true, "HEAD": true}
		for _, method := range rule.RequestMethods { // Assuming RequestMethods is the correct field
			if !validMethods[strings.ToUpper(method)] {
				return errors.New("'CAUTION' action is allowed only for CONNECT/GET/HEAD request methods")
			}
		}
	}

	return nil
}

func validateTimeZone(v interface{}, k string) (ws []string, errors []error) {
	tzStr := v.(string)
	_, err := time.LoadLocation(tzStr)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q is not a valid timezone. Visit https://nodatime.org/TimeZones for the valid IANA list", tzStr))
	}

	return
}

func ConvertRFC1123ToEpoch(timeStr string) (int, error) {
	t, err := time.Parse(time.RFC1123, timeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid time format: %v. Expected format: RFC1123 (Mon, 02 Jan 2006 15:04:05 MST)", err)
	}
	return int(t.Unix()), nil
}

func validateTimesNotInPast(val interface{}, key string) (warns []string, errs []error) {
	timeStr := val.(string)
	timeVal, err := ConvertRFC1123ToEpoch(timeStr)
	if err != nil {
		errs = append(errs, fmt.Errorf("%q: invalid time format: %v. Expected format: RFC1123 (Mon, 02 Jan 2006 15:04:05 MST)", key, err))
		return
	}

	now := time.Now().Unix()

	if int64(timeVal) < now {
		errs = append(errs, fmt.Errorf("%q: time cannot be in the past", key))
	}

	return
}

func convertAndValidateSizeQuota(sizeQuotaMB int) (int, error) {
	const (
		minMB = 10
		maxMB = 100000
	)
	if sizeQuotaMB < minMB || sizeQuotaMB > maxMB {
		return 0, fmt.Errorf("size_quota must be between %d MB and %d MB", minMB, maxMB)
	}
	// Convert MB to KB
	sizeQuotaKB := sizeQuotaMB * 1024
	return sizeQuotaKB, nil
}
