resource "zia_admin_users" "example" {
  login_name = "john.smith@bd-hashicorp.com"
  user_name  = "John Smith Test"
  email      = "john.smith@bd-hashicorp.com"
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

data "zia_admin_users" "john_smith" {
  login_name = "john.smith@bd-hashicorp.com"
}

output "zia_admin_users_john_smith"{
    value = data.zia_admin_users.john_smith
}
