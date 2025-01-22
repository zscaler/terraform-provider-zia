package zia

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/end_user_notification"
)

func dataSourceEndUserNotification() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEndUserNotificationRead,
		Schema: map[string]*schema.Schema{
			"aup_frequency": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The frequency at which the Acceptable Use Policy (AUP) is shown to the end users",
			},
			"aup_custom_frequency": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The custom frequency (in days) for showing the AUP to the end users. Valid range is 1 to 180.",
			},
			"aup_day_offset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies which day of the week or month the AUP is shown for users when aupFrequency is set. Valid range is 1 to 31.",
			},
			"aup_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The acceptable use statement that is shown in the AUP",
			},
			"notification_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of EUN as default or custom",
			},
			"display_reason": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether or not the reason for cautioning or blocking access to a site, file, or application is shown when the respective notification is triggered",
			},
			"display_company_name": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether the organization's name appears in the EUN or not",
			},
			"display_company_logo": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether your organization's logo appears in the EUN or not",
			},
			"custom_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The custom text shown in the EUN",
			},
			"url_cat_review_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether the URL Categorization notification is enabled or disabled",
			},
			"url_cat_review_submit_to_security_cloud": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether users' review requests for possibly misclassified URLs are submitted to the Zscaler service (i.e., Security Cloud) or a custom location.",
			},
			"url_cat_review_custom_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A custom URL location where users' review requests for blocked URLs are sent",
			},
			"url_cat_review_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The message that appears in the URL Categorization notification",
			},
			"security_review_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether the Security Violation notification is enabled or disabled",
			},
			"security_review_submit_to_security_cloud": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether users' review requests for blocked URLs are submitted to the Zscaler service (i.e., Security Cloud) or a custom location.",
			},
			"security_review_custom_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value indicating whether or not to include the ECS option in all DNS queries, originating from all locations and remote users.",
			},
			"security_review_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The message that appears in the Security Violation notification",
			},
			"web_dlp_review_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether the Web DLP Violation notification is enabled or disabled",
			},
			"web_dlp_review_submit_to_security_cloud": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether users' review requests for web DLP policy violation are submitted to the Zscaler service (i.e., Security Cloud) or a custom location.",
			},
			"web_dlp_review_custom_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A custom URL location where users' review requests for the web DLP policy violation are sent",
			},
			"web_dlp_review_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The message that appears in the Web DLP Violation notification",
			},
			"redirect_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The redirect URL for the external site hosting the EUN specified when the custom notification type is selected",
			},
			"support_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address for writing to IT Support",
			},
			"support_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The phone number for contacting IT Support",
			},
			"org_policy_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the organization's policy page. This field is required for the default notification type.",
			},
			"caution_again_after": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The time interval at which the caution notification is shown when users continue browsing a restricted site.",
			},
			"caution_per_domain": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether to display the caution notification at a specific time interval for URLs in the Miscellaneous or Unknown category.",
			},
			"caution_custom_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The custom message that appears in the caution notification",
			},
			"idp_proxy_notification_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The message that appears in the IdP Proxy notification",
			},
			"quarantine_custom_notification_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The message that appears in the quarantine notification",
			},
		},
	}
}

func dataSourceEndUserNotificationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Fetch data from the API
	res, err := end_user_notification.GetUserNotificationSettings(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set ID for the data source
	d.SetId("enduser_notification")

	_ = d.Set("aup_frequency", res.AUPFrequency)
	_ = d.Set("aup_custom_frequency", res.AUPCustomFrequency)
	_ = d.Set("aup_day_offset", res.AUPDayOffset)
	_ = d.Set("aup_message", res.AUPMessage)
	_ = d.Set("notification_type", res.NotificationType)
	_ = d.Set("display_reason", res.DisplayReason)
	_ = d.Set("display_company_name", res.DisplayCompName)
	_ = d.Set("display_company_logo", res.DisplayCompLogo)
	_ = d.Set("custom_text", res.CustomText)
	_ = d.Set("url_cat_review_enabled", res.URLCatReviewEnabled)
	_ = d.Set("url_cat_review_submit_to_security_cloud", res.URLCatReviewSubmitToSecurityCloud)
	_ = d.Set("url_cat_review_custom_location", res.URLCatReviewCustomLocation)
	_ = d.Set("url_cat_review_text", res.URLCatReviewText)
	_ = d.Set("security_review_enabled", res.SecurityReviewEnabled)
	_ = d.Set("security_review_submit_to_security_cloud", res.SecurityReviewSubmitToSecurityCloud)
	_ = d.Set("security_review_custom_location", res.SecurityReviewCustomLocation)
	_ = d.Set("security_review_text", res.SecurityReviewText)
	_ = d.Set("web_dlp_review_enabled", res.WebDLPReviewEnabled)
	_ = d.Set("web_dlp_review_submit_to_security_cloud", res.WebDLPReviewSubmitToSecurityCloud)
	_ = d.Set("web_dlp_review_custom_location", res.WebDLPReviewCustomLocation)
	_ = d.Set("web_dlp_review_text", res.WebDLPReviewText)
	_ = d.Set("redirect_url", res.RedirectURL)
	_ = d.Set("support_email", res.SupportEmail)
	_ = d.Set("support_phone", res.SupportPhone)
	_ = d.Set("org_policy_link", res.OrgPolicyLink)
	_ = d.Set("caution_again_after", res.CautionAgainAfter)
	_ = d.Set("caution_per_domain", res.CautionPerDomain)
	_ = d.Set("caution_custom_text", res.CautionCustomText)
	_ = d.Set("idp_proxy_notification_text", res.IDPProxyNotificationText)
	_ = d.Set("quarantine_custom_notification_text", res.QuarantineCustomNotificationText)

	return nil
}
