data "zia_firewall_filtering_network_service_groups" "example"{
    name = "Corporate Custom SSH TCP_10022"
}

output "zia_firewall_filtering_network_service_groups" {
  value = data.zia_firewall_filtering_network_service_groups.example
}