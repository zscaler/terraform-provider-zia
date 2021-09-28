terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_vpn_credentials" "example"{
    fqdn = "vpn@securitygeek.io"
}

output "zia_vpn_credentials"{
    value = data.zia_vpn_credentials.example
}