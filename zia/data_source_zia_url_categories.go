package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlcategories"
)

func dataSourceURLCategories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceURLCategoriesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configured_name": {
				Type:     schema.TypeString,
				Optional: true,
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
										Type:     schema.TypeString,
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
			"val": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"custom_urls_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"urls_retaining_parent_category_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceURLCategoriesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *urlcategories.URLCategory
	id, _ := d.Get("id").(string)
	if resp == nil && id != "" {
		log.Printf("[INFO] Getting url categories : %s\n", id)
		res, err := zClient.urlcategories.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}

	name, _ := d.Get("configured_name").(string)
	if resp == nil && id != "" {
		log.Printf("[INFO] Getting url categories : %s\n", name)
		res, err := zClient.urlcategories.GetCustomURLCategories(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf(resp.ID))
		_ = d.Set("configured_name", resp.ConfiguredName)
		_ = d.Set("urls", resp.Urls)
		_ = d.Set("db_categorized_urls", resp.DBCategorizedUrls)
		_ = d.Set("custom_category", resp.CustomCategory)
		_ = d.Set("editable", resp.Editable)
		_ = d.Set("description", resp.Description)
		_ = d.Set("type", resp.Type)
		_ = d.Set("val", resp.Val)
		_ = d.Set("custom_urls_count", resp.CustomUrlsCount)
		_ = d.Set("urls_retaining_parent_category_count", resp.UrlsRetainingParentCategoryCount)

		if err := d.Set("scopes", flattenScopes(resp)); err != nil {
			return err
		}

		if err := d.Set("url_keyword_counts", flattenUrlKeywordCounts(resp.URLKeywordCounts)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any category or custom category with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenScopes(scopes *urlcategories.URLCategory) []interface{} {
	scope := make([]interface{}, len(scopes.Scopes))
	for i, val := range scopes.Scopes {
		scope[i] = map[string]interface{}{
			"type":                        val.Type,
			"scope_group_member_entities": flattenScopeGroupMemberEntities(val.ScopeGroupMemberEntities),
			"scope_entities":              flattenScopeEntities(val.ScopeEntities),
		}
	}

	return scope
}

func flattenScopeGroupMemberEntities(scopeGroupMember []urlcategories.ScopeGroupMemberEntities) []interface{} {
	scopeGroups := make([]interface{}, len(scopeGroupMember))
	for i, val := range scopeGroupMember {
		scopeGroups[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return scopeGroups
}

func flattenScopeEntities(scopesEntities []urlcategories.ScopeEntities) []interface{} {
	Entity := make([]interface{}, len(scopesEntities))
	for i, val := range scopesEntities {
		Entity[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return Entity
}

func flattenUrlKeywordCounts(urlKeywords urlcategories.URLKeywordCounts) interface{} {
	return []map[string]interface{}{
		{
			"total_url_count":             urlKeywords.TotalURLCount,
			"retain_parent_url_count":     urlKeywords.RetainParentURLCount,
			"total_keyword_count":         urlKeywords.TotalKeywordCount,
			"retain_parent_keyword_count": urlKeywords.RetainParentKeywordCount,
		},
	}
}
