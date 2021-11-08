resource "zia_firewall_filtering_rule" "example" {
    name = "Example"
    description = "Example"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    groups {
        id = [ data.zia_group_management.normal_internet.id ]
    }
}

data "zia_group_management" "normal_internet" {
 name = "Normal_Internet"
}