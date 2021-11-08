resource "zia_firewall_filtering_rule" "example" {
    name = "Example"
    description = "Example"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    nw_application_groups {
        id = [ data.zia_firewall_filtering_network_application_groups.example.id ]
    }
}

data "zia_firewall_filtering_network_application_groups" "example"{
    name = "Microsoft Office365"
}