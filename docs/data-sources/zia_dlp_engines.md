---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_engines"
description: |-
  Official documentation https://help.zscaler.com/zia/about-dlp-engines
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/dlpEngines-get
  Get information about ZIA DLP Engines.
---

# zia_dlp_engines (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-dlp-engines)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/dlpEngines-get)

Use the **zia_dlp_engines** data source to get information about a ZIA DLP Engines in the Zscaler Internet Access cloud or via the API.

## Example Usage - Retrieve Custom DLP Engine by name

```hcl
#
data "zia_dlp_engines" "this"{
    name = "Example"
}
```

## Example Usage - Retrieve Custom DLP Engine by ID

```hcl
data "zia_dlp_engines" "this"{
    id = 1234567890
}
```

## Example Usage - Retrieve Predefined DLP Engine by Name

```hcl
data "zia_dlp_engines" "this"{
    predefined_engine_name = "PCI"
}

data "zia_dlp_engines" "this"{
    predefined_engine_name = "EXTERNAL"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP engine name as configured by the admin. This attribute is required in POST and PUT requests for custom DLP engines.

* `predefined_engine_name` - (String) To search for predefined DLP Engines use this attribute.

### Optional

* `id` - (Number) The unique identifier for the DLP engine.
* `predefined_engine_name` - (String) The name of the predefined DLP engine.
* `engine_expression` - (String) The boolean logical operator in which various DLP dictionaries are combined within a DLP engine's expression.
* `custom_dlp_engine` - (Bool) Indicates whether this is a custom DLP engine. If this value is set to true, the engine is custom.
* `description` - (String) The DLP engine's description.
