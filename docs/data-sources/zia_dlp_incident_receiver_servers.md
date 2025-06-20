---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_incident_receiver_servers"
description: |-
  Official documentation https://help.zscaler.com/zia/about-zscaler-incident-receiver
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/incidentReceiverServers-get
  Get information about ZIA DLP Incident Receiver Servers.
---

# zia_dlp_incident_receiver_servers (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-zscaler-incident-receiver)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/incidentReceiverServers-get)

Use the **zia_dlp_incident_receiver_servers** data source to get information about a ZIA DLP Incident Receiver Server in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP Incident Receiver Server by name
data "zia_dlp_incident_receiver_servers" "this" {
  name = "ZS_Incident_Receiver"
}
```

```hcl
# Retrieve a DLP Incident Receiver Server by ID
data "zia_dlp_incident_receiver_servers" "this"{
    id = 1234567890
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP Incident Receiver Server name as configured by the admin.

### Optional

* `id` - (Number) The unique identifier for the DLP engine.
* `url` - (String) The Incident Receiver server URL.
* `status` - (String) The status of the Incident Receiver. The returned values are:
  * ``ENABLED``
  * ``DISABLED``
  * ``DISABLED_BY_SERVICE_PROVIDER``
  * ``NOT_PROVISIONED_IN_SERVICE_PROVIDER``

* `flags` - (Number) The Incident Receiver server flag.
