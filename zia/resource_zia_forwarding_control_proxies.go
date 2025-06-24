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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/proxies"
)

func resourceForwardingControlProxies() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceForwardingControlProxiesCreate,
		ReadContext:   resourceForwardingControlProxiesRead,
		UpdateContext: resourceForwardingControlProxiesUpdate,
		DeleteContext: resourceForwardingControlProxiesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("proxy_id", idInt)
				} else {
					resp, err := proxies.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("proxy_id", resp.ID)
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
			"proxy_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the Proxy gateway",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional details about the Proxy gateway",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Gateway type",
				ValidateFunc: validation.StringInSlice([]string{
					"PROXYCHAIN",
					"ZIA",
					"ECSELF",
				}, false),
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address or the FQDN of the third-party proxy service",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port number on which the third-party proxy service listens to the requests forwarded from Zscaler",
			},
			"insert_xau_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag indicating whether X-Authenticated-User header is added by the proxy. Enable to automatically insert authenticated user ID to the HTTP header, X-Authenticated-User.",
			},
			"base64_encode_xau_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag indicating whether the added X-Authenticated-User header is Base64 encoded. When enabled, the user ID is encoded using the Base64 encoding method.",
			},
			"cert": setSingleIDSchemaTypeCustom("The root certificate used by the third-party proxy to perform SSL inspection"),
		},
	}
}

func resourceForwardingControlProxiesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandForwardingControlProxies(d)
	log.Printf("[INFO] Creating ZIA proxies\n%+v\n", req)

	resp, _, err := proxies.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA proxies request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("proxy_id", resp.ID)

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceForwardingControlProxiesRead(ctx, d, meta)
}

func resourceForwardingControlProxiesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "proxy_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no proxies id is set"))
	}
	resp, err := proxies.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia proxies %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia proxies:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("proxy_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("type", resp.Type)
	_ = d.Set("address", resp.Address)
	_ = d.Set("port", resp.Port)
	_ = d.Set("insert_xau_header", resp.InsertXauHeader)
	_ = d.Set("base64_encode_xau_header", resp.Base64EncodeXauHeader)

	if err := d.Set("cert", flattenCert(resp.Cert)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceForwardingControlProxiesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "proxy_id")
	if !ok {
		log.Printf("[ERROR] proxy ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia proxy ID: %v\n", id)
	req := expandForwardingControlProxies(d)
	if _, err := proxies.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := proxies.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceForwardingControlProxiesRead(ctx, d, meta)
}

func resourceForwardingControlProxiesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "proxy_id")
	if !ok {
		log.Printf("[ERROR] proxy ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia proxy ID: %v\n", (d.Id()))

	if _, err := proxies.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia proxy deleted")

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandForwardingControlProxies(d *schema.ResourceData) proxies.Proxies {
	id, _ := getIntFromResourceData(d, "proxy_id")
	result := proxies.Proxies{
		ID:                    id,
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Type:                  d.Get("type").(string),
		Address:               d.Get("address").(string),
		Port:                  d.Get("port").(int),
		InsertXauHeader:       d.Get("insert_xau_header").(bool),
		Base64EncodeXauHeader: d.Get("base64_encode_xau_header").(bool),
		Cert:                  expandCert(d, "cert"),
	}
	return result
}

func flattenCert(customID *common.IDNameExternalID) []interface{} {
	if customID == nil || customID.ID == 0 {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id": customID.ID,
		},
	}
}

func expandCert(d *schema.ResourceData, key string) *common.IDNameExternalID {
	if v, ok := d.GetOk(key); ok {
		setList := v.(*schema.Set).List()
		if len(setList) > 0 {
			if idMap, ok := setList[0].(map[string]interface{}); ok {
				return &common.IDNameExternalID{
					ID: idMap["id"].(int),
				}
			}
		}
	}
	return nil
}
