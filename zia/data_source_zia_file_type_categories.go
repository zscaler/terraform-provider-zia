package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
				Description: "File type ID. If specified, returns a single category matching this ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "File type description. If specified, returns a single category matching this name.",
			},
			"parent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Parent category of the file type (only set when querying by id or name)",
			},
			"enums": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enum value to filter file types for specific policy categories. Valid values: ZSCALERDLP, EXTERNALDLP, FILETYPECATEGORYFORFILETYPECONTROL",
			},
			"exclude_custom_file_types": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value specifying whether custom file types must be excluded from the list or not",
			},
			"categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of file type categories matching the filter criteria (returned when querying by enums without id or name)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File type category ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File type category name",
						},
						"parent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent category of the file type",
						},
					},
				},
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
	if enumsValue, ok := d.GetOk("enums"); ok {
		enumStr := enumsValue.(string)
		if enumStr != "" {
			if opts == nil {
				opts = &filetypecontrol.GetFileTypeCategoriesFilterOptions{}
			}
			opts.Enums = []string{enumStr}
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

	// Check if user wants a specific category by id or name
	id, hasID := getIntFromResourceData(d, "id")
	name, _ := d.Get("name").(string)

	// If id or name is provided, return single result (backward compatible)
	if hasID || name != "" {
		var resp *filetypecontrol.FileTypeCategory

		if hasID {
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
		} else if name != "" {
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

		// Set single result
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("id", resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("parent", resp.Parent)
		// Clear categories list when returning single result
		_ = d.Set("categories", []interface{}{})

		return nil
	}

	// If only filters provided (enums/exclude_custom_file_types), return all matching categories
	log.Printf("[INFO] Returning %d file type categories\n", len(fileTypeCategories))

	// Generate a stable ID for this list query based on the query parameters
	var idBuilder strings.Builder
	if enumsValue, ok := d.GetOk("enums"); ok {
		idBuilder.WriteString(enumsValue.(string))
	}
	if idBuilder.Len() == 0 {
		idBuilder.WriteString("all")
	}
	if excludeCustom, ok := d.GetOk("exclude_custom_file_types"); ok {
		idBuilder.WriteString(fmt.Sprintf("-exclude_%v", excludeCustom))
	}

	// Use hash code to generate a stable numeric ID
	hashCode := schema.HashString(idBuilder.String())
	d.SetId(strconv.Itoa(hashCode))

	// Flatten the categories list
	categoriesList := make([]interface{}, len(fileTypeCategories))
	for i, cat := range fileTypeCategories {
		categoriesList[i] = map[string]interface{}{
			"id":     cat.ID,
			"name":   cat.Name,
			"parent": cat.Parent,
		}
	}

	_ = d.Set("categories", categoriesList)

	return nil
}
