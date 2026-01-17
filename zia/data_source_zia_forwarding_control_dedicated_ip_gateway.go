package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/dedicatedipgateways"
)

func dataSourceForwardingControlDedicatedIPGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceForwardingControlDedicatedIPGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "A unique identifier assigned to the Dedicated IP gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the Dedicated IP gateway",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional details about the Dedicated IP gateway",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the Dedicated IP gateway was created",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the Dedicated IP gateway was last modified",
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
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default Dedicated IP gateway",
			},
		},
	}
}

func dataSourceForwardingControlDedicatedIPGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *dedicatedipgateways.DedicatedIPGateways
	var searchCriteria string

	allGateways, err := dedicatedipgateways.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	id, ok := getIntFromResourceData(d, "id")
	if ok {
		searchCriteria = fmt.Sprintf("id=%d", id)
		for _, gateway := range allGateways {
			if gateway.ID == id {
				resp = &gateway
				break
			}
		}
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		searchCriteria = fmt.Sprintf("name=%s", name)
		for _, gateway := range allGateways {
			if strings.EqualFold(gateway.Name, name) {
				resp = &gateway
				break
			}
		}
	}

	if resp == nil {
		log.Printf("[INFO] Could not find dedicated IP gateway with %s", searchCriteria)
		return diag.FromErr(fmt.Errorf("couldn't find any dedicated IP gateway with %s", searchCriteria))
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("create_time", resp.CreateTime)
	_ = d.Set("last_modified_time", resp.LastModifiedTime)
	_ = d.Set("default", resp.Default)

	if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
