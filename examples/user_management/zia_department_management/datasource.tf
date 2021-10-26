data "zia_department_management" "engineering" {
 name = "Engineering"
}

output "zia_department_management" {
  value = data.zia_department_management.engineering
}

data "zia_department_management" "finance" {
 name = "Finance"
}

output "zia_department_management" {
  value = data.zia_department_management.finance
}