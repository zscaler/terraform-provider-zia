data "zia_department_management" "engineering" {
 name = "Engineering"
}

output "zia_department_management_engineering" {
  value = data.zia_department_management.engineering
}

data "zia_department_management" "finance" {
 name = "Finance"
}

output "zia_department_management_finance" {
  value = data.zia_department_management.finance
}