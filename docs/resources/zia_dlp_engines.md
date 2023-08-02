---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_engines"
description: |-
  Get information about ZIA DLP Engines.
---

# Data Source: zia_dlp_engines

Use the **zia_dlp_engines** resource allows the creation and management of ZIA DLP Engines in the Zscaler Internet Access cloud or via the API.

⚠️ **WARNING:** "Before using the new ``zia_dlp_engines`` resource contact [Zscaler Support](https://help.zscaler.com/login-tickets)." and request the following API methods ``POST``, ``PUT``, and ``DELETE`` to be enabled for your organization.

## Example Usage

```hcl
# Retrieve a DLP Engine by name
resource "zia_dlp_engines" "this" {
  name = "Example"
  description = "Example"
  engine_expression = "((D63.S > 1))"
  custom_dlp_engine = true
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP engine name as configured by the admin. This attribute is required in POST and PUT requests for custom DLP engines.
* `predefined_engine_name` - (String) The name of the predefined DLP engine.
* `engine_expression` - (String) The boolean logical operator in which various DLP dictionaries are combined within a DLP engine's expression.
* `custom_dlp_engine` - (Bool) Indicates whether this is a custom DLP engine. If this value is set to true, the engine is custom.

### Optional

* `description` - (String) The DLP engine's description.
