terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_activation_status" "example"{
}

output "zia_activation_status"{
    value = data.zia_activation_status.example
}

/*
resource "zia_activation_status" "example"{
    status = "ACTIVE"
}

output "zia_activation_status"{
    value = zia_activation_status.example
}
*/