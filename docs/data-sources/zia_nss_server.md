---
subcategory: "NSS Server"
layout: "zscaler"
page_title: "ZIA: nss_server"
description: |-
  Official documentation https://help.zscaler.com/zia/about-nss-servers
  API documentation https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssServers-get
  Get information about NSS Server details.
---

# zia_nss_server (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-nss-servers)
* [API documentation](https://help.zscaler.com/zia/cloud-nanolog-streaming-service-nss#/nssServers-get)

Use the **zia_nss_server** data source to get information about a nss server resource in the Zscaler Internet Access cloud or via the API.
See [Adding NSS Servers](https://help.zscaler.com/zia/adding-nss-servers) for more details.

## Example Usage - Retrieve by Name

```hcl
data "zia_nss_server" "this" {
    name = "NSSServer01"
}
```

## Example Usage - Retrieve by ID

```hcl
data "zia_nss_server" "this" {
    id = "5445585"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the nss server to be exported.
* `id` - (String) System-generated identifier of the NSS server based on the software platform

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - (String) Enables or disables the status of the NSS server. Returned Values: `ENABLED`, `DISABLED`, `DISABLED_BY_SERVICE_PROVIDER`, `NOT_PROVISIONED_IN_SERVICE_PROVIDER`, `IN_TRIAL`
* `state` - (String) The health of the NSS server. Returned Values:  `UNHEALTHY`, `HEALTHY`, `UNKNOWN`
* `type` - (String) Whether you are creating an NSS for web logs or firewall logs. Returned Values:  `NSS_FOR_WEB`, `NSS_FOR_FIREWALL`
* `icap_svr_id` - (integer) The ICAP server ID
