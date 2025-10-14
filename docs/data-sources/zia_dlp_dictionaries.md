---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_dictionaries"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-custom-dlp-dictionary
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/dlpDictionaries-post
  Get information about ZIA DLP Dictionaries.
---

# zia_dlp_dictionaries (Data Source)

* [Official documentation](https://help.zscaler.com/zia/adding-custom-dlp-dictionary)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/dlpDictionaries-post)

Use the **zia_dlp_dictionaries** data source to get information about a DLP dictionary option available in the Zscaler Internet Access.

```hcl
# Retrieve a DLP Dictionary by name
data "zia_dlp_dictionaries" "example"{
    name = "SALESFORCE_REPORT_LEAKAGE"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) DLP dictionary name
* `id` - (Optional) Unique identifier for the DLP dictionary

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `confidence_threshold` - (String) he DLP confidence threshold. [`CONFIDENCE_LEVEL_LOW`, `CONFIDENCE_LEVEL_MEDIUM` `CONFIDENCE_LEVEL_HIGH` ]
* `custom_phrase_match_type` - (String) The DLP custom phrase match type. [ `MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY`, `MATCH_ANY_CUSTOM_PHRASE_PATTERN_DICTIONARY` ]
* `dictionary_type` - (String) The DLP dictionary type. The cloud service API only supports custom DLP dictionaries that are using the `PATTERNS_AND_PHRASES` type.

* `confidence_level_for_predefined_dict` - (Optional) The DLP confidence threshold for predefined dictionaries. The following values are supported:
  * `CONFIDENCE_LEVEL_LOW`
  * `CONFIDENCE_LEVEL_MEDIUM`
  * `CONFIDENCE_LEVEL_HIGH`

* `name_l10n_tag` - (Boolean) Indicates whether the name is localized or not. This is always set to True for predefined DLP dictionaries.

* `hierarchical_identifiers` - (Optional) List of hierarchical identifiers for the DLP dictionary. Supported only for the following Identifiers: `ASPP_LEAKAGE`, `CRED_LEAKAGE`, `EUIBAN_LEAKAGE`, `PPEU_LEAKAGE`, `USDL_LEAKAGE`.

`phrases` - (List of Object) List containing the patterns used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries

* `action` - (String) The action applied to a DLP dictionary using patterns
* `pattern` - (String) DLP dictionary pattern

`patterns` - (List of Object) List containing the patterns used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries

* `action` - (String) The action applied to a DLP dictionary using patterns
* `pattern` - (String) DLP dictionary pattern
* `ignore_exact_match_idm_dict` - (Boolean) Indicates whether to exclude documents that are a 100% match to already-indexed documents from triggering an Indexed Document Match (IDM) Dictionary.
* `include_bin_numbers` - (Boolean) A true value denotes that the specified Bank Identification Number (BIN) values are included in the Credit Cards dictionary. A false value denotes that the specified BIN values are excluded from the Credit Cards dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
* `bin_numbers` - (Boolean) The list of Bank Identification Number (BIN) values that are included or excluded from the Credit Cards dictionary. BIN values can be specified only for Diners Club, Mastercard, RuPay, and Visa cards. Up to 512 BIN values can be configured in a dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
* `dict_template_id` - (Number) ID of the predefined dictionary (original source dictionary) that is used for cloning. This field is applicable only to cloned dictionaries. Only a limited set of identification-based predefined dictionaries (e.g., Credit Cards, Social Security Numbers, National Identification Numbers, etc.) can be cloned. Up to 4 clones can be created from a predefined dictionary.
* `predefined_clone` - (Boolean) This field is set to true if the dictionary is cloned from a predefined dictionary. Otherwise, it is set to false.
* `custom` - (Boolean) This value is set to true for custom DLP dictionaries.
* `proximity` - (Optional, Integer) The DLP dictionary proximity length that defines how close a high confidence phrase must be to an instance of the pattern (that the dictionary detects) to count as a match. Supported values between `0` and `10000`
* `proximity_enabled_for_custom_dictionary` - (Optional, Boolean) A Boolean constant that indicates if proximity length is enabled or disabled for a custom DLP dictionary.
* `proximity_length_enabled` - (Optional, Boolean) A Boolean constant that indicates whether the proximity length option is supported for a DLP dictionary or not. A true value indicates that the proximity length option is supported, whereas a false value indicates that it is not supported.
