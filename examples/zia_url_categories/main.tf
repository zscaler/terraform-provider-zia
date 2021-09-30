terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}


data "zia_url_categories" "example"{
    //id = "SOCIAL_NETWORKING"
    //id = "CUSTOM_10"
    configured_name = "Custom_Category"
}

output "zia_url_categories"{
    value = data.zia_url_categories.example
}
