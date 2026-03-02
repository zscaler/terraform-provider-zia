resource "zia_email_profile" "this" {
    name        = "EmailProfile01"
    description = "Email recipient profile for DLP notifications"
    emails      = ["admin@example.com", "security@example.com"]
}
