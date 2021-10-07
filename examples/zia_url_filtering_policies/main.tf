terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}


resource "zia_url_filtering_rules" "example"{
    name = "Example"
    order = 2
    protocols = ["ANY_RULE"]
    state = "ENABLED"
    request_methods = ["OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "CONNECT", "OTHER"]
    description = "Example"
    action = "ALLOW"
    url_categories = ["FINANCE", "OTHER_BUSINESS_AND_ECONOMY", "CORPORATE_MARKETING"]
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
