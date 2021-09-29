terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_ip_destination_groups" "example"{
    name = "Trusted-Destinations"
}

output "zia_ip_destination_groups"{
    value = data.zia_ip_destination_groups.example
}