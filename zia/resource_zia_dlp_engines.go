package zia

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_engines"
)

func resourceDLPEngines() *schema.Resource {
	return &schema.Resource{
		Create: resourceDLPEnginesCreate,
		Read:   resourceDLPEnginesRead,
		Update: resourceDLPEnginesUpdate,
		Delete: resourceDLPEnginesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)
				service := zClient.dlp_engines

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("engine_id", idInt)
				} else {
					resp, err := dlp_engines.GetByName(service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("engine_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DLP engine name as configured by the admin.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP engine's description.",
			},
			"engine_expression": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The boolean logical operator in which various DLP dictionaries are combined within a DLP engine's expression.",
			},
			"custom_dlp_engine": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether this is a custom DLP engine. If this value is set to true, the engine is custom.",
			},
		},
	}
}

func resourceDLPEnginesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlp_engines

	req := expandDLPEngines(d)
	log.Printf("[INFO] Creating zia dlp engine\n%+v\n", req)

	resp, _, err := dlp_engines.Create(service, &req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia dlp engine request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("engine_id", resp.ID)
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceDLPEnginesRead(d, m)
}

func resourceDLPEnginesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlp_engines

	id, ok := getIntFromResourceData(d, "engine_id")
	if !ok {
		return fmt.Errorf("no dlp engine id is set")
	}
	resp, err := dlp_engines.Get(service, id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia dlp engine%s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting zia dlp engine:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("engine_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("engine_expression", resp.EngineExpression)
	_ = d.Set("custom_dlp_engine", resp.CustomDlpEngine)

	return nil
}

func resourceDLPEnginesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlp_engines

	id, ok := getIntFromResourceData(d, "engine_id")
	if !ok {
		log.Printf("[ERROR] dlp engine ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia dlp engine ID: %v\n", id)
	req := expandDLPEngines(d)
	if _, err := dlp_engines.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := dlp_engines.Update(service, id, &req); err != nil {
		return err
	}
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceDLPEnginesRead(d, m)
}

func resourceDLPEnginesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlp_engines

	id, ok := getIntFromResourceData(d, "engine_id")
	if !ok {
		log.Printf("[ERROR] dlp engine ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia dlp engine ID: %v\n", (d.Id()))

	if _, err := dlp_engines.Delete(service, id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] zia dlp engine deleted")
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandDLPEngines(d *schema.ResourceData) dlp_engines.DLPEngines {
	id, _ := getIntFromResourceData(d, "engine_id")
	result := dlp_engines.DLPEngines{
		ID:               id,
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		EngineExpression: d.Get("engine_expression").(string),
		CustomDlpEngine:  d.Get("custom_dlp_engine").(bool),
	}
	return result
}
