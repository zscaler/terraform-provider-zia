package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol/custom_file_types"
)

func resourceCustomFileTypes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomFileTypesCreate,
		ReadContext:   resourceCustomFileTypesRead,
		UpdateContext: resourceCustomFileTypesUpdate,
		DeleteContext: resourceCustomFileTypesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("file_id", idInt)
				} else {
					resp, err := custom_file_types.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("file_id", resp.ID)
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
				Description: "The unique identifier for the custom file type.",
			},
			"file_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the custom file type.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the custom file type.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the custom file type, if any.",
			},
			"extension": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The file type extension. The maximum extension length is 10 characters.",
			},
			"file_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "File type ID. This ID is assigned and maintained for all file types including predefined and custom file types, and this value is different from the custom file type ID.",
			},
		},
	}
}

func resourceCustomFileTypesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandCustomFileTypes(d)
	log.Printf("[INFO] Creating ZIA custom file type\n%+v\n", req)

	resp, err := custom_file_types.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA custom file types request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("file_id", resp.ID)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceCustomFileTypesRead(ctx, d, meta)
}

func resourceCustomFileTypesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "file_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no custom file type id is set"))
	}
	resp, err := custom_file_types.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia custom file types %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia custom file types:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("file_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("extension", resp.Extension)
	_ = d.Set("file_type_id", resp.FileTypeID)

	return nil
}

func resourceCustomFileTypesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "file_id")
	if !ok {
		log.Printf("[ERROR] custom file type ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia custom file type ID: %v\n", id)
	req := expandCustomFileTypes(d)
	if _, err := custom_file_types.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := custom_file_types.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceCustomFileTypesRead(ctx, d, meta)
}

func resourceCustomFileTypesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "file_id")
	if !ok {
		log.Printf("[ERROR] custom file type ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia custom file type ID: %v\n", (d.Id()))

	if _, err := custom_file_types.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia custom file type deleted")

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandCustomFileTypes(d *schema.ResourceData) custom_file_types.CustomFileTypes {
	id, _ := getIntFromResourceData(d, "file_id")
	result := custom_file_types.CustomFileTypes{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Extension:   d.Get("extension").(string),
		FileTypeID:  d.Get("file_type_id").(int),
	}
	return result
}
