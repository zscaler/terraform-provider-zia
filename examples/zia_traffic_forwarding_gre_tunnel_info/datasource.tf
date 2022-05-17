data "zia_traffic_forwarding_gre_tunnel_info" "example1" {
  ip_address = "50.98.112.169"
}

output "zia_traffic_forwarding_gre_tunnel_info_example1" {
  value = data.zia_traffic_forwarding_gre_tunnel_info.example1
}

data "zia_traffic_forwarding_gre_tunnel_info" "example2" {
  ip_address = "187.22.113.134"
}

output "zia_traffic_forwarding_gre_tunnel_info_example2" {
  value = data.zia_traffic_forwarding_gre_tunnel_info.example2
}

data "zia_traffic_forwarding_gre_tunnel_info" "example3" {
  ip_address = "61.68.118.237"
}

output "zia_traffic_forwarding_gre_tunnel_info_example3" {
  value = data.zia_traffic_forwarding_gre_tunnel_info.example3
}

# :warning: `ip_address` is the public IP address associated with the GRE Tunnel