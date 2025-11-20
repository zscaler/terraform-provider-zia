package zia

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlcategories"
)

// Allows only one API request at a time
// var urlCategoriesSemaphore = make(chan struct{}, 1)

func resourceURLCategories() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceURLCategoriesCreate,
		ReadContext:   resourceURLCategoriesRead,
		UpdateContext: resourceURLCategoriesUpdate,
		DeleteContext: resourceURLCategoriesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				identifier := d.Id()

				// Use the new slice-returning helper
				categories, err := urlcategories.GetAllCustomURLCategories(ctx, service)
				if err != nil {
					return nil, fmt.Errorf("failed retrieving custom URL categories: %w", err)
				}

				var matched *urlcategories.URLCategory
				for i := range categories {
					cat := categories[i]
					if strings.EqualFold(cat.ID, identifier) || strings.EqualFold(cat.ConfiguredName, identifier) {
						matched = &cat
						break
					}
				}

				if matched == nil {
					return nil, fmt.Errorf("no custom URL category found with ID or configuredName: %q", identifier)
				}

				d.SetId(matched.ID)
				_ = d.Set("category_id", matched.ID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"category_id": {
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
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringLenBetween(0, 256),
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"urls": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 25000,
			},
			"keywords": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 2048,
			},
			"keywords_retaining_parent_category": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 2048,
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
			"ip_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 2000,
			},
			"ip_ranges_retaining_parent_category": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 2000,
			},
			"custom_ip_ranges_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ip_ranges_retaining_parent_category_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"super_category": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Super Category of the URL category.
				This field is required when creating custom URL categories..
				See the URL Categories API for the list of available super categories:
				https://help.zscaler.com/zia/url-categories#/urlCategories-get`,
			},
			"val": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique ID for the URL category.",
			},
		},
	}
}

func resourceURLCategoriesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Acquire semaphore before making an API request
	apiSemaphore <- struct{}{}
	defer func() { <-apiSemaphore }() // Release semaphore after the request is done

	zClient := meta.(*Client)
	service := zClient.Service

	req := expandURLCategory(d)
	log.Printf("[INFO] Creating zia url category\n%+v\n", req)

	// Use the existing CreateURLCategories function
	resp, err := urlcategories.CreateURLCategories(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created zia url category request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)
	_ = d.Set("category_id", resp.ID)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceURLCategoriesRead(ctx, d, meta)
}

func resourceURLCategoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getStringFromResourceData(d, "category_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no url category id is set"))
	}

	// Use GetAll() instead of Get() to reduce API calls during terraform refresh
	// customOnly=true to only retrieve custom categories (which are the ones managed by Terraform)
	// includeOnlyUrlKeywordCounts=false to get full category details
	allCategories, err := urlcategories.GetAll(ctx, service, true, false)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the specific category by ID
	var resp *urlcategories.URLCategory
	for i := range allCategories {
		if allCategories[i].ID == id {
			resp = &allCategories[i]
			break
		}
	}

	// If not found in GetAll(), fall back to individual Get() call
	// This handles newly created categories that may not be in the cached list yet
	if resp == nil {
		log.Printf("[DEBUG] Category %s not found in GetAll() response, falling back to individual Get() call", id)
		individualResp, err := urlcategories.Get(ctx, service, id)
		if err != nil {
			if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
				log.Printf("[WARN] Removing url category %s from state because it no longer exists in ZIA", d.Id())
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
		resp = individualResp
	}

	log.Printf("[INFO] Getting url category :\n%+v\n", resp)

	d.SetId(resp.ID)
	_ = d.Set("category_id", resp.ID)
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
	_ = d.Set("ip_ranges", resp.IPRanges)
	_ = d.Set("ip_ranges_retaining_parent_category", resp.IPRangesRetainingParentCategory)
	_ = d.Set("custom_ip_ranges_count", resp.CustomIpRangesCount)
	_ = d.Set("ip_ranges_retaining_parent_category_count", resp.IPRangesRetainingParentCategoryCount)
	_ = d.Set("val", resp.Val)

	if err := d.Set("scopes", flattenScopesLite(resp)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url_keyword_counts", flattenUrlKeywordCounts(resp.URLKeywordCounts)); err != nil {
		return diag.FromErr(err)
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

func resourceURLCategoriesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Acquire semaphore before making an API request
	apiSemaphore <- struct{}{}
	defer func() { <-apiSemaphore }() // Release semaphore after the request is done

	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getStringFromResourceData(d, "category_id")
	if !ok {
		log.Printf("[ERROR] custom url category ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("custom url category ID not set: %v", id))
	}

	log.Printf("[INFO] Updating custom url category ID: %v\n", id)
	req := expandURLCategory(d)

	if _, err := urlcategories.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, _, err := urlcategories.UpdateURLCategories(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceURLCategoriesRead(ctx, d, meta)
}

func resourceURLCategoriesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "category_id")
	if !ok {
		log.Printf("[ERROR] url category id ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting custom url category ID: %v\n", (d.Id()))

	if _, err := urlcategories.DeleteURLCategories(ctx, service, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] custom url category deleted")

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandURLCategory(d *schema.ResourceData) urlcategories.URLCategory {
	id, _ := getStringFromResourceData(d, "category_id")
	result := urlcategories.URLCategory{
		ID:                                   id,
		Val:                                  d.Get("val").(int),
		ConfiguredName:                       d.Get("configured_name").(string),
		Keywords:                             SetToStringList(d, "keywords"),
		KeywordsRetainingParentCategory:      SetToStringList(d, "keywords_retaining_parent_category"),
		Urls:                                 SetToStringList(d, "urls"),
		DBCategorizedUrls:                    SetToStringList(d, "db_categorized_urls"),
		IPRanges:                             SetToStringList(d, "ip_ranges"),
		IPRangesRetainingParentCategory:      SetToStringList(d, "ip_ranges_retaining_parent_category"),
		CustomIpRangesCount:                  d.Get("custom_ip_ranges_count").(int),
		IPRangesRetainingParentCategoryCount: d.Get("ip_ranges_retaining_parent_category_count").(int),
		CustomCategory:                       d.Get("custom_category").(bool),
		SuperCategory:                        d.Get("super_category").(string),
		Editable:                             d.Get("editable").(bool),
		Description:                          d.Get("description").(string),
		Type:                                 d.Get("type").(string),
		CustomUrlsCount:                      d.Get("custom_urls_count").(int),
		UrlsRetainingParentCategoryCount:     d.Get("urls_retaining_parent_category_count").(int),
		Scopes:                               expandURLCategoryScopes(d),
		URLKeywordCounts:                     expandURLKeywordCounts(d),
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
