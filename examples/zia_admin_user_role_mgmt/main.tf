terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_admin_user_role_mgmt" "example1"{
    login_name = "admin@24326813.zscalerthree.net"
}

output "zia_admin_user_role_mgmt_example1"{
    value = data.zia_admin_user_role_mgmt.example1
}

data "zia_admin_user_role_mgmt" "example2"{
    login_name = "wguilherme@securitygeek.io"
}

output "zia_admin_user_role_mgmt_example2"{
    value = data.zia_admin_user_role_mgmt.example2
}