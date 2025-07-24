package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_control_rules"
)

func dataSourceBandwdithControlRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBandwdithControlRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the rule label.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The bandwidth control rule name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the bandwidth control rule",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether the rule is enabled or disabled",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The order of the bandwidth control rule",
			},
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Admin rank of the Bandwidth Control policy rule",
			},
			"access_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access privilege of this rule based on the admin's Role Based Authorization (RBA) state",
			},
			"protocols": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The protocol criteria specified for the bandwidth control rule",
			},
			"min_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum percentage of a location's bandwidth you want to be guaranteed for each selected bandwidth class",
			},
			"max_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum percentage of a location's bandwidth you want to be guaranteed for each selected bandwidth class",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the DLP policy rule was last modified.",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the DLP policy rule last.",
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
							Description: "Identifier that uniquely identifies an entity",
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
			"labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Labels that are applicable to the rule",
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
							Description: "Identifier that uniquely identifies an entity",
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
			"locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations to which the bandwidth control rule must be applied",
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
							Description: "Identifier that uniquely identifies an entity",
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
			"location_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of locations groups to which the bandwidth control rule must be applied.",
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
							Description: "Identifier that uniquely identifies an entity",
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
			"time_windows": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Name-ID pairs of time windows to which the bandwidth control rule must be applied.",
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
							Description: "Identifier that uniquely identifies an entity",
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
			"bandwidth_classes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The bandwidth classes to which you want to apply this rule",
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
							Description: "Identifier that uniquely identifies an entity",
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
	}
}

func dataBandwdithControlRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *bandwidth_control_rules.BandwidthControlRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for bandwidth control rule  id: %d\n", id)
		res, err := bandwidth_control_rules.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for bandwidth control rule name: %s\n", name)
		res, err := bandwidth_control_rules.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("state", resp.State)
		_ = d.Set("order", resp.Order)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("access_control", resp.AccessControl)
		_ = d.Set("protocols", resp.Protocols)
		_ = d.Set("min_bandwidth", resp.MinBandwidth)
		_ = d.Set("max_bandwidth", resp.MaxBandwidth)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("bandwidth_classes", flattenIDExtensions(resp.BandwidthClasses)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("labels", flattenIDExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("locations", flattenIDExtensions(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("location_groups", flattenIDExtensions(resp.LocationGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("time_windows", flattenIDExtensions(resp.TimeWindows)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any bandwidth classname '%s' or id '%d'", name, id))
	}

	return nil
}
