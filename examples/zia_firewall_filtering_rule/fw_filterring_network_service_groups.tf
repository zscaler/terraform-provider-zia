resource "zia_firewall_filtering_rule" "example" {
    name = "Example"
    description = "Example"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    nw_service_groups {
        id = [ data.zia_firewall_filtering_network_service_groups.example.id ]
    }
}

data "zia_firewall_filtering_network_service_groups" "example"{
    name = "Corporate Custom SSH TCP_10022"
}