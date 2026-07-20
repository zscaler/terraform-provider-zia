---
subcategory: "End User Notification"
layout: "zscaler"
page_title: "ZIA: eun_user_confirmation_template_product"
description: |-
  Official documentation https://help.zscaler.com/zia/about-user-confirmation-notification-templates
  API documentation https://help.zscaler.com/legacy-apis/end-user-notifications#/userConfirmation/product/{product}-get
  Retrieves a user confirmation notification template by policy type.
---

# zia_eun_user_confirmation_template_product (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-user-confirmation-notification-templates)
* [API documentation](https://help.zscaler.com/legacy-apis/end-user-notifications#/userConfirmation/product/{product}-get)

Use the **zia_eun_user_confirmation_template_product** data source to retrieve a user confirmation notification template by policy type. User confirmation notifications are supported for Inline Web DLP and Endpoint DLP policies. The retrieved template can then be referenced from rule-based resources such as `zia_endpoint_dlp_rules`.

The service returns the complete list of user confirmation templates for the requested `product`. Because it does not offer a server-side search key or pagination, the template is selected locally by `id` or `name`. When neither `id` nor `name` is provided, the default template for the policy type is returned (falling back to the first template in the list if none is marked as default).

## Example Usage - By Name

```hcl
data "zia_eun_user_confirmation_template_product" "this" {
  product = "ENDPOINT_DLP"
  name    = "Default"
}
```

## Example Usage - By ID

```hcl
data "zia_eun_user_confirmation_template_product" "this" {
  product = "INLINE"
  id      = 12345
}
```

## Example Usage - Default Template

```hcl
# When neither id nor name is set, the default template for the policy type is returned.
data "zia_eun_user_confirmation_template_product" "this" {
  product = "ENDPOINT_DLP"
}
```

## Example Usage - Referenced by zia_endpoint_dlp_rules

```hcl
data "zia_eun_user_confirmation_template_product" "endpoint_dlp" {
  product = "ENDPOINT_DLP"
}

resource "zia_endpoint_dlp_rules" "this" {
  name                 = "Example_Endpoint_DLP_Rule"
  description          = "Example_Endpoint_DLP_Rule"
  order                = 1
  rank                 = 7
  action               = "CONFIRM"
  state                = "ENABLED"
  severity             = "RULE_SEVERITY_HIGH"
  data_transfer_method = "PRINTING"

  # Reference the resolved user confirmation template ID.
  uc_template_id = data.zia_eun_user_confirmation_template_product.endpoint_dlp.id
}
```

## Argument Reference

The following arguments are supported:

### Required

* `product` - (String) The policy type associated with the user confirmation notification template. Supported values are `INLINE`, `ENDPOINT_DLP`, `CLOUDAPP`, `URL`, `FILE_TYPE`, `FIREWALL`, `DNS`, and `IPS`. User confirmation notifications are supported for Inline Web DLP (`INLINE`) and Endpoint DLP (`ENDPOINT_DLP`) policies.

### Optional

* `id` - (Integer) The unique identifier of the user confirmation notification template to look up.
* `name` - (String) The name of the user confirmation notification template to look up.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `channel` - (String) The channel associated with the user confirmation notification template.
* `default` - (Boolean) Indicates whether this is the default user confirmation notification template for the policy type.
* `language_templates` - (List) The list of per-language user confirmation messages associated with the template.
  * `language` - (String) The language of the user confirmation message.
  * `message` - (String) The user confirmation message displayed for the language.
  * `default` - (Boolean) Indicates whether this is the default language template.
