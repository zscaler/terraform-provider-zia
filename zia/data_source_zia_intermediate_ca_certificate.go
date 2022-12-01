package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zia/services/intermediatecacertificates"
)

func dataSourceIntermediateCACertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIntermediateCACertificateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_certificate": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cert_start_date": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cert_exp_date": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"current_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_generation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hsm_attestation_verified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"csr_file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"csr_generation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceIntermediateCACertificateRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *intermediatecacertificates.IntermediateCACertificate
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting intermediate ca certificate id: %d\n", id)
		res, err := zClient.intermediatecacertificates.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting intermediate ca certificate : %s\n", name)
		res, err := zClient.intermediatecacertificates.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("type", resp.Type)
		_ = d.Set("region", resp.Region)
		_ = d.Set("status", resp.Status)
		_ = d.Set("default_certificate", resp.DefaultCertificate)
		_ = d.Set("cert_start_date", resp.CertStartDate)
		_ = d.Set("cert_exp_date", resp.CertExpDate)
		_ = d.Set("current_state", resp.CurrentState)
		_ = d.Set("public_key", resp.PublicKey)
		_ = d.Set("key_generation_time", resp.KeyGenerationTime)
		_ = d.Set("hsm_attestation_verified_time", resp.HSMAttestationVerifiedTime)
		_ = d.Set("csr_file_name", resp.CSRFileName)
		_ = d.Set("csr_generation_time", resp.CSRGenerationTime)

	} else {
		return fmt.Errorf("couldn't find any time window with name '%s' or id '%d'", name, id)
	}

	return nil
}
