package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_engines"
)

func dataSourceDLPEngines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDLPEnginesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier for the DLP engine.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP engine name as configured by the admin. This attribute is required in POST and PUT requests for custom DLP engines.",
			},
			"predefined_engine_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the predefined DLP engine.",
			},
			"engine_expression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The boolean logical operator in which various DLP dictionaries are combined within a DLP engine's expression.",
			},
			"custom_dlp_engine": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is a custom DLP engine. If this value is set to true, the engine is custom.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The DLP engine's description.",
			},
		},
	}
}

func dataSourceDLPEnginesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *dlp_engines.DLPEngines
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for dlp engine id: %d\n", id)
		res, err := zClient.dlp_engines.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for dlp engine name: %s\n", name)
		res, err := zClient.dlp_engines.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("predefined_engine_name", resp.PredefinedEngineName)
		_ = d.Set("engine_expression", resp.EngineExpression)
		_ = d.Set("custom_dlp_engine", resp.CustomDlpEngine)

	} else {
		return fmt.Errorf("couldn't find any dlp engine name '%s' or id '%d'", name, id)
	}

	return nil
}
