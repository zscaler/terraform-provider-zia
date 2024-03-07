---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_dictionaries"
description: |-
  Creates and manages ZIA DLP dictionaries.
---

# Resource: zia_dlp_dictionaries

The **zia_dlp_dictionaries** resource allows the creation and management of ZIA DLP dictionaries in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
resource "zia_dlp_dictionaries" "example"{
    name = "Your Dictionary Name"
    description = "Your Description"
    phrases {
        action = "PHRASE_COUNT_TYPE_ALL"
        phrase = "YourPhrase"
    }
    custom_phrase_match_type = "MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY"
    patterns {
        action = "PATTERN_COUNT_TYPE_UNIQUE"
        pattern = "YourPattern"
    }
    dictionary_type = "PATTERNS_AND_PHRASES"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP dictionary's name
* `dictionary_type` - (Required) The DLP dictionary type. The following values are supported:
  * `PATTERNS_AND_PHRASES`
  * `EXACT_DATA_MATCH`
  * `INDEXED_DATA_MATCH`
* `custom_phrase_match_type` - (Required) The DLP custom phrase match type. Supported values are:
  * `MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY`
  * `MATCH_ANY_CUSTOM_PHRASE_PATTERN_DICTIONARY`
  Note: This attribute should only be set when the dictionary_type is set to ``PATTERNS_AND_PHRASES``

* `phrases` - (Required) List containing the phrases used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries. Required when `dictionary_type` is `PATTERNS_AND_PHRASES`
  * `action` - (Required) The action applied to a DLP dictionary using patterns. The following values are supported:
    * `PATTERN_COUNT_TYPE_ALL`
    * `PATTERN_COUNT_TYPE_UNIQUE`
  * `phrase` - (Required) DLP dictionary phrase

* `patterns` - (Required) List containing the patterns used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries. Required when `dictionary_type` is `PATTERNS_AND_PHRASES`
  * `action` - (Required) The action applied to a DLP dictionary using patterns. The following values are supported:
    * `PATTERN_COUNT_TYPE_ALL`
    * `PATTERN_COUNT_TYPE_UNIQUE`
  * `pattern` - (Required) DLP dictionary pattern

### Optional

* `confidence_threshold` - (Optional) The DLP confidence threshold. The following values are supported:
  * `CONFIDENCE_LEVEL_LOW`
  * `CONFIDENCE_LEVEL_MEDIUM`
  * `CONFIDENCE_LEVEL_HIGH`

* `threshold_type` - (Optional) The DLP threshold type. The following values are supported:
  * `VIOLATION_COUNT_ONLY`
  * `CONFIDENCE_SCORE_ONLY`
  * `VIOLATION_AND_CONFIDENCE`

* `threshold_type` - (Optional) The DLP threshold type. The following values are supported:
  * `VIOLATION_COUNT_ONLY`
  * `CONFIDENCE_SCORE_ONLY`
  * `VIOLATION_AND_CONFIDENCE`

* `exact_data_match_details` - (Optional) Exact Data Match (EDM) related information for custom DLP dictionaries.
  * `dictionary_edm_mapping_id` - (Optional) The unique identifier for the EDM mapping.
  * `schema_id` - (Optional) The unique identifier for the EDM template (or schema).
  * `primary_field` - (Optional) The EDM template's primary field.
  * `secondary_fields` - (Optional) The EDM template's secondary fields.
  * `secondary_field_match_on` - (Optional) The EDM secondary field to match on.
        - `"MATCHON_NONE"`
        - `"MATCHON_ANY_1"`
        - `"MATCHON_ANY_2"`
        - `"MATCHON_ANY_3"`
        - `"MATCHON_ANY_4"`
        - `"MATCHON_ANY_5"`
        - `"MATCHON_ANY_6"`
        - `"MATCHON_ANY_7"`
        - `"MATCHON_ANY_8"`
        - `"MATCHON_ANY_9"`
        - `"MATCHON_ANY_10"`
        - `"MATCHON_ANY_11"`
        - `"MATCHON_ANY_12"`
        - `"MATCHON_ANY_13"`
        - `"MATCHON_ANY_14"`
        - `"MATCHON_ANY_15"`
        - `"MATCHON_ALL"`

* `idm_profile_match_accuracy` - (Optional) List of Indexed Document Match (IDM) profiles and their corresponding match accuracy for custom DLP dictionaries.
  * `adp_idm_profile` - (Optional) The IDM template reference.
  * `match_accuracy` - (Optional) The IDM template match accuracy.
        - `"LOW"`
        - `"MEDIUM"`
        - `"HEAVY"`

* `proximity` - (Optional) The DLP dictionary proximity length.
* `ignore_exact_match_idm_dict` - (Optional) Indicates whether to exclude documents that are a 100% match to already-indexed documents from triggering an Indexed Document Match (IDM) Dictionary.
* `include_bin_numbers` - (Optional) A true value denotes that the specified Bank Identification Number (BIN) values are included in the Credit Cards dictionary. A false value denotes that the specified BIN values are excluded from the Credit Cards dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
* `bin_numbers` - (Optional) The list of Bank Identification Number (BIN) values that are included or excluded from the Credit Cards dictionary. BIN values can be specified only for Diners Club, Mastercard, RuPay, and Visa cards. Up to 512 BIN values can be configured in a dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
* `dict_template_id` - (Optional) ID of the predefined dictionary (original source dictionary) that is used for cloning. This field is applicable only to cloned dictionaries. Only a limited set of identification-based predefined dictionaries (e.g., Credit Cards, Social Security Numbers, National Identification Numbers, etc.) can be cloned. Up to 4 clones can be created from a predefined dictionary.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_dlp_dictionaries** can be imported by using `<DICTIONARY ID>` or `<DICTIONARY_NAME>` as the import ID.

For example:

```shell
terraform import zia_dlp_dictionaries.example <dictionary_id>
```

or

```shell
terraform import zia_dlp_dictionaries.example <dictionary_name>
```
