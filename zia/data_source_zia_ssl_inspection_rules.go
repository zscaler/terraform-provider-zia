package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sslinspection"
)

func dataSourceSSLInspectionRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSSLInspectionRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the SSL Inspection rule",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the SSL Inspection rule",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enables or disables the SSL Inspection rules.",
			},
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Admin rank of the admin who creates this rule",
			},
			// "access_control": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The rule order of execution for the  SSL Inspection rules with respect to other rules.",
			},
			"road_warrior_for_kerberos": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When set to true, the rule is applied to remote users that use PAC with Kerberos authentication.",
			},
			"platforms": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Zscaler Client Connector device platforms for which the rule must be applied.",
			},
			"url_categories": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of URL categories for which rule must be applied",
			},
			"cloud_applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of cloud applications to which the SSL Inspection rule must be applied.",
			},
			"user_agent_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User agent type list",
			},
			"device_trust_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation.",
			},
			"action": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the action configuration for the SSL inspection rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"show_eun": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"show_eunatp": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"override_default_certificate": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ssl_interception_cert": {
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
									"default_certificate": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"decrypt_sub_actions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_certificates": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ocsp_check": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"block_ssl_traffic_with_no_sni_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"min_client_tls_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"min_server_tls_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"block_undecrypt": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"http2_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"do_not_decrypt_sub_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bypass_other_policies": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"server_certificates": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ocsp_check": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"block_ssl_traffic_with_no_sni_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"min_tls_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of locations for the which policy must be applied. If not set, policy is applied for all locations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
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
				Description: "Name-ID pairs of locations groups for which rule must be applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"departments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of department for which the rule is applied. If not set, rule will be applied for all departments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of groups for which the policy must be applied. If not set, policy is applied for all groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of users for the which policy must be applied. If not set, user criteria is not considered for policy enforcement.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
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
				Description: "This is an immutable reference to an entity that mainly consists of id and name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the rule was last modified. Ignored if the request is POST or PUT.",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Admin user that last modified the rule. Ignored if the request is POST or PUT.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
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
				Description: "This is an immutable reference to an entity that mainly consists of id and name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"device_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of device groups for which the rule is applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"devices": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of device for which the rule is applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Additional information about the entity",
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
			"source_ip_groups": {
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
			"proxy_gateways": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier assigned to the Application Segment",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Application Segment",
						},
						// "external_id": {
						// 	Type:        schema.TypeString,
						// 	Computed:    true,
						// 	Description: "Indicates the external ID. Applicable only when this reference is of an external entity.",
						// },
					},
				},
			},
			"zpa_app_segments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier assigned to the Application Segment",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Application Segment",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the external ID. Applicable only when this reference is of an external entity.",
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
						// "default_rule": {
						// 	Type:     schema.TypeBool,
						// 	Computed: true,
						// },
						// "predefined": {
						// 	Type:     schema.TypeBool,
						// 	Computed: true,
						// },
					},
				},
			},
		},
	}
}

func dataSourceSSLInspectionRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *sslinspection.SSLInspectionRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for SSL Inspection rule id: %d\n", id)
		rule, err := sslinspection.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = rule
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data SSL Inspection rule : %s\n", name)
		res, err := sslinspection.GetByName(ctx, service, name)
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
		// _ = d.Set("access_control", resp.AccessControl)
		_ = d.Set("platforms", resp.Platforms)
		_ = d.Set("cloud_applications", resp.CloudApplications)
		_ = d.Set("road_warrior_for_kerberos", resp.RoadWarriorForKerberos)
		_ = d.Set("url_categories", resp.URLCategories)
		_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)
		_ = d.Set("user_agent_types", resp.UserAgentTypes)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		// Now set the action block
		if err := d.Set("action", flattenSSLInspectionAction(resp.Action)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("devices", flattenIDNameExtensions(resp.Devices)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("device_groups", flattenIDNameExtensions(resp.DeviceGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("location_groups", flattenIDNameExtensions(resp.LocationGroups)); err != nil {
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

		if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("source_ip_groups", flattenIDNameExtensions(resp.SourceIPGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dest_ip_groups", flattenIDNameExtensions(resp.DestIpGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("proxy_gateways", flattenIDExtensionsListIDs(resp.ProxyGateways)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("zpa_app_segments", flattenZPAAppSegments(resp.ZPAAppSegments)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any user with name '%s' or id '%d'", name, id))
	}

	return nil
}

func flattenSSLInspectionAction(action sslinspection.Action) []interface{} {
	a := make(map[string]interface{})

	a["type"] = action.Type
	a["show_eun"] = action.ShowEUN
	a["show_eunatp"] = action.ShowEUNATP
	a["override_default_certificate"] = action.OverrideDefaultCertificate

	// ssl_interception_cert
	a["ssl_interception_cert"] = flattenSSLInterceptionCert(action.SSLInterceptionCert)

	// do_not_decrypt_sub_actions
	a["do_not_decrypt_sub_actions"] = flattenDoNotDecryptSubActions(action.DoNotDecryptSubActions)

	// decrypt_sub_actions (NEW)
	a["decrypt_sub_actions"] = flattenDecryptSubActions(action.DecryptSubActions)

	return []interface{}{a}
}

func flattenSSLInterceptionCert(cert *sslinspection.SSLInterceptionCert) []interface{} {
	if cert == nil {
		return []interface{}{}
	}
	c := make(map[string]interface{})
	c["id"] = cert.ID
	c["name"] = cert.Name
	// c["default_certificate"] = cert.DefaultCertificate
	return []interface{}{c}
}

func flattenDoNotDecryptSubActions(subActions *sslinspection.DoNotDecryptSubActions) []interface{} {
	if subActions == nil {
		return []interface{}{}
	}
	sa := make(map[string]interface{})
	sa["bypass_other_policies"] = subActions.BypassOtherPolicies
	sa["server_certificates"] = subActions.ServerCertificates
	sa["ocsp_check"] = subActions.OcspCheck
	sa["block_ssl_traffic_with_no_sni_enabled"] = subActions.BlockSslTrafficWithNoSniEnabled
	sa["min_tls_version"] = subActions.MinTLSVersion // Or minServerTLSVersion if the API expects that exact string
	return []interface{}{sa}
}

func flattenDecryptSubActions(subActions *sslinspection.DecryptSubActions) []interface{} {
	if subActions == nil {
		return []interface{}{}
	}
	sa := make(map[string]interface{})
	sa["server_certificates"] = subActions.ServerCertificates
	sa["ocsp_check"] = subActions.OcspCheck
	sa["block_ssl_traffic_with_no_sni_enabled"] = subActions.BlockSslTrafficWithNoSniEnabled
	sa["min_client_tls_version"] = subActions.MinClientTLSVersion
	sa["min_server_tls_version"] = subActions.MinServerTLSVersion
	sa["block_undecrypt"] = subActions.BlockUndecrypt
	sa["http2_enabled"] = subActions.HTTP2Enabled

	return []interface{}{sa}
}
