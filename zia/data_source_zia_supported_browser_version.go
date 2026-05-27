package zia

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/secure_browsing"
)

func dataSourceSupportedBrowserVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSupportedBrowserVersionRead,
		Schema: map[string]*schema.Schema{
			"browser_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CHROME",
					"FIREFOX",
					"SAFARI",
					"OPERA",
					"MSCHREDGE",
				}, false),
				Description: "Optional filter — return only the supported version entry for this browser type. One of: `CHROME`, `FIREFOX`, `SAFARI`, `OPERA`, `MSCHREDGE`.",
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JMESPath expression applied to the result set client-side. Filter on the API field names: `browserType`, `versions`, `olderVersions`. Example: `\"[?contains(versions, 'C130X')]\"`.",
			},
			"browsers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of supported browser version entries, one element per browser type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"browser_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The browser type.",
						},
						"versions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of currently supported versions for this browser type.",
						},
						"older_versions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of older / no-longer-current versions for this browser type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceSupportedBrowserVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	if searchExpr, ok := d.GetOk("search"); ok {
		ctx = zscaler.ContextWithJMESPath(ctx, searchExpr.(string))
		log.Printf("[INFO] zia_supported_browser_version JMESPath filter: %s\n", searchExpr.(string))
	}

	res, err := secure_browsing.GetSupportedBrowserVersions(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	// The SDK call is a single GET (not paginated), so JMESPath context
	// enrichment isn't auto-applied. Apply it explicitly here.
	res, err = zscaler.ApplyJMESPathFromContext(ctx, res)
	if err != nil {
		return diag.FromErr(err)
	}

	if bt, ok := d.GetOk("browser_type"); ok {
		wanted := bt.(string)
		filtered := res[:0]
		for _, b := range res {
			if b.BrowserType == wanted {
				filtered = append(filtered, b)
			}
		}
		res = filtered
	}

	browsers := make([]map[string]interface{}, 0, len(res))
	for _, b := range res {
		browsers = append(browsers, map[string]interface{}{
			"browser_type":   b.BrowserType,
			"versions":       b.Versions,
			"older_versions": b.OlderVersions,
		})
	}

	d.SetId("supported_browser_versions")
	if err := d.Set("browsers", browsers); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
