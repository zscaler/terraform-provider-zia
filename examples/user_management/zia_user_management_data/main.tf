terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

resource "zia_user_management" "gary_sands" {
 name = "Gary Sands"
 email = "gary.sands@securitygeek.io"
 password = "Wg210680*"
 groups {
   id = data.zia_group_management.normal_internet.id
   name = data.zia_group_management.normal_internet.name
  }
 department {
   id = data.zia_department_management.engineering.id
   name = data.zia_department_management.engineering.name
  }

}

output "zia_user_management" {
  value = zia_user_management.gary_sands
}

data "zia_group_management" "normal_internet" {
 name = "Normal_Internet"
}

data "zia_department_management" "engineering" {
 name = "Engineering"
}

/*
data "zia_group_management" "engineering" {
 name = "Engineering"
}

output "zia_user_management_adam_name" {
  value = data.zia_user_management.adam.name
}

output "zia_user_management_adam_groups" {
  value = data.zia_user_management.adam.groups
}

output "zia_user_management_adam_department" {
  value = data.zia_user_management.adam.department
}
*/