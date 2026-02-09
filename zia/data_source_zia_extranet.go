package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/extranet"
)

func dataSourceExtranet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExtranetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the extranet.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the extranet.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the extranet.",
			},
			"extranet_dns_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the DNS servers specified for the extranet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID generated for the DNS server configuration.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the DNS server.",
						},
						"primary_dns_server": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the primary DNS server.",
						},
						"secondary_dns_server": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the secondary DNS server.",
						},
						"use_as_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the DNS servers specified in the extranet are the designated default servers.",
						},
					},
				},
			},
			"extranet_ip_pool_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the traffic selectors (IP pools) specified for the extranet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID generated for the IP pool configuration.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the IP pool.",
						},
						"ip_start": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The starting IP address of the pool.",
						},
						"ip_end": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ending IP address of the pool.",
						},
						"use_as_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this IP pool is the designated default.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp when the extranet was created.",
			},
			"modified_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp when the extranet was last modified.",
			},
		},
	}
}

// flattenExtranetDNSList converts []extranet.ExtranetDNSList into a slice of maps for schema.TypeList.
func flattenExtranetDNSList(list []extranet.ExtranetDNSList) []interface{} {
	if len(list) == 0 {
		return nil
	}
	out := make([]interface{}, 0, len(list))
	for _, v := range list {
		m := map[string]interface{}{
			"id":                   v.ID,
			"name":                 v.Name,
			"primary_dns_server":   v.PrimaryDNSServer,
			"secondary_dns_server": v.SecondaryDNSServer,
			"use_as_default":       v.UseAsDefault,
		}
		out = append(out, m)
	}
	return out
}

// flattenExtranetIpPoolList converts []extranet.ExtranetPoolList into a slice of maps for schema.TypeList.
func flattenExtranetIpPoolList(list []extranet.ExtranetPoolList) []interface{} {
	if len(list) == 0 {
		return nil
	}
	out := make([]interface{}, 0, len(list))
	for _, v := range list {
		m := map[string]interface{}{
			"id":             v.ID,
			"name":           v.Name,
			"ip_start":       v.IPStart,
			"ip_end":         v.IPEnd,
			"use_as_default": v.UseAsDefault,
		}
		out = append(out, m)
	}
	return out
}

func dataSourceExtranetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Fetching all extranets")
	allExtranets, err := extranet.GetAll(ctx, service, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all extranets: %w", err))
	}

	log.Printf("[DEBUG] Retrieved %d extranets", len(allExtranets))

	var resp *extranet.Extranet
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	// Search by ID first if provided
	if idProvided {
		log.Printf("[INFO] Searching for extranet by ID: %d", id)
		for i := range allExtranets {
			if allExtranets[i].ID == id {
				resp = &allExtranets[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("extranet with ID %d not found", id))
		}
	}

	// Search by name if not found by ID and name is provided (case-insensitive)
	if resp == nil && nameProvided && nameStr != "" {
		log.Printf("[INFO] Searching for extranet by name: %s", nameStr)
		for i := range allExtranets {
			if strings.EqualFold(allExtranets[i].Name, nameStr) {
				resp = &allExtranets[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("extranet with name %q not found", nameStr))
		}
	}

	// If neither ID nor name provided, or no match found
	if resp == nil {
		if idProvided || (nameProvided && nameStr != "") {
			return diag.FromErr(fmt.Errorf("no extranet found with name %q or id %d", nameStr, id))
		}
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	// Set the resource data from Extranet struct
	d.SetId(fmt.Sprintf("%d", resp.ID))
	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting id: %w", err))
	}
	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %w", err))
	}
	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(fmt.Errorf("error setting description: %w", err))
	}
	if err := d.Set("extranet_dns_list", flattenExtranetDNSList(resp.ExtranetDNSList)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting extranet_dns_list: %w", err))
	}
	if err := d.Set("extranet_ip_pool_list", flattenExtranetIpPoolList(resp.ExtranetIpPoolList)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting extranet_ip_pool_list: %w", err))
	}
	if err := d.Set("created_at", resp.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting created_at: %w", err))
	}
	if err := d.Set("modified_at", resp.ModifiedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting modified_at: %w", err))
	}

	log.Printf("[DEBUG] Extranet found: ID=%d, Name=%s", resp.ID, resp.Name)
	return nil
}
