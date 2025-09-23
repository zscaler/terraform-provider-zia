package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/workloadgroups"
)

func resourceWorkloadGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkloadGroupsCreate,
		ReadContext:   resourceWorkloadGroupsRead,
		UpdateContext: resourceWorkloadGroupsUpdate,
		DeleteContext: resourceWorkloadGroupsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("group_id", idInt)
				} else {
					resp, err := workloadgroups.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("group_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier assigned to the workload group",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "A unique identifier assigned to the workload group",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the workload group",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the workload group",
			},
			"expression_json": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression_containers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"ANY",
											"VPC",
											"SUBNET",
											"VM",
											"ENI",
											"ATTR",
										}, false),
									},
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"AND",
											"OR",
											"OPEN_PARENTHESES",
											"CLOSE_PARENTHESES",
										}, false),
									},
									"tag_container": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tags": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"operator": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														"AND",
														"OR",
														"OPEN_PARENTHESES",
														"CLOSE_PARENTHESES",
													}, false),
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
		},
	}
}

func resourceWorkloadGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandWorkloadGroups(d)
	log.Printf("[INFO] Creating ZIA workload group\n%+v\n", req)

	resp, _, err := workloadgroups.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA workload group request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)

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

	return resourceWorkloadGroupsRead(ctx, d, meta)
}

func resourceWorkloadGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no workload group id is set"))
	}
	resp, err := workloadgroups.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia workload group %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia workload group:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("group_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)

	// Flatten expression_json if present
	if err := d.Set("expression_json", flattenExpressionJSON(&resp.WorkloadTagExpression)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWorkloadGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] workload group ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia workload group ID: %v\n", id)
	req := expandWorkloadGroups(d)
	if _, err := workloadgroups.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := workloadgroups.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceWorkloadGroupsRead(ctx, d, meta)
}

func resourceWorkloadGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] workload group ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia workload group ID: %v\n", (d.Id()))

	if _, err := workloadgroups.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia workload group deleted")

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandWorkloadGroups(d *schema.ResourceData) workloadgroups.WorkloadGroup {
	id, _ := getIntFromResourceData(d, "group_id")

	result := workloadgroups.WorkloadGroup{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Expand expression_json if provided
	expressionJSON := expandExpressionJSON(d)
	if expressionJSON != nil {
		result.WorkloadTagExpression = *expressionJSON
	}

	return result
}

func expandExpressionJSON(d *schema.ResourceData) *workloadgroups.WorkloadTagExpression {
	expressionJSONInterface, ok := d.GetOk("expression_json")
	if !ok {
		return nil
	}

	expressionJSONList, ok := expressionJSONInterface.([]interface{})
	if !ok || len(expressionJSONList) == 0 {
		return nil
	}

	expressionJSON := &workloadgroups.WorkloadTagExpression{}
	expressionJSONData := expressionJSONList[0].(map[string]interface{})

	// Expand expression_containers
	if containersInterface, exists := expressionJSONData["expression_containers"]; exists {
		containersList := containersInterface.([]interface{})
		expressionContainers := make([]workloadgroups.ExpressionContainer, len(containersList))

		for i, containerInterface := range containersList {
			containerData := containerInterface.(map[string]interface{})
			container := workloadgroups.ExpressionContainer{
				TagType:  containerData["tag_type"].(string),
				Operator: containerData["operator"].(string),
			}

			// Expand tag_container
			if tagContainerInterface, exists := containerData["tag_container"]; exists {
				tagContainerList := tagContainerInterface.([]interface{})
				if len(tagContainerList) > 0 {
					tagContainerData := tagContainerList[0].(map[string]interface{})
					tagContainer := workloadgroups.TagContainer{
						Operator: tagContainerData["operator"].(string),
					}

					// Expand tags
					if tagsInterface, exists := tagContainerData["tags"]; exists {
						tagsList := tagsInterface.([]interface{})
						tags := make([]workloadgroups.Tags, len(tagsList))

						for k, tagInterface := range tagsList {
							tagData := tagInterface.(map[string]interface{})
							tags[k] = workloadgroups.Tags{
								Key:   tagData["key"].(string),
								Value: tagData["value"].(string),
							}
						}
						tagContainer.Tags = tags
					}

					container.TagContainer = tagContainer
				}
			}

			expressionContainers[i] = container
		}
		expressionJSON.ExpressionContainers = expressionContainers
	}

	return expressionJSON
}

func flattenExpressionJSON(expressionJSON *workloadgroups.WorkloadTagExpression) []interface{} {
	if expressionJSON == nil {
		return nil
	}

	result := map[string]interface{}{}

	// Flatten expression_containers
	if expressionJSON.ExpressionContainers != nil {
		containers := make([]interface{}, len(expressionJSON.ExpressionContainers))

		for i, container := range expressionJSON.ExpressionContainers {
			containerMap := map[string]interface{}{
				"tag_type": container.TagType,
				"operator": container.Operator,
			}

			// Flatten tag_container
			tagContainerMap := map[string]interface{}{
				"operator": container.TagContainer.Operator,
			}

			// Flatten tags
			if container.TagContainer.Tags != nil {
				tags := make([]interface{}, len(container.TagContainer.Tags))

				for k, tag := range container.TagContainer.Tags {
					tags[k] = map[string]interface{}{
						"key":   tag.Key,
						"value": tag.Value,
					}
				}
				tagContainerMap["tags"] = tags
			}

			containerMap["tag_container"] = []interface{}{tagContainerMap}
			containers[i] = containerMap
		}
		result["expression_containers"] = containers
	}

	return []interface{}{result}
}
