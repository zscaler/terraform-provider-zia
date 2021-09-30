terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}


data "zia_network_service_groups" "example"{
    name = "Corporate Remote Access Services"
}

output "zia_network_service_groups"{
    value = data.zia_network_service_groups.example
}


/*
data "zia_network_services" "example"{
    name = "DNS"
}

output "zia_network_services"{
    value = data.zia_network_services.example
}
*/