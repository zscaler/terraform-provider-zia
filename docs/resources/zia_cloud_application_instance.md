---
subcategory: "Cloud App Control Policy"
layout: "zscaler"
page_title: "ZIA: cloud_application_instance"
description: |-
  Official documentation https://help.zscaler.com/zia/about-cloud-application-instances
  API documentation https://help.zscaler.com/zia/cloud-app-control-policy#/cloudApplicationInstances-post
  Add a new cloud application instance
---

# zia_cloud_application_instance (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-cloud-application-instances)
* [API documentation](https://help.zscaler.com/zia/cloud-app-control-policy#/cloudApplicationInstances-post)

The **zia_cloud_application_instance** resource allows the creation and management of cloud application instance.

## Example Usage

```hcl
resource "zia_cloud_application_instance" "this" {
  name          = "SharePointOnline"
  instance_type = "SHAREPOINTONLINE"

  instance_identifiers {
    instance_identifier      = "acme.sharepoint.com"
    instance_identifier_name = "acme"
    identifier_type          = "URL"
  }

  instance_identifiers {
    instance_identifier      = "acme1.sharepoint.com"
    instance_identifier_name = "acme1"
    identifier_type          = "URL"
  }

  instance_identifiers {
    instance_identifier      = "acme2.sharepoint.com"
    instance_identifier_name = "acme2"
    identifier_type          = "URL"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Optional) Name of the cloud application instance
* `instance_type` - (String) Type of the cloud application instance

  * `instance_identifiers` - (List) List of identifiers for the cloud application instance
    * `instance_identifier` - (String) Unique identifying string for the instance
    * `instance_identifier_name` - (String) Unique identifying string for the instance
    * `identifier_type` - (Number) Unique Type of the cloud application instance. Supported Values: `URL`, `REFURL`, `KEYWORD`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_cloud_application_instance** can be imported by using `<INSTANCE_ID>` or `<INSTANCE_NAME>` as the import ID.

For example:

```shell
terraform import zia_cloud_application_instance.example <instance_id>
```

or

```shell
terraform import zia_cloud_application_instance.example <instance_name>
```
