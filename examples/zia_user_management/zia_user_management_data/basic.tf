######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_user_management" "john_ashcroft" {
 name = "John Ashcroft"
 email = "john.ashcroft@acme.com"
 password = "<YOURPASSWORDHERE>"
 auth_methods = ["BASIC"]
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