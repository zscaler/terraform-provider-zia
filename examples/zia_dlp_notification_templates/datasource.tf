data "zia_dlp_notification_templates" "example"{
    name = "DLP Auditor Template"
}

output "zia_dlp_notification_templates"{
    value = data.zia_dlp_notification_templates.example
}