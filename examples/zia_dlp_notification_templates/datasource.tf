terraform {
    required_providers {
        zia = {
            version = "2.0.1"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_dlp_notification_templates" "example"{
    name = "Terraform DLP Template"
}

output "zia_dlp_notification_templates"{
    value = data.zia_dlp_notification_templates.example
}