---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_icap_servers"
description: |-
  Official documentation https://help.zscaler.com/zia/about-icap-communication-between-zscaler-and-dlp-servers
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/icapServers/lite-get
  Gets a the list of DLP servers using ICAP
---

# zia_dlp_engines (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-icap-communication-between-zscaler-and-dlp-servers)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/icapServers/lite-get)

Use the **zia_dlp_engines** data source to get information about a the list of DLP servers using ICAP in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP ICAP Server by name
data "zia_dlp_icap_servers" "example"{
    name = "Example"
}
```

```hcl
# Retrieve a DLP ICAP Server by ID
data "zia_dlp_icap_servers" "example"{
    id = 1234567890
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP server name as configured by the admin.

### Optional

* `id` - (Number) The unique identifier for a DLP server.
* `name` - (String) The DLP server name.
* `url` - (String) The DLP server URL.
* `status` - (Bool) The DLP server status
