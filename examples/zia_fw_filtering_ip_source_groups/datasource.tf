data "zia_firewall_filtering_ip_source_groups" "example" {
    name = "example"
}

output "zia_firewall_filtering_ip_source_groups_example" {
    value = data.zia_firewall_filtering_ip_source_groups.example
}