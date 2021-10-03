terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

resource "zia_ip_source_groups" "example"{
    name = "Example"
    description = "Example"
    ip_addresses = [ "100.100.100.1" ]

}

output "zia_ip_source_groups"{
    value = zia_ip_source_groups.example
}

/*
data "zia_ip_source_groups" "example"{
    name = "Trusted_Sources"
}

output "zia_ip_source_groups"{
    value = data.zia_ip_source_groups.example
}
*/