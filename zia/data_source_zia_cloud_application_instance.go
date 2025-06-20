package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloud_app_instances"
)

func dataSourceCloudApplicationInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudApplicationInstanceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Unique identifier for the cloud application instance",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the cloud application instance",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the cloud application instance",
			},
			"modified_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp of when the cloud application instance was last modified",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the cloud application instance last",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"instance_identifiers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of identifiers for the cloud application instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique identifier for the cloud application instance",
						},
						"instance_identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifying string for the instance",
						},
						"instance_identifier_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifying string for the instance",
						},
						"identifier_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the cloud application instance",
						},
						"modified_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp of when the instance was last modified.",
						},
						"last_modified_by": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The admin that modified the instance last.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Identifier that uniquely identifies an entity",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The configured name of the entity",
									},
									"extensions": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudApplicationInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *cloud_app_instances.CloudApplicationInstances
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for cloud application instance id: %d\n", id)
		res, err := cloud_app_instances.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for cloud application instance name: %s\n", name)
		res, err := cloud_app_instances.GetInstanceByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.InstanceID))
		_ = d.Set("name", resp.InstanceName)
		_ = d.Set("instance_type", resp.InstanceType)
		_ = d.Set("modified_at", resp.ModifiedAt)

		if err := d.Set("instance_identifiers", flattenInstanceIdentifiers(resp.InstanceIdentifiers)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.ModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any cloud application instance name '%s' or id '%d'", name, id))
	}

	return nil
}

func flattenInstanceIdentifiers(items []cloud_app_instances.InstanceIdentifiers) []interface{} {
	var result []interface{}

	for _, item := range items {
		m := map[string]interface{}{
			"instance_id":              item.InstanceID,
			"instance_identifier":      item.InstanceIdentifier,
			"instance_identifier_name": item.InstanceIdentifierName,
			"identifier_type":          item.IdentifierType,
			"modified_at":              item.ModifiedAt,
			"last_modified_by":         flattenLastModifiedBy(item.ModifiedBy),
		}
		result = append(result, m)
	}

	return result
}
