terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_dlp_dictionaries" "example"{
    name = "SGIO-EDM-Test"
}

output "zia_dlp_dictionaries"{
    value = data.zia_dlp_dictionaries.example
}