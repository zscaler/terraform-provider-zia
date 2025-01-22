package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlpdictionaries"
)

func resourceDLPDictionaries() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDLPDictionariesCreate,
		ReadContext:   resourceDLPDictionariesRead,
		UpdateContext: resourceDLPDictionariesUpdate,
		DeleteContext: resourceDLPDictionariesDelete,
		CustomizeDiff: validateDLPHierarchicalIdentifiersDiff,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("dictionary_id", idInt)
				} else {
					resp, err := dlpdictionaries.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("dictionary_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"dictionary_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The DLP dictionary's name",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The desciption of the DLP dictionary",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"custom": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The DLP dictionary proximity length.",
			},
			"confidence_level_for_predefined_dict": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP confidence threshold for predefined dictionaries",
				// Default:     "CONFIDENCE_LEVEL_MEDIUM",
				ValidateFunc: validation.StringInSlice([]string{
					"CONFIDENCE_LEVEL_LOW",
					"CONFIDENCE_LEVEL_MEDIUM",
					"CONFIDENCE_LEVEL_HIGH",
				}, false),
			},
			"confidence_threshold": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP confidence threshold",
				ValidateFunc: validation.StringInSlice([]string{
					"CONFIDENCE_LEVEL_LOW",
					"CONFIDENCE_LEVEL_MEDIUM",
					"CONFIDENCE_LEVEL_HIGH",
				}, false),
			},
			"phrases": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 256,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PHRASE_COUNT_TYPE_UNIQUE",
								"PHRASE_COUNT_TYPE_ALL",
							}, false),
						},
						"phrase": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(0, 128),
						},
					},
				},
			},
			"custom_phrase_match_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY",
					"MATCH_ANY_CUSTOM_PHRASE_PATTERN_DICTIONARY",
				}, false),
			},
			"patterns": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				MaxItems:    8,
				Description: "List containing the patterns used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The action applied to a DLP dictionary using patterns",
							ValidateFunc: validation.StringInSlice([]string{
								"PATTERN_COUNT_TYPE_ALL",
								"PATTERN_COUNT_TYPE_UNIQUE",
							}, false),
						},
						"pattern": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "DLP dictionary pattern",
							ValidateFunc: validation.StringLenBetween(0, 256),
						},
					},
				},
			},
			"dictionary_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP dictionary type.",
				ValidateFunc: validation.StringInSlice([]string{
					"PATTERNS_AND_PHRASES",
					"EXACT_DATA_MATCH",
					"INDEXED_DATA_MATCH",
				}, false),
			},
			"hierarchical_identifiers": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of hierarchical identifiers for the DLP dictionary.",
			},
			"exact_data_match_details": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Exact Data Match (EDM) related information for custom DLP dictionaries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dictionary_edm_mapping_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The unique identifier for the EDM mapping",
						},
						"schema_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The unique identifier for the EDM template (or schema).",
						},
						"primary_field": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The EDM template's primary field.",
						},
						"secondary_fields": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"secondary_field_match_on": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The EDM secondary field to match on.",
							ValidateFunc: validation.StringInSlice([]string{
								"MATCHON_NONE", "MATCHON_ANY_1", "MATCHON_ANY_2",
								"MATCHON_ANY_3", "MATCHON_ANY_4", "MATCHON_ANY_5",
								"MATCHON_ANY_6", "MATCHON_ANY_7", "MATCHON_ANY_8",
								"MATCHON_ANY_9", "MATCHON_ANY_10", "MATCHON_ANY_11",
								"MATCHON_ANY_12", "MATCHON_ANY_13", "MATCHON_ANY_14",
								"MATCHON_ANY_15", "MATCHON_ALL",
							}, false),
						},
					},
				},
			},
			"idm_profile_match_accuracy": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "List of Indexed Document Match (IDM) profiles and their corresponding match accuracy for custom DLP dictionaries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"adp_idm_profile": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							// MaxItems:    1,
							Description: "The action applied to a DLP dictionary using patterns",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
										Optional: true,
									},
									"extensions": {
										Type:     schema.TypeMap,
										Computed: true,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"match_accuracy": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The IDM template match accuracy.",
							ValidateFunc: validation.StringInSlice([]string{
								"LOW", "MEDIUM", "HEAVY",
							}, false),
						},
					},
				},
			},
			"proximity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The DLP dictionary proximity length.",
			},
			"ignore_exact_match_idm_dict": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to exclude documents that are a 100% match to already-indexed documents from triggering an Indexed Document Match (IDM) Dictionary.",
			},
			"include_bin_numbers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "A true value denotes that the specified Bank Identification Number (BIN) values are included in the Credit Cards dictionary. A false value denotes that the specified BIN values are excluded from the Credit Cards dictionary.Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.",
			},
			"bin_numbers": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "The list of Bank Identification Number (BIN) values that are included or excluded from the Credit Cards dictionary. BIN values can be specified only for Diners Club, Mastercard, RuPay, and Visa cards. Up to 512 BIN values can be configured in a dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.",
			},
			"dict_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the predefined dictionary (original source dictionary) that is used for cloning. This field is applicable only to cloned dictionaries. Only a limited set of identification-based predefined dictionaries (e.g., Credit Cards, Social Security Numbers, National Identification Numbers, etc.) can be cloned. Up to 4 clones can be created from a predefined dictionary.",
			},
		},
	}
}

func resourceDLPDictionariesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandDLPDictionaries(d, true)
	log.Printf("[INFO] Creating zia dlp dictionaries\n%+v\n", req)
	if req.DictionaryType != "PATTERNS_AND_PHRASES" && req.CustomPhraseMatchType != "" {
		log.Printf("[ERROR] custom_phrase_match_type should not be set when dictionary_type is not set to 'PATTERNS_AND_PHRASES'")
		return diag.FromErr(fmt.Errorf("[ERROR] custom_phrase_match_type should not be set when dictionary_type is not set to 'PATTERNS_AND_PHRASES'"))
	}
	resp, _, err := dlpdictionaries.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia dlp dictionaries request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("dictionary_id", resp.ID)
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceDLPDictionariesRead(ctx, d, meta)
}

func resourceDLPDictionariesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "dictionary_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no DLP dictionary id is set"))
	}
	resp, err := dlpdictionaries.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing dlp dictionary %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting dlp dictionary :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("dictionary_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("custom", resp.Custom)
	_ = d.Set("confidence_threshold", resp.ConfidenceThreshold)
	_ = d.Set("confidence_level_for_predefined_dict", resp.ConfidenceLevelForPredefinedDict)
	_ = d.Set("custom_phrase_match_type", resp.CustomPhraseMatchType)
	_ = d.Set("dictionary_type", resp.DictionaryType)
	_ = d.Set("hierarchical_identifiers", d.Get("hierarchical_identifiers")) // Keep the user input for hierarchical_identifiers
	_ = d.Set("ignore_exact_match_idm_dict", resp.IgnoreExactMatchIdmDict)
	_ = d.Set("include_bin_numbers", resp.IncludeBinNumbers)
	_ = d.Set("bin_numbers", resp.BinNumbers)
	_ = d.Set("dict_template_id", resp.DictTemplateId)
	_ = d.Set("proximity", resp.Proximity)
	if err := d.Set("phrases", flattenPhrases(resp)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("patterns", flattenPatterns(resp)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("exact_data_match_details", flattenEDMDetails(resp)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("idm_profile_match_accuracy", flattenIDMProfileMatchAccuracySimple(resp)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDLPDictionariesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "dictionary_id")
	if !ok {
		log.Printf("[ERROR] dlp dictionaryID not set: %v\n", id)
	}

	log.Printf("[INFO] Updating dlp dictionary ID: %v\n", id)
	req := expandDLPDictionaries(d, true)
	if _, err := dlpdictionaries.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if req.DictionaryType != "PATTERNS_AND_PHRASES" && req.CustomPhraseMatchType != "" {
		log.Printf("[ERROR] custom_phrase_match_type should not be set when dictionary_type is not set to 'PATTERNS_AND_PHRASES'")
		return diag.FromErr(fmt.Errorf("[ERROR] custom_phrase_match_type should not be set when dictionary_type is not set to 'PATTERNS_AND_PHRASES'"))
	}
	if _, _, err := dlpdictionaries.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceDLPDictionariesRead(ctx, d, meta)
}

func resourceDLPDictionariesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "dictionary_id")
	if !ok {
		log.Printf("[ERROR] dlp dictionary ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting dlp dictionary ID: %v\n", (d.Id()))

	if _, err := dlpdictionaries.DeleteDlpDictionary(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] dlp dictionary deleted")
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandDLPDictionaries(d *schema.ResourceData, isCreate bool) dlpdictionaries.DlpDictionary {
	id, _ := getIntFromResourceData(d, "dictionary_id")
	hierarchicalIdentifiers := expandHierarchicalIdentifiers(d.Get("hierarchical_identifiers").(*schema.Set).List())

	result := dlpdictionaries.DlpDictionary{
		ID:                               id,
		Name:                             d.Get("name").(string),
		Description:                      d.Get("description").(string),
		ConfidenceThreshold:              d.Get("confidence_threshold").(string),
		CustomPhraseMatchType:            d.Get("custom_phrase_match_type").(string),
		ConfidenceLevelForPredefinedDict: d.Get("confidence_level_for_predefined_dict").(string),
		DictionaryType:                   d.Get("dictionary_type").(string),
		Custom:                           d.Get("custom").(bool),
		IgnoreExactMatchIdmDict:          d.Get("ignore_exact_match_idm_dict").(bool),
		IncludeBinNumbers:                d.Get("include_bin_numbers").(bool),
		DictTemplateId:                   d.Get("dict_template_id").(int),
		Proximity:                        d.Get("proximity").(int),
		HierarchicalIdentifiers:          hierarchicalIdentifiers,
	}
	binNumbers := []int{}
	for _, i := range d.Get("bin_numbers").([]interface{}) {
		binNumbers = append(binNumbers, i.(int))
	}
	result.BinNumbers = binNumbers

	phrases := expandDLPDictionariesPhrases(d)
	if phrases != nil {
		result.Phrases = phrases
	}

	// Include phrases and patterns only if it's a create operation or if the dictionary is not cloned
	if isCreate || (result.DictTemplateId == 0 && len(result.HierarchicalIdentifiers) == 0) {
		patterns := expandDLPDictionariesPatterns(d)
		if patterns != nil {
			result.Patterns = patterns
		}
	}

	edmDetails := expandEDMDetails(d)
	if edmDetails != nil {
		result.EDMMatchDetails = edmDetails
	}

	idmProfileMarch := expandIDMProfileMatchAccuracy(d)
	if idmProfileMarch != nil {
		result.IDMProfileMatchAccuracy = idmProfileMarch
	}
	return result
}

func expandHierarchicalIdentifiers(identifiers []interface{}) []string {
	var expandedIdentifiers []string
	for _, identifier := range identifiers {
		if values, exists := predefinedIdentifiersMap[identifier.(string)]; exists {
			expandedIdentifiers = append(expandedIdentifiers, values...)
		} else {
			expandedIdentifiers = append(expandedIdentifiers, identifier.(string))
		}
	}
	return expandedIdentifiers
}

func expandDLPDictionariesPhrases(d *schema.ResourceData) []dlpdictionaries.Phrases {
	var dlpPhraseItems []dlpdictionaries.Phrases
	dlpPhraseInterface, ok := d.GetOk("phrases")
	if !ok {
		return dlpPhraseItems
	}
	dlpPhrases, ok := dlpPhraseInterface.(*schema.Set)
	if !ok {
		return dlpPhraseItems
	}
	for _, dlpItemObj := range dlpPhrases.List() {
		dlpItem, ok := dlpItemObj.(map[string]interface{})
		if !ok {
			return dlpPhraseItems
		}
		dlpPhraseItems = append(dlpPhraseItems, dlpdictionaries.Phrases{
			Action: dlpItem["action"].(string),
			Phrase: dlpItem["phrase"].(string),
		})
	}
	return dlpPhraseItems
}

func expandDLPDictionariesPatterns(d *schema.ResourceData) []dlpdictionaries.Patterns {
	var dlpPatternsItems []dlpdictionaries.Patterns
	dlpPatternsInterface, ok := d.GetOk("patterns")
	if !ok {
		return dlpPatternsItems
	}
	dlpPatterns, ok := dlpPatternsInterface.(*schema.Set)
	if !ok {
		return dlpPatternsItems
	}
	for _, patternObj := range dlpPatterns.List() {
		dlpItem, ok := patternObj.(map[string]interface{})
		if !ok {
			return dlpPatternsItems
		}
		dlpPatternsItems = append(dlpPatternsItems, dlpdictionaries.Patterns{
			Action:  dlpItem["action"].(string),
			Pattern: dlpItem["pattern"].(string),
		})
	}
	return dlpPatternsItems
}

func expandEDMDetails(d *schema.ResourceData) []dlpdictionaries.EDMMatchDetails {
	var dlpEdmDetails []dlpdictionaries.EDMMatchDetails
	dlpEdmInterface, ok := d.GetOk("exact_data_match_details")
	if !ok {
		return dlpEdmDetails
	}
	dlpEdmDetailSet, ok := dlpEdmInterface.(*schema.Set)
	if !ok {
		return dlpEdmDetails
	}
	for _, dlpEdmDetailObj := range dlpEdmDetailSet.List() {
		dlpEdmItem, ok := dlpEdmDetailObj.(map[string]interface{})
		if !ok {
			return dlpEdmDetails
		}
		secFields := []int{}
		for _, i := range dlpEdmItem["secondary_fields"].([]interface{}) {
			value, ok := i.(int)
			if !ok {
				continue
			}
			secFields = append(secFields, value)
		}
		dlpEdmDetails = append(dlpEdmDetails, dlpdictionaries.EDMMatchDetails{
			DictionaryEdmMappingID: dlpEdmItem["dictionary_edm_mapping_id"].(int),
			SchemaID:               dlpEdmItem["schema_id"].(int),
			PrimaryField:           dlpEdmItem["primary_field"].(int),
			SecondaryFields:        secFields,
			SecondaryFieldMatchOn:  dlpEdmItem["secondary_field_match_on"].(string),
		})
	}
	return dlpEdmDetails
}

func expandIDMProfileMatchAccuracy(d *schema.ResourceData) []dlpdictionaries.IDMProfileMatchAccuracy {
	var idmProfileMatchAccuracies []dlpdictionaries.IDMProfileMatchAccuracy
	dlpEdmInterface, ok := d.GetOk("idm_profile_match_accuracy")
	if !ok {
		return idmProfileMatchAccuracies
	}
	dlpEdmDetailSet, ok := dlpEdmInterface.(*schema.Set)
	if !ok {
		return idmProfileMatchAccuracies
	}
	for _, dlpEdmDetailObj := range dlpEdmDetailSet.List() {
		dlpEdmItem, ok := dlpEdmDetailObj.(map[string]interface{})
		if !ok {
			return idmProfileMatchAccuracies
		}
		var profile *common.IDNameExtensions
		profiles := expandIDMProfile(dlpEdmItem, "adp_idm_profile")
		if len(profiles) > 0 {
			profile = &profiles[0]
		}
		idmProfileMatchAccuracies = append(idmProfileMatchAccuracies, dlpdictionaries.IDMProfileMatchAccuracy{
			MatchAccuracy: dlpEdmItem["match_accuracy"].(string),
			AdpIdmProfile: profile,
		})
	}
	return idmProfileMatchAccuracies
}

func expandIDMProfile(m map[string]interface{}, key string) []common.IDNameExtensions {
	setInterface, ok := m[key]
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.IDNameExtensions
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil {
				result = append(result, common.IDNameExtensions{
					ID:         itemMap["id"].(int),
					Extensions: itemMap["extensions"].(map[string]interface{}),
				})
			}
		}
		return result
	}
	return []common.IDNameExtensions{}
}

func flattenIDNameExtensionSimple(list []common.IDNameExtensions) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		r := map[string]interface{}{
			"id": val.ID,
		}
		if val.Extensions != nil {
			r["extensions"] = val.Extensions
		}
		flattenedList[i] = r
	}
	return flattenedList
}

func flattenIDMProfileMatchAccuracySimple(edm *dlpdictionaries.DlpDictionary) []interface{} {
	idmProfileMatchAccuracies := make([]interface{}, len(edm.IDMProfileMatchAccuracy))
	for i, val := range edm.IDMProfileMatchAccuracy {
		exts := []common.IDNameExtensions{}
		if val.AdpIdmProfile != nil {
			exts = append(exts, *val.AdpIdmProfile)
		}

		idmProfileMatchAccuracies[i] = map[string]interface{}{
			"match_accuracy":  val.MatchAccuracy,
			"adp_idm_profile": flattenIDNameExtensionSimple(exts),
		}
	}

	return idmProfileMatchAccuracies
}
