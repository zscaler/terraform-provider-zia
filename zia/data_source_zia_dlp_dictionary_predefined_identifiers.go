package zia

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlpdictionaries"
)

func dataSourceDLPDictionaryPredefinedIdentifiers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDLPDictionaryPredefinedIdentifiersRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePredefinedIdentifier,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"predefined_identifiers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDLPDictionaryPredefinedIdentifiersRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlpdictionaries

	dictionaryName := d.Get("name").(string)

	identifiers, dictionaryID, err := dlpdictionaries.GetPredefinedIdentifiers(service, dictionaryName)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(dictionaryID))

	if err := d.Set("predefined_identifiers", identifiers); err != nil {
		return err
	}

	log.Printf("[INFO] Retrieved predefined identifiers for dictionary %s: %v", dictionaryName, identifiers)

	return nil
}
