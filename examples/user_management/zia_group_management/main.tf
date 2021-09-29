terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_group_management" "devops" {
 name = "DevOps"
}

output "zia_group_management" {
  value = data.zia_group_management.devops
}