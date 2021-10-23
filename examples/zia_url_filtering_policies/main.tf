terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

<<<<<<< HEAD
resource "zia_url_filtering_rules" "sao_paulo_guest_wifi" {
    name = "Sao Paulo Guest WiFi"
    description = "Sao Paulo Guest WiFi"
=======
resource "zia_location_management" "toronto"{
    name = "SGIO-IPSEC-Toronto"
    description = "Created with Terraform"
    ip_addresses = [ zia_traffic_forwarding_static_ip.example.ip_address ]
}

resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "50.98.112.169"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = false
}

resource "zia_url_filtering_rules" "block_innapropriate_contents"{
    name = "Block Inappropriate Contents"
    description = "Block all inappropriate content for all users."
    order = 2
>>>>>>> master
    state = "ENABLED"
    action = "BLOCK"
    order = 1
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
}
resource "zia_url_filtering_rules" "block_streaming" {
    name = "Block Streaming"
    description = "Block Video Streaming."
    state = "ENABLED"
    action = "BLOCK"
    order = 2
    url_categories = ["ANY"]
    protocols = ["ANY_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}

resource "zia_url_filtering_rules" "caution_for_gambling" {
    name = "Caution for Gambling"
    description = "Caution Marketing users going to gambling sites, but allow them after the caution. This is to support marketing ad campaigns on Gambling websites."
    state = "ENABLED"
    action = "CAUTION"
    order = 3
    url_categories = ["ANY"]
    protocols = ["ANY_RULE"]
    request_methods = [ "CONNECT", "GET", "HEAD" ]
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

