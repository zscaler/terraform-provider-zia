package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/extranet"
)

func resourceExtranet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExtranetCreate,
		ReadContext:   resourceExtranetRead,
		UpdateContext: resourceExtranetUpdate,
		DeleteContext: resourceExtranetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("extranet_id", idInt)
				} else {
					resp, err := extranet.GetExtranetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("extranet_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extranet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
				Description:  "The name of the extranet.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
				Description:      "The description of the extranet.",
			},
			"extranet_dns_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "DNS servers specified for the extranet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID generated for the DNS server configuration.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the DNS server.",
						},
						"primary_dns_server": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address of the primary DNS server.",
						},
						"secondary_dns_server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IP address of the secondary DNS server.",
						},
						"use_as_default": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether this DNS server configuration is the designated default.",
						},
					},
				},
			},
			"extranet_ip_pool_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Traffic selector IP pools specified for the extranet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID generated for the IP pool configuration.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the IP pool.",
						},
						"ip_start": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The starting IP address of the pool.",
						},
						"ip_end": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ending IP address of the pool.",
						},
						"use_as_default": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether this IP pool is the designated default.",
						},
					},
				},
			},
		},
	}
}

func resourceExtranetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandExtranet(d)
	log.Printf("[INFO] Creating ZIA extranet\n%+v\n", req)

	resp, _, err := extranet.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA extranet request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("extranet_id", resp.ID)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceExtranetRead(ctx, d, meta)
}

func resourceExtranetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "extranet_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no extranet_id is set"))
	}
	resp, err := extranet.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia_extranet %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("extranet_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("extranet_dns_list", flattenExtranetDNSList(resp.ExtranetDNSList))
	_ = d.Set("extranet_ip_pool_list", flattenExtranetIpPoolList(resp.ExtranetIpPoolList))

	return nil
}

func resourceExtranetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "extranet_id")
	if !ok {
		log.Printf("[ERROR] cloud ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia extranet ID: %v\n", id)
	req := expandExtranet(d)
	if _, err := extranet.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := extranet.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceExtranetRead(ctx, d, meta)
}

func resourceExtranetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "extranet_id")
	if !ok {
		log.Printf("[ERROR] cloud ID not set: %v\n", id)
	}

	if _, err := extranet.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := extranet.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia_extranet %d deleted", id)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandExtranet(d *schema.ResourceData) extranet.Extranet {
	id, _ := getIntFromResourceData(d, "extranet_id")
	result := extranet.Extranet{
		ID:                 id,
		Name:               d.Get("name").(string),
		Description:        getString(d.Get("description")),
		ExtranetDNSList:    expandExtranetDNSList(d.Get("extranet_dns_list")),
		ExtranetIpPoolList: expandExtranetIpPoolList(d.Get("extranet_ip_pool_list")),
	}
	return result
}

func expandExtranetDNSList(v interface{}) []extranet.ExtranetDNSList {
	if v == nil {
		return nil
	}
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	out := make([]extranet.ExtranetDNSList, 0, len(list))
	for _, raw := range list {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		entry := extranet.ExtranetDNSList{
			Name:               item["name"].(string),
			PrimaryDNSServer:   item["primary_dns_server"].(string),
			SecondaryDNSServer: getString(item["secondary_dns_server"]),
			UseAsDefault:       getBool(item["use_as_default"]),
		}
		if id, ok := getIntFromNested(item, "id"); ok && id > 0 {
			entry.ID = id
		}
		out = append(out, entry)
	}
	return out
}

func expandExtranetIpPoolList(v interface{}) []extranet.ExtranetPoolList {
	if v == nil {
		return nil
	}
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	out := make([]extranet.ExtranetPoolList, 0, len(list))
	for _, raw := range list {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		entry := extranet.ExtranetPoolList{
			Name:         item["name"].(string),
			IPStart:      item["ip_start"].(string),
			IPEnd:        item["ip_end"].(string),
			UseAsDefault: getBool(item["use_as_default"]),
		}
		if id, ok := getIntFromNested(item, "id"); ok && id > 0 {
			entry.ID = id
		}
		out = append(out, entry)
	}
	return out
}

func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func getBool(v interface{}) bool {
	if v == nil {
		return false
	}
	b, _ := v.(bool)
	return b
}

func getIntFromNested(m map[string]interface{}, key string) (int, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case int:
		return t, true
	case int64:
		return int(t), true
	case float64:
		return int(t), true
	default:
		return 0, false
	}
}
