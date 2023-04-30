package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlpdictionaries"
)

func dataSourceDLPDictionaries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDLPDictionariesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"confidence_threshold": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phrases": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phrase": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"custom_phrase_match_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"patterns": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name_l10n_tag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"threshold_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dictionary_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exact_data_match_details": {
				Type:        schema.TypeSet,
				Computed:    true,
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
							Computed:    true,
							Description: "The EDM template's primary field.",
						},
						"secondary_fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"secondary_field_match_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The EDM secondary field to match on.",
						},
					},
				},
			},
			"idm_profile_match_accuracy": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of Indexed Document Match (IDM) profiles and their corresponding match accuracy for custom DLP dictionaries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"adp_idm_profile": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The action applied to a DLP dictionary using patterns",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Identifier that uniquely identifies an entity",
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
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
						"match_accuracy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IDM template match accuracy.",
						},
					},
				},
			},
			"proximity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The DLP dictionary proximity length.",
			},
			"ignore_exact_match_idm_dict": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to exclude documents that are a 100% match to already-indexed documents from triggering an Indexed Document Match (IDM) Dictionary.",
			},
			"include_bin_numbers": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A true value denotes that the specified Bank Identification Number (BIN) values are included in the Credit Cards dictionary. A false value denotes that the specified BIN values are excluded from the Credit Cards dictionary.Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.",
			},
			"bin_numbers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "The list of Bank Identification Number (BIN) values that are included or excluded from the Credit Cards dictionary. BIN values can be specified only for Diners Club, Mastercard, RuPay, and Visa cards. Up to 512 BIN values can be configured in a dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.",
			},
			"dict_template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the predefined dictionary (original source dictionary) that is used for cloning. This field is applicable only to cloned dictionaries. Only a limited set of identification-based predefined dictionaries (e.g., Credit Cards, Social Security Numbers, National Identification Numbers, etc.) can be cloned. Up to 4 clones can be created from a predefined dictionary.",
			},
			"predefined_clone": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This field is set to true if the dictionary is cloned from a predefined dictionary. Otherwise, it is set to false.",
			},
			"proximity_length_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This value is set to true if proximity length and high confidence phrases are enabled for the DLP dictionary.",
			},
		},
	}
}

func dataSourceDLPDictionariesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *dlpdictionaries.DlpDictionary
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for dlp dictionary id: %d\n", id)
		res, err := zClient.dlpdictionaries.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for dlp dictionary: %s\n", name)
		res, err := zClient.dlpdictionaries.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("custom", resp.Custom)
		_ = d.Set("confidence_threshold", resp.ConfidenceThreshold)
		_ = d.Set("custom_phrase_match_type", resp.CustomPhraseMatchType)
		_ = d.Set("name_l10n_tag", resp.NameL10nTag)
		_ = d.Set("threshold_type", resp.ThresholdType)
		_ = d.Set("dictionary_type", resp.DictionaryType)
		_ = d.Set("ignore_exact_match_idm_dict", resp.IgnoreExactMatchIdmDict)
		_ = d.Set("include_bin_numbers", resp.IncludeBinNumbers)
		_ = d.Set("bin_numbers", resp.BinNumbers)
		_ = d.Set("dict_template_id", resp.DictTemplateId)
		_ = d.Set("predefined_clone", resp.PredefinedClone)
		_ = d.Set("proximity_length_enabled", resp.ProximityLengthEnabled)
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
		if err := d.Set("idm_profile_match_accuracy", flattenIDMProfileMatchAccuracy(resp)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any dlp dictionary with name '%s' or id '%d'", name, id)
	}

	return nil
}

func flattenPhrases(phrases *dlpdictionaries.DlpDictionary) []interface{} {
	dlpPhrases := make([]interface{}, len(phrases.Phrases))
	for i, val := range phrases.Phrases {
		dlpPhrases[i] = map[string]interface{}{
			"action": val.Action,
			"phrase": val.Phrase,
		}
	}

	return dlpPhrases
}

func flattenPatterns(patterns *dlpdictionaries.DlpDictionary) []interface{} {
	dlpPatterns := make([]interface{}, len(patterns.Patterns))
	for i, val := range patterns.Patterns {
		dlpPatterns[i] = map[string]interface{}{
			"action":  val.Action,
			"pattern": val.Pattern,
		}
	}

	return dlpPatterns
}

func flattenEDMDetails(edm *dlpdictionaries.DlpDictionary) []interface{} {
	edmDetails := make([]interface{}, len(edm.EDMMatchDetails))
	for i, val := range edm.EDMMatchDetails {
		edmDetails[i] = map[string]interface{}{
			"dictionary_edm_mapping_id": val.DictionaryEdmMappingID,
			"schema_id":                 val.SchemaID,
			"primary_field":             val.PrimaryField,
			"secondary_fields":          val.SecondaryFields,
			"secondary_field_match_on":  val.SecondaryFieldMatchOn,
		}
	}

	return edmDetails
}

func flattenIDMProfileMatchAccuracy(edm *dlpdictionaries.DlpDictionary) []interface{} {
	idmProfileMatchAccuracies := make([]interface{}, len(edm.IDMProfileMatchAccuracy))
	for i, val := range edm.IDMProfileMatchAccuracy {
		exts := []common.IDNameExtensions{}
		if val.AdpIdmProfile != nil {
			exts = append(exts, *val.AdpIdmProfile)
		}

		idmProfileMatchAccuracies[i] = map[string]interface{}{
			"match_accuracy":  val.MatchAccuracy,
			"adp_idm_profile": flattenIDNameExtensions(exts),
		}
	}

	return idmProfileMatchAccuracies
}
