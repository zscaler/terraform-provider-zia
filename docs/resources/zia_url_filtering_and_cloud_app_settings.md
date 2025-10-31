---
subcategory: "URL Filtering Rule"
layout: "zscaler"
page_title: "ZIA: url_filtering_and_cloud_app_settings"
description: |-
    Official documentation https://help.zscaler.com/zia/url-cloud-app-control-policy-settings#/advancedUrlFilterAndCloudAppSettings-get
    API documentation https://help.zscaler.com/zia/url-cloud-app-control-policy-settings#/advancedUrlFilterAndCloudAppSettings-get
  Updates the URL and Cloud App Control advanced policy settings
---

# zia_url_filtering_and_cloud_app_settings (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-url-categories)
* [API documentation](https://help.zscaler.com/zia/url-categories#/urlCategories-get)

The **zia_url_filtering_and_cloud_app_settings** resource allows you to updates the the URL and Cloud App Control advanced policy settings To learn more see [Configuring Advanced Policy Settings](https://help.zscaler.com/unified/configuring-advanced-policy-settings)

## Example Usage

```hcl
resource "zia_url_filtering_and_cloud_app_settings" "this" {
    block_skype                             = true
    consider_embedded_sites                 = false
    enable_block_override_for_non_auth_user = false
    enable_chatgpt_prompt                   = false
    enable_cipa_compliance                  = false
    enable_dynamic_content_cat              = true
    enable_gemini_prompt                    = false
    enable_meta_prompt                      = false
    enable_microsoft_copilot_prompt         = false
    enable_msft_o365                        = false
    enable_newly_registered_domains         = false
    enable_office365                        = true
    enable_per_plexity_prompt               = false
    enable_poep_prompt                      = false
    enable_ucaas_logmein                    = false
    enable_ucaas_ring_central               = false
    enable_ucaas_talkdesk                   = false
    enable_ucaas_webex                      = false
    enable_ucaas_zoom                       = false
    enforce_safe_search                     = false
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `enable_dynamic_content_cat` - (Boolean) A Boolean value that indicates if dynamic categorization of URLs by analyzing content of uncategorized websites using AI/ML tools is enabled or not.
* `consider_embedded_sites` - (Boolean) Indicates if URL filtering rules must be applied to sites that are translated using translation services.
* `enforce_safe_search` - (Boolean) Indicates whether only safe content must be returned for web, image, and video search.

* `safe_search_apps` - (List of String) A list of applications for which the SafeSearch enforcement applies. You cannot modify this field when the enforce_safe_search field is disabled. [See the URL & Cloud App Control Policy](https://help.zscaler.com/zia/url-cloud-app-control-policy-settings#/advancedUrlFilterAndCloudAppSettings-get) for the list of available apps

* `enable_office365` - (Boolean) Enables or disables Microsoft Office 365 configuration.
* `enable_msft_o365` - (Boolean) Enables or disables Microsoft-recommended Office 365 one-click configuration.
* `enable_ucaas_zoom` - (Boolean) Indicates if the Zscaler service is allowed to automatically permit secure local breakout for Zoom traffic.
* `enable_ucaas_logmein` - (Boolean) Indicates if the Zscaler service is allowed to automatically permit secure local breakout for GoTo traffic.
* `enable_ucaas_ringcentral` - (Boolean) Indicates if the Zscaler service is allowed to automatically permit secure local breakout for RingCentral traffic.
* `enable_ucaas_webex` - (Boolean) Indicates if the Zscaler service is allowed to automatically permit secure local breakout for Webex traffic.
* `enable_ucaas_talkdesk` - (Boolean) Indicates if the Zscaler service is allowed to automatically permit secure local breakout for Talkdesk traffic.
* `enable_chatgpt_prompt` - (Boolean) Indicates if the use of generative AI prompts with ChatGPT by users should be categorized and logged.
* `enable_microsoft_copilot_prompt` - (Boolean) Indicates if the use of generative AI prompts with Microsoft Copilot by users should be categorized and logged.
* `enable_gemini_prompt` - (Boolean) Indicates if the use of generative AI prompts with Google Gemini by users should be categorized and logged.
* `enable_poeprompt` - (Boolean) Indicates if the use of generative AI prompts with Poe by users should be categorized and logged.
* `enable_meta_prompt` - (Boolean) Indicates if the use of generative AI prompts with Meta AI by users should be categorized and logged.
* `enable_perplexity_prompt` - (Boolean) Indicates if the use of generative AI prompts with Perplexity by users should be categorized and logged.
* `block_skype` - (Boolean) Indicates whether access to Skype is blocked.
* `enable_newly_registered_domains` - (Boolean) Indicates whether newly registered and observed domains identified within hours of going live are allowed or blocked.
* `enable_block_override_for_non_auth_user` - (Boolean) Indicates if authorized users can temporarily override block action on websites by providing their authentication information.
* `enable_cipa_compliance` - (Boolean) Indicates if the predefined CIPA Compliance Rule is enabled.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_url_filtering_and_cloud_app_settings** can be imported by using `app_setting` as the import ID.

For example:

```shell
terraform import zia_url_filtering_and_cloud_app_settings.this "app_setting"
```
