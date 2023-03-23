resource "zia_user_management" "john_ashcroft" {
 name = "John Ashcroft"
 email = "john.ashcroft@acme.com"
 password = "P@ssw0rd123*"
 auth_methods = ["BASIC", "DIGEST"]
 groups {
  id = [ data.zia_group_management.normal_internet.id,
         data.zia_group_management.devops.id ]
  }
 department {
  id = data.zia_department_management.engineering.id
  }

}

data "zia_group_management" "normal_internet" {
 name = "Normal_Internet"
}

data "zia_group_management" "devops" {
 name = "DevOps"
}

data "zia_department_management" "engineering" {
 name = "Engineering"
}