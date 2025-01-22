package zia

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/activation"
)

func resourceActivationStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActivationCreate,
		ReadContext:   resourceActivationRead,
		DeleteContext: resourceActivationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Organization Policy Edit/Update Activation status",
				ValidateFunc: validation.StringInSlice([]string{
					"ACTIVE",
				}, false),
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceActivationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandActivation(d)
	log.Printf("[INFO] Performing configuration activation\n%+v\n", req)

	resp, err := activation.CreateActivation(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Configuration activation successfull. %v\n", resp.Status)
	d.SetId("activation")
	return resourceActivationRead(ctx, d, meta)
}

func resourceActivationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := activation.GetActivationStatus(ctx, service)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Cannot obtain activation %s from ZIA", d.Id())
			// Activation is not an actual object; hence no ID should be set.
			// d.SetId("")
			// return nil
		}

		return diag.FromErr(err)
	}
	log.Printf("[INFO] Reading activation status: %+v\n", resp)
	_ = d.Set("status", resp.Status)

	return nil
}

func resourceActivationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Delete doesn't actually do anything, because an activation can't be deleted.
	return nil
}

func expandActivation(d *schema.ResourceData) activation.Activation {
	return activation.Activation{
		Status: d.Get("status").(string),
	}
}
