terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_department_management" "engineering" {
 name = "Engineering"
}

output "zia_department_management" {
  value = data.zia_department_management.engineering
}