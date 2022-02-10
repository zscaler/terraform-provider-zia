terraform {
    required_providers {
        zia = {
            version = "1.0.4"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_dlp_notification_templates" "example"{
    name = "DLP Auditor Template"
}

output "zia_dlp_notification_templates"{
    value = data.zia_dlp_notification_templates.example
}