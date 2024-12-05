---
subcategory: "URL Filtering Rule"
layout: "zscaler"
page_title: "ZIA: url_filtering_and_cloud_app_settings"
description: |-
    Retrieves information about URL and Cloud App Control advanced policy settings.
---

# Data Source: zia_url_filtering_and_cloud_app_settings

Use the **zia_url_filtering_and_cloud_app_settings** data source to get information about URL and Cloud App Control advanced policy settings.

```hcl
# data "zia_url_filtering_and_cloud_app_settings" "this" {}
```

## Schema

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `enable_dynamic_content_cat` - (Boolean) A Boolean value that indicates if dynamic categorization of URLs by analyzing content of uncategorized websites using AI/ML tools is enabled or not.
* `consider_embedded_sites` - (Boolean) Indicates if URL filtering rules must be applied to sites that are translated using translation services.
* `enforce_safe_search` - (Boolean) Indicates whether only safe content must be returned for web, image, and video search.
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
