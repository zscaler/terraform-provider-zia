package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_engines"
)

func dataSourceDLPEngines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDLPEnginesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier for the DLP engine.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP engine name as configured by the admin.",
			},
			"predefined_engine_name": {
				Type:        schema.TypeString,
				Optional:    true,
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

func dataSourceDLPEnginesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *dlp_engines.DLPEngines
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for dlp engine id: %d\n", id)
		res, err := dlp_engines.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for dlp engine name: %s\n", name)
		res, err := dlp_engines.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	predefined, _ := d.Get("predefined_engine_name").(string)
	if resp == nil && predefined != "" {
		log.Printf("[INFO] Getting data for predefined dlp engine name: %s\n", predefined)
		res, err := dlp_engines.GetByPredefinedEngine(ctx, service, predefined)
		if err != nil {
			return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("couldn't find any dlp engine name '%s' or id '%d'", name, id))
	}

	return nil
}
