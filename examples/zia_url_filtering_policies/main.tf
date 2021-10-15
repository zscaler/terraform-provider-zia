terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

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
    order = 1
    state = "ENABLED"
    locations {
        id = [zia_location_management.toronto.id]
    }

    url_categories = [ "ADULT_SEX_EDUCATION",
                       "ADULT_THEMES",
                       "COMPUTER_HACKING",
                       "LINGERIE_BIKINI",
                       "NUDITY",
                       "OTHER_ADULT_MATERIAL",
                       "OTHER_DRUGS",
                        "OTHER_ILLEGAL_OR_QUESTIONABLE",
                       "QUESTIONABLE",
                       "PORNOGRAPHY",
                       "SEXUALITY"
                    ]
    protocols = ["ANY_RULE"]
    request_methods = ["OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "CONNECT", "OTHER"]
    action = "BLOCK"

}

data "zia_url_filtering_rules" "example"{
    //name = "Block Inappropriate Content"
    name = zia_url_filtering_rules.block_innapropriate_contents.name
}

output "zia_url_filtering_rules"{
    value = data.zia_url_filtering_rules.example
}

