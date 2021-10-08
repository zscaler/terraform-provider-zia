terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_gre_internal_ip_range_list" "example"{
    required_count = 10
}

output "zia_gre_internal_ip_range_list_example"{
    value = data.zia_gre_internal_ip_range_list.example
}