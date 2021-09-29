terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_ip_source_groups" "example"{
    name = "Trusted_Sources"
}

output "zia_ip_source_groups"{
    value = data.zia_ip_source_groups.example
}