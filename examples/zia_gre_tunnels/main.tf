terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_traffic_forwarding_gre_tunnels" "example"{
    id = 61125
}

output "zia_traffic_forwarding_gre_tunnels"{
    value = data.zia_traffic_forwarding_gre_tunnels.example
}