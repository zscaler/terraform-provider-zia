package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol/custom_file_types"
)

func dataSourceCustomFileTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomFileTypesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier for the custom file type.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the custom file type.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the custom file type, if any.",
			},
			"extension": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The file type extension. The maximum extension length is 10 characters.",
			},
			"file_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "File type ID. This ID is assigned and maintained for all file types including predefined and custom file types, and this value is different from the custom file type ID.",
			},
		},
	}
}

func dataSourceCustomFileTypesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *custom_file_types.CustomFileTypes
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for custom file type id: %d\n", id)
		res, err := custom_file_types.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for custom file type: %s\n", name)
		res, err := custom_file_types.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("id", resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("extension", resp.Extension)
		_ = d.Set("file_type_id", resp.FileTypeID)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any custom file type with name '%s' or id '%d'", name, id))
	}

	return nil
}
