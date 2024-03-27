package zia

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
)

func resourceTrafficForwardingStaticIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrafficForwardingStaticIPCreate,
		Read:   resourceTrafficForwardingStaticIPRead,
		Update: resourceTrafficForwardingStaticIPUpdate,
		Delete: resourceTrafficForwardingStaticIPDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("static_ip_id", idInt)
				} else {
					resp, err := zClient.staticips.GetByIPAddress(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("static_ip_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"static_ip_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the Static IP.",
			},
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "The static IP address",
			},
			"geo_override": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "If not set, geographic coordinates and city are automatically determined from the IP address. Otherwise, the latitude and longitude coordinates must be provided.",
			},
			"latitude": {
				Type:             schema.TypeFloat,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     ValidateLongitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Required only if the geoOverride attribute is set. Latitude with 7 digit precision after decimal point, ranges between -90 and 90 degrees.",
			},
			"longitude": {
				Type:             schema.TypeFloat,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     ValidateLongitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Required only if the geoOverride attribute is set. Longitude with 7 digit precision after decimal point, ranges between -180 and 180 degrees.",
			},
			"routable_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether a non-RFC 1918 IP address is publicly routable. This attribute is ignored if there is no ZIA Private Service Edge associated to the organization.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about this static IP address",
			},
		},
	}
}

func resourceTrafficForwardingStaticIPCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	if err := checkGeoOverride(d); err != nil {
		return err
	}
	req := expandTrafficForwardingStaticIP(d)
	log.Printf("[INFO] Creating zia static ip\n%+v\n", req)

	resp, _, err := zClient.staticips.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia static ip request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("static_ip_id", resp.ID)

	// Sleep for 2 seconds before triggering the activation
	time.Sleep(2 * time.Second)

	//Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}

	return resourceTrafficForwardingStaticIPRead(d, m)
}

func resourceTrafficForwardingStaticIPRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "static_ip_id")
	if !ok {
		return fmt.Errorf("no Traffic Forwarding zia static ip id is set")
	}
	resp, err := zClient.staticips.Get(id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing static ip %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting static ip:\n%+v\n", resp)
	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("static_ip_id", resp.ID)
	_ = d.Set("ip_address", resp.IpAddress)
	_ = d.Set("geo_override", resp.GeoOverride)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("routable_ip", resp.RoutableIP)
	_ = d.Set("comment", resp.Comment)

	return nil
}

func resourceTrafficForwardingStaticIPUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "static_ip_id")
	if !ok {
		log.Printf("[ERROR] static ip ID not set: %v\n", id)
	}
	if err := checkGeoOverride(d); err != nil {
		return err
	}

	log.Printf("[INFO] Updating static ip ID: %v\n", id)
	req := expandTrafficForwardingStaticIP(d)
	if _, err := zClient.staticips.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := zClient.staticips.Update(id, &req); err != nil {
		return err
	}
	// Sleep for 2 seconds before triggering the activation
	time.Sleep(2 * time.Second)

	//Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}

	return resourceTrafficForwardingStaticIPRead(d, m)
}

func checkGeoOverride(d *schema.ResourceData) error {
	geoOverride, ok := d.GetOk("geo_override")
	if !ok || !(geoOverride.(bool)) {
		return nil
	}
	_, ok = d.GetOk("latitude")
	if !ok {
		return fmt.Errorf("when geo_override is set to true you must specify the latitude & longitude")
	}
	_, ok = d.GetOk("longitude")
	if !ok {
		return fmt.Errorf("when geo_override is set to true you must specify the longitude & longitude")
	}
	return nil
}

func resourceTrafficForwardingStaticIPDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "static_ip_id")
	if !ok {
		log.Printf("[ERROR] static ip ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting static ip ID: %v\n", (d.Id()))

	if _, err := zClient.staticips.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] static ip deleted")

	//Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}

	return nil
}

func expandTrafficForwardingStaticIP(d *schema.ResourceData) staticips.StaticIP {
	id, _ := getIntFromResourceData(d, "static_ip_id")
	result := staticips.StaticIP{
		ID:          id,
		IpAddress:   d.Get("ip_address").(string),
		GeoOverride: d.Get("geo_override").(bool),
		Latitude:    d.Get("latitude").(float64),
		Longitude:   d.Get("longitude").(float64),
		RoutableIP:  d.Get("routable_ip").(bool),
		Comment:     d.Get("comment").(string),
	}
	return result
}
