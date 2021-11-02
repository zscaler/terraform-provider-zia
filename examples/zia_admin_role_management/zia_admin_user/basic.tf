terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}

provider "zia" {}


resource "zia_admin_users" "example" {
  login_name = "john.smith@securitygeek.io"
  user_name  = "John Smith Test"
  email      = "john.smith@securitygeek.io"
  password   = "AeQ9E5w8B$"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  // admin_scope {
  //   scope_entities {
  //   id = [ data.zia_department_management.engineering.id,
  //          data.zia_department_management.executive.id 
  //   ]
  // }
  //   type = "DEPARTMENT"
  // }
}

output "zia_admin_users_example" {
  value = zia_admin_users.example
}


data "zia_admin_roles" "super_admin"{
    name = "Super Admin"
}

/*
data "zia_department_management" "engineering" {
 name = "Engineering"
}

data "zia_department_management" "executive" {
 name = "Executive"
}
*/

data "zia_admin_users" "john_smith" {
  login_name = "john.smith@securitygeek.io"
}

output "zia_admin_users_john_smith"{
    value = data.zia_admin_users.john_smith
}
