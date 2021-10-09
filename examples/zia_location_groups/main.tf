
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
data "zia_location_groups" "example"{
    name = "Corporate User Traffic Group"
}

output "zia_location_groups"{
    value = data.zia_location_groups.example
}
*/

data "zia_location_groups" "example"{
    name = "Home-Office-Vancouver"
}

output "zia_location_groups"{
    value = data.zia_location_groups.example
}

