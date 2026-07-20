package zia

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/end_user_notification"
)

func dataSourceEUNUserConfirmationTemplateProduct() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEUNUserConfirmationTemplateProductRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"INLINE", "ENDPOINT_DLP", "CLOUDAPP", "URL", "FILE_TYPE", "FIREWALL", "DNS", "IPS"}, false),
				Description:  "The policy type associated with the user confirmation notification template. User confirmation notifications are supported for Inline Web DLP and Endpoint DLP policies. Supported values are `INLINE`, `ENDPOINT_DLP`, `CLOUDAPP`, `URL`, `FILE_TYPE`, `FIREWALL`, `DNS`, and `IPS`.",
			},
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the user confirmation notification template to look up.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the user confirmation notification template to look up.",
			},
			"channel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The channel associated with the user confirmation notification template.",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default user confirmation notification template for the policy type.",
			},
			"language_templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of per-language user confirmation messages associated with the template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"language": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEUNUserConfirmationTemplateProductRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	product := d.Get("product").(string)

	templates, err := end_user_notification.GetEunTemplateByPolicy(ctx, service, product)
	if err != nil {
		return diag.FromErr(err)
	}

	// The service returns the full list for the given policy type without a
	// server-side search key or pagination, so the requested template is
	// resolved locally by id or name.
	var match *end_user_notification.UserConfirmationByPolicyType

	if id, ok := getIntFromResourceData(d, "id"); ok && id != 0 {
		for i := range templates {
			if templates[i].ID == id {
				match = &templates[i]
				break
			}
		}
		if match == nil {
			return diag.FromErr(fmt.Errorf("no user confirmation template found with id '%d' for product '%s'", id, product))
		}
	} else if name, ok := d.Get("name").(string); ok && name != "" {
		for i := range templates {
			if templates[i].Name == name {
				match = &templates[i]
				break
			}
		}
		if match == nil {
			return diag.FromErr(fmt.Errorf("no user confirmation template found with name '%s' for product '%s'", name, product))
		}
	} else {
		// When neither id nor name is provided, resolve to the default template
		// for the policy type, falling back to the first entry returned.
		for i := range templates {
			if templates[i].Default {
				match = &templates[i]
				break
			}
		}
		if match == nil && len(templates) > 0 {
			match = &templates[0]
		}
		if match == nil {
			return diag.FromErr(fmt.Errorf("no user confirmation templates found for product '%s'", product))
		}
	}

	d.SetId(strconv.Itoa(match.ID))
	_ = d.Set("id", match.ID)
	_ = d.Set("name", match.Name)
	_ = d.Set("channel", match.Channel)
	_ = d.Set("product", match.Product)
	_ = d.Set("default", match.Default)
	_ = d.Set("language_templates", flattenEUNUserConfirmationLanguageTemplates(match.LanguageTemplates))

	return nil
}

func flattenEUNUserConfirmationLanguageTemplates(templates []end_user_notification.LanguageTemplatesLite) []interface{} {
	result := make([]interface{}, 0, len(templates))
	for _, t := range templates {
		result = append(result, map[string]interface{}{
			"language": t.Language,
			"message":  t.Message,
			"default":  t.Default,
		})
	}
	return result
}
