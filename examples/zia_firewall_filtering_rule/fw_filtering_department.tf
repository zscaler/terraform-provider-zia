resource "zia_firewall_filtering_rule" "example" {
    name = "Example"
    description = "Example"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    departments {
        id = [ data.zia_department_management.engineering.id ]
    }
}

data "zia_department_management" "engineering" {
 name = "Engineering"
}