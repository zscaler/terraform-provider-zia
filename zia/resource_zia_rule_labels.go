package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/rule_labels"
)

func resourceRuleLabels() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleLabelsCreate,
		Read:   resourceRuleLabelsRead,
		Update: resourceRuleLabelsUpdate,
		Delete: resourceRuleLabelsDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("rule_label_id", idInt)
				} else {
					resp, err := zClient.rule_labels.GetRuleLabelByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("rule_label_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_label_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
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
							Optional: true,
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
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"referenced_rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceRuleLabelsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandRuleLabels(d)
	log.Printf("[INFO] Creating zia rule labels\n%+v\n", req)

	resp, _, err := zClient.rule_labels.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia rule labels request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_label_id", resp.ID)
	return resourceRuleLabelsRead(d, m)
}

func resourceRuleLabelsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_label_id")
	if !ok {
		return fmt.Errorf("no rule labels id is set")
	}
	resp, err := zClient.rule_labels.Get(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia rule labels %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting zia rule labels:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_label_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("referenced_rule_count", resp.ReferencedRuleCount)

	if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
		return err
	}

	if err := d.Set("created_by", flattenCreatedBy(resp.CreatedBy)); err != nil {
		return err
	}

	return nil
}

func resourceRuleLabelsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_label_id")
	if !ok {
		log.Printf("[ERROR] rule label ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia rule label ID: %v\n", id)
	req := expandRuleLabels(d)
	if _, err := zClient.rule_labels.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := zClient.rule_labels.Update(id, &req); err != nil {
		return err
	}

	return resourceRuleLabelsRead(d, m)
}

func resourceRuleLabelsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_label_id")
	if !ok {
		log.Printf("[ERROR] rule label ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia rule label ID: %v\n", (d.Id()))

	if _, err := zClient.rule_labels.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] zia rule label deleted")
	return nil
}

func expandRuleLabels(d *schema.ResourceData) rule_labels.RuleLabels {
	id, _ := getIntFromResourceData(d, "rule_label_id")
	result := rule_labels.RuleLabels{
		ID:                  id,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		LastModifiedTime:    d.Get("last_modified_time").(int),
		LastModifiedBy:      expandIDNameExtensions(d, "last_modified_by"),
		CreatedBy:           expandIDNameExtensions(d, "created_by"),
		ReferencedRuleCount: d.Get("referenced_rule_count").(int),
	}
	return result
}
