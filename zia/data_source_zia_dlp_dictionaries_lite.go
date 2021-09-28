package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/dlpdictionaries"
)

func dataSourceDLPDictionariesLite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDLPDictionariesLiteRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
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
				Computed: true,
			},
			"phrases": {
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
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
			"custom": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dictionary_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDLPDictionariesLiteRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *dlpdictionaries.DlpDictionary
	idObj, idSet := d.GetOk("id")
	id, idIsInt := idObj.(int)
	if idSet && idIsInt && id > 0 {
		log.Printf("[INFO] Getting data for vpn credential id: %d\n", id)
		res, err := zClient.dlpdictionaries.GetDlpDictionaryLite(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for vpn credential fqdn: %s\n", name)
		res, err := zClient.dlpdictionaries.GetDlpDictionaryLiteByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("confidence_threshold", resp.ConfidenceThreshold)
		_ = d.Set("custom_phrase_match_type", resp.CustomPhraseMatchType)
		_ = d.Set("name_l10n_tag", resp.NameL10nTag)
		_ = d.Set("custom", resp.Custom)
		_ = d.Set("threshold_type", resp.ThresholdType)
		_ = d.Set("dictionary_type", resp.DictionaryType)
		if err := d.Set("phrases", flattenPhrases(resp)); err != nil {
			return err
		}

		if err := d.Set("patterns", flattenPatterns(resp)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any dlp dictionary with name '%s' or id '%d'", name, id)
	}

	return nil
}
