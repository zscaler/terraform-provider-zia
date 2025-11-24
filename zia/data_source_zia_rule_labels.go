package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/rule_labels"
)

func dataSourceRuleLabels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRuleLabelsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the rule label.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The rule label name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule label description.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the rule lable was last modified. This is a read-only field. Ignored by PUT and DELETE requests.",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the rule label last. This is a read-only field. Ignored by PUT requests.",
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
			},
			"created_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that created the rule label. This is a read-only field. Ignored by PUT requests.",
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
			},
			"referenced_rule_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of rules that reference the label.",
			},
		},
	}
}

func dataSourceRuleLabelsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Always fetch all rule labels and search locally
	log.Printf("[INFO] Fetching all rule labels\n")
	allRuleLabels, err := rule_labels.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all rule labels: %s", err))
	}

	log.Printf("[DEBUG] Retrieved %d rule labels\n", len(allRuleLabels))

	var resp *rule_labels.RuleLabels
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	// Search by ID first if provided
	if idProvided {
		log.Printf("[INFO] Searching for rule label by ID: %d\n", id)
		for _, label := range allRuleLabels {
			if label.ID == id {
				resp = &label
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting rule label by ID %d: rule label not found", id))
		}
	}

	// Search by name if not found by ID and name is provided
	if resp == nil && nameProvided && nameStr != "" {
		log.Printf("[INFO] Searching for rule label by name: %s\n", nameStr)
		for _, label := range allRuleLabels {
			if label.Name == nameStr {
				resp = &label
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting rule label by name %s: rule label not found", nameStr))
		}
	}

	// If neither ID nor name provided, or no match found
	if resp == nil {
		if idProvided || (nameProvided && nameStr != "") {
			return diag.FromErr(fmt.Errorf("couldn't find any rule label with name '%s' or id '%d'", nameStr, id))
		}
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	// Set the resource data
	d.SetId(fmt.Sprintf("%d", resp.ID))
	err = d.Set("name", resp.Name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}
	err = d.Set("description", resp.Description)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting description: %s", err))
	}
	err = d.Set("referenced_rule_count", resp.ReferencedRuleCount)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting referenced_rule_count: %s", err))
	}
	err = d.Set("last_modified_time", resp.LastModifiedTime)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting last_modified_time: %s", err))
	}

	if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_by", flattenCreatedBy(resp.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Rule label found: ID=%d, Name=%s\n", resp.ID, resp.Name)
	return nil
}
