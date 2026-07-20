package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_custom_apps"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
)

// endpointDLPCustomAppsMutex serializes create/update/delete calls for custom
// endpoint applications. The API rejects concurrent writes to this collection
// with an HTTP 412 (precondition failed), so parallel Terraform operations must
// be funneled through a single writer at a time.
var endpointDLPCustomAppsMutex sync.Mutex

func resourceEndpointDLPCustomApps() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointDLPCustomAppsCreate,
		ReadContext:   resourceEndpointDLPCustomAppsRead,
		UpdateContext: resourceEndpointDLPCustomAppsUpdate,
		DeleteContext: resourceEndpointDLPCustomAppsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("custom_app_id", idInt)
				} else {
					apps, err := endpoint_custom_apps.GetCustomApps(ctx, service, &endpoint_custom_apps.GetCustomAppsFilterOptions{Search: id})
					if err != nil {
						return nil, err
					}
					found := false
					for i := range apps {
						if strings.EqualFold(apps[i].ApplicationName, id) {
							d.SetId(strconv.Itoa(apps[i].ResourceID))
							_ = d.Set("custom_app_id", apps[i].ResourceID)
							found = true
							break
						}
					}
					if !found {
						return nil, fmt.Errorf("couldn't find any endpoint dlp custom app with name %q", id)
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_app_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
				Description:  "The name of the custom endpoint application.",
			},
			"channel": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "APPLICATION_FILE_ACCESS",
				ValidateFunc: validation.StringInSlice([]string{
					"APPLICATION_FILE_ACCESS",
				}, false),
				Description: "The channel of the custom endpoint application. Custom applications only support `APPLICATION_FILE_ACCESS`.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
				Description:  "The description of the custom endpoint application.",
			},
			"app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The application identifier of the custom endpoint application. Assigned by the service.",
			},
			"application": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The application details of the custom endpoint application.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"WINDOWS_OS",
								"MAC_OS",
							}, false),
							Description: "The operating system type of the custom endpoint application. Supported values: `WINDOWS_OS`, `MAC_OS`.",
						},
						"file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The file name of the custom endpoint application.",
						},
						"original_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The original file name of the custom endpoint application.",
						},
						"bundle_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The bundle identifier of the custom endpoint application.",
						},
						"digitally_signed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates whether the custom endpoint application is digitally signed.",
						},
					},
				},
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
		},
	}
}

func resourceEndpointDLPCustomAppsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandEndpointDLPCustomApp(d)
	log.Printf("[INFO] Creating zia endpoint dlp custom app\n%+v\n", req)

	endpointDLPCustomAppsMutex.Lock()
	resp, _, err := endpoint_custom_apps.Create(ctx, service, &req)
	endpointDLPCustomAppsMutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia endpoint dlp custom app. ID: %v\n", resp.ID)
	// The create response carries the full object. State is populated from it
	// directly rather than re-reading from the list endpoint, whose cached
	// response can be stale immediately after a write (the write path targets
	// /customApp while the list is /customApps, so the list cache is not
	// invalidated). Re-reading here intermittently wiped just-created records.
	setEndpointDLPCustomAppFromResource(d, resp)

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

func resourceEndpointDLPCustomAppsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "custom_app_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no endpoint dlp custom app id is set"))
	}

	// The single-object GET endpoint (/customApps/{id}) currently returns HTTP 400,
	// so the record is resolved from the list endpoint and matched on resourceId.
	resp, err := findCustomAppByID(ctx, service, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		log.Printf("[WARN] Removing endpoint dlp custom app %s from state because it no longer exists in ZIA", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Getting endpoint dlp custom app:\n%+v\n", resp)

	d.SetId(strconv.Itoa(resp.ResourceID))
	_ = d.Set("custom_app_id", resp.ResourceID)
	_ = d.Set("name", resp.ApplicationName)
	_ = d.Set("description", resp.Description)
	_ = d.Set("channel", "APPLICATION_FILE_ACCESS")
	_ = d.Set("application_type", resp.ApplicationType)
	_ = d.Set("zapp_id", resp.ZappID)
	// The read endpoint (/endPointApplications/customApps/{id}) does not return
	// channel or appId, so those are preserved from the prior state above.
	if err := d.Set("application", flattenEndpointDLPCustomAppApplication(resp)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEndpointDLPCustomAppsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "custom_app_id")
	if !ok {
		log.Printf("[ERROR] endpoint dlp custom app ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("endpoint dlp custom app ID not set"))
	}
	log.Printf("[INFO] Updating zia endpoint dlp custom app ID: %v\n", id)
	req := expandEndpointDLPCustomApp(d)

	endpointDLPCustomAppsMutex.Lock()
	resp, _, err := endpoint_custom_apps.Update(ctx, service, id, &req)
	endpointDLPCustomAppsMutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	// As with Create, populate state from the update response to avoid a
	// stale-cache read from the list endpoint. Fall back to the request payload
	// if the SDK did not echo the object back.
	if resp == nil {
		resp = &req
	}
	setEndpointDLPCustomAppFromResource(d, resp)

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

func resourceEndpointDLPCustomAppsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "custom_app_id")
	if !ok {
		log.Printf("[ERROR] endpoint dlp custom app ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia endpoint dlp custom app ID: %v\n", d.Id())

	endpointDLPCustomAppsMutex.Lock()
	_, err := endpoint_custom_apps.Delete(ctx, service, id)
	endpointDLPCustomAppsMutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia endpoint dlp custom app deleted")

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

// setEndpointDLPCustomAppFromResource populates state from the write (Create /
// Update) response, which is returned in the endpoint_resource shape. osType is
// normalized back to the short form the user configures.
func setEndpointDLPCustomAppFromResource(d *schema.ResourceData, resp *endpoint_resource.EndpointResource) {
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("custom_app_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("channel", resp.Channel)
	_ = d.Set("app_id", resp.AppID)
	_ = d.Set("application", []interface{}{
		map[string]interface{}{
			"os_type":            flattenCustomAppOsType(resp.Application.OsType),
			"file_name":          resp.Application.FileName,
			"original_file_name": resp.Application.OriginalFileName,
			"bundle_id":          resp.Application.BundleID,
			"digitally_signed":   resp.Application.DigitallySigned,
		},
	})
}

// findCustomAppByID resolves a custom endpoint application from the list
// endpoint and matches it on resourceId. It is used in place of the single
// GET endpoint, which currently returns HTTP 400.
func findCustomAppByID(ctx context.Context, service *zscaler.Service, id int) (*endpoint_custom_apps.EndpointApplications, error) {
	apps, err := endpoint_custom_apps.GetCustomApps(ctx, service, nil)
	if err != nil {
		return nil, err
	}
	for i := range apps {
		if apps[i].ResourceID == id {
			return &apps[i], nil
		}
	}
	return nil, nil
}

func expandEndpointDLPCustomApp(d *schema.ResourceData) endpoint_resource.EndpointResource {
	id, _ := getIntFromResourceData(d, "custom_app_id")
	return endpoint_resource.EndpointResource{
		ID:          id,
		Name:        d.Get("name").(string),
		Channel:     d.Get("channel").(string),
		Description: d.Get("description").(string),
		AppID:       d.Get("app_id").(int),
		Application: expandEndpointDLPCustomAppApplication(d),
	}
}

func expandEndpointDLPCustomAppApplication(d *schema.ResourceData) endpoint_resource.Application {
	raw, ok := d.GetOk("application")
	if !ok {
		return endpoint_resource.Application{}
	}
	list := raw.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return endpoint_resource.Application{}
	}
	m := list[0].(map[string]interface{})
	return endpoint_resource.Application{
		OsType:           expandCustomAppOsType(m["os_type"].(string)),
		FileName:         m["file_name"].(string),
		OriginalFileName: m["original_file_name"].(string),
		BundleID:         m["bundle_id"].(string),
		DigitallySigned:  m["digitally_signed"].(bool),
	}
}

// flattenEndpointDLPCustomAppApplication maps the flat read shape returned by the
// customApps GET endpoint back into the nested application block.
func flattenEndpointDLPCustomAppApplication(resp *endpoint_custom_apps.EndpointApplications) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"os_type":            flattenCustomAppOsType(resp.OsType),
			"file_name":          resp.Filename,
			"original_file_name": resp.OriginalFileName,
			"bundle_id":          resp.Bundle,
			"digitally_signed":   resp.DigitallySigned,
		},
	}
}

// expandCustomAppOsType converts the short osType the user configures
// (WINDOWS_OS / MAC_OS) into the fully-qualified enum the API expects on write.
func expandCustomAppOsType(v string) string {
	switch v {
	case "WINDOWS_OS":
		return "EPTDLP_RESOURCE_APP_OS_WINDOWS_OS"
	case "MAC_OS":
		return "EPTDLP_RESOURCE_APP_OS_MAC_OS"
	default:
		return v
	}
}

// flattenCustomAppOsType maps the API osType back to the short form stored in
// state, tolerating either the short or fully-qualified form on read.
func flattenCustomAppOsType(v string) string {
	switch v {
	case "EPTDLP_RESOURCE_APP_OS_WINDOWS_OS":
		return "WINDOWS_OS"
	case "EPTDLP_RESOURCE_APP_OS_MAC_OS":
		return "MAC_OS"
	default:
		return v
	}
}
