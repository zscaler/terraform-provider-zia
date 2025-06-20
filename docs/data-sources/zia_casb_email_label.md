---
subcategory: "SaaS Security API"
layout: "zscaler"
page_title: "ZIA): casb_email_label"
description: |-
  Official documentation https://help.zscaler.com/zia/about-email-labels
  API documentation https://help.zscaler.com/zia/saas-security-api#/casbEmailLabel/lite-get
  Retrieves the email labels generated for the SaaS Security API policies in a user's email account

---
# zia_casb_email_label (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-email-labels)
* [API documentation](https://help.zscaler.com/zia/saas-security-api#/casbEmailLabel/lite-get)

Use the **zia_casb_email_label** data source to get information about email labels generated for the SaaS Security API policies in a user's email account

## Example Usage - By Name

```hcl

data "zia_casb_email_label" "this" {
  name = "EmailLabel01"
}
```

## Example Usage - By ID

```hcl

data "zia_casb_email_label" "this" {
  id = 154658
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Int) SaaS Security API email label ID
* `name` - (String) SaaS Security API email label name

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `label_desc` - (String) Description of the email label
* `label_color` - (String) Color to apply to the email label
  * `CASB_EMAIL_LABEL_COLOR_RED`,
  * `CASB_EMAIL_LABEL_COLOR_YELLOW`,
  * `CASB_EMAIL_LABEL_COLOR_GREEN`,
  * `CASB_EMAIL_LABEL_COLOR_ORANGE`,
  * `CASB_EMAIL_LABEL_COLOR_BLUE`,
  * `CASB_EMAIL_LABEL_COLOR_PURPLE`
* `label_deleted` - (Boolean) A Boolean value that indicates whether or not the email label is deleted
