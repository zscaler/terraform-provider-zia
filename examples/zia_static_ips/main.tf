terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_traffic_forwarding_static_ip" "example"{
    //id = 61125
    ip_address = "96.53.93.170"
}

output "zia_traffic_forwarding_static_ip"{
    value = data.zia_traffic_forwarding_static_ip.example
}