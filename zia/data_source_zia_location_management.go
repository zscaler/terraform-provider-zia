package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/location/locationmanagement"
)

func dataSourceLocationManagement() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationManagementRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"up_bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dn_bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tz": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpn_credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pre_shared_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"managed_by": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"extensions": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"location": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"extensions": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"extranet": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"extranet_ip_pool": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"extranet_dns": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"auth_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"basic_auth_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Basic Authentication at the location",
			},
			"digest_auth_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable Digest Authentication at the location",
			},
			"kerberos_auth_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable Kerberos Authentication at the location",
			},
			"iot_discovery_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable IOT Discovery at the location",
			},
			"ssl_scan_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"zapp_ssl_scan_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"xff_forward_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"surrogate_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"idle_time_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"display_time_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"surrogate_ip_enforced_for_known_browsers": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"surrogate_refresh_time_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"surrogate_refresh_time_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ofw_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ips_control": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"aup_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"caution_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"aup_block_internet_until_accepted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"aup_force_ssl_inspection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"aup_timeout_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"profile": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_extranet_ts_pool": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates that the traffic selector specified in the extranet is the designated default traffic selector",
			},
			"default_extranet_dns": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates that the DNS server configuration used in the extranet is the designated default DNS server",
			},
			"other_sub_location": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv4 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other and it can be renamed, if required.",
			},
			"other6_sub_location": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv6 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other6 and it can be renamed, if required. This field is applicable only if ipv6Enabled is set is true.",
			},
			"sub_loc_scope_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether defining scopes is allowed for this sublocation.",
			},
			"sub_loc_scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines a scope for the sublocation from the available types to segregate workload traffic from a single sublocation to apply different Cloud Connector and ZIA security policies.",
			},
			"sub_loc_scope_values": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sub_loc_acc_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// Add logic to search SubLocation by Name and ID
// See SDK #PR93 https://github.com/zscaler/zscaler-sdk-go/pull/93
func dataSourceLocationManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *locationmanagement.Locations
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for location id: %d\n", id)
		res, err := locationmanagement.GetLocationOrSublocationByID(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	parentName, _ := d.Get("parent_name").(string)
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		if parentName != "" {
			log.Printf("[INFO] Getting data for location name: %s - parent name:%s\n", name, parentName)
			res, err := locationmanagement.GetSubLocationByNames(ctx, service, parentName, name)
			if err != nil {
				return diag.FromErr(err)
			}
			resp = res
		} else {
			log.Printf("[INFO] Getting data for location name: %s\n", name)
			res, err := locationmanagement.GetLocationOrSublocationByName(ctx, service, name)
			if err != nil {
				return diag.FromErr(err)
			}
			resp = res
		}
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
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
		_ = d.Set("default_extranet_ts_pool", resp.DefaultExtranetTsPool)
		_ = d.Set("default_extranet_dns", resp.DefaultExtranetDns)
		_ = d.Set("other_sub_location", resp.OtherSubLocation)
		_ = d.Set("other6_sub_location", resp.Other6SubLocation)
		_ = d.Set("sub_loc_scope", resp.SubLocScope)
		_ = d.Set("sub_loc_scope_values", resp.SubLocScopeValues)
		_ = d.Set("sub_loc_acc_ids", resp.SubLocAccIDs)
		_ = d.Set("sub_loc_scope_enabled", resp.SubLocScopeEnabled)

		if err := d.Set("vpn_credentials", flattenLocationVPNCredentials(resp.VPNCredentials)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("extranet", flattenCustomIDSet(resp.Extranet)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("extranet_ip_pool", flattenCustomIDSet(resp.ExtranetIpPool)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("extranet_dns", flattenCustomIDSet(resp.ExtranetDns)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any location with name '%s'", name))
	}

	return nil
}

func flattenLocationVPNCredentials(vpnCredential []locationmanagement.VPNCredentials) []interface{} {
	vpnCredentials := make([]interface{}, len(vpnCredential))
	for i, vpnCredential := range vpnCredential {
		vpnCredentials[i] = map[string]interface{}{
			"id":             vpnCredential.ID,
			"type":           vpnCredential.Type,
			"fqdn":           vpnCredential.FQDN,
			"pre_shared_key": vpnCredential.PreSharedKey,
			"comments":       vpnCredential.Comments,
			"managed_by":     flattenLocationManagedBy(vpnCredential),
			"location":       flattenVPNCredentialLocation(vpnCredential),
		}
	}

	return vpnCredentials
}

func flattenLocationManagedBy(managedBy locationmanagement.VPNCredentials) []interface{} {
	managed := make([]interface{}, len(managedBy.ManagedBy))
	for i, val := range managedBy.ManagedBy {
		managed[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return managed
}

func flattenVPNCredentialLocation(location locationmanagement.VPNCredentials) []interface{} {
	locations := make([]interface{}, len(location.Location))
	for i, val := range location.Location {
		locations[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return locations
}
