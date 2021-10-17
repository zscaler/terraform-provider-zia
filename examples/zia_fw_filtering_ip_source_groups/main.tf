terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

resource "zia_firewall_filtering_ip_source_groups" "example1"{
    name = "example1"
    description = "example1"
    ip_addresses = ["1.2.3.4", "1.2.3.5", "1.2.3.6"]
}

resource "zia_firewall_filtering_ip_source_groups" "example2"{
    name = "example2"
    description = "example2"
    ip_addresses = ["1.2.3.7", "1.2.3.8", "1.2.3.9"]
}


