package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/terraform-provider-zia/gozscaler/client"
	"github.com/zscaler/terraform-provider-zia/gozscaler/trafficforwarding/gretunnels"
	"github.com/zscaler/terraform-provider-zia/gozscaler/trafficforwarding/virtualipaddresslist"
)

func resourceTrafficForwardingGRETunnel() *schema.Resource {
	return &schema.Resource{
		Create:   resourceTrafficForwardingGRETunnelCreate,
		Read:     resourceTrafficForwardingGRETunnelRead,
		Update:   resourceTrafficForwardingGRETunnelUpdate,
		Delete:   resourceTrafficForwardingGRETunnelDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"tunnel_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the GRE tunnel.",
			},
			"source_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The source IP address of the GRE tunnel. This is typically a static IP address in the organization or SD-WAN.",
				ValidateFunc: validation.IsIPAddress,
			},
			"primary_dest_vip": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The primary destination data center and virtual IP address (VIP) of the GRE tunnel",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Unique identifer of the GRE virtual IP address (VIP)",
						},
						"virtual_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "GRE cluster virtual IP address (VIP)",
						},
						"private_service_edge": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Set to true if the virtual IP address (VIP) is a ZIA Private Service Edge",
						},
						"datacenter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data center information",
						},
					},
				},
			},
			"secondary_dest_vip": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The secondary destination data center and virtual IP address (VIP) of the GRE tunnel",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Unique identifer of the GRE virtual IP address (VIP)",
						},
						"virtual_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "GRE cluster virtual IP address (VIP)",
						},
						"private_service_edge": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Set to true if the virtual IP address (VIP) is a ZIA Private Service Edge",
						},
						"datacenter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data center information",
						},
					},
				},
			},
			"internal_ip_range": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				Description:  "The start of the internal IP address in /29 CIDR range",
				ValidateFunc: validation.IsIPv4Range,
			},

			"last_modification_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_modified_by": {
				Type:     schema.TypeSet,
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
			"within_country": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Restrict the data center virtual IP addresses (VIPs) only to those within the same country as the source IP address",
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "When within_country is enabled, you must set this to the country code.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about this GRE tunnel",
			},
			"ip_unnumbered": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "This is required to support the automated SD-WAN provisioning of GRE tunnels, when set to true gre_tun_ip and gre_tun_id are set to null",
			},
		},
	}
}

func resourceTrafficForwardingGRETunnelCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandGRETunnel(d)
	log.Printf("[INFO] Creating zia gre tunnel\n%+v\n", req)
	err := asssignVipsIfNotSet(d, zClient, &req)
	if err != nil {
		return err
	}
	resp, _, err := zClient.gretunnels.CreateGreTunnels(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia gre tunnel request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("tunnel_id", resp.ID)
	return resourceTrafficForwardingGRETunnelRead(d, m)
}

func asssignVipsIfNotSet(d *schema.ResourceData, zClient *Client, req *gretunnels.GreTunnels) error {
	if (req.PrimaryDestVip == nil || (req.PrimaryDestVip.VirtualIP == "" && req.PrimaryDestVip.ID == 0)) ||
		(req.SecondaryDestVip == nil || (req.SecondaryDestVip.VirtualIP == "" && req.SecondaryDestVip.ID == 0)) {
		// one of the vips not set, pick 2 from the recommandedVips
		countryCode, ok := getStringFromResourceData(d, "country_code")
		var pair []virtualipaddresslist.GREVirtualIPList
		if ok {
			vips, err := zClient.virtualipaddresslist.GetPairZSGREVirtualIPsWithinCountry(req.SourceIP, countryCode)
			if err != nil {
				log.Printf("[ERROR] Got: %v\n", err)
				vips, err = zClient.virtualipaddresslist.GetZSGREVirtualIPList(req.SourceIP, 2)
				if err != nil {
					return err
				}
			}
			pair = *vips
		} else {
			vips, err := zClient.virtualipaddresslist.GetZSGREVirtualIPList(req.SourceIP, 2)
			if err != nil {
				return err
			}
			pair = *vips
		}
		req.PrimaryDestVip = &gretunnels.PrimaryDestVip{ID: pair[0].ID, VirtualIP: pair[0].VirtualIp}
		req.SecondaryDestVip = &gretunnels.SecondaryDestVip{ID: pair[1].ID, VirtualIP: pair[1].VirtualIp}
	}
	return nil
}
func resourceTrafficForwardingGRETunnelRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id, ok := getIntFromResourceData(d, "tunnel_id")
	if !ok {
		return fmt.Errorf("no Traffic Forwarding GRE Tunnel id is set")
	}
	resp, err := zClient.gretunnels.GetGreTunnels(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing gre tunnel %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting gre tunnel:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("tunnel_id", resp.ID)
	_ = d.Set("source_ip", resp.SourceIP)
	_ = d.Set("internal_ip_range", resp.InternalIpRange)
	_ = d.Set("last_modification_time", resp.LastModificationTime)
	_ = d.Set("within_country", resp.WithinCountry)
	_ = d.Set("comment", resp.Comment)
	_ = d.Set("ip_unnumbered", resp.IPUnnumbered)
	if err := d.Set("primary_dest_vip", flattenGrePrimaryDestVipSimple(resp.PrimaryDestVip)); err != nil {
		return err
	}

	if err := d.Set("secondary_dest_vip", flattenGreSecondaryDestVipSimple(resp.SecondaryDestVip)); err != nil {
		return err
	}

	if err := d.Set("last_modified_by", flattenGreLastModifiedBy(resp.LastModifiedBy)); err != nil {
		return err
	}

	return nil
}

func flattenGrePrimaryDestVipSimple(primaryDestVip *gretunnels.PrimaryDestVip) interface{} {
	return []map[string]interface{}{
		{
			"id":                   primaryDestVip.ID,
			"virtual_ip":           primaryDestVip.VirtualIP,
			"private_service_edge": primaryDestVip.PrivateServiceEdge,
			"datacenter":           primaryDestVip.Datacenter,
		},
	}
}
func flattenGreSecondaryDestVipSimple(secondaryDestVip *gretunnels.SecondaryDestVip) interface{} {
	return []map[string]interface{}{
		{
			"id":                   secondaryDestVip.ID,
			"virtual_ip":           secondaryDestVip.VirtualIP,
			"private_service_edge": secondaryDestVip.PrivateServiceEdge,
			"datacenter":           secondaryDestVip.Datacenter,
		},
	}
}

func resourceTrafficForwardingGRETunnelUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "tunnel_id")
	if !ok {
		log.Printf("[ERROR] gre tunnel ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating gre tunnel ID: %v\n", id)
	req := expandGRETunnel(d)

	err := asssignVipsIfNotSet(d, zClient, &req)
	if err != nil {
		return err
	}
	if _, _, err := zClient.gretunnels.UpdateGreTunnels(id, &req); err != nil {
		return err
	}

	return resourceTrafficForwardingGRETunnelRead(d, m)
}

func resourceTrafficForwardingGRETunnelDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id, ok := getIntFromResourceData(d, "tunnel_id")
	if !ok {
		log.Printf("[ERROR] gre tunnel ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting gre tunnel ID: %v\n", id)

	if _, err := zClient.gretunnels.DeleteGreTunnels(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] gre tunnel deleted")
	return nil
}

func expandGRETunnel(d *schema.ResourceData) gretunnels.GreTunnels {
	id, _ := getIntFromResourceData(d, "tunnel_id")
	result := gretunnels.GreTunnels{
		ID:                   id,
		SourceIP:             d.Get("source_ip").(string),
		InternalIpRange:      d.Get("internal_ip_range").(string),
		LastModificationTime: d.Get("last_modification_time").(int),
		WithinCountry:        d.Get("within_country").(bool),
		Comment:              d.Get("comment").(string),
		IPUnnumbered:         d.Get("ip_unnumbered").(bool),
	}
	primaryDestVip := expandPrimaryDestVip(d)
	if primaryDestVip != nil {
		result.PrimaryDestVip = primaryDestVip
	}
	secondaryDestVip := expandSecondaryDestVip(d)
	if secondaryDestVip != nil {
		result.SecondaryDestVip = secondaryDestVip
	}
	lastModifiedBy := expandLastModifiedByTunnel(d)
	if lastModifiedBy != nil {
		result.LastModifiedBy = lastModifiedBy
	}
	return result
}

func expandPrimaryDestVip(d *schema.ResourceData) *gretunnels.PrimaryDestVip {
	vipsObj, ok := d.GetOk("primary_dest_vip")
	if !ok {
		return nil
	}
	vips, ok := vipsObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(vips.List()) > 0 {
		vipObj := vips.List()[0]
		vip, ok := vipObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &gretunnels.PrimaryDestVip{
			ID:                 vip["id"].(int),
			VirtualIP:          vip["virtual_ip"].(string),
			PrivateServiceEdge: vip["private_service_edge"].(bool),
			Datacenter:         vip["datacenter"].(string),
		}
	}
	return nil
}

func expandSecondaryDestVip(d *schema.ResourceData) *gretunnels.SecondaryDestVip {
	vipsObj, ok := d.GetOk("secondary_dest_vip")
	if !ok {
		return nil
	}
	vips, ok := vipsObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(vips.List()) > 0 {
		vipObj := vips.List()[0]
		vip, ok := vipObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &gretunnels.SecondaryDestVip{
			ID:                 vip["id"].(int),
			VirtualIP:          vip["virtual_ip"].(string),
			PrivateServiceEdge: vip["private_service_edge"].(bool),
			Datacenter:         vip["datacenter"].(string),
		}
	}
	return nil
}

func expandLastModifiedByTunnel(d *schema.ResourceData) *gretunnels.LastModifiedBy {
	lastModifiedByObj, ok := d.GetOk("last_modified_by")
	if !ok {
		return nil
	}
	lastModifiedSet, ok := lastModifiedByObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(lastModifiedSet.List()) > 0 {
		lastModifiedObj := lastModifiedSet.List()[0]
		lastModified, ok := lastModifiedObj.(map[string]interface{})
		if !ok {
			return nil
		}
		result := &gretunnels.LastModifiedBy{
			ID: lastModified["id"].(int),
		}
		if lastModified["extensions"] != nil {
			result.Extensions, _ = lastModified["extensions"].(map[string]interface{})
		}
		return result
	}
	return nil
}
