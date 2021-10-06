terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}


resource "zia_activation_status" "example1"{
    status = "ACTIVE"
}

output "zia_activation_status_example1"{
    value = zia_activation_status.example1
}

data "zia_activation_status" "example2"{
}

output "zia_activation_status_example2"{
    value = data.zia_activation_status.example2
}



