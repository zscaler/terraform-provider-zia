package zia

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_dlp_rules"
)

// endpointReceiverType is the fixed receiver type required by the endpoint DLP
// rule API. The user only supplies the receiver "id".
const endpointReceiverType = "ZIR"

var (
	endpointDlpRulesLock     sync.Mutex
	endpointDlpStartingOrder int
)

func resourceEndpointDLPRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointDLPRulesCreate,
		ReadContext:   resourceEndpointDLPRulesRead,
		UpdateContext: resourceEndpointDLPRulesUpdate,
		DeleteContext: resourceEndpointDLPRulesDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// The endpoint application selectors are only valid when the data transfer
			// method targets application file access.
			hasEndpointApps := endpointBlockHasValues(d.Get("end_point_applications"), "zapp_id")
			hasEndpointAppGroups := endpointBlockHasValues(d.Get("end_point_application_groups"), "group_id")
			if hasEndpointApps || hasEndpointAppGroups {
				if d.Get("data_transfer_method").(string) != "APPLICATION_FILE_ACCESS" {
					return fmt.Errorf("'end_point_applications' and 'end_point_application_groups' can only be set when 'data_transfer_method' is 'APPLICATION_FILE_ACCESS'")
				}
			}

			externalEmail, emailSet := d.GetOk("external_auditor_email")
			auditorRaw, auditorSet := d.GetOk("auditor")
			nt, ntSet := d.GetOk("notification_template")

			isExternalSet := emailSet && externalEmail != ""
			isAuditorSet := auditorSet && auditorRaw != ""
			isNTSet := ntSet && nt != ""

			// Rule 3: Mutually exclusive
			if isExternalSet && isAuditorSet {
				return fmt.Errorf("'external_auditor_email' and 'auditor' cannot both be set")
			}

			// Rule 1: If external_auditor_email is set, notification_template must be set
			if isExternalSet {
				if !isNTSet {
					return fmt.Errorf("when setting 'external_auditor_email', 'notification_template' must also be set")
				}
				return nil // valid as long as no conflict and notification_template is present
			}

			// Rule 2: If external_auditor_email is not set, both auditor and notification_template must be set (or neither)
			if isAuditorSet || isNTSet {
				if !isAuditorSet || !isNTSet {
					return fmt.Errorf("when 'external_auditor_email' is not set, both 'auditor' and 'notification_template' must be set")
				}
			}

			// Rule 4: If none are set, it's valid
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
					resp, err := endpoint_dlp_rules.GetByName(ctx, service, id)
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
				Description:      "The description of the DLP policy rule.",
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"rank": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 7),
				Description:  "Admin rank of the admin who creates this rule.",
			},
			"order": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The rule order of execution for the DLP policy rule with respect to other rules.",
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
			"sub_rules": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The IDs of the exception (sub) rules associated with this parent rule. Sub-rules are managed through the 'zia_endpoint_dlp_sub_rules' resource.",
			},
			"user_risk_score_levels": getUserRiskScoreLevels(),
			"device_trust_levels":    getDeviceTrustLevels(),
			"file_types": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of file types to which the DLP policy rule must be applied.",
			},
			"data_transfer_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The data transfer method to which the DLP policy rule must be applied.",
				ValidateFunc: validation.StringInSlice([]string{
					"APPLICATION_FILE_ACCESS",
					"NETWORK_DRIVE_TRANSFER",
					"PRINTING",
					"REMOVABLE_DRIVE_TRANSFER",
				}, false),
			},
			"network_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The network type to which the DLP policy rule must be applied.",
			},
			"min_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 96000),
				Description:  "The minimum file size (in KB) used for evaluation of the DLP policy rule.",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The action taken when traffic matches the DLP policy rule criteria.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"NONE",
					"BLOCK",
					"CONFIRM",
					"ALLOW",
				}, false),
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enables or disables the DLP policy rule.",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"external_auditor_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address of an external auditor to whom DLP email notifications are sent",
			},
			"without_content_inspection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"eun_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the End User Notification (EUN) is enabled for the DLP policy rule.",
			},
			"eun_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The EUN template ID associated with the rule.",
			},
			"uc_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier of the user confirmation template.",
			},
			"device_groups":                setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":                      setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"users":                        setIDsSchemaTypeCustom(nil, "The Name-ID pairs of users to which the DLP policy rule must be applied"),
			"groups":                       setIDsSchemaTypeCustom(nil, "The Name-ID pairs of groups to which the DLP policy rule must be applied"),
			"departments":                  setIDsSchemaTypeCustom(nil, "The Name-ID pairs of departments to which the DLP policy rule must be applied"),
			"resources":                    setIDsSchemaTypeCustom(nil, "The Name-ID pairs of resources to which the DLP policy rule must be applied"),
			"resource_groups":              setIDsSchemaTypeCustom(nil, "The Name-ID pairs of resource groups to which the DLP policy rule must be applied"),
			"dlp_engines":                  setIDsSchemaTypeCustom(nil, "The list of DLP engines to which the DLP policy rule must be applied"),
			"labels":                       setIDsSchemaTypeCustom(nil, "The Name-ID pairs of rule labels associated to the DLP policy rule"),
			"auditor":                      setSingleIDSchemaTypeCustom("The auditor to which the DLP policy rule must be applied"),
			"notification_template":        setSingleIDSchemaTypeCustom("The template used for DLP notification emails"),
			"end_point_applications":       setCustomKeyIDsSchema("zapp_id", "The endpoint applications to which the DLP policy rule must be applied. Can only be set when 'data_transfer_method' is set to 'APPLICATION_FILE_ACCESS'."),
			"end_point_application_groups": setCustomKeyIDsSchema("group_id", "The endpoint application groups to which the DLP policy rule must be applied. Can only be set when 'data_transfer_method' is set to 'APPLICATION_FILE_ACCESS'."),
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

func resourceEndpointDLPRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandEndpointDLPRules(d)

	log.Printf("[INFO] Creating zia endpoint dlp rule\n%+v\n", req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()
	for {
		endpointDlpRulesLock.Lock()
		if endpointDlpStartingOrder == 0 {
			list, _ := endpoint_dlp_rules.GetAll(ctx, service)
			for _, r := range list {
				if r.Order > endpointDlpStartingOrder {
					endpointDlpStartingOrder = r.Order
				}
			}
			if endpointDlpStartingOrder == 0 {
				endpointDlpStartingOrder = 1
			}
		}
		endpointDlpRulesLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = endpointDlpStartingOrder

		resp, _, err := endpoint_dlp_rules.Create(ctx, service, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") && !strings.Contains(err.Error(), "ICAP Receiver with id") {
				if time.Since(start) < timeout {
					time.Sleep(5 * time.Second)
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia endpoint dlp rule request. Took: %s, without locking: %s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)

		resourceType := "endpoint_dlp_rules"

		reorderWithBeforeReorder(
			OrderRule{Order: order, Rank: req.Rank},
			resp.ID,
			resourceType,
			func() (map[int]OrderRule, error) {
				allRules, err := endpoint_dlp_rules.GetAll(ctx, service)
				if err != nil {
					return nil, err
				}
				m := make(map[int]OrderRule, len(allRules))
				for _, r := range allRules {
					m[r.ID] = OrderRule{Order: r.Order, Rank: r.Rank}
				}
				return m, nil
			},
			func(id int, order OrderRule) error {
				rule, err := endpoint_dlp_rules.Get(ctx, service, id)
				if err != nil {
					return err
				}

				rule.LastModifiedTime = 0
				rule.LastModifiedBy = nil
				rule.Order = order.Order
				rule.Rank = order.Rank

				log.Printf("[DEBUG] Updating parent rule ID %d to order %d", id, order.Order)

				_, _, err = endpoint_dlp_rules.Update(ctx, service, id, rule)
				if err != nil {
					log.Printf("[ERROR] Failed to update order for rule ID %d: %v", id, err)
				}
				return err
			},
			nil, // Remove beforeReorder function to avoid adding too many rules to the map
		)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceEndpointDLPRulesRead(ctx, d, meta); diags.HasError() {
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

func resourceEndpointDLPRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia endpoint dlp rule id is set"))
	}
	resp, err := endpoint_dlp_rules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing endpoint dlp rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting endpoint dlp rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("description", resp.Description)
	_ = d.Set("file_types", resp.FileTypes)
	_ = d.Set("user_risk_score_levels", resp.UserRiskScoreLevels)
	_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)
	_ = d.Set("state", resp.State)
	_ = d.Set("min_size", resp.MinSize)
	_ = d.Set("action", resp.Action)
	_ = d.Set("severity", resp.Severity)
	_ = d.Set("data_transfer_method", resp.DataTransferMethod)
	_ = d.Set("network_type", resp.NetworkType)
	_ = d.Set("external_auditor_email", resp.ExternalAuditorEmail)
	_ = d.Set("without_content_inspection", resp.WithoutContentInspection)
	_ = d.Set("eun_enabled", resp.EunEnabled)
	_ = d.Set("eun_template_id", resp.EunTemplateId)
	_ = d.Set("uc_template_id", resp.UcTemplateId)

	// Flatten sub_rules
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

	if err := d.Set("device_groups", flattenIDExtensionsListIDs(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("devices", flattenIDExtensionsListIDs(resp.Devices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("users", flattenIDExtensionsListIDs(resp.Users)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("resources", flattenIDExtensionsListIDs(resp.Resources)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("resource_groups", flattenIDExtensionsListIDs(resp.ResourceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dlp_engines", flattenIDExtensionsListIDs(resp.DlpEngines)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("auditor", flattenSingleIDNameExtensions(resp.Auditor)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("notification_template", flattenSingleIDNameExtensions(resp.NotificationTemplate)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("receiver", flattenEndpointReceiver(resp.Receiver)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting receiver: %s", err))
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("end_point_applications", flattenEndpointApplications(resp.EndPointApplications)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("end_point_application_groups", flattenEndpointApplicationGroups(resp.EndPointApplicationGroups)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEndpointDLPRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] endpoint dlp rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("endpoint dlp rule ID not set"))
	}

	req := expandEndpointDLPRules(d)

	// Drop the resource from state if it no longer exists.
	if _, err := endpoint_dlp_rules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	// Park the rule at the tail order before reordering. Writing the intended
	// order directly makes parallel updates fight over positions and can leave
	// the rules in a rotated order; the shared reorder loop then moves this rule
	// to its intended position. Mirrors firewall_filtering / ssl_inspection.
	existingRules, err := endpoint_dlp_rules.GetAll(ctx, service)
	if err != nil {
		log.Printf("[ERROR] error getting all endpoint dlp rules: %v", err)
	}
	sort.Slice(existingRules, func(i, j int) bool {
		return existingRules[i].Rank < existingRules[j].Rank || (existingRules[i].Rank == existingRules[j].Rank && existingRules[i].Order < existingRules[j].Order)
	})
	intendedOrder := req.Order
	intendedRank := req.Rank
	if len(existingRules) > 0 {
		req.Order = existingRules[len(existingRules)-1].Order
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, _, err := endpoint_dlp_rules.Update(ctx, service, id, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
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

	// Handle ordering update after a successful rule update
	resourceType := "endpoint_dlp_rules"

	reorderWithBeforeReorder(OrderRule{Order: intendedOrder, Rank: intendedRank}, id, resourceType,
		func() (map[int]OrderRule, error) {
			allRules, err := endpoint_dlp_rules.GetAll(ctx, service)
			if err != nil {
				return nil, err
			}
			m := make(map[int]OrderRule, len(allRules))
			for _, r := range allRules {
				m[r.ID] = OrderRule{Order: r.Order, Rank: r.Rank}
			}
			return m, nil
		},
		func(id int, order OrderRule) error {
			rule, err := endpoint_dlp_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			if rule.Order == order.Order && rule.Rank == order.Rank {
				return nil
			}

			rule.LastModifiedTime = 0
			rule.LastModifiedBy = nil
			rule.Order = order.Order
			rule.Rank = order.Rank

			log.Printf("[DEBUG] Updating parent rule ID %d to order %d", id, order.Order)

			_, _, err = endpoint_dlp_rules.Update(ctx, service, id, rule)
			if err != nil {
				log.Printf("[ERROR] Failed to update order for rule ID %d: %v", id, err)
			}
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	markOrderRuleAsDone(id, resourceType)
	waitForReorder(resourceType)

	return resourceEndpointDLPRulesRead(ctx, d, meta)
}

func resourceEndpointDLPRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] endpoint dlp rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting endpoint dlp rule ID: %v\n", (d.Id()))

	if _, err := endpoint_dlp_rules.Delete(ctx, service, id); err != nil {
		if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
			log.Printf("[INFO] endpoint dlp rule %d not found, skipping deletion", id)
			return nil
		}
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] endpoint dlp rule deleted")

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandEndpointDLPRules(d *schema.ResourceData) endpoint_dlp_rules.EndpointDlpRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandEndpointDLPRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	result := endpoint_dlp_rules.EndpointDlpRules{
		ID:                        id,
		Name:                      d.Get("name").(string),
		Order:                     order,
		Rank:                      d.Get("rank").(int),
		Description:               d.Get("description").(string),
		Action:                    d.Get("action").(string),
		State:                     d.Get("state").(string),
		Severity:                  d.Get("severity").(string),
		ExternalAuditorEmail:      d.Get("external_auditor_email").(string),
		WithoutContentInspection:  d.Get("without_content_inspection").(bool),
		MinSize:                   d.Get("min_size").(int),
		EunEnabled:                d.Get("eun_enabled").(bool),
		EunTemplateId:             d.Get("eun_template_id").(int),
		UcTemplateId:              d.Get("uc_template_id").(int),
		DataTransferMethod:        d.Get("data_transfer_method").(string),
		NetworkType:               d.Get("network_type").(string),
		FileTypes:                 SetToStringList(d, "file_types"),
		UserRiskScoreLevels:       SetToStringList(d, "user_risk_score_levels"),
		DeviceTrustLevels:         SetToStringList(d, "device_trust_levels"),
		Receiver:                  expandEndpointReceiver(d, "receiver"),
		Auditor:                   expandSingleIDNameExtensions(d, "auditor"),
		NotificationTemplate:      expandSingleIDNameExtensions(d, "notification_template"),
		Groups:                    expandIDNameExtensionsSet(d, "groups"),
		Departments:               expandIDNameExtensionsSet(d, "departments"),
		Users:                     expandIDNameExtensionsSet(d, "users"),
		Resources:                 expandIDNameExtensionsSet(d, "resources"),
		ResourceGroups:            expandIDNameExtensionsSet(d, "resource_groups"),
		Devices:                   expandIDNameExtensionsSet(d, "devices"),
		DeviceGroups:              expandIDNameExtensionsSet(d, "device_groups"),
		DlpEngines:                expandIDNameExtensionsSet(d, "dlp_engines"),
		Labels:                    expandIDNameExtensionsSet(d, "labels"),
		EndPointApplications:      expandEndpointApplications(d, "end_point_applications"),
		EndPointApplicationGroups: expandEndpointApplicationGroups(d, "end_point_application_groups"),
	}
	return result
}

// expandEndpointReceiver builds the SDK Receiver from the resource block. The
// user only supplies "id"; the receiver type is always forced to "ZIR".
func expandEndpointReceiver(d *schema.ResourceData, key string) *dlp_web_rules.Receiver {
	set, ok := d.Get(key).(*schema.Set)
	if !ok || set.Len() == 0 {
		return nil
	}
	item := set.List()[0].(map[string]interface{})
	idInt, err := strconv.Atoi(item["id"].(string))
	if err != nil || idInt == 0 {
		return nil
	}
	return &dlp_web_rules.Receiver{
		ID:   idInt,
		Type: endpointReceiverType,
	}
}

// flattenEndpointReceiver renders the SDK Receiver back into the resource block,
// exposing only the "id" attribute.
func flattenEndpointReceiver(receiver *dlp_web_rules.Receiver) []interface{} {
	if receiver == nil || receiver.ID == 0 {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id": strconv.Itoa(receiver.ID),
		},
	}
}

// endpointBlockHasValues reports whether the single-block set (groups/users
// shape) contains at least one value under innerKey.
func endpointBlockHasValues(raw interface{}, innerKey string) bool {
	set, ok := raw.(*schema.Set)
	if !ok || set.Len() == 0 {
		return false
	}
	block, ok := set.List()[0].(map[string]interface{})
	if !ok {
		return false
	}
	inner, ok := block[innerKey].(*schema.Set)
	return ok && inner.Len() > 0
}

func expandEndpointApplications(d *schema.ResourceData, key string) []common.EndPointApplications {
	raw, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	set := raw.(*schema.Set)
	if set.Len() == 0 {
		return nil
	}
	block := set.List()[0].(map[string]interface{})
	zappIDs, ok := block["zapp_id"].(*schema.Set)
	if !ok || zappIDs.Len() == 0 {
		return nil
	}
	out := make([]common.EndPointApplications, 0, zappIDs.Len())
	for _, v := range zappIDs.List() {
		out = append(out, common.EndPointApplications{
			ZappID: v.(string),
		})
	}
	return out
}

func expandEndpointApplicationGroups(d *schema.ResourceData, key string) []common.EndPointApplicationGroups {
	raw, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	set := raw.(*schema.Set)
	if set.Len() == 0 {
		return nil
	}
	block := set.List()[0].(map[string]interface{})
	groupIDs, ok := block["group_id"].(*schema.Set)
	if !ok || groupIDs.Len() == 0 {
		return nil
	}
	out := make([]common.EndPointApplicationGroups, 0, groupIDs.Len())
	for _, v := range groupIDs.List() {
		groupID, err := strconv.Atoi(v.(string))
		if err != nil {
			continue
		}
		out = append(out, common.EndPointApplicationGroups{
			GroupID: groupID,
		})
	}
	return out
}

func flattenEndpointApplications(apps []common.EndPointApplications) []interface{} {
	if len(apps) == 0 {
		return nil
	}
	zappIDs := make([]interface{}, 0, len(apps))
	for _, a := range apps {
		if a.ZappID == "" {
			continue
		}
		zappIDs = append(zappIDs, a.ZappID)
	}
	if len(zappIDs) == 0 {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"zapp_id": zappIDs,
		},
	}
}

func flattenEndpointApplicationGroups(groups []common.EndPointApplicationGroups) []interface{} {
	if len(groups) == 0 {
		return nil
	}
	groupIDs := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		if g.GroupID == 0 {
			continue
		}
		groupIDs = append(groupIDs, strconv.Itoa(g.GroupID))
	}
	if len(groupIDs) == 0 {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"group_id": groupIDs,
		},
	}
}
