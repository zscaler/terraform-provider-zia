package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/proxy_gateways"
)

func dataSourceForwardingControlProxyGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceForwardingControlProxyGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "A unique identifier assigned to the Proxy gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the Proxy gateway",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional details about the Proxy gateway",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returned Values: PROXYCHAIN, ZIA, ECSELF",
			},
			"fail_closed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether fail close is enabled to drop the traffic or disabled to allow the traffic when both primary and secondary proxies defined in this gateway are unreachable.",
			},
			"primary_proxy": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The primary proxy for the gateway. This field is not applicable to the Lite API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"secondary_proxy": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The seconday proxy for the gateway. This field is not applicable to the Lite API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
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

func dataSourceForwardingControlProxyGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *proxy_gateways.ProxyGateways
	var searchCriteria string

	// Check if searching by ID
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting proxy gateway by id: %d\n", id)
		searchCriteria = fmt.Sprintf("id=%d", id)

		// Get all proxy gateways and find the one with matching ID
		allProxyGateways, err := proxy_gateways.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, tw := range allProxyGateways {
			if tw.ID == id {
				resp = &tw
				break
			}
		}
	}

	// Check if searching by name (only if ID search didn't find anything)
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting proxy gateway by name: %s\n", name)
		searchCriteria = fmt.Sprintf("name=%s", name)

		res, err := proxy_gateways.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("fail_closed", resp.FailClosed)
		_ = d.Set("type", resp.Type)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("primary_proxy", flattenIDNameExternalSet(resp.PrimaryProxy)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("secondary_proxy", flattenIDNameExternalSet(resp.SecondaryProxy)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any proxy gateway with %s", searchCriteria))
	}

	return nil
}
