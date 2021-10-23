terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}


data "zia_firewall_filtering_time_windows" "work_hours"{
    name = "Work hours"
}

output "zia_firewall_filtering_time_windows_work_hours"{
    value = data.zia_firewall_filtering_time_windows.work_hours
}

data "zia_firewall_filtering_time_windows" "weekends"{
    name = "Weekends"
}

output "zia_firewall_filtering_time_windows_weekends"{
    value = data.zia_firewall_filtering_time_windows.weekends
}

data "zia_firewall_filtering_time_windows" "off_hours"{
    name = "Off hours"
}

output "zia_firewall_filtering_time_windows_off_hours"{
    value = data.zia_firewall_filtering_time_windows.off_hours
}
