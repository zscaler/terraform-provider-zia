/*
terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}

provider "zia" {}

data "zia_admin_users" "john_ashcroft" {
  login_name = "john.smith@bd-hashicorp.com"
}

output "zia_admin_users_john_ashcroft"{
    value = data.zia_admin_users.john_ashcroft
}
*/