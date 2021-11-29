resource "zia_traffic_forwarding_vpn_credentials" "example"{
    type = "IP"
    fqdn = "1.1.1.1"
    comments = "Created via Terraform"
    pre_shared_key = "newPassword123!"
}

output "zia_traffic_forwarding_vpn_credentials"{
    value = zia_traffic_forwarding_vpn_credentials.example
}