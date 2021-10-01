package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/dlpdictionaries"
)

func resourceDLPDictionaries() *schema.Resource {
	return &schema.Resource{
		Create:   resourceDLPDictionariesCreate,
		Read:     resourceDLPDictionariesRead,
		Update:   resourceDLPDictionariesUpdate,
		Delete:   resourceDLPDictionariesDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			// "id": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"confidence_threshold": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pattern": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name_l10n_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"threshold_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dictionary_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PATTERNS_AND_PHRASES",
					"EXACT_DATA_MATCH",
					"INDEXED_DATA_MATCH",
				}, false),
			},
		},
	}
}

func resourceDLPDictionariesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandDLPDictionaries(d)
	log.Printf("[INFO] Creating zia dlp dictionaries\n%+v\n", req)

	resp, err := zClient.dlpdictionaries.CreateDlpDictionary(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia dlp dictionaries request. ID: %v\n", resp)

	return resourceDLPDictionariesRead(d, m)
}

func resourceDLPDictionariesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.dlpdictionaries.GetDlpDictionaries(d.Id())

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing dlp dictionary %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting dlp dictionary :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("confidence_threshold", resp.ConfidenceThreshold)
	_ = d.Set("custom_phrase_match_type", resp.CustomPhraseMatchType)
	_ = d.Set("name_l10n_tag", resp.NameL10nTag)
	_ = d.Set("threshold_type", resp.ThresholdType)
	_ = d.Set("dictionary_type", resp.DictionaryType)
	if err := d.Set("phrases", flattenPhrases(resp)); err != nil {
		return err
	}

	if err := d.Set("patterns", flattenPatterns(resp)); err != nil {
		return err
	}

	return nil
}

func resourceDLPDictionariesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating vpn credentials ID: %v\n", id)
	req := expandDLPDictionaries(d)

	if _, err := zClient.dlpdictionaries.UpdateDlpDictionary(id, &req); err != nil {
		return err
	}

	return resourceDLPDictionariesRead(d, m)
}

func resourceDLPDictionariesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	// Need to pass the ID (int) of the resource for deletion
	log.Printf("[INFO] Deleting vpn credentials ID: %v\n", (d.Id()))

	if _, err := zClient.dlpdictionaries.DeleteDlpDictionary(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] vpn credentials deleted")
	return nil
}

func expandDLPDictionaries(d *schema.ResourceData) dlpdictionaries.DlpDictionary {
	dlpDictionaries := dlpdictionaries.DlpDictionary{
		Phrases:  expandDLPDictionariesPhrases(d),
		Patterns: expandDLPDictionariesPatterns(d),
	}
	return dlpDictionaries
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
