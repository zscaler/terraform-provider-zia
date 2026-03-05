---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: email_profile"
description: |-
    Official documentation https://help.zscaler.com/zia/about-dlp-notification-templates
    API documentation https://help.zscaler.com/legacy-apis/data-loss-prevention-0#/emailRecipientProfile-get
    Retrieves email recipient profiles configured for the organization.
---

# zia_email_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-dlp-notification-templates)
* [API documentation](https://help.zscaler.com/legacy-apis/data-loss-prevention-0#/emailRecipientProfile-get)

Use the **zia_email_profile** data source to get information about email recipient profiles configured for the organization in the Zscaler Internet Access cloud.

## Example Usage - Retrieve by Name

```hcl
data "zia_email_profile" "this" {
    name = "EmailProfile01"
}
```

## Example Usage - Retrieve by ID

```hcl
data "zia_email_profile" "this" {
    id = 1254674585
}
```

## Argument Reference

The following arguments are supported:

### Required

At least one of the following must be provided:

* `id` - (Integer) Unique identifier for the email recipient profile. Used to look up a single profile when provided.
* `name` - (String) Email recipient profile name. Used to search for a profile by name when provided.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Integer) The unique identifier for the email recipient profile.
* `name` - (String) The name of the email recipient profile.
* `description` - (String) The description of the email recipient profile.
* `emails` - (Set of String) The list of recipient email addresses.
