terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_dlp_dictionaries" "example1"{
    name = "SALESFORCE_REPORT_LEAKAGE"
}

output "zia_dlp_dictionaries_example1"{
    value = data.zia_dlp_dictionaries.example1
}

data "zia_dlp_dictionaries" "example2"{
    name = "TIN_LEAKAGE"
}

output "zia_dlp_dictionaries_example2"{
    value = data.zia_dlp_dictionaries.example2
}

data "zia_dlp_dictionaries" "example3"{
    name = "PESEL_LEAKAGE"
}

output "zia_dlp_dictionaries_example3"{
    value = data.zia_dlp_dictionaries.example3
}

data "zia_dlp_dictionaries" "example4"{
    name = "AHV_LEAKAGE"
}

output "zia_dlp_dictionaries_example4"{
    value = data.zia_dlp_dictionaries.example4
}