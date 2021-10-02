terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_gre_virtual_ip_address_list" "example"{
    //source_ip = "50.98.112.169"
    datacenter = "YVR1"
}

output "zia_gre_virtual_ip_address_list"{
    value = data.zia_gre_virtual_ip_address_list.example
}