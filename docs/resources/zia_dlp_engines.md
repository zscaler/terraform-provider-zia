---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_engines"
description: |-
  Official documentation https://help.zscaler.com/zia/about-dlp-engines
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/dlpEngines-get
  Adds a new custom DLP engine
---

# zia_dlp_engines (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-dlp-engines)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/dlpEngines-get)

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

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_dlp_engines** can be imported by using `<ENGINE_ID>` or `<ENGINE_NAME>` as the import ID.

For example:

```shell
terraform import zia_dlp_engines.example <engine_id>
```

or

```shell
terraform import zia_dlp_engines.example <engine_name>
```
