package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol"
)

func dataSourceFileTypeCategories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFileTypeCategoriesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "File type ID",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "File type description",
			},
			"parent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Parent category of the file type",
			},
			"enums": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Enum values to filter file types for specific policy categories. Valid values: ZSCALERDLP, EXTERNALDLP, FILETYPECATEGORYFORFILETYPECONTROL",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclude_custom_file_types": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value specifying whether custom file types must be excluded from the list or not",
			},
		},
	}
}

func dataSourceFileTypeCategoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Build filter options from schema
	var opts *filetypecontrol.GetFileTypeCategoriesFilterOptions
	var hasFilters bool

	// Get enums from schema
	if enumsInterface, ok := d.GetOk("enums"); ok {
		enumsList := enumsInterface.([]interface{})
		enums := make([]string, len(enumsList))
		for i, v := range enumsList {
			enums[i] = v.(string)
		}
		if len(enums) > 0 {
			if opts == nil {
				opts = &filetypecontrol.GetFileTypeCategoriesFilterOptions{}
			}
			opts.Enums = enums
			hasFilters = true
		}
	}

	// Get exclude_custom_file_types from schema
	if excludeCustom, ok := d.GetOk("exclude_custom_file_types"); ok {
		exclude := excludeCustom.(bool)
		if opts == nil {
			opts = &filetypecontrol.GetFileTypeCategoriesFilterOptions{}
		}
		opts.ExcludeCustomFileTypes = &exclude
		hasFilters = true
	}

	// Only create opts if we have filters
	if !hasFilters {
		opts = nil
	}

	// Get all file type categories with optional filters
	log.Printf("[INFO] Getting file type categories\n")
	fileTypeCategories, err := filetypecontrol.GetFileTypeCategories(ctx, service, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(fileTypeCategories) == 0 {
		return diag.FromErr(fmt.Errorf("no file type categories found"))
	}

	var resp *filetypecontrol.FileTypeCategory
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Searching for file type category by id: %d\n", id)
		for i := range fileTypeCategories {
			if fileTypeCategories[i].ID == id {
				resp = &fileTypeCategories[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("couldn't find any file type category with id '%d'", id))
		}
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Searching for file type category by name: %s\n", name)
		for i := range fileTypeCategories {
			if strings.EqualFold(fileTypeCategories[i].Name, name) {
				resp = &fileTypeCategories[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("couldn't find any file type category with name '%s'", name))
		}
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("parent", resp.Parent)

	return nil
}
