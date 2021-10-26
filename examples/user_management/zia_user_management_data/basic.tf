resource "zia_user_management" "adam_ashcroft" {
 name = "Adam Ashcroft"
 email = "adam.ashcroft@acme.com"
 password = "P@ssw0rd123*"
 groups {
   id = data.zia_group_management.finance.id
   name = data.zia_group_management.normal_internet.name
  }
 department {
   id = data.zia_department_management.engineering.id
   name = data.zia_department_management.engineering.name
  }

}