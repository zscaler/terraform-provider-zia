package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/rule_labels"
)

func dataSourceRuleLabels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRuleLabelsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
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

func dataSourceRuleLabelsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *rule_labels.RuleLabels
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for rule label id: %d\n", id)
		res, err := zClient.rule_labels.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for rule label name: %s\n", name)
		res, err := zClient.rule_labels.GetRuleLabelByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("referenced_rule_count", resp.ReferencedRuleCount)

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return err
		}

		if err := d.Set("created_by", flattenCreatedBy(resp.CreatedBy)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any rule label name '%s' or id '%d'", name, id)
	}

	return nil
}
