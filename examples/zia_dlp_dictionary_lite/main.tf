terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_dlp_dictionaries_lite" "example"{
    name = "SALESFORCE_REPORT_LEAKAGE"
}

output "zia_dlp_dictionaries_lite_example"{
    value = data.zia_dlp_dictionaries_lite.example
}

