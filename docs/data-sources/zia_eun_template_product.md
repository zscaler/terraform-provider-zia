---
subcategory: "End User Notification"
layout: "zscaler"
page_title: "ZIA: eun_template_product"
description: |-
  Official documentation https://help.zscaler.com/zia/about-browser-eun-template
  API documentation https://help.zscaler.com/legacy-apis/end-user-notifications#/eunTemplate/{templateType}/product/{product}-get
  Retrieves a notification template for browser-based and Zscaler Client Connector-based notifications by policy type.
---

# zia_eun_template_product (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-browser-eun-template)
* [API documentation](https://help.zscaler.com/legacy-apis/end-user-notifications#/eunTemplate/{templateType}/product/{product}-get)

Use the **zia_eun_template_product** data source to retrieve a notification template for browser-based and Zscaler Client Connector-based notifications by policy type. The retrieved template can then be referenced from other rule-based resources such as `zia_cloud_app_control_rule`, `zia_file_type_control_rules`, and `zia_url_filtering_rules`.

The service returns the complete list of templates for the requested `template_type` and `product`. Because it does not offer a server-side search key, the template is selected locally by `id` or `name`. When neither `id` nor `name` is provided, the default template for the policy type is returned (falling back to the first template in the list if none is marked as default).

## Example Usage - By Name

```hcl
data "zia_eun_template_product" "this" {
  template_type = "ZCC"
  product       = "ENDPOINT_DLP"
  name          = "Default"
}
```

## Example Usage - By ID

```hcl
data "zia_eun_template_product" "this" {
  template_type = "BROWSER"
  product       = "URL"
  id            = 12345
}
```

## Example Usage - Default Template

```hcl
# When neither id nor name is set, the default template for the policy type is returned.
data "zia_eun_template_product" "this" {
  template_type = "ZCC"
  product       = "ENDPOINT_DLP"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `template_type` - (String) The type of notification template. Supported values are `ZCC` (Zscaler Client Connector-based) and `BROWSER` (browser-based).
* `product` - (String) The policy type associated with the notification template. Supported values are `INLINE`, `ENDPOINT_DLP`, `CLOUDAPP`, `URL`, `FILE_TYPE`, `FIREWALL`, `DNS`, and `IPS`.

### Optional

* `id` - (Integer) The unique identifier of the notification template to look up.
* `name` - (String) The name of the notification template to look up.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `channel` - (String) The channel associated with the notification template.
* `type` - (String) The notification template type as returned by the service.
* `caution_interval` - (String) The interval used for caution notifications.
* `default` - (Boolean) Indicates whether this is the default notification template for the policy type.
* `notification_details` - (List of String) The list of notification details associated with the template.
* `recommended_cloud_app` - (List) The recommended cloud application associated with the notification template.
  * `val` - (Integer) The identifier of the recommended cloud application.
  * `name` - (String) The name of the recommended cloud application.
  * `channel` - (String) The channel of the recommended cloud application.
  * `product` - (String) The product of the recommended cloud application.
  * `type` - (String) The type of the recommended cloud application.
  * `misc` - (String) The caution interval associated with the recommended cloud application.
  * `app_not_ready` - (Boolean) Indicates whether the recommended cloud application is not yet ready.
  * `under_migration` - (Boolean) Indicates whether the recommended cloud application is under migration.
  * `app_cat_modified` - (Boolean) Indicates whether the recommended cloud application category has been modified.
  * `deprecated` - (Boolean) Indicates whether the recommended cloud application is deprecated.
* `language_templates` - (List) The list of per-language notification messages associated with the template.
  * `language` - (String) The language of the notification message.
  * `allow_message` - (String) The message displayed when access is allowed.
  * `block_message` - (String) The message displayed when access is blocked.
  * `encrypt_message` - (String) The message displayed when content is encrypted.
  * `readonly_message` - (String) The message displayed when content is read-only.
  * `caution_message` - (String) The message displayed in the caution notification.
  * `redirect_response_message` - (String) The message displayed for redirect responses.
  * `default` - (Boolean) Indicates whether this is the default language template.
