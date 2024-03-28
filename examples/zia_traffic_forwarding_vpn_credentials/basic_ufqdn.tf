######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_traffic_forwarding_vpn_credentials" "example"{
    type = "UFQDN"
    fqdn = "sjc-1-37@acme.com"
    comments = "Created via Terraform"
    pre_shared_key = "***************"
}

output "zia_traffic_forwarding_vpn_credentials"{
    value = zia_traffic_forwarding_vpn_credentials.example
}