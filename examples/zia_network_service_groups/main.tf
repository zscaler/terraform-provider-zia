terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

/*
data "zia_network_service_groups_lite" "example"{
    name = "Test"
}

output "zia_network_service_groups_lite"{
    value = data.zia_network_service_groups_lite.example
}
*/

data "zia_network_services_lite" "example"{
    name = "DNS"
}

output "zia_network_services_lite"{
    value = data.zia_network_services_lite.example
}