package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_clusters"
)

func dataSourceVZENCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVZENClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "System-generated Virtual Service Edge cluster ID",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of the Virtual Service Edge cluster",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the status of the Virtual Service Edge cluster. The status is set to ENABLED by default",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Virtual Service Edge cluster type",
			},
			"ip_sec_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Virtual Service Edge cluster IP address",
			},
			"subnet_mask": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Virtual Service Edge cluster subnet mask",
			},
			"default_gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the default gateway to the internet",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the rule lable was last modified. This is a read-only field. Ignored by PUT and DELETE requests.",
			},
			"virtual_zen_nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Virtual Service Edge instances you want to include in the cluster.",
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
						"external_id": {
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

func dataSourceVZENClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *vzen_clusters.VZENClusters
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for rule label id: %d\n", id)
		res, err := vzen_clusters.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for rule label name: %s\n", name)
		res, err := vzen_clusters.GetClusterByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("status", resp.Status)
		_ = d.Set("type", resp.Type)
		_ = d.Set("ip_address", resp.IpAddress)
		_ = d.Set("subnet_mask", resp.SubnetMask)
		_ = d.Set("default_gateway", resp.DefaultGateway)
		_ = d.Set("ip_sec_enabled", resp.IpSecEnabled)

		if err := d.Set("virtual_zen_nodes", flattenCommonIDNameExternalID(resp.VirtualZenNodes)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any rule label name '%s' or id '%d'", name, id))
	}

	return nil
}

func flattenCommonIDNameExternalID(list []common.IDNameExternalID) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		r := map[string]interface{}{
			"id":   val.ID,
			"name": val.Name,
			// "external_id": val.ExternalID,
			// "extensions":  val.Extensions,
		}
		flattenedList[i] = r
	}
	return flattenedList
}
