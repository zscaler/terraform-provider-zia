resource "zia_user_management" "john_ashcroft" {
 name = "John Ashcroft"
 email = "john.ashcroft@acme.com"
 password = "P@ssw0rd123*"
 groups {
  id = data.zia_group_management.normal_internet.id
  }
 department {
  id = data.zia_department_management.engineering.id
  }

}

data "zia_group_management" "normal_internet" {
 name = "Normal_Internet"
}

data "zia_department_management" "engineering" {
 name = "Engineering"
}