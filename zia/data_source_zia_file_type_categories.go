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
	var enumSpecified bool

	if enumsValue, ok := d.GetOk("enums"); ok {
		enumStr := enumsValue.(string)
		if enumStr != "" {
			opts = &filetypecontrol.GetFileTypeCategoriesFilterOptions{}
			opts.Enums = []string{enumStr}
			enumSpecified = true
		}
	}

	if excludeCustom, ok := d.GetOk("exclude_custom_file_types"); ok {
		exclude := excludeCustom.(bool)
		if opts == nil {
			opts = &filetypecontrol.GetFileTypeCategoriesFilterOptions{}
		}
		opts.ExcludeCustomFileTypes = &exclude
	}

	// Check if user wants a specific category by id or name
	id, hasID := getIntFromResourceData(d, "id")
	name, _ := d.Get("name").(string)

	// When searching by id or name without a specific enum filter, query all
	// enum scopes because the API requires an enum to return results and a
	// given file type category may only exist within a particular scope.
	if (hasID || name != "") && !enumSpecified {
		log.Printf("[INFO] Searching file type categories across all enum scopes\n")
		allEnums := []string{"ZSCALERDLP", "EXTERNALDLP", "FILETYPECATEGORYFORFILETYPECONTROL"}
		var allCategories []filetypecontrol.FileTypeCategory
		seen := make(map[int]bool)
		for _, enumVal := range allEnums {
			searchOpts := &filetypecontrol.GetFileTypeCategoriesFilterOptions{
				Enums: []string{enumVal},
			}
			if excludeCustom, ok := d.GetOk("exclude_custom_file_types"); ok {
				exclude := excludeCustom.(bool)
				searchOpts.ExcludeCustomFileTypes = &exclude
			}
			categories, err := filetypecontrol.GetFileTypeCategories(ctx, service, searchOpts)
			if err != nil {
				log.Printf("[WARN] Failed to get file type categories for enum %s: %v", enumVal, err)
				continue
			}
			for _, cat := range categories {
				if !seen[cat.ID] {
					seen[cat.ID] = true
					allCategories = append(allCategories, cat)
				}
			}
		}

		if len(allCategories) == 0 {
			return diag.FromErr(fmt.Errorf("no file type categories found across any enum scope"))
		}

		var resp *filetypecontrol.FileTypeCategory
		if hasID {
			log.Printf("[INFO] Searching for file type category by id: %d\n", id)
			for i := range allCategories {
				if allCategories[i].ID == id {
					resp = &allCategories[i]
					break
				}
			}
			if resp == nil {
				return diag.FromErr(fmt.Errorf("couldn't find any file type category with id '%d'", id))
			}
		} else {
			log.Printf("[INFO] Searching for file type category by name: %s\n", name)
			for i := range allCategories {
				if strings.EqualFold(allCategories[i].Name, name) {
					resp = &allCategories[i]
					break
				}
			}
			if resp == nil {
				return diag.FromErr(fmt.Errorf("couldn't find any file type category with name '%s'", name))
			}
		}

		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("id", resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("parent", resp.Parent)
		_ = d.Set("categories", []interface{}{})
		return nil
	}

	// Fetch categories with the specified filter (or nil for list-all mode)
	log.Printf("[INFO] Getting file type categories\n")
	fileTypeCategories, err := filetypecontrol.GetFileTypeCategories(ctx, service, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(fileTypeCategories) == 0 {
		return diag.FromErr(fmt.Errorf("no file type categories found"))
	}

	// If id or name is provided with a specific enum, find the single result
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
		} else {
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

		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("id", resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("parent", resp.Parent)
		_ = d.Set("categories", []interface{}{})
		return nil
	}

	// List mode: return all matching categories
	log.Printf("[INFO] Returning %d file type categories\n", len(fileTypeCategories))

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

	hashCode := schema.HashString(idBuilder.String())
	d.SetId(strconv.Itoa(hashCode))

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
