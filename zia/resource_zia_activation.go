package zia

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/activation"
)

func resourceActivationStatus() *schema.Resource {
	return &schema.Resource{
		Create:   resourceActivationCreate,
		Read:     resourceActivationRead,
		Delete:   resourceActivationDelete,
		Importer: &schema.ResourceImporter{},

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

func resourceActivationCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandActivation(d)
	log.Printf("[INFO] Performing configuration activation\n%+v\n", req)

	resp, err := zClient.activation.CreateActivation(req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Configuration activation successfull. %v\n", resp.Status)
	d.SetId("activation")
	return resourceActivationRead(d, m)
}

func resourceActivationRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.activation.GetActivationStatus()
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Cannot obtain activation %s from ZIA", d.Id())
			// Activation is not an actual object; hence no ID should be set.
			// d.SetId("")
			// return nil
		}

		return err
	}
	log.Printf("[INFO] Reading activation status: %+v\n", resp)
	_ = d.Set("status", resp.Status)

	return nil
}

func resourceActivationDelete(d *schema.ResourceData, m interface{}) error {
	// Delete doesn't actually do anything, because an activation can't be deleted.
	return nil
}

func expandActivation(d *schema.ResourceData) activation.Activation {
	return activation.Activation{
		Status: d.Get("status").(string),
	}
}
