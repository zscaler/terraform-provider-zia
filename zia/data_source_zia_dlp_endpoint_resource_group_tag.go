package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_channel"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_group"
)

func dataSourceDLPEndpointResourceGroupTag() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDLPEndpointResourceGroupTagRead,
		Schema: map[string]*schema.Schema{
			"channel": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DLP endpoint resource channel whose tags are returned.",
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
				Description: "The unique identifier of a specific tag to return. Optional; when set, the tag is looked up directly by ID (returns a name-ID pair without a description).",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a specific tag to return. Optional; when set, the result is filtered to this tag.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of tags, in name-ID pairs, defined for the specified channel.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies the tag.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the tag.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the tag.",
						},
					},
				},
			},
			"tags_by_name": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A convenience map of tag name to tag ID, allowing a single tag ID to be looked up directly by name (e.g. tags_by_name[\"PrinterTag01\"]).",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceDLPEndpointResourceGroupTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	channel := endpoint_resource_channel.Channel(d.Get("channel").(string))

	name, _ := d.Get("name").(string)
	filterID, _ := getIntFromResourceData(d, "id")

	tagsOut := []interface{}{}
	byName := map[string]interface{}{}
	var lastID int
	matched := 0

	if filterID != 0 {
		// Look a single tag up directly by ID. This endpoint returns name-ID
		// pairs only (no description).
		log.Printf("[INFO] Getting dlp endpoint resource group tag by id: %d\n", filterID)
		tags, err := endpoint_resource_group.GetResourceGroupTag(ctx, service, filterID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error retrieving tag id '%d': %s", filterID, err))
		}
		for _, t := range tags {
			if name != "" && !strings.EqualFold(t.Name, name) {
				continue
			}
			tagsOut = append(tagsOut, map[string]interface{}{
				"id":          t.ID,
				"name":        t.Name,
				"description": "",
			})
			if t.Name != "" {
				byName[t.Name] = t.ID
			}
			lastID = t.ID
			matched++
		}
	} else {
		// Retrieve the channel's tag catalog and optionally filter by name.
		opts := &endpoint_resource_group.GetResourceTagsListFilterOptions{}
		if name != "" {
			search := true
			opts.Name = name
			opts.SearchResources = &search
		}

		log.Printf("[INFO] Getting dlp endpoint resource tags for channel: %s\n", channel)
		list, err := endpoint_resource_group.GetResourceGroupTagsList(ctx, service, channel, opts)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error retrieving tags for channel '%s': %s", channel, err))
		}
		for _, t := range list {
			if name != "" && !strings.EqualFold(t.Name, name) {
				continue
			}
			tagsOut = append(tagsOut, map[string]interface{}{
				"id":          t.ID,
				"name":        t.Name,
				"description": t.Description,
			})
			if t.Name != "" {
				byName[t.Name] = t.ID
			}
			lastID = t.ID
			matched++
		}
	}

	if (filterID != 0 || name != "") && matched == 0 {
		return diag.FromErr(fmt.Errorf("couldn't find any dlp endpoint resource tag with name '%s' or id '%d' in channel '%s'", name, filterID, channel))
	}

	if matched == 1 {
		d.SetId(strconv.Itoa(lastID))
	} else {
		d.SetId(string(channel))
	}
	_ = d.Set("channel", string(channel))

	if err := d.Set("tags", tagsOut); err != nil {
		return diag.FromErr(fmt.Errorf("error setting tags: %s", err))
	}

	if err := d.Set("tags_by_name", byName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting tags_by_name: %s", err))
	}

	return nil
}
