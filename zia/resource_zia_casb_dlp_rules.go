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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/saas_security_api/casb_dlp_rules"
)

var (
	cloudCasbDlpRuleLock          sync.Mutex
	cloudCasbDlpRuleStartingOrder int
)

func resourceCasbDlpRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCasbDlpRulesCreate,
		ReadContext:   resourceCasbDlpRulesRead,
		UpdateContext: resourceCasbDlpRulesUpdate,
		DeleteContext: resourceCasbDlpRulesDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			contentLocation := d.Get("content_location").(string)
			domains, domainsSet := d.GetOk("domains")

			if contentLocation != "CONTENT_LOCATION_SHARED_CHANNEL" && domainsSet && len(domains.(*schema.Set).List()) > 0 {
				return fmt.Errorf("'domains' can only be set when 'content_location' is 'CONTENT_LOCATION_SHARED_CHANNEL'")
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
					resp, err := casb_dlp_rules.GetByRuleID(ctx, service, ruleType, idInt)
					if err != nil {
						return nil, err
					}
					d.SetId(strconv.Itoa(resp.ID))
					_ = d.Set("rule_id", resp.ID)
					_ = d.Set("type", ruleType)
				} else {
					// If identifier is a name
					resources, err := casb_dlp_rules.GetByRuleType(ctx, service, ruleType)
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
				Description: "System-generated identifier for the SaaS Security Data at Rest Scanning DLP rule",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "System-generated identifier for the SaaS Security Data at Rest Scanning DLP rule",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule name",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An admin editable text-based description of the rule",
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of SaaS Security Data at Rest Scanning DLP rule",
				ValidateFunc: validation.StringInSlice([]string{
					"OFLCASB_DLP_FILE",
					"OFLCASB_DLP_EMAIL",
					"OFLCASB_DLP_CRM",
					"OFLCASB_DLP_ITSM",
					"OFLCASB_DLP_COLLAB",
					"OFLCASB_DLP_REPO",
					"OFLCASB_DLP_STORAGE",
					"OFLCASB_DLP_GENAI",
				}, false),
			},
			"order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Order of rule execution with respect to other SaaS Security Data at Rest Scanning DLP rules",
			},
			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Admin rank that is assigned to this rule. Mandatory when admin rank-based access restriction is enabled",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ENABLED",
				Description: "Administrative state of the rule",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The configured action for the policy rule",
				ValidateFunc: validation.StringInSlice([]string{
					"OFLCASB_DLP_REPORT_INCIDENT",
					"OFLCASB_DLP_SHARE_READ_ONLY",
					"OFLCASB_DLP_EXTERNAL_SHARE_READ_ONLY",
					"OFLCASB_DLP_INTERNAL_SHARE_READ_ONLY",
					"OFLCASB_DLP_REMOVE_PUBLIC_LINK_SHARE",
					"OFLCASB_DLP_REVOKE_SHARE",
					"OFLCASB_DLP_REMOVE_EXTERNAL_SHARE",
					"OFLCASB_DLP_REMOVE_INTERNAL_SHARE",
					"OFLCASB_DLP_REMOVE_COLLABORATORS",
					"OFLCASB_DLP_REMOVE_INTERNAL_LINK_SHARE",
					"OFLCASB_DLP_REMOVE_DISCOVERABLE",
					"OFLCASB_DLP_NOTIFY_END_USER",
					"OFLCASB_DLP_APPLY_MIP_TAG",
					"OFLCASB_DLP_APPLY_BOX_TAG",
					"OFLCASB_DLP_MOVE_TO_RESTRICTED_FOLDER",
					"OFLCASB_DLP_REMOVE",
					"OFLCASB_DLP_QUARANTINE",
					"OFLCASB_DLP_APPLY_EMAIL_TAG",
					"OFLCASB_DLP_APPLY_GOOGLEDRIVE_LABEL",
					"OFLCASB_DLP_REMOVE_EXT_COLLABORATORS",
					"OFLCASB_DLP_QUARANTINE_TO_USER_ROOT_FOLDER",
					"OFLCASB_DLP_APPLY_WATERMARK",
					"OFLCASB_DLP_REMOVE_WATERMARK",
					"OFLCASB_DLP_APPLY_HEADER",
					"OFLCASB_DLP_APPLY_FOOTER",
					"OFLCASB_DLP_APPLY_HEADER_FOOTER",
					"OFLCASB_DLP_REMOVE_HEADER",
					"OFLCASB_DLP_REMOVE_FOOTER",
					"OFLCASB_DLP_REMOVE_HEADER_FOOTER",
					"OFLCASB_DLP_BLOCK",
					"OFLCASB_DLP_APPLY_ATLASSIAN_CLASSIFICATION_LABEL",
					"OFLCASB_DLP_ALLOW",
					"OFLCASB_DLP_REDACT",
				}, false),
			},
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The severity level of the incidents that match the policy rule",
				ValidateFunc: validation.StringInSlice([]string{
					"RULE_SEVERITY_HIGH",
					"RULE_SEVERITY_MEDIUM",
					"RULE_SEVERITY_LOW",
					"RULE_SEVERITY_INFO",
				}, false),
			},
			"bucket_owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A user who inspect their buckets for sensitive data. When you choose a user, their buckets are available in the Buckets field",
			},
			"external_auditor_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address of the external auditor to whom the DLP email alerts are sent",
			},
			"content_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location for the content that the Zscaler service inspects for sensitive data",
				ValidateFunc: validation.StringInSlice([]string{
					"CONTENT_LOCATION_PRIVATE_CHANNEL",
					"CONTENT_LOCATION_PUBLIC_CHANNEL",
					"CONTENT_LOCATION_SHARED_CHANNEL",
					"CONTENT_LOCATION_DIRECT_MESSAGE",
					"CONTENT_LOCATION_MULTI_PERSON_DIRECT_MESSAGE",
				}, false),
			},
			"recipient": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies if the email recipient is internal or external",
			},
			"quarantine_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location where all the quarantined files are moved and necessary actions are taken by either deleting or restoring the data",
			},
			"watermark_delete_old_version": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to delete an old version of the watermarked file",
			},
			"include_criteria_domain_profile": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, criteriaDomainProfiles is included as part of the criteria, else they are excluded from the criteria.",
			},
			"include_email_recipient_profile": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, emailRecipientProfiles is included as part of the criteria, else they are excluded from the criteria.",
			},
			"without_content_inspection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, Content Matching is set to None",
			},
			"include_entity_groups": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "If true, entityGroups is included as part of the criteria, else are excluded from the criteria",
			},
			"domains": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The domain for the external organization sharing the channel",
			},
			"cloud_app_tenants":         setIDsSchemaTypeCustom(nil, "Name-ID pairs of the cloud application tenants for which the rule is applied"),
			"entity_groups":             setIDsSchemaTypeCustom(nil, "Name-ID pairs of entity groups that are part of the rule criteria"),
			"included_domain_profiles":  setIDsSchemaTypeCustom(nil, "Name-ID pairs of domain profiles included in the criteria for the rule"),
			"excluded_domain_profiles":  setIDsSchemaTypeCustom(nil, "Name-ID pairs of domain profiles excluded in the criteria for the rule"),
			"criteria_domain_profiles":  setIDsSchemaTypeCustom(nil, "Name-ID pairs of domain profiles that are mandatory in the criteria for the rule"),
			"email_recipient_profiles":  setIDsSchemaTypeCustom(nil, "Name-ID pairs of recipient profiles for which the rule is applied"),
			"object_types":              setIDsSchemaTypeCustom(nil, "List of object types for which the rule is applied"),
			"labels":                    setIDsSchemaTypeCustom(intPtr(1), "Name-ID pairs of rule labels associated with the rule"),
			"dlp_engines":               setIDsSchemaTypeCustom(nil, "The list of DLP engines to which the DLP policy rule must be applied"),
			"buckets":                   setIDsSchemaTypeCustom(intPtr(8), "The buckets for the Zscaler service to inspect for sensitive data"),
			"groups":                    setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of groups for which the rule is applied"),
			"departments":               setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of departments for which rule must be applied"),
			"users":                     setIDsSchemaTypeCustom(intPtr(4), "Name-ID pairs of users for which rule must be applied"),
			"zscaler_incident_receiver": setSingleIDSchemaTypeCustom("The Zscaler Incident Receiver details"),
			"auditor_notification":      setSingleIDSchemaTypeCustom("Notification template used for DLP email alerts sent to the auditor"),
			"tag":                       setSingleIDSchemaTypeCustom("Tag applied to the rule"),
			"watermark_profile":         setSingleIDSchemaTypeCustom("Watermark profile applied to the rule"),
			"redaction_profile":         setSingleIDSchemaTypeCustom("Name-ID of the redaction profile in the criteria"),
			"casb_email_label":          setSingleIDSchemaTypeCustom("Name-ID of the email label associated with the rule"),
			"casb_tombstone_template":   setSingleIDSchemaTypeCustom("Name-ID of the quarantine tombstone template associated with the rule"),
			"collaboration_scope":       getCasbRuleCollaborationScope(),
			"file_types":                getFileTypes(),
			"components":                getCasbRuleComponents(),
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

func resourceCasbDlpRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandCasbDlpRules(d)
	log.Printf("[INFO] Creating zia casb dlp rule\n%+v\n", req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		cloudCasbDlpRuleLock.Lock()
		if cloudCasbDlpRuleStartingOrder == 0 {
			rules, _ := casb_dlp_rules.GetByRuleType(ctx, service, req.Type)
			for _, r := range rules {
				if r.Order > cloudCasbDlpRuleStartingOrder {
					cloudCasbDlpRuleStartingOrder = r.Order
				}
			}
			if cloudCasbDlpRuleStartingOrder == 0 {
				cloudCasbDlpRuleStartingOrder = 1
			}
		}
		cloudCasbDlpRuleLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = cloudCasbDlpRuleStartingOrder

		resp, err := casb_dlp_rules.Create(ctx, service, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Creating casb dlp rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.Errorf("error creating resource: %v", err)
		}

		log.Printf("[INFO] Created zia casb dlp rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorderWithBeforeReorder(
			OrderRule{Order: order, Rank: req.Rank},
			resp.ID,
			"casb_dlp_rules",
			func() (int, error) {
				rules, err := casb_dlp_rules.GetByRuleType(ctx, service, req.Type)
				if err != nil {
					return 0, err
				}
				return len(rules), nil
			},
			func(id int, order OrderRule) error {
				rule, err := casb_dlp_rules.GetByRuleID(ctx, service, req.Type, id)
				if err != nil {
					return fmt.Errorf("failed to retrieve rule by ID: %v", err)
				}
				// Optional: avoid unnecessary updates if the current order is already correct
				if rule.Order == order.Order {
					return nil
				}
				rule.Order = order.Order
				_, err = casb_dlp_rules.Update(ctx, service, id, rule)

				if err != nil {
					return fmt.Errorf("failed to update rule order: %v", err)
				}
				return nil
			},
			nil)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceCasbDlpRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(resp.ID, "casb_dlp_rules")
		break
	}

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

	return resourceCasbDlpRulesRead(ctx, d, meta)
}

func resourceCasbDlpRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia casb dlp rule id is set"))
	}
	ruleType, ok := d.Get("type").(string)
	if !ok || ruleType == "" {
		return diag.FromErr(fmt.Errorf("no rule type is set"))
	}
	resp, err := casb_dlp_rules.GetByRuleID(ctx, service, ruleType, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing casb dlp rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting casb dlp rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("type", resp.Type)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("state", resp.State)
	_ = d.Set("action", resp.Action)
	_ = d.Set("severity", resp.Severity)
	_ = d.Set("components", resp.Components)
	_ = d.Set("bucket_owner", resp.BucketOwner)
	_ = d.Set("external_auditor_email", resp.ExternalAuditorEmail)
	_ = d.Set("collaboration_scope", resp.CollaborationScope)
	_ = d.Set("content_location", resp.ContentLocation)
	_ = d.Set("file_types", resp.FileTypes)
	_ = d.Set("domains", resp.Domains)
	_ = d.Set("recipient", resp.Recipient)
	_ = d.Set("quarantine_location", resp.QuarantineLocation)
	_ = d.Set("watermark_delete_old_version", resp.WatermarkDeleteOldVersion)
	_ = d.Set("include_criteria_domain_profile", resp.IncludeCriteriaDomainProfile)
	_ = d.Set("include_email_recipient_profile", resp.IncludeEmailRecipientProfile)
	_ = d.Set("without_content_inspection", resp.WithoutContentInspection)
	_ = d.Set("include_entity_groups", resp.IncludeEntityGroups)

	if err := d.Set("groups", flattenIDExtensionsListIDs(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("departments", flattenIDExtensionsListIDs(resp.Departments)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("users", flattenIDExtensionsListIDs(resp.Users)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dlp_engines", flattenIDExtensionsListIDs(resp.DLPEngines)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("included_domain_profiles", flattenIDExtensionsListIDs(resp.IncludedDomainProfiles)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("excluded_domain_profiles", flattenIDExtensionsListIDs(resp.ExcludedDomainProfiles)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("criteria_domain_profiles", flattenIDExtensionsListIDs(resp.CriteriaDomainProfiles)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("email_recipient_profiles", flattenIDExtensionsListIDs(resp.EmailRecipientProfiles)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entity_groups", flattenIDExtensionsListIDs(resp.EntityGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("object_types", flattenIDExtensionsListIDs(resp.ObjectTypes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("buckets", flattenIDExtensionsListIDs(resp.Buckets)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cloud_app_tenants", flattenIDExtensionsListIDs(resp.CloudAppTenants)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("casb_tombstone_template", flattenCustomIDSet(resp.CasbTombstoneTemplate)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tag", flattenCustomIDSet(resp.Tag)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("watermark_profile", flattenCustomIDSet(resp.WatermarkProfile)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zscaler_incident_receiver", flattenCustomIDSet(resp.ZscalerIncidentReceiver)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("receiver", flattenReceiverCASBResource(resp.Receiver)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting receiver: %s", err))
	}

	if err := d.Set("auditor_notification", flattenCustomIDSet(resp.AuditorNotification)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCasbDlpRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] cloud application control rule ID not set: %v\n", id)
		return diag.Errorf("cloud application control rule ID not set")
	}

	ruleType, ok := d.Get("type").(string)
	if !ok || ruleType == "" {
		return diag.Errorf("no rule type is set")
	}

	log.Printf("[INFO] Updating zia cloud application control rule ID: %v\n", id)
	req := expandCasbDlpRules(d)

	if _, err := casb_dlp_rules.GetByRuleID(ctx, service, ruleType, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := casb_dlp_rules.Update(ctx, service, id, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating casb dlp rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.Errorf("error updating resource: %v", err)
		}

		reorderWithBeforeReorder(OrderRule{Order: req.Order, Rank: req.Rank}, req.ID, "casb_dlp_rules",
			func() (int, error) {
				rules, err := casb_dlp_rules.GetByRuleType(ctx, service, req.Type)
				if err != nil {
					return 0, err
				}
				return len(rules), nil
			},
			func(id int, order OrderRule) error {
				rule, err := casb_dlp_rules.GetByRuleID(ctx, service, req.Type, id)
				if err != nil {
					return fmt.Errorf("failed to retrieve rule by ID: %v", err)
				}
				// Optional: avoid unnecessary updates if the current order is already correct
				if rule.Order == order.Order {
					return nil
				}
				rule.Order = order.Order
				_, err = casb_dlp_rules.Update(ctx, service, id, rule)

				if err != nil {
					return fmt.Errorf("failed to update rule order: %v", err)
				}
				return nil
			},
			nil)

		if diags := resourceCasbDlpRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(req.ID, "casb_dlp_rules")
		break
	}

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

	return resourceCasbDlpRulesRead(ctx, d, meta)
}

func resourceCasbDlpRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] cloud application control rule not set: %v\n", id)
	}
	ruleType, ok := d.Get("type").(string)
	if !ok || ruleType == "" {
		return diag.FromErr(fmt.Errorf("no rule type is set"))
	}
	log.Printf("[INFO] Deleting cloud application control rule ID: %v\n", (d.Id()))

	if _, err := casb_dlp_rules.Delete(ctx, service, ruleType, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] cloud application control rule deleted")
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

func expandCasbDlpRules(d *schema.ResourceData) casb_dlp_rules.CasbDLPRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandCasbDlpRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	result := casb_dlp_rules.CasbDLPRules{
		ID:                           id,
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		Type:                         d.Get("type").(string),
		Order:                        order,
		State:                        d.Get("state").(string),
		Rank:                         d.Get("rank").(int),
		Action:                       d.Get("action").(string),
		Severity:                     d.Get("severity").(string),
		BucketOwner:                  d.Get("bucket_owner").(string),
		ContentLocation:              d.Get("content_location").(string),
		Recipient:                    d.Get("recipient").(string),
		QuarantineLocation:           d.Get("quarantine_location").(string),
		WatermarkDeleteOldVersion:    d.Get("watermark_delete_old_version").(bool),
		IncludeCriteriaDomainProfile: d.Get("include_criteria_domain_profile").(bool),
		IncludeEmailRecipientProfile: d.Get("include_email_recipient_profile").(bool),
		IncludeEntityGroups:          d.Get("include_entity_groups").(bool),
		WithoutContentInspection:     d.Get("without_content_inspection").(bool),
		ExternalAuditorEmail:         d.Get("external_auditor_email").(string),
		Components:                   SetToStringList(d, "components"),
		CollaborationScope:           SetToStringList(d, "collaboration_scope"),
		Domains:                      SetToStringList(d, "domains"),
		FileTypes:                    SetToStringList(d, "file_types"),
		Buckets:                      expandIDNameExtensionsSet(d, "buckets"),
		Groups:                       expandIDNameExtensionsSet(d, "groups"),
		Departments:                  expandIDNameExtensionsSet(d, "departments"),
		Users:                        expandIDNameExtensionsSet(d, "users"),
		Labels:                       expandIDNameExtensionsSet(d, "labels"),
		DLPEngines:                   expandIDNameExtensionsSet(d, "dlp_engines"),
		CloudAppTenants:              expandIDNameExtensionsSet(d, "cloud_app_tenants"),
		EntityGroups:                 expandIDNameExtensionsSet(d, "entity_groups"),
		IncludedDomainProfiles:       expandIDNameExtensionsSet(d, "included_domain_profiles"),
		ExcludedDomainProfiles:       expandIDNameExtensionsSet(d, "exclude_domain_profiles"),
		CriteriaDomainProfiles:       expandIDNameExtensionsSet(d, "criteria_domain_profiles"),
		EmailRecipientProfiles:       expandIDNameExtensionsSet(d, "email_recipient_profiles"),
		ObjectTypes:                  expandIDNameExtensionsSet(d, "object_types"),
		ZscalerIncidentReceiver:      expandIDNameExtensionsSetSingle(d, "zscaler_incident_receiver"),
		AuditorNotification:          expandIDNameExtensionsSetSingle(d, "auditor_notification"),
		WatermarkProfile:             expandIDNameExtensionsSetSingle(d, "watermark_profile"),
		RedactionProfile:             expandIDNameExtensionsSetSingle(d, "redaction_profile"),
		CasbEmailLabel:               expandIDNameExtensionsSetSingle(d, "casb_email_label"),
		CasbTombstoneTemplate:        expandIDNameExtensionsSetSingle(d, "casb_tombstone_template"),
		Tag:                          expandIDNameExtensionsSetSingle(d, "tag"),
		Receiver:                     expandCASBReceiver(d, "receiver"),
	}

	return result
}

func expandCASBReceiver(d *schema.ResourceData, key string) *casb_dlp_rules.Receiver {
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

	receiver := &casb_dlp_rules.Receiver{
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

func flattenReceiverCASBResource(receiver *casb_dlp_rules.Receiver) []interface{} {
	if receiver == nil {
		return nil
	}

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
