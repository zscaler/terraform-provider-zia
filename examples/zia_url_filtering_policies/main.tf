terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}


data "zia_url_filtering_policies" "example"{
    //name = "Block Inappropriate Content"
    name = "Isolate - Allow Paste"
}

output "zia_url_filtering_policies"{
    value = data.zia_url_filtering_policies.example
}
