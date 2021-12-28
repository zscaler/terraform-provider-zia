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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location Name.",
			},
			"parent_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"ip_addresses"},
				Description:  "Parent Location ID. If this ID does not exist or is 0, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: SUB",
			},
			"up_bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 99999999),
				Description:  "Upload bandwidth in bytes. The value 0 implies no Bandwidth Control enforcement.",
			},
			"dn_bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 99999999),
				Description:  "Download bandwidth in bytes. The value 0 implies no Bandwidth Control enforcement.",
			},
			"country": getLocationManagementCountries(),
			"tz":      getLocationManagementTimeZones(),
			"ip_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						validation.IsIPv4Range,
						validation.IsIPv4Address,
					),
				},
				Description: "For locations: IP addresses of the egress points that are provisioned in the Zscaler Cloud. Each entry is a single IP address (e.g., 238.10.33.9).",
			},
			"ports": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP ports that are associated with the location.",
			},
			"vpn_credentials": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"IP",
								"UFQDN",
							}, false),
						},
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_address": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"pre_shared_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ssl_scan_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable SSL Inspection. Set to true in order to apply your SSL Inspection policy to HTTPS traffic in the location and inspect HTTPS transactions for data leakage, malicious content, and viruses.",
				Deprecated:  "This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future. Enable SSL Inspection. Set to true in order to apply your SSL Inspection policy to HTTPS traffic in the location and inspect HTTPS transactions for data leakage, malicious content, and viruses. To learn more, see Deploying SSL Inspection",
			},
			"zapp_ssl_scan_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Zscaler App SSL Setting. When set to true, the Zscaler App SSL Scan Setting will take effect, irrespective of the SSL policy that is configured for the location.",
				Deprecated:  "This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future. Enable Zscaler App SSL Setting. When set to true, the Zscaler App SSL Scan Setting will take effect, irrespective of the SSL policy that is configured for the location. To learn more, see Configuring SSL Inspection for Zscaler App",
			},
			"xff_forward_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable XFF Forwarding. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.",
			},
			"surrogate_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Surrogate IP. When set to true, users are mapped to internal device IP addresses.",
			},
			"auth_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enforce Authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.",
			},
			"idle_time_in_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Idle Time to Disassociation. The user mapping idle time (in minutes) is required if a Surrogate IP is enabled.",
			},
			"display_time_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "MINUTE",
				Description: "Display Time Unit. The time unit to display for IP Surrogate idle time to disassociation.",
				ValidateFunc: validation.StringInSlice([]string{
					"MINUTE",
					"HOUR",
					"DAY",
				}, false),
			},
			"surrogate_ip_enforced_for_known_browsers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enforce Surrogate IP for Known Browsers. When set to true, IP Surrogate is enforced for all known browsers.",
			},
			"surrogate_refresh_time_in_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Refresh Time for re-validation of Surrogacy. The surrogate refresh time (in minutes) to re-validate the IP surrogates.",
			},
			"surrogate_refresh_time_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display Refresh Time Unit. The time unit to display for refresh time for re-validation of surrogacy.",
				ValidateFunc: validation.StringInSlice([]string{
					"MINUTE",
					"HOUR",
					"DAY",
				}, false),
			},
			"ofw_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Firewall. When set to true, Firewall is enabled for the location.",
			},
			"ips_control": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable IPS Control. When set to true, IPS Control is enabled for the location if Firewall is enabled.",
			},
			"aup_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable AUP. When set to true, AUP is enabled for the location.",
			},
			"caution_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Caution. When set to true, a caution notifcation is enabled for the location.",
			},
			"aup_block_internet_until_accepted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "For First Time AUP Behavior, Block Internet Access. When set, all internet access (including non-HTTP traffic) is disabled until the user accepts the AUP.",
			},
			"aup_force_ssl_inspection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "For First Time AUP Behavior, Force SSL Inspection. When set, Zscaler will force SSL Inspection in order to enforce AUP for HTTPS traffic.",
			},
			"aup_timeout_in_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Custom AUP Frequency. Refresh time (in days) to re-validate the AUP.",
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Profile tag that specifies the location traffic type. If not specified, this tag defaults to `Unassigned`.",
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"CORPORATE",
					"SERVER",
					"GUESTWIFI",
					"IOT",
				}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Additional notes or information regarding the location or sub-location. The description cannot exceed 1024 characters.",
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
		},
	}
}

func resourceLocationManagementCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	if parentIDInt, ok := d.GetOk("parent_id"); ok && parentIDInt.(int) != 0 {
		ipInter, ipSet := d.GetOk("ip_addresses")
		if !ipSet || len(removeEmpty(ListToStringSlice(ipInter.([]interface{})))) == 0 {
			return fmt.Errorf("when the location is a sub-location ip_addresses must not be empty: %v", d.Get("name"))
		}
	}

	req := expandLocationManagement(d)
	log.Printf("[INFO] Creating zia location management\n%+v\n", req)
	if err := checkSurrogateIPDependencies(req); err != nil {
		return err
	}
	if err := checkVPNCredentials(req); err != nil {
		return err
	}
	resp, err := zClient.locationmanagement.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia location management request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("location_id", resp.ID)

	return resourceLocationManagementRead(d, m)
}
func checkVPNCredentials(locations locationmanagement.Locations) error {
	for _, vpn := range locations.VPNCredentials {
		if vpn.Type == "IP" && vpn.IPAddress == "" {
			return fmt.Errorf("ip_address is required when VPN credential is of type IP")
		}
	}
	return nil
}
func checkSurrogateIPDependencies(loc locationmanagement.Locations) error {
	if loc.SurrogateIP && loc.IdleTimeInMinutes == 0 {
		return fmt.Errorf("surrogate IP requires setting of an idle timeout")
	}
	if loc.SurrogateIP && !loc.AuthRequired {
		return fmt.Errorf("authentication required must be enabled, when enabling surrogate IP")
	}
	if loc.SurrogateIPEnforcedForKnownBrowsers && !loc.SurrogateIP {
		return fmt.Errorf("surrogate IP must be enabled, when enforcing surrogate IP for known browsers")
	}
	if loc.SurrogateIPEnforcedForKnownBrowsers && loc.SurrogateRefreshTimeInMinutes == 0 && loc.SurrogateRefreshTimeUnit == "" {
		return fmt.Errorf("enforcing surrogate IP for known browsers requires setting of refresh timeout")
	}
	return nil
}

func resourceLocationManagementRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "location_id")
	if !ok {
		return fmt.Errorf("no location management id is set")
	}
	resp, err := zClient.locationmanagement.GetLocation(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
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

	if err := d.Set("vpn_credentials", flattenLocationVPNCredentialsSimple(resp.VPNCredentials)); err != nil {
		return err
	}

	return nil
}

func flattenLocationVPNCredentialsSimple(vpnCredential []locationmanagement.VPNCredentials) []interface{} {
	vpnCredentials := make([]interface{}, len(vpnCredential))
	for i, vpnCredential := range vpnCredential {
		vpnCredentials[i] = map[string]interface{}{
			"id":             vpnCredential.ID,
			"type":           vpnCredential.Type,
			"fqdn":           vpnCredential.FQDN,
			"ip_address":     vpnCredential.IPAddress,
			"pre_shared_key": vpnCredential.PreSharedKey,
			"comments":       vpnCredential.Comments,
		}
	}

	return vpnCredentials
}

func resourceLocationManagementUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "location_id")
	if !ok {
		log.Printf("[ERROR] location ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating location management ID: %v\n", id)
	req := expandLocationManagement(d)
	if err := checkSurrogateIPDependencies(req); err != nil {
		return err
	}

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

func removeEmpty(list []string) []string {
	result := []string{}
	for _, i := range list {
		if i != "" {
			result = append(result, i)
		}
	}
	return result
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
		IPAddresses:                         removeEmpty(ListToStringSlice(d.Get("ip_addresses").([]interface{}))),
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
				IPAddress:    vpnItem["ip_address"].(string),
				PreSharedKey: vpnItem["pre_shared_key"].(string),
				Comments:     vpnItem["comments"].(string),
			}
		}
	}
	return vpnCredentials
}
