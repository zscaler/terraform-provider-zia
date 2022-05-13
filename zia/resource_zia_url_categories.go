package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/terraform-provider-zia/gozscaler/client"
	"github.com/zscaler/terraform-provider-zia/gozscaler/urlcategories"
)

func resourceURLCategories() *schema.Resource {
	return &schema.Resource{
		Create: resourceURLCategoriesCreate,
		Read:   resourceURLCategoriesRead,
		Update: resourceURLCategoriesUpdate,
		Delete: resourceURLCategoriesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("url_category_id", id)
				} else {
					resp, err := zClient.urlcategories.GetCustomURLCategories(id)
					if err == nil {
						d.SetId(resp.ID)
						d.Set("url_category_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

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
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"keywords": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"keywords_retaining_parent_category": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"db_categorized_urls": {
				Type:     schema.TypeSet,
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
						"scope_group_member_entities": listIDsSchemaType("list of scope group member IDs"),
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
						"scope_entities": listIDsSchemaType("list of scope IDs"),
					},
				},
			},
			"editable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_url_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"retain_parent_url_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"total_keyword_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"retain_parent_keyword_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"custom_urls_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"urls_retaining_parent_category_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
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

	if err := d.Set("scopes", flattenScopesLite(resp)); err != nil {
		return err
	}

	if err := d.Set("url_keyword_counts", flattenUrlKeywordCounts(resp.URLKeywordCounts)); err != nil {
		return err
	}

	return nil
}

func flattenScopesLite(scopes *urlcategories.URLCategory) []interface{} {
	scope := make([]interface{}, len(scopes.Scopes))
	for i, val := range scopes.Scopes {
		scope[i] = map[string]interface{}{
			"type":                        val.Type,
			"scope_group_member_entities": flattenIDs(val.ScopeGroupMemberEntities),
			"scope_entities":              flattenIDs(val.ScopeEntities),
		}
	}

	return scope
}

func resourceURLCategoriesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getStringFromResourceData(d, "url_category_id")
	if !ok {
		log.Printf("[ERROR] custom url category ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating custom url category ID: %v\n", id)
	req := expandURLCategory(d)

	if _, _, err := zClient.urlcategories.UpdateURLCategories(id, &req); err != nil {
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
		Keywords:                         SetToStringList(d, "keywords"),
		KeywordsRetainingParentCategory:  SetToStringList(d, "keywords_retaining_parent_category"),
		Urls:                             SetToStringList(d, "urls"),
		DBCategorizedUrls:                SetToStringList(d, "db_categorized_urls"),
		CustomCategory:                   d.Get("custom_category").(bool),
		SuperCategory:                    d.Get("super_category").(string),
		Editable:                         d.Get("editable").(bool),
		Description:                      d.Get("description").(string),
		Type:                             d.Get("type").(string),
		CustomUrlsCount:                  d.Get("custom_urls_count").(int),
		UrlsRetainingParentCategoryCount: d.Get("urls_retaining_parent_category_count").(int),
		Scopes:                           expandURLCategoryScopes(d),
		URLKeywordCounts:                 expandURLKeywordCounts(d),
	}
	return result
}

func expandURLKeywordCounts(d *schema.ResourceData) *urlcategories.URLKeywordCounts {
	keywordCounts := urlcategories.URLKeywordCounts{}
	if keywordCountsInterface, ok := d.GetOk("url_keyword_counts"); ok {
		keywordCountsList := keywordCountsInterface.([]interface{})
		for _, keywordCountsMap := range keywordCountsList {
			keywordCountsItem := keywordCountsMap.(map[string]interface{})
			keywordCounts.TotalURLCount, _ = keywordCountsItem["total_url_count"].(int)
			keywordCounts.RetainParentURLCount, _ = keywordCountsItem["retain_parent_url_count"].(int)
			keywordCounts.TotalKeywordCount, _ = keywordCountsItem["total_keyword_count"].(int)
			keywordCounts.RetainParentKeywordCount, _ = keywordCountsItem["retain_parent_keyword_count"].(int)
			break
		}
	}
	return &keywordCounts
}
func expandURLCategoryScopes(d *schema.ResourceData) []urlcategories.Scopes {
	var scopes []urlcategories.Scopes
	if scopeInterface, ok := d.GetOk("scopes"); ok {
		scopesSet, ok := scopeInterface.(*schema.Set)
		if !ok {
			return scopes
		}
		scopes = make([]urlcategories.Scopes, len(scopesSet.List()))
		for i, val := range scopesSet.List() {
			scopeItem := val.(map[string]interface{})
			scopes[i] = urlcategories.Scopes{
				ScopeGroupMemberEntities: expandIDNameExtensionsMap(scopeItem, "scope_group_member_entities"),
				Type:                     scopeItem["type"].(string),
				ScopeEntities:            expandIDNameExtensionsMap(scopeItem, "scope_entities"),
			}
		}
	}
	return scopes
}
