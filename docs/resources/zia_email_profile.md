---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: email_profile"
description: |-
    Official documentation https://help.zscaler.com/zia/about-dlp-notification-templates
    API documentation https://help.zscaler.com/legacy-apis/data-loss-prevention-0#/emailRecipientProfile-post
    Creates and manages ZIA email recipient profiles.
---

# zia_email_profile (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-dlp-notification-templates)
* [API documentation](https://help.zscaler.com/legacy-apis/data-loss-prevention-0#/emailRecipientProfile-post)

The **zia_email_profile** resource allows the creation and management of email recipient profiles in the Zscaler Internet Access cloud. Email recipient profiles define sets of email addresses that can be referenced in DLP rules and notification configurations.

~> NOTE: This an Early Access feature.

## Example Usage

```hcl
resource "zia_email_profile" "this" {
    name        = "EmailProfile01"
    description = "Email recipient profile for DLP notifications"
    emails      = ["admin@example.com", "security@example.com"]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (String) The name of the email recipient profile.

### Optional

* `description` - (String) The description of the email recipient profile.
* `emails` - (Set of String) The list of recipient email addresses.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The unique identifier for the email recipient profile (Terraform internal).
* `email_profile_id` - (Integer) The unique identifier for the email recipient profile (API).

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_email_profile** can be imported by using `<EMAIL PROFILE ID>` or `<EMAIL PROFILE NAME>` as the import ID.

For example:

```shell
terraform import zia_email_profile.example <email_profile_id>
```

or

```shell
terraform import zia_email_profile.example <email_profile_name>
```
