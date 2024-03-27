######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_admin_users" "john_smith" {
  login_name                      = "john.smith@acme.com"
  username                       = "John Smith"
  email                           = "john.smith@acme.com"
  is_password_login_allowed       = true
  password                        = "<YOURPASSWORDHERE>"
  is_security_report_comm_enabled = true
  is_service_update_comm_enabled  = true
  is_product_update_comm_enabled  = true
  comments                        = "Administrator Group"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  admin_scope_type = "LOCATION_GROUP"
    admin_scope_entities {
        id = [ data.zia_location_groups.corporate_user_traffic_group.id ]
    }
}

data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}

data "zia_location_groups" "corporate_user_traffic_group" {
  name = "Corporate User Traffic Group"
}