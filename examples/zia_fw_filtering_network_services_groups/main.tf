terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

resource "zia_firewall_filtering_network_service_groups" "example"{
    name = "example"
    description = "example"
    services {
        id = 773995
    }
}

