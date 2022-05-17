data "zia_admin_roles" "example1"{
    name = "Super Admin"
}

output "zia_admin_roles_example1"{
    value = data.zia_admin_roles.example1
}

data "zia_admin_roles" "example2"{
    name = "DevOps_Role"
}

output "zia_admin_roles_example2"{
    value = data.zia_admin_roles.example2
}