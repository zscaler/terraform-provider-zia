######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_traffic_forwarding_vpn_credentials" "example"{
    type = "IP"
    ip_address = "1.1.1.1"
    comments = "Created via Terraform"
    pre_shared_key = "***************"
}

output "zia_traffic_forwarding_vpn_credentials"{
    value = zia_traffic_forwarding_vpn_credentials.example
}