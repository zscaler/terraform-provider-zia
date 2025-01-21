package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_rules"
)

func dataSourceSandboxRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSandboxRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Sandbox rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the Sandbox rule",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enables or disables the sandbox rules.",
			},
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Admin rank of the admin who creates this rule",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The rule order of execution for the  sandbox rules with respect to other rules.",
			},
			"ba_rule_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action configured for the rule that must take place if the traffic matches the rule criteria.",
			},
			"first_time_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether a First-Time Action is specifically configured for the rule",
			},
			"first_time_operation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action that must take place when users download unknown files for the first time",
			},
			"ml_action_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When set to true, this indicates that 'Machine Learning Intelligence Action' checkbox has been checked on",
			},
			"by_threat_score": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "",
			},
			"ba_policy_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The threat categories to which the rule applies",
			},
			"file_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "File type categories for which the policy is applied. If not set, the rule is applied across all file types.",
			},
			"protocols": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Protocol for the given rule. This field is not applicable to the Lite API.",
			},
			"url_categories": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The URL categories to which the rule applies",
			},
			"locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of locations for the which policy must be applied. If not set, policy is applied for all locations.",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
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
				Description: "Name-ID pairs of locations groups for which rule must be applied.",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
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
				Description: "Name-ID pairs of department for which the rule is applied. If not set, rule will be applied for all departments.",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
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
				Description: "Name-ID pairs of groups for which the policy must be applied. If not set, policy is applied for all groups.",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
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
				Description: "Name-ID pairs of users for the which policy must be applied. If not set, user criteria is not considered for policy enforcement.",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"time_windows": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "This is an immutable reference to an entity that mainly consists of id and name.",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the rule was last modified. Ignored if the request is POST or PUT.",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Admin user that last modified the rule. Ignored if the request is POST or PUT.",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
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
				Description: "This is an immutable reference to an entity that mainly consists of id and name.",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"device_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of device groups for which the rule is applied",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"devices": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of device for which the rule is applied",
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
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"zpa_app_segments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier assigned to the Application Segment",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Application Segment",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the external ID. Applicable only when this reference is of an external entity.",
						},
					},
				},
			},
		},
	}
}

func dataSourceSandboxRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *sandbox_rules.SandboxRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for sandbox rule id: %d\n", id)
		res, err := sandbox_rules.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data sandbox rule : %s\n", name)
		res, err := sandbox_rules.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("state", resp.State)
		_ = d.Set("order", resp.Order)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("ba_rule_action", resp.BaRuleAction)
		_ = d.Set("first_time_enable", resp.FirstTimeEnable)
		_ = d.Set("first_time_operation", resp.FirstTimeOperation)
		_ = d.Set("ml_action_enabled", resp.MLActionEnabled)
		_ = d.Set("by_threat_score", resp.ByThreatScore)
		_ = d.Set("ba_policy_categories", resp.BaPolicyCategories)
		_ = d.Set("url_categories", resp.URLCategories)
		_ = d.Set("protocols", resp.Protocols)
		_ = d.Set("file_types", resp.FileTypes)

		if err := d.Set("devices", flattenIDNameExtensions(resp.Devices)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("device_groups", flattenIDNameExtensions(resp.DeviceGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("location_groups", flattenIDNameExtensions(resp.LocationGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("departments", flattenIDNameExtensions(resp.Departments)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("groups", flattenIDNameExtensions(resp.Groups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("users", flattenIDNameExtensions(resp.Users)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("zpa_app_segments", flattenZPAAppSegments(resp.ZPAAppSegments)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any user with name '%s' or id '%d'", name, id))
	}

	return nil
}
