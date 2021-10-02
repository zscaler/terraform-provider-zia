package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/gretunnels"
)

func resourceTrafficForwardingGRETunnel() *schema.Resource {
	return &schema.Resource{
		Create:   resourceTrafficForwardingGRETunnelCreate,
		Read:     resourceTrafficForwardingGRETunnelRead,
		Update:   resourceTrafficForwardingGRETunnelUpdate,
		Delete:   resourceTrafficForwardingGRETunnelDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"source_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The source IP address of the GRE tunnel.",
			},
			"primary_dest_vip": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Role of the admin. This is not required for an auditor.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Unique identifer of the GRE virtual IP address (VIP)",
						},
						"virtual_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GRE cluster virtual IP address (VIP)",
						},
					},
				},
			},
			"secondary_dest_vip": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Role of the admin. This is not required for an auditor.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Unique identifer of the GRE virtual IP address (VIP)",
						},
						"virtual_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GRE cluster virtual IP address (VIP)",
						},
					},
				},
			},
			"internal_ip_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The start of the internal IP address in /29 CIDR range",
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
				Type:     schema.TypeBool,
				Optional: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_unnumbered": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "This is required to support the automated SD-WAN provisioning of GRE tunnels, when set to true gre_tun_ip and gre_tun_id are set to null",
			},
		},
	}
}

func resourceTrafficForwardingGRETunnelCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandGRETunnel(d)
	log.Printf("[INFO] Creating zia gre tunnel\n%+v\n", req)

	resp, _, err := zClient.gretunnels.CreateGreTunnels(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia gre tunnel request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))

	return resourceTrafficForwardingGRETunnelRead(d, m)
}

func resourceTrafficForwardingGRETunnelRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id, ok := getIntFromResourceData(d, "id")
	if !ok {
		return fmt.Errorf("no Traffic Forwarding GRE Tunnel id is set")
	}
	resp, err := zClient.gretunnels.GetGreTunnels(id)

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing gre tunnel %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting gre tunnel:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("source_ip", resp.SourceIP)
	_ = d.Set("internal_ip_range", resp.InternalIpRange)
	_ = d.Set("last_modification_time", resp.LastModificationTime)
	_ = d.Set("within_country", resp.WithinCountry)
	_ = d.Set("comment", resp.Comment)
	_ = d.Set("ip_unnumbered", resp.IPUnnumbered)
	if err := d.Set("primary_dest_vip", flattenGrePrimaryDestVip(resp.PrimaryDestVip)); err != nil {
		return err
	}

	if err := d.Set("secondary_dest_vip", flattenGreSecondaryDestVip(resp.SecondaryDestVip)); err != nil {
		return err
	}

	if err := d.Set("managed_by", flattenGreManagedBy(resp.ManagedBy)); err != nil {
		return err
	}

	if err := d.Set("last_modified_by", flattenGreLastModifiedBy(resp.LastModifiedBy)); err != nil {
		return err
	}

	return nil
}

func resourceTrafficForwardingGRETunnelUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating gre tunnel ID: %v\n", id)
	req := expandGRETunnel(d)

	if _, _, err := zClient.gretunnels.UpdateGreTunnels(id, &req); err != nil {
		return err
	}

	return resourceTrafficForwardingGRETunnelRead(d, m)
}

func resourceTrafficForwardingGRETunnelDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting gre tunnel ID: %v\n", (d.Id()))

	if _, err := zClient.gretunnels.DeleteGreTunnels(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] gre tunnel deleted")
	return nil
}

func expandGRETunnel(d *schema.ResourceData) gretunnels.GreTunnels {
	return gretunnels.GreTunnels{
		SourceIP:             d.Get("source_ip").(string),
		InternalIpRange:      d.Get("internal_ip_range").(string),
		LastModificationTime: d.Get("last_modification_time").(int),
		WithinCountry:        d.Get("within_country").(bool),
		Comment:              d.Get("comment").(string),
		IPUnnumbered:         d.Get("ip_unnumbered").(bool),
		PrimaryDestVip:       expandPrimaryDestVip(d),
		SecondaryDestVip:     expandSecondaryDestVip(d),
		LastModifiedBy:       expandLastModifiedBy(d),
	}
}

func expandPrimaryDestVip(d *schema.ResourceData) gretunnels.PrimaryDestVip {
	primaryDestVip := gretunnels.PrimaryDestVip{
		ID:        d.Get("id").(int),
		VirtualIP: d.Get("virtual_ip").(string),
	}

	return primaryDestVip
}

func expandSecondaryDestVip(d *schema.ResourceData) gretunnels.SecondaryDestVip {
	secondaryDestVip := gretunnels.SecondaryDestVip{
		ID:        d.Get("id").(int),
		VirtualIP: d.Get("virtual_ip").(string),
	}

	return secondaryDestVip
}

func expandLastModifiedBy(d *schema.ResourceData) gretunnels.LastModifiedBy {
	lastModifiedBy := gretunnels.LastModifiedBy{
		ID:         d.Get("id").(int),
		Extensions: d.Get("extensions").(map[string]interface{}),
	}

	return lastModifiedBy
}
