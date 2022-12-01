package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/intermediatecacertificates"
)

func resourceIntermediateCACertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceIntermediateCACertificateCreate,
		Read:   resourceIntermediateCACertificateRead,
		Update: resourceIntermediateCACertificateUpdate,
		Delete: resourceIntermediateCACertificateDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("cert_id", idInt)
				} else {
					resp, err := zClient.intermediatecacertificates.GetByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("cert_id", resp.ID)
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
				Description: "Unique identifier for the intermediate CA certificate",
			},
			"cert_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifier for the intermediate CA certificate",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the intermediate CA certificate",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description for the intermediate CA certificate",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the intermediate CA certificate. Available types: Zscaler’s intermediate CA certificate (provided by Zscaler), custom intermediate certificate with software protection, and custom intermediate certificate with cloud HSM protection.",
				ValidateFunc: validation.StringInSlice([]string{
					"ZSCALER",
					"CUSTOM_SW",
					"CUSTOM_HSM",
				}, false),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of the HSM resources. Required for custom intermediate CA certificates with cloud HSM protection.",
				ValidateFunc: validation.StringInSlice([]string{
					"GLOBAL",
					"ASIA",
					"EUROPE",
					"US",
				}, false),
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Determines whether the intermediate CA certificate is enabled or disabled for SSL inspection. Subscription to cloud HSM protection allows a maximum of four active certificates for SSL inspection at a time, whereas software protection subscription allows only one active certificate.",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"default_certificate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, the intermediate CA certificate is the default intermediate certificate. Only one certificate can be marked as the default intermediate certificate at a time.",
			},
			"current_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tracks the progress of the intermediate CA certificate in the configuration workflow",
				ValidateFunc: validation.StringInSlice([]string{
					"GENERAL_DONE",
					"KEYGEN_DONE",
					"PUBKEY_DONE",
					"ATTESTATION_DONE",
					"ATTESTATION_VERIFY_DONE",
					"CSRGEN_DONE",
					"INTCERT_UPLOAD_DONE",
					"CERTCHAIN_UPLOAD_DONE",
					"CERT_READY",
				}, false),
			},
		},
	}
}

func resourceIntermediateCACertificateCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandIntermediateCACertificate(d)
	log.Printf("[INFO] Creating intermediate ca certificate\n%+v\n", req)

	resp, err := zClient.intermediatecacertificates.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia intermediate ca certificate request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("cert_id", resp.ID)
	return resourceIntermediateCACertificateRead(d, m)
}

func resourceIntermediateCACertificateRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "cert_id")
	if !ok {
		return fmt.Errorf("no intermediate ca certificate id is set")
	}
	resp, err := zClient.intermediatecacertificates.Get(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia intermediate ca certificate %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting intermediate ca certificate :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("type", resp.Type)
	_ = d.Set("region", resp.Region)
	_ = d.Set("status", resp.Status)
	_ = d.Set("default_certificate", resp.DefaultCertificate)
	_ = d.Set("current_state", resp.CurrentState)

	return nil
}

func resourceIntermediateCACertificateUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "cert_id")
	if !ok {
		log.Printf("[ERROR] intermediate ca certificate ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating intermediate ca certificate ID: %v\n", id)
	req := expandIntermediateCACertificate(d)
	if _, err := zClient.intermediatecacertificates.Get(req.ID); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := zClient.intermediatecacertificates.Update(id, &req); err != nil {
		return err
	}

	return resourceIntermediateCACertificateRead(d, m)
}

func resourceIntermediateCACertificateDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "cert_id")
	if !ok {
		log.Printf("[ERROR] intermediate ca certificate id ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting intermediate ca certificate ID: %v\n", (d.Id()))

	if _, err := zClient.intermediatecacertificates.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] intermediate ca certificate deleted")
	return nil
}

func expandIntermediateCACertificate(d *schema.ResourceData) intermediatecacertificates.IntermediateCACertificate {
	return intermediatecacertificates.IntermediateCACertificate{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Type:               d.Get("type").(string),
		Region:             d.Get("region").(string),
		Status:             d.Get("status").(string),
		DefaultCertificate: d.Get("default_certificate").(bool),
		CurrentState:       d.Get("current_state").(string),
	}
}
