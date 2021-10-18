terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_firewall_filtering_network_application" "example1"{
    id = "DICT"
    locale="en-US"
}

output "zia_firewall_filtering_network_application1"{
    value = data.zia_firewall_filtering_network_application.example1
}


// data "zia_firewall_filtering_network_application" "example2"{
//     name = "Microsoft Office365"
// }

// output "zia_firewall_filtering_network_application2"{
//     value = data.zia_firewall_filtering_network_application.example2
// }