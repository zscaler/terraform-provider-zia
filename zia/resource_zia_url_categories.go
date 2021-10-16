package zia

import (
	"fmt"
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
			"url_category_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"keywords": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"keywords_retaining_parent_category": {
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
				Default:  false,
			},
			"super_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scopes": {
				Type:     schema.TypeSet,
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
							Type:     schema.TypeSet,
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
				Default:  false,
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
	d.SetId(resp.ID)
	_ = d.Set("url_category_id", resp.ID)
	return resourceURLCategoriesRead(d, m)
}

func resourceURLCategoriesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getStringFromResourceData(d, "url_category_id")
	if !ok {
		return fmt.Errorf("no url category rule id is set")
	}
	resp, err := zClient.urlcategories.Get(id)

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing zia url category %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting url category :\n%+v\n", resp)

	d.SetId(fmt.Sprintf(resp.ID))
	_ = d.Set("url_category_id", resp.ID)
	_ = d.Set("configured_name", resp.ConfiguredName)
	_ = d.Set("keywords", resp.Keywords)
	_ = d.Set("keywords_retaining_parent_category", resp.KeywordsRetainingParentCategory)
	_ = d.Set("urls", resp.Urls)
	_ = d.Set("db_categorized_urls", resp.DBCategorizedUrls)
	_ = d.Set("custom_category", resp.CustomCategory)
	_ = d.Set("super_category", resp.SuperCategory)
	_ = d.Set("editable", resp.Editable)
	_ = d.Set("description", resp.Description)
	_ = d.Set("type", resp.Type)
	_ = d.Set("custom_urls_count", resp.CustomUrlsCount)
	_ = d.Set("urls_retaining_parent_category_count", resp.UrlsRetainingParentCategoryCount)

	if err := d.Set("scopes", flattenScopes(resp)); err != nil {
		return err
	}

	// if err := d.Set("url_keyword_counts", flattenUrlKeywordCounts(resp.URLKeywordCounts)); err != nil {
	// 	return err
	// }

	return nil
}

func resourceURLCategoriesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getStringFromResourceData(d, "url_category_id")
	if !ok {
		log.Printf("[ERROR] vpn credentials ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating custom url category ID: %v\n", id)
	req := expandURLCategory(d)

	if _, err := zClient.urlcategories.UpdateURLCategories(id, &req); err != nil {
		return err
	}

	return resourceURLCategoriesRead(d, m)
}

func resourceURLCategoriesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "url_category_id")
	if !ok {
		log.Printf("[ERROR] url category id ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting custom url category ID: %v\n", (d.Id()))

	if _, err := zClient.urlcategories.DeleteURLCategories(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] custom url category deleted")
	return nil
}

func expandURLCategory(d *schema.ResourceData) urlcategories.URLCategory {
	id, _ := getStringFromResourceData(d, "url_category_id")
	result := urlcategories.URLCategory{
		ID:                               id,
		ConfiguredName:                   d.Get("configured_name").(string),
		Keywords:                         ListToStringSlice(d.Get("keywords").([]interface{})),
		KeywordsRetainingParentCategory:  ListToStringSlice(d.Get("keywords_retaining_parent_category").([]interface{})),
		Urls:                             ListToStringSlice(d.Get("urls").([]interface{})),
		DBCategorizedUrls:                ListToStringSlice(d.Get("db_categorized_urls").([]interface{})),
		CustomCategory:                   d.Get("custom_category").(bool),
		SuperCategory:                    d.Get("super_category").(string),
		Editable:                         d.Get("editable").(bool),
		Description:                      d.Get("description").(string),
		Type:                             d.Get("type").(string),
		CustomUrlsCount:                  d.Get("custom_urls_count").(int),
		UrlsRetainingParentCategoryCount: d.Get("urls_retaining_parent_category_count").(int),
		// Scopes:                           expandURLCategoryScopes(d),
		//URLKeywordCounts: expandURLKeywordCounts(d),
	}
	urlCategoryScopes := expandURLCategoryScopes(d)
	if urlCategoryScopes != nil {
		result.Scopes = urlCategoryScopes
	}
	return result
}

func expandURLCategoryScopes(d *schema.ResourceData) []urlcategories.Scopes {
	var scopes []urlcategories.Scopes
	if scopeInterface, ok := d.GetOk("scopes"); ok {
		scope := scopeInterface.([]interface{})
		scopes = make([]urlcategories.Scopes, len(scope))
		for i, val := range scope {
			scopeItem := val.(map[string]interface{})
			scopes[i] = urlcategories.Scopes{
				ScopeGroupMemberEntities: expandCustomURLScopeGroupMemberEntities(d),
				Type:                     scopeItem["type"].(string),
				ScopeEntities:            expandCustomURLScopeEntities(d),
			}
		}
	}

	return scopes
}

func expandCustomURLScopeGroupMemberEntities(d *schema.ResourceData) []urlcategories.ScopeGroupMemberEntities {
	var scopeGroupMemberEntities []urlcategories.ScopeGroupMemberEntities
	if scopeGroupInterface, ok := d.GetOk("scope_group_member_entities"); ok {
		scopeGroup := scopeGroupInterface.([]interface{})
		scopeGroupMemberEntities = make([]urlcategories.ScopeGroupMemberEntities, len(scopeGroup))
		for i, val := range scopeGroup {
			scopeGroupItem := val.(map[string]interface{})
			scopeGroupMemberEntities[i] = urlcategories.ScopeGroupMemberEntities{
				ID:         scopeGroupItem["id"].(int),
				Extensions: scopeGroupItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return scopeGroupMemberEntities
}

func expandCustomURLScopeEntities(d *schema.ResourceData) []urlcategories.ScopeEntities {
	var scopeEntities []urlcategories.ScopeEntities
	if scopeEntitiesInterface, ok := d.GetOk("scope_entities"); ok {
		scopeEntity := scopeEntitiesInterface.([]interface{})
		scopeEntities = make([]urlcategories.ScopeEntities, len(scopeEntity))
		for i, val := range scopeEntity {
			scopeEntityItem := val.(map[string]interface{})
			scopeEntities[i] = urlcategories.ScopeEntities{
				ID:         scopeEntityItem["id"].(int),
				Extensions: scopeEntityItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return scopeEntities
}
