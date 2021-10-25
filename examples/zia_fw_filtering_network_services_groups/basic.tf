resource "zia_firewall_filtering_network_service_groups" "example"{
    name = "example"
    description = "example"
    services {
        id = [
            data.zia_firewall_filtering_network_service.example1.id,
            data.zia_firewall_filtering_network_service.example2.id,
        ]
    }
}

output  "zia_firewall_filtering_network_service_groups"  {
    value = zia_firewall_filtering_network_service_groups.example
}