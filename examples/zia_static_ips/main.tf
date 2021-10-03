terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "50.98.112.169"
    routable_ip = true
    comment = "Created with Terraform"
}

output "zia_traffic_forwarding_static_ip"{
    value = zia_traffic_forwarding_static_ip.example
}

/*
data "zia_traffic_forwarding_static_ip" "example"{
    ip_address = "96.53.93.170"
}

output "zia_traffic_forwarding_static_ip"{
    value = data.zia_traffic_forwarding_static_ip.example
}
*/