terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}

provider "zia" {}


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
  }

  }

output "zia_admin_users_example1" {
  value = zia_admin_users.example1
}


data "zia_admin_roles" "example1"{
    name = "Super Admin"
}

/*
data "zia_admin_users" "example1" {
  login_name = zia_admin_users.example1.login_name
}
*/


/*
output "zia_admin_roles_example1"{
    value = data.zia_admin_roles.example1
}
*/