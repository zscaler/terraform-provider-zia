terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}

provider "zia" {}


data "zia_admin_roles" "example1"{
    name = "Super Admin"
}

resource "zia_admin_users" "example1" {
  login_name = "zia-api2@securitygeek.io"
  user_name  = "John Smith Test"
  email      = "zia-api2@securitygeek.io"
  password   = "AeQ9E5w8B$"
  role {
    id = data.zia_admin_roles.example1.id
  }
  admin_scope {
    type = "ORGANIZATION"
    scope_entities {
      id = []
    }
    scope_group_member_entities {
      id = [25658881]
    }
  }
}

output "zia_admin_users_example1" {
  value = zia_admin_users.example1.user_name
}


data "zia_admin_users" "example" {
  login_name = "amazzal.elhabib@securitygeek.io"
}

output "zia_admin_users_example"{
    value = data.zia_admin_users.example
}