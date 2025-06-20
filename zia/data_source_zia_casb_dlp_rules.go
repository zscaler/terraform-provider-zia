package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/saas_security_api/casb_dlp_rules"
)

func dataSourceCasbDlpRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCasbDlpRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "System-generated identifier for the SaaS Security Data at Rest Scanning DLP rule",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Rule name",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of SaaS Security Data at Rest Scanning DLP rule",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Order of rule execution with respect to other SaaS Security Data at Rest Scanning DLP rules",
			},
			"rank": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Administrative state of the rule",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configured action for the policy rule",
			},
			"severity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The severity level of the incidents that match the policy rule",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An admin editable text-based description of the rule",
			},
			"bucket_owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A user who inspect their buckets for sensitive data. When you choose a user, their buckets are available in the Buckets field.",
			},
			"external_auditor_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of the external auditor to whom the DLP email alerts are sent",
			},
			"content_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location for the content that the Zscaler service inspects for sensitive data",
			},
			"number_of_internal_collaborators": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Selects the number of internal collaborators for files that are shared with specific collaborators or are discoverable within an organization",
			},
			"number_of_external_collaborators": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Selects the number of external collaborators for files that are shared with specific collaborators outside of an organization",
			},
			"recipient": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies if the email recipient is internal or external",
			},
			"quarantine_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location where all the quarantined files are moved and necessary actions are taken by either deleting or restoring the data",
			},
			"access_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access privilege of this rule based on the admin's Role Based Authorization (RBA) state",
			},
			"watermark_delete_old_version": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether to delete an old version of the watermarked file",
			},
			"include_criteria_domain_profile": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, criteriaDomainProfiles is included as part of the criteria, else they are excluded from the criteria.",
			},
			"include_email_recipient_profile": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, emailRecipientProfiles is included as part of the criteria, else they are excluded from the criteria.",
			},
			"without_content_inspection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, Content Matching is set to None",
			},
			"include_entity_groups": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, entityGroups is included as part of the criteria, else are excluded from the criteria",
			},
			"components": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of components for which the rule is applied. Zscaler service inspects these components for sensitive data",
			},
			"collaboration_scope": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Collaboration scope for the rule",
			},
			"domains": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The domain for the external organization sharing the channel",
			},
			"file_types": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "File types for which the rule is applied. If not set, the rule is applied across all file types",
			},
			"cloud_app_tenants": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the cloud application tenants for which the rule is applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"entity_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of entity groups that are part of the rule criteria",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"included_domain_profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of domain profiles included in the criteria for the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"excluded_domain_profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of domain profiles excluded in the criteria for the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"criteria_domain_profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of domain profiles that are mandatory in the criteria for the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"email_recipient_profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of recipient profiles for which the rule is applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"buckets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The buckets for the Zscaler service to inspect for sensitive data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"object_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of object types for which the rule is applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"dlp_engines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of DLP engines to which the DLP policy rule must be applied",
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
			},
			"labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of rule labels associated with the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of groups for which the rule is applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"departments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of departments for which the rule is applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of users for which rule is applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"zscaler_incident_receiver": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Zscaler Incident Receiver details",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"auditor_notification": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Notification template used for DLP email alerts sent to the auditor",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"tag": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tag applied to the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"watermark_profile": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Watermark profile applied to the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"redaction_profile": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID of the redaction profile in the criteria",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"casb_email_label": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID of the email label associated with the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
			"casb_tombstone_template": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID of the quarantine tombstone template associated with the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
					},
				},
			},
		},
	}
}

func dataSourceCasbDlpRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	ruleType, ok := d.Get("type").(string)
	if !ok || ruleType == "" {
		return diag.FromErr(fmt.Errorf("type must be specified"))
	}

	var resp *casb_dlp_rules.CasbDLPRules
	id, idOk := getIntFromResourceData(d, "id")
	if idOk {
		log.Printf("[INFO] Getting data for casb dlp rule rule id: %d and type: %s\n", id, ruleType)
		res, err := casb_dlp_rules.GetByRuleID(ctx, service, ruleType, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, nameOk := d.Get("name").(string)
	if resp == nil && nameOk && name != "" {
		log.Printf("[INFO] Getting data for casb dlp rule rule: %s and type: %s\n", name, ruleType)
		res, err := casb_dlp_rules.GetByRuleType(ctx, service, ruleType)
		if err != nil {
			return diag.FromErr(err)
		}

		// Look for the rule with the specified name
		for _, rule := range res {
			if rule.Name == name {
				resp = &rule
				break
			}
		}
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("id", resp.ID)
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
		_ = d.Set("number_of_internal_collaborators", resp.NumberOfInternalCollaborators)
		_ = d.Set("number_of_external_collaborators", resp.NumberOfExternalCollaborators)
		_ = d.Set("recipient", resp.Recipient)
		_ = d.Set("quarantine_location", resp.QuarantineLocation)
		_ = d.Set("access_control", resp.AccessControl)
		_ = d.Set("watermark_delete_old_version", resp.WatermarkDeleteOldVersion)
		_ = d.Set("include_criteria_domain_profile", resp.IncludeCriteriaDomainProfile)
		_ = d.Set("include_email_recipient_profile", resp.IncludeEmailRecipientProfile)
		_ = d.Set("without_content_inspection", resp.WithoutContentInspection)
		_ = d.Set("include_entity_groups", resp.IncludeEntityGroups)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("groups", flattenIDNameExtensions(resp.Groups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("departments", flattenIDNameExtensions(resp.Departments)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("users", flattenIDNameExtensions(resp.Users)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dlp_engines", flattenIDExtensions(resp.DLPEngines)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("included_domain_profiles", flattenIDNameExtensions(resp.IncludedDomainProfiles)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("excluded_domain_profiles", flattenIDNameExtensions(resp.ExcludedDomainProfiles)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("criteria_domain_profiles", flattenIDNameExtensions(resp.CriteriaDomainProfiles)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("email_recipient_profiles", flattenIDNameExtensions(resp.EmailRecipientProfiles)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("entity_groups", flattenIDNameExtensions(resp.EntityGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("object_types", flattenIDNameExtensions(resp.ObjectTypes)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("cloud_app_tenants", flattenIDNameExtensions(resp.CloudAppTenants)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("buckets", flattenIDNameExtensions(resp.Buckets)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("casb_tombstone_template", flattenCustomIDNameSet(resp.CasbTombstoneTemplate)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("casb_email_label", flattenCustomIDNameSet(resp.CasbEmailLabel)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("tag", flattenCustomIDNameSet(resp.Tag)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("watermark_profile", flattenCustomIDNameSet(resp.WatermarkProfile)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("zscaler_incident_receiver", flattenCustomIDNameSet(resp.ZscalerIncidentReceiver)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("auditor_notification", flattenCustomIDNameSet(resp.AuditorNotification)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any cloud application rule with name '%s' or id '%d'", name, id))
	}

	return nil
}
