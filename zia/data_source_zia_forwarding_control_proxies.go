package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/proxies"
)

func dataSourceForwardingControlProxies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceForwardingControlProxiesRead,
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
				Description: "Gateway type",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address or the FQDN of the third-party proxy service",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port number on which the third-party proxy service listens to the requests forwarded from Zscaler",
			},
			"insert_xau_header": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether X-Authenticated-User header is added by the proxy. Enable to automatically insert authenticated user ID to the HTTP header, X-Authenticated-User.",
			},
			"base64_encode_xau_header": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether the added X-Authenticated-User header is Base64 encoded. When enabled, the user ID is encoded using the Base64 encoding method.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the ZPA gateway was last modified",
			},
			"cert": {
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

func dataSourceForwardingControlProxiesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *proxies.Proxies
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for proxies id: %d\n", id)
		res, err := proxies.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for proxies name: %s\n", name)
		res, err := proxies.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("type", resp.Type)
		_ = d.Set("address", resp.Address)
		_ = d.Set("port", resp.Port)
		_ = d.Set("insert_xau_header", resp.InsertXauHeader)
		_ = d.Set("base64_encode_xau_header", resp.Base64EncodeXauHeader)

		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("cert", flattenLastModifiedByExternalID(resp.Cert)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("last_modified_by", flattenLastModifiedByExternalID(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any proxy with name '%s'", name))
	}

	return nil
}
