package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol"
)

func dataSourceFileTypeControlRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFileTypeControlRulesRead,
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
				Description: "The File Type Control policy rule name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the File Type Control rule.",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The rule order of execution for the  File Type Control rule with respect to other rules.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enables or disables the File Type Control rule.",
			},
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Admin rank of the admin who creates this rule",
			},
			"filtering_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action taken when traffic matches policy. This field is not applicable to the Lite API.",
			},
			"operation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "File operation performed. This field is not applicable to the Lite API.",
			},
			"access_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access privilege of this rule based on the admin's RBA state. Ignored if the request is POST or PUT.",
			},
			"time_quota": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Time quota in minutes, after which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to 'BLOCK', this field is not applicable.",
			},
			"size_quota": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size quota in KB beyond which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to 'BLOCK', this field is not applicable.",
			},
			"min_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimum file size (in KB) used for evaluation of the FTP rule",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum file size (in KB) used for evaluation of the FTP rule",
			},
			"capture_pcap": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that indicates whether packet capture (PCAP) is enabled or not",
			},
			"active_content": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to check whether a file has active content or not",
			},
			"unscannable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to check whether a file has active content or not",
			},
			"cloud_applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of cloud applications to which the File Type Control rule must be applied.",
			},
			"url_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of URL categories for which the policy must be applied. If not set, policy is applied for all URL categories.",
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
			"device_trust_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of device trust levels for which the rule must be applied. While the High Trust, Medium Trust, or Low Trust evaluation is applicable only to Zscaler Client Connector traffic, Unknown evaluation applies to all traffic.",
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

func dataSourceFileTypeControlRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *filetypecontrol.FileTypeRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for file type control rule id: %d\n", id)
		res, err := filetypecontrol.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data file type control rule : %s\n", name)
		res, err := filetypecontrol.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("order", resp.Order)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("state", resp.State)
		_ = d.Set("access_control", resp.AccessControl)
		_ = d.Set("filtering_action", resp.FilteringAction)
		_ = d.Set("operation", resp.Operation)
		_ = d.Set("time_quota", resp.TimeQuota)
		_ = d.Set("size_quota", resp.SizeQuota)
		_ = d.Set("max_size", resp.MaxSize)
		_ = d.Set("min_size", resp.MinSize)
		_ = d.Set("capture_pcap", resp.CapturePCAP)
		_ = d.Set("active_content", resp.ActiveContent)
		_ = d.Set("unscannable", resp.Unscannable)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)
		_ = d.Set("cloud_applications", resp.CloudApplications)
		_ = d.Set("file_types", resp.FileTypes)
		_ = d.Set("protocols", resp.Protocols)
		_ = d.Set("url_categories", resp.URLCategories)
		_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)

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

		if err := d.Set("time_windows", flattenIDNameExtensions(resp.TimeWindows)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("device_groups", flattenIDNameExtensions(resp.DeviceGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("devices", flattenIDNameExtensions(resp.Devices)); err != nil {
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
