---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_dictionaries"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-custom-dlp-dictionary
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/dlpDictionaries-post
  Creates and manages ZIA DLP dictionaries.
---

# zia_dlp_dictionaries (Resource)

* [Official documentation](https://help.zscaler.com/zia/adding-custom-dlp-dictionary)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/dlpDictionaries-post)

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

## Example Usage - With Hierarchical Identifiers (Clone Predefined Dictionary)

```hcl
data "zia_dlp_dictionary_predefined_identifiers" "this" {
  name = "EUIBAN_LEAKAGE"
}

data "zia_dlp_dictionaries" "this" {
  name = "EUIBAN_LEAKAGE"
}

resource "zia_dlp_dictionaries" "example"{
    name                     = "Example Dictionary Clone"
    description              = "Example Dictionary Clone"
    confidence_level_for_predefined_dict = "CONFIDENCE_LEVEL_MEDIUM"
    hierarchical_identifiers = [data.zia_dlp_dictionary_predefined_identifiers.this.predefined_identifiers]
    confidence_threshold     = "CONFIDENCE_LEVEL_HIGH"
    dict_template_id         = data.zia_dlp_dictionaries.this.id
    phrases {
        action = "PHRASE_COUNT_TYPE_ALL"
        phrase = "YourPhrase1"
    }
    custom_phrase_match_type = "MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY"
    dictionary_type          = "PATTERNS_AND_PHRASES"
}
```

## Example Usage - With Exact Data Match (EDM)

```hcl
data "zia_dlp_edm_schema" "this"{
    project_name = "EDM_TEMPLATE01"
}

resource "zia_dlp_dictionaries" "dlp_dictionaries" {
  name        = "edm_dic_tf"
  description = "edm dictionary"
  dictionary_type = "EXACT_DATA_MATCH"
  custom = true

  exact_data_match_details {
    schema_id = data.zia_dlp_edm_schema.this.schema_id
    primary_fields             = [3]
    secondary_fields          = [1,2]
    secondary_field_match_on  = "MATCHON_ALL"

  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required, String) The DLP dictionary's name
* `dictionary_type` - (Required, String) The DLP dictionary type. The following values are supported:
  * `PATTERNS_AND_PHRASES`
  * `EXACT_DATA_MATCH`
  * `INDEXED_DATA_MATCH`
  * `MIP_TAG`

### Optional

* `description` - (Optional, String) The description of the DLP dictionary
* `custom` - (Optional, Boolean) Indicates whether the DLP dictionary is custom. Defaults to true.
* `confidence_level_for_predefined_dict` - (Optional, String) The DLP confidence threshold for predefined dictionaries. The following values are supported:
  * `CONFIDENCE_LEVEL_LOW`
  * `CONFIDENCE_LEVEL_MEDIUM`
  * `CONFIDENCE_LEVEL_HIGH`

* `confidence_threshold` - (Optional, String) The DLP confidence threshold. The following values are supported:
  * `CONFIDENCE_LEVEL_LOW`
  * `CONFIDENCE_LEVEL_MEDIUM`
  * `CONFIDENCE_LEVEL_HIGH`

* `phrases` - (Optional, Set) List containing the phrases used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries. Maximum 256 items.
  * `action` - (Optional, String) The action applied to a DLP dictionary using phrases. The following values are supported:
    * `PHRASE_COUNT_TYPE_UNIQUE`
    * `PHRASE_COUNT_TYPE_ALL`
  * `phrase` - (Optional, String) DLP dictionary phrase (0-128 characters)

* `custom_phrase_match_type` - (Optional, String) The DLP custom phrase match type. Supported values are:
  * `MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY`
  * `MATCH_ANY_CUSTOM_PHRASE_PATTERN_DICTIONARY`
  Note: This attribute should only be set when the dictionary_type is set to `PATTERNS_AND_PHRASES`

* `patterns` - (Optional, Set) List containing the patterns used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries. Maximum 8 items.
  * `action` - (Optional, String) The action applied to a DLP dictionary using patterns. The following values are supported:
    * `PATTERN_COUNT_TYPE_ALL`
    * `PATTERN_COUNT_TYPE_UNIQUE`
  * `pattern` - (Optional, String) DLP dictionary pattern (0-256 characters)

* `hierarchical_identifiers` - (Optional, Set of Strings) List of hierarchical identifiers for the DLP dictionary.

* `exact_data_match_details` - (Optional, Set) Exact Data Match (EDM) related information for custom DLP dictionaries.
  * `dictionary_edm_mapping_id` - (Optional, Integer) The unique identifier for the EDM mapping.
  * `schema_id` - (Optional, Integer) The unique identifier for the EDM template (or schema). To retrieve the EDM (Exact Data Match) information, use the data source: `zia_dlp_edm_schema`

  * `primary_fields` - (Optional, List of Integers) The EDM template's primary fields.
  * `secondary_fields` - (Optional, List of Integers) The EDM template's secondary fields.
  * `secondary_field_match_on` - (Optional, String) The EDM secondary field to match on. The following values are supported:
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
        - `"MATCHON_ATLEAST_1"`

* `idm_profile_match_accuracy` - (Optional, Set) List of Indexed Document Match (IDM) profiles and their corresponding match accuracy for custom DLP dictionaries.
  * `adp_idm_profile` - (Optional, Set) The IDM template reference. To retrieve the IDM (Index Data Match) information, use the data source: `zia_dlp_idm_profiles`
    * `id` - (Optional, Integer) Identifier that uniquely identifies an entity
    * `extensions` - (Optional, Map of Strings) Extensions map
  * `match_accuracy` - (Optional, String) The IDM template match accuracy. The following values are supported:
    * `"LOW"`
    * `"MEDIUM"`
    * `"HEAVY"`

* `proximity` - (Optional, Integer) The DLP dictionary proximity length.
* `ignore_exact_match_idm_dict` - (Optional, Boolean) Indicates whether to exclude documents that are a 100% match to already-indexed documents from triggering an Indexed Document Match (IDM) Dictionary.
* `include_bin_numbers` - (Optional, Boolean) A true value denotes that the specified Bank Identification Number (BIN) values are included in the Credit Cards dictionary. A false value denotes that the specified BIN values are excluded from the Credit Cards dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
* `bin_numbers` - (Optional, List of Integers) The list of Bank Identification Number (BIN) values that are included or excluded from the Credit Cards dictionary. BIN values can be specified only for Diners Club, Mastercard, RuPay, and Visa cards. Up to 512 BIN values can be configured in a dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
* `dict_template_id` - (Optional, Integer) ID of the predefined dictionary (original source dictionary) that is used for cloning. This field is applicable only to cloned dictionaries. Only a limited set of identification-based predefined dictionaries (e.g., Credit Cards, Social Security Numbers, National Identification Numbers, etc.) can be cloned. Up to 4 clones can be created from a predefined dictionary.

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
