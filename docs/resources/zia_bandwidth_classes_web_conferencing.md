---
subcategory: "Bandwidth Control"
layout: "zscaler"
page_title: "ZIA): bandwidth_classes_web_conferencing"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-bandwidth-classes
  API documentation https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthClasses-post
  Adds a new web conferencing bandwidth class

---
# zia_bandwidth_classes_web_conferencing (Resource)

* [Official documentation](https://help.zscaler.com/zia/adding-bandwidth-classes)
* [API documentation](https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthClasses-post)

Use the **zia_bandwidth_classes_web_conferencing** resource allows the creation and management of ZIA web conferencing bandwidth classes in the Zscaler Internet Access cloud or via the API.

## Example Usage - BANDWIDTH_CAT_WEBCONF

```hcl

resource "zia_bandwidth_classes_web_conferencing" "this" {
    name = "BANDWIDTH_CAT_WEBCONF"
    type = "BANDWIDTH_CAT_WEBCONF"
    applications = ["WEBEX", "GOTOMEETING", "LIVEMEETING", "INTERCALL", "CONNECT"]
}
```

## Example Usage - BANDWIDTH_CAT_VOIP

```hcl

resource "zia_bandwidth_classes_web_conferencing" "this" {
    name = "BANDWIDTH_CAT_VOIP"
    type = "BANDWIDTH_CAT_VOIP"
    applications = ["SKYPE"]
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Optional) Name of the bandwidth class

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `type` - (String) The application type for which the bandwidth class is configured. Supported values: `BANDWIDTH_CAT_WEBCONF` or `BANDWIDTH_CAT_VOIP`

* `applications` - (List of strings) The applications included in the bandwidth class

  * Supported Values for `BANDWIDTH_CAT_WEBCONF`: `WEBEX`, `GOTOMEETING`, `LIVEMEETING`, `INTERCALL`, `CONNECT`
  * Supported Values for `BANDWIDTH_CAT_VOIP`: `SKYPE`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_bandwidth_classes_web_conferencing** can be imported by using `<CLASS_ID>` or `<CLASS_NAME>` as the import ID.

For example:

```shell
terraform import zia_bandwidth_classes_web_conferencing.this <class_id>
```

```shell
terraform import zia_bandwidth_classes_web_conferencing.this <"class_name">
```
