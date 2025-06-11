---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA): forwarding_control_proxies"
description: |-
  Get information about firewall IPS Control policy rule.

---
# Data Source: zia_forwarding_control_proxies

Use the **zia_forwarding_control_proxies** data source to get information about a third-party proxy service available in the Zscaler Internet Access.

## Example Usage - Retrieve By Name

```hcl
data "zia_forwarding_control_proxies" "this" {
    name = "Proxy01"
}
```

## Example Usage - Retrieve By ID

```hcl
data "zia_forwarding_control_proxies" "this" {
    id = "18492370"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Proxy name for the third-party proxy services
* `id` - (Optional) Unique identifier for the third-party proxy services

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Additional notes or information
* `address` - (String) The IP address or the FQDN of the third-party proxy service
* `port` - (integer) The port number on which the third-party proxy service listens to the requests forwarded from Zscaler
* `insert_xau_header` - (Boolean) Flag indicating whether X-Authenticated-User header is added by the proxy. Enable to automatically insert authenticated user ID to the HTTP header, X-Authenticated-User.
* `base64_encode_xau_header` - (Boolean) Flag indicating whether the added X-Authenticated-User header is Base64 encoded. When enabled, the user ID is encoded using the Base64 encoding method.

* `type` - (String) Gateway type. Returned values: `PROXYCHAIN`, `ZIA`, `ECSELF`

* `cert` - (Set of Objects) The root certificate used by the third-party proxy to perform SSL inspection. This root certificate is used by Zscaler to validate the SSL leaf certificates signed by the upstream proxy. The required root certificate appears in this drop-down list only if it is uploaded from the Administration > Root Certificates page.
      - `id` - (Integer) Identifier that uniquely identifies the certificate
