package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
)

func dataSourceDlpWebRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDlpWebRulesRead,
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
			"order": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the DLP policy rule.",
			},
			"access_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access privilege for this DLP policy rule based on the admin's state.",
			},
			"protocols": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The protocol criteria specified for the DLP policy rule.",
			},
			"rank": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations to which the DLP policy rule must be applied.",
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
			"location_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the DLP policy rule must be applied.",
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
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of groups to which the DLP policy rule must be applied.",
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
			"departments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of departments to which the DLP policy rule must be applied.",
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
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of users to which the DLP policy rule must be applied.",
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
			"url_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of URL categories to which the DLP policy rule must be applied.",
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
			"dlp_engines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of DLP engines to which the DLP policy rule must be applied.",
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
			"file_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of file types to which the DLP policy rule must be applied.",
			},
			"cloud_applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of cloud applications to which the DLP policy rule must be applied.",
			},
			"min_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum file size (in KB) used for evaluation of the DLP policy rule.",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action taken when traffic matches the DLP policy rule criteria.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enables or disables the DLP policy rule.",
			},
			"severity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the severity selected for the DLP rule violation",
			},
			"parent_rule": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the parent rule under which an exception rule is added.",
			},
			"sub_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of exception rules added to a parent rule",
			},
			"time_windows": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of time windows to which the DLP policy rule must be applied.",
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
			"external_auditor_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of an external auditor to whom DLP email notifications are sent.",
			},
			"match_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The match only criteria for DLP engines.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the DLP policy rule was last modified.",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the DLP policy rule last.",
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
			"without_content_inspection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of rule labels associated to the DLP policy rule.",
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
			"source_ip_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of Source IP Groups associated to the DLP policy rule.",
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
			"dlp_download_scan_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"zcc_notifications_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"zscaler_incident_receiver": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether a Zscaler Incident Receiver is associated to the DLP policy rule.",
			},
			"excluded_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The name-ID pairs of the groups that are excluded from the DLP policy rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
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
			"excluded_departments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The name-ID pairs of the departments that are excluded from the DLP policy rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
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
			"excluded_users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The name-ID pairs of the users that are excluded from the DLP policy rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
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
			"included_domain_profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The name-ID pairs of the users that are excluded from the DLP policy rule.",
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
							Description: "The name of the workload group",
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
			"workload_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of preconfigured workload groups to which the policy must be applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier assigned to the workload group",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the workload group",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the workload group",
						},
						"last_modified_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_modified_by": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
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
					},
				},
			},
		},
	}
}

func dataSourceDlpWebRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *dlp_web_rules.WebDLPRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for dlp web rule id: %d\n", id)
		res, err := dlp_web_rules.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for dlp web rule: %s\n", name)
		res, err := dlp_web_rules.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("order", resp.Order)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("access_control", resp.AccessControl)
		_ = d.Set("protocols", resp.Protocols)
		_ = d.Set("description", resp.Description)
		_ = d.Set("file_types", resp.FileTypes)
		_ = d.Set("cloud_applications", resp.CloudApplications)
		_ = d.Set("state", resp.State)
		_ = d.Set("min_size", resp.MinSize)
		_ = d.Set("action", resp.Action)
		_ = d.Set("severity", resp.Severity)
		_ = d.Set("parent_rule", resp.ParentRule)
		_ = d.Set("match_only", resp.MatchOnly)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)
		_ = d.Set("external_auditor_email", resp.ExternalAuditorEmail)
		_ = d.Set("without_content_inspection", resp.WithoutContentInspection)
		_ = d.Set("zscaler_incident_receiver", resp.ZscalerIncidentReceiver)
		_ = d.Set("zcc_notifications_enabled", resp.ZCCNotificationsEnabled)
		_ = d.Set("dlp_download_scan_enabled", resp.DLPDownloadScanEnabled)

		// Flatten sub_rules
		subRuleIDs := make([]interface{}, len(resp.SubRules))
		for i, subRule := range resp.SubRules {
			subRuleIDs[i] = strconv.Itoa(subRule.ID)
		}
		if err := d.Set("sub_rules", subRuleIDs); err != nil {
			return diag.FromErr(fmt.Errorf("error setting sub_rules: %s", err))
		}

		if err := d.Set("locations", flattenIDExtensions(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("included_domain_profiles", flattenIDExtensions(resp.IncludedDomainProfiles)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("location_groups", flattenIDExtensions(resp.LocationGroups)); err != nil {
			return diag.FromErr(err)
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

		if err := d.Set("url_categories", flattenIDExtensions(resp.URLCategories)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dlp_engines", flattenIDExtensions(resp.DLPEngines)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("time_windows", flattenIDExtensions(resp.TimeWindows)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("last_modified_by", flattenIDExtensionsList(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("labels", flattenIDExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("source_ip_groups", flattenIDExtensions(resp.SourceIpGroups)); err != nil {
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
		if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting workload_groups: %s", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any web dlp rule with name '%s' or id '%d'", name, id))
	}

	return nil
}
