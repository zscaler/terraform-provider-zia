data "zia_group_management" "devops" {
 name = "DevOps"
}

output "zia_group_management_devops" {
  value = data.zia_group_management.devops
}