resource "zia_firewall_filtering_rule" "example" {
    name = "Example"
    description = "Example"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    time_windows {
        id = [ data.zia_firewall_filtering_time_window.work_hours.id ]
    }
}

data "zia_firewall_filtering_time_window" "work_hours"{
    name = "Work hours"
}
