---
subcategory: "NSS Server"
layout: "zscaler"
page_title: "ZIA: zia_nss_server"
description: |-
  Creates and manages ZIA NSS Servers.
---

# Resource: zia_nss_server

The **zia_nss_server** resource allows the creation and management of NSS Server Objects in the Zscaler Internet Access cloud or via the API.
See [Adding NSS Servers](https://help.zscaler.com/zia/adding-nss-servers) for more details.

## Example Usage - Type NSS_FOR_FIREWALL

```hcl
resource "zia_nss_server" "this" {
    name = "NSSServer01"
    status = "ENABLED"
    type = "NSS_FOR_FIREWALL"
}
```

## Example Usage - Type NSS_FOR_WEB

resource "zia_nss_server" "this" {
    name = "NSSServer01"
    status = "ENABLED"
    type = "NSS_FOR_WEB"
}

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the devices to be created.
* `type` - (String) Whether you are creating an NSS for web logs or firewall logs. Returned Values:  `NSS_FOR_WEB`, `NSS_FOR_FIREWALL`

### Optional

* `description` - (String) The rule label description.
* `status` - (String) Enables or disables the status of the NSS server. Returned Values: `ENABLED`, `DISABLED`, `DISABLED_BY_SERVICE_PROVIDER`, `NOT_PROVISIONED_IN_SERVICE_PROVIDER`, `IN_TRIAL`
* `state` - (String) The health of the NSS server. Returned Values:  `UNHEALTHY`, `HEALTHY`, `UNKNOWN`
* `icap_svr_id` - (integer) The ICAP server ID

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_nss_server** can be imported by using `<NSS_ID>` or `<NSS_NAME>` as the import ID.

For example:

```shell
terraform import zia_nss_server.example <nss_id>
```

or

```shell
terraform import zia_nss_server.example <nss_name>
```
