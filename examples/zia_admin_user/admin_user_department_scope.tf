resource "zia_admin_users" "john_smith" {
  login_name                      = "john.smith@acme.com"
  username                       = "John Smith"
  email                           = "john.smith@acme.com"
  is_password_login_allowed       = true
  password                        = "AeQ9E5w8B$"
  is_security_report_comm_enabled = true
  is_service_update_comm_enabled  = true
  is_product_update_comm_enabled  = true
  comments                        = "Administrator Group"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  admin_scope {
    type = "DEPARTMENT"
    scope_entities {
      id = [data.zia_department_management.engineering.id]
    }
  }
}

data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}

data "zia_department_management" "engineering" {
  name = "Engineering"
}
