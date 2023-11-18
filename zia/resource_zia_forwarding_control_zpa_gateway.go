package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/zpa_gateways"
)

func resourceForwardingControlZPAGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceForwardingControlZPAGatewayCreate,
		Read:   resourceForwardingControlZPAGatewayRead,
		Update: resourceForwardingControlZPAGatewayUpdate,
		Delete: resourceForwardingControlZPAGatewayDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("rule_id", idInt)
				} else {
					resp, err := zClient.zpa_gateways.GetByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("rule_id", resp.ID)
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ZPA gateway",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional details about the ZPA gateway",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates whether the ZPA gateway is configured for Zscaler Internet Access (using option ZPA) or Zscaler Cloud Connector (using option ECZPA)",
				ValidateFunc: validation.StringInSlice([]string{
					"ZPA",
					"ECZPA",
				}, false),
			},
			"zpa_tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the ZPA tenant where Source IP Anchoring is configured",
			},
			"zpa_server_group":         setIDsSchemaTypeCustom(intPtr(1), "The ZPA Server Group that is configured for Source IP Anchoring"),
			"zpa_application_segments": setIDsSchemaTypeCustom(intPtr(1), "The Name-ID pairs of locations groups to which the DLP policy rule must be applied."),
		},
	}
}

func resourceForwardingControlZPAGatewayCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandForwardingControlZPAGateway(d)
	log.Printf("[INFO] Creating forwarding control zpa gateway\n%+v\n", req)

	resp, err := zClient.zpa_gateways.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created forwarding control zpa gateway request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("id", resp.ID)
	return resourceFWIPDestinationGroupsRead(d, m)
}

func resourceForwardingControlZPAGatewayRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "id")
	if !ok {
		return fmt.Errorf("no forwarding control zpa gateway id is set")
	}
	resp, err := zClient.zpa_gateways.Get(id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing forwarding control zpa gateway %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting forwarding control zpa gateway:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("type", resp.Type)
	_ = d.Set("zpa_tenant_id", resp.ZPATenantId)

	if err := d.Set("zpa_server_group", flattenZPAServerGroupID(resp.ZPAServerGroup)); err != nil {
		return err
	}

	if err := d.Set("zpa_app_segments", flattenZPAAppSegmentsID(resp.ZPAAppSegments)); err != nil {
		return err
	}

	return nil
}

func resourceForwardingControlZPAGatewayUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "id")
	if !ok {
		log.Printf("[ERROR] forwarding control zpa gateway ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia forwarding control zpa gateway ID: %v\n", id)
	req := expandForwardingControlZPAGateway(d)
	if _, err := zClient.zpa_gateways.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := zClient.zpa_gateways.Update(id, &req); err != nil {
		return err
	}

	return resourceForwardingControlZPAGatewayRead(d, m)
}

func resourceForwardingControlZPAGatewayDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "id")
	if !ok {
		log.Printf("[ERROR] forwarding control zpa gateway not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting forwarding control zpa gateway ID: %v\n", (d.Id()))

	if _, err := zClient.zpa_gateways.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] forwarding control zpa gateway deleted")
	return nil
}

func expandForwardingControlZPAGateway(d *schema.ResourceData) zpa_gateways.ZPAGateways {
	id, _ := getIntFromResourceData(d, "id")
	result := zpa_gateways.ZPAGateways{
		ID:             id,
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Type:           d.Get("type").(string),
		ZPATenantId:    d.Get("zpa_tenant_id").(int),
		ZPAServerGroup: expandZPAServerGroupID(d, "zpa_server_group"),
		ZPAAppSegments: expandZPAAppSegmentsID(d, "zpa_app_segments"),
	}
	return result
}
