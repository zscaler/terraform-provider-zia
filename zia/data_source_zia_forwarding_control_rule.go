package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/forwarding_rules"
)

func dataSourceForwardingControlRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceForwardingControlRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A unique identifier assigned to the forwarding rule",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the forwarding rule",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the forwarding rule",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The rule type selected from the available options",
			},
			"forward_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of traffic forwarding method selected from the available options",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether the forwarding rule is enabled or disabled",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The order of execution for the forwarding rule order",
			},
			"rank": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Admin rank assigned to the forwarding rule",
			},
			"ec_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies",
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
					},
				},
			},
			"locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.",
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
					},
				},
			},
			"location_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the location groups to which the forwarding rule applies",
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
					},
				},
			},
			"departments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the departments to which the forwarding rule applies. If not set, the rule applies to all departments.",
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
					},
				},
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the user groups to which the forwarding rule applies. If not set, the rule applies to all groups.",
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
					},
				},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the users to which the forwarding rule applies. If not set, user criteria is ignored during policy enforcement.",
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
					},
				},
			},
			"src_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.",
			},
			"src_ip_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.",
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
					},
				},
			},
			"src_ipv6_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.",
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
					},
				},
			},
			"dest_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses. If not set, the rule is not restricted to a specific destination addresses unless specified by destCountries, destIpGroups, or destIpCategories.",
			},
			"dest_ip_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination IP categories to which the rule applies. If not set, the rule is not restricted to specific destination IP categories.",
			},
			"res_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination domain categories to which the rule applies",
			},
			"dest_countries": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.",
			},
			"dest_ip_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User-defined destination IP address groups to which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group.",
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
					},
				},
			},
			"dest_ipv6_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Destination IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.Note: User-defined groups for IPv6 addresses are currently not supported.",
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
					},
				},
			},
			"nw_services": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User-defined network services to which the rule applies. If not set, the rule is not restricted to a specific network service. Note: When the forwarding method is Proxy Chaining, only TCP-based network services are considered for policy match .",
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
					},
				},
			},
			"nw_service_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User-defined network service group to which the rule applies. If not set, the rule is not restricted to a specific network service group.",
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
					},
				},
			},
			"nw_applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User-defined network service applications to which the rule applies. If not set, the rule is not restricted to a specific network service application.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"nw_application_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User-defined network service application groups to which the rule applied. If not set, the rule is not restricted to a specific network service application group.",
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
					},
				},
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
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
					},
				},
			},
			"proxy_gateway": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method.",
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
			"zpa_application_segments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector).",
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
						"ddescription": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Additional information about the Application Segment",
						},
						"zpa_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the ZPA tenant where the Application Segment is configured",
						},
						"deleted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "ID of the ZPA tenant where the Application Segment is configured",
						},
					},
				},
			},
			"zpa_application_segment_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of ZPA Application Segment Groups for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier assigned to the Application Segment Group",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Application Segment Group",
						},
						"zpa_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the external ID. Applicable only when this reference is of an external entity.",
						},
						"deleted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the ZPA Application Segment Group has been deleted",
						},
						"zpa_app_segments_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of ZPA Application Segments in the group",
						},
					},
				},
			},
			"zpa_gateway": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The ZPA Server Group for which this rule is applicable. Only the Server Groups that are associated with the selected Application Segments are allowed. This field is applicable only for the ZPA forwarding method.",
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
			"zpa_broker_rule": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The predefined ZPA Broker Rule generated by Zscaler",
			},
			"devices": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of devices for which the rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.",
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
			"device_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of device groups for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.",
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
	}
}

func dataSourceForwardingControlRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *forwarding_rules.ForwardingRules
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for forwarding control rule id: %d\n", id)
		res, err := forwarding_rules.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for forwarding control rule : %s\n", name)
		res, err := forwarding_rules.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("forward_method", resp.ForwardMethod)
		_ = d.Set("order", resp.Order)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("state", resp.State)
		_ = d.Set("type", resp.Type)
		_ = d.Set("src_ips", resp.SrcIps)
		_ = d.Set("dest_addresses", resp.DestAddresses)
		_ = d.Set("dest_ip_categories", resp.DestIpCategories)
		_ = d.Set("dest_countries", resp.DestCountries)
		_ = d.Set("res_categories", resp.ResCategories)
		_ = d.Set("zpa_broker_rule", resp.ZPABrokerRule)

		if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("location_groups", flattenIDNameExtensions(resp.LocationsGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("ec_groups", flattenIDNameExtensions(resp.ECGroups)); err != nil {
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

		if err := d.Set("src_ip_groups", flattenIDNameExtensions(resp.SrcIpGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("src_ipv6_groups", flattenIDNameExtensions(resp.SrcIpv6Groups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dest_ip_groups", flattenIDNameExtensions(resp.DestIpGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("dest_ipv6_groups", flattenIDNameExtensions(resp.DestIpv6Groups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("nw_services", flattenIDNameExtensions(resp.NwServices)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("nw_service_groups", flattenIDNameExtensions(resp.NwServiceGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("nw_application_groups", flattenIDNameExtensions(resp.NwApplicationGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("proxy_gateway", flattenIDNameSet(resp.ProxyGateway)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("zpa_gateway", flattenIDNameSet(resp.ZPAGateway)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("zpa_app_segments", flattenZPAAppSegments(resp.ZPAAppSegments)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("zpa_application_segments", flattenZPAApplicationSegments(resp.ZPAApplicationSegments)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("zpa_application_segment_groups", flattenZPAApplicationSegmentGroups(resp.ZPAApplicationSegmentGroups)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any forwarding rule with name '%s' or id '%d'", name, id))
	}

	return nil
}

func flattenZPAAppSegments(list []common.ZPAAppSegments) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		r := map[string]interface{}{
			"id":          val.ID,
			"name":        val.Name,
			"external_id": val.ExternalID,
		}
		flattenedList[i] = r
	}
	return flattenedList
}

func flattenZPAApplicationSegments(list []forwarding_rules.ZPAApplicationSegments) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		r := map[string]interface{}{
			"id":          val.ID,
			"name":        val.Name,
			"description": val.Description,
			"zpa_id":      val.ZPAID,
			"deleted":     val.Deleted,
		}
		flattenedList[i] = r
	}
	return flattenedList
}

func flattenZPAApplicationSegmentGroups(list []forwarding_rules.ZPAApplicationSegmentGroups) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		r := map[string]interface{}{
			"id":                     val.ID,
			"name":                   val.Name,
			"zpa_id":                 val.ZPAID,
			"deleted":                val.Deleted,
			"zpa_app_segments_count": val.ZPAAppSegmentsCount,
		}
		flattenedList[i] = r
	}
	return flattenedList
}
