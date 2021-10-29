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