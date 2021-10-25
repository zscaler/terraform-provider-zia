data "zia_firewall_filtering_network_service" "example" {
  name = "ICMP_ANY"
}

output "zia_firewall_filtering_network_service" {
  value = data.zia_firewall_filtering_network_service.example
}