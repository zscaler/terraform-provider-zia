data "zia_firewall_filtering_time_window" "work_hours"{
    name = "Work hours"
}

output "zia_firewall_filtering_time_window_work_hours"{
    value = data.zia_firewall_filtering_time_window.work_hours
}

data "zia_firewall_filtering_time_window" "weekends"{
    name = "Weekends"
}

output "zia_firewall_filtering_time_window_weekends"{
    value = data.zia_firewall_filtering_time_window.weekends
}

data "zia_firewall_filtering_time_window" "off_hours"{
    name = "Off hours"
}

output "zia_firewall_filtering_time_window_off_hours"{
    value = data.zia_firewall_filtering_time_window.off_hours
}
