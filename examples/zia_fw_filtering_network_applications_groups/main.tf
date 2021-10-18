terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

/*
data "zia_firewall_filtering_network_application_groups" "example"{
    name = "Microsoft Office365"
}

output "zia_firewall_filtering_network_application_groups"{
    value = data.zia_firewall_filtering_network_application_groups.example
}
*/

resource "zia_firewall_filtering_network_application_groups" "example"{
    name = "Example"
    description = "Example"
    network_applications = [ "APNS", "APPSTORE", "DICT", "EPM", "GARP", "ICLOUD", "IOS_OTA_UPDATE"]
}

output "zia_firewall_filtering_network_application_groups" {
    value = zia_firewall_filtering_network_application_groups.example
}