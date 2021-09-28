terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_public_node_vips" "example"{
    datacenter = "AKL1"
}

output "zia_public_node_vips"{
    value = data.zia_public_node_vips.example
}