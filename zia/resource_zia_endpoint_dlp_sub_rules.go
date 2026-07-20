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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_dlp_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_dlp_sub_rules"
)

// Sub-rules are scoped to a parent rule. The reorder engine is keyed per parent
// so the contiguous 1..N ordering is maintained independently for each parent's
// exception list. The starting-order map is keyed by parent rule ID.
var (
	endpointDlpSubRulesLock          sync.Mutex
	endpointDlpSubRulesStartingOrder = make(map[int]int)
)

func resourceEndpointDLPSubRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointDLPSubRulesCreate,
		ReadContext:   resourceEndpointDLPSubRulesRead,
		UpdateContext: resourceEndpointDLPSubRulesUpdate,
		DeleteContext: resourceEndpointDLPSubRulesDelete,
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

			if isExternalSet && isAuditorSet {
				return fmt.Errorf("'external_auditor_email' and 'auditor' cannot both be set")
			}

			if isExternalSet {
				if !isNTSet {
					return fmt.Errorf("when setting 'external_auditor_email', 'notification_template' must also be set")
				}
				return nil
			}

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

				// A sub-rule can only be located through its parent, so the import
				// key must carry both: "<parentRuleID>:<subRuleID>" or
				// "<parentRuleID>:<subRuleName>".
				raw := d.Id()
				parts := strings.Split(raw, ":")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid import id %q, expected format <parentRuleID>:<subRuleID> or <parentRuleID>:<subRuleName>", raw)
				}

				parentID, err := strconv.Atoi(parts[0])
				if err != nil {
					return nil, fmt.Errorf("invalid parent rule id %q: %w", parts[0], err)
				}
				_ = d.Set("parent_rule", parentID)

				parent, err := endpoint_dlp_rules.Get(ctx, service, parentID)
				if err != nil {
					return nil, fmt.Errorf("unable to read parent rule %d: %w", parentID, err)
				}

				if subID, convErr := strconv.Atoi(parts[1]); convErr == nil {
					d.SetId(strconv.Itoa(subID))
					_ = d.Set("sub_rule_id", subID)
					return []*schema.ResourceData{d}, nil
				}

				for _, sr := range parent.SubRules {
					if strings.EqualFold(sr.Name, parts[1]) {
						d.SetId(strconv.Itoa(sr.ID))
						_ = d.Set("sub_rule_id", sr.ID)
						return []*schema.ResourceData{d}, nil
					}
				}
				return nil, fmt.Errorf("sub-rule %q not found under parent rule %d", parts[1], parentID)
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parent_rule": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The unique identifier of the parent rule under which this exception (sub) rule is created. Changing this value recreates the sub-rule under the new parent.",
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
				Description:  "The rule order of execution for the sub-rule with respect to other sub-rules of the same parent.",
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
				ValidateFunc: validation.StringInSlice([]string{
					"TRUSTED",
					"OFFTRUSTED",
					"VPN",
					"ANY",
				}, false),
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
					"ALLOW",
					"BLOCK",
					"CONFIRM",
					"PROTECT",
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

func resourceEndpointDLPSubRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	parentID := d.Get("parent_rule").(int)
	req := expandEndpointDLPSubRules(d)

	log.Printf("[INFO] Creating zia endpoint dlp sub-rule under parent %d\n%+v\n", parentID, req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()
	resourceType := fmt.Sprintf("endpoint_dlp_sub_rules_%d", parentID)

	for {
		endpointDlpSubRulesLock.Lock()
		if endpointDlpSubRulesStartingOrder[parentID] == 0 {
			if parent, err := readEndpointDLPParentRule(ctx, service, parentID); err == nil && parent != nil {
				for _, sr := range parent.SubRules {
					if sr.Order > endpointDlpSubRulesStartingOrder[parentID] {
						endpointDlpSubRulesStartingOrder[parentID] = sr.Order
					}
				}
			}
			if endpointDlpSubRulesStartingOrder[parentID] == 0 {
				endpointDlpSubRulesStartingOrder[parentID] = 1
			}
		}
		endpointDlpSubRulesLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = endpointDlpSubRulesStartingOrder[parentID]

		resp, _, err := endpoint_dlp_sub_rules.Create(ctx, service, parentID, &req)

		// If the sub-rule already exists upstream (typically because an earlier
		// apply created it but failed to persist it to state), adopt the
		// existing sub-rule by name and reconcile it to the desired config
		// instead of failing with DUPLICATE_ITEM.
		if err != nil && isDuplicateItemError(err) {
			if existing, found := findEndpointDLPSubRuleByName(ctx, service, parentID, req.Name); found {
				log.Printf("[WARN] endpoint dlp sub-rule %q already exists under parent %d (id %d); adopting and reconciling", req.Name, parentID, existing.ID)
				req.ID = existing.ID
				if _, _, uerr := endpoint_dlp_sub_rules.Update(ctx, service, parentID, existing.ID, &req); uerr == nil {
					created := req
					resp = &created
					err = nil
				} else {
					err = uerr
				}
			}
		}

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

		log.Printf("[INFO] Created zia endpoint dlp sub-rule. Took: %s, without locking: %s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp.ID)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("sub_rule_id", resp.ID)

		// Converge the sub-rule's position within the parent's subRules block.
		// The engine reads the block (via a live parent read) and only moves
		// sub-rules whose order/rank differs from the declared value — exactly
		// the pattern used by the other rule-based resources, applied to the
		// nested subRules list.
		reorderWithBeforeReorder(
			OrderRule{Order: order, Rank: req.Rank},
			resp.ID,
			resourceType,
			endpointDLPSubRulesGetCurrent(ctx, service, parentID),
			endpointDLPSubRulesUpdateOrder(ctx, service, parentID),
			nil,
		)

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

	// Record state from a live read of the parent's subRules block after the
	// reorder has converged, so the persisted order reflects the API (and the
	// declared configuration) rather than an intermediate position.
	return resourceEndpointDLPSubRulesRead(ctx, d, meta)
}

func resourceEndpointDLPSubRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	parentID := d.Get("parent_rule").(int)
	subID, ok := getIntFromResourceData(d, "sub_rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia endpoint dlp sub-rule id is set"))
	}

	// A sub-rule is only retrievable through its parent rule; there is no
	// dedicated GET for an individual sub-rule. Read the parent with a live
	// view of its subRules block so the recorded order reflects the API.
	parent, err := readEndpointDLPParentRule(ctx, service, parentID)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing endpoint dlp sub-rule %s from state because its parent rule %d no longer exists in ZIA", d.Id(), parentID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var resp *endpoint_dlp_rules.EndpointDlpRules
	for i := range parent.SubRules {
		if parent.SubRules[i].ID == subID {
			resp = &parent.SubRules[i]
			break
		}
	}
	if resp == nil {
		log.Printf("[WARN] Removing endpoint dlp sub-rule %s from state because it no longer exists under parent rule %d", d.Id(), parentID)
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Getting endpoint dlp sub-rule:\n%+v\n", resp)

	return setEndpointDLPSubRuleState(d, subRuleFromParentEntry(*resp), parentID)
}

// setEndpointDLPSubRuleState writes a sub-rule (sourced from the parent's
// subRules block) into Terraform state.
func setEndpointDLPSubRuleState(d *schema.ResourceData, r *endpoint_dlp_sub_rules.EndpointDlpSubRules, parentID int) diag.Diagnostics {
	d.SetId(strconv.Itoa(r.ID))
	_ = d.Set("sub_rule_id", r.ID)
	_ = d.Set("parent_rule", parentID)
	_ = d.Set("name", r.Name)
	_ = d.Set("order", r.Order)
	_ = d.Set("rank", r.Rank)
	_ = d.Set("description", r.Description)
	_ = d.Set("file_types", r.FileTypes)
	_ = d.Set("user_risk_score_levels", r.UserRiskScoreLevels)
	_ = d.Set("device_trust_levels", r.DeviceTrustLevels)
	_ = d.Set("state", r.State)
	_ = d.Set("min_size", r.MinSize)
	_ = d.Set("action", r.Action)
	_ = d.Set("severity", r.Severity)
	_ = d.Set("data_transfer_method", r.DataTransferMethod)
	_ = d.Set("network_type", r.NetworkType)
	_ = d.Set("external_auditor_email", r.ExternalAuditorEmail)
	_ = d.Set("without_content_inspection", r.WithoutContentInspection)
	_ = d.Set("eun_enabled", r.EunEnabled)
	_ = d.Set("eun_template_id", r.EunTemplateId)
	_ = d.Set("uc_template_id", r.UcTemplateId)

	if err := d.Set("groups", flattenIDExtensionsListIDs(r.Groups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("departments", flattenIDExtensionsListIDs(r.Departments)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("device_groups", flattenIDExtensionsListIDs(r.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("devices", flattenIDExtensionsListIDs(r.Devices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("users", flattenIDExtensionsListIDs(r.Users)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resources", flattenIDExtensionsListIDs(r.Resources)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resource_groups", flattenIDExtensionsListIDs(r.ResourceGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dlp_engines", flattenIDExtensionsListIDs(r.DlpEngines)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("auditor", flattenSingleIDNameExtensions(r.Auditor)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("notification_template", flattenSingleIDNameExtensions(r.NotificationTemplate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("receiver", flattenEndpointReceiver(r.Receiver)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting receiver: %s", err))
	}
	if err := d.Set("labels", flattenIDExtensionsListIDs(r.Labels)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("end_point_applications", flattenEndpointApplications(r.EndPointApplications)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("end_point_application_groups", flattenEndpointApplicationGroups(r.EndPointApplicationGroups)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// isDuplicateItemError reports whether the API rejected a write because an item
// with the same name already exists.
func isDuplicateItemError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "DUPLICATE_ITEM")
}

// findEndpointDLPSubRuleByName looks a sub-rule up by name under its parent. It
// is used to adopt a sub-rule that already exists upstream (for example, one
// created by a previous apply that failed to persist to state).
func findEndpointDLPSubRuleByName(ctx context.Context, service *zscaler.Service, parentID int, name string) (*endpoint_dlp_sub_rules.EndpointDlpSubRules, bool) {
	parent, err := readEndpointDLPParentRule(ctx, service, parentID)
	if err != nil || parent == nil {
		return nil, false
	}
	for i := range parent.SubRules {
		if strings.EqualFold(parent.SubRules[i].Name, name) {
			return subRuleFromParentEntry(parent.SubRules[i]), true
		}
	}
	return nil, false
}

func resourceEndpointDLPSubRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	parentID := d.Get("parent_rule").(int)
	subID, ok := getIntFromResourceData(d, "sub_rule_id")
	if !ok {
		log.Printf("[ERROR] endpoint dlp sub-rule ID not set: %v\n", subID)
		return diag.FromErr(fmt.Errorf("endpoint dlp sub-rule ID not set"))
	}

	req := expandEndpointDLPSubRules(d)
	req.ID = subID

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, _, err := endpoint_dlp_sub_rules.Update(ctx, service, parentID, subID, &req)

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

	// Converge the sub-rule's position within the parent's subRules block after
	// the update, mirroring the other rule-based resources.
	resourceType := fmt.Sprintf("endpoint_dlp_sub_rules_%d", parentID)
	reorderWithBeforeReorder(
		OrderRule{Order: req.Order, Rank: req.Rank},
		subID,
		resourceType,
		endpointDLPSubRulesGetCurrent(ctx, service, parentID),
		endpointDLPSubRulesUpdateOrder(ctx, service, parentID),
		nil,
	)

	markOrderRuleAsDone(subID, resourceType)
	waitForReorder(resourceType)

	// Record state from a live read of the parent's subRules block after the
	// reorder has converged.
	return resourceEndpointDLPSubRulesRead(ctx, d, meta)
}

func resourceEndpointDLPSubRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	parentID := d.Get("parent_rule").(int)
	subID, ok := getIntFromResourceData(d, "sub_rule_id")
	if !ok {
		log.Printf("[ERROR] endpoint dlp sub-rule ID not set: %v\n", subID)
	}
	log.Printf("[INFO] Deleting endpoint dlp sub-rule ID: %v (parent %d)\n", d.Id(), parentID)

	if _, err := endpoint_dlp_sub_rules.Delete(ctx, service, parentID, subID); err != nil {
		if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
			log.Printf("[INFO] endpoint dlp sub-rule %d not found, skipping deletion", subID)
			return nil
		}
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] endpoint dlp sub-rule deleted")

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

// readEndpointDLPParentRule returns the parent rule (including its subRules
// block) with a live view of the subRules order.
//
// Sub-rule ordering is expressed by the `order` field of each entry in the
// parent's `subRules` list, so every ordering decision reads the parent rule.
// A sub-rule write, however, targets .../endPointDlpRules/{parent}/subRule/{id}
// and does not invalidate the cached GET of the parent rule itself
// (.../endPointDlpRules/{parent}). A plain Get can therefore hand back a stale
// subRules snapshot mid-apply — the reorder engine then under-counts the
// orderable sub-rules and never converges. Appending a unique query parameter
// produces a distinct cache key so this resource's ordering reads always
// reflect the live subRules block. This is scoped to this resource and does not
// alter the shared cache, client, or reorder engine. If the read fails for any
// reason it falls back to the standard accessor.
func readEndpointDLPParentRule(ctx context.Context, service *zscaler.Service, parentID int) (*endpoint_dlp_rules.EndpointDlpRules, error) {
	var rule endpoint_dlp_rules.EndpointDlpRules
	endpoint := fmt.Sprintf("/zia/api/v1/endPointDlpRules/%d?_ts=%d", parentID, time.Now().UnixNano())
	if err := service.Client.Read(ctx, endpoint, &rule); err != nil {
		return endpoint_dlp_rules.Get(ctx, service, parentID)
	}
	return &rule, nil
}

// endpointDLPSubRulesGetCurrent returns the reorder engine's snapshot callback:
// the current Order/Rank of every sub-rule under the given parent.
func endpointDLPSubRulesGetCurrent(ctx context.Context, service *zscaler.Service, parentID int) func() (map[int]OrderRule, error) {
	return func() (map[int]OrderRule, error) {
		parent, err := readEndpointDLPParentRule(ctx, service, parentID)
		if err != nil {
			return nil, err
		}
		m := make(map[int]OrderRule, len(parent.SubRules))
		for _, r := range parent.SubRules {
			m[r.ID] = OrderRule{Order: r.Order, Rank: r.Rank}
		}
		return m, nil
	}
}

// endpointDLPSubRulesUpdateOrder returns the reorder engine's update callback.
// It re-sends a sub-rule through the dedicated sub-rule endpoint with the new
// order/rank, sourcing the full body from the parent's sub-rule list.
func endpointDLPSubRulesUpdateOrder(ctx context.Context, service *zscaler.Service, parentID int) func(id int, order OrderRule) error {
	return func(id int, order OrderRule) error {
		parent, err := readEndpointDLPParentRule(ctx, service, parentID)
		if err != nil {
			return err
		}
		var entry *endpoint_dlp_rules.EndpointDlpRules
		for i := range parent.SubRules {
			if parent.SubRules[i].ID == id {
				entry = &parent.SubRules[i]
				break
			}
		}
		if entry == nil {
			return fmt.Errorf("sub-rule %d not found under parent rule %d during reorder", id, parentID)
		}
		if entry.Order == order.Order && entry.Rank == order.Rank {
			return nil
		}

		// Re-send the full sub-rule body sourced from the parent's subRules
		// block (a live read), changing only its position. Every field the API
		// requires — including data_transfer_method — is present in the block.
		sub := subRuleFromParentEntry(*entry)
		sub.Order = order.Order
		sub.Rank = order.Rank

		log.Printf("[DEBUG] Updating sub-rule ID %d (parent ID: %d) to order %d", id, parentID, order.Order)

		if _, _, err := endpoint_dlp_sub_rules.Update(ctx, service, parentID, id, sub); err != nil {
			log.Printf("[ERROR] Failed to update order for sub-rule ID %d: %v", id, err)
			return err
		}
		return nil
	}
}

// subRuleFromParentEntry converts a sub-rule as returned inside a parent rule
// (endpoint_dlp_rules.EndpointDlpRules) into the write payload used by the
// dedicated sub-rule endpoint. Read-only stamps are dropped.
func subRuleFromParentEntry(r endpoint_dlp_rules.EndpointDlpRules) *endpoint_dlp_sub_rules.EndpointDlpSubRules {
	return &endpoint_dlp_sub_rules.EndpointDlpSubRules{
		ID:                        r.ID,
		Name:                      r.Name,
		State:                     r.State,
		Order:                     r.Order,
		Rank:                      r.Rank,
		FileTypes:                 r.FileTypes,
		DataTransferMethod:        r.DataTransferMethod,
		Description:               r.Description,
		MinSize:                   r.MinSize,
		DeviceTrustLevels:         r.DeviceTrustLevels,
		Action:                    r.Action,
		ExternalAuditorEmail:      r.ExternalAuditorEmail,
		ParentRule:                r.ParentRule,
		Severity:                  r.Severity,
		UserRiskScoreLevels:       r.UserRiskScoreLevels,
		EunEnabled:                r.EunEnabled,
		EunTemplateId:             r.EunTemplateId,
		UcTemplateId:              r.UcTemplateId,
		NetworkType:               r.NetworkType,
		WithoutContentInspection:  r.WithoutContentInspection,
		NotificationTemplate:      r.NotificationTemplate,
		Auditor:                   r.Auditor,
		Resources:                 r.Resources,
		ResourceGroups:            r.ResourceGroups,
		Labels:                    r.Labels,
		DlpEngines:                r.DlpEngines,
		Users:                     r.Users,
		Groups:                    r.Groups,
		Departments:               r.Departments,
		Devices:                   r.Devices,
		DeviceGroups:              r.DeviceGroups,
		Receiver:                  r.Receiver,
		EndPointApplications:      r.EndPointApplications,
		EndPointApplicationGroups: r.EndPointApplicationGroups,
	}
}

func expandEndpointDLPSubRules(d *schema.ResourceData) endpoint_dlp_sub_rules.EndpointDlpSubRules {
	id, _ := getIntFromResourceData(d, "sub_rule_id")

	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandEndpointDLPSubRules: sub-rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	result := endpoint_dlp_sub_rules.EndpointDlpSubRules{
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
		ParentRule:                d.Get("parent_rule").(int),
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
