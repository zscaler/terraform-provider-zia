terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

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
    name_l10n_tag = false
    dictionary_type = "PATTERNS_AND_PHRASES"
}

output "zia_dlp_dictionaries_example"{
    value = zia_dlp_dictionaries.example
}

/*
data "zia_dlp_dictionaries" "example1"{
    name = "SALESFORCE_REPORT_LEAKAGE"
}

output "zia_dlp_dictionaries_example1"{
    value = data.zia_dlp_dictionaries.example1
}

data "zia_dlp_dictionaries" "example2"{
    name = "TIN_LEAKAGE"
}

output "zia_dlp_dictionaries_example2"{
    value = data.zia_dlp_dictionaries.example2
}

data "zia_dlp_dictionaries" "example3"{
    name = "PESEL_LEAKAGE"
}

output "zia_dlp_dictionaries_example3"{
    value = data.zia_dlp_dictionaries.example3
}

data "zia_dlp_dictionaries" "example4"{
    name = "AHV_LEAKAGE"
}

output "zia_dlp_dictionaries_example4"{
    value = data.zia_dlp_dictionaries.example4
}
*/