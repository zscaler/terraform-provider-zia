package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_applications"
)

func dataSourceEndpointDLPApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEndpointDLPApplicationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier (resourceId) of the endpoint application. Used to look up a single application by ID.",
			},
			"application_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the endpoint application. Used to look up a single application by name.",
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search string applied server-side against application names to narrow the result set before matching. Useful because the application list can be extensive.",
			},
			"os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The operating system type of the endpoint application (e.g. WINDOWS_OS, MAC_OS). Also applied server-side as a filter when set.",
			},
			"application_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The application type of the endpoint application (e.g. WELLKNOWN, CUSTOM, DISCOVERED). Also applied server-side as a filter when set.",
			},
			"zapp_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Zscaler Client Connector application identifier.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the endpoint application.",
			},
			"file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The file name of the endpoint application.",
			},
			"original_file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The original file name of the endpoint application.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bundle identifier of the endpoint application.",
			},
			"digitally_signed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the endpoint application is digitally signed.",
			},
			"mod_uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The modification unique identifier of the endpoint application.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the endpoint application was last modified.",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the endpoint application is deleted.",
			},
			"version":  dataSourceEndpointApplicationVersionSchema(),
			"versions": dataSourceEndpointApplicationVersionSchema(),
		},
	}
}

func dataSourceEndpointApplicationVersionSchema() *schema.Schema {
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

func dataSourceEndpointDLPApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, _ := getIntFromResourceData(d, "id")
	name, _ := d.Get("application_name").(string)

	if id == 0 && name == "" {
		return diag.FromErr(fmt.Errorf("one of 'id' or 'application_name' must be set to look up an endpoint application"))
	}

	// The API supports native server-side filtering by name (search), os type,
	// and application type, plus pagination. Push those filters down so the
	// extensive application list is narrowed before the local id/name match,
	// rather than fetching and filtering the entire catalogue client-side.
	opts := &endpoint_applications.GetApplicationCountFilterOptions{}
	if v, ok := d.Get("os_type").(string); ok && v != "" {
		opts.OsType = v
	}
	if v, ok := d.Get("application_type").(string); ok && v != "" {
		opts.ApplicationType = v
	}
	if s, ok := d.Get("search").(string); ok && s != "" {
		opts.Search = s
	} else if name != "" {
		opts.Search = name
	}

	apps, err := endpoint_applications.GetAllEndpointApplications(ctx, service, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp *common.EndPointApplications
	for i := range apps {
		if id != 0 {
			if apps[i].ResourceID == id {
				resp = &apps[i]
				break
			}
			continue
		}
		if strings.EqualFold(apps[i].ApplicationName, name) {
			resp = &apps[i]
			break
		}
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any endpoint application with name '%s' or id '%d'", name, id))
	}

	log.Printf("[INFO] Getting data for endpoint application id: %d name: %s\n", resp.ResourceID, resp.ApplicationName)

	d.SetId(strconv.Itoa(resp.ResourceID))
	_ = d.Set("id", resp.ResourceID)
	_ = d.Set("application_name", resp.ApplicationName)
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

	if err := d.Set("version", flattenEndpointApplicationVersion(resp.Version)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("versions", flattenEndpointApplicationVersions(resp.Versions)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenEndpointApplicationVersion(v common.Version) []interface{} {
	if v == (common.Version{}) {
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

func flattenEndpointApplicationVersions(v common.Versions) []interface{} {
	if v == (common.Versions{}) {
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
