terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

resource "zia_firewall_filtering_destination_groups" "example1"{
    name = "example1"
    description = "example1"
    type = "DSTN_IP"
    addresses = ["1.2.3.4", "1.2.3.5", "1.2.3.6"]
    countries = ["COUNTRY_CA"]
}

resource "zia_firewall_filtering_destination_groups" "example2"{
    name = "example2"
    description = "example2"
    type = "DSTN_IP"
    addresses = ["1.2.3.7", "1.2.3.8", "1.2.3.9"]
    countries = ["COUNTRY_US"]
}


