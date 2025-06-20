package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/saas_security_api"
)

func dataSourceCasbTombstoneTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCasbTombstoneTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Tombstone file template ID",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Tombstone file template name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The text that is included in the tombstone file",
			},
		},
	}
}

func dataSourceCasbTombstoneTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var matched *saas_security_api.QuarantineTombstoneLite

	templates, err := saas_security_api.GetQuarantineTombstoneLite(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve quarantine tombstone template: %w", err))
	}

	id, idOk := getIntFromResourceData(d, "id")
	name, _ := d.Get("name").(string)

	for _, template := range templates {
		if idOk && template.ID == id {
			matched = &template
			break
		}

		if name != "" && template.Name == name {
			matched = &template
			break
		}
	}

	if matched == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any quarantine tombstone template with name '%s' or id '%d'", name, id))
	}

	d.SetId(fmt.Sprintf("%d", matched.ID))
	_ = d.Set("name", matched.Name)
	_ = d.Set("description", matched.Description)

	return nil
}
