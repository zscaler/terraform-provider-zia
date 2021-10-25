data "zia_firewall_filtering_network_application_groups" "example"{
    name = "Microsoft Office365"
}

output "zia_firewall_filtering_network_application_groups"{
    value = data.zia_firewall_filtering_network_application_groups.example
}