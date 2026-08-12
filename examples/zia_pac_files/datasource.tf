# Retrieve all PAC files
data "zia_pac_files" "all" {}

# Retrieve a PAC file by name
data "zia_pac_files" "this" {
  name = "example_pac_file"
}

output "pac_file_url" {
  value = data.zia_pac_files.this.pac_files[0].pac_url
}
