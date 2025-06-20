package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/saas_security_api"
)

func dataSourceCasbEmailLabel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCasbEmailLabelRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "SaaS Security API email label ID",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "SaaS Security API email label name",
			},
			"label_deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Color to apply to the email label",
			},
		},
	}
}

func dataSourceCasbEmailLabelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var matched *saas_security_api.CasbEmailLabel

	labels, err := saas_security_api.GetCasbEmailLabelLite(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve email label: %w", err))
	}

	id, idOk := getIntFromResourceData(d, "id")
	name, _ := d.Get("name").(string)

	for _, label := range labels {
		// Match by ID (if provided)
		if idOk && label.ID == id {
			matched = &label
			break
		}

		// Match by name (if ID not matched or not provided)
		if name != "" && label.Name == name {
			matched = &label
			break
		}
	}

	if matched == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any email label with name '%s' or id '%d'", name, id))
	}

	// Populate the schema fields
	d.SetId(fmt.Sprintf("%d", matched.ID))
	_ = d.Set("name", matched.Name)
	_ = d.Set("label_deleted", matched.LabelDeleted)

	return nil
}
