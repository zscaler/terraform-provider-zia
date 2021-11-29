data "zia_traffic_forwarding_vpn_credentials" "example"{
    fqdn = "sjc-1-37@acme.com"
}

output "zia_vpn_credentials_sjc-1-37"{
    value = data.zia_traffic_forwarding_vpn_credentials.example
}

data "zia_traffic_forwarding_vpn_credentials" "example"{
    ip_address = "1.1.1.1"
}

output "zia_vpn_credentials_sjc-1-37"{
    value = data.zia_traffic_forwarding_vpn_credentials.example
}