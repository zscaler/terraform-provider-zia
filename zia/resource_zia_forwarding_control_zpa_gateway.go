package zia

import (
	"fmt"
	"log"
	"strconv"
	"time"

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
					_ = d.Set("gateway_id", idInt)
				} else {
					resp, err := zClient.zpa_gateways.GetByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("gateway_id", resp.ID)
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
			"gateway_id": {
				Type:     schema.TypeInt,
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
			"zpa_server_group": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The ZPA Server Group that is configured for Source IP Anchoring",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the ZPA Gateway.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the ZPA Gateway.",
						},
					},
				},
			},
			"zpa_app_segments": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "All the Application Segments that are associated with the selected ZPA Server Group for which Source IP Anchoring is enabled",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the application segment.",
						},
						"external_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "External ID of the application segment.",
						},
					},
				},
			},
		},
	}
}

func validatePredefinedObject(req zpa_gateways.ZPAGateways) error {
	if req.Name == "Auto ZPA Gateway" {
		return fmt.Errorf("predefined zpa gateway '%s' cannot be deleted", req.Name)
	}
	return nil
}

func resourceForwardingControlZPAGatewayCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandForwardingControlZPAGateway(d)
	log.Printf("[INFO] Creating forwarding control zpa gateway\n%+v\n", req)

	if err := validatePredefinedObject(req); err != nil {
		return err
	}

	resp, err := zClient.zpa_gateways.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created forwarding control zpa gateway request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("gateway_id", resp.ID)
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceForwardingControlZPAGatewayRead(d, m)
}

func resourceForwardingControlZPAGatewayRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[DEBUG] Current value of gateway_id: %v", d.Get("gateway_id"))

	id, ok := getIntFromResourceData(d, "gateway_id")

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
	_ = d.Set("gateway_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	log.Printf("[DEBUG] Type returned from API: %s", resp.Type)
	if resp.Type == "" {
		resp.Type = d.Get("type").(string)
	}
	_ = d.Set("type", resp.Type)

	if err := d.Set("zpa_server_group", flattenZPAServerGroupSimple(resp.ZPAServerGroup)); err != nil {
		return err
	}
	if err := d.Set("zpa_app_segments", flattenZPAGWAppSegments(resp.ZPAAppSegments)); err != nil {
		return err
	}

	return nil
}

func resourceForwardingControlZPAGatewayUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		log.Printf("[ERROR] forwarding control zpa gateway ID not set: %v\n", id)
	}
	if !d.HasChange("type") || d.Get("type") == "" {
		d.Set("type", "ZPA")
	}

	log.Printf("[INFO] Updating zia forwarding control zpa gateway ID: %v\n", id)
	req := expandForwardingControlZPAGateway(d)

	if err := validatePredefinedObject(req); err != nil {
		return err
	}
	if _, err := zClient.zpa_gateways.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := zClient.zpa_gateways.Update(id, &req); err != nil {
		return err
	}
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceForwardingControlZPAGatewayRead(d, m)
}

func resourceForwardingControlZPAGatewayDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "gateway_id")
	if !ok {
		log.Printf("[ERROR] forwarding control zpa gateway not set: %v\n", id)
	}
	// Retrieve the rule to check if it's a predefined one
	gwObj, err := zClient.zpa_gateways.Get(id)
	if err != nil {
		return fmt.Errorf("error retrieving zpa gateway object %d: %v", id, err)
	}

	// Validate if the rule can be deleted
	if err := validatePredefinedObject(*gwObj); err != nil {
		return err
	}

	log.Printf("[INFO] Deleting forwarding control zpa gateway ID: %v\n", (d.Id()))

	if _, err := zClient.zpa_gateways.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] forwarding control zpa gateway deleted")
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandForwardingControlZPAGateway(d *schema.ResourceData) zpa_gateways.ZPAGateways {
	id, _ := getIntFromResourceData(d, "gateway_id")
	gatewayType, exists := d.GetOk("type")
	if !exists {
		gatewayType = "ZPA"
	}
	result := zpa_gateways.ZPAGateways{
		ID:             id,
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Type:           gatewayType.(string),
		ZPAServerGroup: expandZPAServerGroup(d, "zpa_server_group"),
		ZPAAppSegments: expandZPAGWAppSegment(d, "zpa_app_segments"),
	}
	return result
}

func expandZPAServerGroup(d *schema.ResourceData, key string) zpa_gateways.ZPAServerGroup {
	listInterface, exists := d.GetOk(key)
	if !exists || len(listInterface.([]interface{})) == 0 {
		return zpa_gateways.ZPAServerGroup{}
	}

	groupMap := listInterface.([]interface{})[0].(map[string]interface{})

	return zpa_gateways.ZPAServerGroup{
		ExternalID: groupMap["external_id"].(string),
		Name:       groupMap["name"].(string),
	}
}

func flattenZPAServerGroupSimple(serverGroup zpa_gateways.ZPAServerGroup) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":        serverGroup.Name,
			"external_id": serverGroup.ExternalID,
		},
	}
}

func expandZPAGWAppSegment(d *schema.ResourceData, key string) []zpa_gateways.ZPAAppSegments {
	setInterface, exists := d.GetOk(key)
	if !exists {
		return nil
	}

	inputSet := setInterface.(*schema.Set).List()
	var result []zpa_gateways.ZPAAppSegments
	for _, item := range inputSet {
		itemMap := item.(map[string]interface{})
		name := itemMap["name"].(string)
		externalID := itemMap["external_id"].(string)

		segment := zpa_gateways.ZPAAppSegments{
			Name:       name,
			ExternalID: externalID,
		}
		result = append(result, segment)
	}
	return result
}

func flattenZPAGWAppSegments(list []zpa_gateways.ZPAAppSegments) []interface{} {
	flattenedList := make([]interface{}, 0, len(list))
	for _, val := range list {
		r := map[string]interface{}{
			"name":        val.Name,
			"external_id": val.ExternalID,
		}
		flattenedList = append(flattenedList, r)
	}
	return flattenedList
}
