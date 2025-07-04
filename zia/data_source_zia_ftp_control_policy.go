package zia

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ftp_control_policy"
)

func dataSourceFTPControlPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFTPControlPolicyRead,
		Schema: map[string]*schema.Schema{
			"ftp_over_http_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to enable FTP over HTTP. ",
			},
			"ftp_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to enable native FTP. When enabled, users can connect to native FTP sites and download files.",
			},
			"url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of URL categories that allow FTP traffic",
			},
			"urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Domains or URLs included for the FTP Control settings",
			},
		},
	}
}

func dataSourceFTPControlPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	res, err := ftp_control_policy.GetFTPControlPolicy(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set ID for the data source
	d.SetId("ftp_control")

	// Map values to the Terraform schema
	_ = d.Set("ftp_over_http_enabled", res.FtpOverHttpEnabled)
	_ = d.Set("ftp_enabled", res.FtpEnabled)
	_ = d.Set("url_categories", res.UrlCategories)
	_ = d.Set("urls", res.Urls)
	return nil
}
