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
				ValidateFunc:     ValidateLatitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Latitude with 7 digit precision after decimal point, ranges between -90 and 90 degrees. If not provided, the API will automatically determine it from the IP address.",
			},
			"longitude": {
				Type:             schema.TypeFloat,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     ValidateLongitude,
				DiffSuppressFunc: DiffSuppressFuncCoordinate,
				Description:      "Longitude with 7 digit precision after decimal point, ranges between -180 and 180 degrees. If not provided, the API will automatically determine it from the IP address.",
			},
			"routable_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether a non-RFC 1918 IP address is publicly routable. This attribute is ignored if there is no ZIA Private Service Edge associated to the organization.",
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				// Computed:    true,
				Description: "Additional information about this static IP address",
			},
		},
	}
}

func resourceTrafficForwardingStaticIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Handle geo_override with auto-determined coordinates
	// This may create the IP if coordinates need to be determined
	if err := autoPopulateCoordinates(ctx, d, zClient); err != nil {
		return diag.FromErr(err)
	}

	// Check if IP was already created by autoPopulateCoordinates
	if d.Id() != "" {
		log.Printf("[INFO] Static IP already created during coordinate auto-population. ID: %s", d.Id())
		// Skip to read and activation
	} else {
		// Normal create flow
		req := expandTrafficForwardingStaticIP(d)
		log.Printf("[INFO] Creating zia static ip\n%+v\n", req)

		resp, _, err := staticips.Create(ctx, service, &req)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] Created zia static ip request. ID: %v\n", resp)
		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("static_ip_id", resp.ID)
	}

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
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

	// Use GetAll() instead of individual Get() to reduce API calls during terraform refresh.
	// With hundreds of static IPs, individual GET calls per resource drain rate limits quickly.
	allIPs, err := staticips.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp *staticips.StaticIP
	for i := range allIPs {
		if allIPs[i].ID == id {
			resp = &allIPs[i]
			break
		}
	}

	// Fallback to direct Get if not found in GetAll (may be due to stale SDK cache)
	if resp == nil {
		log.Printf("[WARN] Static IP %d not found in GetAll response, falling back to direct Get", id)
		rule, err := staticips.Get(ctx, service, id)
		if err != nil {
			if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
				log.Printf("[WARN] Removing static ip %s from state because it no longer exists in ZIA", d.Id())
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
		resp = rule
	}

	log.Printf("[INFO] Getting static ip:\n%+v\n", resp)
	log.Printf("[DEBUG] API returned coordinates - Latitude: %.10f, Longitude: %.10f", resp.Latitude, resp.Longitude)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("static_ip_id", resp.ID)
	_ = d.Set("ip_address", resp.IpAddress)
	_ = d.Set("geo_override", resp.GeoOverride)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("routable_ip", resp.RoutableIP)
	_ = d.Set("comment", resp.Comment)

	log.Printf("[DEBUG] State updated with coordinates - Latitude: %.10f, Longitude: %.10f",
		d.Get("latitude").(float64), d.Get("longitude").(float64))

	return nil
}

func resourceTrafficForwardingStaticIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "static_ip_id")
	if !ok {
		log.Printf("[ERROR] static ip ID not set: %v\n", id)
	}

	// Get current resource state to populate coordinates if needed
	currentIP, err := staticips.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// If geo_override=true but coordinates not provided, use existing coordinates from API
	geoOverride, _ := d.Get("geo_override").(bool)
	if geoOverride {
		_, hasLat := d.GetOk("latitude")
		_, hasLon := d.GetOk("longitude")

		if !hasLat || !hasLon {
			// Use existing coordinates from the current resource
			log.Printf("[INFO] geo_override=true but coordinates not provided, using existing: Latitude=%.10f, Longitude=%.10f",
				currentIP.Latitude, currentIP.Longitude)
			_ = d.Set("latitude", currentIP.Latitude)
			_ = d.Set("longitude", currentIP.Longitude)
		}
	}

	log.Printf("[INFO] Updating static ip ID: %v\n", id)
	req := expandTrafficForwardingStaticIP(d)
	if _, _, err := staticips.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}
	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceTrafficForwardingStaticIPRead(ctx, d, meta)
}

func checkGeoOverride(d *schema.ResourceData) error {
	// Validation removed - coordinates will be auto-determined if needed
	return nil
}

// autoPopulateCoordinates auto-determines coordinates when geo_override=true but lat/long not provided
func autoPopulateCoordinates(ctx context.Context, d *schema.ResourceData, service *Client) error {
	svc := service.Service
	geoOverride, _ := d.Get("geo_override").(bool)

	// If geo_override is false, API will auto-determine coordinates - nothing to do
	if !geoOverride {
		return nil
	}

	// Check if user provided coordinates
	lat, hasLat := d.GetOk("latitude")
	lon, hasLon := d.GetOk("longitude")

	// If both provided, use user values
	if hasLat && hasLon {
		log.Printf("[DEBUG] User provided coordinates: Latitude=%.10f, Longitude=%.10f",
			lat.(float64), lon.(float64))
		return nil
	}

	// geo_override=true but coordinates NOT provided
	// Solution: Temporarily use geo_override=false to let API determine coordinates
	ipAddress := d.Get("ip_address").(string)
	log.Printf("[INFO] geo_override=true but coordinates not provided, auto-determining for IP: %s", ipAddress)

	// Check if this IP already exists (could have coordinates)
	existingIP, err := staticips.GetByIPAddress(ctx, svc, ipAddress)
	if err == nil {
		// IP exists - reuse its coordinates
		log.Printf("[INFO] Found existing IP with coordinates: Latitude=%.10f, Longitude=%.10f",
			existingIP.Latitude, existingIP.Longitude)
		_ = d.Set("latitude", existingIP.Latitude)
		_ = d.Set("longitude", existingIP.Longitude)
		return nil
	}

	// IP doesn't exist yet - create with geo_override=false first to get coordinates
	log.Printf("[DEBUG] Creating temporary IP to determine coordinates")
	tempReq := staticips.StaticIP{
		IpAddress:   ipAddress,
		GeoOverride: false, // API will determine coordinates
		RoutableIP:  d.Get("routable_ip").(bool),
		Comment:     d.Get("comment").(string),
	}

	tempResp, _, err := staticips.Create(ctx, svc, &tempReq)
	if err != nil {
		return fmt.Errorf("failed to create static IP to determine coordinates: %w", err)
	}

	// Get the full details including coordinates
	ipWithCoords, err := staticips.Get(ctx, svc, tempResp.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch static IP coordinates: %w", err)
	}

	// Store coordinates for use in subsequent operations
	_ = d.Set("latitude", ipWithCoords.Latitude)
	_ = d.Set("longitude", ipWithCoords.Longitude)
	log.Printf("[INFO] Auto-determined coordinates: Latitude=%.10f, Longitude=%.10f",
		ipWithCoords.Latitude, ipWithCoords.Longitude)

	// Now update with geo_override=true using the determined coordinates
	updateReq := staticips.StaticIP{
		ID:          tempResp.ID,
		IpAddress:   ipAddress,
		GeoOverride: true, // Now we can set it to true
		Latitude:    ipWithCoords.Latitude,
		Longitude:   ipWithCoords.Longitude,
		RoutableIP:  d.Get("routable_ip").(bool),
		Comment:     d.Get("comment").(string),
	}

	_, _, err = staticips.Update(ctx, svc, tempResp.ID, &updateReq)
	if err != nil {
		return fmt.Errorf("failed to update static IP with geo_override: %w", err)
	}

	// Set the ID so the rest of Create knows it's already created
	d.SetId(strconv.Itoa(tempResp.ID))
	_ = d.Set("static_ip_id", tempResp.ID)

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

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
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
