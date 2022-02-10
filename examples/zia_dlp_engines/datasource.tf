data "zia_dlp_engines" "example"{
    name = "Credit Cards"
}

output "zia_dlp_engines"{
    value = data.zia_dlp_engines.example
}