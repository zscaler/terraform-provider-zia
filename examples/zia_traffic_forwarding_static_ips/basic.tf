resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "1.1.1.1"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = true
    latitude = -36.848461
    longitude = 174.763336
}

output "zia_traffic_forwarding_static_ip"{
    value = zia_traffic_forwarding_static_ip.example
}


