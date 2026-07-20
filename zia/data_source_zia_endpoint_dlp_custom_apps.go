package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_custom_apps"
)

func dataSourceEndpointDLPCustomApps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEndpointDLPCustomAppsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier of the custom endpoint application.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the custom endpoint application.",
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search string used to match against custom application names when looking up by name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the custom endpoint application.",
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The operating system type of the custom endpoint application.",
			},
			"file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The file name of the custom endpoint application.",
			},
			"original_file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The original file name of the custom endpoint application.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bundle identifier of the custom endpoint application.",
			},
			"digitally_signed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the custom endpoint application is digitally signed.",
			},
			"application_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The application type of the endpoint application.",
			},
			"zapp_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Zscaler Client Connector application identifier.",
			},
			"mod_uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The modification unique identifier of the custom endpoint application.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the custom endpoint application was last modified.",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the custom endpoint application is deleted.",
			},
			"version":  dataSourceEndpointCustomAppVersionSchema(),
			"versions": dataSourceEndpointCustomAppVersionSchema(),
		},
	}
}

func dataSourceEndpointCustomAppVersionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"version": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"z_ver_id_md32": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"threat_type": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"threat_level": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"bundle_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"code_signing_certificate_status": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"threat_level_updated": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceEndpointDLPCustomAppsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *endpoint_custom_apps.EndpointApplications

	id, ok := getIntFromResourceData(d, "id")
	if ok && id != 0 {
		// The single-object GET endpoint (/customApps/{id}) currently returns
		// HTTP 400, so the record is resolved from the list and matched on id.
		log.Printf("[INFO] Getting data for endpoint dlp custom app id: %d\n", id)
		res, err := findCustomAppByID(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for endpoint dlp custom app: %s\n", name)
		search := name
		if s, ok := d.Get("search").(string); ok && s != "" {
			search = s
		}
		apps, err := endpoint_custom_apps.GetCustomApps(ctx, service, &endpoint_custom_apps.GetCustomAppsFilterOptions{Search: search})
		if err != nil {
			return diag.FromErr(err)
		}
		for i := range apps {
			if strings.EqualFold(apps[i].ApplicationName, name) {
				resp = &apps[i]
				break
			}
		}
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any endpoint dlp custom app with name '%s' or id '%d'", name, id))
	}

	d.SetId(strconv.Itoa(resp.ResourceID))
	_ = d.Set("name", resp.ApplicationName)
	_ = d.Set("description", resp.Description)
	_ = d.Set("os_type", resp.OsType)
	_ = d.Set("file_name", resp.Filename)
	_ = d.Set("original_file_name", resp.OriginalFileName)
	_ = d.Set("bundle_id", resp.Bundle)
	_ = d.Set("digitally_signed", resp.DigitallySigned)
	_ = d.Set("application_type", resp.ApplicationType)
	_ = d.Set("zapp_id", resp.ZappID)
	_ = d.Set("mod_uid", resp.ModUId)
	_ = d.Set("last_modified_time", resp.LastModifiedTime)
	_ = d.Set("deleted", resp.Deleted)

	if err := d.Set("version", flattenEndpointCustomAppVersion(resp.Version)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("versions", flattenEndpointCustomAppVersions(resp.Versions)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenEndpointCustomAppVersion(v endpoint_custom_apps.Version) []interface{} {
	if v == (endpoint_custom_apps.Version{}) {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"version":                         v.Version,
			"z_ver_id_md32":                   v.ZverIDMD32,
			"threat_type":                     v.ThreatType,
			"threat_level":                    v.ThreatLevel,
			"bundle_id":                       v.Bundle,
			"code_signing_certificate_status": v.CodeSigningCertificateStatus,
			"threat_level_updated":            v.ThreatLevelUpdated,
		},
	}
}

func flattenEndpointCustomAppVersions(versions []endpoint_custom_apps.Versions) []interface{} {
	if len(versions) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, len(versions))
	for i, v := range versions {
		out[i] = map[string]interface{}{
			"version":                         v.Version,
			"z_ver_id_md32":                   v.ZverIDMD32,
			"threat_type":                     v.ThreatType,
			"threat_level":                    v.ThreatLevel,
			"bundle_id":                       v.Bundle,
			"code_signing_certificate_status": v.CodeSigningCertificateStatus,
			"threat_level_updated":            v.ThreatLevelUpdated,
		}
	}
	return out
}
