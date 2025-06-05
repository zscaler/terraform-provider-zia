resource "zia_traffic_forwarding_gre_tunnel" "example" {
  source_ip = zia_traffic_forwarding_static_ip.example_1.ip_address
  comment   = "GRE Tunnel Created with Terraform"
  depends_on = [ zia_traffic_forwarding_static_ip.example ]
  within_country = true
  country_code = "CA"
  ip_unnumbered = false
}

output "zia_traffic_forwarding_gre_tunnel_example1" {
  value = zia_traffic_forwarding_gre_tunnel.example1
}

// Create static ip addresses
resource "zia_traffic_forwarding_static_ip" "example_1"{
    ip_address =  "1.1.1.1"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = true
    latitude = -23.548670
    longitude = -46.638248
}