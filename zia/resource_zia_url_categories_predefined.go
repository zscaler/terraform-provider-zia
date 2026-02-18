package zia

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlcategories"
)

func resourceURLCategoriesPredefined() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceURLCategoriesPredefinedCreate,
		ReadContext:   resourceURLCategoriesPredefinedRead,
		UpdateContext: resourceURLCategoriesPredefinedUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				identifier := d.Id()

				allCategories, err := urlcategories.GetAll(ctx, service, false, false, "ALL")
				if err != nil {
					return nil, fmt.Errorf("failed retrieving URL categories: %w", err)
				}

				var matched *urlcategories.URLCategory
				for i := range allCategories {
					cat := allCategories[i]
					if cat.CustomCategory {
						continue
					}
					if strings.EqualFold(cat.ID, identifier) || strings.EqualFold(cat.ConfiguredName, identifier) {
						matched = &cat
						break
					}
				}

				if matched == nil {
					return nil, fmt.Errorf("no predefined URL category found with ID or name: %q", identifier)
				}

				d.SetId(matched.ID)
				_ = d.Set("name", matched.ID)
				_ = d.Set("category_id", matched.ID)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The predefined URL category ID or display name (e.g., `FINANCE` or `Finance`). The provider resolves this to the canonical category ID.",
			},
			"category_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The canonical predefined URL category identifier resolved by the provider.",
			},
			"configured_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the predefined URL category. Read-only.",
			},
			"super_category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The super category of the predefined URL category. Read-only.",
			},
			"url_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL type (e.g., `EXACT`). Read-only for predefined categories.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the URL category. Read-only.",
			},
			"val": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The numeric identifier for the URL category. Read-only.",
			},
			"editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the category is editable. Read-only.",
			},
			"urls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom URLs to add to the predefined URL category.",
			},
			"urls_retaining_parent_category": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URLs that are also retained under the original parent URL category.",
			},
			"keywords": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom keywords associated to the URL category.",
			},
			"keywords_retaining_parent_category": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Retained custom keywords from the parent URL category.",
			},
			"ip_ranges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom IP address ranges associated to the URL category.",
			},
			"ip_ranges_retaining_parent_category": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Retaining parent custom IP address ranges associated to the URL category.",
			},
			"db_categorized_urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URLs categorized by the Zscaler database. Read-only.",
			},
			"custom_urls_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of custom URLs associated to the URL category. Read-only.",
			},
			"urls_retaining_parent_category_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of URLs retaining the parent category. Read-only.",
			},
			"custom_ip_ranges_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of custom IP address ranges associated to the URL category. Read-only.",
			},
			"ip_ranges_retaining_parent_category_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of IP ranges retaining the parent category. Read-only.",
			},
			"url_keyword_counts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "URL and keyword counts for the URL category. Read-only.",
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
		},
	}
}

// resolvePredefinedCategory looks up a predefined category by ID or ConfiguredName
// from the full list of categories. Returns the matched category or nil.
func resolvePredefinedCategory(allCategories []urlcategories.URLCategory, identifier string) *urlcategories.URLCategory {
	for i := range allCategories {
		if allCategories[i].CustomCategory {
			continue
		}
		if strings.EqualFold(allCategories[i].ID, identifier) || strings.EqualFold(allCategories[i].ConfiguredName, identifier) {
			return &allCategories[i]
		}
	}
	return nil
}

func resourceURLCategoriesPredefinedCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	name := d.Get("name").(string)

	allCategories, err := urlcategories.GetAll(ctx, service, false, false, "ALL")
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed retrieving URL categories: %w", err))
	}

	existing := resolvePredefinedCategory(allCategories, name)
	if existing == nil {
		return diag.FromErr(fmt.Errorf("predefined URL category %q not found", name))
	}

	resolvedID := existing.ID
	log.Printf("[INFO] Managing predefined URL category %s (resolved from %q)", resolvedID, name)

	desired := expandURLCategoryPredefined(d, existing)

	if diags := applyPredefinedCategoryDiff(ctx, service, resolvedID, existing, &desired); diags != nil {
		return diags
	}

	d.SetId(resolvedID)
	_ = d.Set("category_id", resolvedID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceURLCategoriesPredefinedRead(ctx, d, meta)
}

func resourceURLCategoriesPredefinedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	categoryID := d.Get("category_id").(string)
	if categoryID == "" {
		categoryID = d.Id()
	}

	allCategories, err := urlcategories.GetAll(ctx, service, false, false, "ALL")
	if err != nil {
		return diag.FromErr(err)
	}

	resp := resolvePredefinedCategory(allCategories, categoryID)

	if resp == nil {
		log.Printf("[DEBUG] Predefined category %s not found in GetAll() response, falling back to individual Get()", categoryID)
		individualResp, err := urlcategories.Get(ctx, service, categoryID)
		if err != nil {
			log.Printf("[WARN] Predefined URL category %s not found, removing from state", categoryID)
			d.SetId("")
			return nil
		}
		if individualResp.CustomCategory {
			log.Printf("[WARN] Category %s is a custom category, not predefined — removing from state", categoryID)
			d.SetId("")
			return nil
		}
		resp = individualResp
	}

	log.Printf("[INFO] Reading predefined URL category: %s", resp.ID)

	d.SetId(resp.ID)
	_ = d.Set("category_id", resp.ID)
	_ = d.Set("configured_name", resp.ConfiguredName)
	_ = d.Set("super_category", resp.SuperCategory)
	_ = d.Set("url_type", resp.UrlType)
	_ = d.Set("type", resp.Type)
	_ = d.Set("val", resp.Val)
	_ = d.Set("editable", resp.Editable)
	_ = d.Set("urls", resp.Urls)
	_ = d.Set("db_categorized_urls", resp.DBCategorizedUrls)
	_ = d.Set("keywords", resp.Keywords)
	_ = d.Set("keywords_retaining_parent_category", resp.KeywordsRetainingParentCategory)
	_ = d.Set("ip_ranges", resp.IPRanges)
	_ = d.Set("ip_ranges_retaining_parent_category", resp.IPRangesRetainingParentCategory)
	_ = d.Set("custom_urls_count", resp.CustomUrlsCount)
	_ = d.Set("urls_retaining_parent_category_count", resp.UrlsRetainingParentCategoryCount)
	_ = d.Set("custom_ip_ranges_count", resp.CustomIpRangesCount)
	_ = d.Set("ip_ranges_retaining_parent_category_count", resp.IPRangesRetainingParentCategoryCount)

	if err := d.Set("url_keyword_counts", flattenUrlKeywordCounts(resp.URLKeywordCounts)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceURLCategoriesPredefinedUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	categoryID := d.Get("category_id").(string)

	currentCategory, err := urlcategories.Get(ctx, service, categoryID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read predefined URL category %s: %w", categoryID, err))
	}

	log.Printf("[INFO] Updating predefined URL category %s", categoryID)

	desired := expandURLCategoryPredefined(d, currentCategory)

	if diags := applyPredefinedCategoryDiff(ctx, service, categoryID, currentCategory, &desired); diags != nil {
		return diags
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceURLCategoriesPredefinedRead(ctx, d, meta)
}

// applyPredefinedCategoryDiff computes the diff between the current API state and
// the desired HCL state, then applies changes using the appropriate API semantics:
//   - urls, ip_ranges, ip_ranges_retaining_parent_category: incremental via
//     ADD_TO_LIST / REMOVE_FROM_LIST (plain PUT is additive for these fields)
//   - keywords, keywords_retaining_parent_category: replacement via plain PUT
//     (the API does not support ADD_TO_LIST / REMOVE_FROM_LIST for keyword fields)
func applyPredefinedCategoryDiff(ctx context.Context, service *zscaler.Service, categoryID string, current *urlcategories.URLCategory, desired *urlcategories.URLCategory) diag.Diagnostics {
	urlsToAdd, urlsToRemove := computeSliceDiff(current.Urls, desired.Urls)
	ipToAdd, ipToRemove := computeSliceDiff(current.IPRanges, desired.IPRanges)
	ipRetToAdd, ipRetToRemove := computeSliceDiff(current.IPRangesRetainingParentCategory, desired.IPRangesRetainingParentCategory)

	incrementalRemoves := len(urlsToRemove) + len(ipToRemove) + len(ipRetToRemove)
	incrementalAdds := len(urlsToAdd) + len(ipToAdd) + len(ipRetToAdd)
	keywordsChanged := !stringSlicesEqual(current.Keywords, desired.Keywords)
	kwRetChanged := !stringSlicesEqual(current.KeywordsRetainingParentCategory, desired.KeywordsRetainingParentCategory)

	log.Printf("[DEBUG] Predefined category diff — Incremental removes: %d (urls=%d, ip=%d, ipRet=%d), Incremental adds: %d (urls=%d, ip=%d, ipRet=%d), Keywords changed: %v, Keywords retaining changed: %v",
		incrementalRemoves, len(urlsToRemove), len(ipToRemove), len(ipRetToRemove),
		incrementalAdds, len(urlsToAdd), len(ipToAdd), len(ipRetToAdd),
		keywordsChanged, kwRetChanged)

	if incrementalRemoves > 0 {
		log.Printf("[INFO] Using REMOVE_FROM_LIST to remove %d items from url/ip fields", incrementalRemoves)
		removeCategory := buildImmutableCategory(current)
		removeCategory.Urls = urlsToRemove
		removeCategory.IPRanges = ipToRemove
		removeCategory.IPRangesRetainingParentCategory = ipRetToRemove
		if _, _, err := urlcategories.UpdateURLCategories(ctx, service, categoryID, &removeCategory, "REMOVE_FROM_LIST"); err != nil {
			return diag.FromErr(fmt.Errorf("failed to remove items from predefined category: %w", err))
		}
	}

	if incrementalAdds > 0 {
		log.Printf("[INFO] Using ADD_TO_LIST to add %d items to url/ip fields", incrementalAdds)
		addCategory := buildImmutableCategory(current)
		addCategory.Urls = urlsToAdd
		addCategory.IPRanges = ipToAdd
		addCategory.IPRangesRetainingParentCategory = ipRetToAdd
		if _, _, err := urlcategories.UpdateURLCategories(ctx, service, categoryID, &addCategory, "ADD_TO_LIST"); err != nil {
			return diag.FromErr(fmt.Errorf("failed to add items to predefined category: %w", err))
		}
	}

	if keywordsChanged || kwRetChanged {
		log.Printf("[INFO] Using plain PUT to update keyword fields (replacement semantics)")
		kwCategory := buildImmutableCategory(current)
		kwCategory.Keywords = desired.Keywords
		kwCategory.KeywordsRetainingParentCategory = desired.KeywordsRetainingParentCategory
		if _, _, err := urlcategories.UpdateURLCategories(ctx, service, categoryID, &kwCategory, ""); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update keywords for predefined category: %w", err))
		}
	}

	return nil
}

func buildImmutableCategory(existing *urlcategories.URLCategory) urlcategories.URLCategory {
	return urlcategories.URLCategory{
		ID:             existing.ID,
		ConfiguredName: existing.ConfiguredName,
		SuperCategory:  existing.SuperCategory,
		UrlType:        existing.UrlType,
		Type:           existing.Type,
		Val:            existing.Val,
		CustomCategory: existing.CustomCategory,
		Editable:       existing.Editable,
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := stringSliceToMap(a)
	bMap := stringSliceToMap(b)
	for k := range aMap {
		if !bMap[k] {
			return false
		}
	}
	return true
}

func computeSliceDiff(current, desired []string) (toAdd, toRemove []string) {
	currentMap := stringSliceToMap(current)
	desiredMap := stringSliceToMap(desired)

	for item := range desiredMap {
		if !currentMap[item] {
			toAdd = append(toAdd, item)
		}
	}
	for item := range currentMap {
		if !desiredMap[item] {
			toRemove = append(toRemove, item)
		}
	}
	return
}

func expandURLCategoryPredefined(d *schema.ResourceData, existing *urlcategories.URLCategory) urlcategories.URLCategory {
	return urlcategories.URLCategory{
		ID:                              existing.ID,
		ConfiguredName:                  existing.ConfiguredName,
		SuperCategory:                   existing.SuperCategory,
		UrlType:                         existing.UrlType,
		Type:                            existing.Type,
		Val:                             existing.Val,
		CustomCategory:                  existing.CustomCategory,
		Editable:                        existing.Editable,
		Urls:                            SetToStringList(d, "urls"),
		Keywords:                        SetToStringList(d, "keywords"),
		KeywordsRetainingParentCategory: SetToStringList(d, "keywords_retaining_parent_category"),
		IPRanges:                        SetToStringList(d, "ip_ranges"),
		IPRangesRetainingParentCategory: SetToStringList(d, "ip_ranges_retaining_parent_category"),
	}
}
