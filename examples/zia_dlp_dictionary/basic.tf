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

output "zia_dlp_dictionaries_example"{
    value = zia_dlp_dictionaries.example
}
