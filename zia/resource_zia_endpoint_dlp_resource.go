package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_channel"
)

func resourceEndpointDLPResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointDLPResourceCreate,
		ReadContext:   resourceEndpointDLPResourceRead,
		UpdateContext: resourceEndpointDLPResourceUpdate,
		DeleteContext: resourceEndpointDLPResourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				// The Read path requires the channel, so imports must carry it.
				// Accept "<CHANNEL>:<id>" or "<CHANNEL>:<name>".
				parts := strings.SplitN(d.Id(), ":", 2)
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("invalid import id %q: expected format \"<CHANNEL>:<id>\" or \"<CHANNEL>:<name>\" (e.g. \"PRINTING:12345\")", d.Id())
				}

				channel := endpoint_resource_channel.Channel(parts[0])
				key := parts[1]

				if idInt, parseErr := strconv.Atoi(key); parseErr == nil {
					_ = d.Set("channel", string(channel))
					_ = d.Set("resource_id", idInt)
					d.SetId(strconv.Itoa(idInt))
					return []*schema.ResourceData{d}, nil
				}

				list, err := endpoint_resource_channel.GetChannelList(ctx, service, channel, &endpoint_resource_channel.GetChannelFilterOptions{Name: key})
				if err != nil {
					return nil, err
				}
				for i := range list {
					if strings.EqualFold(list[i].Name, key) {
						_ = d.Set("channel", string(channel))
						_ = d.Set("resource_id", list[i].ID)
						d.SetId(strconv.Itoa(list[i].ID))
						return []*schema.ResourceData{d}, nil
					}
				}
				return nil, fmt.Errorf("couldn't find any dlp endpoint resource with name %q in channel %q", key, channel)
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"channel": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The channel of the DLP endpoint resource.",
				ValidateFunc: validation.StringInSlice([]string{
					"PRINTING",
					"REMOVABLE_DRIVE_TRANSFER",
					"NETWORK_DRIVE_TRANSFER",
				}, false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the DLP endpoint resource.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The description of the DLP endpoint resource.",
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"network_drive_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The network drive type of the DLP endpoint resource.",
			},
			"server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The server name of the DLP endpoint resource.",
			},
			"app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The application identifier of the DLP endpoint resource.",
			},
			"network_drives": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of network drives associated with the DLP endpoint resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"printer": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The printer details of the DLP endpoint resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unc": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"removable_storage": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The removable storage details of the DLP endpoint resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vendor_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"serial_number": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceEndpointDLPResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandEndpointDLPResource(d)
	log.Printf("[INFO] Creating zia dlp endpoint resource\n%+v\n", req)

	resp, _, err := endpoint_resource.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia dlp endpoint resource. ID: %v\n", resp.ID)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("resource_id", resp.ID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceEndpointDLPResourceRead(ctx, d, meta)
}

func resourceEndpointDLPResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "resource_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no dlp endpoint resource id is set"))
	}

	channel := endpoint_resource_channel.Channel(d.Get("channel").(string))
	if channel == "" {
		return diag.FromErr(fmt.Errorf("no channel is set for dlp endpoint resource id %d", id))
	}

	resp, err := endpoint_resource_channel.GetChannel(ctx, service, channel, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing dlp endpoint resource %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting dlp endpoint resource:\n%+v\n", resp)

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("resource_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("channel", resp.Channel)
	_ = d.Set("network_drive_type", resp.NetworkDriveType)
	_ = d.Set("description", resp.Description)
	_ = d.Set("server_name", resp.ServerName)
	_ = d.Set("app_id", resp.AppID)

	if err := d.Set("network_drives", flattenDLPEndpointNetworkDrives(resp.NetworkDrives)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("printer", flattenDLPEndpointPrinter(resp.Printer)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("removable_storage", flattenDLPEndpointRemovableStorage(resp.RemovableStorage)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEndpointDLPResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "resource_id")
	if !ok {
		log.Printf("[ERROR] dlp endpoint resource ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("dlp endpoint resource ID not set"))
	}
	log.Printf("[INFO] Updating zia dlp endpoint resource ID: %v\n", id)
	req := expandEndpointDLPResource(d)

	if _, _, err := endpoint_resource.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceEndpointDLPResourceRead(ctx, d, meta)
}

func resourceEndpointDLPResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "resource_id")
	if !ok {
		log.Printf("[ERROR] dlp endpoint resource ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia dlp endpoint resource ID: %v\n", d.Id())

	if _, err := endpoint_resource.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia dlp endpoint resource deleted")

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandEndpointDLPResource(d *schema.ResourceData) endpoint_resource.EndpointResource {
	id, _ := getIntFromResourceData(d, "resource_id")
	result := endpoint_resource.EndpointResource{
		ID:               id,
		Name:             d.Get("name").(string),
		Channel:          d.Get("channel").(string),
		NetworkDriveType: d.Get("network_drive_type").(string),
		Description:      d.Get("description").(string),
		ServerName:       d.Get("server_name").(string),
		AppID:            d.Get("app_id").(int),
		NetworkDrives:    expandEndpointDLPNetworkDrives(d),
		Printer:          expandEndpointDLPPrinter(d),
		RemovableStorage: expandEndpointDLPRemovableStorage(d),
	}
	return result
}

func expandEndpointDLPNetworkDrives(d *schema.ResourceData) []endpoint_resource.NetworkDrive {
	raw, ok := d.GetOk("network_drives")
	if !ok {
		return nil
	}
	list := raw.([]interface{})
	drives := make([]endpoint_resource.NetworkDrive, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		drives = append(drives, endpoint_resource.NetworkDrive{
			NetworkPath: m["network_path"].(string),
		})
	}
	return drives
}

func expandEndpointDLPPrinter(d *schema.ResourceData) endpoint_resource.Printer {
	raw, ok := d.GetOk("printer")
	if !ok {
		return endpoint_resource.Printer{}
	}
	list := raw.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return endpoint_resource.Printer{}
	}
	m := list[0].(map[string]interface{})
	return endpoint_resource.Printer{
		Unc:       m["unc"].(string),
		IpAddress: m["ip_address"].(string),
		Domain:    m["domain"].(string),
	}
}

func expandEndpointDLPRemovableStorage(d *schema.ResourceData) endpoint_resource.RemovableStorage {
	raw, ok := d.GetOk("removable_storage")
	if !ok {
		return endpoint_resource.RemovableStorage{}
	}
	list := raw.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return endpoint_resource.RemovableStorage{}
	}
	m := list[0].(map[string]interface{})
	return endpoint_resource.RemovableStorage{
		VendorId:     m["vendor_id"].(string),
		ProductId:    m["product_id"].(string),
		SerialNumber: m["serial_number"].(string),
	}
}
