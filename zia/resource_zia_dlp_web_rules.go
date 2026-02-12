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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
)

var (
	dlpWebRulesLock             sync.Mutex
	dlpWebStartingOrder         int
	dlpWebSubRulesStartingOrder map[int]int = make(map[int]int)
)

func resourceDlpWebRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDlpWebRulesCreate,
		ReadContext:   resourceDlpWebRulesRead,
		UpdateContext: resourceDlpWebRulesUpdate,
		DeleteContext: resourceDlpWebRulesDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
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
					resp, err := dlp_web_rules.GetByName(ctx, service, id)
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
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"protocols": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The protocol criteria specified for the DLP policy rule.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rank": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 7),
				Description:  "Admin rank of the admin who creates this rule",
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
				Computed:    true,
				Description: "Indicates the severity selected for the DLP rule violation",
				ValidateFunc: validation.StringInSlice([]string{
					"RULE_SEVERITY_HIGH",
					"RULE_SEVERITY_MEDIUM",
					"RULE_SEVERITY_LOW",
					"RULE_SEVERITY_INFO",
				}, false),
			},
			"parent_rule": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the parent rule under which an exception rule is added",
			},
			"sub_rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of exception rules added to a parent rule",
			},
			"user_risk_score_levels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "",
			},
			"cloud_applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `The list of cloud applications to which the DLP policy rule must be applied
				Use the data source zia_cloud_applications to get the list of available cloud applications:
				https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_cloud_applications
				`,
			},
			"file_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `The list of file types for which the DLP policy rule must be applied,
				See the Web DLP Rules API for the list of available File types:
				https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-get`,
			},
			"min_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 96000),
				Description:  "The minimum file size (in KB) used for evaluation of the DLP policy rule.",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The action taken when traffic matches the DLP policy rule criteria.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"NONE",
					"BLOCK",
					"CONFIRM",
					"ALLOW",
					"ICAP_RESPONSE",
				}, false),
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
			"match_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The match only criteria for DLP engines.",
			},
			"inspect_http_get_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
			},
			"without_content_inspection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"dlp_download_scan_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"zcc_notifications_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"eun_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The EUN template ID associated with the rule",
			},
			"zscaler_incident_receiver": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether a Zscaler Incident Receiver is associated to the DLP policy rule",
			},
			"locations":                setIDsSchemaTypeCustom(nil, "The Name-ID pairs of locations to which the DLP policy rule must be applied"),
			"location_groups":          setIDsSchemaTypeCustom(nil, "The Name-ID pairs of locations groups to which the DLP policy rule must be applied"),
			"users":                    setIDsSchemaTypeCustom(nil, "The Name-ID pairs of users to which the DLP policy rule must be applied"),
			"groups":                   setIDsSchemaTypeCustom(nil, "The Name-ID pairs of groups to which the DLP policy rule must be applied"),
			"departments":              setIDsSchemaTypeCustom(nil, "The Name-ID pairs of departments to which the DLP policy rule must be applied"),
			"excluded_departments":     setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of departments which the DLP policy rule must exclude"),
			"excluded_users":           setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of users which the DLP policy rule must exclude"),
			"excluded_groups":          setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of groups which the DLP policy rule must exclude"),
			"included_domain_profiles": setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of domain profiiles which the DLP policy rule must include"),
			"excluded_domain_profiles": setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of domain profiles to which the DLP policy rule must exclude"),
			"workload_groups":          setIdNameSchemaCustom(255, "The list of preconfigured workload groups to which the policy must be applied"),
			"dlp_engines":              setIDsSchemaTypeCustom(intPtr(4), "The list of DLP engines to which the DLP policy rule must be applied"),
			"time_windows":             setIDsSchemaTypeCustom(intPtr(2), "list of source ip groups"),
			"file_type_categories":     setIDsSchemaTypeCustom(nil, "The list of file types to which the rule applies"),
			"labels":                   setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule"),
			"source_ip_groups":         setIDsSchemaTypeCustom(nil, "list of source ip groups"),
			"url_categories":           setIDsSchemaTypeCustom(nil, "The list of URL categories to which the DLP policy rule must be applied"),
			"auditor":                  setSingleIDSchemaTypeCustom("The auditor to which the DLP policy rule must be applied"),
			"notification_template":    setSingleIDSchemaTypeCustom("The template used for DLP notification emails"),
			"icap_server":              setSingleIDSchemaTypeCustom("The DLP server, using ICAP, to which the transaction content is forwarded"),
			"receiver": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "The receiver information for the DLP policy rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier for the receiver",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the receiver",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of the receiver",
						},
						"tenant": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Tenant information for the receiver",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Unique identifier for the tenant",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the tenant",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceDlpWebRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandDlpWebRules(d)

	// Validate file types
	if err := validateDLPRuleFileTypes(req); err != nil {
		return diag.FromErr(err)
	}

	// Validate OCR DLP rules
	if err := validateOCRDlpWebRules(req); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating zia web dlp rule\n%+v\n", req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()
	// Determine whether it's a sub-rule
	isSubRule := req.ParentRule != 0
	for {
		dlpWebRulesLock.Lock()
		if isSubRule {
			if dlpWebSubRulesStartingOrder[req.ParentRule] == 0 {
				list, _ := dlp_web_rules.Get(ctx, service, req.ParentRule)
				for _, subRule := range list.SubRules {
					r, err := dlp_web_rules.Get(ctx, service, subRule.ID)
					if err != nil {
						return diag.FromErr(err)
					}
					if r.Order > dlpWebSubRulesStartingOrder[req.ParentRule] {
						dlpWebSubRulesStartingOrder[req.ParentRule] = r.Order
					}
				}
				if dlpWebSubRulesStartingOrder[req.ParentRule] == 0 {
					dlpWebSubRulesStartingOrder[req.ParentRule] = 1
				}
			}
		} else {
			if dlpWebStartingOrder == 0 {
				list, _ := dlp_web_rules.GetAll(ctx, service)
				for _, r := range list {
					if r.Order > dlpWebStartingOrder {
						dlpWebStartingOrder = r.Order
					}
				}
				if dlpWebStartingOrder == 0 {
					dlpWebStartingOrder = 1
				}
			}
		}
		dlpWebRulesLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		if isSubRule {
			req.Order = dlpWebSubRulesStartingOrder[req.ParentRule]
		} else {
			req.Order = dlpWebStartingOrder
		}

		resp, err := dlp_web_rules.Create(ctx, service, &req)

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

		log.Printf("[INFO] Created zia web dlp rule request. Took: %s, without locking: %s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)

		resourceType := "dlp_web_rules"
		if isSubRule {
			resourceType = fmt.Sprintf("dlp_web_rules_sub_%d", req.ParentRule)
		}

		reorderWithBeforeReorder(
			OrderRule{Order: order, Rank: req.Rank},
			resp.ID,
			resourceType,
			func() (int, error) {
				if isSubRule {
					parent, err := dlp_web_rules.Get(ctx, service, req.ParentRule)
					if err != nil {
						return 0, err
					}
					return len(parent.SubRules), nil
				}
				allRules, err := dlp_web_rules.GetAll(ctx, service)
				if err != nil {
					return 0, err
				}
				// Count all rules including predefined ones for proper ordering
				return len(allRules), nil
			},
			func(id int, order OrderRule) error {
				// Custom updateOrder that handles predefined rules
				rule, err := dlp_web_rules.Get(ctx, service, id)
				if err != nil {
					return err
				}

				rule.Order = order.Order
				rule.Rank = order.Rank

				if rule.ParentRule != 0 {
					log.Printf("[DEBUG] Updating sub-rule ID %d (parent ID: %d) to order %d", id, rule.ParentRule, order.Order)
				} else {
					log.Printf("[DEBUG] Updating parent rule ID %d to order %d", id, order.Order)
				}

				_, err = dlp_web_rules.Update(ctx, service, id, rule)
				if err != nil {
					log.Printf("[ERROR] Failed to update order for rule ID %d: %v", id, err)
				}
				return err
			},
			nil, // Remove beforeReorder function to avoid adding too many rules to the map
		)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceDlpWebRulesRead(ctx, d, meta); diags.HasError() {
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

func resourceDlpWebRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia web dlp rule id is set"))
	}
	resp, err := dlp_web_rules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing web dlp rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting web dlp rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("protocols", resp.Protocols)
	_ = d.Set("description", resp.Description)
	_ = d.Set("file_types", resp.FileTypes)
	_ = d.Set("cloud_applications", resp.CloudApplications)
	_ = d.Set("user_risk_score_levels", resp.UserRiskScoreLevels)
	_ = d.Set("state", resp.State)
	_ = d.Set("min_size", resp.MinSize)
	_ = d.Set("action", resp.Action)
	_ = d.Set("severity", resp.Severity)
	_ = d.Set("parent_rule", resp.ParentRule)
	_ = d.Set("match_only", resp.MatchOnly)
	_ = d.Set("inspect_http_get_enabled", resp.InspectHttpGetEnabled)
	_ = d.Set("external_auditor_email", resp.ExternalAuditorEmail)
	_ = d.Set("without_content_inspection", resp.WithoutContentInspection)
	_ = d.Set("dlp_download_scan_enabled", resp.DLPDownloadScanEnabled)
	_ = d.Set("zcc_notifications_enabled", resp.ZCCNotificationsEnabled)
	_ = d.Set("zscaler_incident_receiver", resp.ZscalerIncidentReceiver)
	_ = d.Set("eun_template_id", resp.EUNTemplateID)

	// Flatten sub_rules
	subRuleIDs := make([]interface{}, len(resp.SubRules))
	for i, subRule := range resp.SubRules {
		subRuleIDs[i] = strconv.Itoa(subRule.ID)
	}
	if err := d.Set("sub_rules", subRuleIDs); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sub_rules: %s", err))
	}

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenIDExtensionsListIDs(resp.LocationGroups)); err != nil {
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

	if err := d.Set("included_domain_profiles", flattenIDExtensionsListIDs(resp.IncludedDomainProfiles)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("excluded_domain_profiles", flattenIDExtensionsListIDs(resp.ExcludedDomainProfiles)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url_categories", flattenIDExtensionsListIDs(resp.URLCategories)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dlp_engines", flattenIDExtensionsListIDs(resp.DLPEngines)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_windows", flattenIDExtensionsListIDs(resp.TimeWindows)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("file_type_categories", flattenIDListIDs(resp.FileTypeCategories)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("auditor", flattenCustomIDSet(resp.Auditor)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("notification_template", flattenCustomIDSet(resp.NotificationTemplate)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("icap_server", flattenCustomIDSet(resp.IcapServer)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("receiver", flattenReceiverResource(resp.Receiver)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting receiver: %s", err))
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
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
	if err := d.Set("source_ip_groups", flattenIDExtensionsListIDs(resp.SourceIpGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting workload_groups: %s", err))
	}
	return nil
}

func resourceDlpWebRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] web dlp rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("web dlp rule ID not set"))
	}

	req := expandDlpWebRules(d)

	// Validate file types
	if err := validateDLPRuleFileTypes(req); err != nil {
		return diag.FromErr(err)
	}

	// Validate OCR DLP rules
	if err := validateOCRDlpWebRules(req); err != nil {
		return diag.FromErr(err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := dlp_web_rules.Update(ctx, service, id, &req)

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
	isSubRule := req.ParentRule != 0
	resourceType := "dlp_web_rules"
	if isSubRule {
		resourceType = fmt.Sprintf("dlp_web_rules_sub_%d", req.ParentRule)
	}

	reorderWithBeforeReorder(OrderRule{Order: req.Order, Rank: req.Rank}, id, resourceType,
		func() (int, error) {
			if isSubRule {
				parent, err := dlp_web_rules.Get(ctx, service, req.ParentRule)
				if err != nil {
					return 0, err
				}
				return len(parent.SubRules), nil
			}
			allRules, err := dlp_web_rules.GetAll(ctx, service)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			rule, err := dlp_web_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			// Optional: avoid unnecessary updates if the current order is already correct
			if rule.Order == order.Order && rule.Rank == order.Rank {
				return nil
			}

			// Update the order
			rule.Order = order.Order
			rule.Rank = order.Rank

			// Log and ensure ParentRule is set for sub-rules
			if rule.ParentRule != 0 {
				log.Printf("[DEBUG] Updating sub-rule ID %d (parent ID: %d) to order %d", id, rule.ParentRule, order.Order)
			} else {
				log.Printf("[DEBUG] Updating parent rule ID %d to order %d", id, order.Order)
			}

			_, err = dlp_web_rules.Update(ctx, service, id, rule)
			if err != nil {
				log.Printf("[ERROR] Failed to update order for rule ID %d: %v", id, err)
			}
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	markOrderRuleAsDone(id, resourceType)
	waitForReorder(resourceType)

	return nil
}

func resourceDlpWebRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] web dlp rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting dlp rule ID: %v\n", (d.Id()))

	if _, err := dlp_web_rules.Delete(ctx, service, id); err != nil {
		if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
			log.Printf("[INFO] web dlp rule %d not found, skipping deletion", id)
			return nil
		}
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] web dlp rule deleted")

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

func expandDlpWebRules(d *schema.ResourceData) dlp_web_rules.WebDLPRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandDlpWebRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	result := dlp_web_rules.WebDLPRules{
		ID:                       id,
		Name:                     d.Get("name").(string),
		Order:                    order,
		Rank:                     d.Get("rank").(int),
		Description:              d.Get("description").(string),
		Action:                   d.Get("action").(string),
		State:                    d.Get("state").(string),
		Severity:                 d.Get("severity").(string),
		ExternalAuditorEmail:     d.Get("external_auditor_email").(string),
		MatchOnly:                d.Get("match_only").(bool),
		WithoutContentInspection: d.Get("without_content_inspection").(bool),
		DLPDownloadScanEnabled:   d.Get("dlp_download_scan_enabled").(bool),
		ZCCNotificationsEnabled:  d.Get("zcc_notifications_enabled").(bool),
		ZscalerIncidentReceiver:  d.Get("zscaler_incident_receiver").(bool),
		InspectHttpGetEnabled:    d.Get("inspect_http_get_enabled").(bool),
		MinSize:                  d.Get("min_size").(int),
		ParentRule:               d.Get("parent_rule").(int),
		EUNTemplateID:            d.Get("eun_template_id").(int),
		Protocols:                SetToStringList(d, "protocols"),
		FileTypes:                SetToStringList(d, "file_types"),
		CloudApplications:        SetToStringList(d, "cloud_applications"),
		UserRiskScoreLevels:      SetToStringList(d, "user_risk_score_levels"),
		SubRules:                 expandSubRules(d.Get("sub_rules").(*schema.Set)),
		Auditor:                  expandIDNameExtensionsSetSingle(d, "auditor"),
		NotificationTemplate:     expandIDNameExtensionsSetSingle(d, "notification_template"),
		IcapServer:               expandIDNameExtensionsSetSingle(d, "icap_server"),
		Receiver:                 expandReceiver(d, "receiver"),
		Locations:                expandIDNameExtensionsSet(d, "locations"),
		LocationGroups:           expandIDNameExtensionsSet(d, "location_groups"),
		Groups:                   expandIDNameExtensionsSet(d, "groups"),
		Departments:              expandIDNameExtensionsSet(d, "departments"),
		Users:                    expandIDNameExtensionsSet(d, "users"),
		URLCategories:            expandIDNameExtensionsSet(d, "url_categories"),
		DLPEngines:               expandIDNameExtensionsSet(d, "dlp_engines"),
		TimeWindows:              expandIDNameExtensionsSet(d, "time_windows"),
		Labels:                   expandIDNameExtensionsSet(d, "labels"),
		ExcludedUsers:            expandIDNameExtensionsSet(d, "excluded_users"),
		ExcludedGroups:           expandIDNameExtensionsSet(d, "excluded_groups"),
		ExcludedDepartments:      expandIDNameExtensionsSet(d, "excluded_departments"),
		SourceIpGroups:           expandIDNameExtensionsSet(d, "source_ip_groups"),
		IncludedDomainProfiles:   expandIDNameExtensionsSet(d, "included_domain_profiles"),
		ExcludedDomainProfiles:   expandIDNameExtensionsSet(d, "excluded_domain_profiles"),
		WorkloadGroups:           expandWorkloadGroupsIDName(d, "workload_groups"),
		FileTypeCategories:       expandIDSet(d, "file_type_categories"),
	}
	return result
}

func expandSubRules(set *schema.Set) []dlp_web_rules.SubRule {
	var subRules []dlp_web_rules.SubRule
	for _, item := range set.List() {
		subRuleID, err := strconv.Atoi(item.(string))
		if err == nil {
			subRules = append(subRules, dlp_web_rules.SubRule{ID: subRuleID})
		}
	}
	return subRules
}

func expandReceiver(d *schema.ResourceData, key string) *dlp_web_rules.Receiver {
	receiverSet, ok := d.Get(key).(*schema.Set)
	if !ok || receiverSet.Len() == 0 {
		return nil
	}

	// Since receiver is a single item set, get the first (and only) item
	receiverList := receiverSet.List()
	if len(receiverList) == 0 {
		return nil
	}

	item := receiverList[0].(map[string]interface{})

	// Convert string ID to int for the SDK struct
	idStr := item["id"].(string)
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		// If conversion fails, use 0 as default
		idInt = 0
	}

	receiver := &dlp_web_rules.Receiver{
		ID:   idInt,
		Name: item["name"].(string),
		Type: item["type"].(string),
	}

	// Handle tenant if present
	if tenantList, ok := item["tenant"].([]interface{}); ok && len(tenantList) > 0 {
		if tenantMap, ok := tenantList[0].(map[string]interface{}); ok {
			tenantIDStr := tenantMap["id"].(string)
			tenantIDInt, err := strconv.Atoi(tenantIDStr)
			if err != nil {
				tenantIDInt = 0
			}
			receiver.Tenant = &common.IDNameExtensions{
				ID:   tenantIDInt,
				Name: tenantMap["name"].(string),
			}
		}
	}

	return receiver
}

func flattenReceiverResource(receiver *dlp_web_rules.Receiver) []interface{} {
	if receiver == nil {
		return nil
	}

	// Check if the receiver is actually empty (no meaningful data)
	if receiver.ID == 0 && receiver.Name == "" && receiver.Type == "" && receiver.Tenant == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":   strconv.Itoa(receiver.ID),
		"name": receiver.Name,
		"type": receiver.Type,
	}

	// Handle tenant if present
	if receiver.Tenant != nil {
		tenant := map[string]interface{}{
			"id":   strconv.Itoa(receiver.Tenant.ID),
			"name": receiver.Tenant.Name,
		}
		result["tenant"] = []interface{}{tenant}
	}

	return []interface{}{result}
}

func flattenIDListIDs(list []common.IDName) []interface{} {
	if len(list) == 0 {
		// Return an empty slice instead of nil
		return []interface{}{}
	}

	ids := []int{}
	for _, item := range list {
		if item.ID == 0 && item.Name == "" {
			continue
		}
		ids = append(ids, item.ID)
	}

	if len(ids) == 0 {
		// Again return []interface{}{} instead of nil
		return []interface{}{}
	}

	// The rest remains the same
	return []interface{}{
		map[string]interface{}{
			"id": ids,
		},
	}
}

func expandIDSet(d *schema.ResourceData, key string) []common.IDName {
	setInterface, ok := d.GetOk(key)
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.IDName
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil && itemMap["id"] != nil {
				set := itemMap["id"].(*schema.Set)
				for _, id := range set.List() {
					result = append(result, common.IDName{
						ID: id.(int),
					})
				}
			}
		}
		return result
	}
	return []common.IDName{}
}
