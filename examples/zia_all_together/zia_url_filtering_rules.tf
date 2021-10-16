/*
resource "zia_url_filtering_rules" "isolate_allow_paste" {
    name = "Isolate - Allow Paste"
    description = "Isolate - Allow Paste"
    state = "ENABLED"
    action = "ISOLATE"
    cbi_profile_id = 25625083
    order = 1
    url_categories = ["MISCELLANEOUS_OR_UNKNOWN", "OTHER_MISCELLANEOUS"]
    protocols = ["HTTPS_RULE", "HTTP_RULE"]
    request_methods = [ "CONNECT", "GET", "HEAD", "TRACE"]
}
*/


resource "zia_url_filtering_rules" "mcas_block" {
    name = "MCAS Block"
    description = "MCAS Block"
    state = "ENABLED"
    action = "BLOCK"
    order = 1
    url_categories = [zia_url_categories.mcas_unsanctioned_apps.id]
    protocols = ["HTTPS_RULE", "HTTP_RULE", "HTTP_PROXY", "FTP_RULE", "FOHTTP_RULE", "SSL_RULE", "TUNNEL_RULE", "TUNNELSSL_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}

resource "zia_url_filtering_rules" "sao_paulo_guest_wifi" {
    name = "Sao Paulo Guest WiFi"
    description = "Sao Paulo Guest WiFi"
    state = "ENABLED"
    action = "BLOCK"
    order = 2
    url_categories = [
                        "ADULT_SEX_EDUCATION",
                        "ADULT_THEMES",
                        "ANONYMIZER",
                        "COMPUTER_HACKING",
                        "COPYRIGHT_INFRINGEMENT", 
                        "GAMBLING",
                        "K_12_SEX_EDUCATION",
                        "LINGERIE_BIKINI",
                        "MATURE_HUMOR",
                        "MILITANCY_HATE_AND_EXTREMISM",
                        "NUDITY",
                        "OTHER_ADULT_MATERIAL",
                        "OTHER_ILLEGAL_OR_QUESTIONABLE",
                        "OTHER_SECURITY",
                        "PORNOGRAPHY",
                        "PROFANITY",
                        "QUESTIONABLE",
                        "SEXUALITY",
                        "SOCIAL_NETWORKING",
                        "ADWARE_OR_SPYWARE",
                        "TASTELESS",
                        "VIOLENCE",
                        "WEAPONS_AND_BOMBS",
                    ]
    protocols = ["ANY_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
    locations {
        id = [ zia_location_management.br_sao_paulo_branch01_guest_wifi.id ]
    }
}

resource "zia_url_filtering_rules" "caution_for_gambling" {
    name = "Caution for Gambling"
    description = "Caution Marketing users going to gambling sites, but allow them after the caution. This is to support marketing ad campaigns on Gambling websites."
    state = "ENABLED"
    action = "CAUTION"
    order = 3
    url_categories = ["ANY"]
    protocols = ["ANY_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}


resource "zia_url_filtering_rules" "allow_all_other_traffic" {
    name = "Allow All Other Traffic"
    description = "Allow All Other Traffic"
    state = "ENABLED"
    action = "ALLOW"
    order = 4
    url_categories = ["ANY"]
    protocols = ["HTTPS_RULE", "HTTP_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}
