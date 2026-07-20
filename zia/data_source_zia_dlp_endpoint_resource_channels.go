package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_channel"
)

func dataSourceDLPEndpointResourceChannels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDLPEndpointResourceChannelsRead,
		Schema: map[string]*schema.Schema{
			"channel": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DLP endpoint resource channel to query.",
				ValidateFunc: validation.StringInSlice([]string{
					"PRINTING",
					"REMOVABLE_DRIVE_TRANSFER",
					"NETWORK_DRIVE_TRANSFER",
					"PERSONAL_CLOUD_STORAGE",
				}, false),
			},
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier of the DLP endpoint resource.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the DLP endpoint resource.",
			},
			"is_predefined": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the DLP endpoint resource is predefined.",
			},
			"network_drive_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network drive type of the DLP endpoint resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the DLP endpoint resource.",
			},
			"server_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server name of the DLP endpoint resource.",
			},
			"app_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The application identifier of the DLP endpoint resource.",
			},
			"network_drives": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of network drives associated with the DLP endpoint resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"printer": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The printer details of the DLP endpoint resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"removable_storage": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The removable storage details of the DLP endpoint resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vendor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"application": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The application details of the DLP endpoint resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"original_file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bundle_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"digitally_signed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDLPEndpointResourceChannelsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	channel := endpoint_resource_channel.Channel(d.Get("channel").(string))

	var resp *endpoint_resource.EndpointResource

	id, ok := getIntFromResourceData(d, "id")
	if ok && id != 0 {
		log.Printf("[INFO] Getting data for dlp endpoint resource id: %d in channel: %s\n", id, channel)
		res, err := endpoint_resource_channel.GetChannel(ctx, service, channel, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for dlp endpoint resource: %s in channel: %s\n", name, channel)
		list, err := endpoint_resource_channel.GetChannelList(ctx, service, channel, &endpoint_resource_channel.GetChannelFilterOptions{Name: name})
		if err != nil {
			return diag.FromErr(err)
		}
		for i := range list {
			if strings.EqualFold(list[i].Name, name) {
				resp = &list[i]
				break
			}
		}
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any dlp endpoint resource with name '%s' or id '%d' in channel '%s'", name, id, channel))
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("channel", resp.Channel)
	_ = d.Set("is_predefined", resp.IsPredefined)
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
	if err := d.Set("application", flattenDLPEndpointApplication(resp.Application)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenDLPEndpointNetworkDrives(drives []endpoint_resource.NetworkDrive) []interface{} {
	if len(drives) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, len(drives))
	for i, drive := range drives {
		out[i] = map[string]interface{}{
			"network_path": drive.NetworkPath,
		}
	}
	return out
}

func flattenDLPEndpointPrinter(printer endpoint_resource.Printer) []interface{} {
	if printer == (endpoint_resource.Printer{}) {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"unc":        printer.Unc,
			"ip_address": printer.IpAddress,
			"domain":     printer.Domain,
		},
	}
}

func flattenDLPEndpointRemovableStorage(storage endpoint_resource.RemovableStorage) []interface{} {
	if storage == (endpoint_resource.RemovableStorage{}) {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"vendor_id":     storage.VendorId,
			"product_id":    storage.ProductId,
			"serial_number": storage.SerialNumber,
		},
	}
}

func flattenDLPEndpointApplication(app endpoint_resource.Application) []interface{} {
	if app == (endpoint_resource.Application{}) {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"os_type":            app.OsType,
			"file_name":          app.FileName,
			"original_file_name": app.OriginalFileName,
			"bundle_id":          app.BundleID,
			"digitally_signed":   app.DigitallySigned,
		},
	}
}
