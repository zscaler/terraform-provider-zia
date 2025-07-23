---
subcategory: "Bandwidth Control"
layout: "zscaler"
page_title: "ZIA): bandwidth_classes"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-bandwidth-classes
  API documentation https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthClasses-get
  Retrieves a list of bandwidth classes for an organization

---
# zia_bandwidth_classes (Data Source)

* [Official documentation](https://help.zscaler.com/zia/adding-bandwidth-classes)
* [API documentation](https://help.zscaler.com/zia/bandwidth-control-classes#/bandwidthClasses-post)

Use the **zia_bandwidth_classes** Retrieves all the available bandwidth control classes.

## Example Usage - By Name

```hcl
data "zia_bandwidth_classes" "this" {
    name = "Gen_AI_Classes"
}
```

## Example Usage - By ID

```hcl
data "zia_bandwidth_classes" "this" {
    id = 13
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) System-generated identifier for the bandwidth class
* `name` - (Optional) Name of the bandwidth class

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `urls` - (List of strings) The URLs included in the bandwidth class. You can include multiple entries.
* `url_categories` - (List of strings) The URL categories to add to the bandwidth class. Use the data source `zia_url_categories` to retrieve the available categories or visit the [Help Portal](https://help.zscaler.com/zia/url-categories#/urlCategories-get)
* `web_applications` - (List of strings) The web conferencing applications included in the bandwidth class. Use the data source `zia_cloud_applications` to retrieve the available applications or visit the [Help Portal](https://help.zscaler.com/zia/cloud-applications#/cloudApplications/policy-get)
