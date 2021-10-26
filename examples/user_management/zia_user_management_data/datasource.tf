data "zia_user_management" "adam_ashcroft" {
 name = "Adam Ashcroft"
}

output "zia_user_management" {
  value = data.zia_user_management.adam_ashcroft
}