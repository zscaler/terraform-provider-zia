package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
)

func dataSourceFWNetworkServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFWNetworkServicesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"src_tcp_ports":  dataNetworkPortsSchema("src tcp ports"),
			"dest_tcp_ports": dataNetworkPortsSchema("dest tcp ports"),
			"src_udp_ports":  dataNetworkPortsSchema("src udp ports"),
			"dest_udp_ports": dataNetworkPortsSchema("dest udp ports"),
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_name_l10n_tag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceFWNetworkServicesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *networkservices.NetworkServices
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting network services id: %d\n", id)
		res, err := networkservices.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)

	// Prepare optional parameters
	var protocol, locale *string
	if v, ok := d.GetOk("protocol"); ok && v.(string) != "" {
		protocolStr := v.(string)
		protocol = &protocolStr
		log.Printf("[DEBUG] Protocol parameter set: %s\n", *protocol)
	}
	if v, ok := d.GetOk("locale"); ok && v.(string) != "" {
		localeStr := v.(string)
		locale = &localeStr
		log.Printf("[DEBUG] Locale parameter set: %s\n", *locale)
	}

	log.Printf("[DEBUG] Search parameters - name: '%s', protocol: %v, locale: %v, resp is nil: %v\n", name, protocol, locale, resp == nil)

	// Search using GetByName - supports name, protocol, locale, or any combination
	// The SDK now handles empty name when protocol or locale is provided
	if resp == nil {
		if name != "" || protocol != nil || locale != nil {
			log.Printf("[INFO] Getting network services by name: '%s', protocol: %v, locale: %v\n", name, protocol, locale)
			res, err := networkservices.GetByName(ctx, service, name, protocol, locale)
			if err != nil {
				log.Printf("[ERROR] GetByName failed: %v\n", err)
				return diag.FromErr(err)
			}
			resp = res
			log.Printf("[INFO] Successfully retrieved network service: ID=%d, Name=%s\n", resp.ID, resp.Name)
		} else {
			log.Printf("[DEBUG] No search parameters provided (name, protocol, or locale)\n")
		}
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("tag", resp.Tag)
		_ = d.Set("type", resp.Type)
		_ = d.Set("description", resp.Description)
		_ = d.Set("protocol", resp.Protocol)
		_ = d.Set("is_name_l10n_tag", resp.IsNameL10nTag)

		if err := d.Set("src_tcp_ports", flattenNetwordPorts(resp.SrcTCPPorts)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dest_tcp_ports", flattenNetwordPorts(resp.DestTCPPorts)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("src_udp_ports", flattenNetwordPorts(resp.SrcUDPPorts)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dest_udp_ports", flattenNetwordPorts(resp.DestUDPPorts)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		var searchCriteria string
		if id > 0 {
			searchCriteria = fmt.Sprintf("id '%d'", id)
		} else if name != "" {
			searchCriteria = fmt.Sprintf("name '%s'", name)
			if protocol != nil {
				searchCriteria += fmt.Sprintf(", protocol '%s'", *protocol)
			}
			if locale != nil {
				searchCriteria += fmt.Sprintf(", locale '%s'", *locale)
			}
		} else if protocol != nil {
			searchCriteria = fmt.Sprintf("protocol '%s'", *protocol)
			if locale != nil {
				searchCriteria += fmt.Sprintf(", locale '%s'", *locale)
			}
		} else if locale != nil {
			searchCriteria = fmt.Sprintf("locale '%s'", *locale)
		} else {
			searchCriteria = "the provided criteria"
		}
		return diag.FromErr(fmt.Errorf("couldn't find any network service with %s", searchCriteria))
	}

	return nil
}
