package zia

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_clusters"
)

func resourceVZENCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVZENClusterCreate,
		ReadContext:   resourceVZENClusterRead,
		UpdateContext: resourceVZENClusterUpdate,
		DeleteContext: resourceVZENClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("cluster_id", idInt)
				} else {
					resp, err := vzen_clusters.GetClusterByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("cluster_id", resp.ID)
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
			"cluster_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the Virtual Service Edge cluster",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the status of the Virtual Service Edge cluster. The status is set to ENABLED by default",
			},
			"ip_sec_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Virtual Service Edge cluster IP address",
			},
			"subnet_mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Virtual Service Edge cluster subnet mask",
			},
			"default_gateway": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the default gateway to the internet",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Virtual Service Edge cluster subnet mask",
				ValidateFunc: validation.StringInSlice([]string{
					`ANY`, `NONE`, `SME`, `SMSM`, `SMCA`, `SMUI`, `SMCDS`, `SMDNSD`, `SMAA`, `SMTP`, `SMQTN`, `VIP`,
					`UIZ`, `UIAE`, `SITEREVIEW`, `PAC`, `S_RELAY`, `M_RELAY`, `H_MON`, `SMIKE`, `NSS`, `SMEZA`, `SMLB`,
					`SMFCCLT`, `SMBA`, `SMBAC`, `SMESXI`, `SMBAUI`, `VZEN`, `ZSCMCLT`, `SMDLP`, `ZSQUERY`, `ADP`, `SMCDSDLP`,
					`SMSCIM`, `ZSAPI`, `ZSCMCDSSCLT`, `LOCAL_MTS`, `SVPN`, `SMCASB`, `SMFALCONUI`, `MOBILEAPP_REG`, `SMRESTSVR`,
					`FALCONCA`, `MOBILEAPP_NF`, `ZIRSVR`, `SMEDGEUI`, `ALERTEVAL`, `ALERTNOTIF`, `SMPARTNERUI`, `CQM`, `DATAKEEPER`,
					`SMBAM`, `ZWACLT`,
				}, false),
			},
			"virtual_zen_nodes": setIDExternalIDCustom(intPtr(255), "List of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector)."),
		},
	}
}

func resourceVZENClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandVZENClusters(d)
	log.Printf("[INFO] Creating ZIA vzen clusters\n%+v\n", req)

	resp, _, err := vzen_clusters.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA vzen clusters request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("cluster_id", resp.ID)

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

	return resourceVZENClusterRead(ctx, d, meta)
}

func resourceVZENClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "cluster_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no vzen clusters id is set"))
	}
	resp, err := vzen_clusters.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia vzen clusters %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia vzen clusters:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("status", resp.Status)
	_ = d.Set("type", resp.Type)
	_ = d.Set("ip_address", resp.IpAddress)
	_ = d.Set("subnet_mask", resp.SubnetMask)
	_ = d.Set("default_gateway", resp.DefaultGateway)
	_ = d.Set("ip_sec_enabled", resp.IpSecEnabled)

	if err := d.Set("virtual_zen_nodes", flattenCommonIDNameExternalIDSimple(resp.VirtualZenNodes)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVZENClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "cluster_id")
	if !ok {
		log.Printf("[ERROR] vzen cluster ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia vzen cluster ID: %v\n", id)

	req := expandVZENClusters(d)
	if _, err := vzen_clusters.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := vzen_clusters.Update(ctx, service, id, &req); err != nil {
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

	return resourceVZENClusterRead(ctx, d, meta)
}

func resourceVZENClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "cluster_id")
	if !ok {
		log.Printf("[ERROR] vzen cluster ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia vzen cluster ID: %v\n", (d.Id()))

	if _, err := vzen_clusters.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia vzen cluster deleted")

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

func expandVZENClusters(d *schema.ResourceData) vzen_clusters.VZENClusters {
	id, _ := getIntFromResourceData(d, "cluster_id")
	result := vzen_clusters.VZENClusters{
		ID:              id,
		Name:            d.Get("name").(string),
		Status:          d.Get("status").(string),
		Type:            d.Get("type").(string),
		IpAddress:       d.Get("ip_address").(string),
		SubnetMask:      d.Get("subnet_mask").(string),
		DefaultGateway:  d.Get("default_gateway").(string),
		IpSecEnabled:    d.Get("ip_sec_enabled").(bool),
		VirtualZenNodes: expandCommonIDNameExternalID(d, "virtual_zen_nodes"),
	}
	return result
}

func flattenCommonIDNameExternalIDSimple(list []common.IDNameExternalID) []interface{} {
	if len(list) == 0 {
		return nil
	}

	var ids []int
	for _, v := range list {
		ids = append(ids, v.ID)
	}
	sort.Ints(ids)

	var idInterfaces []interface{}
	for _, id := range ids {
		idInterfaces = append(idInterfaces, id)
	}

	return []interface{}{
		map[string]interface{}{
			"id": idInterfaces,
		},
	}
}

func expandCommonIDNameExternalID(d *schema.ResourceData, key string) []common.IDNameExternalID {
	raw, exists := d.GetOk(key)
	if !exists {
		return nil
	}

	blocks := raw.(*schema.Set).List()
	var result []common.IDNameExternalID

	for _, block := range blocks {
		m := block.(map[string]interface{})

		// Expecting "id" to be a list of ints
		if idList, ok := m["id"].([]interface{}); ok {
			for _, id := range idList {
				result = append(result, common.IDNameExternalID{
					ID: id.(int),
				})
			}
		}
	}
	return result
}
