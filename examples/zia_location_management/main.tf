terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_location_management" "vancouver"{
    name = "SGIO-IPSEC-Vancouver"
}

output "zia_location_management"{
    value = data.zia_location_management.vancouver
}