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
  is_exec_mobile_app_enabled      = false
  comments                        = "Administrator Group"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  admin_scope_type = "ORGANIZATION"
}

data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}