package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/proxies"
)

func dataSourceDedicatedIPProxy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedIPProxyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "A unique identifier assigned to the Dedicated IP Gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the Dedicated IP Gateway",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the Dedicated IP Gateway was created",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the Dedicated IP Gateway was last modified",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default Dedicated IP Gateway",
			},
		},
	}
}

func dataSourceDedicatedIPProxyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	allGateways, err := proxies.GetDedicatedIPGWLite(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp *proxies.DedicatedIPGateways
	var searchCriteria string

	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting dedicated IP proxy by id: %d\n", id)
		searchCriteria = fmt.Sprintf("id=%d", id)
		for i := range allGateways {
			if allGateways[i].Id == id {
				resp = &allGateways[i]
				break
			}
		}
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting dedicated IP proxy by name: %s\n", name)
		searchCriteria = fmt.Sprintf("name=%s", name)
		for i := range allGateways {
			if strings.EqualFold(allGateways[i].Name, name) {
				resp = &allGateways[i]
				break
			}
		}
	}

	if resp == nil {
		if searchCriteria != "" {
			return diag.FromErr(fmt.Errorf("couldn't find any dedicated IP gateway with %s", searchCriteria))
		}
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be specified"))
	}

	d.SetId(fmt.Sprintf("%d", resp.Id))
	_ = d.Set("id", resp.Id)
	_ = d.Set("name", resp.Name)
	_ = d.Set("create_time", resp.CreateTime)
	_ = d.Set("last_modified_time", resp.LastModifiedTime)
	_ = d.Set("default", resp.Default)

	return nil
}
