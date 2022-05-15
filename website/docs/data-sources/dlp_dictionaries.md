---
subcategory: "DLP Dictionaries"
layout: "zia"
page_title: "ZIA: dlp_dictionaries"
description: |-
  Retrieve ZIA DLP Dictionaries details.
  
---

# zia_dlp_dictionaries (Data Source)

The **zia_dlp_dictionaries** data source provides details about a specific DLP dictionary option available in the Zscaler Internet Access.

```hcl
# Retrieve ZIA Location Group
data "zia_dlp_dictionaries" "example"{
    name = "SALESFORCE_REPORT_LEAKAGE"
}

output "zia_dlp_dictionaries_example"{
    value = data.zia_dlp_dictionaries.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) DLP dictionary name
* `id` - (Optional) Unique identifier for the DLP dictionary

## Attribute Reference

The following attributes are supported:

* `confidence_threshold` - (String) he DLP confidence threshold. [`CONFIDENCE_LEVEL_LOW`, `CONFIDENCE_LEVEL_MEDIUM` `CONFIDENCE_LEVEL_HIGH` ]
* `custom_phrase_match_type` - (String) The DLP custom phrase match type. [ `MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY`, `MATCH_ANY_CUSTOM_PHRASE_PATTERN_DICTIONARY` ]
* `dictionary_type` - (String) The DLP dictionary type. The cloud service API only supports custom DLP dictionaries that are using the `PATTERNS_AND_PHRASES` type.
* `name_l10n_tag` - (Boolean) Indicates whether the name is localized or not. This is always set to True for predefined DLP dictionaries.

`phrases` - (List of Object) List containing the patterns used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries

* `action` - (String) The action applied to a DLP dictionary using patterns
* `pattern` - (String) DLP dictionary pattern

`patterns` - (List of Object) List containing the patterns used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries

* `action` - (String) The action applied to a DLP dictionary using patterns
* `pattern` - (String) DLP dictionary pattern
