terraform {
    required_providers {
        zia = {
            version = "1.0.4"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_devices" "model"{
    model = "VMware"
}

output "zia_devices"{
    value = data.zia_devices.model
}