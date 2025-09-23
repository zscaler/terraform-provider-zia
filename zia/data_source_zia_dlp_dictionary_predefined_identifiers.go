package zia

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlpdictionaries"
)

func dataSourceDLPDictionaryPredefinedIdentifiers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDLPDictionaryPredefinedIdentifiersRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePredefinedIdentifier,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"predefined_identifiers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDLPDictionaryPredefinedIdentifiersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	dictionaryName := d.Get("name").(string)

	identifiers, dictionaryID, err := dlpdictionaries.GetPredefinedIdentifiers(ctx, service, dictionaryName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(dictionaryID))

	if err := d.Set("predefined_identifiers", identifiers); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved predefined identifiers for dictionary %s: %v", dictionaryName, identifiers)

	return nil
}
