package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/zpa_gateways"
)

func dataSourceForwardingControlZPAGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataForwardingControlZPAGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "A unique identifier assigned to the ZPA gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the ZPA gateway",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional details about the ZPA gateway",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether the ZPA gateway is configured for Zscaler Internet Access (using option ZPA) or Zscaler Cloud Connector (using option ECZPA)",
			},
			"zpa_tenant_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the ZPA tenant where Source IP Anchoring is configured",
			},
			"zpa_server_group": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The ZPA Server Group that is configured for Source IP Anchoring",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"zpa_app_segments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All the Application Segments that are associated with the selected ZPA Server Group for which Source IP Anchoring is enabled",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the ZPA gateway was last modified",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "This is an immutable reference to an entity. which mainly consists of id and name",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataForwardingControlZPAGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *zpa_gateways.ZPAGateways
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for zpa gateway id: %d\n", id)
		// Use GetAll to avoid API bug where Get by ID returns incorrect app segments
		allGateways, err := zpa_gateways.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, gateway := range allGateways {
			if gateway.ID == id {
				resp = &gateway
				break
			}
		}

		if resp == nil {
			return diag.FromErr(fmt.Errorf("no zpa gateway found with id: %d", id))
		}
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for zpa gateway name: %s\n", name)
		// Use GetAll to avoid API bug where GetByName returns incorrect app segments
		allGateways, err := zpa_gateways.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, gateway := range allGateways {
			if strings.EqualFold(gateway.Name, name) {
				resp = &gateway
				break
			}
		}

		if resp == nil {
			return diag.FromErr(fmt.Errorf("no zpa gateway found with name: %s", name))
		}
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("type", resp.Type)
		_ = d.Set("zpa_tenant_id", resp.ZPATenantId)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("zpa_server_group", flattenZPAServerGroup(resp.ZPAServerGroup)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("zpa_app_segments", flattenFWDZPAAppSegments(resp.ZPAAppSegments)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any zpa gateway name '%s' or id '%d'", name, id))
	}

	return nil
}

func flattenZPAServerGroup(group zpa_gateways.ZPAServerGroup) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":        group.Name,
			"external_id": group.ExternalID,
		},
	}
}

func flattenFWDZPAAppSegments(segments []zpa_gateways.ZPAAppSegments) []map[string]interface{} {
	flattenedSegments := make([]map[string]interface{}, len(segments))

	for i, segment := range segments {
		flattenedSegments[i] = map[string]interface{}{
			"name":        segment.Name,
			"external_id": segment.ExternalID,
		}
	}

	return flattenedSegments
}
