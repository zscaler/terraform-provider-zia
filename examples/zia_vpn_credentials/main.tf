terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

resource "zia_traffic_forwarding_vpn_credentials" "example"{
    type = "UFQDN"
    fqdn = "sjc-1-37@securitygeek.io"
    comments = "created automatically"
    pre_shared_key = "newPassword123!"
}

/*
data "zia_vpn_credentials" "example"{
    fqdn = "vpn@securitygeek.io"
}

output "zia_vpn_credentials"{
    value = data.zia_vpn_credentials.example
}
*/

