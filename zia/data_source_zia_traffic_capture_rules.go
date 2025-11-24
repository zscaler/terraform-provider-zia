package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/traffic_capture"
)

func dataSourceTrafficCaptureRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrafficCaptureRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Unique identifier for the Traffic Capture policy rule",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of the Traffic Capture policy rule",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule order number of the Traffic Capture policy rule",
			},
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Admin rank of the Traffic Capture policy rule",
			},
			"access_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The admin's access privilege to this rule based on the assigned role",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action the Traffic Capture policy rule takes when packets match the rule",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines whether the Traffic Capture policy rule is enabled or disabled",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the rule",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the rule was last modified",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User who last modified the rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"src_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User-defined source IP addresses",
			},
			"dest_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination IP addresses",
			},
			"dest_ip_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "IP address categories of destination",
			},
			"dest_countries": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Destination countries",
			},
			"source_countries": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Source countries",
			},
			"exclude_src_countries": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether source countries are excluded",
			},
			"nw_applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Network service applications",
			},
			"default_rule": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if this is a default rule",
			},
			"predefined": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if this is a predefined rule",
			},
			"txn_size_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The maximum size of traffic to capture per connection",
			},
			"txn_sampling": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The percentage of connections sampled for capturing each time the rule is triggered",
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"device_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"app_service_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"dest_ip_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"nw_services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"nw_service_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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

			"nw_application_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"src_ip_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"departments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"workload_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of preconfigured workload groups to which the policy must be applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier assigned to the workload group",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the workload group",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the workload group",
						},
						"expression": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the workload group",
						},
						"last_modified_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_modified_by": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
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
						"expression_json": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expression_containers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"tag_container": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tags": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"key": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"value": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
															"operator": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"device_trust_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Device trust levels",
			},
		},
	}
}

func dataSourceTrafficCaptureRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *traffic_capture.TrafficCaptureRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting traffic capture rule by ID: %d\n", id)
		res, err := traffic_capture.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting traffic capture rule by name: %s\n", name)
		res, err := traffic_capture.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("order", resp.Order)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("access_control", resp.AccessControl)
		_ = d.Set("action", resp.Action)
		_ = d.Set("state", resp.State)
		_ = d.Set("description", resp.Description)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)
		_ = d.Set("src_ips", resp.SrcIps)
		_ = d.Set("dest_addresses", resp.DestAddresses)
		_ = d.Set("dest_ip_categories", resp.DestIpCategories)
		_ = d.Set("dest_countries", resp.DestCountries)
		_ = d.Set("source_countries", resp.SourceCountries)
		_ = d.Set("exclude_src_countries", resp.ExcludeSrcCountries)
		_ = d.Set("nw_applications", resp.NwApplications)
		_ = d.Set("default_rule", resp.DefaultRule)
		_ = d.Set("predefined", resp.Predefined)
		_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)
		_ = d.Set("txn_size_limit", resp.TxnSizeLimit)
		_ = d.Set("txn_sampling", resp.TxnSampling)

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("location_groups", flattenIDNameExtensions(resp.LocationsGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("departments", flattenIDNameExtensions(resp.Departments)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("groups", flattenIDNameExtensions(resp.Groups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("users", flattenIDNameExtensions(resp.Users)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("time_windows", flattenIDNameExtensions(resp.TimeWindows)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("nw_application_groups", flattenIDNameExtensions(resp.NwApplicationGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("app_service_groups", flattenIDNameExtensions(resp.AppServiceGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dest_ip_groups", flattenIDNameExtensions(resp.DestIpGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("nw_services", flattenIDNameExtensions(resp.NwServices)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("nw_service_groups", flattenIDNameExtensions(resp.NwServiceGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("src_ip_groups", flattenIDNameExtensions(resp.SrcIpGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("device_groups", flattenIDNameExtensions(resp.DeviceGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("devices", flattenIDNameExtensions(resp.Devices)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting workload_groups: %s", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any traffic capture rule with name '%s' or id '%d'", name, id))
	}

	return nil
}
