package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/outbound_email_dlp"
)

var (
	outboundEmailDlpLock                  sync.Mutex
	outboundEmailDlpStartingOrder         int
	outboundEmailDlpSubRulesStartingOrder = make(map[int]int)
)

func resourceOutboundEmailDLP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOutboundEmailDLPCreate,
		ReadContext:   resourceOutboundEmailDLPRead,
		UpdateContext: resourceOutboundEmailDLPUpdate,
		DeleteContext: resourceOutboundEmailDLPDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			externalEmail, emailSet := d.GetOk("external_auditor_email")
			auditorRaw, auditorSet := d.GetOk("auditor")
			nt, ntSet := d.GetOk("notification_template")

			isExternalSet := emailSet && externalEmail != ""
			isAuditorSet := auditorSet && auditorRaw != ""
			isNTSet := ntSet && nt != ""

			// external_auditor_email and auditor are mutually exclusive.
			if isExternalSet && isAuditorSet {
				return fmt.Errorf("'external_auditor_email' and 'auditor' cannot both be set")
			}

			// When external_auditor_email is set, notification_template must also be set.
			if isExternalSet {
				if !isNTSet {
					return fmt.Errorf("when setting 'external_auditor_email', 'notification_template' must also be set")
				}
				return nil
			}

			// When external_auditor_email is not set, auditor and notification_template
			// must be set together (or neither).
			if isAuditorSet || isNTSet {
				if !isAuditorSet || !isNTSet {
					return fmt.Errorf("when 'external_auditor_email' is not set, both 'auditor' and 'notification_template' must be set")
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
					resp, err := outbound_email_dlp.GetByName(ctx, service, id)
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DLP policy rule name.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				Description:      "The description of the DLP policy rule.",
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"order": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The rule order of execution for the DLP policy rule with respect to other rules.",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enables or disables the DLP policy rule.",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
					"ENABLED_AUTHENTICATION_REQUIRED",
				}, false),
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The action taken when traffic matches the DLP policy rule criteria.",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"BLOCK",
					"QUARANTINE",
				}, false),
			},
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates the severity selected for the DLP rule violation.",
				ValidateFunc: validation.StringInSlice([]string{
					"RULE_SEVERITY_HIGH",
					"RULE_SEVERITY_MEDIUM",
					"RULE_SEVERITY_LOW",
					"RULE_SEVERITY_INFO",
				}, false),
			},
			"file_types": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of file types to which the DLP policy rule must be applied.",
			},
			"content_locations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of content locations to which the DLP policy rule must be applied.",
			},
			"min_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The minimum file size (in KB) used for evaluation of the DLP policy rule.",
			},
			"without_content_inspection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"external_auditor_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address of an external auditor to whom DLP email notifications are sent.",
			},
			"custom_header": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The custom header value inserted when the rule action includes custom header insertion.",
			},
			"parent_rule": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the parent rule under which an exception rule is added.",
			},
			"sub_rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of exception rules added to a parent rule.",
			},
			"user_risk_score_levels":   getUserRiskScoreLevels(),
			"groups":                   setIDsSchemaTypeCustom(nil, "The Name-ID pairs of groups to which the DLP policy rule must be applied"),
			"departments":              setIDsSchemaTypeCustom(nil, "The Name-ID pairs of departments to which the DLP policy rule must be applied"),
			"users":                    setIDsSchemaTypeCustom(nil, "The Name-ID pairs of users to which the DLP policy rule must be applied"),
			"excluded_groups":          setIDsSchemaTypeCustom(nil, "The Name-ID pairs of groups that are excluded from the DLP policy rule"),
			"excluded_departments":     setIDsSchemaTypeCustom(nil, "The Name-ID pairs of departments that are excluded from the DLP policy rule"),
			"excluded_users":           setIDsSchemaTypeCustom(nil, "The Name-ID pairs of users that are excluded from the DLP policy rule"),
			"time_windows":             setIDsSchemaTypeCustom(nil, "The Name-ID pairs of time windows to which the DLP policy rule must be applied"),
			"dlp_engines":              setIDsSchemaTypeCustom(nil, "The list of DLP engines to which the DLP policy rule must be applied"),
			"labels":                   setIDsSchemaTypeCustom(nil, "The Name-ID pairs of rule labels associated to the DLP policy rule"),
			"included_domain_profiles": setIDsSchemaTypeCustom(nil, "The Name-ID pairs of domain profiles included in the DLP policy rule"),
			"email_tenants":            setIDsSchemaTypeCustom(nil, "The Name-ID pairs of email tenants to which the DLP policy rule must be applied"),
			"email_recipient_profiles": setIDsSchemaTypeCustom(nil, "The Name-ID pairs of email recipient profiles to which the DLP policy rule must be applied"),
			"auditor":                  setSingleIDSchemaTypeCustom("The auditor to which the DLP policy rule must be applied"),
			"notification_template":    setSingleIDSchemaTypeCustom("The template used for DLP notification emails"),
			"receiver": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "The Zscaler Incident Receiver associated with the DLP policy rule. Only the 'id' must be set; the receiver 'type' is always 'ZIR'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier of the Zscaler Incident Receiver.",
						},
					},
				},
			},
		},
	}
}

func resourceOutboundEmailDLPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandOutboundEmailDLP(d)

	log.Printf("[INFO] Creating zia outbound email dlp rule\n%+v\n", req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()
	isSubRule := req.ParentRule != 0
	for {
		outboundEmailDlpLock.Lock()
		if isSubRule {
			if outboundEmailDlpSubRulesStartingOrder[req.ParentRule] == 0 {
				list, _ := outbound_email_dlp.Get(ctx, service, req.ParentRule)
				if list != nil {
					for _, subRule := range list.SubRules {
						r, err := outbound_email_dlp.Get(ctx, service, subRule.ID)
						if err != nil {
							outboundEmailDlpLock.Unlock()
							return diag.FromErr(err)
						}
						if r.Order > outboundEmailDlpSubRulesStartingOrder[req.ParentRule] {
							outboundEmailDlpSubRulesStartingOrder[req.ParentRule] = r.Order
						}
					}
				}
				if outboundEmailDlpSubRulesStartingOrder[req.ParentRule] == 0 {
					outboundEmailDlpSubRulesStartingOrder[req.ParentRule] = 1
				}
			}
		} else {
			if outboundEmailDlpStartingOrder == 0 {
				list, _ := outbound_email_dlp.GetAll(ctx, service, nil)
				for _, r := range list {
					if r.Order > outboundEmailDlpStartingOrder {
						outboundEmailDlpStartingOrder = r.Order
					}
				}
				if outboundEmailDlpStartingOrder == 0 {
					outboundEmailDlpStartingOrder = 1
				}
			}
		}
		outboundEmailDlpLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		if isSubRule {
			req.Order = outboundEmailDlpSubRulesStartingOrder[req.ParentRule]
		} else {
			req.Order = outboundEmailDlpStartingOrder
		}

		resp, _, err := outbound_email_dlp.Create(ctx, service, &req)

		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if time.Since(start) < timeout {
					time.Sleep(5 * time.Second)
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia outbound email dlp rule request. Took: %s, without locking: %s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)

		resourceType := "outbound_email_dlp"
		if isSubRule {
			resourceType = fmt.Sprintf("outbound_email_dlp_sub_%d", req.ParentRule)
		}

		reorderWithBeforeReorder(
			OrderRule{Order: order, Rank: 0},
			resp.ID,
			resourceType,
			func() (map[int]OrderRule, error) {
				if isSubRule {
					parent, err := outbound_email_dlp.Get(ctx, service, req.ParentRule)
					if err != nil {
						return nil, err
					}
					m := make(map[int]OrderRule, len(parent.SubRules))
					for _, r := range parent.SubRules {
						m[r.ID] = OrderRule{Order: r.Order, Rank: 0}
					}
					return m, nil
				}
				allRules, err := outbound_email_dlp.GetAll(ctx, service, nil)
				if err != nil {
					return nil, err
				}
				m := make(map[int]OrderRule, len(allRules))
				for _, r := range allRules {
					m[r.ID] = OrderRule{Order: r.Order, Rank: 0}
				}
				return m, nil
			},
			func(id int, order OrderRule) error {
				rule, err := outbound_email_dlp.Get(ctx, service, id)
				if err != nil {
					return err
				}
				if rule.Order == order.Order {
					return nil
				}

				rule.LastModifiedTime = 0
				rule.LastModifiedBy = nil
				rule.Order = order.Order

				_, _, err = outbound_email_dlp.Update(ctx, service, id, rule)
				if err != nil {
					log.Printf("[ERROR] Failed to update order for rule ID %d: %v", id, err)
				}
				return err
			},
			nil,
		)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceOutboundEmailDLPRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second)
				continue
			}
			return diags
		}

		markOrderRuleAsDone(resp.ID, resourceType)
		waitForReorder(resourceType)
		break
	}

	time.Sleep(2 * time.Second)

	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func resourceOutboundEmailDLPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia outbound email dlp rule id is set"))
	}
	resp, err := outbound_email_dlp.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing outbound email dlp rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting outbound email dlp rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("order", resp.Order)
	_ = d.Set("description", resp.Description)
	_ = d.Set("state", resp.State)
	_ = d.Set("action", resp.Action)
	_ = d.Set("severity", resp.Severity)
	_ = d.Set("file_types", resp.FileTypes)
	_ = d.Set("content_locations", resp.ContentLocations)
	_ = d.Set("user_risk_score_levels", resp.UserRiskScoreLevels)
	_ = d.Set("min_size", resp.MinSize)
	_ = d.Set("without_content_inspection", resp.WithoutContentInspection)
	_ = d.Set("external_auditor_email", resp.ExternalAuditorEmail)
	_ = d.Set("custom_header", resp.CustomHeader)
	_ = d.Set("parent_rule", resp.ParentRule)

	subRuleIDs := make([]interface{}, len(resp.SubRules))
	for i, subRule := range resp.SubRules {
		subRuleIDs[i] = strconv.Itoa(subRule.ID)
	}
	if err := d.Set("sub_rules", subRuleIDs); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sub_rules: %s", err))
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
	if err := d.Set("excluded_groups", flattenIDExtensionsListIDs(resp.ExcludedGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("excluded_departments", flattenIDExtensionsListIDs(resp.ExcludedDepartments)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("excluded_users", flattenIDExtensionsListIDs(resp.ExcludedUsers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("time_windows", flattenIDExtensionsListIDs(resp.TimeWindows)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dlp_engines", flattenIDExtensionsListIDs(resp.DlpEngines)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("included_domain_profiles", flattenIDExtensionsListIDs(resp.IncludedDomainProfiles)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_tenants", flattenIDExtensionsListIDs(resp.EmailTenants)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_recipient_profiles", flattenIDExtensionsListIDs(resp.EmailRecipientProfiles)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("auditor", flattenSingleIDNameExtensions(resp.Auditor)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("notification_template", flattenSingleIDNameExtensions(resp.NotificationTemplate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("receiver", flattenEndpointReceiver(resp.Receiver)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOutboundEmailDLPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] outbound email dlp rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("outbound email dlp rule ID not set"))
	}

	req := expandOutboundEmailDLP(d)

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, _, err := outbound_email_dlp.Update(ctx, service, id, &req)

		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			log.Printf("[INFO] Retrying due to API error: %s", err)
			if time.Since(start) < timeout {
				time.Sleep(5 * time.Second)
				continue
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		break
	}

	isSubRule := req.ParentRule != 0
	resourceType := "outbound_email_dlp"
	if isSubRule {
		resourceType = fmt.Sprintf("outbound_email_dlp_sub_%d", req.ParentRule)
	}

	reorderWithBeforeReorder(OrderRule{Order: req.Order, Rank: 0}, id, resourceType,
		func() (map[int]OrderRule, error) {
			if isSubRule {
				parent, err := outbound_email_dlp.Get(ctx, service, req.ParentRule)
				if err != nil {
					return nil, err
				}
				m := make(map[int]OrderRule, len(parent.SubRules))
				for _, r := range parent.SubRules {
					m[r.ID] = OrderRule{Order: r.Order, Rank: 0}
				}
				return m, nil
			}
			allRules, err := outbound_email_dlp.GetAll(ctx, service, nil)
			if err != nil {
				return nil, err
			}
			m := make(map[int]OrderRule, len(allRules))
			for _, r := range allRules {
				m[r.ID] = OrderRule{Order: r.Order, Rank: 0}
			}
			return m, nil
		},
		func(id int, order OrderRule) error {
			rule, err := outbound_email_dlp.Get(ctx, service, id)
			if err != nil {
				return err
			}
			if rule.Order == order.Order {
				return nil
			}

			rule.LastModifiedTime = 0
			rule.LastModifiedBy = nil
			rule.Order = order.Order

			_, _, err = outbound_email_dlp.Update(ctx, service, id, rule)
			if err != nil {
				log.Printf("[ERROR] Failed to update order for rule ID %d: %v", id, err)
			}
			return err
		},
		nil,
	)

	markOrderRuleAsDone(id, resourceType)
	waitForReorder(resourceType)

	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	}

	return resourceOutboundEmailDLPRead(ctx, d, meta)
}

func resourceOutboundEmailDLPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] outbound email dlp rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting outbound email dlp rule ID: %v\n", (d.Id()))

	if _, err := outbound_email_dlp.Delete(ctx, service, id); err != nil {
		if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
			log.Printf("[INFO] outbound email dlp rule %d not found, skipping deletion", id)
			return nil
		}
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] outbound email dlp rule deleted")

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

// expandOutboundEmailSubRules builds the list of child rules referenced by a
// parent rule from the "sub_rules" set of string IDs. Only the ID is sent; the
// API resolves the rest of each child rule from its own definition.
func expandOutboundEmailSubRules(set *schema.Set) []outbound_email_dlp.OutboundEmailDlp {
	if set == nil {
		return nil
	}
	var subRules []outbound_email_dlp.OutboundEmailDlp
	for _, item := range set.List() {
		subRuleID, err := strconv.Atoi(item.(string))
		if err == nil {
			subRules = append(subRules, outbound_email_dlp.OutboundEmailDlp{ID: subRuleID})
		}
	}
	return subRules
}

func expandOutboundEmailDLP(d *schema.ResourceData) outbound_email_dlp.OutboundEmailDlp {
	id, _ := getIntFromResourceData(d, "rule_id")

	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandOutboundEmailDLP: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	result := outbound_email_dlp.OutboundEmailDlp{
		ID:                       id,
		Name:                     d.Get("name").(string),
		Order:                    order,
		Description:              d.Get("description").(string),
		State:                    d.Get("state").(string),
		Action:                   d.Get("action").(string),
		Severity:                 d.Get("severity").(string),
		MinSize:                  d.Get("min_size").(int),
		WithoutContentInspection: d.Get("without_content_inspection").(bool),
		ExternalAuditorEmail:     d.Get("external_auditor_email").(string),
		CustomHeader:             d.Get("custom_header").(string),
		ParentRule:               d.Get("parent_rule").(int),
		SubRules:                 expandOutboundEmailSubRules(d.Get("sub_rules").(*schema.Set)),
		FileTypes:                SetToStringList(d, "file_types"),
		ContentLocations:         SetToStringList(d, "content_locations"),
		UserRiskScoreLevels:      SetToStringList(d, "user_risk_score_levels"),
		Groups:                   expandIDNameExtensionsSet(d, "groups"),
		Departments:              expandIDNameExtensionsSet(d, "departments"),
		Users:                    expandIDNameExtensionsSet(d, "users"),
		ExcludedGroups:           expandIDNameExtensionsSet(d, "excluded_groups"),
		ExcludedDepartments:      expandIDNameExtensionsSet(d, "excluded_departments"),
		ExcludedUsers:            expandIDNameExtensionsSet(d, "excluded_users"),
		TimeWindows:              expandIDNameExtensionsSet(d, "time_windows"),
		DlpEngines:               expandIDNameExtensionsSet(d, "dlp_engines"),
		Labels:                   expandIDNameExtensionsSet(d, "labels"),
		IncludedDomainProfiles:   expandIDNameExtensionsSet(d, "included_domain_profiles"),
		EmailTenants:             expandIDNameExtensionsSet(d, "email_tenants"),
		EmailRecipientProfiles:   expandIDNameExtensionsSet(d, "email_recipient_profiles"),
		Auditor:                  expandSingleIDNameExtensions(d, "auditor"),
		NotificationTemplate:     expandSingleIDNameExtensions(d, "notification_template"),
		Receiver:                 expandEndpointReceiver(d, "receiver"),
	}
	return result
}
