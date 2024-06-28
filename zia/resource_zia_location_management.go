package zia

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/location/locationmanagement"
)

func resourceLocationManagement() *schema.Resource {
	return &schema.Resource{
		Create: resourceLocationManagementCreate,
		Read:   resourceLocationManagementRead,
		Update: resourceLocationManagementUpdate,
		Delete: resourceLocationManagementDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)
				service := zClient.locationmanagement

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("location_id", idInt)
				} else {
					resp, err := locationmanagement.GetLocationOrSublocationByName(service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("location_id", resp.ID)
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
			"location_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Location Name.",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Additional notes or information regarding the location or sub-location. The description cannot exceed 1024 characters.",
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Parent Location ID. If this ID does not exist or is 0, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: SUB",
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						validation.IsIPv4Address,
						validation.IsIPv4Range,
						validation.IsIPAddress,
						validation.IsCIDRNetwork(0, 32),
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
			},
			"zapp_ssl_scan_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Zscaler App SSL Setting. When set to true, the Zscaler App SSL Scan Setting will take effect, irrespective of the SSL policy that is configured for the location.",
			},
			"xff_forward_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable XFF Forwarding. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.",
			},
			"other_sublocation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv4 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other and it can be renamed, if required.",
			},
			"other6_sublocation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv6 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other6 and it can be renamed, if required. This field is applicable only if ipv6Enabled is set is true.",
			},
			"surrogate_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Surrogate IP. When set to true, users are mapped to internal device IP addresses.",
			},
			"basic_auth_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Basic Authentication at the location",
			},
			"digest_auth_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Digest Authentication at the location",
			},
			"kerberos_auth_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Kerberos Authentication at the location",
			},
			"iot_discovery_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable IOT Discovery at the location",
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
			"ipv6_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, IPv6 is enabled for the location and IPv6 traffic from the location can be forwarded to the Zscaler service to enforce security policies.",
			},
			"ipv6_dns_64prefix": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "(Optional) Name-ID pair of the NAT64 prefix configured as the DNS64 prefix for the location. If specified, the DNS64 prefix is used for the IP addresses that reside in this location. If not specified, a prefix is selected from the set of supported prefixes. ",
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
					"WORKLOAD",
				}, false),
			},
		},
	}
}

func resourceLocationManagementCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.locationmanagement

	if parentIDInt, ok := d.GetOk("parent_id"); ok && parentIDInt.(int) != 0 {
		ipAddresses := d.Get("ip_addresses").(*schema.Set)
		if len(removeEmpty(ListToStringSlice(ipAddresses.List()))) == 0 {
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
	resp, err := locationmanagement.Create(service, &req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia location management request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("location_id", resp.ID)

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
	service := zClient.locationmanagement

	id, ok := getIntFromResourceData(d, "location_id")
	if !ok {
		return fmt.Errorf("no location management id is set")
	}
	resp, err := locationmanagement.GetLocationOrSublocationByID(service, id)
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
	_ = d.Set("description", resp.Description)
	_ = d.Set("parent_id", resp.ParentID)
	_ = d.Set("up_bandwidth", resp.UpBandwidth)
	_ = d.Set("dn_bandwidth", resp.DnBandwidth)
	_ = d.Set("country", resp.Country)
	_ = d.Set("tz", resp.TZ)
	_ = d.Set("ip_addresses", resp.IPAddresses)
	_ = d.Set("ports", resp.Ports)
	_ = d.Set("auth_required", resp.AuthRequired)
	_ = d.Set("basic_auth_enabled", resp.BasicAuthEnabled)
	_ = d.Set("digest_auth_enabled", resp.DigestAuthEnabled)
	_ = d.Set("kerberos_auth_enabled", resp.KerberosAuth)
	_ = d.Set("iot_discovery_enabled", resp.IOTDiscoveryEnabled)
	_ = d.Set("other_sublocation", resp.OtherSubLocation)
	_ = d.Set("other6_sublocation", resp.Other6SubLocation)
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
	_ = d.Set("ipv6_enabled", resp.IPv6Enabled)
	_ = d.Set("ipv6_dns_64prefix", resp.IPv6Dns64Prefix)

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
	service := zClient.locationmanagement

	id, ok := getIntFromResourceData(d, "location_id")
	if !ok {
		log.Printf("[ERROR] location ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating location management ID: %v\n", id)
	req := expandLocationManagement(d)
	if err := checkSurrogateIPDependencies(req); err != nil {
		return err
	}
	if _, err := locationmanagement.GetLocationOrSublocationByID(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, _, err := locationmanagement.Update(service, id, &req); err != nil {
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

	return resourceLocationManagementRead(d, m)
}

func resourceLocationManagementDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.locationmanagement

	id, ok := getIntFromResourceData(d, "location_id")
	if !ok {
		log.Printf("[ERROR] gre tunnel ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting location management ID: %v\n", (d.Id()))
	err := DetachRuleIDNameExtensions(
		zClient,
		id,
		"Users",
		func(r *filteringrules.FirewallFilteringRules) []common.IDNameExtensions {
			return r.Users
		},
		func(r *filteringrules.FirewallFilteringRules, ids []common.IDNameExtensions) {
			r.Users = ids
		},
	)
	if err != nil {
		return err
	}
	if _, err := locationmanagement.Delete(service, id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] location deleted")

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
		Description:                         d.Get("description").(string),
		ParentID:                            d.Get("parent_id").(int),
		UpBandwidth:                         d.Get("up_bandwidth").(int),
		DnBandwidth:                         d.Get("dn_bandwidth").(int),
		Country:                             d.Get("country").(string),
		TZ:                                  d.Get("tz").(string),
		IPAddresses:                         SetToStringList(d, "ip_addresses"), // removeEmpty(ListToStringSlice(d.Get("ip_addresses").([]interface{}))),
		Ports:                               d.Get("ports").(string),
		AuthRequired:                        d.Get("auth_required").(bool),
		BasicAuthEnabled:                    d.Get("basic_auth_enabled").(bool),
		DigestAuthEnabled:                   d.Get("digest_auth_enabled").(bool),
		KerberosAuth:                        d.Get("kerberos_auth_enabled").(bool),
		IOTDiscoveryEnabled:                 d.Get("iot_discovery_enabled").(bool),
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
		IPv6Enabled:                         d.Get("ipv6_enabled").(bool),
		IPv6Dns64Prefix:                     d.Get("ipv6_dns_64prefix").(bool),
		OtherSubLocation:                    d.Get("other_sublocation").(bool),
		Other6SubLocation:                   d.Get("other6_sublocation").(bool),

		AUPTimeoutInDays: d.Get("aup_timeout_in_days").(int),
		Profile:          d.Get("profile").(string),
		VPNCredentials:   expandLocationManagementVPNCredentials(d),
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
