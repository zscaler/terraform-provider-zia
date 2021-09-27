terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_user_management" "adam" {
 name = "Adam Ashcroft"
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
/*
data "zia_user_management" "albert" {
 id = "29309058"
}


data "zia_user_management" "adam" {
 id = "29309057"
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

output "zia_user_management_albert_name" {
  value = data.zia_user_management.albert.name
}

output "zia_user_management_albert_groups" {
  value = data.zia_user_management.albert.groups
}

output "zia_user_management_albert_department" {
  value = data.zia_user_management.albert.department
}
*/