package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/workloadgroups"
)

func dataSourceWorkloadGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkloadGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "A unique identifier assigned to the workload group",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the workload group",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the workload group",
			},
			"expression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The workload group expression containing tag types, tags, and their relationships.",
			},
			"expression_json": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression_containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_container": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tags": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"value": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
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
						"external_id": {
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
	}
}

func dataSourceWorkloadGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *workloadgroups.WorkloadGroup
	var searchCriteria string

	// Check if searching by ID
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting time window by id: %d\n", id)
		searchCriteria = fmt.Sprintf("id=%d", id)

		// Get all time windows and find the one with matching ID
		allTimeWindows, err := workloadgroups.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, tw := range allTimeWindows {
			if tw.ID == id {
				resp = &tw
				break
			}
		}
	}

	// Check if searching by name (only if ID search didn't find anything)
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting time window by name: %s\n", name)
		searchCriteria = fmt.Sprintf("name=%s", name)

		res, err := workloadgroups.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("expression", resp.Expression)

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

		expressionJson := flattenWorkloadTagExpression(resp.WorkloadTagExpression)
		if err := d.Set("expression_json", expressionJson); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any workload groups with %s", searchCriteria))
	}

	return nil
}

// Flatten the WorkloadTagExpression into a format suitable for Terraform schema
func flattenWorkloadTagExpression(expression workloadgroups.WorkloadTagExpression) []interface{} {
	if len(expression.ExpressionContainers) == 0 {
		return nil
	}

	var flattenedExpression []map[string]interface{}
	for _, container := range expression.ExpressionContainers {
		flattenedContainer := map[string]interface{}{
			"tag_type":      container.TagType,
			"operator":      container.Operator,
			"tag_container": []interface{}{flattenTagContainer(container.TagContainer)},
		}

		flattenedExpression = append(flattenedExpression, flattenedContainer)
	}
	return []interface{}{map[string]interface{}{"expression_containers": flattenedExpression}}
}

// Flatten the TagContainer structure
func flattenTagContainer(container workloadgroups.TagContainer) map[string]interface{} {
	var flattenedTags []map[string]interface{}
	for _, tag := range container.Tags {
		flattenedTag := map[string]interface{}{
			"key":   tag.Key,
			"value": tag.Value,
		}
		flattenedTags = append(flattenedTags, flattenedTag)
	}

	return map[string]interface{}{
		"tags":     flattenedTags,
		"operator": container.Operator,
	}
}
