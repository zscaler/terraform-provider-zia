package zia

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/locationmanagement"
)

func dataSourceLocationManagement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationManagementRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"up_bandwidth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dn_bandwidth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tz": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_addresses": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ports": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vpn_credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
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
					},
				},
			},
			"auth_required": {
				Type:     schema.TypeBool,
				Computed: true,
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
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
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
			"profile": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLocationManagementRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *locationmanagement.Locations
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting location information %s\n", id)
		res, err := zClient.locationmanagement.GetLocations(id)
		if err != nil {
			return err
		}
		resp = res
	}
	/*name, ok := d.Get("name").(string)
	if id == "" && ok && name != "" {
		log.Printf("[INFO] Getting data for server group name %s\n", name)
		res, err := zClient.locationmanagement.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}*/
	if resp != nil {

		d.SetId(resp.ID)
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
	}
	return nil
}

// Need to flatten managedby, lastModifiedBy and vpnCredentials
