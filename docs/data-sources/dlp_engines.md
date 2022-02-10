---
subcategory: "DLP Engines"
layout: "zia"
page_title: "ZIA: dlp_engines"
description: |-
  Retrieve ZIA DLP Engines details.
  
---

# zia_dlp_engines (Data Source)

The **zia_dlp_engines** data source provides details about a specific DLP engine options available in the Zscaler Internet Access.

```hcl
# Retrieve ZIA DLP Engines
data "zia_dlp_engines" "example"{
    name = "Credit Cards"
}

output "zia_dlp_engines"{
    value = data.zia_dlp_engines.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) DLP dictionary name

## Attribute Reference (Read-Only)

The following attributes are supported:

* `id` - (int) The unique identifier for the DLP engine.
* `name` - (string) The DLP engine name as configured by the admin. This attribute is required in POST and PUT requests for custom DLP engines.
* `predefined_engine_name` - (String) The name of the predefined DLP engine.
* `engine_expression` - (String) The boolean logical operator in which various DLP dictionaries are combined within a DLP engine's expression.
`custom_dlp_engine` - (bool) Indicates whether this is a custom DLP engine. If this value is set to true, the engine is custom
* `description` - (String) The DLP engine's description.
