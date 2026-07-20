package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_dlp_rules"
)

func dataSourceEndpointDLPRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEndpointDLPRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP policy rule name.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enables or disables the DLP policy rule.",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The rule order of execution for the DLP policy rule with respect to other rules.",
			},
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Admin rank of the admin who creates this rule.",
			},
			"file_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of file types to which the DLP policy rule must be applied.",
			},
			"data_transfer_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The data transfer method to which the DLP policy rule must be applied.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the DLP policy rule.",
			},
			"min_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum file size (in KB) used for evaluation of the DLP policy rule.",
			},
			"device_trust_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of device trust levels for which the policy rule must be applied.",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action taken when traffic matches the DLP policy rule criteria.",
			},
			"external_auditor_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of an external auditor to whom DLP email notifications are sent.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the DLP policy rule was last modified.",
			},
			"parent_rule": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the parent rule under which an exception rule is added.",
			},
			"severity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the severity selected for the DLP rule violation.",
			},
			"user_risk_score_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The user risk score levels to which the DLP policy rule must be applied.",
			},
			"eun_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the End User Notification (EUN) is enabled for the DLP policy rule.",
			},
			"eun_template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the End User Notification (EUN) template.",
			},
			"uc_template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the user confirmation template.",
			},
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network type to which the DLP policy rule must be applied.",
			},
			"without_content_inspection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"notification_template": dataSourceEndpointDLPIDNameExtensionsSchema("The template used for DLP notification emails."),
			"auditor":               dataSourceEndpointDLPIDNameExtensionsSchema("The auditor to which the DLP policy rule must be applied."),
			"last_modified_by":      dataSourceEndpointDLPIDNameExtensionsSchema("The admin that modified the DLP policy rule last."),
			"receiver":              dataSourceEndpointDLPReceiverSchema(),
			"resources":             dataSourceEndpointDLPIDNameExtensionsSchema("The Name-ID pairs of resources to which the DLP policy rule must be applied."),
			"resource_groups":       dataSourceEndpointDLPIDNameExtensionsSchema("The Name-ID pairs of resource groups to which the DLP policy rule must be applied."),
			"labels":                dataSourceEndpointDLPIDNameExtensionsSchema("The Name-ID pairs of rule labels associated to the DLP policy rule."),
			"dlp_engines":           dataSourceEndpointDLPIDNameExtensionsSchema("The list of DLP engines to which the DLP policy rule must be applied."),
			"users":                 dataSourceEndpointDLPIDNameExtensionsSchema("The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"groups":                dataSourceEndpointDLPIDNameExtensionsSchema("The Name-ID pairs of groups to which the DLP policy rule must be applied."),
			"departments":           dataSourceEndpointDLPIDNameExtensionsSchema("The Name-ID pairs of departments to which the DLP policy rule must be applied."),
			"devices":               dataSourceEndpointDLPIDNameExtensionsSchema("The Name-ID pairs of devices to which the DLP policy rule must be applied."),
			"device_groups":         dataSourceEndpointDLPIDNameExtensionsSchema("The Name-ID pairs of device groups to which the DLP policy rule must be applied."),
			"sub_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of exception rules added to a parent rule.",
			},
			"end_point_applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of endpoint applications to which the DLP policy rule must be applied.",
				Elem:        dataSourceEndpointDLPApplicationsSchema(),
			},
			"end_point_application_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of endpoint application groups to which the DLP policy rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mod_uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_modified_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_point_applications": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataSourceEndpointDLPApplicationsSchema(),
						},
					},
				},
			},
		},
	}
}

func dataSourceEndpointDLPIDNameExtensionsSchema(description string) *schema.Schema {
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
					Description: "Identifier that uniquely identifies an entity",
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

func dataSourceEndpointDLPApplicationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bundle_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"original_file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"digitally_signed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mod_uid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"application_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zapp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceEndpointDLPVersionSchema(),
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceEndpointDLPVersionSchema(),
			},
		},
	}
}

func dataSourceEndpointDLPReceiverSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The receiver information for the DLP policy rule",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Unique identifier for the receiver",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the receiver",
				},
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Type of the receiver",
				},
				"tenant": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Tenant information for the receiver",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:        schema.TypeInt,
								Computed:    true,
								Description: "Unique identifier for the tenant",
							},
							"name": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Name of the tenant",
							},
							"external_id": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "External identifier for the tenant",
							},
							"extensions": {
								Type:        schema.TypeMap,
								Computed:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Description: "Additional properties for the tenant",
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceEndpointDLPVersionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"z_ver_id_md32": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"threat_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"threat_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bundle_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code_signing_certificate_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"threat_level_updated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceEndpointDLPRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *endpoint_dlp_rules.EndpointDlpRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for endpoint dlp rule id: %d\n", id)
		res, err := endpoint_dlp_rules.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for endpoint dlp rule: %s\n", name)
		res, err := endpoint_dlp_rules.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("state", resp.State)
		_ = d.Set("order", resp.Order)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("file_types", resp.FileTypes)
		_ = d.Set("data_transfer_method", resp.DataTransferMethod)
		_ = d.Set("description", resp.Description)
		_ = d.Set("min_size", resp.MinSize)
		_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)
		_ = d.Set("action", resp.Action)
		_ = d.Set("external_auditor_email", resp.ExternalAuditorEmail)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)
		_ = d.Set("parent_rule", resp.ParentRule)
		_ = d.Set("severity", resp.Severity)
		_ = d.Set("user_risk_score_levels", resp.UserRiskScoreLevels)
		_ = d.Set("eun_enabled", resp.EunEnabled)
		_ = d.Set("eun_template_id", resp.EunTemplateId)
		_ = d.Set("uc_template_id", resp.UcTemplateId)
		_ = d.Set("network_type", resp.NetworkType)
		_ = d.Set("without_content_inspection", resp.WithoutContentInspection)

		// Flatten sub_rules
		subRuleIDs := make([]interface{}, len(resp.SubRules))
		for i, subRule := range resp.SubRules {
			subRuleIDs[i] = strconv.Itoa(subRule.ID)
		}
		if err := d.Set("sub_rules", subRuleIDs); err != nil {
			return diag.FromErr(fmt.Errorf("error setting sub_rules: %s", err))
		}

		if err := d.Set("notification_template", flattenIDExtensionsList(resp.NotificationTemplate)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("auditor", flattenIDExtensionsList(resp.Auditor)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("last_modified_by", flattenIDExtensionsList(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("receiver", flattenReceiver(resp.Receiver)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting receiver: %s", err))
		}
		if err := d.Set("resources", flattenIDExtensions(resp.Resources)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("resource_groups", flattenIDExtensions(resp.ResourceGroups)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("labels", flattenIDExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("dlp_engines", flattenIDExtensions(resp.DlpEngines)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("users", flattenIDExtensions(resp.Users)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("groups", flattenIDExtensions(resp.Groups)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("departments", flattenIDExtensions(resp.Departments)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("devices", flattenIDExtensions(resp.Devices)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("device_groups", flattenIDExtensions(resp.DeviceGroups)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("end_point_applications", flattenEndpointDLPApplications(resp.EndPointApplications)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting end_point_applications: %s", err))
		}
		if err := d.Set("end_point_application_groups", flattenEndpointDLPApplicationGroups(resp.EndPointApplicationGroups)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting end_point_application_groups: %s", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any endpoint dlp rule with name '%s' or id '%d'", name, id))
	}

	return nil
}

func flattenEndpointDLPApplications(apps []common.EndPointApplications) []interface{} {
	if len(apps) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, len(apps))
	for i, a := range apps {
		out[i] = map[string]interface{}{
			"resource_id":        a.ResourceID,
			"description":        a.Description,
			"os_type":            a.OsType,
			"application_name":   a.ApplicationName,
			"bundle_id":          a.Bundle,
			"filename":           a.Filename,
			"original_file_name": a.OriginalFileName,
			"digitally_signed":   a.DigitallySigned,
			"mod_uid":            a.ModUId,
			"last_modified_time": a.LastModifiedTime,
			"application_type":   a.ApplicationType,
			"zapp_id":            a.ZappID,
			"deleted":            a.Deleted,
			"version":            flattenEndpointDLPVersion(a.Version),
			"versions":           flattenEndpointDLPVersions(a.Versions),
		}
	}
	return out
}

func flattenEndpointDLPApplicationGroups(groups []common.EndPointApplicationGroups) []interface{} {
	if len(groups) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, len(groups))
	for i, g := range groups {
		out[i] = map[string]interface{}{
			"group_id":               g.GroupID,
			"name":                   g.Name,
			"description":            g.Description,
			"mod_uid":                g.ModUId,
			"last_modified_time":     g.LastModifiedTime,
			"end_point_applications": flattenEndpointDLPApplications(g.EndPointApplications),
		}
	}
	return out
}

func flattenEndpointDLPVersion(v common.Version) []interface{} {
	if v == (common.Version{}) {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"version":                         v.Version,
			"z_ver_id_md32":                   v.ZverIDMD32,
			"threat_type":                     v.ThreatType,
			"threat_level":                    v.ThreatLevel,
			"bundle_id":                       v.Bundle,
			"code_signing_certificate_status": v.CodeSigningCertificateStatus,
			"threat_level_updated":            v.ThreatLevelUpdated,
		},
	}
}

func flattenEndpointDLPVersions(v common.Versions) []interface{} {
	if v == (common.Versions{}) {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"version":                         v.Version,
			"z_ver_id_md32":                   v.ZverIDMD32,
			"threat_type":                     v.ThreatType,
			"threat_level":                    v.ThreatLevel,
			"bundle_id":                       v.Bundle,
			"code_signing_certificate_status": v.CodeSigningCertificateStatus,
			"threat_level_updated":            v.ThreatLevelUpdated,
		},
	}
}
