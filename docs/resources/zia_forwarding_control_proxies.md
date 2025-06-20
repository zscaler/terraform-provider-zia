---
subcategory: "Forwarding Control Policy"
layout: "zscaler"
page_title: "ZIA: zia_forwarding_control_proxies"
description: |-
  Official documentation https://help.zscaler.com/zia/about-third-party-proxies
  API documentation https://help.zscaler.com/zia/forwarding-control-policy#/proxies-get
  Creates and manages ZIA forwarding control proxies for third-party proxy services.
---

# zia_forwarding_control_proxies (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-third-party-proxies)
* [API documentation](https://help.zscaler.com/zia/forwarding-control-policy#/proxies-get)

Use the **zia_forwarding_control_proxies** resource allows the creation and management of ZIA forwarding control Proxies for third-party proxy services integration between Zscaler Internet Access and Zscaler Private Access. This resource can then be associated with a ZIA Forwarding Control Rule.

## Example Usage - No Certificate

```hcl
resource "zia_forwarding_control_proxies" "this" {
  name  = "Proxy01_Terraform"
  description = "Proxy01_Terraform"
  type = "PROXYCHAIN"
  address = "192.168.1.150"
  port = 5000
  insert_xau_header = true
  base64_encode_xau_header = true
}
```

## Example Usage - With Certificate

```hcl
resource "zia_forwarding_control_proxies" "this" {
  name  = "Proxy01_Terraform"
  description = "Proxy01_Terraform"
  type = "PROXYCHAIN"
  address = "192.168.1.150"
  port = 5000
  insert_xau_header = true
  base64_encode_xau_header = true
  cert {
    id = 18492369
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` (Required) Proxy name for the third-party proxy services
* `address` - (Required) The IP address or the FQDN of the third-party proxy service
* `port` - (Required) The port number on which the third-party proxy service listens to the requests forwarded from Zscaler
* `type` - (Required) Gateway type. Supported values: `PROXYCHAIN`, `ZIA`, `ECSELF`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Additional notes or information
* `insert_xau_header` - (Boolean) Flag indicating whether X-Authenticated-User header is added by the proxy. Enable to automatically insert authenticated user ID to the HTTP header, X-Authenticated-User.
* `base64_encode_xau_header` - (Boolean) Flag indicating whether the added X-Authenticated-User header is Base64 encoded. When enabled, the user ID is encoded using the Base64 encoding method.

* `cert` - (Set of Objects) The root certificate used by the third-party proxy to perform SSL inspection. This root certificate is used by Zscaler to validate the SSL leaf certificates signed by the upstream proxy. The required root certificate appears in this drop-down list only if it is uploaded from the Administration > Root Certificates page.
      - `id` - (Integer) Identifier that uniquely identifies the certificate

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_forwarding_control_proxies** can be imported by using `<PROXY_ID>` or `<PROXY_NAME>` as the import ID.

For example:

```shell
terraform import zia_forwarding_control_proxies.example <proxy_id>
```

or

```shell
terraform import zia_forwarding_control_proxies.example <proxy_name>
```
