---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_notification_templates"
description: |-
  Official documentation https://help.zscaler.com/zia/about-dlp-notification-templates
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/dlpNotificationTemplates-get
  Get information about DLP Notification Templates.
---

# zia_dlp_notification_templates (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-dlp-notification-templates)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/dlpNotificationTemplates-get)

Use the **zia_dlp_notification_templates** data source to get information about a ZIA DLP Notification Templates in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP Template by name
data "zia_dlp_notification_templates" "example"{
    name = "DLP Auditor Template Test"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP policy rule name.

### Optional

* `id` - (Optional) The unique identifier for a DLP notification template.
* `plain_text_message` - (Optional) The template for the plain text UTF-8 message body that must be displayed in the DLP notification email.
* `html_message` - (Optional) The template for the HTML message body that must be displayed in the DLP notification email.
* `subject` - (Optional) The Subject line that is displayed within the DLP notification email.
* `attach_content` - (Optional) If set to true, the content that is violation is attached to the DLP notification email.
* `tls_enabled` - (Optional) If set to true, the content that is violation is attached to the DLP notification email.
