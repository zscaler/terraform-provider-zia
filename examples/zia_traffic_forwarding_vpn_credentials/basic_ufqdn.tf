resource "zia_traffic_forwarding_vpn_credentials" "example"{
    type = "UFQDN"
    fqdn = "sjc-1-37@acme.com"
    comments = "Created via Terraform"
    pre_shared_key = "newPassword123!"
}

output "zia_traffic_forwarding_vpn_credentials"{
    value = zia_traffic_forwarding_vpn_credentials.example
}