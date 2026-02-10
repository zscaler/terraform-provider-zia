package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlcategories"
)

func dataSourceURLCategories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceURLCategoriesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"configured_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category_group": {
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
			"super_category": {
				Type:     schema.TypeString,
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
							Optional: true,
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
				Optional: true,
				Default:  "ALL",
				ValidateFunc: validation.StringInSlice([]string{
					"ALL",
					"URL_CATEGORY",
					"TLD_CATEGORY",
				}, false),
				Description: "Type of URL categories to retrieve. Valid values: ALL (default - includes all types), URL_CATEGORY, TLD_CATEGORY",
			},
			"url_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EXACT",
					"REGEX",
				}, false),
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
			"val": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
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
			"regex_patterns": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"regex_patterns_retaining_parent_category": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceURLCategoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *urlcategories.URLCategory
	var err error

	id, _ := d.Get("id").(string)
	name, _ := d.Get("configured_name").(string)
	categoryType, _ := d.Get("type").(string)

	// Default to "ALL" if not specified
	if categoryType == "" {
		categoryType = "ALL"
	}

	// Ensure either ID or name is provided
	if id == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'configured_name' must be specified"))
	}

	if id != "" {
		log.Printf("[INFO] Getting URL categories by ID: %s\n", id)
		resp, err = urlcategories.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if name != "" {
		log.Printf("[INFO] Getting URL categories by name: %s (type: %s)\n", name, categoryType)
		resp, err = urlcategories.GetCustomURLCategories(ctx, service, name, true, true, categoryType)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// After attempting to fetch, check if resp is still nil, indicating no data was found.
	if resp == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any URL category with ID '%s' or name '%s'", id, name))
	}

	// Set the data source fields with the response
	d.SetId(resp.ID)
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
	_ = d.Set("regex_patterns", resp.RegexPatterns)
	_ = d.Set("regex_patterns_retaining_parent_category", resp.RegexPatternsRetainingParentCategory)
	_ = d.Set("url_type", resp.UrlType)
	_ = d.Set("category_group", resp.CategoryGroup)
	_ = d.Set("super_category", resp.SuperCategory)

	if err := d.Set("scopes", flattenScopes(resp)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url_keyword_counts", flattenUrlKeywordCounts(resp.URLKeywordCounts)); err != nil {
		return diag.FromErr(err)
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
