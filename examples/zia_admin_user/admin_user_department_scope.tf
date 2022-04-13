terraform {
  required_providers {
    zia = {
      version = "2.0.1"
      source  = "zscaler.com/zia/zia"
    }
  }
}

provider "zia" {}

resource "zia_admin_users" "john_smith" {
  login_name                      = "john.smith@securitygeek.io"
  username                        = "John Smith"
  email                           = "john.smith@securitygeek.io"
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
    scope_group_member_entities {
      id = []
    }
  }
}

data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}

data "zia_department_management" "engineering" {
  name = "Engineering"
}
