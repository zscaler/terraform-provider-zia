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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloud_app_instances"
)

func resourceCloudApplicationInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudApplicationInstanceCreate,
		ReadContext:   resourceCloudApplicationInstanceRead,
		UpdateContext: resourceCloudApplicationInstanceUpdate,
		DeleteContext: resourceCloudApplicationInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("instance_id", idInt)
				} else {
					resp, err := cloud_app_instances.GetInstanceByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.InstanceID))
						_ = d.Set("instance_id", resp.InstanceID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the cloud application instance.",
			},
			"instance_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the cloud application instance.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the cloud application instance.",
			},
			"instance_identifiers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of identifiers for the cloud application instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique identifier for the cloud application instance.",
						},
						"instance_identifier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unique identifying string for the instance.",
						},
						"instance_identifier_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unique identifying string for the instance.",
						},
						"identifier_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of the cloud application instance.",
						},
					},
				},
			},
		},
	}
}

func resourceCloudApplicationInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandCloudApplicationInstance(d)
	log.Printf("[INFO] Creating ZIA cloud application instances\n%+v\n", req)

	resp, _, err := cloud_app_instances.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA cloud application instances request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.InstanceID))
	_ = d.Set("instance_id", resp.InstanceID)

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceCloudApplicationInstanceRead(ctx, d, meta)
}

func resourceCloudApplicationInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "instance_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no cloud application instances id is set"))
	}
	resp, err := cloud_app_instances.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia cloud application instances %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia cloud application instances:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.InstanceID))
	_ = d.Set("instance_name", resp.InstanceName)
	_ = d.Set("instance_type", resp.InstanceType)

	if err := d.Set("instance_identifiers", flattenInstanceIdentifiersSimple(resp.InstanceIdentifiers)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCloudApplicationInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "instance_id")
	if !ok {
		log.Printf("[ERROR] cloud application instance ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia cloud application instance ID: %v\n", id)
	req := expandCloudApplicationInstance(d)
	if _, err := cloud_app_instances.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := cloud_app_instances.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceCloudApplicationInstanceRead(ctx, d, meta)
}

func resourceCloudApplicationInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "instance_id")
	if !ok {
		log.Printf("[ERROR] cloud application instance ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia cloud application instance ID: %v\n", (d.Id()))

	if _, err := cloud_app_instances.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia cloud application instance deleted")

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandCloudApplicationInstance(d *schema.ResourceData) cloud_app_instances.CloudApplicationInstances {
	id, _ := getIntFromResourceData(d, "instance_id")
	result := cloud_app_instances.CloudApplicationInstances{
		InstanceID:          id,
		InstanceName:        d.Get("name").(string),
		InstanceType:        d.Get("instance_type").(string),
		InstanceIdentifiers: expandInstanceIdentifiers(d),
	}
	return result
}

func expandInstanceIdentifiers(d *schema.ResourceData) []cloud_app_instances.InstanceIdentifiers {
	rawList, ok := d.Get("instance_identifiers").([]interface{})
	if !ok || len(rawList) == 0 {
		return nil
	}

	var result []cloud_app_instances.InstanceIdentifiers
	for _, raw := range rawList {
		if raw == nil {
			continue
		}
		item := raw.(map[string]interface{})

		instanceID := 0
		if v, ok := item["instance_id"].(int); ok {
			instanceID = v
		}

		result = append(result, cloud_app_instances.InstanceIdentifiers{
			InstanceID:             instanceID,
			InstanceIdentifier:     item["instance_identifier"].(string),
			IdentifierType:         item["identifier_type"].(string),
			InstanceIdentifierName: item["instance_identifier_name"].(string),
		})
	}
	return result
}

func flattenInstanceIdentifiersSimple(items []cloud_app_instances.InstanceIdentifiers) []interface{} {
	var result []interface{}

	for _, item := range items {
		m := map[string]interface{}{
			"instance_id":              item.InstanceID,
			"instance_identifier":      item.InstanceIdentifier,
			"instance_identifier_name": item.InstanceIdentifierName,
			"identifier_type":          item.IdentifierType,
		}
		result = append(result, m)
	}

	return result
}
