---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_notification_templates"
description: |-
  Creates and manages ZIA DLP Notification Templates.
---

# Resource: zia_dlp_notification_templates

The **zia_dlp_notification_templates** resource allows the creation and management of ZIA DLP Notification Templates in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
resource "zia_dlp_notification_templates" "example" {
    name                = "DLP Auditor Template Test"
    subject             = "DLP Violation: ${TRANSACTION_ID} ${ENGINES}"
    attach_content      = true
    tls_enabled         = true
    html_message        = file("./index.html")
    plain_text_message = file("./dlp.txt")
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP policy rule name.
* `plain_text_message` - (Required) The template for the plain text UTF-8 message body that must be displayed in the DLP notification email.
* `html_message` - (Required) The template for the HTML message body that must be displayed in the DLP notification email.

### Optional

* `subject` - (Optional) The Subject line that is displayed within the DLP notification email.
* `attach_content` - (Optional) If set to true, the content that is violation is attached to the DLP notification email.
* `tls_enabled` - (Optional) If set to true, the content that is violation is attached to the DLP notification email.
