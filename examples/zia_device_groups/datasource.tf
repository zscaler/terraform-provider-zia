terraform {
    required_providers {
        zia = {
            version = "1.0.4"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_device_groups" "ios"{
    name = "IOS"
}

output "zia_device_groups_ios"{
    value = data.zia_device_groups.ios
}

data "zia_device_groups" "android"{
    name = "Android"
}

output "zia_device_groups_android"{
    value = data.zia_device_groups.android
}

data "zia_device_groups" "windows"{
    name = "Windows"
}

output "zia_device_groups_windows"{
    value = data.zia_device_groups.windows
}

data "zia_device_groups" "mac"{
    name = "Mac"
}

output "zia_device_groups_mac"{
    value = data.zia_device_groups.mac
}

data "zia_device_groups" "linux"{
    name = "Linux"
}

output "zia_device_groups_linux"{
    value = data.zia_device_groups.linux
}

data "zia_device_groups" "no_client_connector"{
    name = "No Client Connector"
}

output "zia_device_groups_no_client_connector"{
    value = data.zia_device_groups.no_client_connector
}

data "zia_device_groups" "cloud_browser_isolation"{
    name = "Cloud Browser Isolation"
}

output "zia_device_groups_cloud_browser_isolation"{
    value = data.zia_device_groups.cloud_browser_isolation
}