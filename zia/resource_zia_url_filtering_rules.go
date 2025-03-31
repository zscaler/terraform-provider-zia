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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

var (
	urlFilteringLock          sync.Mutex
	urlFilteringStartingOrder int
)

func resourceURLFilteringRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceURLFilteringRulesCreate,
		ReadContext:   resourceURLFilteringRulesRead,
		UpdateContext: resourceURLFilteringRulesUpdate,
		DeleteContext: resourceURLFilteringRulesDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			action := d.Get("action").(string)

			// Common validation for actions other than BLOCK
			if action != "BLOCK" {
				if blockOverride, blockOverrideOk := d.GetOk("block_override"); blockOverrideOk && blockOverride.(bool) {
					return errors.New("block_override can only be set when action is BLOCK")
				}

				if _, overrideUsersOk := d.GetOk("override_users"); overrideUsersOk {
					return errors.New("override_users can only be set when action is BLOCK and block_override is true")
				}

				if _, overrideGroupsOk := d.GetOk("override_groups"); overrideGroupsOk {
					return errors.New("override_groups can only be set when action is BLOCK and block_override is true")
				}
			}

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
						return errors.New("user_agent_types should not contain 'OTHER' when action is ISOLATE. Valid options are: CHROME, FIREFOX, MSIE, MSEDGE, MSCHREDGE, OPERA, SAFARI")
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

			case "BLOCK":
				// Validation for BLOCK action
				blockOverride, blockOverrideOk := d.GetOk("block_override")
				if blockOverrideOk && blockOverride.(bool) {
					// If block_override is true, override_users and override_groups can be set but are optional
					// No further checks needed here as block_override being true allows override_users and override_groups to be set optionally
				} else {
					// If block_override is false, override_users and override_groups must not be set
					if _, overrideUsersOk := d.GetOk("override_users"); overrideUsersOk {
						return errors.New("override_users can only be set when block_override is true")
					}

					if _, overrideGroupsOk := d.GetOk("override_groups"); overrideGroupsOk {
						return errors.New("override_groups can only be set when block_override is true")
					}
				}
			}

			// Validate enforce_time_validity, validity_start_time, validity_end_time, and validity_time_zone_id
			if enforceTimeValidity, ok := d.GetOk("enforce_time_validity"); ok && enforceTimeValidity.(bool) {
				if _, ok := d.GetOk("validity_start_time"); !ok {
					return errors.New("validity_start_time must be set when enforce_time_validity is true")
				} else {
					validityStartTimeStr := d.Get("validity_start_time").(string)
					if isSingleDigitDay(validityStartTimeStr) {
						return errors.New("validity_start_time must have a two-digit day (e.g., 02 instead of 2)")
					}
				}
				if _, ok := d.GetOk("validity_end_time"); !ok {
					return errors.New("validity_end_time must be set when enforce_time_validity is true")
				} else {
					validityEndTimeStr := d.Get("validity_end_time").(string)
					if isSingleDigitDay(validityEndTimeStr) {
						return errors.New("validity_end_time must have a two-digit day (e.g., 02 instead of 2)")
					}
				}
				if _, ok := d.GetOk("validity_time_zone_id"); !ok {
					return errors.New("validity_time_zone_id must be set when enforce_time_validity is true")
				}
			} else {
				// If enforce_time_validity is false, ensure validity attributes are not set
				if _, ok := d.GetOk("validity_start_time"); ok {
					return errors.New("validity_start_time can only be set when enforce_time_validity is true")
				}
				if _, ok := d.GetOk("validity_end_time"); ok {
					return errors.New("validity_end_time can only be set when enforce_time_validity is true")
				}
				if _, ok := d.GetOk("validity_time_zone_id"); ok {
					return errors.New("validity_time_zone_id can only be set when enforce_time_validity is true")
				}
			}

			return nil
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("rule_id", idInt)
				} else {
					resp, err := urlfilteringpolicies.GetByName(ctx, service, id)
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
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Order of execution of rule with respect to other URL Filtering rules",
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

func currentOrderVsRankWording(ctx context.Context, zClient *Client) string {
	service := zClient.Service

	list, err := urlfilteringpolicies.GetAll(ctx, service)
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

func resourceURLFilteringRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandURLFilteringRules(d)
	log.Printf("[INFO] Creating url filtering rule\n%+v\n", req)

	// Validate URL Filtering Actions
	if err := validateURLFilteringActions(req); err != nil {
		return diag.FromErr(err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		urlFilteringLock.Lock()
		if urlFilteringStartingOrder == 0 {
			list, _ := urlfilteringpolicies.GetAll(ctx, service)
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
		resp, err := urlfilteringpolicies.Create(ctx, service, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") && reg.MatchString(err.Error()) {
				return diag.FromErr(fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, order, req.Rank, currentOrderVsRankWording(ctx, zClient), err))
			}

			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Creating url filtering rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created url filtering rule request. took:%s, without locking:%s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorder(order, resp.ID, "url_filtering_rules", func() (int, error) {
			list, err := urlfilteringpolicies.GetAll(ctx, service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := urlfilteringpolicies.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, _, err = urlfilteringpolicies.Update(ctx, service, id, rule)
			return err
		})

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceURLFilteringRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(resp.ID, "url_filtering_rules")
		break
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func resourceURLFilteringRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no url filtering rule id is set"))
	}
	resp, err := urlfilteringpolicies.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia url filtering rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
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
	_ = d.Set("action", resp.Action)
	_ = d.Set("ciparule", resp.Ciparule)

	// Update the cbi_profile block in the state
	if resp.CBIProfile.ID != "" {
		if err := d.Set("cbi_profile", flattenCBIProfileSimple(&resp.CBIProfile)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("groups", flattenIDExtensionsListIDs(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("departments", flattenIDExtensionsListIDs(resp.Departments)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("users", flattenIDExtensionsListIDs(resp.Users)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_windows", flattenIDs(resp.TimeWindows)); err != nil {
		return diag.FromErr(err)
	}

	// Ensure override_users and override_groups are only set when block_override is true and action is BLOCK
	if resp.Action == "BLOCK" && resp.BlockOverride {
		if err := d.Set("override_users", flattenIDExtensionsListIDs(resp.OverrideUsers)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("override_groups", flattenIDExtensionsListIDs(resp.OverrideGroups)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		// Remove override_users and override_groups from state if block_override is not true or action is not BLOCK
		_ = d.Set("override_users", nil)
		_ = d.Set("override_groups", nil)
	}

	if err := d.Set("location_groups", flattenIDExtensionsListIDs(resp.LocationGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("device_groups", flattenIDExtensionsListIDs(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("devices", flattenIDExtensionsListIDs(resp.Devices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_ip_groups", flattenIDExtensionsListIDs(resp.SourceIPGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting workload_groups: %s", err))
	}
	return nil
}

func resourceURLFilteringRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] url filtering rule ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating url filtering rule ID: %v\n", id)
	req := expandURLFilteringRules(d)

	// Validate URL Filtering Actions
	if err := validateURLFilteringActions(req); err != nil {
		return diag.FromErr(err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, _, err := urlfilteringpolicies.Update(ctx, service, id, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating url filtering rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(5 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		reorder(req.Order, req.ID, "url_filtering_rules", func() (int, error) {
			list, err := urlfilteringpolicies.GetAll(ctx, service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := urlfilteringpolicies.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, _, err = urlfilteringpolicies.Update(ctx, service, id, rule)
			return err
		})

		if diags := resourceURLFilteringRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(req.ID, "url_filtering_rules")
		break
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func resourceURLFilteringRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] url filtering rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting url filtering rule ID: %v\n", (d.Id()))

	if _, err := urlfilteringpolicies.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] url filtering rule deleted")
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandURLFilteringRules(d *schema.ResourceData) urlfilteringpolicies.URLFilteringRule {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandSSLInspectionRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

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

	result := urlfilteringpolicies.URLFilteringRule{
		ID:                     id,
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		Order:                  order,
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
		},
	}
}
