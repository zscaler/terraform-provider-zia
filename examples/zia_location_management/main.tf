terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

resource "zia_location_management" "toronto"{
    name = "SGIO-IPSEC-Toronto"
    description = "Created with Terraform"
    ip_addresses = [ zia_traffic_forwarding_static_ip.example.ip_address ]
    surrogate_ip = true
    idle_time_in_minutes = 480
    auth_required = true
    surrogate_ip_enforced_for_known_browsers = true
    surrogate_refresh_time_in_minutes = 240
    surrogate_refresh_time_unit = "MINUTE"
    vpn_credentials {
       id = zia_traffic_forwarding_vpn_credentials.example.vpn_credental_id
       type = zia_traffic_forwarding_vpn_credentials.example.type
    }
}

resource "zia_traffic_forwarding_vpn_credentials" "example"{
    type = "UFQDN"
    fqdn = "sjc-1-373@securitygeek.io"
    comments = "created automatically"
    pre_shared_key = "newPassword123!"
}

resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "50.98.112.169"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = false
}