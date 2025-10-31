package zia

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/end_user_notification"
)

func resourceEndUserNotification() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceEndUserNotificationRead,
		CreateContext: resourceEndUserNotificationCreate,
		UpdateContext: resourceEndUserNotificationUpdate,
		DeleteContext: resourceFuncNoOp,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// notificationType := d.Get("notification_type").(string)

			// Validation for notification_type = "DEFAULT"
			// if notificationType == "DEFAULT" {
			// 	forbiddenFields := []string{"redirect_url", "security_review_custom_location", "url_cat_review_custom_location"}
			// 	for _, field := range forbiddenFields {
			// 		if d.HasChange(field) && d.Get(field) != "" {
			// 			return fmt.Errorf("attribute '%s' cannot be set when 'notification_type' is 'DEFAULT'", field)
			// 		}
			// 	}
			// }

			// Validation for notification_type = "CUSTOM"
			// if notificationType == "CUSTOM" {
			// 	forbiddenFields := []string{"display_reason", "display_comp_name", "display_company_logo"}
			// 	for _, field := range forbiddenFields {
			// 		if d.HasChange(field) && d.Get(field) != "" {
			// 			return fmt.Errorf("attribute '%s' cannot be set when 'notification_type' is 'CUSTOM'", field)
			// 		}
			// 	}
			// }

			// Validation for url_cat_review_submit_to_security_cloud
			// if d.Get("url_cat_review_submit_to_security_cloud").(bool) {
			// 	if d.HasChange("url_cat_review_custom_location") && d.Get("url_cat_review_custom_location") != "" {
			// 		return fmt.Errorf("attribute 'url_cat_review_custom_location' cannot be set when 'url_cat_review_submit_to_security_cloud' is true")
			// 	}
			// }

			// Validation for security_review_submit_to_security_cloud
			// if d.Get("security_review_submit_to_security_cloud").(bool) {
			// 	if d.HasChange("security_review_custom_location") && d.Get("security_review_custom_location") != "" {
			// 		return fmt.Errorf("attribute 'security_review_custom_location' cannot be set when 'security_review_submit_to_security_cloud' is true")
			// 	}
			// }

			// Validation for aup_frequency
			aupFrequency := d.Get("aup_frequency").(string)

			switch aupFrequency {
			case "CUSTOM":
				if d.Get("aup_day_offset").(int) == 0 {
					return fmt.Errorf("attribute 'aup_day_offset' must be set when 'aup_frequency' is 'CUSTOM'")
				}
				if d.Get("aup_custom_frequency").(int) == 0 {
					return fmt.Errorf("attribute 'aup_custom_frequency' must be set when 'aup_frequency' is 'CUSTOM'")
				}

			case "DAILY", "WEEKLY", "ON_DATE", "ON_WEEKDAY":
				if d.Get("aup_day_offset").(int) == 0 {
					return fmt.Errorf("attribute 'aup_day_offset' must be set when 'aup_frequency' is '%s'", aupFrequency)
				}
			}

			return nil
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				diags := resourceEndUserNotificationRead(ctx, d, meta)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to read advanced settings during import: %s", diags[0].Summary)
				}
				d.SetId("enduser_notification")
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"aup_frequency": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "The frequency at which the Acceptable Use Policy (AUP) is shown to the end users",
				ValidateFunc: validation.StringInSlice([]string{
					"NEVER",
					"SESSION",
					"ONLOGIN",
					"CUSTOM",
					"DAILY",
					"WEEKLY",
					"ON_DATE",
					"ON_WEEKDAY",
				}, false),
			},
			"aup_custom_frequency": {
				Type: schema.TypeInt,
				// Computed:    true,
				Optional:     true,
				Description:  "The custom frequency (in days) for showing the AUP to the end users. Valid range is 0 to 180.",
				ValidateFunc: validation.IntBetween(0, 180),
			},
			"aup_day_offset": {
				Type: schema.TypeInt,
				// Computed:    true,
				Optional:     true,
				Description:  "Specifies which day of the week or month the AUP is shown for users when aupFrequency is set. Valid range is 0 to 31.",
				ValidateFunc: validation.IntBetween(0, 31),
			},
			"aup_message": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "The acceptable use statement that is shown in the AUP",
			},
			"notification_type": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "The type of EUN as default or custom",
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"CUSTOM",
				}, false),
			},
			"display_reason": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether or not the reason for cautioning or blocking access to a site, file, or application is shown when the respective notification is triggered",
			},
			"display_company_name": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether the organization's name appears in the EUN or not",
			},
			"display_company_logo": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether your organization's logo appears in the EUN or not",
			},
			"custom_text": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:         true,
				Description:      "The custom text shown in the EUN",
				ValidateDiagFunc: stringIsMultiLine,
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"url_cat_review_enabled": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether the URL Categorization notification is enabled or disabled",
			},
			"url_cat_review_submit_to_security_cloud": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether users' review requests for possibly misclassified URLs are submitted to the Zscaler service (i.e., Security Cloud) or a custom location.",
			},
			"url_cat_review_custom_location": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "A custom URL location where users' review requests for blocked URLs are sent",
			},
			"url_cat_review_text": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:         true,
				Description:      "The message that appears in the URL Categorization notification",
				ValidateDiagFunc: stringIsMultiLine,
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"security_review_enabled": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether the Security Violation notification is enabled or disabled",
			},
			"security_review_submit_to_security_cloud": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether users' review requests for blocked URLs are submitted to the Zscaler service (i.e., Security Cloud) or a custom location.",
			},
			"security_review_custom_location": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "Value indicating whether or not to include the ECS option in all DNS queries, originating from all locations and remote users.",
			},
			"security_review_text": {
				Type: schema.TypeString,
				// Computed:         true,
				Optional:         true,
				Description:      "The message that appears in the Security Violation notification",
				ValidateDiagFunc: stringIsMultiLine,
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"web_dlp_review_enabled": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether the Web DLP Violation notification is enabled or disabled",
			},
			"web_dlp_review_submit_to_security_cloud": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether users' review requests for web DLP policy violation are submitted to the Zscaler service (i.e., Security Cloud) or a custom location.",
			},
			"web_dlp_review_custom_location": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "A custom URL location where users' review requests for the web DLP policy violation are sent",
			},
			"web_dlp_review_text": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:         true,
				Description:      "The message that appears in the Web DLP Violation notification",
				ValidateDiagFunc: stringIsMultiLine,
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"redirect_url": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "The redirect URL for the external site hosting the EUN specified when the custom notification type is selected",
			},
			"support_email": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "The email address for writing to IT Support",
			},
			"support_phone": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "The phone number for contacting IT Support",
			},
			"org_policy_link": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:    true,
				Description: "The URL of the organization's policy page. This field is required for the default notification type.",
			},
			"caution_again_after": {
				Type: schema.TypeInt,
				// Computed:    true,
				Optional:     true,
				Description:  "The time interval at which the caution notification is shown when users continue browsing a restricted site.",
				ValidateFunc: validation.IntAtLeast(5),
			},
			"caution_per_domain": {
				Type: schema.TypeBool,
				// Computed:    true,
				Optional:    true,
				Description: "Specifies whether to display the caution notification at a specific time interval for URLs in the Miscellaneous or Unknown category.",
			},
			"caution_custom_text": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:         true,
				Description:      "The custom message that appears in the caution notification",
				ValidateDiagFunc: stringIsMultiLine,
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"idp_proxy_notification_text": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional:         true,
				Description:      "The message that appears in the IdP Proxy notification",
				ValidateDiagFunc: stringIsMultiLine,
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"quarantine_custom_notification_text": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The message that appears in the quarantine notification",
				ValidateDiagFunc: stringIsMultiLine,
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
		},
	}
}

func resourceEndUserNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandEndUserNotification(d)
	_, _, err := end_user_notification.UpdateUserNotificationSettings(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("enduser_notification")

	// Sleep for 1 seconds before potentially triggering the activation
	time.Sleep(1 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceEndUserNotificationRead(ctx, d, meta)
}

func resourceEndUserNotificationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	// Apply formatting fixes before storing in state
	_ = d.Set("quarantine_custom_notification_text", normalizeMultiLineString(res.QuarantineCustomNotificationText))

	return nil
}

func resourceEndUserNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandEndUserNotification(d)

	_, _, err := end_user_notification.UpdateUserNotificationSettings(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("enduser_notification")

	// Sleep for 1 seconds before potentially triggering the activation
	time.Sleep(1 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceEndUserNotificationRead(ctx, d, meta)
}

func expandEndUserNotification(d *schema.ResourceData) end_user_notification.UserNotificationSettings {
	result := end_user_notification.UserNotificationSettings{
		AUPFrequency:                        d.Get("aup_frequency").(string),
		AUPCustomFrequency:                  d.Get("aup_custom_frequency").(int),
		AUPDayOffset:                        d.Get("aup_day_offset").(int),
		AUPMessage:                          d.Get("aup_message").(string),
		NotificationType:                    d.Get("notification_type").(string),
		DisplayReason:                       getBoolFromResourceData(d, "display_reason"),
		DisplayCompName:                     getBoolFromResourceData(d, "display_company_name"),
		DisplayCompLogo:                     getBoolFromResourceData(d, "display_company_logo"),
		CustomText:                          d.Get("custom_text").(string),
		URLCatReviewEnabled:                 getBoolFromResourceData(d, "url_cat_review_enabled"),
		URLCatReviewSubmitToSecurityCloud:   getBoolFromResourceData(d, "url_cat_review_submit_to_security_cloud"),
		URLCatReviewCustomLocation:          d.Get("url_cat_review_custom_location").(string),
		URLCatReviewText:                    d.Get("url_cat_review_text").(string),
		SecurityReviewEnabled:               getBoolFromResourceData(d, "security_review_enabled"),
		SecurityReviewSubmitToSecurityCloud: getBoolFromResourceData(d, "security_review_submit_to_security_cloud"),
		SecurityReviewCustomLocation:        d.Get("security_review_custom_location").(string),
		SecurityReviewText:                  d.Get("security_review_text").(string),
		WebDLPReviewEnabled:                 getBoolFromResourceData(d, "web_dlp_review_enabled"),
		WebDLPReviewSubmitToSecurityCloud:   getBoolFromResourceData(d, "web_dlp_review_submit_to_security_cloud"),
		WebDLPReviewCustomLocation:          d.Get("web_dlp_review_custom_location").(string),
		WebDLPReviewText:                    d.Get("web_dlp_review_text").(string),
		RedirectURL:                         d.Get("redirect_url").(string),
		SupportEmail:                        d.Get("support_email").(string),
		SupportPhone:                        d.Get("support_phone").(string),
		OrgPolicyLink:                       d.Get("org_policy_link").(string),
		CautionAgainAfter:                   d.Get("caution_again_after").(int),
		CautionPerDomain:                    getBoolFromResourceData(d, "caution_per_domain"),
		CautionCustomText:                   d.Get("caution_custom_text").(string),
		IDPProxyNotificationText:            d.Get("idp_proxy_notification_text").(string),
		QuarantineCustomNotificationText:    d.Get("quarantine_custom_notification_text").(string),
	}
	return result
}
