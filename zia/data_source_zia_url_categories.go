package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/urlcategories"
)

func dataSourceURLCategories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceURLCategoriesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"val": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configured_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"keywords": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"keywords_retaining_parent_category": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"db_categorized_urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"custom_category": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope_group_member_entities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_entities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
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
					},
				},
			},
			"editable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url_keyword_counts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_url_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"retain_parent_url_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_keyword_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"retain_parent_keyword_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"custom_urls_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"urls_retaining_parent_category_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ip_ranges": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_ranges_retaining_parent_category": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"custom_ip_ranges_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip_ranges_retaining_parent_category_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceURLCategoriesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.urlcategories

	var resp *urlcategories.URLCategory
	var err error

	id, _ := d.Get("id").(string)
	name, _ := d.Get("name").(string)
	configuredName, _ := d.Get("configured_name").(string)
	val, _ := d.Get("val").(int)

	// Process specific parameters
	if id != "" {
		log.Printf("[INFO] Getting URL categories by ID: %s\n", id)
		resp, err = urlcategories.Get(service, id)
		if err != nil {
			return err
		}
	} else if name != "" {
		log.Printf("[INFO] Getting URL categories by name: %s\n", name)
		resp, err = getPredefinedCategoryByName(service, name)
		if err != nil {
			return err
		}
	} else if configuredName != "" {
		log.Printf("[INFO] Getting URL categories by configured name: %s\n", configuredName)
		resp, err = urlcategories.GetCustomURLCategories(service, configuredName, true, true)
		if err != nil {
			return err
		}
	} else if val != 0 {
		log.Printf("[INFO] Getting URL categories by val: %d\n", val)
		categories, err := urlcategories.GetAll(service)
		if err != nil {
			return err
		}

		// Iterate through all categories and match by val
		for _, category := range categories {
			if category.Val == val {
				resp = &category
				break
			}
		}

		if resp == nil {
			return fmt.Errorf("URL category with val '%d' not found", val)
		}
	} else {
		// No parameters were provided, returning an error
		return fmt.Errorf("either 'id', 'name', 'configured_name', or 'val' must be specified")
	}

	// After attempting to fetch, check if resp is still nil, indicating no data was found.
	if resp == nil {
		return fmt.Errorf("couldn't find any URL category with ID '%s', name '%s', configured name '%s', or val '%d'", id, name, configuredName, val)
	}

	// Set the data source fields with the response
	d.SetId(fmt.Sprintf("%d", resp.Val)) // Use the 'val' as the unique ID
	_ = d.Set("configured_name", resp.ConfiguredName)
	_ = d.Set("keywords", resp.Keywords)
	_ = d.Set("keywords_retaining_parent_category", resp.KeywordsRetainingParentCategory)
	_ = d.Set("urls", resp.Urls)
	_ = d.Set("db_categorized_urls", resp.DBCategorizedUrls)
	_ = d.Set("custom_category", resp.CustomCategory)
	_ = d.Set("editable", resp.Editable)
	_ = d.Set("description", resp.Description)
	_ = d.Set("type", resp.Type)
	_ = d.Set("val", resp.Val)
	_ = d.Set("custom_urls_count", resp.CustomUrlsCount)
	_ = d.Set("urls_retaining_parent_category_count", resp.UrlsRetainingParentCategoryCount)
	_ = d.Set("ip_ranges", resp.IPRanges)
	_ = d.Set("ip_ranges_retaining_parent_category", resp.IPRangesRetainingParentCategory)
	_ = d.Set("custom_ip_ranges_count", resp.CustomIpRangesCount)
	_ = d.Set("ip_ranges_retaining_parent_category_count", resp.IPRangesRetainingParentCategoryCount)

	if err := d.Set("scopes", flattenScopes(resp)); err != nil {
		return err
	}

	if err := d.Set("url_keyword_counts", flattenUrlKeywordCounts(resp.URLKeywordCounts)); err != nil {
		return err
	}

	return nil
}

func flattenScopes(scopes *urlcategories.URLCategory) []interface{} {
	scope := make([]interface{}, len(scopes.Scopes))
	for i, val := range scopes.Scopes {
		scope[i] = map[string]interface{}{
			"type":                        val.Type,
			"scope_group_member_entities": flattenIDNameExtensions(val.ScopeGroupMemberEntities),
			"scope_entities":              flattenIDNameExtensions(val.ScopeEntities),
		}
	}

	return scope
}

func flattenUrlKeywordCounts(urlKeywords *urlcategories.URLKeywordCounts) []interface{} {
	if urlKeywords == nil {
		return nil
	}
	m := map[string]interface{}{
		"total_url_count":             urlKeywords.TotalURLCount,
		"retain_parent_url_count":     urlKeywords.RetainParentURLCount,
		"total_keyword_count":         urlKeywords.TotalKeywordCount,
		"retain_parent_keyword_count": urlKeywords.RetainParentKeywordCount,
	}

	return []interface{}{m}
}

// Helper function to get predefined category by name
func getPredefinedCategoryByName(service *services.Service, name string) (*urlcategories.URLCategory, error) {
	// Fetch all predefined categories using GetAll
	categories, err := urlcategories.GetAll(service) // Correct usage of the GetAll method
	if err != nil {
		return nil, err
	}

	// Iterate over the fetched categories to find the one with the matching name
	for _, category := range categories {
		if category.ID == name { // Comparing with the ID field as it represents the "name" for predefined categories
			return &category, nil
		}
	}

	return nil, fmt.Errorf("predefined category with name '%s' not found", name)
}
