package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/locationmanagement"
)

func resourceLocationManagement() *schema.Resource {
	return &schema.Resource{
		Create:   resourceLocationManagementCreate,
		Read:     resourceLocationManagementRead,
		Update:   resourceLocationManagementUpdate,
		Delete:   resourceLocationManagementDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"up_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dn_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"country": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tz": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ports": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_credentials": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CN",
								"IP",
								"UFQDN",
								"XAUTH",
							}, false),
						},
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pre_shared_key": {
							Type:     schema.TypeString,
							Optional: true,
							// Sensitive: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"auth_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssl_scan_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"zapp_ssl_scan_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"xff_forward_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"surrogate_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"idle_time_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"display_time_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MINUTE",
					"HOUR",
					"DAY",
				}, false),
			},
			"surrogate_ip_enforced_for_known_browsers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"surrogate_refresh_time_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"surrogate_refresh_time_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MINUTE",
					"HOUR",
					"DAY",
				}, false),
			},
			"ofw_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ips_control": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"caution_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aup_block_internet_until_accepted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aup_force_ssl_inspection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aup_timeout_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"profile": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"CORPORATE",
					"SERVER",
					"GUESTWIFI",
					"IOT",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceLocationManagementCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandLocationManagement(d)
	log.Printf("[INFO] Creating zia location management\n%+v\n", req)

	resp, err := zClient.locationmanagement.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia location management request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("location_id", resp.ID)

	return resourceLocationManagementRead(d, m)
}

func resourceLocationManagementRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "location_id")
	if !ok {
		return fmt.Errorf("no Traffic Forwarding zia static ip id is set")
	}
	resp, err := zClient.locationmanagement.Get(id)

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing location management %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting location management:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("location_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("parent_id", resp.ParentID)
	_ = d.Set("up_bandwidth", resp.UpBandwidth)
	_ = d.Set("dn_bandwidth", resp.DnBandwidth)
	_ = d.Set("country", resp.Country)
	_ = d.Set("tz", resp.TZ)
	_ = d.Set("ip_addresses", resp.IPAddresses)
	_ = d.Set("ports", resp.Ports)
	_ = d.Set("auth_required", resp.AuthRequired)
	_ = d.Set("ssl_scan_enabled", resp.SSLScanEnabled)
	_ = d.Set("zapp_ssl_scan_enabled", resp.ZappSSLScanEnabled)
	_ = d.Set("xff_forward_enabled", resp.XFFForwardEnabled)
	_ = d.Set("surrogate_ip", resp.SurrogateIP)
	_ = d.Set("idle_time_in_minutes", resp.IdleTimeInMinutes)
	_ = d.Set("display_time_unit", resp.DisplayTimeUnit)
	_ = d.Set("surrogate_ip_enforced_for_known_browsers", resp.SurrogateIPEnforcedForKnownBrowsers)
	_ = d.Set("surrogate_refresh_time_in_minutes", resp.SurrogateRefreshTimeInMinutes)
	_ = d.Set("surrogate_refresh_time_unit", resp.SurrogateRefreshTimeUnit)
	_ = d.Set("ofw_enabled", resp.OFWEnabled)
	_ = d.Set("ips_control", resp.IPSControl)
	_ = d.Set("aup_enabled", resp.AUPEnabled)
	_ = d.Set("caution_enabled", resp.CautionEnabled)
	_ = d.Set("aup_block_internet_until_accepted", resp.AUPBlockInternetUntilAccepted)
	_ = d.Set("aup_force_ssl_inspection", resp.AUPForceSSLInspection)
	_ = d.Set("aup_timeout_in_days", resp.AUPTimeoutInDays)
	_ = d.Set("profile", resp.Profile)
	_ = d.Set("description", resp.Description)

	if err := d.Set("vpn_credentials", flattenLocationVPNCredentials(resp.VPNCredentials)); err != nil {
		return err
	}

	return nil
}

func resourceLocationManagementUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "location_id")
	if !ok {
		log.Printf("[ERROR] location ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating location management ID: %v\n", id)
	req := expandLocationManagement(d)

	if _, _, err := zClient.locationmanagement.Update(id, &req); err != nil {
		return err
	}

	return resourceLocationManagementRead(d, m)
}

func resourceLocationManagementDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "location_id")
	if !ok {
		log.Printf("[ERROR] gre tunnel ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting location management ID: %v\n", (d.Id()))

	if _, err := zClient.locationmanagement.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] location deleted")
	return nil
}

func expandLocationManagement(d *schema.ResourceData) locationmanagement.Locations {
	id, _ := getIntFromResourceData(d, "location_id")
	result := locationmanagement.Locations{
		ID:                                  id,
		Name:                                d.Get("name").(string),
		ParentID:                            d.Get("parent_id").(int),
		UpBandwidth:                         d.Get("up_bandwidth").(int),
		DnBandwidth:                         d.Get("dn_bandwidth").(int),
		Country:                             d.Get("country").(string),
		TZ:                                  d.Get("tz").(string),
		IPAddresses:                         ListToStringSlice(d.Get("ip_addresses").([]interface{})),
		Ports:                               d.Get("ports").(string),
		AuthRequired:                        d.Get("auth_required").(bool),
		SSLScanEnabled:                      d.Get("ssl_scan_enabled").(bool),
		ZappSSLScanEnabled:                  d.Get("zapp_ssl_scan_enabled").(bool),
		XFFForwardEnabled:                   d.Get("xff_forward_enabled").(bool),
		SurrogateIP:                         d.Get("surrogate_ip").(bool),
		IdleTimeInMinutes:                   d.Get("idle_time_in_minutes").(int),
		DisplayTimeUnit:                     d.Get("display_time_unit").(string),
		SurrogateIPEnforcedForKnownBrowsers: d.Get("surrogate_ip_enforced_for_known_browsers").(bool),
		SurrogateRefreshTimeInMinutes:       d.Get("surrogate_refresh_time_in_minutes").(int),
		SurrogateRefreshTimeUnit:            d.Get("surrogate_refresh_time_unit").(string),
		OFWEnabled:                          d.Get("ofw_enabled").(bool),
		IPSControl:                          d.Get("ips_control").(bool),
		AUPEnabled:                          d.Get("aup_enabled").(bool),
		CautionEnabled:                      d.Get("caution_enabled").(bool),
		AUPBlockInternetUntilAccepted:       d.Get("aup_block_internet_until_accepted").(bool),
		AUPForceSSLInspection:               d.Get("aup_force_ssl_inspection").(bool),
		AUPTimeoutInDays:                    d.Get("aup_timeout_in_days").(int),
		Profile:                             d.Get("profile").(string),
		Description:                         d.Get("description").(string),
		VPNCredentials:                      expandLocationManagementVPNCredentials(d),
	}
	vpnCredentials := expandLocationManagementVPNCredentials(d)
	if vpnCredentials != nil {
		result.VPNCredentials = vpnCredentials
	}
	return result
}

func expandLocationManagementVPNCredentials(d *schema.ResourceData) []locationmanagement.VPNCredentials {
	var vpnCredentials []locationmanagement.VPNCredentials
	if vpnCredentialsInterface, ok := d.GetOk("vpn_credentials"); ok {
		vpnCredential := vpnCredentialsInterface.([]interface{})
		vpnCredentials = make([]locationmanagement.VPNCredentials, len(vpnCredential))
		for i, vpn := range vpnCredential {
			vpnItem := vpn.(map[string]interface{})
			vpnCredentials[i] = locationmanagement.VPNCredentials{
				ID:           vpnItem["id"].(int),
				Type:         vpnItem["type"].(string),
				FQDN:         vpnItem["fqdn"].(string),
				PreSharedKey: vpnItem["pre_shared_key"].(string),
				Comments:     vpnItem["comments"].(string),
			}
		}
	}

	return vpnCredentials
}
