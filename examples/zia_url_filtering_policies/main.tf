terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}


resource "zia_url_filtering_rules" "block_innapropriate_contents"{
    name = "Block Inappropriate Contents"
    description = "Block all inappropriate content for all users."
    order = 1
    state = "ENABLED"
    locations {
        id = 
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
/*
data "zia_url_filtering_policies" "example"{
    //name = "Block Inappropriate Content"
    name = "Isolate - Allow Paste"
}

output "zia_url_filtering_policies"{
    value = data.zia_url_filtering_policies.example
}
*/
