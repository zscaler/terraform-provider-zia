package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudnss/nss_servers"
)

func resourceNSSServers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNSSServersCreate,
		ReadContext:   resourceNSSServersRead,
		UpdateContext: resourceNSSServersUpdate,
		DeleteContext: resourceNSSServersDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("nss_id", idInt)
				} else {
					resp, err := nss_servers.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("nss_id", resp.ID)
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
			"nss_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The NSS server name",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enables or disables the status of the NSS server",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the NSS Server",
			},
			"icap_svr_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ICAP server ID",
			},
		},
	}
}

func resourceNSSServersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandNSSServers(d)
	log.Printf("[INFO] Creating ZIA nss servers\n%+v\n", req)

	resp, err := nss_servers.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA nss servers request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("nss_id", resp.ID)

	time.Sleep(2 * time.Second)

	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceNSSServersRead(ctx, d, meta)
}

func resourceNSSServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "nss_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no nss servers id is set"))
	}
	resp, err := nss_servers.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia nss servers %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia nss servers:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("nss_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("status", resp.Status)
	_ = d.Set("type", resp.Type)
	_ = d.Set("icap_svr_id", resp.IcapSvrId)

	return nil
}

func resourceNSSServersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "nss_id")
	if !ok {
		log.Printf("[ERROR] nss servers ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia nss server ID: %v\n", id)
	req := expandNSSServers(d)
	if _, err := nss_servers.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := nss_servers.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	time.Sleep(2 * time.Second)

	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceNSSServersRead(ctx, d, meta)
}

func resourceNSSServersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "nss_id")
	if !ok {
		log.Printf("[ERROR] nss server ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia nss server ID: %v\n", (d.Id()))

	if _, err := nss_servers.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia nss server deleted")

	time.Sleep(2 * time.Second)

	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandNSSServers(d *schema.ResourceData) nss_servers.NSSServers {
	id, _ := getIntFromResourceData(d, "nss_id")
	result := nss_servers.NSSServers{
		ID:        id,
		Name:      d.Get("name").(string),
		Status:    d.Get("status").(string),
		Type:      d.Get("type").(string),
		IcapSvrId: d.Get("icap_svr_id").(int),
	}
	return result
}
