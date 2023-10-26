package zia

/*
import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/zpa_gateways"
)

func dataSourceForwardingControlZPAGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataForwardingControlZPAGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "A unique identifier assigned to the ZPA gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the ZPA gateway",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional details about the ZPA gateway",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether the ZPA gateway is configured for Zscaler Internet Access (using option ZPA) or Zscaler Cloud Connector (using option ECZPA)",
			},
			"zpa_tenant_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the ZPA tenant where Source IP Anchoring is configured",
			},
			"zpa_server_group": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The ZPA Server Group that is configured for Source IP Anchoring",
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
			"zpa_application_segments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All the Application Segments that are associated with the selected ZPA Server Group for which Source IP Anchoring is enabled",
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
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the ZPA gateway was last modified",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "This is an immutable reference to an entity. which mainly consists of id and name",
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

func dataForwardingControlZPAGatewayRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *zpa_gateways.ZPAGateways
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for zpa gateway id: %d\n", id)
		res, err := zClient.zpa_gateways.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for zpa gateway name: %s\n", name)
		res, err := zClient.zpa_gateways.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("type", resp.Type)
		_ = d.Set("zpa_tenant_id", resp.ZPATenantId)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("zpa_server_group", flattenIDNameExtensions(resp.ZPAServerGroup)); err != nil {
			return err
		}

		if err := d.Set("zpa_application_segments", flattenIDNameExtensions(resp.ZPAAppSegments)); err != nil {
			return err
		}

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any zpa gateway name '%s' or id '%d'", name, id)
	}

	return nil
}
*/
