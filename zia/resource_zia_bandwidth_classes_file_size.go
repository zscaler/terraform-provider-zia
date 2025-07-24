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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_classes"
)

func resourceBandwdithClassesFileSize() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceBandwdithClassesFileSizeRead,
		CreateContext: resourceBandwdithClassesFileSizeCreate,
		UpdateContext: resourceBandwdithClassesFileSizeUpdate,
		DeleteContext: resourceFuncNoOp,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("class_id", idInt)
				} else {
					resp, err := bandwidth_classes.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("class_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "System-generated identifier for the bandwidth class",
			},
			"class_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "System-generated identifier for the bandwidth class",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the bandwidth class",
				Default:     "BANDWIDTH_CAT_LARGE_FILE",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "BANDWIDTH_CAT_LARGE_FILE",
				Description: "The application type for which the bandwidth class is configured",
			},
			"file_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The file size for a bandwidth class",
				ValidateFunc: validation.StringInSlice([]string{
					"FILE_5MB",
					"FILE_10MB",
					"FILE_50MB",
					"FILE_100MB",
					"FILE_250MB",
					"FILE_500MB",
					"FILE_1GB",
				}, false),
			},
		},
	}
}

func resourceBandwdithClassesFileSizeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	name := d.Get("name").(string)
	if name == "" {
		return diag.Errorf("'name' must be provided to locate the existing bandwidth class")
	}

	// Lookup by name to get the ID
	existing, err := bandwidth_classes.GetByName(ctx, service, name)
	if err != nil {
		return diag.Errorf("failed to find existing bandwidth class by name %q: %v", name, err)
	}

	log.Printf("[INFO] Found existing bandwidth class %q with ID %d", existing.Name, existing.ID)
	d.SetId(strconv.Itoa(existing.ID))
	_ = d.Set("class_id", existing.ID)

	// Expand request using the retrieved ID
	req := expandBandwidthClassesFileSize(d)
	req.ID = existing.ID // Ensure ID is set

	log.Printf("[INFO] Updating bandwidth class %q via PUT:\n%+v\n", name, req)

	if _, _, err := bandwidth_classes.Update(ctx, service, req.ID, &req); err != nil {
		return diag.FromErr(err)
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceBandwdithClassesFileSizeRead(ctx, d, meta)
}

func resourceBandwdithClassesFileSizeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "class_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no bandwidth class file size id is set"))
	}
	resp, err := bandwidth_classes.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia bandwidth class file size %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia bandwidth class file size:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("class_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("type", resp.Type)
	_ = d.Set("file_size", resp.FileSize)

	return nil
}

func resourceBandwdithClassesFileSizeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var id int
	var ok bool

	// Get ID from state
	id, ok = getIntFromResourceData(d, "class_id")
	if !ok || id == 0 {
		name := d.Get("name").(string)
		if name == "" {
			return diag.Errorf("either 'class_id' or 'name' must be set for update")
		}

		existing, err := bandwidth_classes.GetByName(ctx, service, name)
		if err != nil {
			return diag.Errorf("failed to find bandwidth class with name %q: %v", name, err)
		}

		id = existing.ID
		d.SetId(strconv.Itoa(id))
		_ = d.Set("class_id", id)
		log.Printf("[INFO] Retrieved class ID %d for update by name %q", id, name)
	}

	log.Printf("[INFO] Updating ZIA bandwidth class ID: %d\n", id)

	req := expandBandwidthClassesFileSize(d)
	req.ID = id

	if _, err := bandwidth_classes.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if _, _, err := bandwidth_classes.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceBandwdithClassesFileSizeRead(ctx, d, meta)
}

func expandBandwidthClassesFileSize(d *schema.ResourceData) bandwidth_classes.BandwidthClasses {
	return bandwidth_classes.BandwidthClasses{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		FileSize: d.Get("file_size").(string),
	}
}
