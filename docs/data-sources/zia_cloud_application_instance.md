---
subcategory: "Cloud App Control Policy"
layout: "zscaler"
page_title: "ZIA: cloud_application_instance"
description: |-
  Official documentation https://help.zscaler.com/zia/about-cloud-application-instances
  API documentation https://help.zscaler.com/zia/cloud-app-control-policy#/cloudApplicationInstances-get
  Retrieves the list of cloud application instances
---

# zia_cloud_application_instance (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-cloud-application-instances)
* [API documentation](https://help.zscaler.com/zia/cloud-app-control-policy#/cloudApplicationInstances-get)

Use the **zia_cloud_application_instance** data source to get information about cloud application instances in the Zscaler Internet Access cloud or via the API.

## Example Usage - By Name

```hcl
data "zia_cloud_application_instance" "this"{
    name = "SharePointOnline"
}
```

## Example Usage - By ID

```hcl
data "zia_cloud_application_instance" "this"{
    id = "11743520"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Optional) Name of the cloud application instance
* `id` - (Optional) Unique identifier for the cloud application instance

### Optional Filtering Options

* `instance_type` - (String) Type of the cloud application instance
* `include_deleted` - (Boolean) Include Deleted tenants
* `modified_at` - (String) Timestamp of when the cloud application instance was last modified
* `last_modified_by` - (Number)  The admin that modified the cloud application instance last
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity

    * `instance_identifiers` - (List) List of identifiers for the cloud application instance
      * `instance_id` - (Number) Unique identifier for the cloud application instance
      * `instance_identifier` - (String) Unique identifying string for the instance
      * `instance_identifier_name` - (String) Unique identifying string for the instance
      * `identifier_type` - (Number) Unique Type of the cloud application instance
      * `modified_at` - (String) Timestamp of when the cloud application instance was last modified
      * `last_modified_by` - (Number)  The admin that modified the cloud application instance last
      * `id` - (Number) Identifier that uniquely identifies an entity
      * `name` - (String) The configured name of the entity
