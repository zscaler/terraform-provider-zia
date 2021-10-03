package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/staticips"
)

func resourceTrafficForwardingStaticIP() *schema.Resource {
	return &schema.Resource{
		Create:   resourceTrafficForwardingStaticIPCreate,
		Read:     resourceTrafficForwardingStaticIPRead,
		Update:   resourceTrafficForwardingStaticIPUpdate,
		Delete:   resourceTrafficForwardingStaticIPDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"geo_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"latitude": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"longitude": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"routable_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTrafficForwardingStaticIPCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandTrafficForwardingStaticIP(d)
	log.Printf("[INFO] Creating zia static ip\n%+v\n", req)

	resp, _, err := zClient.staticips.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia static ip request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))

	return resourceTrafficForwardingStaticIPRead(d, m)
}

func resourceTrafficForwardingStaticIPRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "id")
	if !ok {
		return fmt.Errorf("no Traffic Forwarding zia static ip id is set")
	}
	resp, err := zClient.staticips.Get(id)

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing static ip %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting static ip:\n%+v\n", resp)
	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("ip_address", resp.IpAddress)
	_ = d.Set("geo_override", resp.GeoOverride)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("routable_ip", resp.RoutableIP)

	return nil
}

func resourceTrafficForwardingStaticIPUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating static ip ID: %v\n", id)
	req := expandTrafficForwardingStaticIP(d)

	if _, err := zClient.staticips.Update(id, &req); err != nil {
		return err
	}

	return resourceTrafficForwardingStaticIPRead(d, m)
}

func resourceTrafficForwardingStaticIPDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	// Need to pass the ID (int) of the resource for deletion
	log.Printf("[INFO] Deleting static ip ID: %v\n", (d.Id()))

	if err := zClient.staticips.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] static ip deleted")
	return nil
}

func expandTrafficForwardingStaticIP(d *schema.ResourceData) staticips.StaticIP {
	return staticips.StaticIP{
		IpAddress:   d.Get("ip_address").(string),
		GeoOverride: d.Get("geo_override").(bool),
		Latitude:    d.Get("latitude").(int),
		Longitude:   d.Get("longitude").(int),
		RoutableIP:  d.Get("routable_ip").(bool),
		Comment:     d.Get("comment").(string),
	}
}
