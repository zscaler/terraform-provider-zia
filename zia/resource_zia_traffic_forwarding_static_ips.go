package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
)

func resourceTrafficForwardingStaticIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrafficForwardingStaticIPCreate,
		ReadContext:   resourceTrafficForwardingStaticIPRead,
		UpdateContext: resourceTrafficForwardingStaticIPUpdate,
		DeleteContext: resourceTrafficForwardingStaticIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("static_ip_id", idInt)
				} else {
					resp, err := staticips.GetByIPAddress(ctx, service, id)
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
				Computed:    true,
				Description: "Additional information about this static IP address",
			},
		},
	}
}

func resourceTrafficForwardingStaticIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	if err := checkGeoOverride(d); err != nil {
		return diag.FromErr(err)
	}
	req := expandTrafficForwardingStaticIP(d)
	log.Printf("[INFO] Creating zia static ip\n%+v\n", req)

	resp, _, err := staticips.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia static ip request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("static_ip_id", resp.ID)

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceTrafficForwardingStaticIPRead(ctx, d, meta)
}

func resourceTrafficForwardingStaticIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "static_ip_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no Traffic Forwarding zia static ip id is set"))
	}
	resp, err := staticips.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing static ip %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
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

func resourceTrafficForwardingStaticIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "static_ip_id")
	if !ok {
		log.Printf("[ERROR] static ip ID not set: %v\n", id)
	}
	if err := checkGeoOverride(d); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating static ip ID: %v\n", id)
	req := expandTrafficForwardingStaticIP(d)
	if _, err := staticips.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := staticips.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceTrafficForwardingStaticIPRead(ctx, d, meta)
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

func resourceTrafficForwardingStaticIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "static_ip_id")
	if !ok {
		log.Printf("[ERROR] static ip ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting static ip ID: %v\n", (d.Id()))

	if _, err := staticips.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] static ip deleted")

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
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
