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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
)

var (
	dlpWebRulesLock     sync.Mutex
	dlpWebStartingOrder int
)

func resourceDlpWebRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDlpWebRulesCreate,
		ReadContext:   resourceDlpWebRulesRead,
		UpdateContext: resourceDlpWebRulesUpdate,
		DeleteContext: resourceDlpWebRulesDelete,
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
				Description:  "The description of the DLP policy rule.",
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
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The rule order of execution for the DLP policy rule with respect to other rules.",
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
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of cloud applications to which the DLP policy rule must be applied.",
			},
			"file_types": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of file types for which the DLP policy rule must be applied.",
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
				Computed:    true,
				Description: "The email address of an external auditor to whom DLP email notifications are sent.",
			},
			"match_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The match only criteria for DLP engines.",
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
			"zscaler_incident_receiver": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether a Zscaler Incident Receiver is associated to the DLP policy rule.",
			},
			"locations":                setIDsSchemaTypeCustom(intPtr(8), "The Name-ID pairs of locations to which the DLP policy rule must be applied."),
			"location_groups":          setIDsSchemaTypeCustom(intPtr(32), "The Name-ID pairs of locations groups to which the DLP policy rule must be applied."),
			"users":                    setIDsSchemaTypeCustom(intPtr(4), "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"groups":                   setIDsSchemaTypeCustom(intPtr(8), "The Name-ID pairs of groups to which the DLP policy rule must be applied."),
			"departments":              setIDsSchemaTypeCustom(intPtr(8), "The Name-ID pairs of departments to which the DLP policy rule must be applied."),
			"excluded_departments":     setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"excluded_users":           setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"excluded_groups":          setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"included_domain_profiles": setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"excluded_domain_profiles": setIDsSchemaTypeCustom(intPtr(256), "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"workload_groups":          setIdNameSchemaCustom(255, "The list of preconfigured workload groups to which the policy must be applied"),
			"dlp_engines":              setIDsSchemaTypeCustom(intPtr(4), "The list of DLP engines to which the DLP policy rule must be applied."),
			"time_windows":             setIDsSchemaTypeCustom(intPtr(2), "list of time interval during which rule must be enforced."),
			"labels":                   setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule."),
			"source_ip_groups":         setIDsSchemaTypeCustom(nil, "list of source ip groups"),
			"url_categories":           setIDsSchemaTypeCustom(nil, "The list of URL categories to which the DLP policy rule must be applied."),
			"auditor":                  setSingleIDSchemaTypeCustom("The auditor to which the DLP policy rule must be applied."),
			"notification_template":    setSingleIDSchemaTypeCustom("The template used for DLP notification emails."),
			"icap_server":              setSingleIDSchemaTypeCustom("The DLP server, using ICAP, to which the transaction content is forwarded."),
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

	// Validate the OCR DLP web rules (assuming this is another validation function you have)
	if err := validateOCRDlpWebRules(req); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating zia web dlp rule\n%+v\n", req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		dlpWebRulesLock.Lock()
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
		dlpWebRulesLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = dlpWebStartingOrder

		resp, err := dlp_web_rules.Create(ctx, service, &req)
		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") && !strings.Contains(err.Error(), "ICAP Receiver with id") {
				if time.Since(start) < timeout {
					time.Sleep(5 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia web dlp rule request. Took: %s, without locking: %s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorder(order, resp.ID, "dlp_web_rules", func() (int, error) {
			list, err := dlp_web_rules.GetAll(ctx, service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := dlp_web_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = dlp_web_rules.Update(ctx, service, id, rule)
			return err
		})

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceDlpWebRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(resp.ID, "dlp_web_rules")
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
	_ = d.Set("external_auditor_email", resp.ExternalAuditorEmail)
	_ = d.Set("without_content_inspection", resp.WithoutContentInspection)
	_ = d.Set("dlp_download_scan_enabled", resp.DLPDownloadScanEnabled)
	_ = d.Set("zcc_notifications_enabled", resp.ZCCNotificationsEnabled)
	_ = d.Set("zscaler_incident_receiver", resp.ZscalerIncidentReceiver)

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

	if err := d.Set("auditor", flattenCustomIDSet(resp.Auditor)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("notification_template", flattenCustomIDSet(resp.NotificationTemplate)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("icap_server", flattenCustomIDSet(resp.IcapServer)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("excluded_groups", flattenIDExtensions(resp.ExcludedGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("excluded_departments", flattenIDExtensions(resp.ExcludedDepartments)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("excluded_users", flattenIDExtensions(resp.ExcludedUsers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_ip_groups", flattenIDs(resp.SourceIpGroups)); err != nil {
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
	}
	log.Printf("[INFO] Updating web dlp rule ID: %v\n", id)
	req := expandDlpWebRules(d)
	// Validate file types
	if err := validateDLPRuleFileTypes(req); err != nil {
		return diag.FromErr(err)
	}

	// Validate the OCR DLP web rules (assuming this is another validation function you have)
	if err := validateOCRDlpWebRules(req); err != nil {
		return diag.FromErr(err)
	}

	if _, err := dlp_web_rules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := dlp_web_rules.Update(ctx, service, id, &req)
		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if time.Since(start) < timeout {
					time.Sleep(5 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		reorder(req.Order, req.ID, "dlp_web_rules", func() (int, error) {
			list, err := dlp_web_rules.GetAll(ctx, service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := dlp_web_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = dlp_web_rules.Update(ctx, service, id, rule)
			return err
		})

		if diags := resourceDlpWebRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(req.ID, "dlp_web_rules")
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

func resourceDlpWebRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] web dlp rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting dlp rule ID: %v\n", (d.Id()))

	if _, err := dlp_web_rules.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] web dlp rule deleted")

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

func expandDlpWebRules(d *schema.ResourceData) dlp_web_rules.WebDLPRules {
	id, _ := getIntFromResourceData(d, "rule_id")
	result := dlp_web_rules.WebDLPRules{
		ID:                       id,
		Name:                     d.Get("name").(string),
		Order:                    d.Get("order").(int),
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
		MinSize:                  d.Get("min_size").(int),
		ParentRule:               d.Get("parent_rule").(int),
		Protocols:                SetToStringList(d, "protocols"),
		FileTypes:                SetToStringList(d, "file_types"),
		CloudApplications:        SetToStringList(d, "cloud_applications"),
		UserRiskScoreLevels:      SetToStringList(d, "user_risk_score_levels"),
		SubRules:                 expandSubRules(d.Get("sub_rules").(*schema.Set)),
		Auditor:                  expandIDNameExtensionsSetSingle(d, "auditor"),
		NotificationTemplate:     expandIDNameExtensionsSetSingle(d, "notification_template"),
		IcapServer:               expandIDNameExtensionsSetSingle(d, "icap_server"),
		Locations:                expandIDNameExtensionsSet(d, "locations"),
		LocationGroups:           expandIDNameExtensionsSet(d, "location_groups"),
		Groups:                   expandIDNameExtensionsSet(d, "groups"),
		Departments:              expandIDNameExtensionsSet(d, "departments"),
		Users:                    expandIDNameExtensionsSet(d, "users"),
		URLCategories:            expandIDNameExtensionsSet(d, "url_categories"),
		DLPEngines:               expandIDNameExtensionsSet(d, "dlp_engines"),
		TimeWindows:              expandIDNameExtensionsSet(d, "time_windows"),
		Labels:                   expandIDNameExtensionsSet(d, "labels"),
		ExcludedUsers:            expandIDNameExtensionsSet(d, "excluded_groups"),
		ExcludedGroups:           expandIDNameExtensionsSet(d, "excluded_departments"),
		ExcludedDepartments:      expandIDNameExtensionsSet(d, "excluded_users"),
		SourceIpGroups:           expandIDNameExtensionsSet(d, "source_ip_groups"),
		IncludedDomainProfiles:   expandIDNameExtensionsSet(d, "included_domain_profiles"),
		ExcludedDomainProfiles:   expandIDNameExtensionsSet(d, "excluded_domain_profiles"),
		WorkloadGroups:           expandWorkloadGroups(d, "workload_groups"),
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
