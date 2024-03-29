package zia

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/vpncredentials"
)

func resourceTrafficForwardingVPNCredentials() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrafficForwardingVPNCredentialsCreate,
		Read:   resourceTrafficForwardingVPNCredentialsRead,
		Update: resourceTrafficForwardingVPNCredentialsUpdate,
		Delete: resourceTrafficForwardingVPNCredentialsDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()

				// Try to import by FQDN
				vpnCredential, err := zClient.vpncredentials.GetByFQDN(id)
				if err == nil {
					d.SetId(strconv.Itoa(vpnCredential.ID))
					_ = d.Set("vpn_id", vpnCredential.ID)
					return []*schema.ResourceData{d}, nil
				}

				// Try to import by IP
				vpnCredential, err = zClient.vpncredentials.GetByIP(id)
				if err == nil {
					d.SetId(strconv.Itoa(vpnCredential.ID))
					_ = d.Set("vpn_id", vpnCredential.ID)
					return []*schema.ResourceData{d}, nil
				}

				// Try to import by VPN Type
				vpnCredential, err = zClient.vpncredentials.GetVPNByType(id)
				if err == nil {
					d.SetId(strconv.Itoa(vpnCredential.ID))
					_ = d.Set("vpn_id", vpnCredential.ID)
					return []*schema.ResourceData{d}, nil
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"vpn_id": {
				Type:     schema.TypeInt,
				Computed: true,
				// ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IP",
					"UFQDN",
				}, false),
			},
			"fqdn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ForceNew:     true,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				ForceNew:     true,
			},
			"pre_shared_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"comments": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
		},
	}
}

func resourceTrafficForwardingVPNCredentialsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandVPNCredentials(d)
	log.Printf("[INFO] Creating zia vpn credentials\n%+v\n", req)

	if err := validateVpnCredentialType(req); err != nil {
		return err
	}

	resp, _, err := zClient.vpncredentials.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia vpn credentials request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("vpn_id", resp.ID)

	// Sleep for 2 seconds before triggering the activation
	time.Sleep(2 * time.Second)

	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return resourceTrafficForwardingVPNCredentialsRead(d, m)
}

func validateVpnCredentialType(vpn vpncredentials.VPNCredentials) error {
	if vpn.Type == "IP" && vpn.IPAddress == "" {
		return fmt.Errorf("invalid input argument, ip_address is required")
	}

	if vpn.Type == "UFQDN" && vpn.FQDN == "" {
		return fmt.Errorf("invalid input argument, fqdn attribute is required")
	}
	return nil
}

func resourceTrafficForwardingVPNCredentialsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "vpn_id")
	if !ok {
		return fmt.Errorf("no Traffic Forwarding zia vpn credentials id is set")
	}
	resp, err := zClient.vpncredentials.Get(id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing vpn credentials %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting vpn credentials:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("vpn_id", resp.ID)
	_ = d.Set("type", resp.Type)
	_ = d.Set("fqdn", resp.FQDN)
	_ = d.Set("ip_address", resp.IPAddress)
	_ = d.Set("comments", resp.Comments)

	return nil
}

func resourceTrafficForwardingVPNCredentialsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "vpn_id")
	if !ok {
		log.Printf("[ERROR] vpn credentials ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating vpn credentials ID: %v\n", id)
	req := expandVPNCredentials(d)
	if _, err := zClient.vpncredentials.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := zClient.vpncredentials.Update(id, &req); err != nil {
		return err
	}
	// Sleep for 2 seconds before triggering the activation
	time.Sleep(2 * time.Second)

	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}

	return resourceTrafficForwardingVPNCredentialsRead(d, m)
}

func resourceTrafficForwardingVPNCredentialsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "vpn_id")
	if !ok {
		log.Printf("[ERROR] vpn credentials ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting vpn credentials ID: %v\n", (d.Id()))

	if err := zClient.vpncredentials.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] vpn credentials deleted")

	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return nil
}

func expandVPNCredentials(d *schema.ResourceData) vpncredentials.VPNCredentials {
	id, _ := getIntFromResourceData(d, "vpn_id")
	result := vpncredentials.VPNCredentials{
		ID:           id,
		Type:         d.Get("type").(string),
		FQDN:         d.Get("fqdn").(string),
		IPAddress:    d.Get("ip_address").(string),
		PreSharedKey: d.Get("pre_shared_key").(string),
		Comments:     d.Get("comments").(string),
	}

	return result
}
