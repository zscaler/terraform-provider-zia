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

func dataSourceEUNTemplateProduct() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEUNTemplateProductRead,
		Schema: map[string]*schema.Schema{
			"template_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ZCC", "BROWSER"}, false),
				Description:  "The type of notification template. Supported values are `ZCC` (Zscaler Client Connector-based) and `BROWSER` (browser-based).",
			},
			"product": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"INLINE", "ENDPOINT_DLP", "CLOUDAPP", "URL", "FILE_TYPE", "FIREWALL", "DNS", "IPS"}, false),
				Description:  "The policy type associated with the notification template. Supported values are `INLINE`, `ENDPOINT_DLP`, `CLOUDAPP`, `URL`, `FILE_TYPE`, `FIREWALL`, `DNS`, and `IPS`.",
			},
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the notification template to look up.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the notification template to look up.",
			},
			"channel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The channel associated with the notification template.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The notification template type as returned by the service.",
			},
			"caution_interval": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The interval used for caution notifications.",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default notification template for the policy type.",
			},
			"notification_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of notification details associated with the template.",
			},
			"recommended_cloud_app": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The recommended cloud application associated with the notification template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"val": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"misc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_not_ready": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"under_migration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"app_cat_modified": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deprecated": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"language_templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of per-language notification messages associated with the template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"language": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allow_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"block_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypt_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"readonly_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"caution_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"redirect_response_message": {
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

func dataSourceEUNTemplateProductRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	templateType := d.Get("template_type").(string)
	product := d.Get("product").(string)

	templates, err := end_user_notification.GetEunTemplateBrowserBasedZCC(ctx, service, templateType, product)
	if err != nil {
		return diag.FromErr(err)
	}

	// The service returns the full list for the given template type and policy
	// type without a server-side search key or pagination, so the requested
	// template is resolved locally by id or name.
	var match *end_user_notification.EunTemplateProduct

	if id, ok := getIntFromResourceData(d, "id"); ok && id != 0 {
		for i := range templates {
			if templates[i].ID == id {
				match = &templates[i]
				break
			}
		}
		if match == nil {
			return diag.FromErr(fmt.Errorf("no notification template found with id '%d' for template_type '%s' and product '%s'", id, templateType, product))
		}
	} else if name, ok := d.Get("name").(string); ok && name != "" {
		for i := range templates {
			if templates[i].Name == name {
				match = &templates[i]
				break
			}
		}
		if match == nil {
			return diag.FromErr(fmt.Errorf("no notification template found with name '%s' for template_type '%s' and product '%s'", name, templateType, product))
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
			return diag.FromErr(fmt.Errorf("no notification templates found for template_type '%s' and product '%s'", templateType, product))
		}
	}

	d.SetId(strconv.Itoa(match.ID))
	_ = d.Set("id", match.ID)
	_ = d.Set("name", match.Name)
	_ = d.Set("channel", match.Channel)
	_ = d.Set("product", match.Product)
	_ = d.Set("type", match.Type)
	_ = d.Set("caution_interval", match.CautionInterval)
	_ = d.Set("default", match.Default)
	_ = d.Set("notification_details", match.NotificationDetails)
	_ = d.Set("recommended_cloud_app", flattenEUNRecommendedCloudApp(match.RecommendedCloudApp))
	_ = d.Set("language_templates", flattenEUNLanguageTemplates(match.LanguageTemplates))

	return nil
}

func flattenEUNRecommendedCloudApp(app end_user_notification.RecommendedCloudApp) []interface{} {
	if (app == end_user_notification.RecommendedCloudApp{}) {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"val":              app.Val,
			"name":             app.Name,
			"channel":          app.Channel,
			"product":          app.Product,
			"type":             app.Type,
			"misc":             app.CautionMiscInterval,
			"app_not_ready":    app.AppNotReady,
			"under_migration":  app.UnderMigration,
			"app_cat_modified": app.AppCatModified,
			"deprecated":       app.Deprecated,
		},
	}
}

func flattenEUNLanguageTemplates(templates []end_user_notification.LanguageTemplates) []interface{} {
	result := make([]interface{}, 0, len(templates))
	for _, t := range templates {
		result = append(result, map[string]interface{}{
			"language":                  t.Language,
			"allow_message":             t.AllowMessage,
			"block_message":             t.BlockMessage,
			"encrypt_message":           t.EncryptMessage,
			"readonly_message":          t.ReadonlyMessage,
			"caution_message":           t.CautionMessage,
			"redirect_response_message": t.RedirectResponseMessage,
			"default":                   t.Default,
		})
	}
	return result
}
