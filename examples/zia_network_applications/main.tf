terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_network_application" "example"{
    id = "DICT"
}

output "zia_network_application"{
    value = data.zia_network_application.example
}


data "zia_network_application_groups_lite" "example"{
    name = "Microsoft Office365"
}

output "zia_network_application_groups_lite"{
    value = data.zia_network_application_groups_lite.example
}