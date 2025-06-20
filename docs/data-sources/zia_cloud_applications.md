---
subcategory: "Cloud Applications"
layout: "zscaler"
page_title: "ZIA: cloud_applications"
description: |-
  Official documentation https://help.zscaler.com/zia/cloud-applications#/cloudApplications/sslPolicy-get
  API documentation https://help.zscaler.com/zia/cloud-applications#/cloudApplications/sslPolicy-get
  Retrieves a list of Predefined and User Defined Cloud Applications associated with the DLP rules, Cloud App Control rules, Advanced Settings, Bandwidth Classes, File Type Control and SSL Inspection rules.
---

# zia_cloud_applications (Data Source)

* [Official documentation](https://help.zscaler.com/zia/cloud-applications#/cloudApplications/sslPolicy-get)
* [API documentation](https://help.zscaler.com/zia/cloud-applications#/cloudApplications/sslPolicy-get)

Use the **zia_cloud_applications** data source to Retrieves a list of Predefined and User Defined Cloud Applications associated with the DLP rules, Cloud App Control rules, Advanced Settings, Bandwidth Classes, File Type Control and SSL Inspection rules. The returned information can be associated with the attribute `cloud_applications` on supported rules.

```hcl
# Retrieves a list of Predefined and User Defined Cloud Applications associated with the DLP rules, Cloud App Control rules, Advanced Settings, Bandwidth Classes, and File Type Control rules.
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_policy"
}

output "zia_cloud_applications" {
  value = data.zia_cloud_applications.this
}

# Retrieves and Filter Cloud Application by application category
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_policy"
  app_class   = ["AI_ML"]
}

output "app_ids" {
  value = [for app in data.zia_cloud_applications.this.applications : app["app"]]
}

# Retrieves specific application by name and category
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_ssl_policy"
  app_class = ["SOCIAL_NETWORKING"]
  app_name = "Nebenan"
}

output "zia_cloud_applications" {
    value = data.zia_cloud_applications.this
}


# Retrieves a list of Predefined and User Defined Cloud Applications associated with the SSL Inspection rules.
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_ssl_policy"
}

output "zia_cloud_applications" {
  value = data.zia_cloud_applications.this
}

#Retrieves and Filter Cloud Application associated with a SSL inspection rule by application category
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_ssl_policy"
  app_class   = ["AI_ML"]
}

output "app_ids" {
  value = [for app in data.zia_cloud_applications.this.applications : app["app"]]
}

# Retrieves specific application associated with a SSL inspection rule by name and category
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_ssl_policy"
  app_class = ["SOCIAL_NETWORKING"]
  app_name = "Nebenan"
}

output "zia_cloud_applications" {
    value = data.zia_cloud_applications.this
}
```

## Argument Reference

The following arguments are supported:

* `cloud_application_policy` - (Required) Retrieves a list of Predefined and User Defined Cloud Applications associated with the DLP rules, Cloud App Control rules, Advanced Settings, Bandwidth Classes, and File Type Control rules.
* `cloud_application_ssl_policy` - (Optional) Retrieves a list of Predefined and User Defined Cloud Applications associated with the SSL Inspection rules.

    **NOTE** You may use `cloud_application_policy` or `cloud_application_ssl_policy` but not both at the same time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `app_class` - (Set of Strings) Filter application by application category
* `app` - (String) Application enum constant
* `app_name` - (String) Cloud application name
* `parent` - (String) pplication category enum constant
* `parent_name` - (String) Name of the cloud application category

## Cloud Application Category App Class Matrix

**Note**: Refer to this matrix when configuring types vs actions for each specific rules

|             App Class                       |
|:-------------------------------------------:|
|---------------------------------------------|
|               `WEBMAIL`                     |
|           `SOCIAL_NETWORKING`               |
|              `STREAMING`                    |
|                 `P2P`                       |
|            `INSTANT_MESSAGING`              |
|               `WEB_SEARCH`                  |
|            `GENERAL_BROWSING`               |
|               `ADMINISTRATION`              |
|               `ENTERPRISE_COLLABORATION`    |
|               `BUSINESS_PRODUCTIVITY`       |
|               `SALES_AND_MARKETING`         |
|               `SYSTEM_AND_DEVELOPMENT`      |
|               `CONSUMER`                    |
|               `FILE_SHARE`                  |
|               `HOSTING_PROVIDER`            |
|               `IT_SERVICES`                 |
|               `DNS_OVER_HTTPS`              |
|               `HUMAN_RESOURCES`             |
|               `LEGAL`                       |
|               `HEALTH_CARE`                 |
|               `FINANCE`                     |
|               `CUSTOM_CAPP`                 |
|               `AI_ML`                       |
|---------------------------------------------|
