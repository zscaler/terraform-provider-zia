---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_engines"
description: |-
  Get information about ZIA DLP Engines.
---

# Data Source: zia_dlp_engines

Use the **zia_dlp_engines** data source to get information about a ZIA DLP Engines in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP Engine by name
data "zia_dlp_engines" "this"{
    name = "Example"
}
```

```hcl
# Retrieve a DLP Engine by ID
data "zia_dlp_engines" "this"{
    id = 1234567890
}
```

```hcl
# Retrieve a Predefined DLP Engine
data "zia_dlp_engines" "this"{
    predefined = "EXTERNAL"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP engine name as configured by the admin. This attribute is required in POST and PUT requests for custom DLP engines.

### Optional

* `id` - (Number) The unique identifier for the DLP engine.
* `predefined_engine_name` - (String) The name of the predefined DLP engine.
* `engine_expression` - (String) The boolean logical operator in which various DLP dictionaries are combined within a DLP engine's expression.
* `custom_dlp_engine` - (Bool) Indicates whether this is a custom DLP engine. If this value is set to true, the engine is custom.
* `description` - (String) The DLP engine's description.
