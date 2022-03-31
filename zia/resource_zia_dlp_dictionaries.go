package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/dlpdictionaries"
)

func resourceDLPDictionaries() *schema.Resource {
	return &schema.Resource{
		Create: resourceDLPDictionariesCreate,
		Read:   resourceDLPDictionariesRead,
		Update: resourceDLPDictionariesUpdate,
		Delete: resourceDLPDictionariesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("dictionary_id", id)
				} else {
					resp, err := zClient.dlpdictionaries.GetByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						d.Set("dictionary_id", resp.ID)
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP dictionary's name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The desciption of the DLP dictionary",
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"phrase": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"custom_phrase_match_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY",
					"MATCH_ANY_CUSTOM_PHRASE_PATTERN_DICTIONARY",
				}, false),
			},
			"patterns": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List containing the patterns used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The action applied to a DLP dictionary using patterns",
						},
						"pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "DLP dictionary pattern",
						},
					},
				},
			},
			"name_l10n_tag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the name is localized or not. This is always set to True for predefined DLP dictionaries.",
			},
			"threshold_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DLP threshold type",
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
			"exact_data_match_details": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Exact Data Match (EDM) related information for custom DLP dictionaries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dictionary_edm_mapping_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier for the EDM mapping",
						},
						"schema_id": {
							Type:        schema.TypeInt,
							Computed:    true,
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
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of Indexed Document Match (IDM) profiles and their corresponding match accuracy for custom DLP dictionaries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"adp_idm_profile": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The action applied to a DLP dictionary using patterns",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Identifier that uniquely identifies an entity",
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
			"match_accuracy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IDM template match accuracy.",
				ValidateFunc: validation.StringInSlice([]string{
					"LOW", "MEDIUM", "HEAVY",
				}, false),
			},
			"proximity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The DLP dictionary proximity length.",
			},
		},
	}

}

func resourceDLPDictionariesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandDLPDictionaries(d)
	log.Printf("[INFO] Creating zia dlp dictionaries\n%+v\n", req)

	resp, _, err := zClient.dlpdictionaries.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia dlp dictionaries request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("dictionary_id", resp.ID)

	return resourceDLPDictionariesRead(d, m)
}

func resourceDLPDictionariesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "dictionary_id")
	if !ok {
		return fmt.Errorf("no DLP dictionary id is set")
	}
	resp, err := zClient.dlpdictionaries.Get(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing dlp dictionary %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting dlp dictionary :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("dictionary_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("confidence_threshold", resp.ConfidenceThreshold)
	_ = d.Set("custom_phrase_match_type", resp.CustomPhraseMatchType)
	_ = d.Set("name_l10n_tag", resp.NameL10nTag)
	_ = d.Set("threshold_type", resp.ThresholdType)
	_ = d.Set("dictionary_type", resp.DictionaryType)
	_ = d.Set("proximity", resp.Proximity)
	if err := d.Set("phrases", flattenPhrases(resp)); err != nil {
		return err
	}

	if err := d.Set("patterns", flattenPatterns(resp)); err != nil {
		return err
	}
	if err := d.Set("exact_data_match_details", flattenEDMDetails(resp)); err != nil {
		return err
	}

	// Need to fully flatten and expand this menu
	// if err := d.Set("idm_profile_match_accuracy", flattenIDMProfileMatch(resp)); err != nil {
	// 	return err
	// }

	return nil
}

func resourceDLPDictionariesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "dictionary_id")
	if !ok {
		log.Printf("[ERROR] dlp dictionaryID not set: %v\n", id)
	}

	log.Printf("[INFO] Updating dlp dictionary ID: %v\n", id)
	req := expandDLPDictionaries(d)

	if _, _, err := zClient.dlpdictionaries.Update(id, &req); err != nil {
		return err
	}

	return resourceDLPDictionariesRead(d, m)
}

func resourceDLPDictionariesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "dictionary_id")
	if !ok {
		log.Printf("[ERROR] dlp dictionary ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting dlp dictionary ID: %v\n", (d.Id()))

	if _, err := zClient.dlpdictionaries.DeleteDlpDictionary(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] dlp dictionary deleted")
	return nil
}

// Need to make all below expand functions as SchemaSet

func expandDLPDictionaries(d *schema.ResourceData) dlpdictionaries.DlpDictionary {
	id, _ := getIntFromResourceData(d, "dictionary_id")
	result := dlpdictionaries.DlpDictionary{
		ID:                    id,
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		ConfidenceThreshold:   d.Get("confidence_threshold").(string),
		CustomPhraseMatchType: d.Get("custom_phrase_match_type").(string),
		DictionaryType:        d.Get("dictionary_type").(string),
	}
	phrases := expandDLPDictionariesPhrases(d)
	if phrases != nil {
		result.Phrases = phrases
	}

	patterns := expandDLPDictionariesPatterns(d)
	if phrases != nil {
		result.Patterns = patterns
	}

	edmDetails := expandEDMDetails(d)
	if edmDetails != nil {
		result.EDMMatchDetails = edmDetails
	}
	return result
}

func expandDLPDictionariesPhrases(d *schema.ResourceData) []dlpdictionaries.Phrases {
	var dlpPhraseItems []dlpdictionaries.Phrases
	if dlpPhraseInterface, ok := d.GetOk("phrases"); ok {
		dlpPhrase := dlpPhraseInterface.([]interface{})
		dlpPhraseItems = make([]dlpdictionaries.Phrases, len(dlpPhrase))
		for i, phrase := range dlpPhrase {
			dlpItem := phrase.(map[string]interface{})
			dlpPhraseItems[i] = dlpdictionaries.Phrases{
				Action: dlpItem["action"].(string),
				Phrase: dlpItem["phrase"].(string),
			}
		}
	}

	return dlpPhraseItems
}

func expandDLPDictionariesPatterns(d *schema.ResourceData) []dlpdictionaries.Patterns {
	var dlpPatternsItems []dlpdictionaries.Patterns
	if dlpPatternsInterface, ok := d.GetOk("patterns"); ok {
		dlpPattern := dlpPatternsInterface.([]interface{})
		dlpPatternsItems = make([]dlpdictionaries.Patterns, len(dlpPattern))
		for i, pattern := range dlpPattern {
			dlpItem := pattern.(map[string]interface{})
			dlpPatternsItems[i] = dlpdictionaries.Patterns{
				Action:  dlpItem["action"].(string),
				Pattern: dlpItem["pattern"].(string),
			}
		}
	}

	return dlpPatternsItems
}

func expandEDMDetails(d *schema.ResourceData) []dlpdictionaries.EDMMatchDetails {
	var dlpEdmDetails []dlpdictionaries.EDMMatchDetails
	if dlpEdmInterface, ok := d.GetOk("exact_data_match_details"); ok {
		dlpEdmDetail := dlpEdmInterface.([]interface{})
		dlpEdmDetails = make([]dlpdictionaries.EDMMatchDetails, len(dlpEdmDetail))
		for i, pattern := range dlpEdmDetail {
			dlpEdmItem := pattern.(map[string]interface{})
			dlpEdmDetails[i] = dlpdictionaries.EDMMatchDetails{
				DictionaryEdmMappingID: dlpEdmItem["dictionaryEdmMappingId"].(int),
				SchemaID:               dlpEdmItem["schema_id"].(int),
				PrimaryField:           dlpEdmItem["primary_field"].(int),
				SecondaryFields:        dlpEdmItem["secondary_fields"].([]int),
				SecondaryFieldMatchOn:  dlpEdmItem["secondary_field_match_on"].(string),
			}
		}
	}

	return dlpEdmDetails
}
