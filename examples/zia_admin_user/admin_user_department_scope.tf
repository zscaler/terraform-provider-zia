######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_admin_users" "john_smith" {
  login_name                      = "john.smith@acme.com"
  username                       = "John Smith"
  email                           = "john.smith@acme.com"
  is_password_login_allowed       = true
  password                        = "********************"
  is_security_report_comm_enabled = true
  is_service_update_comm_enabled  = true
  is_product_update_comm_enabled  = true
  comments                        = "Administrator Group"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  admin_scope_type = "DEPARTMENT"
    admin_scope_entities {
        id = [ data.zia_department_management.engineering.id, data.zia_department_management.sales.id ]
    }
}

data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}

data "zia_department_management" "engineering" {
  name = "Engineering"
}

data "zia_department_management" "sales" {
  name = "Sales"
}

