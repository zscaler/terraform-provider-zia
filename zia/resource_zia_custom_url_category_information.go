package zia

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlcategories"
)

func resourceURLCategories() *schema.Resource {
	return &schema.Resource{
		Create:   resourceURLCategoriesCreate,
		Read:     resourceURLCategoriesRead,
		Update:   resourceURLCategoriesUpdate,
		Delete:   resourceURLCategoriesDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configured_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"db_categorized_urls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"custom_category": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope_group_member_entities": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Optional: true,
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
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ORGANIZATION",
								"DEPARTMENT",
								"LOCATION",
								"LOCATION_GROUP",
							}, false),
						},
						"scope_entities": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Optional: true,
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
					},
				},
			},
			"editable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"URL_CATEGORY",
					"TLD_CATEGORY",
					"ALL",
				}, false),
			},
			"url_keyword_counts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_url_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"retain_parent_url_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"total_keyword_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"retain_parent_keyword_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"custom_urls_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"urls_retaining_parent_category_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceURLCategoriesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandURLCategory(d)
	log.Printf("[INFO] Creating zia url category\n%+v\n", req)

	resp, err := zClient.urlcategories.CreateURLCategories(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia url category request. ID: %v\n", resp)

	return resourceURLCategoriesRead(d, m)
}

func resourceURLCategoriesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.urlcategories.Get(d.Id())

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing zia url category %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting url category :\n%+v\n", resp)

	d.SetId(resp.ID)
	_ = d.Set("configured_name", resp.ConfiguredName)
	_ = d.Set("urls", resp.Urls)
	_ = d.Set("db_categorized_urls", resp.DBCategorizedUrls)
	_ = d.Set("custom_category", resp.CustomCategory)
	_ = d.Set("editable", resp.Editable)
	_ = d.Set("description", resp.Description)
	_ = d.Set("type", resp.Type)
	_ = d.Set("custom_urls_count", resp.CustomUrlsCount)
	_ = d.Set("urls_retaining_parent_category_count", resp.UrlsRetainingParentCategoryCount)

	if err := d.Set("scopes", flattenScopes(resp)); err != nil {
		return err
	}

	if err := d.Set("url_keyword_counts", flattenUrlKeywordCounts(resp.URLKeywordCounts)); err != nil {
		return err
	}

	return nil
}

func resourceURLCategoriesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating custom url category ID: %v\n", id)
	req := expandURLCategory(d)

	if _, err := zClient.urlcategories.UpdateURLCategories(id, &req); err != nil {
		return err
	}

	return resourceURLCategoriesRead(d, m)
}

func resourceURLCategoriesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	// Need to pass the ID (int) of the resource for deletion
	log.Printf("[INFO] Deleting custom url category ID: %v\n", (d.Id()))

	if _, err := zClient.urlcategories.DeleteURLCategories(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] custom url category deleted")
	return nil
}

func expandURLCategory(d *schema.ResourceData) urlcategories.URLCategory {
	urlCategory := urlcategories.URLCategory{
		ConfiguredName:                   d.Get("login_name").(string),
		Urls:                             d.Get("urls").([]string),
		DBCategorizedUrls:                d.Get("db_categorized_urls").([]string),
		CustomCategory:                   d.Get("custom_category").(bool),
		Editable:                         d.Get("editable").(bool),
		Description:                      d.Get("description").(string),
		Type:                             d.Get("type").(string),
		CustomUrlsCount:                  d.Get("custom_urls_count").(int),
		UrlsRetainingParentCategoryCount: d.Get("urls_retaining_parent_category_count").(int),
		Scopes:                           expandURLCategoryScopes(d.Get("scopes").([]interface{})),
		URLKeywordCounts:                 expandURLKeywordCounts(d),
	}
	return urlCategory
}

func expandURLCategoryScopes(scopesUrlCategory []interface{}) []urlcategories.Scopes {
	scopes := make([]urlcategories.Scopes, len(scopesUrlCategory))

	for i, scope := range scopesUrlCategory {
		scopeItem := scope.(map[string]interface{})
		scopes[i] = urlcategories.Scopes{
			Type:                     scopeItem["type"].(string),
			ScopeGroupMemberEntities: expandCustomURLScopeGroupMemberEntities(scopeItem["scope_group_member_entities"].([]interface{})),
			ScopeEntities:            expandCustomURLScopeEntities(scopeItem["scope_entities"].([]interface{})),
		}
	}

	return scopes
}

func expandCustomURLScopeGroupMemberEntities(scopeGroupMember []interface{}) []urlcategories.ScopeGroupMemberEntities {
	scopeGroups := make([]urlcategories.ScopeGroupMemberEntities, len(scopeGroupMember))

	for i, scope := range scopeGroupMember {
		scopeGroup := scope.(map[string]interface{})
		scopeGroups[i] = urlcategories.ScopeGroupMemberEntities{
			ID:         scopeGroup["id"].(int),
			Extensions: scopeGroup["extensions"].(map[string]interface{}),
		}
	}

	return scopeGroups
}

func expandCustomURLScopeEntities(scopeEntity []interface{}) []urlcategories.ScopeEntities {
	scopeEntities := make([]urlcategories.ScopeEntities, len(scopeEntity))

	for i, scope := range scopeEntity {
		scopeEntity := scope.(map[string]interface{})
		scopeEntities[i] = urlcategories.ScopeEntities{
			ID:         scopeEntity["id"].(int),
			Extensions: scopeEntity["extensions"].(map[string]interface{}),
		}
	}

	return scopeEntities
}

func expandURLKeywordCounts(d *schema.ResourceData) urlcategories.URLKeywordCounts {
	keyword := urlcategories.URLKeywordCounts{
		TotalURLCount:            d.Get("total_url_count").(int),
		RetainParentURLCount:     d.Get("retain_parent_url_count").(int),
		TotalKeywordCount:        d.Get("total_keyword_count").(int),
		RetainParentKeywordCount: d.Get("retain_parent_keyword_count").(int),
	}
	return keyword
}
