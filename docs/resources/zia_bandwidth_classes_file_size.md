---
subcategory: "Bandwidth Control"
layout: "zscaler"
page_title: "ZIA): bandwidth_classes_file_size"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-bandwidth-classes
  API documentation https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthClasses-post
  Configure file size bandwidth class

---
# zia_bandwidth_classes_file_size (Resource)

* [Official documentation](https://help.zscaler.com/zia/adding-bandwidth-classes)
* [API documentation](https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthClasses-post)

Use the **zia_bandwidth_classes_file_size** resource allows the creation and management of ZIA file size bandwidth class in the Zscaler Internet Access cloud or via the API.

## Example Usage - BANDWIDTH_CAT_LARGE_FILE

```hcl

resource "zia_bandwidth_classes_file_size" "this1" {
    file_size = "FILE_5MB"
}
```

## Schema

### Required

The following arguments are supported:

* `file_size` - (String) The file size for a bandwidth class. Supported Values for `FILE_5MB`: `FILE_10MB`, `FILE_50MB`, `FILE_100MB`, `FILE_250MB`, `FILE_500MB`, `FILE_1GB`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_bandwidth_classes_file_size** can be imported by using `<CLASS_ID>` or `<CLASS_NAME>` as the import ID.

For example:

```shell
terraform import zia_bandwidth_classes_file_size.this <class_id>
```

```shell
terraform import zia_bandwidth_classes_file_size.this <"class_name">
```
