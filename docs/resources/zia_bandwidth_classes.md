---
subcategory: "Bandwidth Control"
layout: "zscaler"
page_title: "ZIA): bandwidth_classes"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-bandwidth-classes
  API documentation https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthClasses-post
  Adds a new bandwidth class

---
# zia_bandwidth_classes (Resource)

* [Official documentation](https://help.zscaler.com/zia/adding-bandwidth-classes)
* [API documentation](https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthClasses-post)

Use the **zia_bandwidth_classes** resource allows the creation and management of ZIA Bandwidth Control classes in the Zscaler Internet Access cloud or via the API.

## Example Usage - By Name

```hcl

resource "zia_bandwidth_classes" "this" {
    name = "Gen_AI_Classes"
    web_applications = [
        "ACADEMICGPT",
        "AD_CREATIVES",
        "AGENTGPT",
        "AI_ART_GENERATOR",
        "AI_CHAT_ROBOT",
        "AI_COPYWRITING_TREASURE",
        "AI_FOR_SEO",
        "ONE_MIN_AI"
    ]
    urls = ["chatgpt.com", "chatgpt1.com"]
    url_categories = [
        "AI_ML_APPS",
        "GENERAL_AI_ML"
    ]
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Optional) Name of the bandwidth class

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `urls` - (List of strings) The URLs included in the bandwidth class. You can include multiple entries.
* `url_categories` - (List of strings) The URL categories to add to the bandwidth class. Use the data source `zia_url_categories` to retrieve the available categories or visit the [Help Portal](https://help.zscaler.com/zia/url-categories#/urlCategories-get)
* `web_applications` - (List of strings) The web conferencing applications included in the bandwidth class. Use the data source `zia_cloud_applications` to retrieve the available applications or visit the [Help Portal](https://help.zscaler.com/zia/cloud-applications#/cloudApplications/policy-get)

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_bandwidth_classes** can be imported by using `<CLASS_ID>` or `<CLASS_NAME>` as the import ID.

For example:

```shell
terraform import zia_bandwidth_classes.this <class_id>
```

```shell
terraform import zia_bandwidth_classes.this <"class_name">
```
