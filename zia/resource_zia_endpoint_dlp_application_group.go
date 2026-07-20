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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_application_groups"
)

// applicationGroupChannel is the only channel supported by endpoint application
// groups. It is fixed here because the UI/API scope this object to
// application-file-access exclusively.
const applicationGroupChannel = "APPLICATION_FILE_ACCESS"

func resourceEndpointDLPApplicationGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointDLPApplicationGroupCreate,
		ReadContext:   resourceEndpointDLPApplicationGroupRead,
		UpdateContext: resourceEndpointDLPApplicationGroupUpdate,
		DeleteContext: resourceEndpointDLPApplicationGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				// The channel is not returned by the list endpoint, so pin it to
				// the only supported value so the Read that follows import works.
				_ = d.Set("channel", applicationGroupChannel)

				key := d.Id()
				if idInt, parseErr := strconv.Atoi(key); parseErr == nil {
					_ = d.Set("group_id", idInt)
					d.SetId(strconv.Itoa(idInt))
					return []*schema.ResourceData{d}, nil
				}

				group, err := endpoint_application_groups.GetByName(ctx, service, key)
				if err != nil {
					return nil, err
				}
				_ = d.Set("group_id", group.GroupID)
				d.SetId(strconv.Itoa(group.GroupID))
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"channel": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     applicationGroupChannel,
				Description: "The channel the endpoint application group belongs to. Only `APPLICATION_FILE_ACCESS` is supported.",
				ValidateFunc: validation.StringInSlice([]string{
					applicationGroupChannel,
				}, false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the endpoint application group.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The description of the endpoint application group.",
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"resources": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The set of endpoint application IDs (the endpoint applications' zappId) that belong to this group.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceEndpointDLPApplicationGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandEndpointDLPApplicationGroup(d)
	log.Printf("[INFO] Creating zia endpoint application group\n%+v\n", req)

	resp, err := endpoint_application_groups.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia endpoint application group. ID: %v\n", resp.ID)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceEndpointDLPApplicationGroupRead(ctx, d, meta)
}

func resourceEndpointDLPApplicationGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no endpoint application group id is set"))
	}

	// There is no get-group-by-id endpoint, so resolve the group from the list.
	list, err := endpoint_application_groups.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	for i := range list {
		if list[i].GroupID != id {
			continue
		}
		found = true

		d.SetId(strconv.Itoa(list[i].GroupID))
		_ = d.Set("group_id", list[i].GroupID)
		_ = d.Set("name", list[i].Name)
		_ = d.Set("description", list[i].Description)

		// The list endpoint expresses membership through endPointApplications,
		// each carrying the application's zappId.
		resources := make([]string, 0, len(list[i].EndPointApplications))
		for _, app := range list[i].EndPointApplications {
			appID := app.ZappID
			if appID == "" && app.ResourceID != 0 {
				appID = strconv.Itoa(app.ResourceID)
			}
			resources = append(resources, appID)
		}
		if err := d.Set("resources", resources); err != nil {
			return diag.FromErr(err)
		}
		break
	}

	if !found {
		log.Printf("[WARN] Removing endpoint application group %s from state because it no longer exists in ZIA", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceEndpointDLPApplicationGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] endpoint application group ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("endpoint application group ID not set"))
	}
	log.Printf("[INFO] Updating zia endpoint application group ID: %v\n", id)

	req := expandEndpointDLPApplicationGroup(d)
	req.ID = id
	if _, err := endpoint_application_groups.Update(ctx, service, id, &req); err != nil {
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

	return resourceEndpointDLPApplicationGroupRead(ctx, d, meta)
}

func resourceEndpointDLPApplicationGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] endpoint application group ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia endpoint application group ID: %v\n", d.Id())

	if _, err := endpoint_application_groups.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia endpoint application group deleted")

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

// expandEndpointDLPApplicationGroup builds the create/update body. Membership is
// carried inline through the resources list (appId + name), and resourceCount
// mirrors its length as the API expects.
func expandEndpointDLPApplicationGroup(d *schema.ResourceData) endpoint_application_groups.EndpointApplicationGroups {
	resources := expandApplicationGroupResources(d)
	return endpoint_application_groups.EndpointApplicationGroups{
		Channel:       d.Get("channel").(string),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		ResourceCount: len(resources),
		Resources:     resources,
	}
}

func expandApplicationGroupResources(d *schema.ResourceData) []endpoint_application_groups.ApplicationGroupResource {
	raw, ok := d.GetOk("resources")
	if !ok {
		return nil
	}
	set, ok := raw.(*schema.Set)
	if !ok {
		return nil
	}
	out := make([]endpoint_application_groups.ApplicationGroupResource, 0, set.Len())
	for _, v := range set.List() {
		appID, ok := v.(string)
		if !ok || appID == "" {
			continue
		}
		out = append(out, endpoint_application_groups.ApplicationGroupResource{
			AppID: appID,
		})
	}
	return out
}
