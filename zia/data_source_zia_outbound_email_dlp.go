package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/outbound_email_dlp"
)

func dataSourceOutboundEmailDLP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOutboundEmailDLPRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier of the DLP policy rule.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP policy rule name.",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The rule order of execution for the DLP policy rule with respect to other rules.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the DLP policy rule.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enables or disables the DLP policy rule.",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action taken when traffic matches the DLP policy rule criteria.",
			},
			"severity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the severity selected for the DLP rule violation.",
			},
			"file_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of file types to which the DLP policy rule must be applied.",
			},
			"content_locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of content locations to which the DLP policy rule must be applied.",
			},
			"user_risk_score_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The user risk score levels to which the DLP policy rule must be applied.",
			},
			"min_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum file size (in KB) used for evaluation of the DLP policy rule.",
			},
			"without_content_inspection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"external_auditor_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of an external auditor to whom DLP email notifications are sent.",
			},
			"custom_header": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The custom header value inserted when the rule action includes custom header insertion.",
			},
			"parent_rule": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the parent rule under which an exception rule is added.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the DLP policy rule was last modified.",
			},
			"sub_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of exception rules added to a parent rule.",
			},
			"groups":                   dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of groups to which the DLP policy rule must be applied."),
			"departments":              dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of departments to which the DLP policy rule must be applied."),
			"users":                    dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"excluded_groups":          dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of groups that are excluded from the DLP policy rule."),
			"excluded_departments":     dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of departments that are excluded from the DLP policy rule."),
			"excluded_users":           dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of users that are excluded from the DLP policy rule."),
			"time_windows":             dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of time windows to which the DLP policy rule must be applied."),
			"dlp_engines":              dataSourceOutboundEmailDLPIDNameExtensionsSchema("The list of DLP engines to which the DLP policy rule must be applied."),
			"labels":                   dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of rule labels associated to the DLP policy rule."),
			"included_domain_profiles": dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of domain profiles included in the DLP policy rule."),
			"email_tenants":            dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of email tenants to which the DLP policy rule must be applied."),
			"email_recipient_profiles": dataSourceOutboundEmailDLPIDNameExtensionsSchema("The Name-ID pairs of email recipient profiles to which the DLP policy rule must be applied."),
			"auditor":                  dataSourceOutboundEmailDLPIDNameExtensionsSchema("The auditor to which the DLP policy rule must be applied."),
			"notification_template":    dataSourceOutboundEmailDLPIDNameExtensionsSchema("The template used for DLP notification emails."),
			"receiver":                 dataSourceEndpointDLPReceiverSchema(),
			"last_modified_by":         dataSourceOutboundEmailDLPIDNameExtensionsSchema("The admin that modified the DLP policy rule last."),
		},
	}
}

func dataSourceOutboundEmailDLPIDNameExtensionsSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Identifier that uniquely identifies an entity",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The configured name of the entity",
				},
				"extensions": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func dataSourceOutboundEmailDLPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *outbound_email_dlp.OutboundEmailDlp
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for outbound email dlp rule id: %d\n", id)
		res, err := outbound_email_dlp.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for outbound email dlp rule: %s\n", name)
		res, err := outbound_email_dlp.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any outbound email dlp rule with name '%s' or id '%d'", name, id))
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
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
	_ = d.Set("last_modified_time", resp.LastModifiedTime)

	subRuleIDs := make([]interface{}, len(resp.SubRules))
	for i, subRule := range resp.SubRules {
		subRuleIDs[i] = strconv.Itoa(subRule.ID)
	}
	if err := d.Set("sub_rules", subRuleIDs); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sub_rules: %s", err))
	}

	if err := d.Set("groups", flattenIDExtensions(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("departments", flattenIDExtensions(resp.Departments)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("users", flattenIDExtensions(resp.Users)); err != nil {
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
	if err := d.Set("time_windows", flattenIDExtensions(resp.TimeWindows)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dlp_engines", flattenIDExtensions(resp.DlpEngines)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("labels", flattenIDExtensions(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("included_domain_profiles", flattenIDExtensions(resp.IncludedDomainProfiles)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_tenants", flattenIDExtensions(resp.EmailTenants)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_recipient_profiles", flattenIDExtensions(resp.EmailRecipientProfiles)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("auditor", flattenIDExtensionsList(resp.Auditor)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("notification_template", flattenIDExtensionsList(resp.NotificationTemplate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("receiver", flattenReceiver(resp.Receiver)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_modified_by", flattenIDExtensionsList(resp.LastModifiedBy)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
