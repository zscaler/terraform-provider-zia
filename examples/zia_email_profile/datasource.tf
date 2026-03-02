data "zia_email_profile" "this" {
    name = "EmailProfile01"
}

output "zia_email_profile" {
    value = data.zia_email_profile.this
}
