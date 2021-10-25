data "zia_firewall_filtering_destination_groups" "example" {
    name = "example"
}

output "zia_firewall_filtering_destination_groups_example" {
    value = data.zia_firewall_filtering_destination_groups.example
}