data "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "1.1.1.1"
}

output "zia_traffic_forwarding_static_ip_available"{
    value = data.zia_traffic_forwarding_static_ip.example
}