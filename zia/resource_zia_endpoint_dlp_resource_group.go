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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_channel"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_group"
)

func resourceEndpointDLPResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointDLPResourceGroupCreate,
		ReadContext:   resourceEndpointDLPResourceGroupRead,
		UpdateContext: resourceEndpointDLPResourceGroupUpdate,
		DeleteContext: resourceEndpointDLPResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				// The Read path is channel-scoped, so imports must carry the
				// channel. Accept "<CHANNEL>:<id>" or "<CHANNEL>:<name>".
				parts := strings.SplitN(d.Id(), ":", 2)
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("invalid import id %q: expected format \"<CHANNEL>:<id>\" or \"<CHANNEL>:<name>\" (e.g. \"PRINTING:367\")", d.Id())
				}

				channel := endpoint_resource_channel.Channel(parts[0])
				key := parts[1]

				if idInt, parseErr := strconv.Atoi(key); parseErr == nil {
					_ = d.Set("channel", string(channel))
					_ = d.Set("group_id", idInt)
					d.SetId(strconv.Itoa(idInt))
					return []*schema.ResourceData{d}, nil
				}

				list, err := endpoint_resource_group.GetResourceGroupTagsList(ctx, service, channel, nil)
				if err != nil {
					return nil, err
				}
				for i := range list {
					if strings.EqualFold(list[i].Name, key) {
						_ = d.Set("channel", string(channel))
						_ = d.Set("group_id", list[i].ID)
						d.SetId(strconv.Itoa(list[i].ID))
						return []*schema.ResourceData{d}, nil
					}
				}
				return nil, fmt.Errorf("couldn't find any dlp endpoint resource group with name %q in channel %q", key, channel)
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
				Required:    true,
				Description: "The channel the DLP endpoint resource group belongs to. A resource group cannot be moved between channels after creation.",
				ValidateFunc: validation.StringInSlice([]string{
					"PRINTING",
					"REMOVABLE_DRIVE_TRANSFER",
					"NETWORK_DRIVE_TRANSFER",
				}, false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the DLP endpoint resource group (tag).",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The description of the DLP endpoint resource group.",
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"resources": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The set of DLP endpoint resource IDs associated with this group.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceEndpointDLPResourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := endpoint_resource_group.DlpEndpointResourceGroups{
		Channel:     d.Get("channel").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	log.Printf("[INFO] Creating zia dlp endpoint resource group\n%+v\n", req)

	resp, _, err := endpoint_resource_group.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia dlp endpoint resource group. ID: %v\n", resp.ID)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)

	// Resource membership is managed through a dedicated association endpoint,
	// not the group create body. Associate the initial resources (if any) here.
	if adds := setToIntSlice(d, "resources"); len(adds) > 0 {
		assoc := &endpoint_resource_group.EndpointDlpGroupToResourceAssociation{ResourcesToBeAdded: adds}
		if _, err := endpoint_resource_group.UpdateDlpResourcesByTag(ctx, service, resp.ID, assoc); err != nil {
			return diag.FromErr(fmt.Errorf("error associating resources with dlp endpoint resource group %d: %w", resp.ID, err))
		}
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceEndpointDLPResourceGroupRead(ctx, d, meta)
}

func resourceEndpointDLPResourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no dlp endpoint resource group id is set"))
	}

	channel := endpoint_resource_channel.Channel(d.Get("channel").(string))
	if channel == "" {
		return diag.FromErr(fmt.Errorf("no channel is set for dlp endpoint resource group id %d", id))
	}

	// There is no get-group-by-id endpoint, so resolve the group metadata from
	// the channel's group list.
	list, err := endpoint_resource_group.GetResourceGroupTagsList(ctx, service, channel, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var found *endpoint_resource_group.DlpEndpointResourceGroups
	for i := range list {
		if list[i].ID == id {
			found = &list[i]
			break
		}
	}
	if found == nil {
		log.Printf("[WARN] Removing dlp endpoint resource group %s from state because it no longer exists in ZIA", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(strconv.Itoa(found.ID))
	_ = d.Set("group_id", found.ID)
	_ = d.Set("name", found.Name)
	_ = d.Set("channel", found.Channel)
	_ = d.Set("description", found.Description)

	// Member resources come from a separate association endpoint.
	members, err := endpoint_resource_group.GetDlpResourceByTag(ctx, service, id)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceIDs := make([]int, 0, len(members))
	for _, m := range members {
		resourceIDs = append(resourceIDs, m.ID)
	}
	if err := d.Set("resources", resourceIDs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEndpointDLPResourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] dlp endpoint resource group ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("dlp endpoint resource group ID not set"))
	}
	log.Printf("[INFO] Updating zia dlp endpoint resource group ID: %v\n", id)

	// Group metadata (name/description/channel) is updated on the group
	// endpoint. Membership is handled separately below so we never re-send the
	// full resource list on a metadata change.
	if d.HasChanges("name", "description", "channel") {
		req := endpoint_resource_group.DlpEndpointResourceGroups{
			ID:          id,
			Channel:     d.Get("channel").(string),
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}
		if _, _, err := endpoint_resource_group.Update(ctx, service, id, &req); err != nil {
			return diag.FromErr(err)
		}
	}

	// Apply only the resource membership delta via the association endpoint.
	if d.HasChange("resources") {
		oldRaw, newRaw := d.GetChange("resources")
		oldSet := oldRaw.(*schema.Set)
		newSet := newRaw.(*schema.Set)

		toAdd := intSetDifference(newSet, oldSet)
		toDelete := intSetDifference(oldSet, newSet)

		if len(toAdd) > 0 || len(toDelete) > 0 {
			assoc := &endpoint_resource_group.EndpointDlpGroupToResourceAssociation{
				ResourcesToBeAdded:   toAdd,
				ResourcesToBeDeleted: toDelete,
			}
			if _, err := endpoint_resource_group.UpdateDlpResourcesByTag(ctx, service, id, assoc); err != nil {
				return diag.FromErr(fmt.Errorf("error updating resource membership for dlp endpoint resource group %d: %w", id, err))
			}
		}
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceEndpointDLPResourceGroupRead(ctx, d, meta)
}

func resourceEndpointDLPResourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] dlp endpoint resource group ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia dlp endpoint resource group ID: %v\n", d.Id())

	if _, err := endpoint_resource_group.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia dlp endpoint resource group deleted")

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

// setToIntSlice extracts a TypeSet of ints from the resource data as a []int.
func setToIntSlice(d *schema.ResourceData, key string) []int {
	raw, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	set, ok := raw.(*schema.Set)
	if !ok {
		return nil
	}
	out := make([]int, 0, set.Len())
	for _, v := range set.List() {
		out = append(out, v.(int))
	}
	return out
}

// intSetDifference returns the ints present in a but not in b.
func intSetDifference(a, b *schema.Set) []int {
	diff := a.Difference(b).List()
	out := make([]int, 0, len(diff))
	for _, v := range diff {
		out = append(out, v.(int))
	}
	return out
}
