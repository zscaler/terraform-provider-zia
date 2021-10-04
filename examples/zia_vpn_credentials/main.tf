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
output "zia_traffic_forwarding_vpn_credentials"{
    value = zia_traffic_forwarding_vpn_credentials.example
    // sensitive = true
}
*/


data "zia_traffic_forwarding_vpn_credentials" "example"{
    fqdn = "sjc-1-37@securitygeek.io"
}

output "zia_vpn_credentials_sjc-1-37"{
    value = data.zia_traffic_forwarding_vpn_credentials.example
    sensitive = true
}

