data "zia_group_management" "devops" {
 name = "DevOps"
}

output "zia_group_management" {
  value = data.zia_group_management.devops
}